package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func iotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling IoT")
	if r.Method == "POST" {
		fmt.Println("Post")
		// receive posted data
		body, err := ioutil.ReadAll(r.Body)
		fmt.Println("Body: " + string(body))
		if err != nil {
			fmt.Println(err)
			fmt.Fprint(w, "IoT failed")
		}
	}
}
