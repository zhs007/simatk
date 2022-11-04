package battle5

import (
	"sort"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type Battle struct {
	Scene             *Scene
	mapTeams          map[int]*Team
	Log               *BattleLog
	mapHeros          map[HeroInstanceID]*Hero
	curHeroInstanceID HeroInstanceID
	curBuffInstanceID BuffInstanceID
	CurTurn           int
	isEnd             bool
}

func (battle *Battle) GenAliveHeroList(onhero FuncEachHeroBool) *HeroList {
	hl := NewHeroList()

	if onhero == nil {
		for _, v := range battle.mapHeros {
			if v.IsAlive() {
				hl.AddHero(v)
			}
		}

		return hl
	}

	for _, v := range battle.mapHeros {
		if v.IsAlive() && onhero(v) {
			hl.AddHero(v)
		}
	}

	return hl
}

func (battle *Battle) GenHeroInstanceID() HeroInstanceID {
	id := battle.curHeroInstanceID

	battle.curHeroInstanceID++

	return id
}

func (battle *Battle) GenBuffInstanceID() BuffInstanceID {
	id := battle.curBuffInstanceID

	battle.curBuffInstanceID++

	return id
}

func (battle *Battle) GetTeam(index int) *Team {
	return battle.mapTeams[index]
}

func (battle *Battle) SetTeam(index int, lst []*HeroData, autoSetPos bool) {
	battle.mapTeams[index] = NewTeam(battle, index, lst)

	if autoSetPos {
		battle.mapTeams[index].AutoSetPos()
	}

	for _, v := range battle.mapTeams[index].Heros.Heros {
		v.InstanceID = battle.GenHeroInstanceID()

		battle.Scene.AddHero(v)

		battle.mapHeros[v.InstanceID] = v
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

	return lst.Format()
}

func (battle *Battle) StartBattle() {
	root := battle.Log.StartBattle()

	battle.battleReady(root)

	for i := 0; i < MaxTurn; i++ {
		battle.startTurn(root, i)

		if battle.isEnd {
			break
		}
	}

	battle.Log.BattleEnd(root)
}

func (battle *Battle) battleReady(parent *BattleLogNode) {
	ready := battle.Log.BattleReady(parent)

	lst := battle.GenCurHeroList()

	lst.ForEach(func(h *Hero) {
		battle.Log.HeroComeIn(ready, h)
	})
}

func (battle *Battle) startTurn(parent *BattleLogNode, turnindex int) {
	battle.CurTurn = turnindex

	turn := battle.Log.StartTurn(parent, turnindex+1)

	lst := battle.GenCurHeroList()

	// 找目标
	lst.ForEach(func(h *Hero) {
		target := h.FindTarget()
		if target == nil || target.IsEmpty() {
			h.LastTarget = nil

			battle.Log.FindTarget(turn, h, nil)
		} else {
			h.LastTarget = target.Heros[0]

			battle.Log.FindTarget(turn, h, target.Heros[0])
		}
	})

	// 移动 v2
	// 先判断角色是否需要移动，如果必须要移动，则至少会移动一格
	// 这里主要为了后行动的角色也移动
	lst1 := NewHeroList()
	lst.ForEach(func(h *Hero) {
		if h.LastTarget != nil {
			if h.onMoveStepStart() {
				lst1.AddHero(h)
			}
		}
	})

	if !lst1.IsEmpty() {
		for {
			lst2 := NewHeroList()

			lst1.ForEach(func(h *Hero) {
				// if h.LastTarget != nil {
				if h.move2TargetStep(h.LastTarget) {
					if h.CanAttackWithDistance(h.LastTarget) {
						h.onMoveStepEnd(turn)
					} else {
						lst2.AddHero(h)
					}
				} else {
					h.onMoveStepEnd(turn)
				}
				// }
			})

			if lst2.GetNum() <= 0 {
				break
			}

			lst1 = lst2
		}
	}

	// // 移动
	// lst.ForEach(func(h *Hero) {
	// 	if h.targetMove != nil && !h.targetMove.IsEmpty() {
	// 		if h.CanMove() {
	// 			p := h.Move2Target(h.targetMove.Heros[0])
	// 			if p != nil {
	// 				battle.Log.HeroMove(turn, h, p)
	// 				h.Pos.Set(p)
	// 			}
	// 		}
	// 	}
	// })

	// 攻击
	lst.ForEachWithBreak(func(h *Hero) bool {
		if h.IsAlive() {
			h.ForEachSkills(func(ch *Hero, s *Skill) bool {
				ch.UseSkill(turn, s)

				return battle.isEnd
			})

			if battle.isEnd {
				return false
			}
		}

		return true
	})

	battle.Log.EndTurn(parent, turnindex+1)
}

func (battle *Battle) onHeroBeSkilled(h *Hero, fd *BattleActionFromData) {
	battle.mapTeams[h.TeamIndex].onHeroBeSkilled(h, fd)

	if !h.IsAlive() {
		battle.Log.KillHero(fd.Parent, fd.Hero, h, fd.Skill)

		battle.checkGameEnd()
	}
}

func (battle *Battle) checkGameEnd() {
	aliveteams := 0

	for _, v := range battle.mapTeams {
		if v.IsAlive() {
			aliveteams++
		}
	}

	if aliveteams <= 1 {
		battle.isEnd = true
	}
}

func (battle *Battle) NewBuff(buffid BuffID, from *Hero, skill *Skill) (*Buff, error) {
	bd := MgrStatic.MgrBuffData.GetBuffData(buffid)
	if bd == nil {
		goutils.Error("Battle.NewBuff",
			zap.Int("buffid", int(buffid)),
			zap.Error(ErrInvalidBuffID))

		return nil, ErrInvalidBuffID
	}

	buff := &Buff{
		Data:       bd,
		InstanceID: battle.GenBuffInstanceID(),
		From:       from,
		FromSkill:  skill,
	}

	return buff, nil
}

func NewBattle(w, h int) *Battle {
	scene := NewScene(w, h)
	battle := &Battle{
		Scene:             scene,
		mapTeams:          make(map[int]*Team),
		Log:               NewBattleLog(),
		curHeroInstanceID: 1,
		curBuffInstanceID: 1,
		CurTurn:           1,
		mapHeros:          make(map[HeroInstanceID]*Hero),
	}

	return battle
}

func NewBattleEx(mgr *StaticMgr, team0 []HeroID, team1 []HeroID, w, h int) *Battle {
	battle := NewBattle(w, h)

	lst0 := []*HeroData{}
	for _, v := range team0 {
		hd := mgr.MgrHeroData.GetHeroData(v)
		if hd != nil {
			lst0 = append(lst0, hd)
		}
	}

	battle.SetTeam(0, lst0, true)

	lst1 := []*HeroData{}
	for _, v := range team1 {
		hd := mgr.MgrHeroData.GetHeroData(v)
		if hd != nil {
			lst1 = append(lst1, hd)
		}
	}

	battle.SetTeam(1, lst1, true)

	return battle
}
