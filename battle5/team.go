package battle5

type Team struct {
	Heros     *HeroList
	TeamIndex int
}

func (team *Team) AutoSetPos() {
	lst := [][]*Hero{
		{},
		{},
		{},
	}

	for _, v := range team.Heros.Heros {
		lst[v.Props[PropTypePlace]-1] = append(lst[v.Props[PropTypePlace]-1], v)
	}

	for _, arr := range lst {
		if len(arr) == 3 {
			for i, v := range arr {
				v.SX = v.Props[PropTypePlace]
				v.SY = i + 1
			}
		} else if len(arr) == 2 {
			arr[0].SX = arr[0].Props[PropTypePlace]
			arr[0].SY = 1

			arr[1].SX = arr[1].Props[PropTypePlace]
			arr[1].SY = 1
		} else if len(arr) == 1 {
			arr[0].SX = arr[0].Props[PropTypePlace]
			arr[0].SY = 2
		}
	}
}

func NewTeam(index int, lst []*HeroData) *Team {
	team := &Team{
		Heros:     NewHeroList(),
		TeamIndex: index,
	}

	team.Heros.Init(lst)

	for _, v := range team.Heros.Heros {
		v.TeamIndex = index
	}

	return team
}
