package main

import (
    jwt "github.com/dgrijalva/jwt-go"
    "crypto/rsa"
    "io/ioutil"
    "log"
    "time"
    "net/http"
    "encoding/json"
)


const (
    privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa 2048
    pubKeyPath = "keys/app.rsa.pub"  // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
    verifyKey *rsa.PublicKey
    signKey *rsa.PrivateKey
)

func initKeys() {
    signBytes, err := ioutil.ReadFile(privKeyPath)
    fatal(err)

    signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
    fatal(err)

    verifyBytes, err := ioutil.ReadFile(pubKeyPath)
    fatal(err)

    verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
    fatal(err)
}

func CreateTokenString(user string) (string, error) {
    // create a signer for rsa256
    t := jwt.New(jwt.GetSigningMethod("RS256"))

    // set our claims
    t.Claims["AccesToken"] = "level1"
    t.Claims["CustomUserInfo"] = struct {
        Name string
        Kind string
    }{user, "human"}

    //DEBUG
    log.Printf("Created claims")
    // set expire time
    t.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
    return t.SignedString(signKey)
}

func IsAllowed(w http.ResponseWriter, r *http.Request) bool {
    client_token := r.Header.Get("admintoken")
    if client_token == "" {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusForbidden)
        json.NewEncoder(w).Encode(map[string]string{"error": "token not set"})
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
            return false
        }
        return true

    case *jwt.ValidationError: // something was wrong during the validation
        vErr := err.(*jwt.ValidationError)

        switch vErr.Errors {
        case jwt.ValidationErrorExpired:
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "token expired"})
            return false

        default:
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "error while parsing token "})
            return false
        }

    default: // something else went wrong
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string {"error": "something went wrong"})
        return false
    }
}

// Dummy Authentication function
// yet to be implemented
func Authenticate(user string, pass string) bool {
    if user == "admin" && pass == "mypassword" {
        return true
    } else {
        return false
    }
}

func fatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}