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
		SPECIAL   int
		STRATEGIC int
		PHYSICAL  int
		SECRET    int
		WILD      int
	}
	// ATTRIBUTE enum definition
	ATTRIBUTE struct {
		HEALTH int
	}
)

var (
	target     = TARGET{ONE: 1, TWO: 2, ALL: 3, RANDOM: 0}
	effectType = TYPE{DAMAGE: 1}
	cost       = COSTINDEX{SPECIAL: 0, STRATEGIC: 1, PHYSICAL: 2, SECRET: 3, WILD: 4}
	attr       = ATTRIBUTE{HEALTH: 0}
)
