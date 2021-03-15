package deezer
import "testing"

func TestSongList(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   user, err := newUser("sid", p.Results.Session)
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
