package main

import (
   "bytes"
   "fmt"
   "io"
   "os"
   "strings"
   "sync"
)

func decryptMedia(stream io.Reader, id, FName string, streamLen int64) error {
   var (
      bfKey = getBlowFishKey(id)
      chunkSize = 2048
      destBuffer bytes.Buffer
      errc = make(chan error)
   )
   for position, i := 0, 0; position < int(streamLen); position, i = position+chunkSize, i+1 {
      func(i, position int, streamLen int64, stream io.Reader) {
         if (int(streamLen) - position) >= 2048 {
            chunkSize = 2048
         } else {
            chunkSize = int(streamLen) - position
         }
         buf := make([]byte, chunkSize) // The "chunk" of data
         if _, err := io.ReadFull(stream, buf); err != nil {
            errc <- fmt.Errorf("loop %v %v", i, err)
         }
         var (
            chunkString []byte
            err error
         )
         if i % 3 > 0 || chunkSize < 2048 {
            chunkString = buf
         } else { //Decrypt and then write to destBuffer
            chunkString, err = blowfishDecrypt(buf, bfKey)
            if err != nil {
               errc <- fmt.Errorf("loop %v %v", i, err)
            }
         }
         if _, err := destBuffer.Write(chunkString); err != nil {
            errc <- fmt.Errorf("loop %v %v", i, err)
         }
      }(i, position, streamLen, stream)
   }
   var wg sync.WaitGroup
   for {
      select {
      case err := <-errc:
         return err
      default:
         wg.Wait()
         NameWithoutSlash := strings.ReplaceAll(FName, "/", " ")
         out, err := os.Create(NameWithoutSlash)
         if err != nil {
            return err
         }
         _, err = destBuffer.WriteTo(out)
         if err != nil {
            return err
         }
         return nil
      }
   }
}
