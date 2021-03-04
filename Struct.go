package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "os"
)

var cfg = new(Config)

func ErrorUsage() {
   fmt.Println(`Guide: go-decrypt-deezer [--debug --id --usertoken`)
   fmt.Println(`How Do I Get My UserToken?: https://notabug.org/RemixDevs/DeezloaderRemix/wiki/Login+via+userToken`)
   fmt.Println(`Example: go-decrypt-deezer --id 3135556 --usertoken UserToken_here`)
   flag.PrintDefaults()
   os.Exit(1)
}

func init() {
   flag.BoolVar(&cfg.Debug, "debug", false, "Turn on debuging mode.")
   flag.StringVar(&cfg.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&cfg.ID, "id", "", "Deezer Track ID")
   flag.Parse()
   if cfg.ID == "" {
      fmt.Println("Error: Must have Deezer Track(Song) ID")
      ErrorUsage()
   }
}

type Config struct {
   Debug     bool
   ID        string
   UserToken string
}

type Data struct {
   DATA *TrackData `json:"DATA"`
}

type DeezStruct struct {
   Error   []string    `json:"error,omitempty"`
   Results *ResultList `json:"results,omitempty"`
}

type DeezTrack struct {
   Error   []string `json:"error,omitempty"`
   Results *Data    `json:"results,omitempty"`
}

type OnError struct {
   Error   error
   Message string
}

type ResultList struct {
   DeezToken      string `json:"checkForm,omitempty"`
   CheckFormLogin string `json:"checkFormLogin,omitempty"`
}

type TrackData struct {
   ID           json.Number `json:"SNG_ID"`
   MD5Origin    string      `json:"MD5_ORIGIN"`
   FileSize320  json.Number `json:"FILESIZE_MP3_320"`
   FileSize256  json.Number `json:"FILESIZE_MP3_256"`
   FileSize128  json.Number `json:"FILESIZE_MP3_128"`
   MediaVersion json.Number `json:"MEDIA_VERSION"`
   SngTitle     string      `json:"SNG_TITLE"`
   ArtName      string      `json:"ART_NAME"`
}
