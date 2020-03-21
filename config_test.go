package main

import (
	"strings"
	"testing"
)

func TestParseConf(t *testing.T) {
	yml := `# Server configuration
server:
  host: "localhost"
  port: 9990

weather_api_key: "put_your_key_here"
`

	r := strings.NewReader(yml)
	config, err := ParseConf(r)

	if err != nil {
		t.Errorf("got err = %v", err)
	}
	if config.Server.Host != "localhost" {
		t.Errorf("got config.Server.Host = %s; want %s", config.Server.Host, "localhost")
	}
	if config.Server.Port != "9990" {
		t.Errorf("got config.Server.Port = %s; want %s", config.Server.Port, "9990")
	}
	if config.WeatherAPIKey != "put_your_key_here" {
		t.Errorf("got config.WeatherAPIKey = %s; want %s", config.Server.Host, "put_your_key_here")
	}
}

func TestInvalidYaml(t *testing.T) {
	yml := `# Server configuration
server:
  host: "localhost"
  port: 9990
  oops

weather_api_key: "put_your_key_here"
`

	r := strings.NewReader(yml)
	_, err := ParseConf(r)

	if err == nil {
		t.Errorf("got err = nil")
	}
}
