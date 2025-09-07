package config

import (
	"bytes"
	"os"

	"github.com/pyr33x/benchmq/pkg/er"
	"gopkg.in/yaml.v3"
)

// Config represents the entire yaml config file fields
type Config struct {
	Name   string `yaml:"name"`
	Server server `yaml:"server"`
	Client Client `yaml:"client"`
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
func InitializeCfg() (*Config, error) {
	rawCfg, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, &er.Error{
			Package: "Config",
			Func:    "InitializeCfg",
			Message: er.ErrConfigReadFailed,
		}
	}

	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(rawCfg))
	dec.KnownFields(true)
	if err = dec.Decode(&cfg); err != nil {
		return nil, &er.Error{
			Package: "Config",
			Func:    "InitializeCfg",
			Message: er.ErrUnmarshalFailed,
		}
	}

	cfg.SetDefaults()
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
func (c *Config) SetDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "localhost"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 1883
	}
	if c.Client.KeepAlive == 0 {
		c.Client.KeepAlive = 60
	}
}
