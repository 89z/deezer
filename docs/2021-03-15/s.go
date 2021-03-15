package main
import "fmt"

type (
   body struct {
      License_Token int
      Track_Tokens []int
      Media []media
   }
   media struct {
      Type int
      Formats []formats
   }
   formats struct { Cipher, Format int }
)

func main() {
   b := body{
      License_Token: 9,
      Track_Tokens: []int{9},
      Media: []media{{
         Type: 9,
         Formats: []formats{{
            Cipher: 9, Format: 9,
         }},
      }},
   }
   fmt.Printf("%+v\n", b)
}
