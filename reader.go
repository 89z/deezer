package main

import (
   "crypto/cipher"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "io"
   "net/http"
)

func newReader(sngId, source string) (io.Reader, error) {
   var (
      bfKey []byte
      trackHash = md5Hash(sngId)
   )
   for n := 0; n < 16; n++ {
      bfKey = append(bfKey, trackHash[n] ^ trackHash[n + 16] ^ deezerCBC[n])
   }
   block, err := blowfish.NewCipher(bfKey)
   if err != nil {
      return nil, err
   }
   fmt.Println("Get", source)
   get, err := http.Get(source)
   if err != nil {
      return nil, err
   }
   return &reader{Cipher: block, Reader: get.Body, size: 2048}, nil
}

type reader struct {
   *blowfish.Cipher
   io.Reader
   loop int
   size int
}

func (r *reader) Read(data []byte) (int, error) {
   d, err := r.Reader.Read(data)
   for e := 0; d - e >= r.size; e += r.size {
      if r.loop % 3 == 0 {
         text := data[e : e + r.size]
         cipher.NewCBCDecrypter(r.Cipher, deezerIv).CryptBlocks(text, text)
      }
      r.loop++
   }
   return d, err
}
