package battle5funcs

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

const (
	basicatkAtkPer       = 0 // 伤害率
	basicatkIgnoreDodge  = 1 // 无视闪避
	basicatkIgnoreCrit   = 2 // 无视暴击
	basicatkIgnoreDef    = 3 // 无视防御
	basicatkIgnoreDamage = 4 // 无视伤害附加
)

// 物理攻击
func basicAtkRun(fd *battle5.FuncData, params *battle5.LibFuncParams) (bool, error) {
	hero := params.Src

	params.Target.ForEach(func(v *battle5.Hero) {
		atk := hero.Props[battle5.PropTypeCurAtk] * hero.Props[battle5.PropTypeCurAtk] / (hero.Props[battle5.PropTypeCurAtk] + v.Props[battle5.PropTypeCurDef])
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

// 物理攻击
func basicAtkInit(fd *battle5.FuncData) error {
	fd.Vals = nil

	if len(fd.InVals) >= 1 {
		if fd.InVals[basicatkAtkPer] < 0 {
			goutils.Error("basicAtkInit",
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
			goutils.Error("basicAtkInit",
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
			goutils.Error("basicAtkInit",
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
			goutils.Error("basicAtkInit",
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
			goutils.Error("basicAtkInit",
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

func regBasicAtk(mgr *battle5.FuncMgr) {
	mgr.RegFunc("basicatk", battle5.FuncLib{
		OnProc: basicAtkRun,
		OnInit: basicAtkInit,
	})
}
