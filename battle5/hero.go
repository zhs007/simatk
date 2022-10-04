package battle5

type Hero struct {
	ID               HeroID
	Props            map[PropType]int
	SX, SY           int       // 初始坐标，按本地坐标来的，也就是2队人，这个坐标都是对自己在左边的
	X, Y             int       // 坐标
	TeamIndex        int       // 队伍索引，0-进攻方，1-防守方
	RealBattleHeroID int       // 战斗里hero的唯一标识
	Data             *HeroData // 直接读表数据
	Skills           []*Skill  // 技能
	battle           *Battle
	targetSkills     *HeroList // 技能目标
	tmpDistance      int       // 临时距离，按距离排序用
}

func (hero *Hero) IsAlive() bool {
	return hero.Props[PropTypeCurHP] > 0
}

func (hero *Hero) UseSkill(skill *Skill) {
	if skill != nil {
		// 找目标
		if skill.Data.Find != nil {
			MgrStatic.MgrFunc.Run(skill.Data.Find, NewLibFuncParams(hero.battle, hero, nil))
		}

		// 伤害
		if skill.Data.Atk != nil {
			MgrStatic.MgrFunc.Run(skill.Data.Atk, NewLibFuncParams(hero.battle, hero, nil))
		}
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
	if h0.X == h1.X {
		return 0
	}

	if hero.TeamIndex == 0 {
		if h0.X < h1.X {
			return 1
		}

		return -1
	}

	if h0.X < h1.X {
		return -1
	}

	return 1
}

// 在队列里找最近的多少个，一定会返回一个新的herolist
func (hero *Hero) FindNear(lst *HeroList, num int) *HeroList {
	if num >= lst.GetNum() {
		return lst.Clone()
	}

	lst.ForEach(func(h *Hero) {
		if hero.RealBattleHeroID == h.RealBattleHeroID {
			h.tmpDistance = 0
		} else {
			ox := h.X - hero.X
			oy := h.Y - hero.Y

			if ox < 0 {
				ox = -ox
			}

			if oy < 0 {
				oy = -oy
			}

			if ox >= oy {
				h.tmpDistance = ox
			} else {
				h.tmpDistance = oy
			}
		}
	})

	// 由进及远
	lst.Sort(func(i, j int) bool {
		if lst.Heros[i].tmpDistance == lst.Heros[j].tmpDistance {
			// 优先距离自己这边近的
			nearmyside := hero.cmpNearMySide(lst.Heros[i], lst.Heros[j])
			if nearmyside == 0 {
				// 优先y轴小的
				return lst.Heros[i].Y <= lst.Heros[j].Y
			}

			return nearmyside == 1
		}

		return lst.Heros[i].tmpDistance < lst.Heros[j].tmpDistance
	})

	return NewHeroListEx(lst.Heros[0:num])
}

func (hero *Hero) Clone() *Hero {
	return NewHero(hero.Props[PropTypeHP],
		hero.Props[PropTypeAtk],
		hero.Props[PropTypeDef],
		hero.Props[PropTypeMagic],
		hero.Props[PropTypeSpeed],
		hero.Props[PropTypeAttackType] == 1)
}

func NewHero(hp int, atk int, def int, magic int, speed int, isMagicAtk bool) *Hero {
	hero := &Hero{
		Props: make(map[PropType]int),
		SX:    -1,
		SY:    -1,
		X:     -1,
		Y:     -1,
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
		SX:     -1,
		SY:     -1,
		X:      -1,
		Y:      -1,
		Data:   hd,
		battle: battle,
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
