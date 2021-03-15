package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "os"
)

type (
   a []interface{}
   m map[string]interface{}
)

func getUrl(licenseTok string, trackToks ...string) error {
   in, out := m{
      "license_token": licenseTok,
      "track_tokens": a{trackToks},
      "media": a{m{
         "type": "FULL",
         "formats": a{m{"cipher": "BF_CBC_STRIPE", "format": "MP3_128"}},
      }},
   }, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayMedia, out)
   if err != nil {
      return err
   }
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return err
   }
   os.Stdout.ReadFrom(res.Body)
   return nil
}
