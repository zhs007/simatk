package battle3

import (
	"io/ioutil"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type GenMapParams struct {
	Floor *RandWeights `yaml:"floor"`
	Wall  *RandWeights `yaml:"wall"`
}

func (params *GenMapParams) onLoad() {
	params.Floor.onLoad()
	params.Wall.onLoad()
}

func (params *GenMapParams) GenFloor() int {
	return params.Floor.GenVal()
}

func (params *GenMapParams) GenWall() int {
	return params.Wall.GenVal()
}

func LoadGenMapParams(fn string) (*GenMapParams, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	params := &GenMapParams{}
	err = yaml.Unmarshal(data, params)
	if err != nil {
		goutils.Error("LoadGenMapParams:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	params.onLoad()

	return params, nil
}
