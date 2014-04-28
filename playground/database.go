package main

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
)

var DB *sql.DB

func dbinit() error {
	db, err := sql.Open("postgres", "user=pgtestuser dbname=pgtestdb sslmode=disable")
	if err != nil {
		log.Println(err)
	} else {
		DB = db
	}
	return err
}

func dbclose() error {
	DB.Close()
	return nil
}

func getPassword(username string) string {
	var password string
	log.Printf("Querying password for " + username)
	err := DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&password)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Password is %s\n", password)
	}
	return password
}
