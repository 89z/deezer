package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
   "os"
)

type song struct {
   Results struct { Track_Token string }
}

func newSong(apiToken, sid string, sngId int) (song, error) {
   in, out := map[string]int{"SNG_ID": sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayAPI, out)
   if err != nil {
      return song{}, err
   }
   val := url.Values{
      "api_key": {apiToken},
      "method": {"song.getData"},
      "output": {"3"},
      "sid": {sid},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return song{}, err
   }
   
   os.Stdout.ReadFrom(res.Body)
   return song{}, nil
   
   var data song
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}
