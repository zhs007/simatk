package battle5funcs

import (
	"github.com/zhs007/simatk/battle5"
)

// 找自己
func findMeRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	params.Results = battle5.NewHeroListEx2(params.Src)

	return true, nil
}

// 找自己
func findMeInit(fd *battle5.FuncData) error {
	return nil
}

func regFindMe(mgr *battle5.FuncMgr) {
	mgr.RegFunc("findme", battle5.FuncLib{
		OnProc: findMeRun,
		OnInit: findMeInit,
	})
}
