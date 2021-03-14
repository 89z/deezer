package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

const gateway = "http://www.deezer.com/ajax/gw-light.php"

func deezerPing() (*http.Response, error) {
   req, err := http.NewRequest("GET", gateway, nil)
   if err != nil { return nil, err }
   val := url.Values{
      "api_token": {""}, "api_version": {"1.0"}, "method": {"deezer.ping"},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   return client.Do(req)
}

func userData(name, value string) (*http.Response, error) {
   req, err := http.NewRequest("GET", gateway, nil)
   if err != nil { return nil, err }
   val := url.Values{
      "api_token": {""},
      "api_version": {"1.0"},
      "method": {"deezer.getUserData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: name, Value: value}
   req.AddCookie(&cookie)
   var client http.Client
   return client.Do(req)
}

func pageTrack(sngId, apiToken, sid string) (*http.Response, error) {
   in, out := struct{SNG_ID string}{sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gateway, out)
   if err != nil { return nil, err }
   val := url.Values{
      "api_token": {apiToken},
      "api_version": {"1.0"},
      "method": {"deezer.pageTrack"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: "sid", Value: sid}
   req.AddCookie(&cookie)
   var client http.Client
   return client.Do(req)
}
