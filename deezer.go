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
   deezer320 = '3'
   deezerFlac = '9'
)

var (
   deezerAES = []byte("jo6aey6haid2Teih")
   deezerCBC = []byte("g4el58wc0zvf9na1")
   deezerIv = []byte{0, 1, 2, 3, 4, 5, 6, 7}
)

var deezerAPI = url.URL{
   Scheme: "http", Host: "www.deezer.com", Path: "/ajax/gw-light.php",
}

func decryptAudio(trackId string, from io.Reader, to io.Writer) error {
   var (
      bfKey []byte
      trackHash = md5Hash(trackId)
   )
   for n := 0; n < 16; n++ {
      bfKey = append(bfKey, trackHash[n] ^ trackHash[n + 16] ^ deezerCBC[n])
   }
   block, err := blowfish.NewCipher(bfKey)
   if err != nil {
      return err
   }
   var d int
   for {
      size := 2048
      text := make([]byte, size)
      n, err := from.Read(text)
      if err == io.EOF {
         break
      } else if err != nil {
         return err
      }
      if n < size {
         text = text[:n]
      }
      if d % 3 == 0 && n == size {
         cipher.NewCBCDecrypter(block, deezerIv).CryptBlocks(text, text)
      }
      _, err = to.Write(text)
      if err != nil {
         return err
      }
      d++
   }
   return nil
}

func getData(token, trackId string) (deezData, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return deezData{}, err
   }
   http.DefaultClient.Jar = jar
   val, req := url.Values{}, &http.Request{URL: &deezerAPI}
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
      return deezData{}, err
   }
   defer resp.Body.Close()
   // JSON
   var check deezCheck
   err = json.NewDecoder(resp.Body).Decode(&check)
   if err != nil {
      return deezData{}, err
   }
   // POST
   val.Set("api_token", check.Results.CheckForm)
   val.Set("method", "deezer.pageTrack")
   req.URL.RawQuery = val.Encode()
   req.Method = "POST"
   req.Body = io.NopCloser(strings.NewReader(
      fmt.Sprintf(`{"sng_id": "%v"}`, trackId),
   ))
   fmt.Println(req.Method, req.URL)
   resp, err = http.DefaultClient.Do(req)
   if err != nil {
      return deezData{}, err
   }
   defer resp.Body.Close()
   // JSON
   var track deezTrack
   err = json.NewDecoder(resp.Body).Decode(&track)
   if err != nil {
      return deezData{}, err
   }
   return track.Results.Data, nil
}

func getSource(data deezData, format rune) (string, error) {
   block, err := aes.NewCipher(deezerAES)
   if err != nil {
      return "", err
   }
   plain := fmt.Sprint(
      data.MD5Origin, "\xa4",
      string(format), "\xa4",
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
      x.Encrypt(dst, src)
      src, dst = src[size:], dst[size:]
   }
}
