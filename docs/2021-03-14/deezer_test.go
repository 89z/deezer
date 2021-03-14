package deezer
import "testing"

const arl = "0e21c80ef0b963e68cf5d0a951fc918def86c2188a44b33ab353088f15d7b4" +
"087ed699e6dcd6293514f49439a7d2a7c86bdbcb6e0efae1acd029ec4f267a07b541bfe13872" +
"c5e5715db846bc784701c3794c328411b5cca332d695b37c1946c1"

func _TestPing(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   if p.body.Results.Session == "" {
      t.Error()
   }
}

func TestUserDataArl(t *testing.T) {
   data, err := newUserData("arl", arl)
   if err != nil {
      t.Error(err)
   }
   var sid string
   for _, each := range data.head {
      if each.Name == "sid" { sid = each.Value }
   }
   if sid == "" {
      t.Error()
   }
   if data.body.Results.CheckForm == "" {
      t.Error()
   }
   if data.body.Results.User.Options.License_Token == "" {
      t.Error()
   }
}
