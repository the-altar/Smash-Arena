package arena

// CHARACTER defines the character struct
type CHARACTER struct {
	Health int     `json:"health"`
	ID     int     `json:"id"`
	Skills []SKILL `json:"skills"`
}

// FLAGS defines the flag struct
type FLAGS struct {
	Counterable bool `json:"counterable"`
	Reflectable bool `json:"reflectable"`
}

// COST defines the cost struct
type COST struct {
	Physical  int `json:"physical"`
	Special   int `json:"special"`
	Strategic int `json:"strategic"`
	Secret    int `json:"secret"`
	Wild      int `json:"wild"`
}

// EFFECT defines the effect struct
type EFFECT struct {
	Type  int `json:"type"`
	Value int `json:"value"`
	Attr  int `json:"attr"`
}

// SKILL define the skill struct
type SKILL struct {
	Usable   bool     `json:"usable"`
	Targets  int      `json:"targets"`
	Flags    FLAGS    `json:"flags"`
	Effects  []EFFECT `json:"effects"`
	Cooldown int      `json:"cooldown"`
	Cost     COST     `json:"cost"`
}

// GAME defines our game struct
type GAME struct {
	Characters [3]CHARACTER
	Enemy      [3]CHARACTER
}

// parameters: isCounterable, isReflectable
func buildFlags(isCounterable bool, isReflectable bool)

// parameters: spirit, strength, wisdom, secret, random
func buildCost(spirit int, strength int, wisdom int, secret int, random int)

// parameters: EffectType, Value, Attribute
func buildEffect(eType int, value int, attr int)

// parameters: isUsable, targets, cooldown
func buildSkill(isUsable bool, targets int, cooldown int, effect []EFFECT, flags FLAGS, cost COST)

// parameters: health, id
// returns a CHARACTER
func buildCharacter(health int, id int) CHARACTER
