package deezer
import "testing"

func TestSong(t *testing.T) {
   user, err := newUser("arl", arl)
   if err != nil {
      t.Error(err)
   }
   s, err := newSong(user.Results.CheckForm, user.sid, felix)
   if err != nil {
      t.Error(err)
   }
   if s.Results.Track_Token == "" {
      t.Error()
   }
}
