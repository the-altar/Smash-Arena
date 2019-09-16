package engine

import "fmt"

type skill struct {
	id      int
	name    string
	effects []effect
}

func (s *skill) build(id int, name string, effects []effect) {
	s.id = id
	s.name = name
	for i := 0; i < len(effects); i++ {
		s.effects = append(s.effects, effects[i])
	}
	fmt.Println(s)
}
