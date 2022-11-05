package battle5funcs

import (
	"github.com/zhs007/simatk/battle5"
)

// 附加buff
func attachBuffsRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	// hero := params.Src

	// ti := hero.GetEnemyTeamIndex()

	// if ti >= 0 {
	// 	lst := hero.FindNear(hero.Battle.GetTeam(ti).Heros, fd.Vals[0])

	// 	params.Results = lst
	// }

	return true, nil
}

// 附加buff
func attachBuffsInit(fd *battle5.FuncData) error {
	// fd.Vals = nil

	// if len(fd.InVals) >= 1 {
	// 	if fd.InVals[findnearNum] < 0 {
	// 		goutils.Error("findNearInit",
	// 			goutils.JSON("params", fd),
	// 			zap.Error(ErrInvalidVals))

	// 		return ErrInvalidVals
	// 	}

	// 	fd.Vals = append(fd.Vals, fd.InVals[findnearNum])
	// } else {
	// 	fd.Vals = append(fd.Vals, 1)
	// }

	// if len(fd.InStrVals) >= 1 {
	// 	if !isValidTargetType(fd.InStrVals[0]) {
	// 		goutils.Error("findNearInit",
	// 			goutils.JSON("params", fd),
	// 			zap.Error(ErrInvalidVals))

	// 		return ErrInvalidVals
	// 	}

	// 	fd.Vals = append(fd.Vals, int(toTargetType(fd.InStrVals[0])))
	// } else {
	// 	fd.Vals = append(fd.Vals, int(battle5.TargetTypeEnemy))
	// }

	return nil
}

func regAttachBuffs(mgr *battle5.FuncMgr) {
	mgr.RegFunc("attachbuffs", battle5.FuncLib{
		OnProc: attachBuffsRun,
		OnInit: attachBuffsInit,
	})
}
