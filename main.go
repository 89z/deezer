package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "io/ioutil"
   "log"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "strconv"
   "time"
)

const (
   // APIUrl is the deezer API
   APIUrl = "http://www.deezer.com/ajax/gw-light.php"
   // LoginURL is the API for deezer login
   LoginURL = "https://www.deezer.com/ajax/action.php"
   // deezer domain for cookie check
   Domain = "https://www.deezer.com"
)

var cfg = new(Config)

// GetAudioFile gets the audio file from deezer server
func GetAudioFile(downloadURL, id, FName string, client *http.Client) error {
   req, e := newRequest(downloadURL, "GET", nil)
   if e != nil {
      return e
   }
   resp, e := client.Do(req)
   if e != nil {
      return e
   }
   e = DecryptMedia(resp.Body, id, FName, resp.ContentLength)
   if e != nil {
      return e
   }
   return resp.Body.Close()
}

// GetUrlDownload get the url for the requested track
func GetUrlDownload(id string, client *http.Client) (string, string, *http.Client, error) {
   jsonTrack := &DeezTrack{}
   APIToken, _ := GetToken(client)
   jsonPrep := `{"sng_id":"` + id + `"}`
   jsonStr := []byte(jsonPrep)
   req, e := newRequest(APIUrl, "POST", jsonStr)
   if e != nil {
      return "", "", nil, e
   }
   qs := url.Values{}
   qs.Add("api_version", "1.0")
   qs.Add("api_token", APIToken)
   qs.Add("input", "3")
   qs.Add("method", "deezer.pageTrack")
   req.URL.RawQuery = qs.Encode()
   resp, _ := client.Do(req)
   body, _ := ioutil.ReadAll(resp.Body)
   defer resp.Body.Close()
   e = json.Unmarshal(body, &jsonTrack)
   if e != nil {
      return "", "", nil, e
   }
   FileSize320, _ := jsonTrack.Results.DATA.FileSize320.Int64()
   FileSize256, _ := jsonTrack.Results.DATA.FileSize256.Int64()
   FileSize128, _ := jsonTrack.Results.DATA.FileSize128.Int64()
   var format string
   switch {
   case FileSize320 > 0:
      format = "3"
   case FileSize256 > 0:
      format = "5"
   case FileSize128 > 0:
      format = "1"
   default:
      format = "8"
   }
   songID := jsonTrack.Results.DATA.ID.String()
   md5Origin := jsonTrack.Results.DATA.MD5Origin
   mediaVersion := jsonTrack.Results.DATA.MediaVersion.String()
   songTitle := jsonTrack.Results.DATA.SngTitle
   artName := jsonTrack.Results.DATA.ArtName
   FName := fmt.Sprintf("%s - %s.mp3", songTitle, artName)
   downloadURL, e := DecryptDownload(md5Origin, songID, format, mediaVersion)
   if e != nil {
      return "", "", nil, e
   }
   return downloadURL, FName, client, nil
}

// Login will login the user with the provided credentials
func Login() (*http.Client, error) {
   CookieJar, _ := cookiejar.New(nil)
   client := &http.Client{Jar: CookieJar}
   Deez := &DeezStruct{}
   req, e := newRequest(APIUrl, "POST", nil)
   args := []string{"null", "deezer.getUserData"}
   req = addQs(req, args...)
   resp, e := client.Do(req)
   body, _ := ioutil.ReadAll(resp.Body)
   e = json.Unmarshal(body, &Deez)
   if e != nil {
      return nil, e
   }
   CookieURL, _ := url.Parse(Domain)
   resp.Body.Close()
   form := url.Values{}
   form.Add("type", "login")
   form.Add("checkFormLogin", Deez.Results.CheckFormLogin)
   req, e = newRequest(LoginURL, "POST", form.Encode())
   if e != nil {
      return nil, e
   }
   req.Header.Set("Content-type", "application/x-www-form-urlencoded")
   req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
   resp, e = client.Do(req)
   if e != nil {
      return nil, e
   }
   defer resp.Body.Close()
   body, e = ioutil.ReadAll(resp.Body)
   if e != nil {
      return nil, e
   }
   if resp.StatusCode == 200 {
      addCookies(client, CookieURL)
      return client, nil
   }
   return nil, fmt.Errorf("Statuscode %v %v", resp.StatusCode, e)
}

func addCookies(client *http.Client, CookieURL *url.URL) {
   expire := time.Now().Add(time.Hour * 24 * 180)
   expire.Format("2006-01-02T15:04:05.999Z07:00")
   creation := time.Now().Format("2006-01-02T15:04:05.999Z07:00")
   lastUsed := time.Now().Format("2006-01-02T15:04:05.999Z07:00")
   rawcookie := fmt.Sprintf(
      "arl=%s; expires=%v; %s creation=%v; lastAccessed=%v;",
      cfg.UserToken,
      expire,
      "path=/; domain=deezer.com; max-age=15552000; httponly=true; hostonly=false;",
      creation,
      lastUsed,
   )
   cookies := []*http.Cookie{{
      Domain:   ".deezer.com",
      Expires:  expire,
      HttpOnly: true,
      MaxAge:   15552000,
      Name:     "arl",
      Path:     "/",
      Raw:      rawcookie,
      Value:    cfg.UserToken,
   }}
   client.Jar.SetCookies(CookieURL, cookies)
}

func main() {
   flag.StringVar(&cfg.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&cfg.ID, "id", "", "Deezer Track ID")
   flag.Parse()
   if cfg.ID == "" {
      fmt.Println("Error: Must have Deezer Track(Song) ID")
      fmt.Println("deezer --id 3135556 --usertoken UserToken_here")
      flag.PrintDefaults()
      os.Exit(1)
   }
   client, e := Login()
   if e != nil {
      log.Fatal(e)
   }
   downloadURL, FName, client, e := GetUrlDownload(cfg.ID, client)
   if e != nil {
      log.Fatal(e)
   }
   e = GetAudioFile(downloadURL, cfg.ID, FName, client)
   if e != nil {
      log.Fatal(e)
   }
}
