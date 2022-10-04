package battle5

import (
	"path"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type StaticMgr struct {
	MgrHeroData *HeroDataMgr
	MgrFunc     *FuncMgr
}

var MgrStatic *StaticMgr

func NewStaticMgr(dir string) (*StaticMgr, error) {
	mgr := &StaticMgr{
		MgrFunc: NewFuncMgr(),
	}

	mgr.MgrFunc.Init()

	mgrherodata, err := LoadHeroData(path.Join(dir, "heros.xlsx"))
	if err != nil {
		goutils.Error("NewStaticMgr:LoadHeroData",
			zap.String("dir", dir),
			zap.Error(err))

		return nil, err
	}

	mgr.MgrHeroData = mgrherodata

	return mgr, nil
}
