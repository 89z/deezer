package deezer
import "testing"

func TestUser(t *testing.T) {
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
