// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"database/sql"
)

type Customer struct {
	ID       int32
	Username string
	Password string
	Email    string
	Phone    sql.NullString
	Address  sql.NullString
}
