package main

import (
   "fmt"
   "log"
)

func main() {
   api_token, sid, err := harDecode("deezer.har")
   if err != nil {
      log.Fatal(err)
   }
   track, err := newPageTrack(api_token, sid, "75498415")
   if err != nil {
      log.Fatal(err)
   }
   fmt.Printf("%q\n", track.Results.Data.MD5Origin)
   fmt.Println(track.Results.Data.TrackTokenExpire)
}
