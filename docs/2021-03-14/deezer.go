package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

const gateway = "http://www.deezer.com/ajax/gw-light.php"

type pageTrack struct {
   Results struct {
      Data struct { MD5_Origin string }
   }
}

func newPageTrack(sngId, apiToken, sid string) (pageTrack, error) {
   in, out := struct{SNG_ID string}{sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gateway, out)
   if err != nil {
      return pageTrack{}, err
   }
   val := url.Values{
      "api_token": {apiToken},
      "api_version": {"1.0"},
      "method": {"deezer.pageTrack"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: "sid", Value: sid}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return pageTrack{}, err
   }
   var track pageTrack
   json.NewDecoder(res.Body).Decode(&track)
   return track, nil
}

type ping struct {
   Results struct { Session string }
}

func newPing() (ping, error) {
   req, err := http.NewRequest("GET", gateway, nil)
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
   var p ping
   json.NewDecoder(res.Body).Decode(&p)
   return p, nil
}

type userData struct {
   sid string
   body struct {
      Results struct {
         CheckForm string
         User struct {
            Options struct { License_Token string }
         }
      }
   }
}

func newUserData(name, value string) (userData, error) {
   req, err := http.NewRequest("GET", gateway, nil)
   if err != nil {
      return userData{}, err
   }
   val := url.Values{
      "api_token": {""},
      "api_version": {"1.0"},
      "method": {"deezer.getUserData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: name, Value: value}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return userData{}, err
   }
   var data userData
   for _, each := range res.Cookies() {
      if each.Name == "sid" { data.sid = each.Value }
   }
   json.NewDecoder(res.Body).Decode(&data.body)
   return data, nil
}
