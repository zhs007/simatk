package battle5funcs

import "github.com/zhs007/simatk/battle5"

// 普攻，物理攻击
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

// 普攻，物理攻击
func basicAtkInit(fd *battle5.FuncData) error {
	return nil
	// hero := params.Src

	// params.Target.ForEach(func(v *Hero) {
	// 	atk := hero.Props[PropTypeCurAtk] * hero.Props[PropTypeCurAtk] / (hero.Props[PropTypeCurAtk] + v.Props[PropTypeCurDef])
	// 	if atk > 0 {
	// 		shp := v.Props[PropTypeCurHP]

	// 		v.Props[PropTypeCurHP] -= atk

	// 		ehp := v.Props[PropTypeCurHP]

	// 		// if v.Props[PropTypeCurHP] <= 0 {
	// 		// 	params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
	// 		// }

	// 		params.Battle.Log.SkillAttack(params.LogParent, params.Src, v, params.Skill, shp, ehp)

	// 		v.onPropChg(PropTypeCurHP, shp, ehp, NewBattleActionFromData(params.LogParent, params.Src, params.Skill))
	// 	}
	// })

	// return true, nil
}

func RegBasicAtk(mgr *battle5.FuncMgr) {
	mgr.RegFunc("basicatk", battle5.FuncLib{
		OnProc: basicAtkRun,
		OnInit: basicAtkInit,
	})
}

// // 普攻，魔法攻击
// func basicMAtk(fd *FuncData, params *LibFuncParams) (bool, error) {
// 	hero := params.Src

// 	params.Target.ForEach(func(v *Hero) {
// 		atk := hero.Props[PropTypeCurMagic] * hero.Props[PropTypeCurMagic] / (hero.Props[PropTypeCurMagic] + v.Props[PropTypeCurMagic])
// 		if atk > 0 {
// 			shp := v.Props[PropTypeCurHP]

// 			v.Props[PropTypeCurHP] -= atk

// 			ehp := v.Props[PropTypeCurHP]

// 			// if v.Props[PropTypeCurHP] <= 0 {
// 			// 	params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
// 			// }

// 			params.Battle.Log.SkillAttack(params.LogParent, params.Src, v, params.Skill, shp, ehp)

// 			v.onPropChg(PropTypeCurHP, shp, ehp, NewBattleActionFromData(params.LogParent, params.Src, params.Skill))
// 		}
// 	})

// 	return true, nil
// }
