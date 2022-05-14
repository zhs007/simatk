package main

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/simatk/battle3"
)

func main() {
	f := excelize.NewFile()

	ret := battle3.GenEnemy(battle3.NewUnit(120, 80),
		&battle3.GenEnemyParam{
			Title:           "偏肉",
			UnitType:        battle3.UnitTypeHP,
			MinTurns:        2,
			MaxTurns:        2,
			MinLastHP:       81,
			MaxLastHP:       90,
			StartTotalVal:   100,
			EndTotalVal:     500,
			IsWinner:        true,
			DetailTurnOff:   0,
			DetailLastHPOff: 0,
		})

	ret.OutputExcel(f)

	ret = battle3.GenEnemy(battle3.NewUnit(120, 80),
		&battle3.GenEnemyParam{
			Title:           "精英肉",
			UnitType:        battle3.UnitTypeMoreHP,
			MinTurns:        4,
			MaxTurns:        5,
			MinLastHP:       21,
			MaxLastHP:       40,
			StartTotalVal:   100,
			EndTotalVal:     500,
			IsWinner:        true,
			DetailTurnOff:   0,
			DetailLastHPOff: 0,
		})

	ret.OutputExcel(f)

	ret = battle3.GenEnemy(battle3.NewUnit(120, 80),
		&battle3.GenEnemyParam{
			Title:           "精英输出",
			UnitType:        battle3.UnitTypeMoreDPS,
			MinTurns:        1,
			MaxTurns:        1,
			MinLastHP:       21,
			MaxLastHP:       30,
			StartTotalVal:   100,
			EndTotalVal:     500,
			IsWinner:        true,
			DetailTurnOff:   0,
			DetailLastHPOff: 0,
		})

	ret.OutputExcel(f)

	ret = battle3.GenEnemy(battle3.NewUnit(120, 80),
		&battle3.GenEnemyParam{
			Title:           "均衡",
			UnitType:        battle3.UnitTypeNormal,
			MinTurns:        2,
			MaxTurns:        2,
			MinLastHP:       51,
			MaxLastHP:       60,
			StartTotalVal:   100,
			EndTotalVal:     500,
			IsWinner:        true,
			DetailTurnOff:   0,
			DetailLastHPOff: 0,
		})

	ret.OutputExcel(f)

	// ret.Output("genenemy3.json")

	f.DeleteSheet("Sheet1")

	f.SaveAs("genenemy3.xlsx")
}
