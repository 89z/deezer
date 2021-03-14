package deezer
import "testing"

const (
   apiToken = "4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE"
   sngId = "75498418"
)

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func TestSongData(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   data, err := newSongData(sngId, apiToken, p.Results.Session)
   if err != nil {
      t.Error(err)
   }
   if data.Results.Track_Token == "" {
      t.Error()
   }
}

func _TestPageTrack(t *testing.T) {
   data, err := newUserData("arl", arl)
   if err != nil {
      t.Error(err)
   }
   track, err := newPageTrack(sngId, data.Results.CheckForm, data.sid)
   if err != nil {
      t.Error(err)
   }
   if track.Results.Data.MD5_Origin != "9da3d60b427e895a0f1446a76b3d1488" {
      t.Error()
   }
}

func _TestUserDataArl(t *testing.T) {
   data, err := newUserData("arl", arl)
   if err != nil {
      t.Error(err)
   }
   if data.sid == "" {
      t.Error()
   }
   if data.Results.CheckForm == "" {
      t.Error()
   }
   if data.Results.User.Options.License_Token == "" {
      t.Error()
   }
}

func _TestUserDataSid(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   data, err := newUserData("sid", p.Results.Session)
   if err != nil {
      t.Error(err)
   }
   if data.Results.User.Options.License_Token == "" {
      t.Error()
   }
}

func _TestPing(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   if p.Results.Session == "" {
      t.Error()
   }
}
