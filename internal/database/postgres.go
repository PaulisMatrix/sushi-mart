package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(dbname, user, password string) (*Postgres, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Try connecting to the database
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(30 * time.Second)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("connection timeout")

		case <-ticker.C:
			if err := db.Ping(); err == nil {
				return &Postgres{DB: db}, nil
			}
		}
	}
}
