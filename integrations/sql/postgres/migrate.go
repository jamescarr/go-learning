package main

import (
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func failOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	m, err := migrate.New(
		"file://migrations",
		"postgres://test:test@localhost:5432/test?sslmode=disable")
	failOnError(err)

	err = m.Steps(2)
	failOnError(err)
}
