package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "log"

    "labix.org/v2/mgo/bson"
    "github.com/gorilla/mux"
    "io"
    "io/ioutil"
    "github.com/dgrijalva/jwt-go"
)

type AdminUser struct {
    Username string     `json: "username"`
    Password string     `json: "password"`
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func ServerIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    session, collection := InitServerCollection()
    defer session.Close()
    servers := Servers{}

    SearchAll(bson.M{}, collection, &servers)

    if err := json.NewEncoder(w).Encode(servers); err != nil {
        panic(err)
    }
}

func ServerShow(w http.ResponseWriter, r *http.Request) {
    oid := getOID("serverId", r)

    session, collection := InitServerCollection()
    defer session.Close()
    server := Server{}

    SearchOne(bson.M{"_id": oid}, collection, &server)

    if err := json.NewEncoder(w).Encode(server); err != nil {
        panic(err)
    }
}

func ServerCreate(w http.ResponseWriter, r *http.Request) {
    var server Server
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }

    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    if err := json.Unmarshal(body, &server); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)  // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    session, collection := InitServerCollection()
    defer session.Close()
    server.Added = time.Now()
    err = collection.Insert(server)
    if err != nil {
        log.Println(err)
    }
    w.Header().Set("Content-Type","application/json; charset=UTF=8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(server); err != nil {
        panic(err)
    }
}

func ServerDelete(w http.ResponseWriter, r *http.Request) {
    oid := getOID("serverId", r)
    session, collection := InitServerCollection()
    defer session.Close()
    err := collection.Remove(bson.M{"_id": oid})
    if err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusInternalServerError)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
        }
    }
}

// Handles Authentication requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
    // has to be POST
    if r.Method != "POST" {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var user AdminUser
    json.NewDecoder(r.Body).Decode(&user)

    log.Printf("Authenticate: user[%s] pass[%s]\n", user.Username, user.Password)

    if Authenticate(user.Username, user.Password) {
        log.Printf("Authenticate: user[%s] valid", user.Username)
    } else {
        log.Printf("Authenticate: user[%s] invalid", user.Username)
        w.WriteHeader(http.StatusForbidden)
        return
    }

    tokenString, err := CreateTokenString(user.Username)

    //DEBUG
    log.Println("Created token stirng")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "Sorry error while Signing Key")
        log.Printf("Token Signing error: %s", err)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    //DEBUG
    log.Println("About to send the OK")
    log.Printf("Sending back token: %s\n", tokenString)
    if err := json.NewEncoder(w).Encode(map[string]string{"token":tokenString}); err != nil {
        log.Println(err)
    }
}

// Example restricted handler
func RestrictedHandler(w http.ResponseWriter, r *http.Request) {
    client_token := r.Header.Get("admintoken")
    if client_token == "" {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusForbidden)
        json.NewEncoder(w).Encode(struct{error string}{"token not set"})
    }

    // validate the token
    token, err := jwt.Parse(client_token, func(token *jwt.Token) (interface{}, error) {
        return verifyKey, nil
    })

    switch err.(type) {
    case nil: //no error
        if !token.Valid { // may be invalid
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "token is invalid"})
            log.Println("Token is invalid")
            return
        }
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string {"secret": "you made it to the secret area"})

    case *jwt.ValidationError: // something was wrong during the validation
        vErr := err.(*jwt.ValidationError)

        switch vErr.Errors {
        case jwt.ValidationErrorExpired:
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "token expired"})
            return

        default:
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "error while parsing token "})
            return
        }

    default: // something else went wrong
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string {"error": "something went wrong"})
        return
    }


}

func getOID(name string, r *http.Request) bson.ObjectId {
    vars := mux.Vars(r)
    id := vars[name]
    return bson.ObjectIdHex(id)
}