package battle3

import (
	"io/ioutil"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type GenMapStatic struct {
	Floor *RandWeights `yaml:"floor"`
	Wall  *RandWeights `yaml:"wall"`
}

func (params *GenMapStatic) onLoad() {
	params.Floor.onLoad()
	params.Wall.onLoad()
}

func (params *GenMapStatic) GenFloor() int {
	return params.Floor.GenVal()
}

func (params *GenMapStatic) GenWall() int {
	return params.Wall.GenVal()
}

func (params *GenMapStatic) IsFloor(v int) bool {
	return goutils.IndexOfIntSlice(params.Floor.Vals, v, 0) >= 0
}

func LoadGenMapStatic(fn string) (*GenMapStatic, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	params := &GenMapStatic{}
	err = yaml.Unmarshal(data, params)
	if err != nil {
		goutils.Error("LoadGenMapStatic:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	params.onLoad()

	return params, nil
}
