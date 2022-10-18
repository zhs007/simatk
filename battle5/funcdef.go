package battle5

// 简单遍历hero list
type FuncEachHero func(*Hero)

// 遍历hero list，返回false会break
type FuncEachHeroBreak func(*Hero) bool

// 遍历hero skills，返回false会break
type FuncEachHeroSkill func(*Hero, *Skill) bool

type FuncOnText func(string)

type FuncIsLess func(i, j int) bool

type LibFuncParams struct {
	Src       *Hero
	Target    *HeroList
	Battle    *Battle
	LogParent *BattleLogNode
	Skill     *Skill
	Results   *HeroList
}

func NewLibFuncParams(battle *Battle, src *Hero, target *HeroList, skill *Skill, parent *BattleLogNode) *LibFuncParams {
	return &LibFuncParams{
		Battle:    battle,
		Target:    target,
		Src:       src,
		LogParent: parent,
		Skill:     skill,
	}
}

type FuncInitAllFuncs func(*FuncMgr) error

type FuncLibProc func(*FuncData, *LibFuncParams) (bool, error)
type FuncLibInit func(*FuncData) error

type FuncLib struct {
	OnProc FuncLibProc
	OnInit FuncLibInit
}

type FuncData struct {
	FuncName  string   `json:"name,omitempty"`
	InVals    []int    `json:"vals,omitempty"`
	InStrVals []string `json:"strvals,omitempty"`
	Vals      []int    `json:"-"`
	StrVals   []string `json:"-"`
}
