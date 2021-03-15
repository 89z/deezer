package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

const (
   gatewayAPI = "http://api.deezer.com/1.0/gateway.php"
   gatewayMedia = "https://media.deezer.com/v1/get_url"
   gatewayWWW = "https://www.deezer.com/ajax/gw-light.php"
)

type (
   a []interface{}
   m map[string]interface{}
)

type getUrl struct {
   Data []struct {
      Media []struct {
         Sources []struct { Url string }
      }
   }
}

func newGetUrl(licenseTok string, trackToks ...string) (getUrl, error) {
   in, out := m{
      "license_token": licenseTok,
      "media": a{m{
         "type": "FULL",
         "formats": a{m{"cipher": "BF_CBC_STRIPE", "format": "MP3_128"}},
      }},
      "track_tokens": trackToks,
   }, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayMedia, out)
   if err != nil {
      return getUrl{}, err
   }
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return getUrl{}, err
   }
   var data getUrl
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}

type ping struct {
   Results struct { Session string }
}

func newPing() (ping, error) {
   req, err := http.NewRequest("GET", gatewayWWW, nil)
   if err != nil {
      return ping{}, err
   }
   val := url.Values{
      "api_token": {""}, "api_version": {"1.0"}, "method": {"deezer.ping"},
   }
   req.URL.RawQuery = val.Encode()
   var client http.Client
   res, err := client.Do(req)
   if err != nil {
      return ping{}, err
   }
   var data ping
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}

type song struct {
   Results struct { Track_Token string }
}

func newSong(apiToken, sid string, sngId int) (song, error) {
   in, out := map[string]int{"SNG_ID": sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayAPI, out)
   if err != nil {
      return song{}, err
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
      return song{}, err
   }
   var data song
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}

type songList struct {
   Results struct {
      Data []struct { Track_Token string }
   }
}

func newSongList(apiToken, sid string, sngIds ...int) (songList, error) {
   in, out := map[string][]int{"SNG_IDS": sngIds}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayWWW, out)
   if err != nil {
      return songList{}, err
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
      return songList{}, err
   }
   var list songList
   json.NewDecoder(res.Body).Decode(&list)
   return list, nil
}

type track struct {
   Results struct {
      Data struct { MD5_Origin string }
   }
}

func newTrack(apiToken, sid string, sngId int) (track, error) {
   in, out := map[string]int{"SNG_ID": sngId}, new(bytes.Buffer)
   json.NewEncoder(out).Encode(in)
   req, err := http.NewRequest("POST", gatewayWWW, out)
   if err != nil {
      return track{}, err
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
      return track{}, err
   }
   var page track
   json.NewDecoder(res.Body).Decode(&page)
   return page, nil
}

type user struct {
   sid string
   Results struct {
      CheckForm string
      User struct {
         Options struct { License_Token string }
      }
   }
}

func newUser(name, value string) (user, error) {
   req, err := http.NewRequest("GET", gatewayWWW, nil)
   if err != nil {
      return user{}, err
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
      return user{}, err
   }
   var data user
   for _, each := range res.Cookies() {
      if each.Name == "sid" { data.sid = each.Value }
   }
   json.NewDecoder(res.Body).Decode(&data)
   return data, nil
}
