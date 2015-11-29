package main

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/ampersand8/pulp_controller/config"
    "github.com/ampersand8/pulp_controller/logger"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    logger.Log("Router API Version: "+config.Settings.ApiVersion, logger.INFO)
    routes := SetRoutes()
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = logger.Logger(handler, route.Name)
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
    return router
}