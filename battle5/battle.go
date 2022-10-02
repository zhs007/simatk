package battle5

type Battle struct {
	Scene *Scene
	Teams []*Team
}

func (battle *Battle) SetTeam(index int, lst []*HeroData, autoSetPos bool) {
	battle.Teams[index] = NewTeam(index, lst)

	if autoSetPos {
		battle.Teams[index].AutoSetPos()
	}

	for _, v := range battle.Teams[index].Heros.Heros {
		battle.Scene.AddHero(v)
	}
}

func NewBattle(w, h int) *Battle {
	scene := NewScene(w, h)
	battle := &Battle{
		Scene: scene,
		Teams: []*Team{nil, nil},
	}

	return battle
}
