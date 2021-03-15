package deezer
import "testing"

func TestGetUrl(t *testing.T) {
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
   get, err := newGetUrl(
      user.Results.User.Options.License_Token,
      list.Results.Data[0].Track_Token,
      list.Results.Data[1].Track_Token,
   )
   if err != nil {
      t.Error(err)
   }
   if get.Data[0].Media[0].Sources[0].Url == "" {
      t.Error()
   }
}
