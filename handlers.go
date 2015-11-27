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

    session, collection := InitServerCollection()
    defer session.Close()
    servers := Servers{}

    SearchAll(bson.M{}, collection, &servers)

    if err := json.NewEncoder(w).Encode(servers); err != nil {
        panic(err)
    }
}

func ServerShow(w http.ResponseWriter, r *http.Request) {
    oid := getOID("serverId", r)

    session, collection := InitServerCollection()
    defer session.Close()
    server := Server{}

    SearchOne(bson.M{"_id": oid}, collection, &server)

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
    session, collection := InitServerCollection()
    defer session.Close()
    server.Added = time.Now()
    err = collection.Insert(server)
    if err != nil {
        log.Println(err)
    }
    w.Header().Set("Content-Type","application/json; charset=UTF=8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(server); err != nil {
        panic(err)
    }
}

func ServerDelete(w http.ResponseWriter, r *http.Request) {
    oid := getOID("serverId", r)
    session, collection := InitServerCollection()
    defer session.Close()
    err := collection.Remove(bson.M{"_id": oid})
    if err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusInternalServerError)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
        }
    }
}

func getOID(name string, r *http.Request) bson.ObjectId {
    vars := mux.Vars(r)
    id := vars[name]
    return bson.ObjectIdHex(id)
}