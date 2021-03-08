package main

import (
   "flag"
   "fmt"
   "io"
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
   data, err := getData(token, sngId)
   check(err)
   source, err := getSource(sngId, data, deezer320)
   check(err)
   from, err := newReader(sngId, source)
   check(err)
   to, err := os.Create(
      fmt.Sprintf("%s - %s.mp3", data.ArtName, data.SngTitle),
   )
   check(err)
   io.Copy(to, from)
}
