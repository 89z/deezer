package main

import (
   "flag"
   "fmt"
   "github.com/89z/deezer"
   "io/ioutil"
   "log"
   "net/http"
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
   get, err := http.Get(source)
   check(err)
   body, err := ioutil.ReadAll(get.Body)
   check(err)
   deezer.Decrypt(sngId, body)
   ioutil.WriteFile(
      fmt.Sprintf("%s - %s.mp3", data.ArtName, data.SngTitle), body, os.ModePerm,
   )
}
