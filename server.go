package main

import "time"

type Server struct  {
    Id          int         `json:"id"`
    Name        string      `json:"name"`
    Url         string      `json:"url"`
    Username    string      `json:"username"`
    Password    string      `json:"password"`
    Added       time.Time   `json:"added"`
}

type Servers []Server