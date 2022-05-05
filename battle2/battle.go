package battle2

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
		Units:           lst,
		ForceFirstIndex: -1,
	}

	if src[0].Props[PropTypeCurSpeed] == src[1].Props[PropTypeCurSpeed] {
		ret.ForceFirstIndex = firstIndex
	} else if src[0].Props[PropTypeCurSpeed] > src[1].Props[PropTypeCurSpeed] {
		firstIndex = 0
	} else if src[0].Props[PropTypeCurSpeed] < src[1].Props[PropTypeCurSpeed] {
		firstIndex = 1
	}

	ret.FirstIndex = firstIndex

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
