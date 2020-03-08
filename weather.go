package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// CurrentConditions holds the weather conditions response from openweathermap.org
type CurrentConditions struct {
	Name    string `json:"name"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
}

// FetchCurrentConditions attempts to fetch current weather conditions from openweathermap.org
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

	return cc
}
