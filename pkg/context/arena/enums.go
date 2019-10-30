package arena

type (
	// TARGET enum definition
	TARGET struct {
		ONE    int
		TWO    int
		ALL    int
		RANDOM int
	}
	// TYPE is the effect enum definition
	TYPE struct {
		DAMAGE int
	}
	// COSTINDEX enum definition
	COSTINDEX struct {
		SPIRIT   int
		STRENGTH int
		WISDOM   int
		SECRET   int
		RANDOM   int
	}
	// ATTRIBUTE enum definition
	ATTRIBUTE struct {
		HEALTH int
	}
)

var (
	target     = TARGET{ONE: 1, TWO: 2, ALL: 3, RANDOM: 0}
	effectType = TYPE{DAMAGE: 1}
	cost       = COSTINDEX{SPIRIT: 0, STRENGTH: 1, WISDOM: 2, SECRET: 3, RANDOM: 4}
	attr       = ATTRIBUTE{HEALTH: 0}
)
