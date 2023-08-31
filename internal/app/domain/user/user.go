package user

import "github.com/gocql/gocql"

type User struct {
	ID        gocql.UUID `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	NIM       string     `json:"nim"`
	Year      string     `json:"year"`
}
