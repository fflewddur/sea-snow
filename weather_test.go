package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExpired(t *testing.T) {
	past, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 15:04:05 -0700")
	cc := CurrentConditions{
		RequestTimestamp: past,
	}
	if cc.Expired() != true {
		t.Error("old cc.Expired() = false")
	}

	future, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2106 15:04:05 -0700")
	cc.RequestTimestamp = future
	if cc.Expired() != false {
		t.Error("future cc.Expired() = true")
	}
}

func TestSnowing(t *testing.T) {
	jsonSnowing := `{"coord":{"lon":-122.33,"lat":47.61},"weather":[{"id":615,"main":"Snow","description":"Light rain and snow","icon":"13d"}],"base":"stations","main":{"temp":280.57,"feels_like":276.69,"temp_min":279.15,"temp_max":282.15,"pressure":1020,"humidity":57},"visibility":16093,"wind":{"speed":2.6},"clouds":{"all":20},"dt":1583705816,"sys":{"type":1,"id":3417,"country":"US","sunrise":1583678113,"sunset":1583719494},"timezone":-25200,"id":5809844,"name":"Seattle","cod":200}`
	jsonNotSnowing := `{"coord":{"lon":-122.33,"lat":47.61},"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}],"base":"stations","main":{"temp":280.57,"feels_like":276.69,"temp_min":279.15,"temp_max":282.15,"pressure":1020,"humidity":57},"visibility":16093,"wind":{"speed":2.6},"clouds":{"all":20},"dt":1583705816,"sys":{"type":1,"id":3417,"country":"US","sunrise":1583678113,"sunset":1583719494},"timezone":-25200,"id":5809844,"name":"Seattle","cod":200}`
	
	var cc CurrentConditions
	err := json.Unmarshal([]byte(jsonSnowing), &cc)
	if err != nil {
		t.Error(err.Error())
	}
	if cc.Snowing() != true {
		t.Errorf("jsonSnowing cc.Snowing() = false")
	}

	err = json.Unmarshal([]byte(jsonNotSnowing), &cc)
	if err != nil {
		t.Error(err.Error())
	}
	if cc.Snowing() != false {
		t.Errorf("jsonNotSnowing cc.Snowing() = true")
	}
}