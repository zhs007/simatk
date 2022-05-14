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
			UnitType:        battle3.UnitTypeMoreHP,
			MinTurns:        1,
			MaxTurns:        10,
			MinLastHP:       1,
			MaxLastHP:       99,
			StartTotalVal:   100,
			EndTotalVal:     500,
			IsWinner:        true,
			DetailTurnOff:   1,
			DetailLastHPOff: 10,
		})

	ret.OutputExcel(f)

	ret.Output("genenemy3.json")

	f.DeleteSheet("Sheet1")

	f.SaveAs("genenemy3.xlsx")
}
