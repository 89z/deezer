package main

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/json"
   "errors"
   "flag"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "os"
)

var cfg = new(Config)

func BFDecrypt(buf []byte, bfKey string) ([]byte, error) {
   decrypter, err := blowfish.NewCipher([]byte(bfKey)) // 8bytes
   if err != nil {
      return nil, err
   }
   IV := []byte{0, 1, 2, 3, 4, 5, 6, 7} //8 bytes
   if len(buf)%blowfish.BlockSize != 0 {
      return nil, errors.New("The Buf is not a multiple of 8")
   }
   cbcDecrypter := cipher.NewCBCDecrypter(decrypter, IV)
   cbcDecrypter.CryptBlocks(buf, buf)
   return buf, nil
}

func ErrorUsage() {
   fmt.Println(`Guide: go-decrypt-deezer [--debug --id --usertoken`)
   fmt.Println(`How Do I Get My UserToken?: https://notabug.org/RemixDevs/DeezloaderRemix/wiki/Login+via+userToken`)
   fmt.Println(`Example: go-decrypt-deezer --id 3135556 --usertoken UserToken_here`)
   flag.PrintDefaults()
   os.Exit(1)
}

func Pad(src []byte) []byte {
   padding := aes.BlockSize - len(src)%aes.BlockSize
   padtext := bytes.Repeat([]byte{byte(padding)}, padding)
   return append(src, padtext...)
}

func init() {
   flag.BoolVar(&cfg.Debug, "debug", false, "Turn on debuging mode.")
   flag.StringVar(&cfg.UserToken, "usertoken", "", "Your Unique User Token")
   flag.StringVar(&cfg.ID, "id", "", "Deezer Track ID")
   flag.Parse()
   if cfg.ID == "" {
      fmt.Println("Error: Must have Deezer Track(Song) ID")
      ErrorUsage()
   }
}

type Config struct {
   Debug     bool
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

type OnError struct {
   Error   error
   Message string
}

type ResultList struct {
   DeezToken      string `json:"checkForm,omitempty"`
   CheckFormLogin string `json:"checkFormLogin,omitempty"`
}

type TrackData struct {
   ID           json.Number `json:"SNG_ID"`
   MD5Origin    string      `json:"MD5_ORIGIN"`
   FileSize320  json.Number `json:"FILESIZE_MP3_320"`
   FileSize256  json.Number `json:"FILESIZE_MP3_256"`
   FileSize128  json.Number `json:"FILESIZE_MP3_128"`
   MediaVersion json.Number `json:"MEDIA_VERSION"`
   SngTitle     string      `json:"SNG_TITLE"`
   ArtName      string      `json:"ART_NAME"`
}

type ecb struct {
   b         cipher.Block
   blockSize int
}

func newECB(b cipher.Block) *ecb {
   return &ecb{
      b, b.BlockSize(),
   }
}

type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
   return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int {
   return x.blockSize
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
   if len(src)%x.blockSize != 0 {
      panic("crypto/cipher: input not full blocks")
   }
   if len(dst) < len(src) {
      panic("crypto/cipher: output smaller than input")
   }
   for len(src) > 0 {
      x.b.Decrypt(dst, src[:x.blockSize])
      src = src[x.blockSize:]
      dst = dst[x.blockSize:]
   }
}

type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
   return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int {
   return x.blockSize
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
   if len(src)%x.blockSize != 0 {
      panic("crypto/cipher: input not full blocks")
   }
   if len(dst) < len(src) {
      panic("crypto/cipher: output smaller than input")
   }
   for len(src) > 0 {
      x.b.Encrypt(dst, src[:x.blockSize])
      src = src[x.blockSize:]
      dst = dst[x.blockSize:]
   }
}
