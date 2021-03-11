// Deezer
package deezer

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
   "strings"
)

const (
   FLAC = '9'
   MP3_320 = '3'
)

var API = url.URL{
   Scheme: "http", Host: "www.deezer.com", Path: "/ajax/gw-light.php",
}

var (
   iv = []byte{0, 1, 2, 3, 4, 5, 6, 7}
   keyAES = []byte("jo6aey6haid2Teih")
   keyBlowfish = []byte("g4el58wc0zvf9na1")
)

// Given SNG_ID and byte slice, decrypt byte slice in place.
func Decrypt(sngId string, data []byte) error {
   hash := md5Hash(sngId)
   for n := range keyBlowfish {
      keyBlowfish[n] ^= hash[n] ^ hash[n + len(keyBlowfish)]
   }
   block, err := blowfish.NewCipher(keyBlowfish)
   if err != nil {
      return err
   }
   size := 2048
   for pos := 0; len(data) - pos >= size; pos += size {
      if pos / size % 3 == 0 {
         text := data[pos : pos + size]
         cipher.NewCBCDecrypter(block, iv).CryptBlocks(text, text)
      }
   }
   return nil
}

func logInfo(s string, a ...interface{}) {
   fmt.Print("\x1b[30;106m ", s, " \x1b[m ")
   fmt.Println(a...)
}

func md5Hash(s string) string {
   b := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(b),
   )
}

type Track struct {
   ArtName      string `json:"ART_NAME"`
   MD5Origin    string `json:"MD5_ORIGIN"`
   MediaVersion string `json:"MEDIA_VERSION"`
   SngTitle     string `json:"SNG_TITLE"`
}

// Given a SNG_ID and arl strings, make a "deezer.pageTrack" request and return
// the result.
func NewTrack(sngId, arl string) (Track, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return Track{}, err
   }
   http.DefaultClient.Jar = jar
   val, req := url.Values{}, &http.Request{URL: &API}
   val.Set("api_version", "1.0")
   // GET
   val.Set("api_token", "")
   val.Set("method", "deezer.getUserData")
   req.URL.RawQuery = val.Encode()
   req.Header = http.Header{}
   req.Header.Set("Cookie", "arl=" + arl)
   logInfo("Get", req.URL)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return Track{}, err
   }
   defer resp.Body.Close()
   // JSON
   var data userData
   err = json.NewDecoder(resp.Body).Decode(&data)
   if err != nil {
      return Track{}, err
   }
   // POST
   val.Set("api_token", data.Results.CheckForm)
   val.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = val.Encode()
   req.Method = "POST"
   req.Body = io.NopCloser(strings.NewReader(
      fmt.Sprintf(`{"sng_id": "%v"}`, sngId),
   ))
   logInfo("Post", req.URL)
   resp, err = http.DefaultClient.Do(req)
   if err != nil {
      return Track{}, err
   }
   defer resp.Body.Close()
   // JSON
   var page pageTrack
   err = json.NewDecoder(resp.Body).Decode(&page)
   if err != nil {
      return Track{}, err
   }
   return page.Results.Data, nil
}

// Given SNG_ID and file format, return audio URL.
func (t Track) Source(sngId string, format rune) (string, error) {
   block, err := aes.NewCipher(keyAES)
   if err != nil {
      return "", err
   }
   plain := fmt.Sprint(
      t.MD5Origin, "\xa4", string(format), "\xa4", sngId, "\xa4", t.MediaVersion,
   )
   text := []byte(
      md5Hash(plain) + "\xa4" + plain,
   )
   for len(text) % aes.BlockSize > 0 {
      text = append(text, 0)
   }
   newECBEncrypter(block).CryptBlocks(text, text)
   source := url.URL{
      Scheme: "https",
      Host: fmt.Sprintf("e-cdns-proxy-%c.dzcdn.net", t.MD5Origin[0]),
      Path: fmt.Sprintf("mobile/1/%x", text),
   }
   return source.String(), nil
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
      x.Encrypt(dst, src)
      src, dst = src[size:], dst[size:]
   }
}

type pageTrack struct {
   Results struct {
      Data Track
   }
}

type userData struct {
   Results struct {
      CheckForm string
   }
}
