package battle5

import (
	"sort"
)

type HeroList struct {
	Heros []*Hero
}

func (hl *HeroList) GetNum() int {
	return len(hl.Heros)
}

func (hl *HeroList) Clone() *HeroList {
	lst := NewHeroList()

	for _, v := range hl.Heros {
		lst.AddHero(v)
	}

	return lst
}

func (hl *HeroList) ForEach(oneach FuncEachHero) {
	for _, v := range hl.Heros {
		oneach(v)
	}
}

func (hl *HeroList) Init(battle *Battle, lst []*HeroData) {
	hl.Heros = []*Hero{}

	for _, v := range lst {
		h := NewHeroEx(battle, v)

		hl.Heros = append(hl.Heros, h)
	}
}

func (hl *HeroList) AddHero(h *Hero) {
	hl.Heros = append(hl.Heros, h)
}

func (hl *HeroList) Sort(isless FuncIsLess) {
	sort.Slice(hl.Heros, isless)
}

func (hl *HeroList) SortInBattle() {
	// 这里比快
	sort.Slice(hl.Heros, func(i, j int) bool {
		if hl.Heros[i].Props[PropTypeCurSpeed] == hl.Heros[j].Props[PropTypeCurSpeed] {
			// 如果队伍速度也一样，只可能是同一队，那么优先x小的，再优先y小的
			if hl.Heros[i].Props[PropTypeTeamSpeedVal] == hl.Heros[j].Props[PropTypeTeamSpeedVal] {
				if hl.Heros[i].StaticPos.X == hl.Heros[j].StaticPos.X {
					return hl.Heros[i].StaticPos.Y <= hl.Heros[j].StaticPos.Y
				}

				return hl.Heros[i].StaticPos.X < hl.Heros[j].StaticPos.X
			}

			return hl.Heros[i].Props[PropTypeTeamSpeedVal] > hl.Heros[j].Props[PropTypeTeamSpeedVal]
		}

		return hl.Heros[i].Props[PropTypeCurSpeed] > hl.Heros[j].Props[PropTypeCurSpeed]
	})
}

func (hl *HeroList) SortInAutoSetPos() {
	// 这里比快
	sort.Slice(hl.Heros, func(i, j int) bool {
		if hl.Heros[i].Props[PropTypeCurSpeed] == hl.Heros[j].Props[PropTypeCurSpeed] {
			// 如果速度一样，优先站位靠前的
			return hl.Heros[i].Props[PropTypePlace] <= hl.Heros[j].Props[PropTypePlace]
		}

		return hl.Heros[i].Props[PropTypeCurSpeed] > hl.Heros[j].Props[PropTypeCurSpeed]
	})
}

func (hl *HeroList) SortInAutoSetPosSlow() {
	// 这里比慢
	sort.Slice(hl.Heros, func(i, j int) bool {
		if hl.Heros[i].Props[PropTypeCurSpeed] == hl.Heros[j].Props[PropTypeCurSpeed] {
			// 如果速度一样，优先站位靠前的
			return hl.Heros[i].Props[PropTypePlace] >= hl.Heros[j].Props[PropTypePlace]
		}

		return hl.Heros[i].Props[PropTypeCurSpeed] < hl.Heros[j].Props[PropTypeCurSpeed]
	})
}

func (hl *HeroList) IsEmpty() bool {
	return len(hl.Heros) == 0
}

func (hl *HeroList) Size() int {
	return len(hl.Heros)
}

func NewHeroList() *HeroList {
	heros := &HeroList{}

	return heros
}

func NewHeroListEx(lst []*Hero) *HeroList {
	heros := &HeroList{}

	heros.Heros = append(heros.Heros, lst...)

	return heros
}
