package main

import (
    "log"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
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

func InitServerCollection() (*mgo.Session, *mgo.Collection) {
    session := Connect()
    collection := session.DB("pulp_manager_api_test").C("servers")
    return session, collection
}

func SearchOne(search bson.M, collection *mgo.Collection, result interface{}) {
    err := collection.Find(search).One(result)
    if err != nil {
        log.Println(err)
    }
}

func SearchAll(search bson.M, collection *mgo.Collection, result interface{}) {
    err := collection.Find(search).All(result)
    if err != nil {
        log.Println(err)
    }
}