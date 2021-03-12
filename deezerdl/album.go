package deezer

import (
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strings"
   "time"
)

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

type Track struct {
	ID           int     `json:"SNG_ID,string"`
	Title        string  `json:"SNG_TITLE"`
	TrackNumber  int     `json:"TRACK_NUMBER,string"`
	Gain         float32 `json:"GAIN,string"`
	MD5          string  `json:"MD5_ORIGIN"`
	MediaVersion int     `json:"MEDIA_VERSION,string"`
	api          *API
}

var NoMD5Error = errors.New("no MD5 hash -- try authenticating")

// GetDownloadURL gets the download url (as a *url.URL) for a given
// format
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

// GetMD5 uses an alternative API to get the MD5 of the track
func (track *Track) GetMD5() error {
	resp, err := track.api.MobileApiRequest(getSongMobileMethod,
		strings.NewReader(fmt.Sprintf(`{"SNG_ID":%d}`, track.ID)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // uses the body directly
		return err
	}
	// and then get the MD5
	var results struct {
		MD5 string `json:"MD5_ORIGIN"`
	}
	if err := json.Unmarshal(data.Results, &results); err != nil {
		return err
	}

	if results.MD5 == "" {
		return NoMD5Error
	}
	track.MD5 = results.MD5

	return nil

}

// GetSongData gets a track
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
