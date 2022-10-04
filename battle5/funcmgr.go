package battle5

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type FuncMgr struct {
	MapFunc map[string]FuncLib
}

func (mgr *FuncMgr) RegFunc(name string, f FuncLib) {
	mgr.MapFunc[name] = f
}

func (mgr *FuncMgr) Init() {
	mgr.RegFunc("basicatk", basicAtk)
	mgr.RegFunc("basicmatk", basicMAtk)
}

func (mgr *FuncMgr) Run(name string, params *LibFuncParams) (bool, error) {
	f, isok := mgr.MapFunc[name]
	if isok {
		return f(params)
	}

	goutils.Error("FuncMgr.Run",
		zap.String("name", name),
		zap.Error(ErrInvalidFuncName))

	return false, ErrInvalidFuncName
}

func NewFuncMgr() *FuncMgr {
	return &FuncMgr{
		MapFunc: make(map[string]FuncLib),
	}
}
