package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

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
