package battle3

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type EventFunc func(id int, params []int, strParams []string, unit *Unit, root *Event, cur *EventPool) bool

type EventFuncMgr struct {
	mapFunc map[string]EventFunc
}

func (mgr *EventFuncMgr) Reg(name string, funcEvent EventFunc) {
	mgr.mapFunc[name] = funcEvent
}

func (mgr *EventFuncMgr) Run(name string,
	id int, params []int, strParams []string, unit *Unit, root *Event, cur *EventPool) bool {
	f, isok := mgr.mapFunc[name]
	if isok {
		return f(id, params, strParams, unit, root, cur)
	}

	goutils.Error("EventFuncMgr.Run",
		zap.String("name", name),
		zap.Error(ErrInvalidEventFunc))

	return false
}

func newEventFuncMgr() *EventFuncMgr {
	return &EventFuncMgr{
		mapFunc: make(map[string]EventFunc),
	}
}

// needids -
// 需要哪些事件，各至少需要多少个
func eventNeedIDs(id int, params []int, strParams []string, unit *Unit, root *Event, cur *EventPool) bool {
	// params 一定是偶数长度
	if len(params)%2 != 0 {
		goutils.Error("eventNeedIDs",
			goutils.JSON("params", params),
			zap.Error(ErrInvalidFuncEventParam))

		return false
	}

	for i := 0; i < len(params)/2; i++ {
		n := root.CountID(params[i*2])
		if n < params[i*2+1] {
			return false
		}
	}

	return true
}

// canwin -
// 能战胜，默认的怪物应该都需要这个事件
func eventCanWin(id int, params []int, strParams []string, unit *Unit, root *Event, cur *EventPool) bool {
	monster, err := MgrStatic.MgrCharacter.NewUnit(id)
	if err != nil {
		goutils.Error("eventCanWin",
			goutils.JSON("id", id),
			zap.Error(ErrInvalidCharacterID))

		return false
	}

	ret := startBattle([]*Unit{unit.Clone(), monster})
	return ret.WinIndex == 0
}

// check2prop -
// 检查2个prop之间的值和一个数比较
// strParam[0]是第一个prop，strParam[1]是操作符，strParam[2]是第二个prop，strParam[3]是比较符
// param[0]是数值
// 如果操作符是 除，则结果要x100
func eventCheck2Prop(id int, params []int, strParams []string, unit *Unit, root *Event, cur *EventPool) bool {
	if len(params) < 1 || len(strParams) < 4 {
		goutils.Error("eventCheckPropOff",
			goutils.JSON("params", params),
			goutils.JSON("strParams", strParams),
			zap.Error(ErrInvalidFuncEventParam))

		return false
	}

	propid1, err := Str2Prop(strParams[0])
	if err != nil {
		goutils.Error("eventCheckPropOff:Str2Prop",
			goutils.JSON("strParams", strParams),
			zap.Error(ErrInvalidFuncEventParam))

		return false
	}

	propid2, err := Str2Prop(strParams[2])
	if err != nil {
		goutils.Error("eventCheckPropOff:Str2Prop",
			goutils.JSON("strParams", strParams),
			zap.Error(ErrInvalidFuncEventParam))

		return false
	}

	val := 0

	if strParams[1] == "+" {
		val = unit.Props[propid1] + unit.Props[propid2]
	} else if strParams[1] == "-" {
		val = unit.Props[propid1] - unit.Props[propid2]
	} else if strParams[1] == "*" {
		val = unit.Props[propid1] * unit.Props[propid2]
	} else if strParams[1] == "/" {
		val = unit.Props[propid1] * 100 / unit.Props[propid2]
	}

	if strParams[3] == "<" {
		return val < params[0]
	} else if strParams[3] == ">" {
		return val > params[0]
	} else if strParams[3] == "==" {
		return val == params[0]
	} else if strParams[3] == "<=" {
		return val <= params[0]
	} else if strParams[3] == ">=" {
		return val >= params[0]
	}

	// 如果在前面没有返回，一定是参数错误
	goutils.Error("eventCheckPropOff:non-return",
		goutils.JSON("params", params),
		goutils.JSON("strParams", strParams),
		zap.Error(ErrInvalidFuncEventParam))

	return false
}

// empty -
// 一个事件都不符合条件
func eventEmpty(id int, params []int, strParams []string, unit *Unit, root *Event, cur *EventPool) bool {
	return len(cur.Events) == 0
}
