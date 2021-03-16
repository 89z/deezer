package deezer
import "testing"

func TestTrackArl(t *testing.T) {
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

func TestTrackSid(t *testing.T) {
   ping, err := newPing()
   if err != nil {
      t.Error(err)
   }
   user, err := newUser("sid", ping.Results.Session)
   if err != nil {
      t.Error(err)
   }
   track, err := newTrack(user.Results.CheckForm, user.sid, felix)
   if err != nil {
      t.Error(err)
   }
   get, err := newGetUrl(
      user.Results.User.Options.License_Token, track.Results.Data.Track_Token,
   )
   if err != nil {
      t.Error(err)
   }
   if get.Data[0].Media[0].Sources[0].Url == "" {
      t.Error()
   }
}
