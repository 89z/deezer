package main

import (
   "bytes"
   "encoding/json"
   "log"
   "net/http"
   "net/url"
)

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

const gateway = "http://www.deezer.com/ajax/gw-light.php"

func getUserData(arl string) (*http.Response, error) {
   req, err := http.NewRequest("GET", gateway, nil)
   if err != nil {
      return nil, err
   }
   val := url.Values{
      "api_token": {""},
      "api_version": {"1.0"},
      "method": {"deezer.getUserData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: "arl", Value: arl}
   req.AddCookie(&cookie)
   var client http.Client
   return client.Do(req)
}

func pageTrack(sngId, apiToken string) (*http.Response, error) {
   in, out := struct{SngId string}{sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gateway, out)
   if err != nil {
      return nil, err
   }
   val := url.Values{
      "api_token": {apiToken}, "method": {"deezer.pageTrack"},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   return client.Do(req)
}

func main() {
   res, err := getUserData(arl)
   if err != nil {
      log.Fatal(err)
   }
   println(res)
}
