package battle3

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type PropFuncState int

const (
	PropFuncStateOn  = 1
	PropFuncStateOff = 2
)

type PropFunc func(target *Unit, addon *AddOn, state PropFuncState, targetProp []int, param []int, strParam []string) error

type PropFuncMgr struct {
	MapFunc      map[string]PropFunc
	MapBasicFunc map[int]BasicPropFunc
}

func newPropFuncMgr() *PropFuncMgr {
	mgr := &PropFuncMgr{
		MapFunc:      make(map[string]PropFunc),
		MapBasicFunc: make(map[int]BasicPropFunc),
	}

	return mgr
}

func (mgr *PropFuncMgr) Reg(name string, funcItem PropFunc) {
	mgr.MapFunc[name] = funcItem
}

func (mgr *PropFuncMgr) RegBasic(prop int, funcBasic BasicPropFunc) {
	mgr.MapBasicFunc[prop] = funcBasic
}

func (mgr *PropFuncMgr) Run(name string, target *Unit, addon *AddOn, state PropFuncState,
	targetProp []int, param []int, strParam []string) error {

	f, isok := mgr.MapFunc[name]
	if !isok {
		return ErrInvalidPropFunc
	}

	return f(target, addon, state, targetProp, param, strParam)
}

func (mgr *PropFuncMgr) ChgProp(unit *Unit, prop int, val int) (int, error) {
	f, isok := mgr.MapBasicFunc[prop]
	if !isok {
		return 0, ErrInvalidBasicPropFunc
	}

	sv := unit.Props[prop]
	f(unit, prop, val)

	return unit.Props[prop] - sv, nil
}

// addper
// targetProp[i] += param[i] * Prop[strParam[i]] / 100
func propAddPer(target *Unit, addon *AddOn, state PropFuncState, targetProp []int, param []int, strParam []string) error {
	if len(targetProp) != len(param) || len(targetProp) != len(strParam) {
		goutils.Error("propAddPer",
			goutils.JSON("targetProp", targetProp),
			goutils.JSON("param", param),
			goutils.JSON("strParam", strParam),
			zap.Error(ErrInvalidFuncPropParam))

		return ErrInvalidFuncPropParam
	}

	if state == PropFuncStateOn {
		arr := []int{}

		for i, prop := range targetProp {
			prop2, err := Str2Prop(strParam[i])
			if err != nil {
				goutils.Error("propAddPer:Str2Prop",
					goutils.JSON("strParam", strParam),
					zap.Error(err))

				return err
			}

			cv, _ := target.ChgProp(prop, param[i]*target.Props[prop2]/100)
			arr = append(arr, cv)
		}

		if addon != nil {
			addon.CacheVals = arr
		}
	} else if state == PropFuncStateOff && addon != nil && len(addon.CacheVals) == len(targetProp) {
		for i, prop := range targetProp {
			target.ChgProp(prop, -addon.CacheVals[i])
		}
	}

	return nil
}

// add
// targetProp[i] += param[i]
func propAdd(target *Unit, addon *AddOn, state PropFuncState, targetProp []int, param []int, strParam []string) error {
	if len(targetProp) != len(param) {
		goutils.Error("propAdd",
			goutils.JSON("targetProp", targetProp),
			goutils.JSON("param", param),
			goutils.JSON("strParam", strParam),
			zap.Error(ErrInvalidFuncPropParam))

		return ErrInvalidFuncPropParam
	}

	if state == PropFuncStateOn {
		arr := []int{}

		for i, prop := range targetProp {
			cv, _ := target.ChgProp(prop, param[i])
			arr = append(arr, cv)
		}

		if addon != nil {
			addon.CacheVals = arr
		}
	} else if state == PropFuncStateOff && addon != nil && len(addon.CacheVals) == len(targetProp) {
		for i, prop := range targetProp {
			target.ChgProp(prop, -addon.CacheVals[i])
		}
	}

	return nil
}

// rampage
// 狂暴，降低一个属性的50%，并把该数值的80%加到另外一个属性上
func propRampage(target *Unit, addon *AddOn, state PropFuncState, targetProp []int, param []int, strParam []string) error {
	if len(targetProp) != 2 || len(param) != 2 {
		goutils.Error("propAdd",
			goutils.JSON("targetProp", targetProp),
			goutils.JSON("param", param),
			zap.Error(ErrInvalidFuncPropParam))

		return ErrInvalidFuncPropParam
	}

	if state == PropFuncStateOn {
		v0 := target.Props[targetProp[0]] * param[0] / 100
		cv1, _ := target.ChgProp(targetProp[0], -v0)
		cv2, _ := target.ChgProp(targetProp[1], v0*param[1]/100)

		if addon != nil {
			addon.CacheVals = []int{cv1, cv2}
		}
	} else if state == PropFuncStateOff && addon != nil && len(addon.CacheVals) == 2 {
		target.ChgProp(targetProp[0], -addon.CacheVals[0])
		target.ChgProp(targetProp[1], -addon.CacheVals[1])
	}

	return nil
}
