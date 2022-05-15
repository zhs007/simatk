package battle3

import (
	"io/ioutil"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type EventPool struct {
	Events []*GenEventData `yaml:"events"`
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

	return ep, nil
}
