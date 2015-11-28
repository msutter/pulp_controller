package main

import (
    jwt "github.com/dgrijalva/jwt-go"
    "crypto/rsa"
    "io/ioutil"
    "time"
    "net/http"
    "encoding/json"
)


const (
    privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa 2048
    pubKeyPath = "keys/app.rsa.pub"  // openssl rsa -in app.rsa -pubout > app.rsa.pub
    TOKENEXPIRATION = 1          // Token expiration in hours
)

var (
    verifyKey *rsa.PublicKey
    signKey *rsa.PrivateKey
)

func initKeys() {
    signBytes, err := ioutil.ReadFile(privKeyPath)
    if err != nil {
        Log("Reading private Key File " + privKeyPath + " failed, Error: " + err.Error(), ERROR)
    }

    signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
    if err != nil {
        Log("Parsing private Key File failed, Error: " + err.Error(), ERROR)
    }

    verifyBytes, err := ioutil.ReadFile(pubKeyPath)
    if err != nil {
        Log("Reading public Key File " + pubKeyPath + " failed, Error: " + err.Error(), ERROR)
    }

    verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
    if err != nil {
        Log("Parsing public Key File failed, Error: " + err.Error(), ERROR)
    }
}

func CreateTokenString(user string) (string, error) {
    // create a signer for rsa256
    t := jwt.New(jwt.GetSigningMethod("RS256"))

    // set our claims
    t.Claims["AccesToken"] = "level1"
    t.Claims["CustomUserInfo"] = struct {
        Name string
        Role string
    }{user, "admin"}

    Log("Created claims", DEBUG)
    // set expire time
    t.Claims["exp"] = time.Now().Add(time.Hour * TOKENEXPIRATION).Unix()
    return t.SignedString(signKey)
}

func IsAllowed(w http.ResponseWriter, r *http.Request) bool {
    Log("Checking whether admintoken is set", DEBUG)
    client_token := r.Header.Get("admintoken")
    if client_token == "" {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusForbidden)
        json.NewEncoder(w).Encode(map[string]string{"error": "token not set"})
        Log("Token 'admintoken' not set", WARN)
        return false
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
            Log("Token is invalid", WARN)
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
            Log("Token expired", WARN)
            return false

        default:
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "error while parsing token "})
            Log("Error while parsing token: " + token.Raw, WARN)
            return false
        }

    default: // something else went wrong
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string {"error": "something went wrong"})
        Log("Something with this token is wrong: " + token.Raw, ERROR)
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
        Log(err.Error(), ERROR)
    }
}