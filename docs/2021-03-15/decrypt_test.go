package deezer

import (
   "io/ioutil"
   "os"
   "testing"
)

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

const (
   felix = 75498418
   maria = 75498415
)

func _TestDecrypt(t *testing.T) {
   data, err := ioutil.ReadFile("maria-in.mp3")
   if err != nil {
      t.Error(err)
   }
   err = decrypt("75498415", data)
   if err != nil {
      t.Error(err)
   }
   ioutil.WriteFile("maria-out.mp3", data, os.ModePerm)
}
