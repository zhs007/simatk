package battle5

type Hero struct {
	Props map[PropType]int
}

func (hero *Hero) Attack(toHero *Hero) bool {
	atk := hero.Props[PropTypeCurAtk] * hero.Props[PropTypeCurAtk] / (hero.Props[PropTypeCurAtk] + toHero.Props[PropTypeCurDef])
	toHero.Props[PropTypeCurHP] -= atk

	return toHero.Props[PropTypeCurHP] <= 0
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
