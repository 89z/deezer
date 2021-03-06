package main

import (
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

func blowfishDecrypt(buf []byte, blowfishKey string) ([]byte, error) {
   if len(buf) % blowfish.BlockSize != 0 {
      return nil, fmt.Errorf("The Buf is not a multiple of 8")
   }
   decrypter, err := blowfish.NewCipher([]byte(blowfishKey))
   if err != nil {
      return nil, err
   }
   cipher.NewCBCDecrypter(
      decrypter, []byte{0, 1, 2, 3, 4, 5, 6, 7},
   ).CryptBlocks(buf, buf)
   return buf, nil
}

func decryptDownload(md5Origin, songID, format, version string) (string, error) {
   block, err := aes.NewCipher([]byte(deezerAES))
   if err != nil {
      return "", err
   }
   src := md5Origin + "\xa4" + format + "\xa4" + songID + "\xa4" + version
   content := []byte(
      md5Hash(src) + "\xa4" + src,
   )
   for len(content) % aes.BlockSize > 0 {
      content = append(content, 0)
   }
   newECBEncrypter(block).CryptBlocks(content, content)
   return fmt.Sprintf(
      "https://e-cdns-proxy-%v.dzcdn.net/mobile/1/%x",
      md5Origin[:1],
      content,
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
      blowfishKey string
      chunkSize = 2048
      idM5 = md5Hash(id)
   )
   for i := 0; i < 16; i++ {
      blowfishKey += string(idM5[i] ^ idM5[i + 16] ^ deezerCBC[i])
   }
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
         chunk, err = blowfishDecrypt(buf, blowfishKey)
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

func getAudioFile(downloadURL, id, FName string) error {
   resp, err := http.Get(downloadURL)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return decryptMedia(resp.Body, id, FName, resp.ContentLength)
}

func getToken(client http.Client) (string, error) {
   req, err := http.NewRequest("GET", APIUrl, nil)
   if err != nil {
      return "", err
   }
   resp, err := client.Do(
      addQs(req),
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

func getUrl(id string, httpClient http.Client) (string, string, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return "", "", err
   }
   jar.SetCookies(
      &cookieURL, []*http.Cookie{{Name: "arl", Value: cfg.UserToken}},
   )
   httpClient.Jar = jar
   APIToken, err := getToken(httpClient)
   if err != nil {
      return "", "", err
   }
   sng := fmt.Sprintf(`{"sng_id": "%v"}`, id)
   req, err := http.NewRequest("POST", APIUrl, strings.NewReader(sng))
   if err != nil {
      return "", "", err
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", APIToken)
   qs.Set("input", "3")
   qs.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = qs.Encode()
   resp, err := httpClient.Do(req)
   if err != nil {
      return "", "", err
   }
   defer resp.Body.Close()
   var jsonTrack DeezTrack
   err = json.NewDecoder(resp.Body).Decode(&jsonTrack)
   if err != nil {
      return "", "", err
   }
   downloadURL, err := decryptDownload(
      jsonTrack.Results.DATA.MD5Origin,
      jsonTrack.Results.DATA.ID.String(),
      "3", // 320
      jsonTrack.Results.DATA.MediaVersion,
   )
   if err != nil {
      return "", "", err
   }
   fName := fmt.Sprintf(
      "%s - %s.mp3",
      jsonTrack.Results.DATA.SngTitle,
      jsonTrack.Results.DATA.ArtName,
   )
   return downloadURL, fName, nil
}

func login() (http.Client, error) {
   req, err := http.NewRequest("POST", APIUrl, nil)
   var httpClient http.Client
   resp, err := httpClient.Do(
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
   req, err = http.NewRequest("POST", LoginURL, strings.NewReader(form.Encode()))
   if err != nil {
      return http.Client{}, err
   }
   resp, err = httpClient.Do(req)
   if err != nil {
      return http.Client{}, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return http.Client{}, fmt.Errorf("StatusCode %v", resp.StatusCode)
   }
   return httpClient, nil
}

func md5Hash(s string) string {
   data := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(data),
   )
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

type ecbEncrypter struct {
   cipher.Block
   blockSize int
}

func newECBEncrypter(b cipher.Block) cipher.BlockMode {
   return ecbEncrypter{
      b, b.BlockSize(),
   }
}

func (x ecbEncrypter) BlockSize() int {
   return x.blockSize
}

func (x ecbEncrypter) CryptBlocks(dst, src []byte) {
   size := x.BlockSize()
   if len(src) % size != 0 {
      panic("crypto/cipher: input not full blocks")
   }
   if len(dst) < len(src) {
      panic("crypto/cipher: output smaller than input")
   }
   for len(src) > 0 {
      x.Encrypt(dst, src[:size])
      src = src[size:]
      dst = dst[size:]
   }
}
