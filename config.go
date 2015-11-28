package main
import (
    "os"
    "encoding/json"
    "flag"
    "strconv"
)

var settings struct {
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

func readConfiguration() {
    // First load defaults
    loadDefaultConfig()

    // Second load config file
    configFile, err := os.Open("config.json")
    if err != nil {
        Log("Could not load configuration file - Loading defaults, Error: " + err.Error(), ERROR)
    } else {
        jsonParser := json.NewDecoder(configFile)
        if err = jsonParser.Decode(&settings); err != nil {
            Log("Parsing config file config.json failed, Error:" + err.Error(), ERROR)
        }
    }

    // Third load config from flags
    setFlags()

    Log("API Version: "+settings.ApiVersion, INFO)
}

func setFlags() {
    flag.IntVar(&settings.Loglevel, "loglevel", settings.Loglevel, "loglevel - "+strconv.Itoa(ERROR)+": ERROR, "+
    strconv.Itoa(WARN)+": WARN, "+strconv.Itoa(INFO)+": INFO, "+strconv.Itoa(DEBUG)+": DEBUG")
    flag.StringVar(&settings.Port, "port", settings.Port, "port")
    flag.IntVar(&settings.Tokenexpiration, "expiration", settings.Tokenexpiration, "token expiration in hours")
    flag.StringVar(&settings.PrivKeyPath, "privatekey", settings.PrivKeyPath, "private key file")
    flag.StringVar(&settings.PubKeyPath, "publickey", settings.PubKeyPath, "public key file")
    flag.StringVar(&settings.DbUser, "dbuser", settings.DbUser, "database username")
    flag.StringVar(&settings.DbPassword, "dbpass", settings.DbPassword, "database password")
    flag.StringVar(&settings.DbPort, "dbport", settings.DbPort, "database port")
    flag.StringVar(&settings.DbName, "db", settings.DbName, "database name")
    flag.StringVar(&settings.ApiVersion, "apiversion", settings.ApiVersion, "API version, e.g. /v1 or /v2")
    flag.Parse()
}

func loadDefaultConfig() {
    settings.Loglevel        = INFO
    settings.Port            = "8080"
    settings.Server          = "localhost"
    settings.Tokenexpiration = 1
    settings.PrivKeyPath     = "keys/app.rsa"
    settings.PubKeyPath      = "keys/app.rsa.pub"
    settings.DbUser          = "testuser"
    settings.DbPassword      = "testpassword"
    settings.DbPort          = "27017"
    settings.DbHost          = "localhost"
    settings.DbName          = "pulp_manager_api_test"
    settings.ApiVersion      = "/v1"
}