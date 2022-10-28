package battle5funcs

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

const (
	findnearNum  = 0 // 目标数量，0表示全部
	findnearType = 1 // 目标类型，TargetType
)

// 物理攻击
func findNearRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	hero := params.Src

	ti := hero.GetEnemyTeamIndex()

	if ti >= 0 {
		lst := hero.FindNear(hero.Battle.GetTeam(ti).Heros, fd.Vals[0])

		params.Results = lst
	}

	return true, nil
}

// 物理攻击
func findNearInit(fd *battle5.FuncData) error {
	fd.Vals = nil

	if len(fd.InVals) >= 1 {
		if fd.InVals[findnearNum] < 0 {
			goutils.Error("findNearInit",
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
			goutils.Error("findNearInit",
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

func regFindNear(mgr *battle5.FuncMgr) {
	mgr.RegFunc("findnear", battle5.FuncLib{
		OnProc: findNearRun,
		OnInit: findNearInit,
	})
}
