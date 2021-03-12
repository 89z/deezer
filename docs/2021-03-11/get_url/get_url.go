package main

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
)

func harDecode(har string) (string, error) {
   data, err := ioutil.ReadFile(har)
   if err != nil {
      return "", err
   }
   var archive httpArchive
   json.Unmarshal(data, &archive)
   for _, entry := range archive.Log.Entries {
      var sid string
      for _, cookie := range entry.Request.Cookies {
         if cookie.Name == "sid" {
            sid = cookie.Value
         }
      }
      if sid != "" {
         println("StartedDateTime", entry.StartedDateTime)
         println("sid", sid)
         return sid, nil
      }
   }
   return "", fmt.Errorf("sid not found")
}

type httpArchive struct {
   Log struct {
      Entries []struct {
         Request struct {
            Cookies []struct {
               Name, Value string
            }
            QueryString []struct {
               Name, Value string
            }
         }
         StartedDateTime string
      }
   }
}

type userData struct {
   Results struct {
      UserToken string `json:"USER_TOKEN"`
   }
}
