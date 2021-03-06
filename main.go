package main

import (
   "flag"
   "log"
   "net/http"
   "os"
)

func main() {
   var config configuration
   flag.StringVar(&config.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&config.trackId, "id", "", "Deezer Track ID")
   flag.Parse()
   if config.trackId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   var client http.Client
   downloadURL, FName, err := getUrl(config, client)
   if err != nil {
      log.Fatal(err)
   }
   resp, err := http.Get(downloadURL)
   if err != nil {
      log.Fatal(err)
   }
   defer resp.Body.Close()
   err = decryptMedia(resp, config.trackId, FName)
   if err != nil {
      log.Fatal(err)
   }
}
