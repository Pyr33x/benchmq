package cmd

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the entire yaml config file fields
type Config struct {
	Name   string `yaml:"name"`
	Server Server `yaml:"server"`
	Client Client `yaml:"client"`
}

// Server represents the server configuration fields
type Server struct {
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
func InitializeCfg() *Config {
	rawCfg, err := os.ReadFile("config.yml")
	if err != nil {
		log.Println("Failed to read the configuration file")
	}

	var cfg Config
	err = yaml.Unmarshal(rawCfg, &cfg)
	if err != nil {
		log.Println("Failed to unmarshal the yaml configuration file")
	}

	return &cfg
}
