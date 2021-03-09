// Command line tool to download Deezer track
package main

import (
   "encoding/json"
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

func colorGreen(s string) string {
   return "\x1b[92m" + s + "\x1b[m"
}

func getArl(har string) (string, error) {
   data, err := ioutil.ReadFile(har)
   if err != nil {
      return "", err
   }
   var archive httpArchive
   json.Unmarshal(data, &archive)
   for _, entry := range archive.Log.Entries {
      for _, cookie := range entry.Request.Cookies {
         if cookie.Name == "arl" {
            return cookie.Value, nil
         }
      }
   }
   return "", fmt.Errorf("Arl cookie not found")
}

func main() {
   var format, sngId string
   flag.StringVar(&format, "f", "mp3", "format")
   flag.StringVar(&sngId, "s", "", "SNG_ID")
   har, err := os.UserConfigDir()
   check(err)
   har = filepath.Join(har, "deezer.har")
   flag.StringVar(&har, "h", har, "HTTP archive")
   flag.Parse()
   if sngId == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   arl, err := getArl(har)
   check(err)
   track, err := deezer.NewTrack(sngId, arl)
   check(err)
   var source string
   if format == "flac" {
      source, err = track.Source(sngId, deezer.FLAC)
   } else {
      source, err = track.Source(sngId, deezer.MP3_320)
   }
   check(err)
   fmt.Println(colorGreen("Get"), source)
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

type httpArchive struct {
   Log struct {
      Entries []struct {
         Request struct {
            Cookies []struct {
               Name, Value string
            }
         }
      }
   }
}
