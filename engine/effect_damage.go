package engine

import "fmt"

// Damage is a damaging effect struct
type Damage struct {
	Value    int
	Tick     bool
	Duration int
}

// Exe executes the effect
func (d Damage) Exe() {
	fmt.Println("HELLO!")
}
