package battle5funcs

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

// 魔法攻击
func basicMAtkRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	hero := params.Src

	params.Target.ForEach(func(v *battle5.Hero) {
		atk := hero.Props[battle5.PropTypeCurMagic] * hero.Props[battle5.PropTypeCurMagic] / (hero.Props[battle5.PropTypeCurMagic] + v.Props[battle5.PropTypeCurMagic])
		if atk > 0 {
			shp := v.Props[battle5.PropTypeCurHP]

			v.Props[battle5.PropTypeCurHP] -= atk

			ehp := v.Props[battle5.PropTypeCurHP]

			// if v.Props[PropTypeCurHP] <= 0 {
			// 	params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
			// }

			params.Battle.Log.SkillAttack(params.LogParent, params.Src, v, params.Skill, shp, ehp)

			v.OnPropChg(battle5.PropTypeCurHP, shp, ehp, battle5.NewBattleActionFromData(params.LogParent, params.Src, params.Skill))
		}
	})

	return true, nil
}

// 魔法攻击
func basicMAtkInit(fd *battle5.FuncData) error {
	fd.Vals = nil

	if len(fd.InVals) >= 1 {
		if fd.InVals[basicatkAtkPer] < 0 {
			goutils.Error("basicMAtkInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[basicatkAtkPer])
	} else {
		fd.Vals = append(fd.Vals, 100)
	}

	if len(fd.InVals) >= 2 {
		if !isValidBool(fd.InVals[basicatkIgnoreDodge]) {
			goutils.Error("basicMAtkInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[basicatkIgnoreDodge])
	} else {
		fd.Vals = append(fd.Vals, 1)
	}

	if len(fd.InVals) >= 3 {
		if !isValidBool(fd.InVals[basicatkIgnoreCrit]) {
			goutils.Error("basicMAtkInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[basicatkIgnoreCrit])
	} else {
		fd.Vals = append(fd.Vals, 1)
	}

	if len(fd.InVals) >= 4 {
		if !isValidBool(fd.InVals[basicatkIgnoreDef]) {
			goutils.Error("basicMAtkInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[basicatkIgnoreDef])
	} else {
		fd.Vals = append(fd.Vals, 1)
	}

	if len(fd.InVals) >= 5 {
		if !isValidBool(fd.InVals[basicatkIgnoreDamage]) {
			goutils.Error("basicMAtkInit",
				goutils.JSON("params", fd),
				zap.Error(ErrInvalidVals))

			return ErrInvalidVals
		}

		fd.Vals = append(fd.Vals, fd.InVals[basicatkIgnoreDamage])
	} else {
		fd.Vals = append(fd.Vals, 1)
	}

	return nil
}

func regMBasicAtk(mgr *battle5.FuncMgr) {
	mgr.RegFunc("basicmatk", battle5.FuncLib{
		OnProc: basicMAtkRun,
		OnInit: basicMAtkInit,
	})
}
