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
   return &reader{
      BlockMode: cipher.NewCBCDecrypter(block, deezerIv),
      Reader: get.Body,
      text: make([]byte, 2048),
   }, nil
}

type reader struct {
   cipher.BlockMode
   io.Reader
   loop int
   text []byte
}

func (r *reader) Read(b []byte) (int, error) {
   n, err := r.Reader.Read(r.text)
   if err != nil {
      return 0, err
   }
   if n < len(r.text) {
      r.text = r.text[:n]
   } else if r.loop % 3 == 0 {
      r.CryptBlocks(r.text, r.text)
   }
   r.loop++
   return copy(b, r.text), nil
}

/*
16384
16384

32768
32768

last
19434
*/
