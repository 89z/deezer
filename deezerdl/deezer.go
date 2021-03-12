package deezer

import (
   "crypto/aes"
   "crypto/cipher"
   "crypto/md5"
   "encoding/json"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "io"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "strconv"
   "strings"
)

const (
	getTokenMethod      = "deezer.getUserData"
	getSongMethod       = "song.getData"
	getSongMobileMethod = "song_getData"
)

var apiUrl = url.URL{
	Scheme: "https",
	Host:   "www.deezer.com",
	Path:   "/ajax/gw-light.php",
}

var mobileApiUrl = url.URL{
	Scheme: "https",
	Host:   "api.deezer.com",
	Path:   "/1.0/gateway.php",
}

var deezerUrl = url.URL{
	Scheme: "http",
	Host:   ".deezer.com",
	Path:   "/",
}

type API struct {
	APIToken  string
	client    *http.Client
	DebugMode bool
}

// NewAPI creates a new API with a http Client with cookie jar
func NewAPI(debugMode bool) (*API, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Jar: cookieJar,
	}
	api := API{
		client:    &client,
		DebugMode: debugMode,
	}
	return &api, nil
}

// ApiRequest performs an API request
func (api *API) ApiRequest(method string, body io.Reader) (*http.Response, error) {
	// add the required parameters to the URL
	u := apiUrl
	q := url.Values{
		"api_version": {"1.0"},
		"input":       {"3"},
		"method":      {method},
	}

	// add the token only if we know the token
	if method == getTokenMethod {
		q.Set("api_token", "null")
	} else {
		q.Set("api_token", api.APIToken)
	}
	u.RawQuery = q.Encode()

	// construct the request
	req, err := http.NewRequest(http.MethodPost,
		u.String(),
		body)
	if err != nil {
		return nil, err
	}

	// change user agent
	req.Header.Add("User-Agent", "PostmanRuntime/7.21.0")

	// send
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// MobileApiRequest performs a mobile API request
func (api *API) MobileApiRequest(method string, body io.Reader) (*http.Response, error) {
	// add the required parameters to the URL
	u := mobileApiUrl
	q := url.Values{
		"api_key": {"4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE"},
		"input":   {"3"},
		"output":  {"3"},
		"method":  {method},
	}

	// get the current sid from the cookie jar
	var sid string
	for _, cookie := range api.client.Jar.Cookies(&deezerUrl) {
		if cookie.Name == "sid" {
			sid = cookie.Value
			break
		}
	}
	q.Set("sid", sid)

	u.RawQuery = q.Encode()

	// construct the request
	req, err := http.NewRequest(http.MethodPost,
		u.String(),
		body)
	if err != nil {
		return nil, err
	}

	// change user agent
	req.Header.Add("User-Agent", "PostmanRuntime/7.21.0")

	// send
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CookieLogin allows the user to log in using their arl cookie taken
// from a browser
func (api *API) CookieLogin(arl string) error {
	// add the cookie to the jar
	cookie := http.Cookie{
		Name:   "arl",
		Value:  arl,
		Domain: ".deezer.com",
		Path:   "/",
	}
	api.client.Jar.SetCookies(&deezerUrl, []*http.Cookie{&cookie})

	// get a session
	err := api.getSession()
	if err != nil {
		return err
	}

	// try to get the token
	api.APIToken, err = api.getToken()
	if err != nil {
		return err
	}

	return nil
}

// GetToken gets the user's API token
func (api *API) getToken() (string, error) {
	// make the request
	resp, err := api.ApiRequest(getTokenMethod, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if api.DebugMode {
	}

	// decode result key into a struct from the body
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // uses the body directly
		return "", err
	}
	// and then the checkForm key (the token)
	var results struct {
		Token string `json:"checkForm"`
	}
	if err := json.Unmarshal(data.Results, &results); err != nil {
		return "", err
	}

	if api.DebugMode {
	}
	return results.Token, nil
}

func (api *API) getSession() error {
	req, err := http.NewRequest(http.MethodPost,
		"https://www.deezer.com",
		nil)
	if err != nil {
		return err
	}
	_, err = api.client.Do(req)
      return err
}



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
