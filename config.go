package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config stores runtime configuration options for sea-snow.
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	WeatherAPIKey string `yaml:"weather_api_key"`
}

// ReadConf attempts to read a `config.yml` file in the working directory.
// If it cannot do so, it returns a default configuration.
func ReadConf() Config {
	path := "config.yml"
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	cfg := ParseConf(f)
	return cfg
}

// ParseConf parses the YAML configuration data in r.
func ParseConf(r io.Reader) Config {
	var cfg Config
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("decode config.yml: %v", err.Error())
	}

	return cfg
}
