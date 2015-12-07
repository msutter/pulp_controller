package pulp

import (
    "net/http"
    "github.com/ampersand8/pulp_controller/logger"
    "io/ioutil"
    "io"
    "encoding/json"
    "encoding/base64"
)

func Login(server Server, r *http.Request) {
    token := string(base64.StdEncoding.EncodeToString([]byte(server.Username + ":" + server.Password)))
    r.Header.Add("Authorization", "Basic " + token)
}

func ListRepos(server Server) Repos {
    var repos Repos
    client := http.Client{}
    r, err := http.NewRequest("GET", server.Url+"/repositories/", nil)
    Login(server, r)
    resp, err := client.Do(r)
    if err != nil {
        logger.Log(err.Error(), logger.WARN)
    }

    body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
    if err != nil {
        logger.Log(err.Error(), logger.ERROR)
    }

    if err := resp.Body.Close(); err != nil {
        logger.Log(err.Error(), logger.ERROR)
    }

    if err := json.Unmarshal(body, &repos); err != nil {
        return nil
        logger.Log(err.Error(), logger.WARN)
    }

    return repos
}