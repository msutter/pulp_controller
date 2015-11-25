package main

import (
    "log"
    "labix.org/v2/mgo"
)

func Connect() *mgo.Session{
    url := "mongodb://testuser:testpassword@localhost/pulp_manager_api_test"
    session, err := mgo.Dial(url)
    if err != nil {
        log.Fatal(err)
        panic(err)
    }

    return session
}