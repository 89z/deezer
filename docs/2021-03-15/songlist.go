package deezer

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

type songList struct {
   Results struct {
      Data []struct { MD5_Origin, Track_Token string }
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
