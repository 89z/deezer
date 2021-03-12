package deezer

import (
   "crypto/aes"
   "crypto/cipher"
   "crypto/md5"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "io"
   "os"
   "strconv"
   "strings"
)

const blowfishKey = "g4el58wc0zvf9na1"
const blowfishIV = "\x00\x01\x02\x03\x04\x05\x06\x07"

const fileChunkSize = 2048

// decryptBlowfish decrypts blowfish data
func decryptBlowfish(key, data []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(blowfishIV))

	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	return decrypted, nil
}

// getBlowfishKey calculates the key required to decrypt the
// blowfish-encrypted file
func (track *Track) GetBlowfishKey() []byte {
	hash := MD5Hash([]byte(fmt.Sprintf("%d", track.ID)))
	key := []byte(blowfishKey)

	output := make([]byte, 16)
	for i := 0; i < 16; i++ {
		output[i] = hash[i] ^ hash[i+16] ^ key[i]
	}
	return output
}

// DecryptSongFile decrypts the encrypted chunks of a song downloaded
// from deezer
func DecryptSongFile(key []byte, inputPath, outputPath string) error {
	// open files
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	buf := make([]byte, fileChunkSize)
	n, err := inFile.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	for chunk := 0; n > 0; chunk++ {
		// only decrypt every third chunk (including first
		// chunk)
		encrypted := (chunk%3 == 0)

		// only decrypt if encrypted and whole chunk
		if encrypted && n == fileChunkSize {
			buf, err = decryptBlowfish(key, buf)
			if err != nil {
				return err
			}
		}

		// write the chunk back
		n, err = outFile.Write(buf)
		if err != nil {
			return err
		}

		// read next chunk
		n, err = inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}



var combineChar = []byte("\xa4")
var deezerKey = []byte("jo6aey6haid2Teih")

func MD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

// MakeURLPath generates the path of the download URL
func MakeURLPath(track *Track, format Format) (string, error) {
	// generate MD5 data
	chars := []string{
		track.MD5,
		strconv.Itoa(int(format)),
		strconv.Itoa(track.ID),
		strconv.Itoa(track.MediaVersion),
	}
	md5Data := strings.Join(chars, string(combineChar))
	hash := []byte(MD5Hash([]byte(md5Data)))

	// generate and return hex of encrypted data
	encData := append(hash, combineChar...)
	encData = append(encData, md5Data...)
	encData = append(encData, combineChar...)
	ecb, err := ECB(deezerKey, encData)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", ecb), nil
}

func ECB(key, data []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	size := cipher.BlockSize()
	for len(data)%size != 0 {
		data = append(data, '\x00')
	}

	encrypted := make([]byte, len(data))
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Encrypt(encrypted[bs:be], data[bs:be])
	}

	return encrypted, nil
}
