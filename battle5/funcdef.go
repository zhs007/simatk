package battle5

type FuncEachHero func(*Hero)
type FuncEachHeroBreak func(*Hero) bool

type FuncOnText func(string)

type FuncIsLess func(i, j int) bool

type LibFuncParams struct {
	Src    *Hero
	Target *HeroList
	Battle *Battle
}

func NewLibFuncParams(battle *Battle, src *Hero, target *HeroList) *LibFuncParams {
	return &LibFuncParams{
		Battle: battle,
		Target: target,
		Src:    src,
	}
}

type FuncLib func(*FuncData, *LibFuncParams) (bool, error)

type FuncData struct {
	FuncName string   `json:"name,omitempty"`
	Vals     []int    `json:"vals,omitempty"`
	StrVals  []string `json:"strvals,omitempty"`
}
