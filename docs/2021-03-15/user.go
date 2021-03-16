package deezer

import (
   "encoding/json"
   "net/http"
   "net/url"
)

type user struct {
   sid string
   Results struct {
      CheckForm string
      User struct {
         Options struct { License_Token string }
      }
   }
}

func newUser(name, value string) (user, error) {
   req, err := http.NewRequest("GET", gatewayWWW, nil)
   if err != nil {
      return user{}, err
   }
   val := url.Values{
      "api_token": {""},
      "api_version": {"1.0"},
      "input": {"3"},
      "method": {"deezer.getUserData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: name, Value: value}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return user{}, err
   }
   var data user
   for _, each := range res.Cookies() {
      if each.Name == "sid" { data.sid = each.Value }
   }
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}
