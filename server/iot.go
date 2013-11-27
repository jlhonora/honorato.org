package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
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
