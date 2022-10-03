package battle5

import "sort"

type Battle struct {
	Scene     *Scene
	mapTeams  map[int]*Team
	Log       *BattleLog
	curHeroID int
	mapHeros  map[int]*Hero
}

func (battle *Battle) GenRealHeroID() int {
	id := battle.curHeroID

	battle.curHeroID++

	return id
}

func (battle *Battle) SetTeam(index int, lst []*HeroData, autoSetPos bool) {
	battle.mapTeams[index] = NewTeam(battle, index, lst)

	if autoSetPos {
		battle.mapTeams[index].AutoSetPos()
	}

	for _, v := range battle.mapTeams[index].Heros.Heros {
		battle.Scene.AddHero(v)
	}
}

func (battle *Battle) GenCurHeroList() *HeroList {
	teams := []*Team{}
	for _, v := range battle.mapTeams {
		if v.IsAlive() {
			v.CountSpeed()

			teams = append(teams, v)
		}
	}

	// 这里排序要从小到大，这样索引就是从慢到快了，索引可以直接作为teamspeedval
	sort.Slice(teams, func(i, j int) bool {
		// 如果速度一样，先手队快，其实就是比较teamindex，小的快
		// 这里其实是比谁慢
		if teams[i].Speed == teams[j].Speed {
			return teams[i].TeamIndex > teams[j].TeamIndex
		}

		return teams[i].Speed < teams[j].Speed
	})

	for i, v := range teams {
		v.SetTeamSpeedVal(i)
	}

	lst := NewHeroList()

	for _, v := range battle.mapHeros {
		if v.IsAlive() {
			lst.AddHero(v)
		}
	}

	lst.SortInBattle()

	return lst
}

func (battle *Battle) StartBattle() {
	root := battle.Log.StartBattle()

	battle.battleReady(root)

	for i := 0; i < MaxTurn; i++ {
		battle.startTurn(root, i)
	}
}

func (battle *Battle) battleReady(parent *BattleLogNode) {
	lst := battle.GenCurHeroList()

	lst.ForEach(func(h *Hero) {
		battle.Log.HeroComeIn(parent, h)
	})
}

func (battle *Battle) startTurn(parent *BattleLogNode, turnindex int) {
	battle.Log.StartTurn(parent, turnindex+1)
}

func NewBattle(w, h int) *Battle {
	scene := NewScene(w, h)
	battle := &Battle{
		Scene:     scene,
		mapTeams:  make(map[int]*Team),
		Log:       NewBattleLog(),
		curHeroID: 1,
		mapHeros:  make(map[int]*Hero),
	}

	return battle
}
