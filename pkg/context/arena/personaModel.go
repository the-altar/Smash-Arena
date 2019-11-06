package arena

import "github.com/the-altar/Smash-Arena/pkg/config"

// Persona struct
type Persona struct {
	ID      int
	Name    string
	Profile string
}

// AllPersona from database
func allPersona() []Persona {
	p := make([]Persona, 0)

	sql := "SELECT * from persona"
	rows, err := config.DB.Query(sql)

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		p1 := Persona{}
		if err = rows.Scan(&p1.ID, &p1.Name, &p1.Profile); err != nil {
			panic(err)
		}
		p = append(p, p1)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return p
}
