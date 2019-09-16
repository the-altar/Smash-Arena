package engine

type damage struct {
	value   int
	targets []character
}

func (d *damage) construct(value int, targets []character) {
	d.targets = targets
	d.value = value
}

func (d damage) getTargets() []character {
	return d.targets
}

func (d *damage) setTarget(targets []character) {
	d.targets = targets
}

func (d damage) getDamage() int {
	return d.value
}

func (d *damage) setDamage(value int) {
	d.value = value
}

func (d damage) apply() {
	for i := 0; i < len(d.targets); i++ {
		d.targets[i].takeDamage(d.value)
	}
}
