package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

func CreateConfig() (*Config, error) {
	// Read and process config file =========================
	configPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
		return &Config{}, err
	}
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return config, err
	}

	fmt.Printf("Auth token: %s\n", config.Awesome.WBAuthToken)
	fmt.Printf("Auth token: %s\n", config.Awesome.YAMusicAuthToken)

	return config, nil
}

// NewConfig returns a new decoded Config struct
func CreatePlaylistsMap() (*PlaylistConfig, error) {
	// Create config structure
	config := &PlaylistConfig{}

	// Open config file
	file, err := os.Open("./playlists_map.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
