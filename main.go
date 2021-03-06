package main

import (
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
)

func main() {
   var conf config
   flag.StringVar(&conf.userToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&conf.trackId, "id", "", "Deezer Track ID")
   flag.Parse()
   if conf.trackId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   data, err := getData(conf)
   if err != nil {
      log.Fatal(err)
   }
   source, err := decryptDownload(data)
   if err != nil {
      log.Fatal(err)
   }
   get, err := http.Get(source)
   if err != nil {
      log.Fatal(err)
   }
   defer get.Body.Close()
   create, err := os.Create(
      fmt.Sprintf("%s - %s.mp3", data.ArtName, data.SngTitle),
   )
   if err != nil {
      log.Fatal(err)
   }
   defer create.Close()
   err = decryptMedia(conf, get, create)
   if err != nil {
      log.Fatal(err)
   }
}
