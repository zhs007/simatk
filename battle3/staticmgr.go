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

	mgr := &StaticMgr{
		MgrCharacter: mgrCharacter,
		MgrItem:      mgrItem,
		MgrPropFunc:  newPropFuncMgr(),
	}

	mgr.MgrPropFunc.RegBasic(PropTypeHP, funcNormal)
	mgr.MgrPropFunc.RegBasic(PropTypeDPS, funcNormal)

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

	return mgr, nil
}
