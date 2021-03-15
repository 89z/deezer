package deezer

import (
   "net/http"
   "net/url"
)

const (
   gatewayAPI = "http://api.deezer.com/1.0/gateway.php"
   gatewayMedia = "https://media.deezer.com/v1/get_url"
   gatewayWWW = "https://www.deezer.com/ajax/gw-light.php"
)

type ping struct {
   Results struct { Session string }
}

func newPing() (ping, error) {
   req, err := http.NewRequest("GET", gatewayWWW, nil)
   if err != nil {
      return ping{}, err
   }
   val := url.Values{
      "api_token": {""}, "api_version": {"1.0"}, "method": {"deezer.ping"},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return ping{}, err
   }
   var data ping
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}
