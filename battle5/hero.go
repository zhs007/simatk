package battle5

type Hero struct {
	Props map[PropType]int
}

func (hero *Hero) Attack(toHero *Hero) bool {
	if hero.Props[PropTypeCurAtk] > hero.Props[PropTypeCurMagic] {
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
		hero.Props[PropTypeSpeed])
}

func NewHero(hp int, atk int, def int, magic int, speed int) *Hero {
	hero := &Hero{
		Props: make(map[PropType]int),
	}

	hero.Props[PropTypeHP] = hp
	hero.Props[PropTypeAtk] = atk
	hero.Props[PropTypeDef] = def
	hero.Props[PropTypeMagic] = magic
	hero.Props[PropTypeSpeed] = speed

	hero.Props[PropTypeMaxHP] = hp
	hero.Props[PropTypeCurHP] = hp
	hero.Props[PropTypeCurAtk] = atk
	hero.Props[PropTypeCurDef] = def
	hero.Props[PropTypeCurMagic] = magic
	hero.Props[PropTypeCurSpeed] = speed

	return hero
}
