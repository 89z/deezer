package main

import (
   "bytes"
   "crypto/aes"
   "crypto/md5"
   "encoding/json"
   "fmt"
   "io"
   "io/ioutil"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "strconv"
   "strings"
   "sync"
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

func addQs(req *http.Request, args ...string) *http.Request {
   qs := url.Values{}
   qs.Add("api_version", "1.0")
   qs.Add("api_token", args[0]) //args[0] always token
   qs.Add("input", "3")
   qs.Add("method", args[1]) //args[1] always method
   req.URL.RawQuery = qs.Encode()
   return req
}

func decryptDownload(md5Origin, songID, format, mediaVersion string) (string, error) {
   urlPart := md5Origin + "¤" + format + "¤" + songID + "¤" + mediaVersion
   data := bytes.Replace([]byte(urlPart), []byte("¤"), []byte{164}, -1)
   md5SumVal := fmt.Sprintf("%x", md5.Sum(data))
   urlPart = md5SumVal + "¤" + urlPart + "¤"
   key := []byte("jo6aey6haid2Teih")
   plaintext := Pad(bytes.Replace([]byte(urlPart), []byte("¤"), []byte{164}, -1))
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   encryptText := make([]byte, len(plaintext))
   mode := NewECBEncrypter(block) // return ECB encryptor
   mode.CryptBlocks(encryptText, plaintext)
   return fmt.Sprintf(
      "https://e-cdns-proxy-%v.dzcdn.net/mobile/1/%x",
      md5Origin[:1],
      encryptText,
   ), nil
}

func decryptMedia(stream io.Reader, id, FName string, streamLen int64) error {
   var wg sync.WaitGroup
   chunkSize := 2048
   bfKey := getBlowFishKey(id)
   errc := make(chan error)
   var err error
   var destBuffer bytes.Buffer // final Product
   for position, i := 0, 0; position < int(streamLen); position, i = position+chunkSize, i+1 {
      func(i, position int, streamLen int64, stream io.Reader) {
         var chunkString []byte
         if (int(streamLen) - position) >= 2048 {
            chunkSize = 2048
         } else {
            chunkSize = int(streamLen) - position
         }
         buf := make([]byte, chunkSize) // The "chunk" of data
         if _, err = io.ReadFull(stream, buf); err != nil {
            errc <- fmt.Errorf("loop %v %v", i, err)
         }
         if i % 3 > 0 || chunkSize < 2048 {
            chunkString = buf
         } else { //Decrypt and then write to destBuffer
            chunkString, err = BFDecrypt(buf, bfKey)
            if err != nil {
               errc <- fmt.Errorf("loop %v %v", i, err)
            }
         }
         if _, err := destBuffer.Write(chunkString); err != nil {
            errc <- fmt.Errorf("loop %v %v", i, err)
         }
      }(i, position, streamLen, stream)
   }
   for {
      select {
      case err = <-errc:
         return err
      default:
         wg.Wait()
         NameWithoutSlash := strings.ReplaceAll(FName, "/", "∕")
         out, err := os.Create(NameWithoutSlash)
         if err != nil {
            return err
         }
         _, err = destBuffer.WriteTo(out)
         if err != nil {
            return err
         }
         return nil
      }
   }
}

func getAudioFile(downloadURL, id, FName string, client *http.Client) error {
   req, err := newRequest(downloadURL, "GET", nil)
   if err != nil {
      return err
   }
   resp, err := client.Do(req)
   if err != nil {
      return err
   }
   err = decryptMedia(resp.Body, id, FName, resp.ContentLength)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return nil
}

func getBlowFishKey(id string) string {
   Secret := "g4el58wc0zvf9na1"
   md5Sum := md5.Sum([]byte(id))
   idM5 := fmt.Sprintf("%x", md5Sum)
   var BFKey string
   for i := 0; i < 16; i++ {
      BFKey += string(idM5[i] ^ idM5[i + 16] ^ Secret[i])
   }
   return BFKey
}

func getToken(client *http.Client) (string, error) {
   Deez := &DeezStruct{}
   args := []string{"null", "deezer.getUserData"}
   reqs, err := newRequest(APIUrl, "GET", nil)
   if err != nil {
   return "", err
   }
   reqs = addQs(reqs, args...)
   resp, err := client.Do(reqs)
   if err != nil {
   return "", err
   }
   defer resp.Body.Close()
   body, _ := ioutil.ReadAll(resp.Body)
   err = json.Unmarshal(body, &Deez)
   if err != nil {
   return "", err
   }
   APIToken := Deez.Results.DeezToken
   return APIToken, nil
}

func getUrl(id string, client *http.Client) (string, string, *http.Client, error) {
   jsonTrack := &DeezTrack{}
   APIToken, _ := getToken(client)
   jsonPrep := `{"sng_id":"` + id + `"}`
   jsonStr := []byte(jsonPrep)
   req, err := newRequest(APIUrl, "POST", jsonStr)
   if err != nil {
      return "", "", nil, err
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
   err = json.Unmarshal(body, &jsonTrack)
   if err != nil {
      return "", "", nil, err
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
   downloadURL, err := decryptDownload(md5Origin, songID, format, mediaVersion)
   if err != nil {
      return "", "", nil, err
   }
   return downloadURL, FName, client, nil
}

func login() (*http.Client, error) {
   CookieJar, _ := cookiejar.New(nil)
   client := &http.Client{Jar: CookieJar}
   Deez := &DeezStruct{}
   req, err := newRequest(APIUrl, "POST", nil)
   args := []string{"null", "deezer.getUserData"}
   req = addQs(req, args...)
   resp, err := client.Do(req)
   body, _ := ioutil.ReadAll(resp.Body)
   err = json.Unmarshal(body, &Deez)
   if err != nil {
      return nil, err
   }
   CookieURL, _ := url.Parse(Domain)
   resp.Body.Close()
   form := url.Values{}
   form.Add("type", "login")
   form.Add("checkFormLogin", Deez.Results.CheckFormLogin)
   req, err = newRequest(LoginURL, "POST", form.Encode())
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-type", "application/x-www-form-urlencoded")
   req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
   resp, err = client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err = ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode == 200 {
      addCookies(client, CookieURL)
      return client, nil
   }
   return nil, fmt.Errorf("StatusCode %v %v", resp.StatusCode, err)
}

func newRequest(enPoint, method string, bodyEntity interface{}) (*http.Request, error) {
   var req *http.Request
   var err error
   switch val := bodyEntity.(type) {
   case []byte:
      req, err = http.NewRequest(method, enPoint, bytes.NewBuffer(val))
   case string:
      req, err = http.NewRequest(method, enPoint, strings.NewReader(val))
   default:
      req, err = http.NewRequest(method, enPoint, nil)
   }
   if bodyEntity == nil {
      req, err = http.NewRequest(method, enPoint, nil)
   }
   if err != nil {
      return nil, err
   }
   req.Header.Add(
      "User-Agent",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36",
   )
   req.Header.Add("Content-Language", "en-US")
   req.Header.Add("Cache-Control", "max-age=0")
   req.Header.Add("Accept", "*/*")
   req.Header.Add("Accept-Charset", "utf-8,ISO-8859-1;q=0.7,*;q=0.3")
   req.Header.Add("Accept-Language", "de-DE,de;q=0.8,en-US;q=0.6,en;q=0.4")
   req.Header.Add("Content-type", "application/json")
   return req, nil
}

type Config struct {
   ID        string
   UserToken string
}

type Data struct {
   DATA *TrackData `json:"DATA"`
}

type DeezStruct struct {
   Error   []string    `json:"error,omitempty"`
   Results *ResultList `json:"results,omitempty"`
}

type DeezTrack struct {
   Error   []string `json:"error,omitempty"`
   Results *Data    `json:"results,omitempty"`
}

type ResultList struct {
   CheckFormLogin string `json:"checkFormLogin,omitempty"`
   DeezToken      string `json:"checkForm,omitempty"`
}

type TrackData struct {
   ArtName      string      `json:"ART_NAME"`
   FileSize128  json.Number `json:"FILESIZE_MP3_128"`
   FileSize256  json.Number `json:"FILESIZE_MP3_256"`
   FileSize320  json.Number `json:"FILESIZE_MP3_320"`
   ID           json.Number `json:"SNG_ID"`
   MD5Origin    string      `json:"MD5_ORIGIN"`
   MediaVersion json.Number `json:"MEDIA_VERSION"`
   SngTitle     string      `json:"SNG_TITLE"`
}
