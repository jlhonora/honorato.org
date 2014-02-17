package main

import (
	"database/sql"
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"strings"
	"time"
)

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

// Handles Github activities API request
func githubHandler(w http.ResponseWriter, r *http.Request) {
	body, err := getGithubDbData()
	fmt.Fprint(w, string(body))
	if err != nil {
		fmt.Fprint(w, "Github failed")
	}
}

// Queries the database for github activities and transforms them
// to JSON format
func getGithubDbData() ([]byte, error) {
	fmt.Println("Querying GitHub DB")
	db, err := sql.Open("postgres", "user=pgmainuser dbname=pgmaindb sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	// Select the most recent 10 entries
	rows, err := db.Query("SELECT * FROM activities WHERE activity_type LIKE 'github' ORDER BY created_at DESC LIMIT 10")
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

func updateGithubEvents() (error) {
	body, err := getGithubData()
	json, err := simplejson.NewJson(body)

	arr, err := json.Array()
	if err != nil {
		return err
	}
	deleteRows("activities", len(arr))
	var index int
	for index < len(arr) {
		var total_err error
		body, err := json.GetIndex(index).Get("body").String()
		if err != nil {
			total_err = err
		}
		target_name, err := json.GetIndex(index).Get("target").Get("name").String()
		if err != nil {
			total_err = err
		}
		target_url, err := json.GetIndex(index).Get("target").Get("name_url").String()
		if err != nil {
			total_err = err
		}
		created_at, err := json.GetIndex(index).Get("created_at").String()
		if err != nil {
			total_err = err
		}
		// If there was an error or the body is empty then just return
		if total_err != nil || len(body) == 0 {
			fmt.Println("Error processing entry")
		} else {
			fmt.Println("Event: " + body + " " + created_at)
			insertActivity("github", body, target_name, target_url, created_at)
		}
		// Process the next element
		index++
	}
	return err
}
