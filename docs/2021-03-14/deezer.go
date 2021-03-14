package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

const (
   gatewayAPI = "http://api.deezer.com/1.0/gateway.php"
   gatewayWWW = "https://www.deezer.com/ajax/gw-light.php"
)

type pingRes struct {
   Results struct { Session string }
}

func newPingRes() (pingRes, error) {
   req, err := http.NewRequest("GET", gatewayWWW, nil)
   if err != nil {
      return pingRes{}, err
   }
   val := url.Values{
      "api_token": {""}, "api_version": {"1.0"}, "method": {"deezer.ping"},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return pingRes{}, err
   }
   var ping pingRes
   json.NewDecoder(res.Body).Decode(&ping)
   return ping, nil
}

type songListReq struct {
   Sng_Ids []int
}

type songListRes struct {
   Results struct {
      Data []struct { Track_Token string }
   }
}

func newSongListRes(apiToken, sid string, sngIds ...int) (songListRes, error) {
   in, out := songListReq{sngIds}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayWWW, out)
   if err != nil {
      return songListRes{}, err
   }
   val := url.Values{
      "api_token": {apiToken},
      "api_version": {"1.0"},
      "input": {"3"},
      "method": {"song.getListData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: "sid", Value: sid}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return songListRes{}, err
   }
   var list songListRes
   json.NewDecoder(res.Body).Decode(&list)
   return list, nil
}

type songReq struct {
   Sng_Id int
}

type songRes struct {
   Results struct { Track_Token string }
}

func newSongRes(apiToken, sid string, sngId int) (songRes, error) {
   in, out := songReq{sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayAPI, out)
   if err != nil {
      return songRes{}, err
   }
   val := url.Values{
      "api_key": {apiToken},
      "method": {"song.getData"},
      "output": {"3"},
      "sid": {sid},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return songRes{}, err
   }
   var song songRes
   json.NewDecoder(res.Body).Decode(&song)
   return song, nil
}

type trackRes struct {
   Results struct {
      Data struct { MD5_Origin string }
   }
}

func newTrackRes(apiToken, sid string, sngId int) (trackRes, error) {
   in, out := songReq{sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayWWW, out)
   if err != nil {
      return trackRes{}, err
   }
   val := url.Values{
      "api_token": {apiToken},
      "api_version": {"1.0"},
      "method": {"deezer.pageTrack"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: "sid", Value: sid}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return trackRes{}, err
   }
   var track trackRes
   json.NewDecoder(res.Body).Decode(&track)
   return track, nil
}

type urlReq struct {
   License_Token string
   Media []struct {
      Type string
      Formats []struct { Cipher, Format string }
   }
   Track_Tokens []string
}

type userRes struct {
   sid string
   Results struct {
      CheckForm string
      User struct {
         Options struct { License_Token string }
      }
   }
}

func newUserRes(name, value string) (userRes, error) {
   req, err := http.NewRequest("GET", gatewayWWW, nil)
   if err != nil {
      return userRes{}, err
   }
   val := url.Values{
      "api_token": {""},
      "api_version": {"1.0"},
      "input": {"3"},
      "method": {"deezer.getUserData"},
   }
   req.URL.RawQuery = val.Encode()
   cookie := http.Cookie{Name: name, Value: value}
   req.AddCookie(&cookie)
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return userRes{}, err
   }
   var user userRes
   for _, each := range res.Cookies() {
      if each.Name == "sid" { user.sid = each.Value }
   }
   json.NewDecoder(res.Body).Decode(&user)
   return user, nil
}
