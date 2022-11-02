package battle5

import (
	"path"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type StaticMgr struct {
	MgrHeroData  *HeroDataMgr
	MgrSkillData *SkillDataMgr
	MgrBuffData  *BuffDataMgr
	MgrFunc      *FuncMgr
}

var MgrStatic *StaticMgr

func NewStaticMgr(dir string, oninitallfuncs FuncInitAllFuncs) (*StaticMgr, error) {
	mgr := &StaticMgr{
		MgrFunc: NewFuncMgr(),
	}

	MgrStatic = mgr

	// mgr.MgrFunc.Init()
	oninitallfuncs(mgr.MgrFunc)

	mgrherodata, err := LoadHeroData(path.Join(dir, "heros.xlsx"))
	if err != nil {
		goutils.Error("NewStaticMgr:LoadHeroData",
			zap.String("dir", dir),
			zap.Error(err))

		return nil, err
	}

	mgr.MgrHeroData = mgrherodata

	mgrskilldata, err := LoadSkillData(path.Join(dir, "skills.xlsx"))
	if err != nil {
		goutils.Error("NewStaticMgr:LoadSkillData",
			zap.String("dir", dir),
			zap.Error(err))

		return nil, err
	}

	mgr.MgrSkillData = mgrskilldata

	mgrbuffdata, err := LoadBuffData(path.Join(dir, "buffs.xlsx"))
	if err != nil {
		goutils.Error("NewStaticMgr:LoadBuffData",
			zap.String("dir", dir),
			zap.Error(err))

		return nil, err
	}

	mgr.MgrBuffData = mgrbuffdata

	return mgr, nil
}
