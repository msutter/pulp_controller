package main

import (
    "time"
    "labix.org/v2/mgo/bson"
)

type Server struct  {
    Id          bson.ObjectId  `json:"id" bson:"_id,omitempty"`
    Name        string         `json:"name"`
    Url         string         `json:"url"`
    Username    string         `json:"username"`
    Password    string         `json:"password"`
    Added       time.Time      `json:"added"`
}

type Servers []Server