package main

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "crypto/md5"
   "encoding/json"
   "fmt"
   "golang.org/x/crypto/blowfish"
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
   APIUrl = "http://www.deezer.com/ajax/gw-light.php"
   LoginURL = "https://www.deezer.com/ajax/action.php"
   deezerAES = "jo6aey6haid2Teih"
   deezerCBC = "g4el58wc0zvf9na1"
)

var (
   cfg = new(Config)
   cookieURL = url.URL{Scheme: "https", Host: "www.deezer.com"}
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

func blowfishDecrypt(buf []byte, bfKey string) ([]byte, error) {
   if len(buf) % blowfish.BlockSize != 0 {
      return nil, fmt.Errorf("The Buf is not a multiple of 8")
   }
   decrypter, err := blowfish.NewCipher([]byte(bfKey))
   if err != nil {
      return nil, err
   }
   cipher.NewCBCDecrypter(
      decrypter, []byte{0, 1, 2, 3, 4, 5, 6, 7},
   ).CryptBlocks(buf, buf)
   return buf, nil
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
   block, err := aes.NewCipher([]byte(deezerAES))
   if err != nil {
      return "", err
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
   var (
      bfKey = getBlowFishKey(id)
      chunkSize = 2048
      destBuffer bytes.Buffer
      errc = make(chan error)
      wg sync.WaitGroup
   )
   for position, i := 0, 0; position < int(streamLen); position, i = position+chunkSize, i+1 {
      func(i, position int, streamLen int64, stream io.Reader) {
         if (int(streamLen) - position) >= 2048 {
            chunkSize = 2048
         } else {
            chunkSize = int(streamLen) - position
         }
         buf := make([]byte, chunkSize) // The "chunk" of data
         if _, err := io.ReadFull(stream, buf); err != nil {
            errc <- fmt.Errorf("loop %v %v", i, err)
         }
         var (
            chunkString []byte
            err error
         )
         if i % 3 > 0 || chunkSize < 2048 {
            chunkString = buf
         } else { //Decrypt and then write to destBuffer
            chunkString, err = blowfishDecrypt(buf, bfKey)
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
      case err := <-errc:
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
   md5Sum := md5.Sum([]byte(id))
   idM5 := fmt.Sprintf("%x", md5Sum)
   var BFKey string
   for i := 0; i < 16; i++ {
      BFKey += string(idM5[i] ^ idM5[i + 16] ^ deezerCBC[i])
   }
   return BFKey
}

func getToken(client *http.Client) (string, error) {
   Deez := &DeezStruct{}
   reqs, err := newRequest(APIUrl, "GET", nil)
   if err != nil {
      return "", err
   }
   reqs = addQs(reqs)
   resp, err := client.Do(reqs)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   err = json.Unmarshal(body, &Deez)
   if err != nil {
      return "", err
   }
   APIToken := Deez.Results.DeezToken
   return APIToken, nil
}

func getUrl(id string, client *http.Client) (string, string, *http.Client, error) {
   APIToken, err := getToken(client)
   if err != nil {
      return "", "", nil, err
   }
   jsonPrep := `{"sng_id":"` + id + `"}`
   jsonStr := []byte(jsonPrep)
   req, err := newRequest(APIUrl, "POST", jsonStr)
   if err != nil {
      return "", "", nil, err
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", APIToken)
   qs.Set("input", "3")
   qs.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = qs.Encode()
   resp, err := client.Do(req)
   if err != nil {
      return "", "", nil, err
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return "", "", nil, err
   }
   defer resp.Body.Close()
   var jsonTrack DeezTrack
   err = json.Unmarshal(body, &jsonTrack)
   if err != nil {
      return "", "", nil, err
   }
   downloadURL, err := decryptDownload(
      jsonTrack.Results.DATA.MD5Origin,
      jsonTrack.Results.DATA.ID.String(),
      "3", // 320
      jsonTrack.Results.DATA.MediaVersion,
   )
   if err != nil {
      return "", "", nil, err
   }
   fName := fmt.Sprintf(
      "%s - %s.mp3",
      jsonTrack.Results.DATA.SngTitle,
      jsonTrack.Results.DATA.ArtName,
   )
   return downloadURL, fName, client, nil
}

func login() (*http.Client, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   client := &http.Client{Jar: jar}
   Deez := &DeezStruct{}
   req, err := newRequest(APIUrl, "POST", nil)
   req = addQs(req)
   resp, err := client.Do(req)
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   err = json.Unmarshal(body, &Deez)
   if err != nil {
      return nil, err
   }
   resp.Body.Close()
   form := url.Values{}
   form.Set("type", "login")
   form.Set("checkFormLogin", Deez.Results.CheckFormLogin)
   req, err = newRequest(LoginURL, "POST", form.Encode())
   if err != nil {
      return nil, err
   }
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
      cookies := []*http.Cookie{{
         Name: "arl",
         Value: cfg.UserToken,
      }}
      client.Jar.SetCookies(&cookieURL, cookies)
      return client, nil
   }
   return nil, fmt.Errorf("StatusCode %v %v", resp.StatusCode, err)
}

func newRequest(enPoint, method string, bodyEntity interface{}) (*http.Request, error) {
   var (
      err error
      req *http.Request
   )
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
   ArtName      string `json:"ART_NAME"`
   ID           json.Number `json:"SNG_ID"`
   MD5Origin    string `json:"MD5_ORIGIN"`
   MediaVersion string `json:"MEDIA_VERSION"`
   SngTitle     string `json:"SNG_TITLE"`
}

type ecb struct {
   cipher.Block
}

func newECB(b cipher.Block) ecb {
   return ecb{b}
}

func (x ecb) cryptBlocks(dst, src []byte) error {
   size := x.BlockSize()
   if len(src) % size != 0 {
      return fmt.Errorf("crypto/cipher: input not full blocks")
   }
   if len(dst) < len(src) {
      return fmt.Errorf("crypto/cipher: output smaller than input")
   }
   for len(src) > 0 {
      x.Encrypt(dst, src[:size])
      src = src[size:]
      dst = dst[size:]
   }
   return nil
}
