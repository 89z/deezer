package main

import (
   "deezer"
   "flag"
   "fmt"
   "log"
   "os"
)

func check(err error) {
   if err != nil {
      log.Fatal(err)
   }
}

func main() {
   var sngId, token string
   flag.StringVar(&sngId, "id", "", "Deezer Track ID")
   flag.StringVar(&token, "usertoken", "", "Your Unique User Token")
   flag.Parse()
   if sngId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   data, err := deezer.GetData(sngId, token)
   check(err)
   source, err := deezer.GetSource(sngId, data, deezer.MP3_320)
   check(err)
   from, err := deezer.NewReader(sngId, source)
   check(err)
   to, err := os.Create(
      fmt.Sprintf("%s - %s.mp3", data.ArtName, data.SngTitle),
   )
   check(err)
   to.ReadFrom(from)
}
