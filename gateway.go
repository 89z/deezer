package deezer

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "strings"
)

func gateway(sngId, apiKey string) error {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return err
   }
   http.DefaultClient.Jar = jar
   // GET
   val := url.Values{}
   req, err := http.NewRequest("GET", API.String(), nil)
   if err != nil {
      return err
   }
   val.Set("api_version", "1.0")
   val.Set("api_token", "null")
   val.Set("method", "deezer.ping")
   req.URL.RawQuery = val.Encode()
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   // JSON
   body, err := ioutil.ReadAll(res.Body)
   if err != nil {
      return err
   }
   var ping struct {
      Results struct {
         Session string
      }
   }
   json.Unmarshal(body, &ping)
   // POST
   req, err = http.NewRequest(
      "POST",
      "https://api.deezer.com/1.0/gateway.php",
      strings.NewReader(fmt.Sprintf(`{"SNG_ID": %v}`, sngId)),
   )
   if err != nil {
      return err
   }
   val.Set("api_key", apiKey)
   val.Set("method", "song_getData")
   val.Set("output", "3")
   val.Set("input", "3")
   println(ping.Results.Session)
   val.Set("sid", ping.Results.Session)
   req.URL.RawQuery = val.Encode()
   res, err = http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   // JSON
   body, err = ioutil.ReadAll(res.Body)
   if err != nil {
      return err
   }
   // memcpy(md5, find(page, "PUID\":\""), 32);
   fmt.Printf("%s\n", body)
   return nil
}
