package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "log"

    "labix.org/v2/mgo/bson"
    "github.com/gorilla/mux"
    "io"
    "io/ioutil"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func ServerIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    session := Connect()
    defer session.Close()
    c := session.DB("pulp_manager_api_test").C("servers")
    servers := Servers{}
    err := c.Find(bson.M{}).All(&servers)

    if err != nil {
        log.Fatal(err)
    }

    if err := json.NewEncoder(w).Encode(servers); err != nil {
        panic(err)
    }
}

func ServerShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serverId := vars["serverId"]
    oid := bson.ObjectIdHex(serverId)

    session := Connect()
    defer session.Close()
    c := session.DB("pulp_manager_api_test").C("servers")
    server := Server{}


    err := c.Find(bson.M{"_id": oid}).One(&server)

    if err != nil {
        log.Print(err)
    }

    if err := json.NewEncoder(w).Encode(server); err != nil {
        panic(err)
    }
}

func ServerCreate(w http.ResponseWriter, r *http.Request) {
    var server Server
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }

    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    if err := json.Unmarshal(body, &server); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)  // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    session := Connect()
    defer session.Close()
    server.Added = time.Now()
    c := session.DB("pulp_manager_api_test").C("servers")
    err = c.Insert(server)
    if err != nil {
        log.Fatal(err)
    }
    w.Header().Set("Content-Type","application/json; charset=UTF=8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(server); err != nil {
        panic(err)
    }
}