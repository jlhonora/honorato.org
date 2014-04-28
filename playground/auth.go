package main

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"log"
	"net/http"
)

func Secret(user, realm string) string {
	return getPassword(user)
}

func handle(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "<html><body><h1>Hello, %s!</h1></body></html>", r.Username)
}

func main() {
	log.Printf("Init")
	dbinit()
	defer dbclose()
	authenticator := auth.NewBasicAuthenticator("127.0.0.1", Secret)
	http.HandleFunc("/", authenticator.Wrap(handle))
	log.Printf("Ready")
	http.ListenAndServe(":8080", nil)
}
