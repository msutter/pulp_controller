package main

import (
    jwt "github.com/dgrijalva/jwt-go"
    "crypto/rsa"
    "io/ioutil"
    "log"
    "time"
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