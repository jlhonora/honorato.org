package main

import (
	auth "github.com/abbot/go-http-auth"
	"log"
)

func main() {
	// pw, salt, magic
	log.Printf(string(auth.MD5Crypt([]byte("hello"), []byte("J.w5a/.."), []byte("$apr1$"))))
}
