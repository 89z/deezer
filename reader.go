package main

import (
   "crypto/cipher"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "io"
   "net/http"
)

type reader struct {
   cipher.BlockMode
   io.Reader
   loop int
   size int
}

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
   return &reader{
      BlockMode: cipher.NewCBCDecrypter(block, deezerIv),
      Reader: get.Body,
      size: 2048,
   }, nil
}

func (r *reader) Read(text []byte) (int, error) {
   n, err := r.Reader.Read(text)
   if err != nil {
      return 0, err
   }
   if n < r.size {
      text = text[:n]
   }
   if r.loop % 3 == 0 && n == r.size {
      r.CryptBlocks(text, text)
   }
   r.loop++
   return n, nil
}

/*
16384
16384

32768
32768

last
19434
*/
