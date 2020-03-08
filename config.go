package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config stores runtime configuration options for sea-snow
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
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("decode config.yml: %v", err.Error())
	}

	return cfg
}
