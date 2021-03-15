package deezer
import "testing"

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
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   user, err := newUser("sid", p.Results.Session)
   if err != nil {
      t.Error(err)
   }
   if user.Results.User.Options.License_Token == "" {
      t.Error()
   }
}
