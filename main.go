package main

import (
    "log"
    "net/http"
    "github.com/ampersand8/pulp_controller/config"
    "github.com/ampersand8/pulp_controller/logger"
)

func main() {
    logger.SetLogLevel(config.Settings.Loglevel)
    router := NewRouter()
    log.Fatal(http.ListenAndServe(config.Settings.Server+":"+config.Settings.Port, router))
}