package battle5

// 普攻，物理攻击
func basicAtk(params *LibFuncParams) (bool, error) {
	hero := params.Src
	for _, v := range params.Target {
		atk := hero.Props[PropTypeCurAtk] * hero.Props[PropTypeCurAtk] / (hero.Props[PropTypeCurAtk] + v.Props[PropTypeCurDef])
		if atk > 0 {
			v.Props[PropTypeCurHP] -= atk

			if v.Props[PropTypeCurHP] <= 0 {
				params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
			}
		}
	}

	return true, nil
}

// 普攻，魔法攻击
func basicMAtk(params *LibFuncParams) (bool, error) {
	hero := params.Src
	for _, v := range params.Target {
		atk := hero.Props[PropTypeCurMagic] * hero.Props[PropTypeCurMagic] / (hero.Props[PropTypeCurMagic] + v.Props[PropTypeCurMagic])
		if atk > 0 {
			v.Props[PropTypeCurHP] -= atk

			if v.Props[PropTypeCurHP] <= 0 {
				params.Battle.mapTeams[v.TeamIndex].needUpdAlive = true
			}
		}
	}

	return true, nil
}
