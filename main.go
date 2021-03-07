package main

import (
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "os"
)

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
   if err != nil {
      log.Fatal(err)
   }
   source, err := getSource(sngId, data, deezer320)
   if err != nil {
      log.Fatal(err)
   }
   create, err := os.Create(
      fmt.Sprintf("%s - %s.mp3", data.ArtName, data.SngTitle),
   )
   if err != nil {
      log.Fatal(err)
   }
   defer create.Close()
   err = oldRead(sngId, source, create)
   if err != nil {
      log.Fatal(err)
   }
}

func oldRead(sngId, from string, to io.Writer) error {
   get, err := http.Get(from)
   if err != nil {
      return err
   }
   return decryptAudio(sngId, get.Body, to)
}

func newRead(sngId, from string, to io.Writer) error {
   source, err := newReader(sngId, from)
   if err != nil {
      return err
   }
   _, err = io.Copy(to, source)
   return err
}
