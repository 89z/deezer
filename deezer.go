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

const deezerAPI = "http://www.deezer.com/ajax/gw-light.php"

const (
   deezer320 = '3'
   deezerFlac = '9'
)

var (
   deezerAES = []byte("jo6aey6haid2Teih")
   deezerCBC = []byte("g4el58wc0zvf9na1")
   deezerIv = []byte{0, 1, 2, 3, 4, 5, 6, 7}
)

func getData(conf config) (deezData, error) {
   // client
   jar, err := cookiejar.New(nil)
   if err != nil {
      return deezData{}, err
   }
   client := http.Client{Jar: jar}
   // we must use Request, as cookies are required
   req, err := http.NewRequest("GET", deezerAPI, nil)
   if err != nil {
      return deezData{}, err
   }
   qs := url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", "null")
   qs.Set("method", "deezer.getUserData")
   req.URL.RawQuery = qs.Encode()
   req.AddCookie(&http.Cookie{Name: "arl", Value: conf.userToken})
   // send
   resp, err := client.Do(req)
   if err != nil {
      return deezData{}, err
   }
   defer resp.Body.Close()
   var check deezCheck
   err = json.NewDecoder(resp.Body).Decode(&check)
   if err != nil {
      return deezData{}, err
   }
   sng := fmt.Sprintf(`{"sng_id": "%v"}`, conf.trackId)
   // we must use Request, as cookies are required
   req, err = http.NewRequest("POST", deezerAPI, strings.NewReader(sng))
   if err != nil {
      return deezData{}, err
   }
   qs = url.Values{}
   qs.Set("api_version", "1.0")
   qs.Set("api_token", check.Results.CheckForm)
   qs.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = qs.Encode()
   // read cookies
   resp, err = client.Do(req)
   if err != nil {
      return deezData{}, err
   }
   defer resp.Body.Close()
   var track deezTrack
   err = json.NewDecoder(resp.Body).Decode(&track)
   if err != nil {
      return deezData{}, err
   }
   return track.Results.Data, nil
}

func decryptDownload(data deezData) (string, error) {
   block, err := aes.NewCipher(deezerAES)
   if err != nil {
      return "", err
   }
   plain := fmt.Sprint(
      data.MD5Origin, "\xa4",
      string(deezer320), "\xa4",
      data.SngId, "\xa4",
      data.MediaVersion,
   )
   text := []byte(
      md5Hash(plain) + "\xa4" + plain,
   )
   for len(text) % aes.BlockSize > 0 {
      text = append(text, 0)
   }
   newECBEncrypter(block).CryptBlocks(text, text)
   return fmt.Sprintf(
      "https://e-cdns-proxy-%c.dzcdn.net/mobile/1/%x", data.MD5Origin[0], text,
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
   block, err := blowfish.NewCipher(bfKey)
   if err != nil {
      return err
   }
   var size int64 = 2048
   for n := 0; read.ContentLength > 0; n++ {
      if read.ContentLength < size {
         size = read.ContentLength
      }
      text := make([]byte, size)
      _, err := read.Body.Read(text)
      if err != nil {
         return err
      }
      if n % 3 == 0 && read.ContentLength > size {
         cipher.NewCBCDecrypter(block, deezerIv).CryptBlocks(text, text)
      }
      _, err = write.Write(text)
      if err != nil {
         return err
      }
      read.ContentLength -= size
   }
   return nil
}


func md5Hash(s string) string {
   b := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(b),
   )
}

type config struct {
   userToken string
   trackId string
}

type deezCheck struct {
   Results struct {
      CheckForm string
   }
}

type deezData struct {
   ArtName      string `json:"ART_NAME"`
   MD5Origin    string `json:"MD5_ORIGIN"`
   MediaVersion string `json:"MEDIA_VERSION"`
   SngId        string `json:"SNG_ID"`
   SngTitle     string `json:"SNG_TITLE"`
}

type deezTrack struct {
   Results struct {
      Data deezData
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
      x.Encrypt(dst, src[:size])
      src = src[size:]
      dst = dst[size:]
   }
}
