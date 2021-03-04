package main

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/json"
   "errors"
   "golang.org/x/crypto/blowfish"
)

func BFDecrypt(buf []byte, bfKey string) ([]byte, error) {
   decrypter, err := blowfish.NewCipher([]byte(bfKey)) // 8 bytes
   if err != nil {
      return nil, err
   }
   IV := []byte{0, 1, 2, 3, 4, 5, 6, 7} // 8 bytes
   if len(buf)%blowfish.BlockSize != 0 {
      return nil, errors.New("The Buf is not a multiple of 8")
   }
   cbcDecrypter := cipher.NewCBCDecrypter(decrypter, IV)
   cbcDecrypter.CryptBlocks(buf, buf)
   return buf, nil
}

func Pad(src []byte) []byte {
   padding := aes.BlockSize - len(src)%aes.BlockSize
   padtext := bytes.Repeat([]byte{byte(padding)}, padding)
   return append(src, padtext...)
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
   DeezToken      string `json:"checkForm,omitempty"`
   CheckFormLogin string `json:"checkFormLogin,omitempty"`
}

type TrackData struct {
   ArtName      string `json:"ART_NAME"`
   MD5Origin    string `json:"MD5_ORIGIN"`
   SngTitle     string `json:"SNG_TITLE"`
   FileSize256  int `json:"FILESIZE_MP3_256"`
   FileSize320  int `json:"FILESIZE_MP3_320"`
   ID           json.Number `json:"SNG_ID"`
   MediaVersion json.Number `json:"MEDIA_VERSION"`
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
