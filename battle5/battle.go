package battle5

type Battle struct {
	Scene *Scene
	Teams []*Team
}

func (battle *Battle) SetTeam(index int, lst []*HeroData, autoSetPos bool) {
	battle.Teams[index] = NewTeam(lst)
	battle.Teams[index].AutoSetPos()
}

func NewBattle(w, h int) *Battle {
	scene := NewScene(w, h)
	battle := &Battle{
		Scene: scene,
		Teams: []*Team{nil, nil},
	}

	return battle
}
