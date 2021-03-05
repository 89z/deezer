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
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "strings"
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
   urlPart := md5Origin + "\xa4" + format + "\xa4" + songID + "\xa4" + version
   src := []byte(
      md5Hash(urlPart) + "\xa4" + urlPart + "\xa4",
   )
   padtext := bytes.Repeat(
      []byte{0}, aes.BlockSize - len(src) % aes.BlockSize,
   )
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
   create, err := os.Create(
      strings.ReplaceAll(FName, "/", " "),
   )
   if err != nil {
      return err
   }
   defer create.Close()
   var (
      bfKey = getBlowFishKey(id)
      chunkSize = 2048
   )
   for pos, i := 0, 0; pos < int(streamLen); pos, i = pos + chunkSize, i + 1 {
      if (int(streamLen) - pos) >= 2048 {
         chunkSize = 2048
      } else {
         chunkSize = int(streamLen) - pos
      }
      buf := make([]byte, chunkSize)
      _, err := io.ReadFull(stream, buf)
      if err != nil {
         return fmt.Errorf("loop %v %v", i, err)
      }
      var chunk []byte
      if i % 3 > 0 || chunkSize < 2048 {
         chunk = buf
      } else {
         chunk, err = blowfishDecrypt(buf, bfKey)
         if err != nil {
            return fmt.Errorf("loop %v %v", i, err)
         }
      }
      _, err = create.Write(chunk)
      if err != nil {
         return fmt.Errorf("loop %v %v", i, err)
      }
   }
   return nil
}

func getAudioFile(downloadURL, id, FName string, client http.Client) error {
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
   var (
      BFKey string
      idM5 = md5Hash(id)
   )
   for i := 0; i < 16; i++ {
      BFKey += string(idM5[i] ^ idM5[i + 16] ^ deezerCBC[i])
   }
   return BFKey
}

func getToken(client http.Client) (string, error) {
   reqs, err := newRequest(APIUrl, "GET", nil)
   if err != nil {
      return "", err
   }
   resp, err := client.Do(
      addQs(reqs),
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var deez DeezStruct
   err = json.NewDecoder(resp.Body).Decode(&deez)
   if err != nil {
      return "", err
   }
   return deez.Results.DeezToken, nil
}

func getUrl(id string, client http.Client) (string, string, http.Client, error) {
   APIToken, err := getToken(client)
   if err != nil {
      return "", "", http.Client{}, err
   }
   sng := []byte(
      fmt.Sprintf(`{"sng_id": "%v"}`, id)
   )
   req, err := newRequest(APIUrl, "POST", sng)
   if err != nil {
      return "", "", http.Client{}, err
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", APIToken)
   qs.Set("input", "3")
   qs.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = qs.Encode()
   resp, err := client.Do(req)
   if err != nil {
      return "", "", http.Client{}, err
   }
   defer resp.Body.Close()
   var jsonTrack DeezTrack
   err = json.NewDecoder(resp.Body).Decode(&jsonTrack)
   if err != nil {
      return "", "", http.Client{}, err
   }
   downloadURL, err := decryptDownload(
      jsonTrack.Results.DATA.MD5Origin,
      jsonTrack.Results.DATA.ID.String(),
      "3", // 320
      jsonTrack.Results.DATA.MediaVersion,
   )
   if err != nil {
      return "", "", http.Client{}, err
   }
   fName := fmt.Sprintf(
      "%s - %s.mp3",
      jsonTrack.Results.DATA.SngTitle,
      jsonTrack.Results.DATA.ArtName,
   )
   return downloadURL, fName, client, nil
}

func login() (http.Client, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return http.Client{}, err
   }
   client := http.Client{Jar: jar}
   req, err := newRequest(APIUrl, "POST", nil)
   resp, err := client.Do(
      addQs(req),
   )
   if err != nil {
      return http.Client{}, err
   }
   defer resp.Body.Close()
   var deez DeezStruct
   err = json.NewDecoder(resp.Body).Decode(&deez)
   if err != nil {
      return http.Client{}, err
   }
   form := url.Values{}
   form.Set("type", "login")
   form.Set("checkFormLogin", deez.Results.CheckFormLogin)
   req, err = newRequest(LoginURL, "POST", form.Encode())
   if err != nil {
      return http.Client{}, err
   }
   resp, err = client.Do(req)
   if err != nil {
      return http.Client{}, err
   }
   defer resp.Body.Close()
   if resp.StatusCode == 200 {
      cookies := []*http.Cookie{{
         Name: "arl",
         Value: cfg.UserToken,
      }}
      client.Jar.SetCookies(&cookieURL, cookies)
      return client, nil
   }
   return http.Client{}, fmt.Errorf("StatusCode %v %v", resp.StatusCode, err)
}

func md5Hash(s string) string {
   data := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(data),
   )
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
