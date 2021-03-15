package deezer
import "testing"

const (
   apiToken = "4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE"
   felix = 75498418
   maria = 75498415
)

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func TestGetUrl(t *testing.T) {
}

func TestPing(t *testing.T) {
   ping, err := newPing()
   if err != nil {
      t.Error(err)
   }
   if ping.Results.Session == "" {
      t.Error()
   }
}

func TestSong(t *testing.T) {
   ping, err := newPing()
   if err != nil {
      t.Error(err)
   }
   song, err := newSong(apiToken, ping.Results.Session, felix)
   if err != nil {
      t.Error(err)
   }
   if song.Results.Track_Token == "" {
      t.Error()
   }
}

func TestSongList(t *testing.T) {
   ping, err := newPing()
   if err != nil {
      t.Error(err)
   }
   user, err := newUser("sid", ping.Results.Session)
   if err != nil {
      t.Error(err)
   }
   list, err := newSongList(user.Results.CheckForm, user.sid, felix, maria)
   if err != nil {
      t.Error(err)
   }
   if list.Results.Data[0].Track_Token == "" {
      t.Error(err)
   }
}

func TestTrack(t *testing.T) {
   user, err := newUser("arl", arl)
   if err != nil {
      t.Error(err)
   }
   track, err := newTrack(user.Results.CheckForm, user.sid, felix)
   if err != nil {
      t.Error(err)
   }
   if track.Results.Data.MD5_Origin != "9da3d60b427e895a0f1446a76b3d1488" {
      t.Error()
   }
}

func TestUserArl(t *testing.T) {
   user, err := newUser("arl", arl)
   if err != nil {
      t.Error(err)
   }
   if user.sid == "" {
      t.Error()
   }
   if user.Results.CheckForm == "" {
      t.Error()
   }
   if user.Results.User.Options.License_Token == "" {
      t.Error()
   }
}

func TestUserSid(t *testing.T) {
   ping, err := newPing()
   if err != nil {
      t.Error(err)
   }
   user, err := newUser("sid", ping.Results.Session)
   if err != nil {
      t.Error(err)
   }
   if user.Results.User.Options.License_Token == "" {
      t.Error()
   }
}
