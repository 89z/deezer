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
   if track.Results.Data.MD5_Origin == "" {
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
   if track.Results.Data.Track_Token == "" {
      t.Error()
   }
}
