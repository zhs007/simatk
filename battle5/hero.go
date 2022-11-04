package battle5

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type Hero struct {
	Battle         *Battle
	ID             HeroID
	Props          map[PropType]int
	StaticPos      *Pos                 // 初始坐标，按本地坐标来的，也就是2队人，这个坐标都是对自己在左边的
	Pos            *Pos                 // 坐标
	TeamIndex      int                  // 队伍索引，0-进攻方，1-防守方
	InstanceID     HeroInstanceID       // 战斗里hero的唯一标识
	Data           *HeroData            // 直接读表数据
	Skills         []*Skill             // 技能
	tmpDistance    int                  // 临时距离，按距离排序用
	movePos        *Pos                 // 移动位置，修正位移时用
	lastMoveVal    int                  // 剩余的移动距离，修正位移时用
	LastTarget     *Hero                // 上一次的目标
	MapSingleState map[BuffEffect]*Buff // 简单状态
	Trigger        *BuffTriggerMap
	// targetMove       *HeroList // 移动目标
	// targetSkills     *HeroList // 技能目标
}

func (hero *Hero) IsMe(h *Hero) bool {
	return hero.InstanceID == h.InstanceID
}

func (hero *Hero) IsAlive() bool {
	return hero.Props[PropTypeCurHP] > 0
}

func (hero *Hero) UseSkill(parent *BattleLogNode, skill *Skill) {
	if skill != nil {
		logSkill := hero.Battle.Log.UseSkill(parent, hero, skill)

		if skill.canUseSkill() {
			targets := skill.findTarget(hero)
			if targets != nil && targets.GetNum() > 0 {
				targets.ForEach(func(th *Hero) {
					hero.Battle.Log.FindSkillTarget(logSkill, hero, th)
				})

				targets.ForEach(func(th *Hero) {
					skill.attack(logSkill, hero, th)
				})
			} else {
				hero.Battle.Log.FindSkillTarget(logSkill, hero, nil)
			}
		}

		// // 找目标
		// if skill.Data.Find != nil {
		// 	MgrStatic.MgrFunc.Run(skill.Data.Find, NewLibFuncParams(hero.battle, hero, nil))
		// }

		// // 伤害
		// if skill.Data.Atk != nil {
		// 	MgrStatic.MgrFunc.Run(skill.Data.Atk, NewLibFuncParams(hero.battle, hero, nil))
		// }
	}
}

func (hero *Hero) Attack(toHero *Hero) bool {
	if hero.Props[PropTypeAttackType] == 0 {
		atk := hero.Props[PropTypeCurAtk] * hero.Props[PropTypeCurAtk] / (hero.Props[PropTypeCurAtk] + toHero.Props[PropTypeCurDef])
		if atk > 0 {
			toHero.Props[PropTypeCurHP] -= atk
		}
	} else {
		atk := hero.Props[PropTypeCurMagic] * hero.Props[PropTypeCurMagic] / (hero.Props[PropTypeCurMagic] + toHero.Props[PropTypeCurMagic])
		if atk > 0 {
			toHero.Props[PropTypeCurHP] -= atk
		}
	}

	return toHero.Props[PropTypeCurHP] <= 0
}

// 获取攻击队伍，可能会被混乱等状态影响，如果返回-1，表示所有队伍都需要选择
func (hero *Hero) GetEnemyTeamIndex() int {
	if hero.TeamIndex == 0 {
		return 1
	}

	return 0
}

// 判断谁距离本边更近，1表示h0更近，0表示一样近，-1表示h1更近
func (hero *Hero) cmpNearMySide(h0 *Hero, h1 *Hero) int {
	if h0.Pos.X == h1.Pos.X {
		return 0
	}

	if hero.TeamIndex == 0 {
		if h0.Pos.X < h1.Pos.X {
			return 1
		}

		return -1
	}

	if h0.Pos.X < h1.Pos.X {
		return -1
	}

	return 1
}

// 判断谁距离敌方更近，1表示h0更近，0表示一样近，-1表示h1更近
func (hero *Hero) cmpNearEnemySide(h0 *Hero, h1 *Hero) int {
	if h0.Pos.X == h1.Pos.X {
		return 0
	}

	if hero.TeamIndex == 1 {
		if h0.Pos.X < h1.Pos.X {
			return 1
		}

		return -1
	}

	if h0.Pos.X < h1.Pos.X {
		return -1
	}

	return 1
}

// 在队列里找最近的多少个，一定会返回一个新的herolist
func (hero *Hero) FindNear(lst0 *HeroList, num int) *HeroList {
	if num <= 0 {
		return nil
	}

	lst1 := lst0.GetAliveHeros()
	if lst1 == nil {
		return nil
	}

	if num >= lst1.GetNum() {
		return lst1
	}

	lst1.ForEach(func(h *Hero) {
		if hero.InstanceID == h.InstanceID {
			h.tmpDistance = 0
		} else {
			h.tmpDistance = hero.Pos.CalcDistance(h.Pos)
		}
	})

	// 由近及远
	lst1.Sort(func(i, j int) bool {
		if lst1.Heros[i].tmpDistance == lst1.Heros[j].tmpDistance {
			// 优先上一次目标
			if hero.LastTarget != nil && hero.LastTarget.InstanceID == lst1.Heros[i].InstanceID {
				return true
			}

			if hero.LastTarget != nil && hero.LastTarget.InstanceID == lst1.Heros[j].InstanceID {
				return false
			}

			// 优先距离自己这边近的
			nearmyside := hero.cmpNearMySide(lst1.Heros[i], lst1.Heros[j])
			if nearmyside == 0 {
				// 优先y轴小的
				return lst1.Heros[i].Pos.Y <= lst1.Heros[j].Pos.Y
			}

			return nearmyside == 1
		}

		return lst1.Heros[i].tmpDistance < lst1.Heros[j].tmpDistance
	})

	return NewHeroListEx(lst1.Heros[0:num])
}

// 在队列里找最远的多少个，一定会返回一个新的herolist
func (hero *Hero) FindFar(lst0 *HeroList, num int) *HeroList {
	if num <= 0 {
		return nil
	}

	lst1 := lst0.GetAliveHeros()

	if num >= lst1.GetNum() {
		return lst1 //.Clone()
	}

	lst1.ForEach(func(h *Hero) {
		if hero.IsMe(h) {
			h.tmpDistance = 0
		} else {
			h.tmpDistance = hero.Pos.CalcDistance(h.Pos)
		}
	})

	// 由远及近
	lst1.Sort(func(i, j int) bool {
		if lst1.Heros[i].tmpDistance == lst1.Heros[j].tmpDistance {
			// 优先上一次目标
			if hero.LastTarget != nil && hero.LastTarget.IsMe(lst1.Heros[i]) {
				return true
			}

			if hero.LastTarget != nil && hero.LastTarget.IsMe(lst1.Heros[j]) {
				return false
			}

			// 优先距离敌方近的
			nearmyside := hero.cmpNearEnemySide(lst1.Heros[i], lst1.Heros[j])
			if nearmyside == 0 {
				// 优先y轴小的
				return lst1.Heros[i].Pos.Y <= lst1.Heros[j].Pos.Y
			}

			return nearmyside == 1
		}

		return lst1.Heros[i].tmpDistance > lst1.Heros[j].tmpDistance
	})

	return NewHeroListEx(lst1.Heros[0:num])
}

// 选择目标
func (hero *Hero) FindTarget() *HeroList {
	params := NewLibFuncParams(hero.Battle, hero, nil, nil, nil,
		NewTriggerDataFind(hero.Battle.CurTurn, hero))

	ret, err := hero.Trigger.OnTrigger(TriggerTypeFind, params)
	if err != nil {
		goutils.Error("Hero.FindTarget",
			zap.Error(err))

		return nil
	}

	if !ret {
		return params.Results.Format()
	}

	return hero.findTargetWithFuncData(hero.Data.Find)

	// return lst.Format()
	// hero.targetMove = lst
	// hero.targetSkills = nil

	// return hero.targetMove
}

// 选择目标
func (hero *Hero) findTargetWithFuncData(fd *FuncData) *HeroList {
	// 如果自己的格子上有人，在格子内找目标
	lstp := hero.Battle.Scene.GetHerosWithPos(hero.Pos)
	if lstp.GetNum() > 1 {
		var last *Hero
		var first *Hero

		lstp.ForEachWithBreak(func(h *Hero) bool {
			if !h.IsAlive() {
				return true
			}

			if !h.IsMe(hero) {
				first = h

				if hero.LastTarget != nil {
					if hero.LastTarget.IsMe(h) {
						last = h

						return false
					}
				} else {
					return false
				}
			}

			return true
		})

		if last != nil {
			// hero.targetMove = NewHeroListEx2(last)
			return NewHeroListEx2(last)
		}

		if first != nil {
			return NewHeroListEx2(first)
		}

		return nil
		// hero.targetMove = NewHeroListEx2(first)
		// return NewHeroListEx2(first)

		// return hero.targetMove
	}

	params := NewLibFuncParams(hero.Battle, hero, nil, nil, nil, nil)

	MgrStatic.MgrFunc.Run(fd,
		params)

	return params.Results.Format() //hero.targetSkills
}

func (hero *Hero) CanMove() bool {
	if hero.Props[PropTypeCurMovDistance] <= 0 {
		return false
	}

	if hero.LastTarget != nil {
		return !hero.CanAttackWithDistance(hero.LastTarget)
	}

	return true
}

func (hero *Hero) moveY(target *Hero, mov int) *Pos {
	oy := target.Pos.Y - hero.Pos.Y
	if oy == 0 {
		return nil
	}

	// 直接移到位
	if Abs(oy) <= mov {
		// hero.Pos.Set(target.Pos)

		return target.Pos.Clone()
	}

	if oy < 0 {
		// hero.Pos.SetXY(target.Pos.X, hero.Pos.Y-mov)

		return NewPos(target.Pos.X, hero.Pos.Y-mov)
	}

	// hero.Pos.SetXY(target.Pos.X, hero.Pos.Y+mov)

	return NewPos(target.Pos.X, hero.Pos.Y+mov)
}

// 移动
func (hero *Hero) Move2Target(target *Hero) *Pos {
	// 优先x轴移动
	ox := target.Pos.X - hero.Pos.X
	if ox == 0 {
		return hero.moveY(target, hero.Props[PropTypeCurMovDistance])
	}

	if Abs(ox) <= hero.Props[PropTypeCurMovDistance] {
		o := hero.Props[PropTypeCurMovDistance] - Abs(ox)
		if o > 0 {
			return hero.moveY(target, o)
		}

		// hero.Pos.SetXY(hero.Pos.X+ox, hero.Pos.Y)

		return NewPos(hero.Pos.X+ox, hero.Pos.Y)
	}

	if ox < 0 {
		// hero.Pos.SetXY(hero.Pos.X-hero.Props[PropTypeCurMovDistance], hero.Pos.Y)

		return NewPos(hero.Pos.X-hero.Props[PropTypeCurMovDistance], hero.Pos.Y)
	}

	// hero.Pos.SetXY(hero.Pos.X+hero.Props[PropTypeCurMovDistance], hero.Pos.Y)

	return NewPos(hero.Pos.X+hero.Props[PropTypeCurMovDistance], hero.Pos.Y)
}

func (hero *Hero) GetRealPos() *Pos {
	if hero.movePos != nil {
		return hero.movePos
	}

	return hero.Pos
}

// 是否可以攻击到
func (hero *Hero) CanAttackWithDistance(toHero *Hero) bool {
	tmpDistance := hero.GetRealPos().CalcDistance(toHero.GetRealPos())

	return tmpDistance < hero.Props[PropTypeAtkDistance]
}

// 移动一步
func (hero *Hero) move2TargetStepY(target *Hero) bool {
	oy := target.Pos.Y - hero.movePos.Y
	if oy == 0 {
		return false
	}

	hero.lastMoveVal--

	// 直接移到位
	if Abs(oy) < 1 {
		hero.movePos.Y = target.Pos.Y

		return hero.lastMoveVal > 0
	}

	if oy < 0 {
		hero.movePos.Y--

		return hero.lastMoveVal > 0
	}

	hero.movePos.Y++

	return hero.lastMoveVal > 0
}

// 移动一步
func (hero *Hero) move2TargetStep(target *Hero) bool {
	// if hero.movePos == nil {
	// 	hero.onMoveStepStart()
	// }

	if hero.Props[PropTypeCurMovDistance] == 99 {
		hero.movePos.Set(target.Pos)

		return false
	}

	if hero.lastMoveVal <= 0 {
		return false
	}

	// 优先x轴移动
	ox := target.Pos.X - hero.movePos.X
	if ox == 0 {
		return hero.move2TargetStepY(target)
	}

	hero.lastMoveVal--

	if ox < 0 {
		hero.movePos.X--

		return hero.lastMoveVal > 0
	}

	hero.movePos.X++

	return hero.lastMoveVal > 0
}

func (hero *Hero) onMoveStepStart() bool {
	if hero.CanAttackWithDistance(hero.LastTarget) {
		return false
	}

	hero.movePos = NewPos(hero.Pos.X, hero.Pos.Y)
	hero.lastMoveVal = hero.Props[PropTypeMovDistance]

	return true
}

func (hero *Hero) onMoveStepEnd(parent *BattleLogNode) {
	if !hero.Pos.Equal(hero.movePos) {
		hero.Battle.Log.HeroMove(parent, hero, hero.movePos)

		hero.Battle.Scene.HeroMove(hero, hero.movePos)
		hero.Pos.Set(hero.movePos)
	}

	hero.movePos = nil
}

func (hero *Hero) ForEachSkills(oneach FuncEachHeroSkill) bool {
	for _, v := range hero.Skills {
		if !oneach(hero, v) {
			return false
		}
	}

	return true
}

func (hero *Hero) OnPropChg(pt PropType, startval int, endval int, fd *BattleActionFromData) {
	if pt == PropTypeCurHP {
		if endval <= 0 {
			hero.Battle.onHeroBeSkilled(hero, fd)
		}
	}
}

func (hero *Hero) Clone() *Hero {
	if hero.Data == nil {
		return NewHero(hero.Props[PropTypeHP],
			hero.Props[PropTypeAtk],
			hero.Props[PropTypeDef],
			hero.Props[PropTypeMagic],
			hero.Props[PropTypeSpeed],
			hero.Props[PropTypeAttackType] == 1)
	}

	nh := &Hero{
		Props:   make(map[PropType]int),
		Trigger: NewBuffTriggerMap(),
	}

	nh.ID = hero.ID
	nh.StaticPos = hero.StaticPos.Clone()
	nh.Pos = hero.Pos.Clone()
	nh.TeamIndex = hero.TeamIndex
	nh.InstanceID = hero.InstanceID
	nh.Data = hero.Data
	nh.Battle = hero.Battle
	// nh.targetMove = hero.targetMove
	// nh.targetSkills = hero.targetSkills
	nh.tmpDistance = hero.tmpDistance
	// nh.Trigger = NewBuffTriggerMap()

	for k, v := range hero.Props {
		nh.Props[k] = v
	}

	for _, v := range hero.Skills {
		nh.Skills = append(nh.Skills, v.Clone())
	}

	return nh
}

func NewHero(hp int, atk int, def int, magic int, speed int, isMagicAtk bool) *Hero {
	hero := &Hero{
		Props:          make(map[PropType]int),
		MapSingleState: make(map[BuffEffect]*Buff),
		Trigger:        NewBuffTriggerMap(),
	}

	hero.Props[PropTypeHP] = hp
	hero.Props[PropTypeAtk] = atk
	hero.Props[PropTypeDef] = def
	hero.Props[PropTypeMagic] = magic
	hero.Props[PropTypeSpeed] = speed

	hero.Props[PropTypeMovDistance] = 1
	hero.Props[PropTypeAtkDistance] = 1
	hero.Props[PropTypePlace] = 1

	if isMagicAtk {
		hero.Props[PropTypeAttackType] = 1
	} else {
		hero.Props[PropTypeAttackType] = 0
	}

	hero.Props[PropTypeMaxHP] = hp
	hero.Props[PropTypeCurHP] = hp
	hero.Props[PropTypeCurAtk] = atk
	hero.Props[PropTypeCurDef] = def
	hero.Props[PropTypeCurMagic] = magic
	hero.Props[PropTypeCurSpeed] = speed

	hero.Props[PropTypeCurMovDistance] = 1
	hero.Props[PropTypeCurAtkDistance] = 1

	return hero
}

func NewHeroEx(battle *Battle, hd *HeroData) *Hero {
	hero := &Hero{
		ID:     HeroID(hd.ID),
		Props:  make(map[PropType]int),
		Data:   hd,
		Battle: battle,
		StaticPos: &Pos{
			X: -1,
			Y: -1,
		},
		Pos: &Pos{
			X: -1,
			Y: -1,
		},
		InstanceID:     battle.GenHeroInstanceID(),
		MapSingleState: make(map[BuffEffect]*Buff),
		Trigger:        NewBuffTriggerMap(),
	}

	hero.Props[PropTypeHP] = hd.HP
	hero.Props[PropTypeAtk] = hd.Atk
	hero.Props[PropTypeDef] = hd.Def
	hero.Props[PropTypeMagic] = hd.Magic
	hero.Props[PropTypeSpeed] = hd.Speed

	hero.Props[PropTypeMovDistance] = hd.MovDistance
	hero.Props[PropTypeAtkDistance] = hd.AtkDistance
	hero.Props[PropTypePlace] = hd.Place

	hero.Props[PropTypeMaxHP] = hd.HP
	hero.Props[PropTypeCurHP] = hd.HP
	hero.Props[PropTypeCurAtk] = hd.Atk
	hero.Props[PropTypeCurDef] = hd.Def
	hero.Props[PropTypeCurMagic] = hd.Magic
	hero.Props[PropTypeCurSpeed] = hd.Speed

	hero.Props[PropTypeCurMovDistance] = hd.MovDistance
	hero.Props[PropTypeCurAtkDistance] = hd.AtkDistance

	for _, v := range hd.Skills {
		sd := MgrStatic.MgrSkillData.GetSkillData(v)
		skill := NewSkill(sd)

		hero.Skills = append(hero.Skills, skill)
	}

	return hero
}
