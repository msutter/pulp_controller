package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func SetRoutes() Routes {
    var routes = Routes{
        Route{
            "Index",
            "GET",
            settings.ApiVersion,
            Index,
        },
        Route{
            "ServerIndex",
            "GET",
            settings.ApiVersion + "/servers",
            ServerIndex,
        },
        Route{
            "ServerShow",
            "GET",
            settings.ApiVersion + "/servers/{serverId}",
            ServerShow,
        },
        Route{
            "ServerCreate",
            "POST",
            settings.ApiVersion + "/servers",
            ServerCreate,
        },
        Route{
            "ServerDelete",
            "DELETE",
            settings.ApiVersion + "/servers/{serverId}",
            ServerDelete,
        },
        Route{
            "Authenticate",
            "POST",
            settings.ApiVersion + "/authenticate",
            AuthHandler,
        },
        Route{
            "Restricted",
            "GET",
            settings.ApiVersion + "/restricted",
            RestrictedHandler,
        },
    }
    return routes
}
