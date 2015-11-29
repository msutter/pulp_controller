package main

import (
    "log"
    "net/http"
    "github.com/ampersand8/pulp_controller/config"
)

func main() {
    router := NewRouter()
    log.Fatal(http.ListenAndServe(config.Settings.Server+":"+config.Settings.Port, router))
}