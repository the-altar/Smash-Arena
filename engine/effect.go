package engine

// effect is an interface for skill effects
type effect interface {
	apply()
}

func execute(e effect) {
	return
}
