package battle5funcs

import "github.com/zhs007/simatk/battle5"

func InitAllFuncs(mgr *battle5.FuncMgr) error {
	RegBasicAtk(mgr)

	return nil
}
