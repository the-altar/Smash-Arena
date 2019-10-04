package engine

import "fmt"

// Damage is a damaging effect struct
type Damage struct {
	Value    int  `json:"value"`
	Tick     bool `json:"tick"`
	Duration int  `json:"duration"`
}

// Exe executes the effect
func (d Damage) Exe() {
	fmt.Println("HELLO!")
}
