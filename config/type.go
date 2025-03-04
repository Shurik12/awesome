package config

import "time"

// Config struct for webapp config
type Config struct {
	Awesome struct {
		// Host is the local machine IP Address to bind the HTTP Server to
		Host string `yaml:"host"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		Port string `yaml:"port"`

		// Authentication token to access wildberies api
		WBAuthToken string `yaml:"wb_auth_token"`

		// Authentication token to access yandex music api
		YAMusicAuthToken string `yaml:"yamusic_auth_token"`

		Timeout struct {
			// Server is the general server timeout to use
			// for graceful shutdowns
			Server time.Duration `yaml:"server"`

			// Write is the amount of time to wait until an HTTP server
			// write opperation is cancelled
			Write time.Duration `yaml:"write"`

			// Read is the amount of time to wait until an HTTP server
			// read operation is cancelled
			Read time.Duration `yaml:"read"`

			// Read is the amount of time to wait
			// until an IDLE HTTP session is closed
			Idle time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"awesome"`
}

type PlaylistConfig struct {
	Playlists []struct {
		Title   string   `yaml:"title"`
		Kind    int      `yaml:"kind"`
		Authors []string `yaml:"authors"`
	} `yaml:"playlists"`
}
