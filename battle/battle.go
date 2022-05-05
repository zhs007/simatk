package battle

func getOther(index int) int {
	if index == 0 {
		return 1
	}

	return 0
}

func StartBattle(src []*Unit, firstIndex int) *BattleResult {
	lst := []*Unit{src[0].ResetAndClone(), src[1].ResetAndClone()}
	// units[0].Reset()
	// units[1].Reset()

	ret := &BattleResult{
		Units:      lst,
		FirstIndex: firstIndex,
	}

	otherIndex := getOther(firstIndex)

	for {
		lst[firstIndex].Attack(lst[otherIndex])
		if !lst[otherIndex].IsAlive() {
			ret.WinIndex = firstIndex

			break
		}

		lst[otherIndex].Attack(lst[firstIndex])
		if !lst[firstIndex].IsAlive() {
			ret.WinIndex = otherIndex

			break
		}
	}

	return ret
}
