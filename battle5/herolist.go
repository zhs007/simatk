package battle5

type HeroList struct {
	Heros []*Hero
}

func (hl *HeroList) Init(lst []*HeroData) {
	hl.Heros = []*Hero{}

	for _, v := range lst {
		h := NewHeroEx(v)

		hl.Heros = append(hl.Heros, h)
	}
}

func (hl *HeroList) AddHero(h *Hero) {
	hl.Heros = append(hl.Heros, h)
}

func NewHeroList() *HeroList {
	heros := &HeroList{}

	return heros
}
