package providers

import (
	"database/sql"
	"fmt"
	"os"
)

// DB is needed in a global scope
var DB DatabaseProvider

// DatabaseProvider will establish a connection to our database and manage it
type DatabaseProvider struct {
	driver *sql.DB
}

// Open connection to database and handles error
func (d *DatabaseProvider) Open() {
	var err error
	d.driver, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
}

// FindUserByName by its username
func (d *DatabaseProvider) FindUserByName(name string) (string, error) {
	var password string

	sqlStatement := "SELECT u.user_password FROM public.user as u WHERE u.user_name = $1"

	row := d.driver.QueryRow(sqlStatement, name)

	switch err := row.Scan(&password); err {
	case sql.ErrNoRows:
		return password, fmt.Errorf("No rows were found")
	case nil:
		return password, nil
	default:
		panic(err)
	}
}

// CreateUser creates a new user
func (d *DatabaseProvider) CreateUser(u string, p string) {
	sqlStatement := "INSERT INTO public.user(user_name, user_password) VALUES ($1,$2)"

	if _, err := d.driver.Query(sqlStatement, u, p); err != nil {
		panic(err)
	}
}
