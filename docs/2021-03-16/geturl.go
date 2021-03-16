package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
)

type (
   a []interface{}
   m map[string]interface{}
)

type getUrl struct {
   Data []struct {
      Media []struct {
         Sources []struct { Url string }
      }
   }
   Errors []struct { Message string }
}

func newGetUrl(licenseTok string, trackToks ...string) (getUrl, error) {
   in, out := m{
      "license_token": licenseTok,
      "media": a{m{
         "type": "FULL",
         "formats": a{m{"cipher": "BF_CBC_STRIPE", "format": "MP3_128"}},
      }},
      "track_tokens": trackToks,
   }, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayMedia, out)
   if err != nil {
      return getUrl{}, err
   }
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return getUrl{}, err
   }
   var data getUrl
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}
