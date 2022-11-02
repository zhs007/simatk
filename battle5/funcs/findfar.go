package battle5funcs

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

// 找最远的
func findFarRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	hero := params.Src

	ti := hero.GetEnemyTeamIndex()

	if ti >= 0 {
		lst := hero.FindFar(hero.Battle.GetTeam(ti).Heros, fd.Vals[0])

		params.Results = lst
	}

	return true, nil
}

// 找最远的
func findFarInit(fd *battle5.FuncData) error {
	fd.Vals = nil

	if len(fd.InVals) >= 1 {
		if fd.InVals[findnearNum] < 0 {
			goutils.Error("findFarInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[findnearNum])
	} else {
		fd.Vals = append(fd.Vals, 1)
	}

	if len(fd.InStrVals) >= 1 {
		if !isValidTargetType(fd.StrVals[0]) {
			goutils.Error("findFarInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, int(toTargetType(fd.InStrVals[0])))
	} else {
		fd.Vals = append(fd.Vals, int(battle5.TargetTypeEnemy))
	}

	return nil
}

func regFindFar(mgr *battle5.FuncMgr) {
	mgr.RegFunc("findfar", battle5.FuncLib{
		OnProc: findFarRun,
		OnInit: findFarInit,
	})
}
