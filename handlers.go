package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "labix.org/v2/mgo/bson"
    "github.com/gorilla/mux"
    "io"
    "io/ioutil"
    "github.com/ampersand8/pulp_controller/logger"
    "github.com/ampersand8/pulp_controller/authentication"
    "github.com/ampersand8/pulp_controller/pulp"
    db "github.com/ampersand8/pulp_controller/database"
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

    session, collection := db.InitServerCollection()
    defer session.Close()
    servers := pulp.Servers{}

    db.SearchAll(bson.M{}, collection, &servers)

    if err := json.NewEncoder(w).Encode(servers); err != nil {
        logger.Log("could not encode servers struct, Error: " + err.Error(), logger.ERROR)
    }
}

func ServerShow(w http.ResponseWriter, r *http.Request) {
    oid := getOID("serverId", r)

    session, collection := db.InitServerCollection()
    defer session.Close()
    server := pulp.Server{}

    db.SearchOne(bson.M{"_id": oid}, collection, &server)

    if err := json.NewEncoder(w).Encode(server); err != nil {
        logger.Log("could not encode server struct, Error: " + err.Error(), logger.ERROR)
    }
}

func ServerCreate(w http.ResponseWriter, r *http.Request) {
    if authentication.IsAllowed(w, r) {
        var server pulp.Server
        body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
        if err != nil {
            logger.Log("could not read POST body, Error: " + err.Error(), logger.ERROR)
        }

        if err := r.Body.Close(); err != nil {
            logger.Log("could not close POST body, Error: " + err.Error(), logger.ERROR)
        }

        if err := json.Unmarshal(body, &server); err != nil {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(422)  // unprocessable entity
            if err := json.NewEncoder(w).Encode(err); err != nil {
                logger.Log("could not json/encode error, Error: " + err.Error(), logger.ERROR)
            }
        }
        session, collection := db.InitServerCollection()
        defer session.Close()
        server.Added = time.Now()
        err = collection.Insert(server)
        if err != nil {
            logger.Log("could not insert server to DB, Error: " + err.Error(), logger.ERROR)
        }
        w.Header().Set("Content-Type", "application/json; charset=UTF=8")
        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(server); err != nil {
            panic(err)
        }
    }
}

func ServerDelete(w http.ResponseWriter, r *http.Request) {
    if authentication.IsAllowed(w, r) {
        oid := getOID("serverId", r)
        session, collection := db.InitServerCollection()
        defer session.Close()
        err := collection.Remove(bson.M{"_id": oid})
        if err != nil {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusInternalServerError)
            if err := json.NewEncoder(w).Encode(err); err != nil {
                logger.Log("could not json/encode error, Error: " + err.Error(), logger.ERROR)
            }
        }
    }
}

func RepoList(w http.ResponseWriter, r *http.Request) {
    oid := getOID("serverId", r)

    session, collection := db.InitServerCollection()
    defer session.Close()
    server := pulp.Server{}

    db.SearchOne(bson.M{"_id": oid}, collection, &server)

    repos := pulp.ListRepos(server)

    if err := json.NewEncoder(w).Encode(repos); err != nil {
        logger.Log("could not encode repos struct, Error: " + err.Error(), logger.ERROR)
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

    logger.Log("Authenticate: user["+ user.Username + "] pass[" + user.Password + "]", logger.INFO)

    if authentication.Authenticate(user.Username, user.Password) {
        logger.Log("Authenticating user[" + user.Username + "] successful", logger.INFO)
    } else {
        logger.Log("Authenticating user[" + user.Username + "] not successful", logger.WARN)
        w.WriteHeader(http.StatusForbidden)
        return
    }

    tokenString, err := authentication.CreateTokenString(user.Username)
    logger.Log("Created token string", logger.DEBUG)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "Sorry error while Signing Key")
        logger.Log("Token Signing error: %s" + err.Error(), logger.ERROR)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    logger.Log("Sending back token: " + tokenString, logger.DEBUG)
    if err := json.NewEncoder(w).Encode(map[string]string{"token":tokenString}); err != nil {
        logger.Log("Sending token unsuccessful, Error: " + err.Error(), logger.ERROR)
    }
}

// Example restricted handler
func RestrictedHandler(w http.ResponseWriter, r *http.Request) {
    // if user has no or invalid token error message is generated by Authenticate
    // if token is valid the content below is shown
    if authentication.IsAllowed(w, r) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string {"secret": "you made it to the secret area"})
    }
}

func getOID(name string, r *http.Request) bson.ObjectId {
    vars := mux.Vars(r)
    id := vars[name]
    return bson.ObjectIdHex(id)
}