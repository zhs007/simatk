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
}

func (hero *Hero) IsAlive() bool {
	return hero.Props[PropTypeCurHP] > 0
}

func (hero *Hero) UseSkill(skill *Skill) {
	if skill != nil {
		if skill.Data.Atk != nil {
			MgrStatic.MgrFunc.Run(skill.Data.Atk.FuncName, NewLibFuncParams(hero.battle, hero, nil))
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
