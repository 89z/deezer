package deezer

import (
   "io/ioutil"
   "net/http"
   "testing"
)

const sngId = "75498418"

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func TestArl(t *testing.T) {
   track, err := NewTrack(sngId, arl)
   if err != nil {
      t.Error(err)
   }
   source, err := track.Source(sngId, MP3_320)
   if err != nil {
      t.Error(err)
   }
   get, err := http.Get(source)
   if err != nil {
      t.Error(err)
   }
   body, err := ioutil.ReadAll(get.Body)
   if err != nil {
      t.Error(err)
   }
   Decrypt(sngId, body)
   testHash := md5Hash(string(body))
   if testHash != "87207d3416377217f835b887c74f4300" {
      t.Error(testHash)
   }
}
