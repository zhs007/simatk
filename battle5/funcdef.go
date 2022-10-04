package battle5

type FuncEachHero func(*Hero)

type FuncOnText func(string)

type LibFuncParams struct {
	Src    *Hero
	Target []*Hero
	Battle *Battle
}

func NewLibFuncParams(battle *Battle, src *Hero, target []*Hero) *LibFuncParams {
	return &LibFuncParams{
		Battle: battle,
		Src:    src,
		Target: target,
	}
}

type FuncLib func(*LibFuncParams) (bool, error)
