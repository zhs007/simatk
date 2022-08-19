package battle3

import (
	"path"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type StaticMgr struct {
	MgrCharacter *CharacterDataMgr
	MgrItem      *ItemDataMgr
	MgrPropFunc  *PropFuncMgr
	MgrEventFunc *EventFuncMgr
	MgrStageDev  *StageDevDataMgr
	MgrStage     *StageDataMgr
	ParamsGenMap *GenMapParams
	CfgPath      string
}

func LoadAllStatic(fnpath string) (*StaticMgr, error) {

	mgrCharacter, err := LoadCharacter(path.Join(fnpath, "character.xlsx"))
	if err != nil {
		goutils.Error("LoadAllStatic:LoadCharacter",
			zap.Error(err))

		return nil, err
	}

	mgrItem, err := LoadItem(path.Join(fnpath, "item.xlsx"))
	if err != nil {
		goutils.Error("LoadAllStatic:LoadItem",
			zap.Error(err))

		return nil, err
	}

	mgrStageDev, err := LoadStageDevData(path.Join(fnpath, "stagedev.xlsx"))
	if err != nil {
		goutils.Error("LoadAllStatic:LoadStageDevData",
			zap.Error(err))

		return nil, err
	}

	mgrStage, err := LoadStageData(path.Join(fnpath, "stage.xlsx"))
	if err != nil {
		goutils.Error("LoadAllStatic:LoadStageData",
			zap.Error(err))

		return nil, err
	}

	paramsGenMap, err := LoadGenMapParams(path.Join(fnpath, "genmapparams.yaml"))
	if err != nil {
		goutils.Error("LoadAllStatic:LoadGenMapParams",
			zap.Error(err))

		return nil, err
	}

	mgr := &StaticMgr{
		MgrCharacter: mgrCharacter,
		MgrItem:      mgrItem,
		MgrPropFunc:  newPropFuncMgr(),
		MgrEventFunc: newEventFuncMgr(),
		MgrStageDev:  mgrStageDev,
		MgrStage:     mgrStage,
		CfgPath:      fnpath,
		ParamsGenMap: paramsGenMap,
	}

	mgr.MgrPropFunc.RegBasic(PropTypeHP, funcHP)
	mgr.MgrPropFunc.RegBasic(PropTypeDPS, funcDPS)

	mgr.MgrPropFunc.RegBasic(PropTypeMaxHP, funcMaxHP)
	mgr.MgrPropFunc.RegBasic(PropTypeCurHP, funcCurHP)
	mgr.MgrPropFunc.RegBasic(PropTypeCurDPS, funcNormal)

	mgr.MgrPropFunc.RegBasic(PropTypeIsFirst, funcState)
	mgr.MgrPropFunc.RegBasic(PropTypeIsDouble, funcState)
	mgr.MgrPropFunc.RegBasic(PropTypeIsFightBack, funcState)
	mgr.MgrPropFunc.RegBasic(PropTypeIsLeech, funcState)
	mgr.MgrPropFunc.RegBasic(PropTypeIsReturnDamage, funcState)

	mgr.MgrPropFunc.RegBasic(PropTypeReturnDamageVal, funcNormal)
	mgr.MgrPropFunc.RegBasic(PropTypeLeechVal, funcNormal)
	mgr.MgrPropFunc.RegBasic(PropTypeUpAtk, funcNormal)
	mgr.MgrPropFunc.RegBasic(PropTypeDownAtk, funcNormal)
	mgr.MgrPropFunc.RegBasic(PropTypeUpDamage, funcNormal)
	mgr.MgrPropFunc.RegBasic(PropTypeDownDamage, funcNormal)

	mgr.MgrPropFunc.Reg("addper", propAddPer)
	mgr.MgrPropFunc.Reg("add", propAdd)
	mgr.MgrPropFunc.Reg("rampage", propRampage)

	mgr.MgrEventFunc.Reg("needids", eventNeedIDs)
	mgr.MgrEventFunc.Reg("canwin", eventCanWin)
	mgr.MgrEventFunc.Reg("check2prop", eventCheck2Prop)
	mgr.MgrEventFunc.Reg("empty", eventEmpty)

	return mgr, nil
}

var MgrStatic *StaticMgr

func InitSystem(fnpath string) error {
	mgr, err := LoadAllStatic(fnpath)
	if err != nil {
		goutils.Error("InitSystem:LoadAllStatic",
			zap.Error(err))

		return err
	}

	MgrStatic = mgr

	return nil
}
