package yamusic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type (
	// TracksService is a service to deal with tracks
	TracksService struct {
		client *Client
	}
	// TracksResp describes get user's tracks/like tracks/ response
	TrackResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         []Track        `json:"result"`
	}
	// TracksResp describes get user's tracks/like tracks/ response
	LikeTracksResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         struct {
			Library struct {
				UID          int         `json:"uid"`
				Revision     int         `json:"revision"`
				PlaylistUuid string      `json:"playlistUuid"`
				Tracks       []TrackLike `json:"tracks"`
			} `json:"library"`
		} `json:"result"`
	}
	// TracksResp describes get user's tracks/like tracks/ response
	Supplement struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         struct {
			ID     string `json:"id"`
			Lyrics struct {
				ID              int    `json:"id"`
				Lyrics          int    `json:"lyrics"`
				FullLyrics      string `json:"fullLyrics"`
				HasRights       bool   `json:"hasRights"`
				ShowTranslation bool   `json:"showTranslation"`
			} `json:"lyrics"`
		} `json:"result"`
	}
	// Response of track/%d/download_info
	DownloadInfoResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         []struct {
			Codec           string `json:"codec"`
			Gain            bool   `json:"gain"`
			Preview         bool   `json:"preview"`
			DownloadInfoURL string `json:"downloadInfoUrl"`
			Direct          bool   `json:"direct"`
			BitrateInKbps   int    `json:"bitrateInKbps"`
		} `json:"result"`
	}
	// DownloadInfo is a response of URL from DownloadInfoResp's `DownloadInfoURL` field
	DownloadInfo struct {
		XMLName xml.Name `xml:"download-info"`
		Text    string   `xml:",chardata"`
		Host    string   `xml:"host"`
		Path    string   `xml:"path"`
		TS      string   `xml:"ts"`
		Region  string   `xml:"region"`
		S       string   `xml:"s"`
	}
)

type TrackError string

func (te TrackError) Error() string { return string(te) }

var (
	ErrNilDownloadInfoResp = TrackError("got nil download info resp pointer")
	ErrNilDownloadInfo     = TrackError("got nil download info pointer")
	ErrNilPath             = TrackError("got nil path")
	ErrEmptyPath           = TrackError("got empty path")
	ErrZeroResultLen       = TrackError("len of download inf response's result field is zero")
)

// Get returns track by its ID
func (t *TracksService) GetOne(ctx context.Context, id int) (*TrackResp, *http.Response, error) {
	uri := fmt.Sprintf("tracks/%v", id)
	req, err := t.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}
	track := new(TrackResp)
	resp, err := t.client.Do(ctx, req, track)
	return track, resp, err
}

// Get returns track by its ID
func (t *TracksService) GetAll(ctx context.Context, track_ids []string) (*TrackResp, *http.Response, error) {
	uri := "tracks"

	form := url.Values{}
	form.Set("track-ids", strings.Join(track_ids[:], ","))
	form.Set("with-positions", "false")

	req, err := t.client.NewRequest(http.MethodPost, uri, form)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tracks := new(TrackResp)
	resp, err := t.client.Do(ctx, req, tracks)

	return tracks, resp, err
}

// List returns playlists of the user
func (t *TracksService) GetLike(ctx context.Context) (*LikeTracksResp, *http.Response, error) {
	uri := fmt.Sprintf("users/%v/likes/tracks", t.client.userID)
	req, err := t.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	like_tracks := new(LikeTracksResp)
	resp, err := t.client.Do(ctx, req, like_tracks)
	return like_tracks, resp, err
}

// List returns playlists of the user
func (t *TracksService) GetSupplement(ctx context.Context, id string) (*Supplement, *http.Response, error) {
	uri := fmt.Sprintf("tracks/%v/supplement", id)
	req, err := t.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	supplement := new(Supplement)
	resp, err := t.client.Do(ctx, req, supplement)

	return supplement, resp, err
}

// GetDownloadInfoResp returns DownloadInfoResp byt track's ID
// Be careful: you can get DownloadInfo by DownloadInfoURL only
// for one minute since you called GetDownloadInfoResp
func (t *TracksService) GetDownloadInfoResp(ctx context.Context, id int) (*DownloadInfoResp, *http.Response, error) {
	uri := fmt.Sprintf("tracks/%v/download-info", id)
	req, err := t.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	dlInfoResp := new(DownloadInfoResp)
	resp, err := t.client.Do(ctx, req, dlInfoResp)
	return dlInfoResp, resp, err
}

// GetDownloadInfo returns DownloadInfo by id of track.
// Be careful: it uses the same context for GetDownloadInfoResp
// and its request
func (t *TracksService) GetDownloadInfo(ctx context.Context, id int) (*DownloadInfo, *http.Response, error) {
	dlInfoResp, dirResp, err := t.GetDownloadInfoResp(ctx, id)
	if err != nil {
		return nil, dirResp, err
	}

	if dlInfoResp == nil {
		return nil, nil, ErrNilDownloadInfoResp
	}

	if len(dlInfoResp.Result) == 0 {
		return nil, nil, ErrZeroResultLen
	}

	req, err := t.client.NewRequest(http.MethodGet, dlInfoResp.Result[0].DownloadInfoURL, nil)
	if err != nil {
		return nil, nil, err
	}

	dlInfo := new(DownloadInfo)
	resp, err := t.client.Do(ctx, req, dlInfo)
	return dlInfo, resp, err
}

// GetDownloadURL computes path to track by ID
func (t *TracksService) GetDownloadURL(ctx context.Context, id int) (string, error) {
	dlInfo, _, err := t.GetDownloadInfo(ctx, id)
	if err != nil {
		return "", err
	}
	if dlInfo == nil {
		return "", ErrNilDownloadInfo
	} else if len(dlInfo.Path) == 0 {
		return "", ErrEmptyPath
	}
	// a bit of magic
	const signPrefix = "XGRlBW9FXlekgbPrRHuSiA"
	sign := md5.Sum([]byte(signPrefix + dlInfo.Path[1:] + dlInfo.S))
	uri := fmt.Sprintf(
		"https://%s/get-mp3/%s/%s%s",
		dlInfo.Host,
		hex.EncodeToString(sign[:]),
		dlInfo.TS, dlInfo.Path,
	)
	return uri, nil
}

// Download track by DownloadURL by path on fs
func (t *TracksService) DownloadAll(ctx context.Context, tracks []Track, path string) {
	/// Get list directory
	entries, err := os.ReadDir(path + "/tracks")
	if err != nil {
		log.Fatal(err)
	}

	var entry_name string
	tracks_on_fs := map[string]bool{}

	for _, entry := range entries {
		entry_name = entry.Name()
		entry_name = strings.ReplaceAll(entry_name, ".mp3", "")
		logInfo.Println(entry_name)
		tracks_on_fs[entry_name] = true
	}

	logInfo.Printf("Already loaded tracks by path %s: %d", path, len(entries))
	for _, track := range tracks {
		file_name := t.GetFileName(ctx, track)
		if !tracks_on_fs[file_name] {
			logInfo.Printf("%s", file_name)
			t.Download(ctx, track, path)
		}
	}
}

// Download track by DownloadURL by path on fs
func (t *TracksService) Download(ctx context.Context, track Track, path string) {

	// load track mp3
	file_name := path + "/tracks/" + t.GetFileName(ctx, track) + ".mp3"

	track_id, _ := strconv.Atoi(track.ID)
	uri, _ := t.GetDownloadURL(context.Background(), track_id)
	logInfo.Println(uri)

	output_file, _ := os.Create(file_name)
	defer output_file.Close()

	resp, err := http.Get(uri)
	if err == nil {
		defer resp.Body.Close()
		io.Copy(output_file, resp.Body)
	} else {
		logInfo.Println(err)
		logInfo.Println("Cannot load " + file_name)
	}

	// load track lyrics txt
	if track.LyricsAvailable {
		file_name = path + "/lyrics/" + t.GetFileName(ctx, track) + ".txt"
		supplement, _, _ := t.GetSupplement(ctx, track.ID)
		// open input file
		fi, err := os.Create(file_name)
		if err != nil {
			panic(err)
		}
		// close fi on exit and check for its returned error
		defer func() {
			if err := fi.Close(); err != nil {
				panic(err)
			}
		}()
		fmt.Fprintln(fi, supplement.Result.Lyrics.FullLyrics)
	}

}

func (t *TracksService) GetFileName(ctx context.Context, track Track) string {
	var file_name string
	if len(track.Artists) > 0 {
		file_name += track.Artists[0].Name + " - " + track.Title
	} else {
		file_name += track.Title
	}
	file_name = strings.ReplaceAll(file_name, "/", "|")
	if len(file_name) > 30 {
		file_name = file_name[0:30]
	}
	return file_name
}
