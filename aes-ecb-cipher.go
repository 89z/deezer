package main

import (
   "crypto/cipher"
   "errors"
   "golang.org/x/crypto/blowfish"
)

func blowfishDecrypt(buf []byte, bfKey string) ([]byte, error) {
   decrypter, err := blowfish.NewCipher([]byte(bfKey)) // 8bytes
   if err != nil {
      return nil, err
   }
   IV := []byte{0, 1, 2, 3, 4, 5, 6, 7} //8 bytes
   if len(buf) % blowfish.BlockSize != 0 {
      return nil, errors.New("The Buf is not a multiple of 8")
   }
   cipher.NewCBCDecrypter(decrypter, IV).CryptBlocks(buf, buf)
   return buf, nil
}

type ecb struct {
   cipher.Block
}

func newECB(b cipher.Block) ecb {
   return ecb{b}
}

func (x ecb) cryptBlocks(dst, src []byte) error {
   size := x.BlockSize()
   if len(src) % size != 0 {
      return errors.New("crypto/cipher: input not full blocks")
   }
   if len(dst) < len(src) {
      return errors.New("crypto/cipher: output smaller than input")
   }
   for len(src) > 0 {
      x.Encrypt(dst, src[:size])
      src = src[size:]
      dst = dst[size:]
   }
   return nil
}
