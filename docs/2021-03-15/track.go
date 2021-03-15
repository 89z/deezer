package deezer

import (
   "encoding/json"
   "net/http"
)

type track struct {
   Results struct {
      Data struct { MD5_Origin string }
   }
}

func newTrack(apiToken, sid string, sngId int) (track, error) {
   in, out := map[string]int{"SNG_ID": sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayWWW, out)
   if err != nil {
      return track{}, err
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
      return track{}, err
   }
   var page track
   json.NewDecoder(res.Body).Decode(&page)
   return page, nil
}
