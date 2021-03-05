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
   "strings"
   "sync"
)

const (
   // APIUrl is the deezer API
   APIUrl = "http://www.deezer.com/ajax/gw-light.php"
   // LoginURL is the API for deezer login
   LoginURL = "https://www.deezer.com/ajax/action.php"
   // deezer domain for cookie check
   Domain = "https://www.deezer.com"
)

var (
   cfg = new(Config)
   deezerKey = []byte("jo6aey6haid2Teih")
)


func addQs(req *http.Request) *http.Request {
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", "null")
   qs.Set("input", "3")
   qs.Set("method", "deezer.getUserData")
   req.URL.RawQuery = qs.Encode()
   return req
}

func decryptDownload(md5Origin, songID, format, version string) (string, error) {
   urlPart := md5Origin + "¤" + format + "¤" + songID + "¤" + version
   data := bytes.Replace([]byte(urlPart), []byte("¤"), []byte{164}, -1)
   md5SumVal := fmt.Sprintf("%x", md5.Sum(data))
   urlPart = md5SumVal + "¤" + urlPart + "¤"
   src := bytes.Replace([]byte(urlPart), []byte("¤"), []byte{164}, -1)
   padding := aes.BlockSize - len(src) % aes.BlockSize
   padtext := bytes.Repeat([]byte{byte(padding)}, padding)
   plaintext := append(src, padtext...)
   block, e := aes.NewCipher(deezerKey)
   if e != nil {
      return "", e
   }
   encryptText := make([]byte, len(plaintext))
   newECB(block).cryptBlocks(encryptText, plaintext)
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
   var e error
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
         if _, e = io.ReadFull(stream, buf); e != nil {
            errc <- fmt.Errorf("loop %v %v", i, e)
         }
         if i % 3 > 0 || chunkSize < 2048 {
            chunkString = buf
         } else { //Decrypt and then write to destBuffer
            chunkString, e = blowfishDecrypt(buf, bfKey)
            if e != nil {
               errc <- fmt.Errorf("loop %v %v", i, e)
            }
         }
         if _, e := destBuffer.Write(chunkString); e != nil {
            errc <- fmt.Errorf("loop %v %v", i, e)
         }
      }(i, position, streamLen, stream)
   }
   for {
      select {
      case e = <-errc:
         return e
      default:
         wg.Wait()
         NameWithoutSlash := strings.ReplaceAll(FName, "/", "∕")
         out, e := os.Create(NameWithoutSlash)
         if e != nil {
            return e
         }
         _, e = destBuffer.WriteTo(out)
         if e != nil {
            return e
         }
         return nil
      }
   }
}

func getAudioFile(downloadURL, id, FName string, client *http.Client) error {
   req, e := newRequest(downloadURL, "GET", nil)
   if e != nil {
      return e
   }
   resp, e := client.Do(req)
   if e != nil {
      return e
   }
   e = decryptMedia(resp.Body, id, FName, resp.ContentLength)
   if e != nil {
      return e
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
   reqs, e := newRequest(APIUrl, "GET", nil)
   if e != nil {
      return "", e
   }
   reqs = addQs(reqs)
   resp, e := client.Do(reqs)
   if e != nil {
      return "", e
   }
   defer resp.Body.Close()
   body, _ := ioutil.ReadAll(resp.Body)
   e = json.Unmarshal(body, &Deez)
   if e != nil {
      return "", e
   }
   APIToken := Deez.Results.DeezToken
   return APIToken, nil
}

func getUrl(id string, client *http.Client) (string, string, *http.Client, error) {
   jsonTrack := &DeezTrack{}
   APIToken, _ := getToken(client)
   jsonPrep := `{"sng_id":"` + id + `"}`
   jsonStr := []byte(jsonPrep)
   req, e := newRequest(APIUrl, "POST", jsonStr)
   if e != nil {
      return "", "", nil, e
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", APIToken)
   qs.Set("input", "3")
   qs.Set("method", "deezer.pageTrack")
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
   downloadURL, e := decryptDownload(md5Origin, songID, format, mediaVersion)
   if e != nil {
      return "", "", nil, e
   }
   return downloadURL, FName, client, nil
}

func login() (*http.Client, error) {
   CookieJar, _ := cookiejar.New(nil)
   client := &http.Client{Jar: CookieJar}
   Deez := &DeezStruct{}
   req, e := newRequest(APIUrl, "POST", nil)
   req = addQs(req)
   resp, e := client.Do(req)
   body, _ := ioutil.ReadAll(resp.Body)
   e = json.Unmarshal(body, &Deez)
   if e != nil {
      return nil, e
   }
   CookieURL, _ := url.Parse(Domain)
   resp.Body.Close()
   form := url.Values{}
   form.Set("type", "login")
   form.Set("checkFormLogin", Deez.Results.CheckFormLogin)
   req, e = newRequest(LoginURL, "POST", form.Encode())
   if e != nil {
      return nil, e
   }
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
      cookies := []*http.Cookie{{
         Name: "arl",
         Value: cfg.UserToken,
      }}
      client.Jar.SetCookies(CookieURL, cookies)
      return client, nil
   }
   return nil, fmt.Errorf("StatusCode %v %v", resp.StatusCode, e)
}

func newRequest(enPoint, method string, bodyEntity interface{}) (*http.Request, error) {
   var req *http.Request
   var e error
   switch val := bodyEntity.(type) {
   case []byte:
      req, e = http.NewRequest(method, enPoint, bytes.NewBuffer(val))
   case string:
      req, e = http.NewRequest(method, enPoint, strings.NewReader(val))
   default:
      req, e = http.NewRequest(method, enPoint, nil)
   }
   if bodyEntity == nil {
      req, e = http.NewRequest(method, enPoint, nil)
   }
   if e != nil {
      return nil, e
   }
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
