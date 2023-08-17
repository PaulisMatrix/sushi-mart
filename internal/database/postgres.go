package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgres(dbname, user, password string) (*Postgres, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", user, password, dbname)
	//db, err := sql.Open("postgres", connStr)
	dbpool, err := pgxpool.New(context.Background(), connStr)

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
			if err := dbpool.Ping(context.Background()); err == nil {
				return &Postgres{DB: dbpool}, nil
			}
		}
	}
}
