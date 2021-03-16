package deezer
import "testing"

func TestGetUrlArl(t *testing.T) {
   user, err := newUser("arl", arl)
   if err != nil {
      t.Error(err)
   }
   list, err := newSongList(user.Results.CheckForm, user.sid, felix)
   if err != nil {
      t.Error(err)
   }
   get, err := newGetUrl(
      user.Results.User.Options.License_Token, list.Results.Data[0].Track_Token,
   )
   if err != nil {
      t.Error(err)
   }
   if len(get.Data) == 0 {
      t.Error(get.Errors[0].Message)
   }
}

func TestGetUrlSid(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   user, err := newUser("sid", p.Results.Session)
   if err != nil {
      t.Error(err)
   }
   list, err := newSongList(user.Results.CheckForm, user.sid, felix)
   if err != nil {
      t.Error(err)
   }
   get, err := newGetUrl(
      user.Results.User.Options.License_Token, list.Results.Data[0].Track_Token,
   )
   if err != nil {
      t.Error(err)
   }
   if len(get.Data) == 0 {
      t.Error(get.Errors[0].Message)
   }
}
