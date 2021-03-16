package deezer

import (
   "io/ioutil"
   "os"
   "testing"
)

func TestDecrypt(t *testing.T) {
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
