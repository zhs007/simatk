package battle5funcs

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

const (
	findallExcludeMe  = 0 // 排除自己
	findallTargetType = 1 // 目标类型
)

// 找全部
func findAllRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	hero := params.Src

	if battle5.TargetType(fd.Vals[findallTargetType]) == battle5.TargetTypeEnemy {
		ti := hero.GetEnemyTeamIndex()

		params.Results = hero.Battle.GetTeam(ti).Heros.Clone()
	} else if battle5.TargetType(fd.Vals[findallTargetType]) == battle5.TargetTypeFriend {
		ti := hero.TeamIndex

		if fd.Vals[findallExcludeMe] == 0 {
			params.Results = hero.Battle.GetTeam(ti).Heros.Clone()
		} else {
			params.Results = hero.Battle.GetTeam(ti).Heros.CloneEx(func(h *battle5.Hero) bool {
				return !h.IsMe(hero)
			})
		}
	} else {
		if fd.Vals[findallExcludeMe] == 0 {
			params.Results = hero.Battle.GenAliveHeroList(nil)
		} else {
			params.Results = hero.Battle.GenAliveHeroList(func(h *battle5.Hero) bool {
				return !h.IsMe(hero)
			})
		}
	}

	return true, nil
}

// 找全部
func findAllInit(fd *battle5.FuncData) error {
	fd.Vals = nil

	if len(fd.InVals) >= 1 {
		if !isValidBool(fd.InVals[0]) {
			goutils.Error("findAllInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[0])
	} else {
		fd.Vals = append(fd.Vals, 1)
	}

	if len(fd.InStrVals) >= 1 {
		if !isValidTargetType(fd.StrVals[0]) {
			goutils.Error("findAllInit",
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

func regFindAll(mgr *battle5.FuncMgr) {
	mgr.RegFunc("findall", battle5.FuncLib{
		OnProc: findAllRun,
		OnInit: findAllInit,
	})
}
