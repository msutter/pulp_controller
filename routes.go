package main

import "net/http"

const api_version string = "/v1"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "Index",
        "GET",
        api_version,
        Index,
    },
    Route{
        "ServerIndex",
        "GET",
        api_version + "/servers",
        ServerIndex,
    },
    Route{
        "ServerShow",
        "GET",
        api_version + "/servers/{serverId}",
        ServerShow,
    },
    Route{
        "ServerCreate",
        "POST",
        api_version + "/servers",
        ServerCreate,
    },
    Route{
        "ServerDelete",
        "DELETE",
        api_version + "/servers/{serverId}",
        ServerDelete,
    },
}
