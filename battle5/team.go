package battle5

type Team struct {
	Heros        *HeroList
	TeamIndex    int
	Speed        int
	SpeedVal     int
	battle       *Battle
	needUpdAlive bool // 是否需要刷新alive状态，一般来说，如果有角色阵亡，就需要设置这个状态
	isAlive      bool // 是否存活
}

func (team *Team) IsAlive() bool {
	if team.needUpdAlive {
		alivenum := 0

		team.Heros.ForEach(func(h *Hero) {
			if h.IsAlive() {
				alivenum++
			}
		})

		if alivenum > 0 {
			team.isAlive = true
		} else {
			team.isAlive = false
		}
	}

	return team.isAlive
}

func (team *Team) SetTeamSpeedVal(speedval int) {
	team.SpeedVal = speedval

	team.Heros.ForEach(func(h *Hero) {
		if h.IsAlive() {
			h.Props[PropTypeTeamSpeedVal] = speedval
		}
	})
}

func (team *Team) CountSpeed() int {
	team.Speed = 0

	team.Heros.ForEach(func(h *Hero) {
		if h.IsAlive() {
			team.Speed += h.Props[PropTypeCurSpeed]
		}
	})

	return team.Speed
}

func (team *Team) AutoSetPos() {
	// 算法上，按前中后排处理，如果后排人数多于3个，则把最快的放中排，如果中排人数多与3个则将最快的中排放前排
	// 如果前排人数多于3个，则速度慢的后移
	// 如果一开始中排人数多于3个，则需要判断前后排情况，根据需要前移或后移
	lst := []*HeroList{
		NewHeroList(),
		NewHeroList(),
		NewHeroList(),
	}

	for _, v := range team.Heros.Heros {
		lst[v.Props[PropTypePlace]-1].AddHero(v)
	}

	// 因为一共就上场5个人
	if len(lst[0].Heros) > 3 {
		lst[0].SortInAutoSetPos()

		for i := 3; i < len(lst[0].Heros); i++ {
			lst[1].AddHero(lst[0].Heros[i])
		}

		lst[0].Heros = lst[0].Heros[0:3]

		// if len(lst[1].Heros) > 3 {
		// 	lst[1].SortInAutoSetPos()

		// 	for i := 3; i < len(lst[1].Heros); i++ {
		// 		lst[2].AddHero(lst[1].Heros[i])
		// 	}

		// 	lst[1].Heros = lst[1].Heros[0:3]
		// }
	} else if len(lst[2].Heros) > 3 {
		lst[2].SortInAutoSetPosSlow()

		for i := 3; i < len(lst[2].Heros); i++ {
			lst[1].AddHero(lst[2].Heros[i])
		}

		lst[2].Heros = lst[2].Heros[0:3]

		// if len(lst[1].Heros) > 3 {
		// 	lst[1].SortInAutoSetPosSlow()

		// 	for i := 3; i < len(lst[1].Heros); i++ {
		// 		lst[0].AddHero(lst[1].Heros[i])
		// 	}

		// 	lst[1].Heros = lst[1].Heros[0:3]
		// }
	} else if len(lst[1].Heros) > 3 {
		if len(lst[1].Heros) == 5 {
			lst[1].SortInAutoSetPos()

			lst[0].AddHero(lst[1].Heros[0])
			lst[2].AddHero(lst[1].Heros[4])

			lst[1].Heros = lst[1].Heros[1:4]
		} else {
			if len(lst[0].Heros) == 0 {
				lst[1].SortInAutoSetPos()

				lst[0].AddHero(lst[1].Heros[0])

				lst[1].Heros = lst[1].Heros[1:]
			} else {
				lst[2].SortInAutoSetPosSlow()

				lst[2].AddHero(lst[1].Heros[0])

				lst[1].Heros = lst[1].Heros[1:]
			}
		}
	}

	for x, arr := range lst {
		arr.SortInAutoSetPos()

		if len(arr.Heros) == 3 {
			for i, v := range arr.Heros {
				v.StaticPos.X = x + 1
				v.StaticPos.Y = i + 1
			}
		} else if len(arr.Heros) == 2 {
			arr.Heros[0].StaticPos.X = x + 1
			arr.Heros[0].StaticPos.Y = 1

			arr.Heros[1].StaticPos.X = x + 1
			arr.Heros[1].StaticPos.Y = 1
		} else if len(arr.Heros) == 1 {
			arr.Heros[0].StaticPos.X = x + 1
			arr.Heros[0].StaticPos.Y = 2
		}
	}
}

func NewTeam(battle *Battle, index int, lst []*HeroData) *Team {
	team := &Team{
		Heros:        NewHeroList(),
		TeamIndex:    index,
		battle:       battle,
		needUpdAlive: true,
	}

	team.Heros.Init(battle, lst)

	for _, v := range team.Heros.Heros {
		v.TeamIndex = index
	}

	return team
}
