package battle3

import (
	"io/ioutil"
	"math/rand"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type EventPool struct {
	Events []*GenEventData `yaml:"events"`
}

func (ep *EventPool) onInit() {
	for _, v := range ep.Events {
		if IsMonster(v.ID) {
			v.EventFunc = append(v.EventFunc, &GenEventFuncData{
				PreFunc: "canwin",
			})
		}
	}
}

func (ep *EventPool) GenEvent(root *Event, unit *Unit) (*Event, bool) {
	curpool := &EventPool{}

	for _, v := range ep.Events {
		if v.TotalNum > 0 {
			num := root.CountID(v.ID)
			if num >= v.TotalNum {
				continue
			}
		}

		isok := true
		for _, f := range v.EventFunc {
			if !MgrStatic.MgrEventFunc.Run(f.PreFunc, v.ID, f.PreFuncParams, f.PreFuncStrParams, unit, root, curpool) {
				isok = false

				break
			}
		}

		if isok {
			curpool.Events = append(curpool.Events, v)
		}
	}

	if len(curpool.Events) > 0 {
		cr := rand.Int() % len(curpool.Events)

		return &Event{
			ID: curpool.Events[cr].ID,
		}, curpool.Events[cr].IsEnding
	}

	return nil, false
}

func LoadEventPool(fn string) (*EventPool, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ep := &EventPool{}
	err = yaml.Unmarshal(data, ep)
	if err != nil {
		goutils.Error("LoadEventPool:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	ep.onInit()

	return ep, nil
}
