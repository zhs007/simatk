package battle5

// // 找最近的
// func findNear(fd *FuncData, params *LibFuncParams) (bool, error) {
// 	// params.Src.targetSkills = nil

// 	if len(fd.Vals) != 1 {
// 		goutils.Error("findNear",
// 			goutils.JSON("fd", fd),
// 			zap.Error(ErrInvalidFuncValsLength))

// 		return false, ErrInvalidFuncValsLength
// 	}

// 	hero := params.Src

// 	ti := hero.GetEnemyTeamIndex()

// 	if ti >= 0 {
// 		lst := hero.FindNear(hero.battle.mapTeams[ti].Heros, fd.Vals[0])

// 		params.Results = lst
// 		// hero.targetSkills = lst
// 	}

// 	return true, nil
// }

// // 找最远的
// func findFar(fd *FuncData, params *LibFuncParams) (bool, error) {
// 	// params.Src.targetSkills = nil

// 	if len(fd.Vals) != 1 {
// 		goutils.Error("findFar",
// 			goutils.JSON("fd", fd),
// 			zap.Error(ErrInvalidFuncValsLength))

// 		return false, ErrInvalidFuncValsLength
// 	}

// 	hero := params.Src

// 	ti := hero.GetEnemyTeamIndex()

// 	if ti >= 0 {
// 		lst := hero.FindFar(hero.battle.mapTeams[ti].Heros, fd.Vals[0])

// 		params.Results = lst
// 		// hero.targetSkills = lst
// 	}

// 	return true, nil
// }
