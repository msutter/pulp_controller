package logger

import (
    "log"
    "net/http"
    "time"
)

// The smaller the loglevel the less you see
const (
    ERROR int = 10
    WARN  int = 20
    INFO  int = 30
    DEBUG int = 40
)

var LOGLEVEL int = 30

func SetLogLevel(loglevel int) {
    LOGLEVEL = loglevel
}

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        inner.ServeHTTP(w, r)
        Log(r.Method + ", " + r.RequestURI + ", " + name + ", " + time.Since(start).String(), INFO)
    })
}

func Log(msg string, level int) {
    if LOGLEVEL >= level {
        timestamp := time.Now()
        var namelevel string
        switch level {
        case ERROR: namelevel = "ERROR"
        case WARN: namelevel = "WARN"
        case INFO: namelevel = "INFO"
        case DEBUG: namelevel = "DEBUG"
        }
        log.Printf("%s - %s: %s\n", namelevel, timestamp, msg)
    }
}