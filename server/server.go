package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
	"fmt"
	"time"
	"log"
	"database/sql"
	"github.com/bitly/go-simplejson"
	_ "github.com/lib/pq"
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

// This function fetch the content of a URL will return it as an
// array of bytes if retrieved successfully.
func getContent(url string) ([]byte, []byte, error) {
    // Build the request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      return nil, nil, err
    }
    // Send the request via a client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      return nil, nil, err
    }
    // Defer the closing of the body
    defer resp.Body.Close()
    // Read the content into a byte array
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return nil, nil, err
    }
	// Read the header into a byte array
	/*defer resp.Header.Close()
	header, err := ioutil.ReadAll(resp.Header)
	if err != nil {
		return nil, nil, err
	}*/
    // At this point we're done - simply return the bytes
    return nil, body, nil
}

func handleResources() {
	fmt.Println("Handling resources")
	http.Handle("/dist/",
		http.StripPrefix("/dist/",
		http.FileServer(http.Dir("dist"))))
	http.Handle("/dist/css/",
		http.StripPrefix("/dist/css/",
		http.FileServer(http.Dir("dist/css"))))
	http.Handle("/dist/js/",
		http.StripPrefix("/dist/js/",
		http.FileServer(http.Dir("dist/js"))))
	http.Handle("/dist/fonts/",
		http.StripPrefix("/dist/fonts/",
		http.FileServer(http.Dir("dist/fonts"))))
	http.Handle("/static/",
		http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	fmt.Println("Done handling resources")
}

func periodic() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
		   select {
			case <- ticker.C:
				db, err := sql.Open("postgres", "user=pgmainuser dbname=pgmaindb sslmode=disable")
				if err != nil {
					log.Fatal(err)
				}
				// do stuff
				fmt.Println("Tick")
				body, err := getGithubData()
				json, err := simplejson.NewJson(body)

				arr, err := json.Array()
				if err != nil {
					return
				}
				deleteRows(db, "activities", len(arr))
				var index int
				for index < len(arr) {
					var total_err error
					body, err := json.GetIndex(index).Get("body").String()
					if err != nil { total_err = nil }
					target_name, err := json.GetIndex(index).Get("target").Get("name").String()
					if err != nil { total_err = nil }
					target_url, err := json.GetIndex(index).Get("target").Get("name_url").String()
					if err != nil { total_err = nil }
					created_at, err := json.GetIndex(index).Get("created_at").String()
					if err != nil { total_err = nil }
					if total_err != nil {
						continue;
					}
					fmt.Println("Event: " + body + " " + created_at)
					_, err = db.Exec(`INSERT INTO activities (activity_type, body, target_name, target_url, created_at) ` +
							`VALUES (` + "'github', '" +
								body		+ "', '" +
								target_name + "', '" +
								target_url	+ "', '" +
								created_at	+	 "')")
					if err != nil {
						fmt.Println("Error inserting")
						fmt.Println(err)
					}
					index++
				}

			case <- quit:
				ticker.Stop()
				return
			}
		}
	 }()
}

func main() {
	periodic()
    http.HandleFunc("/github", githubHandler)
    http.HandleFunc("/", indexHandler)
	handleResources()
	log.Fatal(http.ListenAndServe(":7654", nil))
}
