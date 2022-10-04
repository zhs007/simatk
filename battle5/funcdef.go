package battle5

type FuncEachHero func(*Hero)

type FuncOnText func(string)

type FuncIsLess func(i, j int) bool

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

type FuncLib func(*FuncData, *LibFuncParams) (bool, error)

type FuncData struct {
	FuncName string   `json:"name,omitempty"`
	Vals     []int    `json:"vals,omitempty"`
	StrVals  []string `json:"strvals,omitempty"`
}
