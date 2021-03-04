package main

import (
   "flag"
   "fmt"
   "os"
)

type Config struct {
   Debug     bool
   ID        string
   UserToken string
}

var cfg = &Config{false, "", ""}

func debug(msg string, params ...interface{}) {
   if cfg.Debug {
      fmt.Printf("\n"+msg+"\n\n", params...)
   }
}

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
