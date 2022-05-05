package battle2

func getOther(index int) int {
	if index == 0 {
		return 1
	}

	return 0
}

func StartBattle(units []*Unit, firstIndex int) *BattleResult {
	units[0].Reset()
	units[1].Reset()

	ret := &BattleResult{
		Units:      units,
		FirstIndex: firstIndex,
	}

	otherIndex := getOther(firstIndex)

	for {
		units[firstIndex].Attack(units[otherIndex])
		if !units[otherIndex].IsAlive() {
			ret.WinIndex = firstIndex

			break
		}

		units[otherIndex].Attack(units[firstIndex])
		if !units[firstIndex].IsAlive() {
			ret.WinIndex = otherIndex

			break
		}
	}

	return ret
}
