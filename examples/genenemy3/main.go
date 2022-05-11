package main

import "github.com/zhs007/simatk/battle3"

func main() {
	ret := battle3.GenEnemy(battle3.NewUnit(120, 80),
		&battle3.GenEnemyParam{
			UnitType:      battle3.UnitTypeNormal,
			MinTurns:      1,
			MaxTurns:      10,
			MinLastHP:     10,
			MaxLastHP:     120,
			StartTotalVal: 100,
			EndTotalVal:   500,
			IsWinner:      true,
		})

	ret.Output("genenemy3.json")
}
