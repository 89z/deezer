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

func Login() (*http.Client, error) {
   cookieJar, _ := cookiejar.New(nil)
   client := &http.Client{Jar: cookieJar}
   req, e := newRequest(APIUrl, "POST", nil)
   args := []string{"null", "deezer.getUserData"}
   req = addQs(req, args...)
   resp, e := client.Do(req)
   if e != nil {
      return nil, e
   }
   body, e := ioutil.ReadAll(resp.Body)
   if e != nil {
      return nil, e
   }
   var deez DeezStruct
   e = json.Unmarshal(body, &deez)
   if e != nil {
      return nil, e
   }
   CookieURL, _ := url.Parse(Domain)
   resp.Body.Close()
   form := url.Values{}
   form.Add("type", "login")
   form.Add("checkFormLogin", deez.Results.CheckFormLogin)
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

func getAudioFile(downloadURL, id, fName string, client *http.Client) error {
   req, e := newRequest(downloadURL, "GET", nil)
   if e != nil {
      return e
   }
   resp, e := client.Do(req)
   if e != nil {
      return e
   }
   e = DecryptMedia(resp.Body, id, fName, resp.ContentLength)
   if e != nil {
      return e
   }
   return resp.Body.Close()
}

func getUrl(id string, client *http.Client) (string, string, *http.Client, error) {
   APIToken, e := GetToken(client)
   if e != nil {
      return "", "", nil, e
   }
   qs := url.Values{}
   qs.Add("api_token", APIToken)
   qs.Add("api_version", "1.0")
   qs.Add("input", "3")
   qs.Add("method", "deezer.pageTrack")
   jsonPrep := fmt.Sprintf(`{"sng_id": "%v"}`, id)
   req, e := newRequest(
      APIUrl, "POST", []byte(jsonPrep),
   )
   if e != nil {
      return "", "", nil, e
   }
   req.URL.RawQuery = qs.Encode()
   resp, e := client.Do(req)
   if e != nil {
      return "", "", nil, e
   }
   defer resp.Body.Close()
   body, e := ioutil.ReadAll(resp.Body)
   if e != nil {
      return "", "", nil, e
   }
   var jsonTrack DeezTrack
   e = json.Unmarshal(body, &jsonTrack)
   if e != nil {
      return "", "", nil, e
   }
   var format string
   switch {
   case jsonTrack.Results.DATA.FileSize320 > 0:
      format = "3"
   case jsonTrack.Results.DATA.FileSize256 > 0:
      format = "5"
   default:
      format = "8"
   }
   downloadURL, e := DecryptDownload(
      jsonTrack.Results.DATA.MD5Origin,
      jsonTrack.Results.DATA.ID.String(),
      format,
      jsonTrack.Results.DATA.MediaVersion.String(),
   )
   if e != nil {
      return "", "", nil, e
   }
   fName := fmt.Sprintf(
      "%s - %s.mp3",
      jsonTrack.Results.DATA.SngTitle,
      jsonTrack.Results.DATA.ArtName,
   )
   return downloadURL, fName, client, nil
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
   downloadURL, fName, client, e := getUrl(cfg.ID, client)
   if e != nil {
      log.Fatal(e)
   }
   e = getAudioFile(downloadURL, cfg.ID, fName, client)
   if e != nil {
      log.Fatal(e)
   }
}
