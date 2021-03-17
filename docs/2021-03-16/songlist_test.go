package deezer
import "testing"

func TestSongList(t *testing.T) {
   user, err := newUser("arl", arl)
   if err != nil {
      t.Error(err)
   }
   list, err := newSongList(user.Results.CheckForm, user.sid, felix, maria)
   if err != nil {
      t.Error(err)
   }
   if list.Results.Data[0].MD5_Origin == "" {
      t.Error(err)
   }
   if list.Results.Data[0].Track_Token == "" {
      t.Error(err)
   }
}
