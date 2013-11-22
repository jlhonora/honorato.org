package main

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"regexp"
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

var validPath = regexp.MustCompile("^/(edit|save|view|)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
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
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
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
	http.HandleFunc("/github", githubHandler)
	http.HandleFunc("/", indexHandler)
	handleResources()
	log.Fatal(http.ListenAndServe(":7654", nil))
}
