package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

type song struct {
   Results struct { Track_Token string }
}

func newSong(apiToken, sid string, sngId int) (song, error) {
   in, out := map[string]int{"SNG_ID": sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayWWW, out)
   if err != nil {
      return song{}, err
   }
   val := url.Values{
      "api_token": {apiToken},
      "api_version": {"1.0"},
      "input": {"3"},
      "method": {"song.getData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: "sid", Value: sid}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return song{}, err
   }
   var data song
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}
