package main

import (
   "crypto/aes"
   "crypto/cipher"
   "crypto/md5"
   "encoding/json"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "io"
   "strings"
)

const (
   APIURL = "http://www.deezer.com/ajax/gw-light.php"
   deezerAES = "jo6aey6haid2Teih"
   deezerCBC = "g4el58wc0zvf9na1"
)

func decryptDownload(origin, songID, format, version string) (string, error) {
   block, err := aes.NewCipher([]byte(deezerAES))
   if err != nil {
      return "", err
   }
   src := origin + "\xa4" + format + "\xa4" + songID + "\xa4" + version
   content := []byte(
      md5Hash(src) + "\xa4" + src,
   )
   for len(content) % aes.BlockSize > 0 {
      content = append(content, 0)
   }
   newECBEncrypter(block).CryptBlocks(content, content)
   return fmt.Sprintf(
      "https://e-cdns-proxy-%v.dzcdn.net/mobile/1/%x",
      origin[:1],
      content,
   ), nil
}

func decryptMedia(conf config, read *http.Response, write io.Writer) error {
   var (
      bfKey []byte
      trackHash = md5Hash(conf.trackId)
   )
   for n := 0; n < 16; n++ {
      bfKey = append(bfKey, trackHash[n] ^ trackHash[n + 16] ^ deezerCBC[n])
   }
   decrypter, err := blowfish.NewCipher(bfKey)
   if err != nil {
      return err
   }
   var size int64 = 2048
   for n := 0; read.ContentLength > 0; n++ {
      if read.ContentLength < size {
         size = read.ContentLength
      }
      buf := make([]byte, size)
      _, err := read.Body.Read(buf)
      if err != nil {
         return err
      }
      if n % 3 == 0 && read.ContentLength > size {
         cipher.NewCBCDecrypter(
            decrypter, []byte{0, 1, 2, 3, 4, 5, 6, 7},
         ).CryptBlocks(buf, buf)
      }
      _, err = write.Write(buf)
      if err != nil {
         return err
      }
      read.ContentLength -= size
   }
   return nil
}

func getToken(conf config, client http.Client) (string, error) {
   // we must use Request, as cookies are required
   req, err := http.NewRequest("GET", APIURL, nil)
   if err != nil {
      return "", err
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", "null")
   qs.Set("input", "3")
   qs.Set("method", "deezer.getUserData")
   req.URL.RawQuery = qs.Encode()
   req.AddCookie(&http.Cookie{Name: "arl", Value: conf.UserToken})
   resp, err := client.Do(req)
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

func getUrl(conf config) (string, string, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return "", "", err
   }
   client := http.Client{Jar: jar}
   // write cookies
   APIToken, err := getToken(conf, client)
   if err != nil {
      return "", "", err
   }
   sng := fmt.Sprintf(`{"sng_id": "%v"}`, conf.trackId)
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", APIToken)
   qs.Set("input", "3")
   qs.Set("method", "deezer.pageTrack")
   // we must use Request, as cookies are required
   req, err := http.NewRequest("POST", APIURL, strings.NewReader(sng))
   if err != nil {
      return "", "", err
   }
   req.URL.RawQuery = qs.Encode()
   // read cookies
   resp, err := client.Do(req)
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

func md5Hash(s string) string {
   data := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(data),
   )
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

type config struct {
   UserToken string
   trackId string
}

type ecbEncrypter struct {
   cipher.Block
}

func newECBEncrypter(b cipher.Block) cipher.BlockMode {
   return ecbEncrypter{b}
}

func (x ecbEncrypter) BlockSize() int {
   return x.Block.BlockSize()
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
