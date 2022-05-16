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

func (mgr *PropFuncMgr) ChgProp(unit *Unit, prop int, val int) error {
	f, isok := mgr.MapBasicFunc[prop]
	if !isok {
		return ErrInvalidBasicPropFunc
	}

	f(unit, prop, val)

	return nil
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

	for i, prop := range targetProp {
		prop2, err := Str2Prop(strParam[i])
		if err != nil {
			goutils.Error("propAddPer:Str2Prop",
				goutils.JSON("strParam", strParam),
				zap.Error(err))

			return err
		}

		target.ChgProp(prop, param[i]*target.Props[prop2]/100)
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

	for i, prop := range targetProp {
		target.ChgProp(prop, param[i])
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

	v0 := target.Props[targetProp[0]] * param[0] / 100
	target.ChgProp(targetProp[0], -v0)
	target.ChgProp(targetProp[1], v0*param[1]/100)

	return nil
}
