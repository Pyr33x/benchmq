package cmd

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name   string `yaml:"name"`
	Server Server `yaml:"server"`
	Client Client `yaml:"client"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Client struct {
	ClientID     string `yaml:"client_id"`
	KeepAlive    uint16 `yaml:"keep_alive"`
	CleanSession bool   `yaml:"clean_session"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

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
