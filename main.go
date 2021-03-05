package main

import (
   "flag"
   "log"
   "os"
)

func main() {
   flag.StringVar(&cfg.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&cfg.ID, "id", "", "Deezer Track ID")
   flag.Parse()
   if cfg.ID == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   id := cfg.ID
   client, err := login()
   if err != nil {
      log.Fatal(err)
   }
   downloadURL, FName, client, err := getUrl(id, client)
   if err != nil {
      log.Fatal(err)
   }
   err = getAudioFile(downloadURL, id, FName, client)
   if err != nil {
      log.Fatal(err)
   }
}
