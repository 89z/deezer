package main

import (
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
)

func main() {
   var token, trackId string
   flag.StringVar(&token, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&trackId, "id", "", "Deezer Track ID")
   flag.Parse()
   if trackId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   data, err := getData(token, trackId)
   if err != nil {
      log.Fatal(err)
   }
   source, err := getSource(data, deezer320)
   if err != nil {
      log.Fatal(err)
   }
   fmt.Println("GET", source)
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
   err = decryptAudio(trackId, get.Body, create)
   if err != nil {
      log.Fatal(err)
   }
}
