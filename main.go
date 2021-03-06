package main

import (
   "flag"
   "log"
   "net/http"
   "os"
)

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
