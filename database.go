package main

import (
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
)

func Connect() *mgo.Session{
    url := "mongodb://testuser:testpassword@localhost/pulp_manager_api_test"
    session, err := mgo.Dial(url)
    if err != nil {
        Log("Connection to mongodb failed, Error: " + err.Error(), ERROR)
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
        Log("MongoDB search failed, Error: " + err.Error(), ERROR)
    }
}

func SearchAll(search bson.M, collection *mgo.Collection, result interface{}) {
    err := collection.Find(search).All(result)
    if err != nil {
        Log("MongoDB search failed, Error: " + err.Error(), ERROR)
    }
}