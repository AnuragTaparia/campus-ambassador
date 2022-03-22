package data

import(
	"database/sql"
)

type User struct{
	Email string
	UserID string
	Password string
	Phone sql.NullString
	Address sql.NullString
	Name sql.NullString
}