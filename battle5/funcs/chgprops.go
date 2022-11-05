package battle5funcs

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

// 改变属性
func chgPropsRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	hero := params.Src

	for i := 0; i < len(fd.Vals)/2; i++ {
		params.CacheVals = append(params.CacheVals, hero.Props[battle5.PropType(fd.Vals[i*2])])

		hero.Props[battle5.PropType(fd.Vals[i*2])] += fd.Vals[i*2+1]
	}

	return true, nil
}

// 改变属性
func chgPropsInit(fd *battle5.FuncData) error {
	fd.Vals = nil

	if len(fd.InStrVals) != len(fd.InVals) {
		goutils.Error("chgPropsInit",
			goutils.JSON("funcdata", fd),
			zap.Error(ErrInvalidValsOrStrVals))

		return ErrInvalidValsOrStrVals
	}

	for i, v := range fd.InStrVals {
		cpt := battle5.Str2PropType(v)
		if cpt == battle5.PropTypeNone {
			goutils.Error("chgPropsInit",
				zap.String("str", v),
				zap.Error(ErrInvalidStrVals))

			return ErrInvalidStrVals
		}

		fd.Vals = append(fd.Vals, int(cpt), fd.InVals[i])
	}

	return nil
}

func regChgProps(mgr *battle5.FuncMgr) {
	mgr.RegFunc("chgprops", battle5.FuncLib{
		OnProc: chgPropsRun,
		OnInit: chgPropsInit,
	})
}
