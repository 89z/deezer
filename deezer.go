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
   keyAES = []byte("jo6aey6haid2Teih")
   keyBlowfish = []byte("g4el58wc0zvf9na1")
   iv = []byte{0, 1, 2, 3, 4, 5, 6, 7}
)

// Given SNG_ID and byte slice, decrypt byte slice in place
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

// Given a SNG_ID, return Results.Data
func GetData(sngId, token string) (ResultData, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return ResultData{}, err
   }
   http.DefaultClient.Jar = jar
   val, req := url.Values{}, &http.Request{URL: &API}
   val.Set("api_version", "1.0")
   // GET
   val.Set("api_token", "")
   val.Set("method", "deezer.getUserData")
   req.URL.RawQuery = val.Encode()
   req.Header = http.Header{}
   req.Header.Set("Cookie", "arl=" + token)
   fmt.Println("GET", req.URL)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return ResultData{}, err
   }
   defer resp.Body.Close()
   // JSON
   var check deezCheck
   err = json.NewDecoder(resp.Body).Decode(&check)
   if err != nil {
      return ResultData{}, err
   }
   // POST
   val.Set("api_token", check.Results.CheckForm)
   val.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = val.Encode()
   req.Method = "POST"
   req.Body = io.NopCloser(strings.NewReader(
      fmt.Sprintf(`{"sng_id": "%v"}`, sngId),
   ))
   fmt.Println(req.Method, req.URL)
   resp, err = http.DefaultClient.Do(req)
   if err != nil {
      return ResultData{}, err
   }
   defer resp.Body.Close()
   // JSON
   var track deezTrack
   err = json.NewDecoder(resp.Body).Decode(&track)
   if err != nil {
      return ResultData{}, err
   }
   return track.Results.Data, nil
}

// Given SNG_ID, Results.Data and quality, return audio URL
func GetSource(sngId string, data ResultData, format rune) (string, error) {
   block, err := aes.NewCipher(keyAES)
   if err != nil {
      return "", err
   }
   plain := fmt.Sprint(
      data.MD5Origin, "\xa4",
      string(format), "\xa4",
      sngId, "\xa4",
      data.MediaVersion,
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
      Host: fmt.Sprintf("e-cdns-proxy-%c.dzcdn.net", data.MD5Origin[0]),
      Path: fmt.Sprintf("mobile/1/%x", text),
   }
   return source.String(), nil
}

func md5Hash(s string) string {
   b := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(b),
   )
}

type deezCheck struct {
   Results struct {
      CheckForm string
   }
}

type ResultData struct {
   ArtName      string `json:"ART_NAME"`
   MD5Origin    string `json:"MD5_ORIGIN"`
   MediaVersion string `json:"MEDIA_VERSION"`
   SngTitle     string `json:"SNG_TITLE"`
}

type deezTrack struct {
   Results struct {
      Data ResultData
   }
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
