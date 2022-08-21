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
	Road  *RandWeights `yaml:"road"`
	Start *RandWeights `yaml:"start"`
	Exit  *RandWeights `yaml:"exit"`
	Door  *RandWeights `yaml:"door"`
}

func (params *GenMapStatic) onLoad() {
	params.Floor.onLoad()
	params.Wall.onLoad()
	params.Road.onLoad()
	params.Start.onLoad()
	params.Exit.onLoad()
	params.Door.onLoad()
}

func (params *GenMapStatic) IsFloor(v int) bool {
	return goutils.IndexOfIntSlice(params.Floor.Vals, v, 0) >= 0
}

func (params *GenMapStatic) IsWall(v int) bool {
	return goutils.IndexOfIntSlice(params.Wall.Vals, v, 0) >= 0
}

func (params *GenMapStatic) IsRoomFloor(v int) bool {
	return goutils.IndexOfIntSlice(params.Wall.Vals, v, 0) >= 0 ||
		goutils.IndexOfIntSlice(params.Start.Vals, v, 0) >= 0 ||
		goutils.IndexOfIntSlice(params.Exit.Vals, v, 0) >= 0
}

func (params *GenMapStatic) IsDoor(v int) bool {
	return goutils.IndexOfIntSlice(params.Door.Vals, v, 0) >= 0
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
