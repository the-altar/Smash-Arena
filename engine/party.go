package engine

type party struct {
	player string
	team   [3]character
}

func (p party) getTeam() [3]character {
	return p.team
}

func (p *party) form(ids [3]string, player string) {
	var members [3]character

	for i := 0; i < 3; i++ {
		members[i].build(ids[i])
	}

	p.player = player
	p.team = members
}
