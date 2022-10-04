package battle5

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// 找最近的
func findNear(fd *FuncData, params *LibFuncParams) (bool, error) {
	params.Src.targetSkills = nil

	if len(fd.Vals) != 1 {
		goutils.Error("findNear",
			goutils.JSON("fd", fd),
			zap.Error(ErrInvalidFuncValsLength))

		return false, ErrInvalidFuncValsLength
	}

	hero := params.Src

	ti := hero.GetEnemyTeamIndex()

	if ti >= 0 {
		lst := hero.FindNear(hero.battle.mapTeams[ti].Heros, fd.Vals[0])

		hero.targetSkills = lst
	}

	return true, nil
}
