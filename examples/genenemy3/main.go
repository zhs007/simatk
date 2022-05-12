package main

import "github.com/zhs007/simatk/battle3"

func main() {
	ret := battle3.GenEnemy(battle3.NewUnit(120, 80),
		&battle3.GenEnemyParam{
			UnitType:        battle3.UnitTypeNormal,
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

	ret.Output("genenemy3.json")
}
