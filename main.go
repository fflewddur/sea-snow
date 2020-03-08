package main

import (
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
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		snowing := false
		testing := false

		params := r.URL.Query()
		if params.Get("testSnow") != "" {
			snowing = true
			testing = true
		} else if params.Get("testNoSnow") != "" {
			testing = true
		}

		data := SnowData{
			Snowing: snowing,
			Testing: testing,
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})

	err := http.ListenAndServe(":9990", nil)
	if err != nil {
		log.Fatal(err)
	}
}
