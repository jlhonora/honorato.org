package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

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

var db *sql.DB
func dbinit() (error) {
	db, err := sql.Open("postgres", "user=pgmainuser dbname=pgmaindb sslmode=disable")
	if err != nil {
		log.Println(err)
	}
}

func dbclose() (error) {
	db.Close()
}

func deleteRows(table_name string, num int) {
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

func insertActivity(activity_type string, body string, target_name string, target_url string, created_at string) (error) {
	_, err := db.Exec(`INSERT INTO activities (activity_type, body, target_name, target_url, created_at) ` +
		`VALUES ('` + activity_type + "', '" +
		body + "', '" +
		target_name + "', '" +
		target_url + "', '" +
		created_at + "')")
	if err != nil {
		fmt.Println("Error inserting")
		fmt.Println(err)
	}
	return err;
}



