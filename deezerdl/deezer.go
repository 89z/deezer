package main

import (
   "crypto/aes"
   "crypto/cipher"
   "crypto/md5"
   "encoding/json"
   "errors"
   "fmt"
   "golang.org/x/crypto/blowfish"
   "io"
   "io/ioutil"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
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
}

func NewAPI() (*API, error) {
   cookieJar, err := cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   return &API{
      client: &http.Client{Jar: cookieJar},
   }, nil
}

func (api *API) MobileApiRequest(method string, body io.Reader) (*http.Response, error) {
   u := mobileApiUrl
   q := url.Values{
      "api_key": {"4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE"},
      "input":   {"3"},
      "method":  {method},
      "output":  {"3"},
   }
   var sid string
   for _, cookie := range api.client.Jar.Cookies(&mobileApiUrl) {
      if cookie.Name == "sid" {
         sid = cookie.Value
         break
      }
   }
   q.Set("sid", sid)
   u.RawQuery = q.Encode()
   req, err := http.NewRequest(
      http.MethodPost, u.String(), body,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Add("User-Agent", "PostmanRuntime/7.21.0")
   resp, err := api.client.Do(req)
   if err != nil {
      return nil, err
   }
   return resp, nil
}

func (api *API) ApiRequest(method string, body io.Reader) (*http.Response, error) {
   u := apiUrl
   q := url.Values{
      "api_version": {"1.0"}, "input": {"3"}, "method": {method},
   }
   if method == getTokenMethod {
      q.Set("api_token", "null")
   } else {
      q.Set("api_token", api.APIToken)
   }
   u.RawQuery = q.Encode()
   req, err := http.NewRequest(
      http.MethodPost, u.String(), body,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Add("User-Agent", "PostmanRuntime/7.21.0")
   return api.client.Do(req)
}

func (api *API) CookieLogin(arl string) error {
   api.client.Jar.SetCookies(&deezerUrl, []*http.Cookie{
      {Domain: ".deezer.com", Name: "arl", Path: "/", Value: arl},
   })
   err := api.getSession()
   if err != nil {
      return err
   }
   api.APIToken, err = api.getToken()
   return err
}

func (api *API) getToken() (string, error) {
   resp, err := api.ApiRequest(getTokenMethod, nil)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var data struct{Results json.RawMessage `json:"results"`}
   if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // uses the body directly
   return "", err
   }
   var results struct{Token string `json:"checkForm"`}
   if err := json.Unmarshal(data.Results, &results); err != nil {
      return "", err
   }
   return results.Token, nil
}

func (api *API) getSession() error {
   req, err := http.NewRequest(http.MethodPost, "https://www.deezer.com", nil)
   if err != nil {
      return err
   }
   _, err = api.client.Do(req)
   return err
}

const (
   blowfishKey = "g4el58wc0zvf9na1"
   blowfishIV = "\x00\x01\x02\x03\x04\x05\x06\x07"
   fileChunkSize = 2048
)

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

func (track *Track) GetBlowfishKey() []byte {
   hash := MD5Hash([]byte(fmt.Sprintf("%d", track.ID)))
   key := []byte(blowfishKey)
   output := make([]byte, 16)
   for i := 0; i < 16; i++ {
      output[i] = hash[i] ^ hash[i+16] ^ key[i]
   }
   return output
}

func DecryptSongFile(key []byte, inputPath, outputPath string) error {
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
      encrypted := (chunk%3 == 0)
      if encrypted && n == fileChunkSize {
         buf, err = decryptBlowfish(key, buf)
         if err != nil {
            return err
         }
      }
      n, err = outFile.Write(buf)
      if err != nil {
         return err
      }
      n, err = inFile.Read(buf)
      if err != nil && err != io.EOF {
         return err
      }
   }
   return nil
}

var (
   combineChar = []byte("\xa4")
   deezerKey = []byte("jo6aey6haid2Teih")
)

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



const AlbumAPIFormat = "https://api.deezer.com/album/%d"

type AlbumTrack struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// AlbumResponse is an intermediate format for getting album data that
// stores the data before putting it in an Album struct
type AlbumResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	CoverURL    string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXL     string `json:"cover_xl"`
	Date        string `json:"release_date"`
	Tracks      struct {
		Data []AlbumTrack `json:"data"`
	} `json:tracks"`
}

// Album stores the data for the album of interest
type Album struct {
	ID        int
	Title     string
	Link      string
	CoverURL  string
	Covers    Covers
	Date      time.Time
	Tracklist []AlbumTrack
	Tracks    []*Track
	api       *API
}

// Covers stores the different cover sizes available
type Covers struct {
	Small  string
	Medium string
	Big    string
	XL     string
}

// NewAlbum create an Album from an AlbumResponse
func NewAlbum(response *AlbumResponse, api *API) (*Album, error) {
	date, err := time.Parse("2006-01-02", response.Date)
	if err != nil {
		return nil, err
	}
	album := Album{
		ID:       response.ID,
		Title:    response.Title,
		Link:     response.Link,
		CoverURL: response.CoverURL,
		Covers: Covers{
			Small:  response.CoverSmall,
			Medium: response.CoverMedium,
			Big:    response.CoverBig,
			XL:     response.CoverXL,
		},
		Date:      date,
		Tracklist: response.Tracks.Data,
		api:       api,
	}

	return &album, nil
}

// albumRequest performs the API request for the deezer album
// remember to close the body
func (api *API) albumRequest(ID int) (*http.Response, error) {
	// construct the request
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(AlbumAPIFormat, ID),
		nil)

	if err != nil {
		return nil, err
	}

	// send the request
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetAlbum gets the album based on its ID
func (api *API) GetAlbumData(ID int) (*Album, error) {
	// make a request to the public API
	resp, err := api.albumRequest(ID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode the json
	var response AlbumResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// convert to album
	album, err := NewAlbum(&response, api)
	if err != nil {
		return nil, err
	}

	return album, nil
}

// GetTracks gets all tracks in an album and store them in
// album.Tracks. Also return the slice.
func (album *Album) GetTracks() ([]*Track, error) {
	for _, t := range album.Tracklist {
		track, err := album.api.GetSongData(t.ID)
		if err != nil {
			return []*Track{}, err
		}
		album.Tracks = append(album.Tracks, track)
	}

	return album.Tracks, nil
}



type Format int

const (
	FLAC    Format = 9
	MP3_320        = 3
	MP3_256        = 5
)

const (
	downloadHostFormat = "e-cdns-proxy-%c.dzcdn.net"
	downloadPathFormat = "/mobile/1/%s"
)

var NoMD5Error = errors.New("no MD5 hash -- try authenticating")

type Track struct {
   ID           int     `json:"SNG_ID,string"`
   Title        string  `json:"SNG_TITLE"`
   TrackNumber  int     `json:"TRACK_NUMBER,string"`
   Gain         float32 `json:"GAIN,string"`
   MD5          string  `json:"MD5_ORIGIN"`
   MediaVersion int     `json:"MEDIA_VERSION,string"`
   api          *API
}

func (track *Track) GetDownloadURL(format Format) (*url.URL, error) {
	if len(track.MD5) == 0 {
		if err := track.GetMD5(); err != nil {
			return nil, NoMD5Error
		}
	}
	path, err := MakeURLPath(track, format)
	if err != nil {
		return nil, err
	}
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf(downloadHostFormat, track.MD5[0]),
		Path:   fmt.Sprintf(downloadPathFormat, path),
	}
	return &u, nil
}

func (track *Track) GetMD5() error {
   resp, err := track.api.MobileApiRequest(
      getSongMobileMethod,
      strings.NewReader(fmt.Sprintf(`{"SNG_ID":%d}`, track.ID)),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   var data struct{Results json.RawMessage `json:"results"`}
   json.Unmarshal(body, &data)
   var results struct {MD5 string `json:"MD5_ORIGIN"`}
   if err := json.Unmarshal(data.Results, &results); err != nil {
      return err
   }
   if results.MD5 == "" {
      return NoMD5Error
   }
   track.MD5 = results.MD5
   return nil
}

func (api *API) GetSongData(ID int) (*Track, error) {
	// make the request
	body := strings.NewReader(fmt.Sprintf(`{"SNG_ID":%d}`, ID))
	resp, err := api.ApiRequest(getSongMethod, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// decode track from results
	var track Track
	if err := json.Unmarshal(data.Results, &track); err != nil {
		return nil, err
	}
	track.api = api

	return &track, nil
}
