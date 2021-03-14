package deezer

import (
   "encoding/json"
   "net/http"
   "net/url"
)

const gateway = "http://www.deezer.com/ajax/gw-light.php"

type ping struct { *http.Response }

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
   return ping{res}, err
}

func (p ping) session() string {
   var body struct {
      Results struct { Session string }
   }
   json.NewDecoder(p.Body).Decode(&body)
   return body.Results.Session
}

type userData struct { *http.Response }

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
   return userData{res}, err
}

func (d userData) sid() string {
   for _, each := range d.Cookies() {
      if each.Name == "sid" { return each.Value }
   }
   return ""
}

func (d userData) token() string {
   var body struct {
      Results struct {
         User struct {
            Options struct { License_Token string }
         }
      }
   }
   json.NewDecoder(d.Body).Decode(&body)
   return body.Results.User.Options.License_Token
}
