package deezer
import "testing"

func TestSong(t *testing.T) {
   user, err := newUser("arl", arl)
   if err != nil {
      t.Error(err)
   }
   song, err := newSong(user.Results.CheckForm, user.sid, felix)
   if err != nil {
      t.Error(err)
   }
   if song.Results.MD5_Origin == "" {
      t.Error()
   }
   if song.Results.Track_Token == "" {
      t.Error()
   }
}
