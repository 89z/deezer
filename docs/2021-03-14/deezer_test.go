package deezer
import "testing"

func TestPing(t *testing.T) {
   p, err := newPing()
   if err != nil {
      t.Error(err)
   }
   if p.session() == "" {
      t.Error()
   }
}
