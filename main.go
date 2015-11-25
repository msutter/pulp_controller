package main

import (
    "log"
    "net/http"
)

func main() {
    Connect()
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8080", router))
}