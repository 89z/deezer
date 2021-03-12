package main

import (
   "fmt"
   "io/ioutil"
   "log"
   "net/http"
   "net/url"
)

func main() {
   sid, err := harDecode("deezer.har")
   if err != nil {
      log.Fatal(err)
   }
   req, err := http.NewRequest(
      "GET", "http://www.deezer.com/ajax/gw-light.php", nil,
   )
   if err != nil {
      log.Fatal(err)
   }
   val := url.Values{}
   val.Set("method", "deezer.getUserData")
   val.Set("input", "3")
   val.Set("api_version", "1.0")
   val.Set("api_token", "")
   req.URL.RawQuery = val.Encode()
   req.Header.Set("Cookie", "sid=" + sid)
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      log.Fatal(err)
   }
   var data userData
   json.NewDecoder(get.Body).Decode(&data)
   fmt.Println(data)
}
