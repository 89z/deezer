package deezer

import (
   "encoding/json"
   "net/http"
   "net/url"
)

const gateway = "http://www.deezer.com/ajax/gw-light.php"

type ping struct {
   body struct {
      Results struct { Session string }
   }
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
   json.NewDecoder(res.Body).Decode(&p.body)
   return p, nil
}

type userData struct {
   head []*http.Cookie
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
   data.head = res.Cookies()
   json.NewDecoder(res.Body).Decode(&data.body)
   return data, nil
}
