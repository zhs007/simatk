package battle3

// 基本属性修改，val为改变量
// 注意，如果是状态数据，val为最终值
type BasicPropFunc func(unit *Unit, prop int, val int)

// 普通属性，一定大于0
func funcNormal(unit *Unit, prop int, val int) {
	if val > 0 {
		unit.Props[prop] += val
	}

	if unit.Props[prop]+val < 0 {
		unit.Props[prop] = 0
	}
}

// 状态属性，只有 0 和 1
func funcState(unit *Unit, prop int, val int) {
	if val != 0 {
		unit.Props[prop] = 1
	} else {
		unit.Props[prop] = 0
	}
}

// CurHP，需要 >= 0 && <= MaxHP
func funcCurHP(unit *Unit, prop int, val int) {
	if val > 0 {
		if unit.Props[PropTypeCurHP]+val > unit.Props[PropTypeMaxHP] {
			unit.Props[PropTypeCurHP] = unit.Props[PropTypeMaxHP]

			return
		}
	}

	if val < 0 {
		if unit.Props[PropTypeCurHP]+val < 0 {
			unit.Props[PropTypeCurHP] = 0

			return
		}
	}

	unit.Props[PropTypeCurHP] += val
}

// MaxHP，需要 >= 1，且减小时，还要考虑 CurHP 的被动溢出
func funcMaxHP(unit *Unit, prop int, val int) {
	if val < 0 {
		// 最大HP不能为0
		if unit.Props[PropTypeMaxHP]+val < 1 {
			unit.Props[PropTypeMaxHP] = 1
		} else {
			unit.Props[PropTypeMaxHP] += val
		}

		if unit.Props[PropTypeCurHP] > unit.Props[PropTypeMaxHP] {
			unit.Props[PropTypeCurHP] = unit.Props[PropTypeMaxHP]
		}

		return
	}

	unit.Props[PropTypeMaxHP] += val
}
