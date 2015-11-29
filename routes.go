package main

import (
    "net/http"
    "github.com/ampersand8/pulp_controller/config"
)

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
            config.Settings.ApiVersion,
            Index,
        },
        Route{
            "ServerIndex",
            "GET",
            config.Settings.ApiVersion + "/servers",
            ServerIndex,
        },
        Route{
            "ServerShow",
            "GET",
            config.Settings.ApiVersion + "/servers/{serverId}",
            ServerShow,
        },
        Route{
            "ServerCreate",
            "POST",
            config.Settings.ApiVersion + "/servers",
            ServerCreate,
        },
        Route{
            "ServerDelete",
            "DELETE",
            config.Settings.ApiVersion + "/servers/{serverId}",
            ServerDelete,
        },
        Route{
            "Authenticate",
            "POST",
            config.Settings.ApiVersion + "/authenticate",
            AuthHandler,
        },
        Route{
            "Restricted",
            "GET",
            config.Settings.ApiVersion + "/restricted",
            RestrictedHandler,
        },
    }
    return routes
}
