package pulp

import (
    "net/http"
    "github.com/ampersand8/pulp_controller/logger"
    "io/ioutil"
    "io"
    "encoding/json"
)

func ListRepos(server Server) Repos {
    var repos Repos
    r, err := http.Get(server.Url)
    if err != nil {
        logger.Log(err.Error(), logger.WARN)
    }

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        logger.Log(err.Error(), logger.ERROR)
    }

    if err := r.Body.Close(); err != nil {
        logger.Log(err.Error(), logger.ERROR)
    }

    if err := json.Unmarshal(body, &repos); err != nil {
        return nil
        logger.Log(err.Error(), logger.WARN)
    }

    return repos
}