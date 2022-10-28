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

// func (mgr *FuncMgr) Init() {
// 	// battle5func.RegBasicAtk(mgr)

// 	// mgr.RegFunc("basicatk", FuncLib{

// 	// })
// 	// mgr.RegFunc("basicmatk", basicMAtk)

// 	// mgr.RegFunc("findnear", findNear)
// 	// mgr.RegFunc("findfar", findFar)
// }

func (mgr *FuncMgr) Run(fd *FuncData, params *LibFuncParams) (bool, error) {
	f, isok := mgr.MapFunc[fd.FuncName]
	if isok {
		return f.OnProc(fd, params)
	}

	goutils.Error("FuncMgr.Run",
		goutils.JSON("duncdata", fd),
		zap.Error(ErrInvalidFuncName))

	return false, ErrInvalidFuncName
}

func (mgr *FuncMgr) InitFuncData(fd *FuncData) error {
	f, isok := mgr.MapFunc[fd.FuncName]
	if isok {
		return f.OnInit(fd)
	}

	goutils.Error("FuncMgr.InitFuncData",
		goutils.JSON("duncdata", fd),
		zap.Error(ErrInvalidFuncName))

	return ErrInvalidFuncName
}

func NewFuncMgr() *FuncMgr {
	return &FuncMgr{
		MapFunc: make(map[string]FuncLib),
	}
}
