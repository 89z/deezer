package deezer

import (
   "crypto/cipher"
   "crypto/md5"
   "fmt"
   "golang.org/x/crypto/blowfish"
)

var (
   iv = []byte{0, 1, 2, 3, 4, 5, 6, 7}
   keyBlowfish = []byte("g4el58wc0zvf9na1")
)

func decrypt(sngId string, data []byte) error {
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

func md5Hash(s string) string {
   b := []byte(s)
   return fmt.Sprintf(
      "%x", md5.Sum(b),
   )
}
