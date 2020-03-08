package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// UpdateInterval is the time in minutes to wait before refreshing weather data.
const UpdateInterval = time.Minute * 6

// CurrentConditions holds the weather conditions response from openweathermap.org.
type CurrentConditions struct {
	Name    string `json:"name"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	DataTimestamp    int       `json:"dt"` // UNIX timestamp, UTC
	RequestTimestamp time.Time `json:"-"`
}

// Expired returns true if we should check for updated weather conditions.
func (cc *CurrentConditions) Expired() bool {
	now := time.Now()
	expired := now.After(cc.RequestTimestamp.Add(UpdateInterval))
	return expired
}

// Snowing returns true if it appears to be snowing, and false otherwise.
func (cc *CurrentConditions) Snowing() bool {
	if cc != nil && len(cc.Weather) > 0 {
		s := cc.Weather[0].Main
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		return s == "snow"
	}

	return false
}

// FetchCurrentConditions attempts to fetch current weather conditions from openweathermap.org.
func FetchCurrentConditions(apiKey string) CurrentConditions {
	url := "https://api.openweathermap.org/data/2.5/weather?id=5809844&appid=" + apiKey

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	var cc CurrentConditions
	err = json.NewDecoder(res.Body).Decode(&cc)
	if err != nil {
		log.Println(err.Error())
	}
	cc.RequestTimestamp = time.Now()

	return cc
}
