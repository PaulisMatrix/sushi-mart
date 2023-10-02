package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrationsTest() {
	db, err := sql.Open("postgres", "postgres://sushimart:sushimartpass@localhost/migtesting?sslmode=disable")
	if err != nil {
		log.Fatalln("failed to open db conn", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln("failed to init db instance", err)

	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalln("failed to apply the migrations", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatalln("failed to up the migration", err)
	}

	err = m.Drop()
	if err != nil {
		log.Fatalln("failed to drop the db", err)
	}

}
