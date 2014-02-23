package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func genericHttpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling HTTP request")
	fmt.Println(r.Method)
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("Body: " + string(body))
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "HTTP handler failed")
	}
}
