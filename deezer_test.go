package main

import (
   "os"
   "os/exec"
   "testing"
)

const token = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func TestMain(t *testing.T) {
   c := exec.Command("deezer.exe", "-id", "75498415", "-usertoken", token)
   c.Stderr, c.Stdout = os.Stderr, os.Stdout
   e := c.Run()
   if e != nil {
      t.Error(e)
   }
}
