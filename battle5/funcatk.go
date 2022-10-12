package battle5

// 普攻，物理攻击
func basicAtk(fd *FuncData, params *LibFuncParams) (bool, error) {
	hero := params.Src

	params.Target.ForEach(func(v *Hero) {
		atk := hero.Props[PropTypeCurAtk] * hero.Props[PropTypeCurAtk] / (hero.Props[PropTypeCurAtk] + v.Props[PropTypeCurDef])
		if atk > 0 {
			shp := v.Props[PropTypeCurHP]

			v.Props[PropTypeCurHP] -= atk

			ehp := v.Props[PropTypeCurHP]

			// if v.Props[PropTypeCurHP] <= 0 {
			// 	params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
			// }

			params.Battle.Log.SkillAttack(params.LogParent, params.Src, v, params.Skill, shp, ehp)

			v.onPropChg(PropTypeCurHP, shp, ehp, NewBattleActionFromData(params.LogParent, params.Src, params.Skill))
		}
	})

	return true, nil
}

// 普攻，魔法攻击
func basicMAtk(fd *FuncData, params *LibFuncParams) (bool, error) {
	hero := params.Src

	params.Target.ForEach(func(v *Hero) {
		atk := hero.Props[PropTypeCurMagic] * hero.Props[PropTypeCurMagic] / (hero.Props[PropTypeCurMagic] + v.Props[PropTypeCurMagic])
		if atk > 0 {
			shp := v.Props[PropTypeCurHP]

			v.Props[PropTypeCurHP] -= atk

			ehp := v.Props[PropTypeCurHP]

			// if v.Props[PropTypeCurHP] <= 0 {
			// 	params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
			// }

			params.Battle.Log.SkillAttack(params.LogParent, params.Src, v, params.Skill, shp, ehp)

			v.onPropChg(PropTypeCurHP, shp, ehp, NewBattleActionFromData(params.LogParent, params.Src, params.Skill))
		}
	})

	return true, nil
}
