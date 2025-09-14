package config

import (
	"bytes"
	"os"

	"github.com/pyr33x/benchmq/pkg/er"
	"gopkg.in/yaml.v3"
)

// Config represents the entire yaml config file fields
type Config struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Server      server `yaml:"server"`
	Client      Client `yaml:"client"`
}

// Server represents the server configuration fields
type server struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

// Client represents the client configuration fields
type Client struct {
	ClientID     string `yaml:"client_id"`
	KeepAlive    uint16 `yaml:"keep_alive"`
	CleanSession bool   `yaml:"clean_session"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

// InitializeCfg reads the config file and returns a pointer to the Config struct
// If config.yml doesn't exist, it returns a config with default values
func InitializeCfg() (*Config, error) {
	var cfg Config
	configFileExists := false

	// Try to read config file, but don't fail if it doesn't exist
	rawCfg, err := os.ReadFile("config.yml")
	if err == nil {
		configFileExists = true
		// Config file exists, parse it
		dec := yaml.NewDecoder(bytes.NewReader(rawCfg))
		dec.KnownFields(true)
		if err = dec.Decode(&cfg); err != nil {
			return nil, &er.Error{
				Package: "Config",
				Func:    "InitializeCfg",
				Message: er.ErrUnmarshalFailed,
			}
		}
	} else if !os.IsNotExist(err) {
		// File exists but we can't read it (permissions, etc.)
		return nil, &er.Error{
			Package: "Config",
			Func:    "InitializeCfg",
			Message: er.ErrConfigReadFailed,
		}
	}
	// If config file doesn't exist, cfg remains with zero values which will be filled by SetDefaults()

	cfg.SetDefaults(configFileExists)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate does the validation over server configuration
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return &er.Error{
			Package: "Config",
			Func:    "Validate",
			Message: er.ErrEmptyServerHost,
		}
	}
	if c.Server.Port == 0 {
		return &er.Error{
			Package: "Config",
			Func:    "Validate",
			Message: er.ErrInvalidServerPort,
		}
	}
	return nil
}

// SetDefaults sets the default values to fields when values are not acceptable
// configFileExists indicates whether a config file was successfully loaded
func (c *Config) SetDefaults(configFileExists bool) {
	if c.Name == "" {
		c.Name = "BenchMQ"
	}
	if c.Version == "" {
		c.Version = "1.0.0"
	}
	if c.Environment == "" {
		c.Environment = "development"
	}
	if c.Server.Host == "" {
		c.Server.Host = "localhost"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 1883
	}
	if c.Client.ClientID == "" {
		c.Client.ClientID = "benchmq-client"
	}
	if c.Client.KeepAlive == 0 {
		c.Client.KeepAlive = 60
	}
	// Only set CleanSession default when no config file exists
	// If config file exists, respect the explicit value (even if false)
	if !configFileExists && !c.Client.CleanSession {
		c.Client.CleanSession = true
	}
}
