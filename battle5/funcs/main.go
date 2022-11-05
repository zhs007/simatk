package battle5funcs

import "github.com/zhs007/simatk/battle5"

func InitAllFuncs(mgr *battle5.FuncMgr) error {
	regBasicAtk(mgr)
	regMBasicAtk(mgr)

	regFindNear(mgr)
	regFindFar(mgr)
	regFindArea(mgr)
	regFindMe(mgr)
	regFindAll(mgr)
	regFindBuffFrom(mgr)

	regAttachBuffs(mgr)

	regChgProps(mgr)

	return nil
}
