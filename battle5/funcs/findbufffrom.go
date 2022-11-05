package battle5funcs

import (
	"github.com/zhs007/simatk/battle5"
)

// 找buff源
func findBuffFromRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	if params.TriggerData.SrcBuff != nil && params.TriggerData.SrcBuff.From != nil {
		params.Results = battle5.NewHeroListEx2(params.TriggerData.SrcBuff.From)

		// 这个接口一般用于 TriggerTypeFind 里，所以如果找到目标，需要返回 true 中断 find 进程
		return true, nil
	}

	return false, nil
}

// 找buff源
func findBuffFromInit(fd *battle5.FuncData) error {
	return nil
}

func regFindBuffFrom(mgr *battle5.FuncMgr) {
	mgr.RegFunc("findbufffrom", battle5.FuncLib{
		OnProc: findBuffFromRun,
		OnInit: findBuffFromInit,
	})
}
