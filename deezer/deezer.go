// Command line tool to download Deezer track
package main

import (
   "flag"
   "fmt"
   "github.com/89z/deezer"
   "io/ioutil"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func check(err error) {
   if err != nil {
      log.Fatal(err)
   }
}

func main() {
   config, err := os.UserConfigDir()
   check(err)
   config = filepath.Join(config, "deezer", "deezer.txt")
   var arl, format, sngId string
   flag.StringVar(&format, "f", "mp3", "format")
   flag.StringVar(&sngId, "s", "", "SNG_ID")
   flag.StringVar(&arl, "a", config, "Arl cookie value")
   flag.Parse()
   if sngId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   track, err := deezer.NewTrack(sngId, arl)
   check(err)
   var source string
   if format == "flac" {
      source, err = track.GetSource(sngId, deezer.FLAC)
   } else {
      source, err = track.GetSource(sngId, deezer.MP3_320)
   }
   check(err)
   get, err := http.Get(source)
   check(err)
   body, err := ioutil.ReadAll(get.Body)
   check(err)
   deezer.Decrypt(sngId, body)
   ioutil.WriteFile(
      fmt.Sprintf("%v - %v.%v", track.ArtName, track.SngTitle, format),
      body,
      os.ModePerm,
   )
}
