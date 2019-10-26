package user

import (
	"database/sql"
	"fmt"

	"github.com/the-altar/Smash-Arena/pkg/config"
)

// User defines a user struct
type User struct {
	username string
	password string
	ID       int
}

// OneUserByName finds a user
func OneUserByName(name string) (User, error) {
	user := User{}

	sqlStatement := "SELECT * FROM public.user as u WHERE u.user_name = $1"

	row := config.DB.QueryRow(sqlStatement, name)

	switch err := row.Scan(user.ID, user.username, user.password); err {
	case sql.ErrNoRows:
		return user, fmt.Errorf("No rows were found")
	case nil:
		return user, nil
	default:
		panic(err)
	}
}

// CreateUser creates a new user
func CreateUser(u string, p string) error {
	sqlStatement := "INSERT INTO public.user(user_name, user_password) VALUES ($1,$2)"

	if _, err := config.DB.Query(sqlStatement, u, p); err != nil {
		return err
	}
	return nil
}
