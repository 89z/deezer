package main

import (
   "encoding/json"
   "fmt"
   "io"
   "io/ioutil"
   "net/http"
   "net/url"
   "strings"
)

var API = url.URL{
   Scheme: "http", Host: "www.deezer.com", Path: "/ajax/gw-light.php",
}

func harDecode(har string) (string, string, error) {
   data, err := ioutil.ReadFile(har)
   if err != nil {
      return "", "", err
   }
   var archive httpArchive
   json.Unmarshal(data, &archive)
   for _, entry := range archive.Log.Entries {
      var api_token, sid string
      for _, query := range entry.Request.QueryString {
         if query.Name == "api_token" {
            api_token = query.Value
         }
      }
      for _, cookie := range entry.Request.Cookies {
         if cookie.Name == "sid" {
            sid = cookie.Value
         }
      }
      if api_token != "" && sid != "" {
         println("StartedDateTime", entry.StartedDateTime)
         println("api_token", api_token)
         println("sid", sid)
         return api_token, sid, nil
      }
   }
   return "", "", fmt.Errorf("token not found")
}

type Track struct {
   ArtName      string `json:"ART_NAME"`
   MD5Origin    string `json:"MD5_ORIGIN"`
   MediaVersion string `json:"MEDIA_VERSION"`
   SngTitle     string `json:"SNG_TITLE"`
   TrackTokenExpire int `json:"TRACK_TOKEN_EXPIRE"`
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

type pageTrack struct {
   Results struct {
      Data Track
   }
}

func newPageTrack(api_token, sid, sngId string) (pageTrack, error) {
   val, req := url.Values{}, &http.Request{URL: &API}
   val.Set("api_version", "1.0")
   val.Set("api_token", api_token)
   val.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = val.Encode()
   req.Method = "POST"
   req.Body = io.NopCloser(strings.NewReader(
      fmt.Sprintf(`{"sng_id": "%v"}`, sngId),
   ))
   req.Header = http.Header{}
   req.AddCookie(&http.Cookie{Name: "sid", Value: sid})
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return pageTrack{}, err
   }
   data, err := ioutil.ReadAll(res.Body)
   if err != nil {
      return pageTrack{}, err
   }
   fmt.Printf("%s\n", data)
   var page pageTrack
   return page, json.Unmarshal(data, &page)
}
