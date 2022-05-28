package battle3

import "math/rand"

// 跑流程，返回 -1 表示失败，0 表示路径打通，
func calcWithAI1(event *Event, unit *Unit) bool {
	for {
		lst := event.BuildNextEvents()
		if len(lst) == 0 {
			return true
		}

		cr := rand.Int() % len(lst)
		unit.ProcEvent(lst[cr].ID)
		if unit.Props[PropTypeCurHP] <= 0 {
			return false
		}

		lst[cr].isFinished = true
		if lst[cr].IsEnding {
			return true
		}
	}
}

// 返回AI1玩num次，成功的次数
func CalcWinTimesWithAI1(event *Event, num int, unit *Unit) int {
	winnum := 0

	for i := 0; i < num; i++ {
		if calcWithAI1(event.Clone(), unit.Clone()) {
			winnum++
		}
	}

	return winnum
}
