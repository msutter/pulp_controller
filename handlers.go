package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func ServerIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    servers := Servers{
        Server{Id: 1, Name: "test01", Url: "http://test.domain.com", Username: "testuser", Password: "highsecure", Added: time.Now() },
    }

    if err := json.NewEncoder(w).Encode(servers); err != nil {
        panic(err)
    }
}

func ServerShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serverId := vars["serverId"]
    fmt.Fprintln(w, "Server show:", serverId)
}