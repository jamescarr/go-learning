package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func failOnError(e error) {
	if e != nil {
		panic(e)
	}
}

type User struct {
	UserId int
	Name   string
	Email  string
}

func main() {
	db, err := sql.Open("postgres", "user=test password=test dbname=test sslmode=disable")
	failOnError(err)

	db.QueryRow(`DELETE FROM users`)

	// Insert
	var userid int
	err = db.QueryRow(`INSERT INTO users(name, email)
		VALUES('james', 'james@example.com') RETURNING user_id`).Scan(&userid)
	log.Println("James inserted with user_id", userid)
	failOnError(err)

	rows, err := db.Query(`SELECT user_id, name, email FROM users`)
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.Name, &user.Email)
		failOnError(err)
		users = append(users, user)
	}

	log.Println(users)
}
