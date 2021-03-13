package main

import (
   "fmt"
   "log"
)

const cookieArl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15" +
"d7b4087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe1" +
"3872c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func main() {
   api, e := NewAPI()
   if e != nil {
      log.Fatal(e)
   }
   e = api.CookieLogin(cookieArl)
   if e != nil {
      log.Fatal(e)
   }
   track, e := api.GetSongData(75498415)
   if e != nil {
      log.Fatal(e)
   }
   fmt.Printf("%+v\n", track)
   /*
   download, e := track.GetDownloadURL(MP3_320)
   e = track.GetMD5()
   MobileApiRequest
   */
}
