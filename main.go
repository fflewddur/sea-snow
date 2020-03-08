package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// SnowData holds a summary of local weather conditions
type SnowData struct {
	Snowing bool
	Testing bool
}

func main() {
	cfg := ReadConf()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %v\n", r.URL)
		snowing := false
		testing := false

		params := r.URL.Query()
		if params.Get("testSnow") != "" {
			snowing = true
			testing = true
		} else if params.Get("testNoSnow") != "" {
			testing = true
		}

		FetchCurrentConditions(cfg.WeatherAPIKey)

		data := SnowData{
			Snowing: snowing,
			Testing: testing,
		}

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
	})

	serverAndPort := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(serverAndPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
