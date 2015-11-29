package main

import (
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "github.com/ampersand8/pulp_controller/config"
    "github.com/ampersand8/pulp_controller/logger"
)

func Connect() *mgo.Session{
    url := "mongodb://"+config.Settings.DbUser+":"+config.Settings.DbPassword+"@"+config.Settings.DbHost+":"+config.Settings.DbPort+"/"+config.Settings.DbName
    logger.Log("Connecting to "+url, logger.INFO)
    session, err := mgo.Dial(url)
    if err != nil {
        logger.Log("Connection to mongodb failed, Error: " + err.Error(), logger.ERROR)
        panic(err)
    }
    return session
}

func InitServerCollection() (*mgo.Session, *mgo.Collection) {
    session := Connect()
    collection := session.DB(config.Settings.DbName).C("servers")
    return session, collection
}

func SearchOne(search bson.M, collection *mgo.Collection, result interface{}) {
    err := collection.Find(search).One(result)
    if err != nil {
        logger.Log("MongoDB search failed, Error: " + err.Error(), logger.ERROR)
    }
}

func SearchAll(search bson.M, collection *mgo.Collection, result interface{}) {
    err := collection.Find(search).All(result)
    if err != nil {
        logger.Log("MongoDB search failed, Error: " + err.Error(), logger.ERROR)
    }
}