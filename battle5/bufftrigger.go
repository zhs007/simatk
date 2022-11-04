package battle5

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type BuffTriggerMap struct {
	MapBuffs map[TriggerType]*BuffList
}

func NewBuffTriggerMap() *BuffTriggerMap {
	return &BuffTriggerMap{
		MapBuffs: make(map[TriggerType]*BuffList),
	}
}

func (btm *BuffTriggerMap) Add(trigger TriggerType, buff *Buff) {
	lst, isok := btm.MapBuffs[trigger]
	if !isok {
		lst = NewBuffList()

		btm.MapBuffs[trigger] = lst
	}

	lst.Add(buff)
}

func (btm *BuffTriggerMap) AddEx(buff *Buff) {
	for _, v := range buff.Data.Triggers {
		btm.Add(v, buff)
	}
}

func (btm *BuffTriggerMap) RemoveAll() {
	for _, v := range btm.MapBuffs {
		v.RemoveAll()
	}
}

func (btm *BuffTriggerMap) OnTrigger(trigger TriggerType, params *LibFuncParams) (bool, error) {
	lst, isok := btm.MapBuffs[trigger]
	if isok {
		for _, v := range lst.Buffs {
			ret, err := MgrStatic.MgrFunc.Run(v.Data.Trigger, params)
			if err != nil {
				goutils.Error("BuffTriggerMap.OnTrigger:Run",
					zap.Error(err))

				return false, err
			}

			if trigger == TriggerTypeFind {
				if !ret {
					return false, nil
				}
			}
		}
	}

	return true, nil
}
