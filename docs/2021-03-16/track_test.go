package deezer
import "testing"

func TestTrack(t *testing.T) {
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
