package yamusic

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

const (
	apiURL = "https://api.music.yandex.net"
)

type (
	// Doer is an interface that can do http request
	Doer interface {
		Do(req *http.Request) (*http.Response, error)
	}
	// A Client manages communication with the Yandex.Music API.
	Client struct {
		// HTTP client used to communicate with the API.
		client Doer
		// Base URL for API requests.
		baseURL *url.URL
		// Access token to Yandex.Music API
		accessToken string
		userID      int

		config struct {
			Token  string `yaml:"token"`
			Output string `yaml:"output"`
			Log    string `yaml:"log"`
			Host   string `yaml:"host"`
			Port   string `yaml:"port"`
		}

		// Debug sets should library print debug messages or not
		Debug bool
		// Services
		genres    *GenresService
		search    *SearchService
		account   *AccountService
		feed      *FeedService
		playlists *PlaylistsService
		tracks    *TracksService
	}
)

var logDebug = log.New(os.Stdout, "[DEBUG]\t", log.Ldate|log.Ltime|log.Lshortfile)
var logInfo = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)

// NewClient returns a new API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(options ...func(*Client)) *Client {
	baseURL, _ := url.Parse(apiURL)

	c := &Client{
		client:  http.DefaultClient,
		baseURL: baseURL,
	}

	for _, option := range options {
		option(c)
	}

	c.genres = &GenresService{client: c}
	c.search = &SearchService{client: c}
	c.account = &AccountService{client: c}
	c.feed = &FeedService{client: c}
	c.playlists = &PlaylistsService{client: c}
	c.tracks = &TracksService{client: c}

	return c
}

// HTTPClient sets http client for Yandex.Music client
func HTTPClient(httpClient Doer) func(*Client) {
	return func(c *Client) {
		if httpClient != nil {
			c.client = httpClient
		}
	}
}

// BaseURL sets base API URL for Yandex.Music client
func BaseURL(baseURL *url.URL) func(*Client) {
	return func(c *Client) {
		if baseURL != nil {
			c.baseURL = baseURL
		}
	}
}

// NewConfig reads config from provided path for Yandex.Music client
func NewConfig(configPath string) func(*Client) {
	return func(c *Client) {
		// Create config structure
		// config := &Config{}

		// Open config file
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		// Init new YAML decode
		d := yaml.NewDecoder(file)

		// Start YAML decoding from file
		if err := d.Decode(&c.config); err != nil {
			fmt.Println(err)
		}

		// Create dirs if they are not exist
		if _, err := os.Stat(c.config.Output); os.IsNotExist(err) {
			err := os.MkdirAll(c.config.Output, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}
		if _, err := os.Stat(c.config.Log); os.IsNotExist(err) {
			err := os.MkdirAll(c.config.Log, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// AccessToken sets user_id and access token for Yandex.Music client
func AccessToken(userID int) func(*Client) {
	return func(c *Client) {
		if userID != 0 {
			c.userID = userID
		}

		c.accessToken = c.config.Token
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body, except when body is url.Values. If it is url.Values, it is
// encoded as application/x-www-form-urlencoded and included in request
// headers.
func (c *Client) NewRequest(
	method,
	urlStr string,
	body interface{},
) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	var reader io.Reader
	var isForm bool
	if body != nil {
		switch v := body.(type) {
		case url.Values:
			reader = strings.NewReader(v.Encode())
			isForm = true
		default:
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}

			reader = buf
		}
	}

	req, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "OAuth "+c.accessToken)
	if isForm && method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(
	ctx context.Context,
	req *http.Request,
	v interface{},
) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			dat, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewReader(dat))
			err = json.Unmarshal(dat, v)
			if err == io.EOF {
				if c.Debug {
					logDebug.Println("Got empty")
				}
				// Ignore EOF errors caused by empty response body.
				err = nil //nolint:ineffassign
			} else if err != nil {
				// Try parse XML if it's not JSON.
				err = xml.Unmarshal(dat, v) //nolint:ineffassign,staticcheck
			}
		}
	}

	return resp, err
}

// SetUserID sets user's id in client
func (c *Client) SetUserID(nID int) {
	c.userID = nID
}

// UserID returns id of authorized user. If wasn't authorized returns 0.
func (c *Client) UserID() int {
	return c.userID
}

// Genres returns genres service
func (c *Client) Genres() *GenresService {
	return c.genres
}

// Search returns genres service
func (c *Client) Search() *SearchService {
	return c.search
}

// Account returns account service
func (c *Client) Account() *AccountService {
	return c.account
}

// Feed returns feed service
func (c *Client) Feed() *FeedService {
	return c.feed
}

// Playlists returns playlists service
func (c *Client) Playlists() *PlaylistsService {
	return c.playlists
}

// Tracks returns feed service
func (c *Client) Tracks() *TracksService {
	return c.tracks
}

// General types
type (
	// InvocationInfo is base info in all requests
	InvocationInfo struct {
		Hostname string `json:"hostname"`
		ReqID    string `json:"req-id"`
		// ExecDurationMillis sometimes int, sometimes string so saving interface{}
		ExecDurationMillis interface{} `json:"exec-duration-millis"`
	}
	// Error is struct with error type and message.
	Error struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
)

// print playlist structure to console or file
func (c *Client) PrintPlaylists() {
	result, _, _ := c.Playlists().List(context.Background(), 0)
	playlists := result.Result
	fmt.Println("User ", c.UserID(), " playlists:")
	for _, playlist := range playlists {
		fmt.Println("\t", playlist.Kind, ": ", playlist.Title)
	}
}

// print playlist structure to console or file
func (c *Client) GetTracksWithoutPlaylist() []Track {
	// Get like tracks
	var tracks_without_playlist []Track
	var track_ids []string
	tracks_in_playlist := map[string]bool{}
	res1, _, _ := c.Tracks().GetLike(context.Background())
	for _, track := range res1.Result.Library.Tracks {
		track_ids = append(track_ids, track.ID)
	}
	res2, _, _ := c.Tracks().GetAll(context.Background(), track_ids)
	like_tracks := res2.Result // []Track

	// Search for playlists
	result, _, _ := c.Playlists().List(context.Background(), 0)
	playlists := result.Result
	bar := progressbar.Default(int64(len(playlists)))
	for _, playlist := range playlists {
		bar.Add(1)
		result, _, _ := c.Playlists().Get(context.Background(), c.UserID(), playlist.Kind)
		tracks := result.Result.Tracks
		for _, track := range tracks {
			tracks_in_playlist[track.Track.ID] = true
		}
	}
	fmt.Println(len(tracks_in_playlist))
	fmt.Println(len(like_tracks))

	for _, track := range like_tracks {
		if !tracks_in_playlist[track.ID] {
			tracks_without_playlist = append(tracks_without_playlist, track)
		}
	}

	fmt.Println(len(tracks_without_playlist))

	return tracks_without_playlist
}

// Test
func (c *Client) TestIneraction() {
	client := NewClient(
		// read app config
		NewConfig("yamusic_config.yaml"),
		// create default user with your auth token
		// provide user_id and access_token (needed by some methods)
		AccessToken(0),
	)
	// get user by auth token
	accountStatus, _, _ := client.Account().GetUser(context.Background())
	// add userID to client
	client.SetUserID(accountStatus.Result.UID)
	option := 0
	switch option {
	case 1: // Download Like playlist (kind - 3)
		client.Playlists().DownloadOne(context.Background(), 3)
	case 2: // Create, Rename, Delete playlist
		res1, _, _ := client.Playlists().Create(context.Background(), "Test1", true)
		log.Println(res1.Result.Title)
		kind := res1.Result.Kind
		res2, _, _ := client.Playlists().Rename(context.Background(), kind, "Test2")
		log.Println(res2.Result.Title)
		client.Playlists().Delete(context.Background(), kind)
	case 3:
		client.Playlists().DistributeTracksByPlaylists()

	case 4:
		client.Playlists().DeleteTracksFromPlaylists()

	default:
		log.Printf("Don`t use yamusic\n")
	}
}
