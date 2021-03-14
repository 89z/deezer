package deezer

import (
   "encoding/json"
   "testing"
)

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func TestUserDataSid(t *testing.T){
   res, err := deezerPing()
   if err != nil {
      t.Error(err)
   }
   var ping struct {
      Results struct { Session string }
   }
   json.NewDecoder(res.Body).Decode(&ping)
   res, err = userData("sid", ping.Results.Session)
   if err != nil {
      t.Error(err)
   }
   var data struct {
      Results struct {
         User struct {
            Options struct { License_Token string }
         }
      }
   }
   json.NewDecoder(res.Body).Decode(&data)
   if data.Results.User.Options.License_Token == "" {
      t.Error()
   }
}

func TestUserDataArl(t *testing.T) {
   res, err := userData("arl", arl)
   if err != nil {
      t.Error(err)
   }
   var data struct {
      Results struct { CheckForm string }
   }
   json.NewDecoder(res.Body).Decode(&data)
   if data.Results.CheckForm == "" {
      t.Error()
   }
}

func TestPageTrack(t *testing.T) {
   res, err := userData("arl", arl)
   if err != nil {
      t.Error(err)
   }
   var sid string
   for _, each := range res.Cookies() {
      if each.Name == "sid" { sid = each.Value }
   }
   var userData struct {
      Results struct { CheckForm string }
   }
   json.NewDecoder(res.Body).Decode(&userData)
   res, err = pageTrack("75498418", userData.Results.CheckForm, sid)
   if err != nil {
      t.Error(err)
   }
   var track struct {
      Results struct {
         Data struct { MD5_Origin string }
      }
   }
   json.NewDecoder(res.Body).Decode(&track)
   md5 := track.Results.Data.MD5_Origin
   if md5 != "9da3d60b427e895a0f1446a76b3d1488" {
      t.Error(md5)
   }
}
