package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func main() {
   req, err := http.NewRequest(
      "GET", "http://www.deezer.com/ajax/gw-light.php", nil,
   )
   if err != nil {
      log.Fatal(err)
   }
   val := url.Values{}
   var c http.Client
   val.Set("api_version", "1.0")
   // getUserData
   val.Set("api_token", "")
   val.Set("method", "deezer.getUserData")
   req.URL.RawQuery = val.Encode()
   sid, err := harDecode("deezer.har")
   if err != nil {
      log.Fatal(err)
   }
   req.Header.Set("Cookie", "sid=" + sid)
   res, err := c.Do(req)
   if err != nil {
      log.Fatal(err)
   }
   var user userData
   json.NewDecoder(res.Body).Decode(&user)
   // getListData
   val.Set("api_token", user.Results.CheckForm)
   val.Set("method", "song.getListData")
   req.URL.RawQuery = val.Encode()
   req.Method = "POST"
   req.Body = io.NopCloser(strings.NewReader(`{"sng_ids": [137955757,960539]}`))
   res, err = c.Do(req)
   if err != nil {
      log.Fatal(err)
   }
   var list listData
   json.NewDecoder(res.Body).Decode(&list)
   // get_url
   type (
      a []interface{}
      m map[string]interface{}
   )
   request := m{
      "license_token": user.Results.User.Options.LicenseToken,
      "track_tokens": a{list.Results.Data[0].TrackToken},
      "media": a{m{
         "type": "FULL",
         "formats": a{m{"cipher": "BF_CBC_STRIPE", "format": "MP3_128"}},
      }},
   }
   data, err := json.Marshal(request)
   if err != nil {
      log.Fatal(err)
   }
   req, err = http.NewRequest(
      "POST", "https://media.deezer.com/v1/get_url", bytes.NewReader(data),
   )
   if err != nil {
      log.Fatal(err)
   }
   res, err = c.Do(req)
   var response urlResponse
   json.NewDecoder(res.Body).Decode(&response)
   fmt.Println(response.Data[0].Media[0].Sources[0].Url)
}

type urlResponse struct {
   Data []struct {
      Media []struct {
         Sources []struct {
            Url string
         }
      }
   }
}
