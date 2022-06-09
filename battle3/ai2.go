package battle3

import "math/rand"

// 随机跑流程，一定打不过的怪就不要打了
func calcWithAI2(event *Event, unit *Unit) bool {
	for {
		lst := event.BuildNextEventsEx(func(e *Event) bool {
			nu := unit.Clone()
			nu.ProcEvent(e.ID)
			return nu.Props[PropTypeCurHP] > 0
		})
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

// 返回AI2玩num次，成功的次数
func CalcWinTimesWithAI2(event *Event, num int, unit *Unit) int {
	winnum := 0

	for i := 0; i < num; i++ {
		if calcWithAI2(event.Clone(), unit.Clone()) {
			winnum++
		}
	}

	return winnum
}
