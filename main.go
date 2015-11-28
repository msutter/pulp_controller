package main

import (
    "log"
    "net/http"
)

func main() {
    Log("Starting APP", INFO)
    initKeys()
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8080", router))
}