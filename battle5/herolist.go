package battle5

import (
	"sort"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
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

func (hl *HeroList) ForEachWithBreak(oneach FuncEachHeroBreak) bool {
	for _, v := range hl.Heros {
		if !oneach(v) {
			return false
		}
	}

	return true
}

func (hl *HeroList) Init(battle *Battle, lst []*HeroData) {
	hl.Heros = []*Hero{}

	for _, v := range lst {
		h := NewHeroEx(battle, v)

		hl.Heros = append(hl.Heros, h)
	}
}

func (hl *HeroList) RemoveHero(h *Hero) {
	for i, v := range hl.Heros {
		if v.RealBattleHeroID == h.RealBattleHeroID {
			hl.Heros = append(hl.Heros[0:i], hl.Heros[i+1:]...)

			return
		}
	}
}

func (hl *HeroList) Find(h *Hero) int {
	for i, v := range hl.Heros {
		if v.RealBattleHeroID == h.RealBattleHeroID {
			return i
		}
	}

	return -1
}

func (hl *HeroList) AddHero(h *Hero) error {
	if h == nil {
		goutils.Error("HeroList.AddHero",
			zap.Error(ErrHeroIsNull))

		return ErrHeroIsNull
	}

	if hl.Find(h) >= 0 {
		goutils.Error("HeroList.AddHero",
			zap.Error(ErrDuplicateHero))

		return ErrDuplicateHero
	}

	hl.Heros = append(hl.Heros, h)

	return nil
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

func (hl *HeroList) GetAliveHeros() *HeroList {
	lst := NewHeroList()

	hl.ForEach(func(h *Hero) {
		if h.IsAlive() {
			lst.AddHero(h)
		}
	})

	return lst.Format()
}

func (hl *HeroList) Format() *HeroList {
	if hl == nil {
		return nil
	}

	if hl.IsEmpty() {
		return nil
	}

	return hl
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

func NewHeroListEx2(hs ...*Hero) *HeroList {
	heros := &HeroList{}

	heros.Heros = append(heros.Heros, hs...)

	return heros
}
