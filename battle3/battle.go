package battle3

func getOther(index int) int {
	if index == 0 {
		return 1
	}

	return 0
}

func getWinner(units []*Unit) int {
	for i, v := range units {
		if !v.IsAlive() {
			return getOther(i)
		}
	}

	return -1
}

func StartBattle(src []*Unit) *BattleResult {
	lst := []*Unit{src[0].ResetAndClone(), src[1].ResetAndClone()}

	return startBattle(lst)
}

func startBattle(src []*Unit) *BattleResult {
	ret := &BattleResult{
		Units:           src,
		ForceFirstIndex: -1, // 默认没有强制先手
	}

	firstIndex := 0 // 默认发起方先手

	// 如果发起者携带先手，则发起者先手
	if src[0].Props[PropTypeIsFirst] == 1 {
		ret.ForceFirstIndex = 0
		firstIndex = 0
	} else if src[1].Props[PropTypeIsFirst] == 1 { // 如果发起者不携带先手，且受击者携带先手，则受击者先手
		ret.ForceFirstIndex = 1
		firstIndex = 1
	}

	ret.FirstIndex = firstIndex

	otherIndex := getOther(firstIndex)

	curTurns := 0
	for {
		curTurns++

		// 先手攻击
		iskilled := src[firstIndex].Attack(src[otherIndex], true)
		if iskilled {
			winner := getWinner(src)
			if winner >= 0 {
				ret.WinIndex = winner

				break
			}
		}

		// 后手攻击
		iskilled = src[otherIndex].Attack(src[firstIndex], true)
		if iskilled {
			winner := getWinner(src)
			if winner >= 0 {
				ret.WinIndex = winner

				break
			}
		}

		// 超过最大回合数，先手负
		if curTurns >= MaxTurns {
			ret.WinIndex = otherIndex

			break
		}
	}

	ret.Turns = curTurns

	return ret
}
