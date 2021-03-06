package main

import (
   "encoding/json"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
)

var cfg = new(Config)

func getToken(httpClient http.Client) (string, error) {
   // we must use Request, as cookies are required
   req, err := http.NewRequest("GET", APIURL, nil)
   if err != nil {
      return "", err
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", "null")
   qs.Set("input", "3")
   qs.Set("method", "deezer.getUserData")
   req.URL.RawQuery = qs.Encode()
   req.AddCookie(&http.Cookie{Name: "arl", Value: cfg.UserToken})
   resp, err := httpClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var deez DeezStruct
   err = json.NewDecoder(resp.Body).Decode(&deez)
   if err != nil {
      return "", err
   }
   return deez.Results.DeezToken, nil
}

func main() {
   flag.StringVar(&cfg.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&cfg.trackId, "id", "", "Deezer Track ID")
   flag.Parse()
   if cfg.trackId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   var client http.Client
   downloadURL, FName, err := getUrl(cfg.trackId, client)
   if err != nil {
      log.Fatal(err)
   }
   err = getAudioFile(downloadURL, cfg.trackId, FName)
   if err != nil {
      log.Fatal(err)
   }
}

type Config struct {
   UserToken string
   trackId string
}
