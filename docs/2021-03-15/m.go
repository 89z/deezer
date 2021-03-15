package main
import "fmt"

type (
   a []interface{}
   m map[string]interface{}
)

func main() {
   b := m{
      "License_Token": 9,
      "Track_Tokens": a{9},
      "Media": a{m{
         "Type": 9,
         "Formats": a{m{
            "Cipher": 9, "Format": 9,
         }},
      }},
   }
   fmt.Printf("%+v\n", b)
}
