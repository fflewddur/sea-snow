package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// SnowData holds a summary of local weather conditions.
type SnowData struct {
	Snowing bool
	Testing bool
}

type JsonResponse struct {
	Success bool `json:"success"`
	Snowing bool `json:"snowing"`
}

func main() {
	cfg := ReadConf()

	conditions := make(chan CurrentConditions)
	update := make(chan bool)

	go UpdateAsNeeded(update, conditions, cfg.WeatherAPIKey)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := currentWeather(r, update, conditions)
		writeIndex(w, data)
	})

	http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
		data := currentWeather(r, update, conditions)
		response := JsonResponse{
			Success: true,
			Snowing: data.Snowing,
		}
		b, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	serverAndPort := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(serverAndPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateAsNeeded handles requests for current weather conditions and returns the latest values it has.
// If needed, it will check for updated conditions from the weather service.
func UpdateAsNeeded(update <-chan bool, conditions chan<- CurrentConditions, apiKey string) {
	var cc CurrentConditions
	for {
		select {
		case <-update:
			if cc.Expired() {
				cc = FetchCurrentConditions(apiKey)
			}
			conditions <- cc
		}
	}
}

func currentWeather(r *http.Request, update chan<- bool, conditions <-chan CurrentConditions) SnowData {
	snowing := false
	testing := false

	params := r.URL.Query()
	if params.Get("testSnow") != "" {
		snowing = true
		testing = true
	} else if params.Get("testNoSnow") != "" {
		testing = true
	} else {
		update <- true
		cc := <-conditions
		snowing = cc.Snowing()
	}

	data := SnowData{
		Snowing: snowing,
		Testing: testing,
	}

	return data
}

func writeIndex(w http.ResponseWriter, data SnowData) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
