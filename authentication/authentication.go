// Everything about authentication and authorization
package authentication

import (
    jwt "github.com/dgrijalva/jwt-go"
    "github.com/ampersand8/pulp_controller/config"
    "github.com/ampersand8/pulp_controller/logger"
    "crypto/rsa"
    "io/ioutil"
    "time"
    "net/http"
    "encoding/json"
)

var (
    verifyKey *rsa.PublicKey
    signKey *rsa.PrivateKey
)

func init() {
    signBytes, err := ioutil.ReadFile(config.Settings.PrivKeyPath)
    if err != nil {
        logger.Log("Reading private Key File " + config.Settings.PrivKeyPath + " failed, Error: " + err.Error(), logger.ERROR)
    }

    signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
    if err != nil {
        logger.Log("Parsing private Key File failed, Error: " + err.Error(), logger.ERROR)
    }

    verifyBytes, err := ioutil.ReadFile(config.Settings.PubKeyPath)
    if err != nil {
        logger.Log("Reading public Key File " + config.Settings.PubKeyPath + " failed, Error: " + err.Error(), logger.ERROR)
    }

    verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
    if err != nil {
        logger.Log("Parsing public Key File failed, Error: " + err.Error(), logger.ERROR)
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

    logger.Log("Created claims", logger.DEBUG)
    // set expire time
    t.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Settings.Tokenexpiration)).Unix()
    return t.SignedString(signKey)
}

func IsAllowed(w http.ResponseWriter, r *http.Request) bool {
    logger.Log("Checking whether admintoken is set", logger.DEBUG)
    client_token := r.Header.Get("admintoken")
    if client_token == "" {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusForbidden)
        json.NewEncoder(w).Encode(map[string]string{"error": "token not set"})
        logger.Log("Token 'admintoken' not set", logger.WARN)
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
            logger.Log("Token is invalid", logger.WARN)
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
            logger.Log("Token expired", logger.WARN)
            return false

        default:
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string {"error": "error while parsing token "})
            logger.Log("Error while parsing token: " + token.Raw, logger.WARN)
            return false
        }

    default: // something else went wrong
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string {"error": "something went wrong"})
        logger.Log("Something with this token is wrong: " + token.Raw, logger.ERROR)
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
        logger.Log(err.Error(), logger.ERROR)
    }
}
