package deezer

import (
   "net/http"
   "net/url"
)

const gateway = "http://www.deezer.com/ajax/gw-light.php"

type ping *http.Response

func newPing() (ping, error) {
   req, err := http.NewRequest("GET", gateway, nil)
   if err != nil { return nil, err }
   val := url.Values{
      "api_token": {""}, "api_version": {"1.0"}, "method": {"deezer.ping"},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   return client.Do(req)
}
