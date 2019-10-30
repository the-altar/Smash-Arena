package account

import (
	"database/sql"
	"fmt"

	"github.com/the-altar/Smash-Arena/pkg/config"
)

// Account defines a user struct
type Account struct {
	Username string
	Password string
	ID       int
}

// OneAccountByName finds a user
func OneAccountByName(name string) (Account, error) {
	user := Account{}

	sqlStatement := "SELECT * FROM account WHERE account.username = $1"

	row := config.DB.QueryRow(sqlStatement, name)

	switch err := row.Scan(&user.ID, &user.Username, &user.Password); err {
	case sql.ErrNoRows:
		return user, fmt.Errorf("No rows were found")
	case nil:
		return user, nil
	default:
		panic(err)
	}
}

// OneAccountByID fetches a single user from our database using its ID
func OneAccountByID(id int) (Account, error) {
	user := Account{}

	sqlStatement := "SELECT * FROM account WHERE account.account_key = $1"

	row := config.DB.QueryRow(sqlStatement, id)

	switch err := row.Scan(&user.ID, &user.Username, &user.Password); err {
	case sql.ErrNoRows:
		return user, fmt.Errorf("No rows were found")
	case nil:
		return user, nil
	default:
		panic(err)
	}
}

// CreateAccount creates a new user
func CreateAccount(u string, p string) error {
	sqlStatement := "INSERT INTO account(username, keyword) VALUES ($1,$2)"

	if _, err := config.DB.Query(sqlStatement, u, p); err != nil {
		return err
	}
	return nil
}
