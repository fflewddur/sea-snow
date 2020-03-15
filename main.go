package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// SnowData holds a summary of local weather conditions.
type SnowData struct {
	Snowing bool
	Testing bool
}

// JSONResponse holds our... JSON response.
type JSONResponse struct {
	Success bool `json:"success"`
	Snowing bool `json:"snowing"`
}

func main() {
	cfg := ReadConf()

	test := flag.Bool("t", false, "test mode")

	flag.Parse()

	if *test {
		log.Printf("Running in test mode")
		rand.Seed(time.Now().UnixNano())
	}
	conditions := make(chan CurrentConditions)
	update := make(chan bool)

	go fetchConditionsAsNeeded(update, conditions, cfg.WeatherAPIKey)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := currentWeather(r, update, conditions)
		if *test {
			data.Testing = true
			data.Snowing = testSnowing()
		}
		writeIndex(w, data)
	})

	http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
		data := currentWeather(r, update, conditions)
		response := JSONResponse{
			Success: true,
			Snowing: data.Snowing,
		}
		if *test {
			response.Snowing = testSnowing()
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

// fetchConditionsAsNeeded handles requests for current weather conditions and returns the latest values it has.
// If needed, it will check for updated conditions from the weather service.
func fetchConditionsAsNeeded(update <-chan bool, conditions chan<- CurrentConditions, apiKey string) {
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

func testSnowing() bool {
	n := rand.Intn(99)
	return n < 50
}
