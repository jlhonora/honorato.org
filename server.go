package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
	"fmt"
	"time"
	"strings"
	"strconv"
	"log"
	"database/sql"
	"github.com/bitly/go-simplejson"
	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseFiles("index.html"))

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

func formatGithubEvent(event *simplejson.Json) ([]byte, error) {
	var msg_json []byte
	var target_name, target_url, body, msg string
	event_type, _ := event.Get("type").String()
	created_at, _ := event.Get("created_at").String()
	switch event_type {
		case "FollowEvent":
			target_name, _ = event.Get("payload").Get("target").Get("login").String()
			target_url, _ = event.Get("payload").Get("target").Get("html_url").String()
			body = "Followed"
			break
		case "PushEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url = "https://github.com/" + target_name
			body = "Pushed to"
			break
		case "PullRequestEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url, _ = event.Get("payload").Get("pull_request").Get("html_url").String()
			action, _ := event.Get("payload").Get("action").String()
			body = action + " pull request on"
			break
		case "IssueCommentEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url, _ = event.Get("payload").Get("issue").Get("html_url").String()
			body = "Commented on"
			break
		case "ReleaseEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url, _ = event.Get("payload").Get("release").Get("html_url").String()
			ref, _ := event.Get("payload").Get("release").Get("name").String()
			ref_type, _ := event.Get("payload").Get("action").String()
			body = ref_type + " " + ref + " on"
			break
		case "CreateEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url = "https://github.com/" + target_name
			ref, _ := event.Get("payload").Get("ref").String()
			ref_type, _ := event.Get("payload").Get("ref_type").String()
			body = "Created "
			if ref_type != "" {
				body += ref_type + " "
			}
			if ref != "" {
				body += ref + " "
			}
			break
		case "ForkEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url, _ = event.Get("payload").Get("forkee").Get("html_url").String()
			body = "Forked"
			break
		case "PullRequestReviewCommentEvent":
			target_name, _ = event.Get("repo").Get("name").String()
			target_url, _ = event.Get("payload").Get("comment").Get("html_url").String()
			body = "Commented on pull request for"
			break
		default:
			str, _ := event.Get("type").String()
			fmt.Print("Unknown type " + str)
	}
	msg_json = []byte(`{
		"body": "` + body + `",
		"target": {
		  "name": "` + target_name + `",
		  "name_url": "` + target_url + `"
		},
		"created_at": "` + created_at + `"
	}`)

	fmt.Print("Json format:")
	fmt.Print(string(msg_json))
	msg = created_at + " - " + body + " " + target_name
	fmt.Println("\nMessage: " + msg)
	return msg_json, nil
}

func getGithubData() ([]byte, error) {
	fmt.Println("Getting latest events from github")
	_, body, err := getContent("https://api.github.com/users/jlhonora/events?page=1&per_page=10")
	var all_events_json []string
	if err != nil {
		fmt.Println("Error retrieving the events")
		fmt.Println(err)
	} else {
		fmt.Println("Successfully retrieved events")
		fmt.Println(string(body))
		json, err := simplejson.NewJson(body)
		if err != nil {
			fmt.Println("error in NewJson:", err)
		}
		fmt.Println("Events:")
		//fmt.Printf("\t%+v", json)
		// index is the index where we are
		// element is the element from someSlice for where we are
		arr, err := json.Array()
		var index int
		for index < len(arr) {
			//fmt.Printf("Formatting event %d\n", index)
			event_json, _ := formatGithubEvent(json.GetIndex(index))
			all_events_json = append(all_events_json, string(event_json))
			index++
		}
	}
	json_str := "[" + strings.Join(all_events_json, ",") + "]"
	return []byte(json_str), err
}

func githubHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "This is the index")
	//body, err := getGithubData()
	body, err := getGithubDbData()
	fmt.Fprint(w, string(body))
	if err != nil {
		fmt.Fprint(w, "Github failed")
	}
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

func getGithubDbData() ([]byte, error) {
	fmt.Println("Querying GitHub DB")
	db, err := sql.Open("postgres", "user=pgmainuser dbname=pgmaindb sslmode=disable")
	if err != nil {
		log.Println(err)
	}
    rows, err := db.Query("SELECT * FROM activities WHERE activity_type LIKE 'github' LIMIT 10")
    if err != nil {
            log.Println(err)
    }
	var all_events_json []string
	var msg_json string
    for rows.Next() {
		var id int
		var activity_type string
		var body string
		var target_name string
		var target_url string
		var created_at time.Time
		fmt.Println("Scanning row")
		if err := rows.Scan(&id, &activity_type, &body, &target_name, &target_url, &created_at); err != nil {
			log.Println(err)
		}
		msg_json = `{
			"body": "` + body + `",
			"target": {
			  "name": "` + target_name + `",
			  "name_url": "` + target_url + `"
			},
			"created_at": "` + created_at.Format(time.RFC3339) + `"
		}`
		fmt.Printf("JSON: %s\n", msg_json)
		all_events_json = append(all_events_json, msg_json)
    }
    if err := rows.Err(); err != nil {
        log.Println(err)
    }

	json_str := "[" + strings.Join(all_events_json, ",") + "]"
	return []byte(json_str), err
}

func dbTest() {
	fmt.Println("Testing db")
	db, err := sql.Open("postgres", "user=pgmainuser dbname=pgmaindb sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
    rows, err := db.Query("SELECT * FROM activities LIMIT 10")
    if err != nil {
            log.Fatal(err)
    }
    for rows.Next() {
            var activity_type string
            var body string
            var target_name string
            var target_url string
            var created_at time.Time
            if err := rows.Scan(&activity_type, &body, &target_name, &target_url, &created_at); err != nil {
                    log.Fatal(err)
            }
            fmt.Printf("Body is %s\n", body)
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }
	fmt.Println("Testing insert")
	_, err = db.Exec(`INSERT INTO activities (activity_type, body, target_name, target_url, created_at) ` +
					 `VALUES (` +
						"'github', '" +
						"test_body" +
						"', '" + "target_name" +
						"', '" + "target_url" +
						"', '" + time.Now().Format(time.RFC3339) +
					  "')")
	if err != nil {
		fmt.Println("Error inserting")
		fmt.Println(err)
	} else {
		fmt.Println("Insert OK")
	}
}

func deleteRows(db *sql.DB, table_name string, num int) {
	fmt.Println("Deleting " + strconv.Itoa(num) + " rows from " + table_name)
	_, err := db.Exec(`DELETE FROM ` + table_name +
						` WHERE ctid IN (
							SELECT ctid
							FROM ` + table_name +
							` ORDER BY id` +
							` LIMIT ` + strconv.Itoa(num) +
						`)`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Delete OK")
	}
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
