package main

import (
    "log"
    "net/http"
)

func main() {
    readConfiguration()
    Log("My Config: "+settings.ApiVersion, INFO)
    Log("Starting APP", INFO)
    initKeys()
    router := NewRouter()
    log.Fatal(http.ListenAndServe(settings.Server+":"+settings.Port, router))
}