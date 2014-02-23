package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"time"
)

var templates = template.Must(template.ParseFiles("../index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling index")
	renderTemplate(w, "index")
	fmt.Println("Done")
}

func handleResources() {
	fmt.Println("Handling resources")
	http.Handle("/dist/",
		http.StripPrefix("/dist/",
			http.FileServer(http.Dir("../dist"))))
	http.Handle("/dist/css/",
		http.StripPrefix("/dist/css/",
			http.FileServer(http.Dir("../dist/css"))))
	http.Handle("/dist/js/",
		http.StripPrefix("/dist/js/",
			http.FileServer(http.Dir("../dist/js"))))
	http.Handle("/dist/fonts/",
		http.StripPrefix("/dist/fonts/",
			http.FileServer(http.Dir("../dist/fonts"))))
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("../static"))))
	fmt.Println("Done handling resources")
}

func periodic() {
	ticker := time.NewTicker(120 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Tick")
				updateGithubEvents()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	dbinit()
	defer dbclose()
	periodic()
	http.HandleFunc("/api/github", githubHandler)
	http.HandleFunc("/api/iot", iotHandler)
	http.HandleFunc("/broadcast", genericHttpHandler)
	http.HandleFunc("/", indexHandler)
	handleResources()
	log.Fatal(http.ListenAndServe(":7654", nil))
}
