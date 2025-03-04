package yamusic

import (
	"awesome/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type (
	// PlaylistsService is a service to deal with accounts.
	PlaylistsService struct {
		client *Client
	}

	// PlaylistsListResp describes get user's playlists response
	PlaylistsListResp struct {
		InvocationInfo InvocationInfo    `json:"invocationInfo"`
		Error          Error             `json:"error"`
		Result         []PlaylistsResult `json:"result"`
	}

	Artist struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Various  bool   `json:"various"`
		Composer bool   `json:"composer"`
		Cover    struct {
			Type   string `json:"type"`
			Prefix string `json:"prefix"`
			URI    string `json:"uri"`
		} `json:"cover"`
		Genres []interface{} `json:"genres"`
	}

	Artists []Artist

	Label struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}

	Labels []Label

	Album struct {
		ID                       int           `json:"id"`
		Title                    string        `json:"title"`
		Type                     string        `json:"type,omitempty"`
		MetaType                 string        `json:"metaType"`
		ContentWarning           string        `json:"contentWarning,omitempty"`
		Year                     int           `json:"year"`
		ReleaseDate              time.Time     `json:"releaseDate"`
		CoverURI                 string        `json:"coverUri"`
		OgImage                  string        `json:"ogImage"`
		Genre                    string        `json:"genre"`
		Buy                      []interface{} `json:"buy"`
		TrackCount               int           `json:"trackCount"`
		LikesCount               int           `json:"likesCount"`
		Recent                   bool          `json:"recent"`
		VeryImportant            bool          `json:"veryImportant"`
		Available                bool          `json:"available"`
		AvailableForPremiumUsers bool          `json:"availableForPremiumUsers"`
		AvailableForOptions      []string      `json:"availableForOptions"`
		AvailableForMobile       bool          `json:"availableForMobile"`
		AvailablePartially       bool          `json:"availablePartially"`
		Bests                    []int         `json:"bests"`
		Artists                  Artists       `json:"artists"`
		Labels                   Labels        `json:"labels"`
		TrackPosition            struct {
			Volume int `json:"volume"`
			Index  int `json:"index"`
		} `json:"trackPosition"`
	}

	Albums []Album

	Track struct {
		ID             string `json:"id"`
		RealID         string `json:"realId"`
		Title          string `json:"title"`
		ContentWarning string `json:"contentWarning"`
		TrackSource    string `json:"trackSource"`
		Major          struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"major"`
		Available                      bool     `json:"available"`
		AvailableForPremiumUsers       bool     `json:"availableForPremiumUsers"`
		AvailableFullWithoutPermission bool     `json:"availableFullWithoutPermission"`
		AvailableForOptions            []string `json:"availableForOptions"`
		DurationMs                     int      `json:"durationMs"`
		StorageDir                     string   `json:"storageDir"`
		FileSize                       int      `json:"fileSize"`
		R128                           struct {
			I  float64 `json:"i"`
			Tp float64 `json:"tp"`
		} `json:"r128"`
		PreviewDurationMs int     `json:"previewDurationMs"`
		Artists           Artists `json:"artists"`
		Albums            Albums  `json:"albums"`
		CoverURI          string  `json:"coverUri"`
		OgImage           string  `json:"ogImage"`
		LyricsAvailable   bool    `json:"lyricsAvailable"`
		LyricsInfo        struct {
			HasAvailableSyncLyrics bool `json:"hasAvailableSyncLyrics"`
			HasAvailableTextLyrics bool `json:"hasAvailableTextLyrics"`
		} `json:"lyricsInfo"`
		Type             string `json:"type"`
		RememberPosition bool   `json:"rememberPosition"`
		TrackSharingFlag string `json:"trackSharingFlag"`
	}

	TrackFull struct {
		ID        int       `json:"id"`
		Timestamp time.Time `json:"timestamp"`
		Recent    bool      `json:"recent"`
		Track     Track     `json:"track"`
	}

	TrackLike struct {
		ID        string    `json:"id"`
		AlbumId   string    `json:"albumId"`
		Timestamp time.Time `json:"timestamp"`
	}

	Tracks []TrackFull

	PlaylistWithTracks struct {
		PlaylistsResult
		Tracks Tracks `json:"tracks"`
	}

	// PlaylistsGetResp describes get user's playlist by kind response
	PlaylistsGetResp struct {
		InvocationInfo InvocationInfo     `json:"invocationInfo"`
		Error          Error              `json:"error"`
		Result         PlaylistWithTracks `json:"result"`
	}

	// PlaylistsGetByKindsResp describes get user's playlists by kinds response
	PlaylistsGetByKindsResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         []struct {
			PlaylistsResult
			Tracks []struct {
				ID        int       `json:"id"`
				AlbumID   int       `json:"albumId"`
				Timestamp time.Time `json:"timestamp"`
			} `json:"tracks"`
		} `json:"result"`
	}

	// PlaylistsRenameResp describes method rename playlist response
	PlaylistsRenameResp struct {
		InvocationInfo InvocationInfo  `json:"invocationInfo"`
		Error          Error           `json:"error"`
		Result         PlaylistsResult `json:"result"`
	}
	// PlaylistsCreateResp describes method create playlist response
	PlaylistsCreateResp struct {
		InvocationInfo InvocationInfo  `json:"invocationInfo"`
		Error          Error           `json:"error"`
		Result         PlaylistsResult `json:"result"`
	}
	// PlaylistsDeleteResp describes method delete playlist response
	PlaylistsDeleteResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         string         `json:"result"`
	}
	// PlaylistsAddTracksResp describes method add tracks response
	PlaylistsAddTracksResp struct {
		InvocationInfo InvocationInfo  `json:"invocationInfo"`
		Error          Error           `json:"error"`
		Result         PlaylistsResult `json:"result"`
	}
	// PlaylistsRemoveTracksResp describes method add tracks response
	PlaylistsRemoveTracksResp struct {
		InvocationInfo InvocationInfo  `json:"invocationInfo"`
		Error          Error           `json:"error"`
		Result         PlaylistsResult `json:"result"`
	}
	// PlaylistsResult is base result of methods AddTracks and RemoveTracks
	PlaylistsResult struct {
		UID                int            `json:"uid"`
		Kind               int            `json:"kind"`
		Revision           int            `json:"revision"`
		TrackCount         int            `json:"trackCount"`
		DurationMs         int            `json:"durationMs"`
		Collective         bool           `json:"collective"`
		Available          bool           `json:"available"`
		IsBanner           bool           `json:"isBanner"`
		IsPremiere         bool           `json:"isPremiere"`
		Title              string         `json:"title"`
		Description        string         `json:"description"`
		Visibility         string         `json:"visibility"`
		OgImage            string         `json:"ogImage"`
		Created            time.Time      `json:"created"`
		Modified           time.Time      `json:"modified"`
		Cover              PlaylistsCover `json:"cover"`
		Owner              PlaylistsOwner `json:"owner"`
		Tags               []interface{}  `json:"tags"`
		LastOwnerPlaylists []struct {
			UID        int            `json:"uid"`
			Kind       int            `json:"kind"`
			Revision   int            `json:"revision"`
			TrackCount int            `json:"trackCount"`
			DurationMs int            `json:"durationMs"`
			Collective bool           `json:"collective"`
			Available  bool           `json:"available"`
			IsBanner   bool           `json:"isBanner"`
			IsPremiere bool           `json:"isPremiere"`
			Title      string         `json:"title"`
			Visibility string         `json:"visibility"`
			OgImage    string         `json:"ogImage"`
			Created    time.Time      `json:"created"`
			Modified   time.Time      `json:"modified"`
			Tags       []interface{}  `json:"tags"`
			Owner      PlaylistsOwner `json:"owner"`
			Cover      PlaylistsCover `json:"cover"`
		} `json:"lastOwnerPlaylists"`
	}
	// PlaylistsCover is cover of playlist response
	PlaylistsCover struct {
		Error    string   `json:"error"`
		Type     string   `json:"type"`
		ItemsURI []string `json:"itemsUri"`
		Custom   bool     `json:"custom"`
		Dir      string   `json:"dir"`
		Version  string   `json:"version"`
		URI      string   `json:"uri"`
	}
	// PlaylistsOwner is owner of playlist response
	PlaylistsOwner struct {
		UID      int    `json:"uid"`
		Login    string `json:"login"`
		Name     string `json:"name"`
		Sex      string `json:"sex"`
		Verified bool   `json:"verified"`
	}
)

// List returns playlists of the user
func (s *PlaylistsService) List(
	ctx context.Context,
	userID int,
) (*PlaylistsListResp, *http.Response, error) {
	if userID == 0 {
		userID = s.client.userID
	}

	uri := fmt.Sprintf("users/%v/playlists/list", userID)
	req, err := s.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	playlists := new(PlaylistsListResp)
	resp, err := s.client.Do(ctx, req, playlists)
	return playlists, resp, err
}

// Get returns playlist of the user by kind
func (s *PlaylistsService) Get(
	ctx context.Context,
	userID int,
	kind int,
) (*PlaylistsGetResp, *http.Response, error) {
	if userID == 0 {
		userID = s.client.userID
	}

	uri := fmt.Sprintf("users/%v/playlists/%v", userID, kind)
	req, err := s.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	// var playlist interface{}
	playlist := new(PlaylistsGetResp)
	resp, err := s.client.Do(ctx, req, playlist)
	// b, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(b))

	return playlist, resp, err
}

func (s *PlaylistsService) GetByUserIDAndKind(
	ctx context.Context,
	userID string,
	kind int,
) (*PlaylistsGetResp, *http.Response, error) {
	if len(userID) == 0 {
		userID = strconv.Itoa(s.client.userID)
	}

	uri := fmt.Sprintf("users/%v/playlists/%v", userID, kind)
	req, err := s.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	playlist := new(PlaylistsGetResp)
	resp, err := s.client.Do(ctx, req, playlist)
	return playlist, resp, err
}

type (
	// PlaylistsGetByKindOptions options for GetByKinds method
	PlaylistsGetByKindOptions struct {
		Kinds      []int
		Mixed      bool
		RichTracks bool
	}
)

// GetByKinds returns several playlists by kinds with track ids
func (s *PlaylistsService) GetByKinds(
	ctx context.Context,
	userID int,
	opts *PlaylistsGetByKindOptions,
) (*PlaylistsGetByKindsResp, *http.Response, error) {
	if userID == 0 {
		userID = s.client.userID
	}

	if opts == nil {
		opts = &PlaylistsGetByKindOptions{}
	}

	queryParams := url.Values{}
	queryParams.Set("kinds", func() string {
		stringKinds := make([]string, 0, len(opts.Kinds))
		for _, kind := range opts.Kinds {
			stringKinds = append(stringKinds, strconv.Itoa(kind))
		}
		return strings.Join(stringKinds, ",")
	}())
	queryParams.Set("mixed", strconv.FormatBool(opts.Mixed))
	queryParams.Set("rich-tracks", strconv.FormatBool(opts.RichTracks))

	uri := fmt.Sprintf("users/%v/playlists?%v", userID, queryParams.Encode())
	req, err := s.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	playlists := new(PlaylistsGetByKindsResp)
	resp, err := s.client.Do(ctx, req, playlists)
	return playlists, resp, err
}

// Rename renames playlist of current user
func (s *PlaylistsService) Rename(
	ctx context.Context,
	kind int,
	newName string,
) (*PlaylistsRenameResp, *http.Response, error) {
	uri := fmt.Sprintf("users/%v/playlists/%v/name", s.client.userID, kind)

	form := url.Values{}
	form.Set("value", newName)

	req, err := s.client.NewRequest(http.MethodPost, uri, form)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	renamedPlaylist := new(PlaylistsRenameResp)
	resp, err := s.client.Do(ctx, req, renamedPlaylist)
	return renamedPlaylist, resp, err
}

// Create creates playlist for current user
func (s *PlaylistsService) Create(
	ctx context.Context,
	title string,
	isPublic bool,
) (*PlaylistsCreateResp, *http.Response, error) {
	var visibility string
	if isPublic {
		visibility = "public"
	} else {
		visibility = "private"
	}

	form := url.Values{}
	form.Set("title", title)
	form.Set("visibility", visibility)

	uri := fmt.Sprintf("users/%v/playlists/create", s.client.userID)

	req, err := s.client.NewRequest(http.MethodPost, uri, form)
	if err != nil {
		return nil, nil, err
	}

	createdPlaylist := new(PlaylistsCreateResp)
	resp, err := s.client.Do(ctx, req, createdPlaylist)
	return createdPlaylist, resp, err
}

// Delete deletes playlist for current user
func (s *PlaylistsService) Delete(
	ctx context.Context,
	kind int,
) (*PlaylistsDeleteResp, *http.Response, error) {
	uri := fmt.Sprintf("users/%v/playlists/%v/delete", s.client.userID, kind)
	req, err := s.client.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	deletedPlaylist := new(PlaylistsDeleteResp)
	resp, err := s.client.Do(ctx, req, deletedPlaylist)
	return deletedPlaylist, resp, err
}

type (
	// PlaylistsTrack is track object with trackId and albumId
	// that is used in add tracks and remove tracks requests
	PlaylistsTrack struct {
		ID      int `json:"id"`
		AlbumID int `json:"albumId"`
	}
	// PlaylistsAddTracksOptions are options for method AddTracks
	PlaylistsAddTracksOptions struct {
		At int
	}
)

// AddTracks adds tracks to playlist
func (s *PlaylistsService) AddTracks(
	ctx context.Context,
	kind int,
	revision int,
	tracks []PlaylistsTrack,
	opts *PlaylistsAddTracksOptions,
) (*PlaylistsAddTracksResp, *http.Response, error) {
	if opts == nil {
		opts = &PlaylistsAddTracksOptions{
			At: 0,
		}
	}

	diff := []struct {
		Op     string           `json:"op"`
		At     int              `json:"at"`
		Tracks []PlaylistsTrack `json:"tracks"`
	}{
		{
			Op:     "insert",
			At:     opts.At,
			Tracks: tracks,
		},
	}

	b, err := json.Marshal(diff)
	if err != nil {
		return nil, nil, err
	}

	form := url.Values{}
	form.Set("diff", string(b))
	form.Set("revision", strconv.Itoa(revision))

	uri := fmt.Sprintf(
		"users/%v/playlists/%v/change-relative",
		s.client.userID,
		kind,
	)

	req, err := s.client.NewRequest(http.MethodPost, uri, form)
	if err != nil {
		return nil, nil, err
	}

	addTracksResp := new(PlaylistsAddTracksResp)
	resp, err := s.client.Do(ctx, req, addTracksResp)
	return addTracksResp, resp, err
}

type (
	// PlaylistsRemoveTracksOptions are options for method RemoveTracks
	PlaylistsRemoveTracksOptions struct {
		From int
		To   int
	}
)

// RemoveTracks removes tracks from playlist
func (s *PlaylistsService) RemoveTracks(
	ctx context.Context,
	kind int,
	revision int,
	tracks []PlaylistsTrack,
	opts *PlaylistsRemoveTracksOptions,
) (*PlaylistsRemoveTracksResp, *http.Response, error) {
	if opts == nil {
		opts = &PlaylistsRemoveTracksOptions{
			From: 0,
			To:   len(tracks),
		}
	}

	diff := []struct {
		Op     string           `json:"op"`
		From   int              `json:"from"`
		To     int              `json:"to"`
		Tracks []PlaylistsTrack `json:"tracks"`
	}{
		{
			Op:     "delete",
			From:   opts.From,
			To:     opts.To,
			Tracks: tracks,
		},
	}

	b, err := json.Marshal(diff)
	if err != nil {
		return nil, nil, err
	}

	form := url.Values{}
	form.Set("diff", string(b))
	form.Set("revision", strconv.Itoa(revision))

	uri := fmt.Sprintf(
		"users/%v/playlists/%v/change-relative",
		s.client.userID,
		kind,
	)

	req, err := s.client.NewRequest(http.MethodPost, uri, form)
	if err != nil {
		return nil, nil, err
	}

	addTracksResp := new(PlaylistsRemoveTracksResp)
	resp, err := s.client.Do(ctx, req, addTracksResp)
	return addTracksResp, resp, err
}

func (s *PlaylistsService) DownloadOne(ctx context.Context, kind int) {
	resp, _, _ := s.Get(ctx, 0, kind)
	playlist := resp.Result

	/// Get already loaded tracks (already on yandex disk and file system)
	if len(playlist.Tracks) < 1 {
		logInfo.Println("No tracks in playlist: ", playlist.PlaylistsResult.Title)
		return
	}

	logInfo.Printf("Count tracks in playlist %s: %d", playlist.PlaylistsResult.Title, len(playlist.Tracks))
	logInfo.Printf("Playlist to download: %s", playlist.PlaylistsResult.Title)

	playlist_folder := s.client.config.Output + "/" + playlist.PlaylistsResult.Title
	// Create dirs if they are not exist
	if _, err := os.Stat(playlist_folder + "/tracks"); os.IsNotExist(err) {
		err := os.MkdirAll(playlist_folder+"/tracks", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	if _, err := os.Stat(playlist_folder + "/lyrics"); os.IsNotExist(err) {
		err := os.MkdirAll(playlist_folder+"/lyrics", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	/// Get list directory
	var tracks []Track
	for _, track := range playlist.Tracks {
		tracks = append(tracks, track.Track)
	}

	s.client.tracks.DownloadAll(ctx, tracks, playlist_folder)
}

func (s *PlaylistsService) DownloadAll(ctx context.Context, kinds []int) {

	if len(kinds) < 1 {
		result, _, _ := s.client.Playlists().List(context.Background(), 0)
		playlists := result.Result
		for _, playlist := range playlists {
			kinds = append(kinds, playlist.Kind)
		}
	}
	for _, kind := range kinds {
		s.DownloadOne(ctx, kind)
	}
}

// AddTracks adds tracks to playlist
func (s *PlaylistsService) DistributeTracksByPlaylists() {
	var id int
	var add_tracks []PlaylistsTrack
	playlists_map, _ := config.CreatePlaylistsMap()
	tracks_out_playlist := s.client.GetTracksWithoutPlaylist()
	for _, iter := range playlists_map.Playlists {
		add_tracks = nil
		for _, track := range tracks_out_playlist {
			if slices.Contains(iter.Authors, track.Artists[0].Name) {
				id, _ = strconv.Atoi(track.ID)
				add_tracks = append(add_tracks, PlaylistsTrack{
					ID: id, AlbumID: track.Albums[0].ID,
				})
			}
		}
		res1, _, _ := s.client.Playlists().Get(context.Background(), 0, iter.Kind) // for getting revision
		log.Printf("%d\n", len(add_tracks))
		s.client.Playlists().AddTracks(context.Background(), iter.Kind, res1.Result.Revision, add_tracks, nil)
	}
}

// Add tracks to playlist
func (s *PlaylistsService) AddTracksToPlaylist() {
	// kind 1069 - testing playlist
	track_ids := []int{20345844, 21825973}
	playlists_track := []PlaylistsTrack{}
	for _, track_id := range track_ids {
		track, _, _ := s.client.Tracks().GetOne(context.Background(), track_id)
		id, _ := strconv.Atoi(track.Result[0].ID)
		playlists_track = append(playlists_track, PlaylistsTrack{
			ID: id, AlbumID: track.Result[0].Albums[0].ID,
		})
	}
	res1, _, _ := s.client.Playlists().Get(context.Background(), 0, 1069) // for getting revision
	s.client.Playlists().AddTracks(context.Background(), 1069, res1.Result.Revision, playlists_track, nil)
}

// Delete track from playlists
func (s *PlaylistsService) DeleteTracksFromPlaylists() {
	var delete_tracks []PlaylistsTrack
	var id int
	playlists_map, _ := config.CreatePlaylistsMap()
	for _, playlist := range playlists_map.Playlists {
		delete_tracks = nil
		res1, _, _ := s.client.Playlists().Get(context.Background(), 0, playlist.Kind)
		for _, track := range res1.Result.Tracks {
			id, _ = strconv.Atoi(track.Track.ID)
			// log.Printf("Title: %s, albums len: %d", track.Track.Title, len(track.Track.Albums))
			delete_tracks = append(delete_tracks, PlaylistsTrack{
				ID: id, AlbumID: track.Track.Albums[0].ID,
			})
		}
		s.client.Playlists().RemoveTracks(context.Background(), playlist.Kind, res1.Result.Revision, delete_tracks, nil)
	}
}

// Delete tracks from playlist
func (s *PlaylistsService) DeleteTracksFromPlaylist() {
	track_ids := []int{20345844, 21825973}
	playlists_track := []PlaylistsTrack{}
	for _, track_id := range track_ids {
		track, _, _ := s.client.Tracks().GetOne(context.Background(), track_id)
		id, _ := strconv.Atoi(track.Result[0].ID)
		playlists_track = append(playlists_track, PlaylistsTrack{
			ID: id, AlbumID: track.Result[0].Albums[0].ID,
		})
	}
	res1, _, _ := s.client.Playlists().Get(context.Background(), 0, 1069) // for getting revision
	s.client.Playlists().RemoveTracks(context.Background(), 1069, res1.Result.Revision, playlists_track, nil)
}
