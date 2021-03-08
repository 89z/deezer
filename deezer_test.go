package deezer

import (
   "io/ioutil"
   "os"
   "os/exec"
   "testing"
)

const token = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

const (
   hash = "87207d3416377217f835b887c74f4300"
   sngId = "75498418"
)

func _TestOld(t *testing.T) {
   c := exec.Command("download/download", "-id", sngId, "-usertoken", token)
   c.Stderr, c.Stdout = os.Stderr, os.Stdout
   err := c.Run()
   if err != nil {
      t.Error(err)
   }
   data, err := ioutil.ReadFile("Julia Holter - Für Felix.mp3")
   if err != nil {
      t.Error(err)
   }
   testHash := md5Hash(string(data))
   if testHash != hash {
      t.Error(testHash)
   }
}

func TestNew(t *testing.T) {
   data, err := GetData(sngId, token)
   if err != nil {
      t.Error(err)
   }
   source, err := GetSource(sngId, data, MP3_320)
   if err != nil {
      t.Error(err)
   }
   from, err := NewReader(sngId, source)
   if err != nil {
      t.Error(err)
   }
   to, err := ioutil.ReadAll(from)
   if err != nil {
      t.Error(err)
   }
   testHash := md5Hash(string(to))
   if testHash != hash {
      t.Error(testHash)
   }
   ioutil.WriteFile("fail.mp3", to, os.ModePerm)
}
