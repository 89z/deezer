package main

import (
   "flag"
   "log"
   "net/http"
   "os"
)

func main() {
   var conf config
   flag.StringVar(&conf.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&conf.trackId, "id", "", "Deezer Track ID")
   flag.Parse()
   if conf.trackId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   source, dest, err := getUrl(conf)
   if err != nil {
      log.Fatal(err)
   }
   get, err := http.Get(source)
   if err != nil {
      log.Fatal(err)
   }
   defer get.Body.Close()
   create, err := os.Create(dest)
   if err != nil {
      log.Fatal(err)
   }
   defer create.Close()
   err = decryptMedia(conf, get, create)
   if err != nil {
      log.Fatal(err)
   }
}
