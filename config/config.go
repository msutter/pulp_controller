package config
import (
    "os"
    "encoding/json"
    "flag"
    "strconv"
    "github.com/ampersand8/pulp_controller/logger"
)

var Settings struct {
    Loglevel           int               `json: "loglevel"`
    Port               string            `json: "port"`
    Server             string            `json: "server"`
    Tokenexpiration    int               `json: "tokenexpiration"`
    PrivKeyPath        string            `json: "privkeypath"`  // openssl genrsa -out app.rsa 2048
    PubKeyPath         string            `json: "pubkeypath"`   // openssl rsa -in app.rsa -pubout > app.rsa.pub
    DbUser             string            `json: "dbuser"`
    DbPassword         string            `json: "dbpassword"`
    DbPort             string            `json: "dbport"`
    DbHost             string            `json: "dbhost"`
    DbName             string            `json: "dbname"`
    ApiVersion         string            `json: "apiversion"`
}

func init() {
    // First load defaults
    loadDefaultConfig()

    // Second load config file
    configFile, err := os.Open("config.json")
    if err != nil {
        logger.Log("Could not load configuration file - Loading defaults, Error: " + err.Error(), logger.ERROR)
    } else {
        jsonParser := json.NewDecoder(configFile)
        if err = jsonParser.Decode(&Settings); err != nil {
            logger.Log("Parsing config file config.json failed, Error:" + err.Error(), logger.ERROR)
        }
    }

    // Third load config from flags
    setFlags()

    logger.Log("API Version: "+Settings.ApiVersion, logger.INFO)
}

func setFlags() {
    flag.IntVar(&Settings.Loglevel, "loglevel", Settings.Loglevel, "loglevel - "+strconv.Itoa(logger.ERROR)+": ERROR, "+
    strconv.Itoa(logger.WARN)+": WARN, "+strconv.Itoa(logger.INFO)+": INFO, "+strconv.Itoa(logger.DEBUG)+": DEBUG")
    flag.StringVar(&Settings.Port, "port", Settings.Port, "port")
    flag.IntVar(&Settings.Tokenexpiration, "expiration", Settings.Tokenexpiration, "token expiration in hours")
    flag.StringVar(&Settings.PrivKeyPath, "privatekey", Settings.PrivKeyPath, "private key file")
    flag.StringVar(&Settings.PubKeyPath, "publickey", Settings.PubKeyPath, "public key file")
    flag.StringVar(&Settings.DbUser, "dbuser", Settings.DbUser, "database username")
    flag.StringVar(&Settings.DbPassword, "dbpass", Settings.DbPassword, "database password")
    flag.StringVar(&Settings.DbPort, "dbport", Settings.DbPort, "database port")
    flag.StringVar(&Settings.DbName, "db", Settings.DbName, "database name")
    flag.StringVar(&Settings.ApiVersion, "apiversion", Settings.ApiVersion, "API version, e.g. /v1 or /v2")
    flag.Parse()
}

func loadDefaultConfig() {
    Settings.Loglevel        = logger.INFO
    Settings.Port            = "8080"
    Settings.Server          = "localhost"
    Settings.Tokenexpiration = 1
    Settings.PrivKeyPath     = "keys/app.rsa"
    Settings.PubKeyPath      = "keys/app.rsa.pub"
    Settings.DbUser          = "testuser"
    Settings.DbPassword      = "testpassword"
    Settings.DbPort          = "27017"
    Settings.DbHost          = "localhost"
    Settings.DbName          = "pulp_manager_api_test"
    Settings.ApiVersion      = "/v1"
}