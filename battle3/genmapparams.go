package battle3

import (
	"io/ioutil"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type GenMapRoomData struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
	Num    int `yaml:"num"`
}

type GenMapParams struct {
	Width           int               `yaml:"width"`
	Height          int               `yaml:"height"`
	Rooms           []*GenMapRoomData `yaml:"rooms"`
	IsEnclosingWall bool              `yaml:"isEnclosingWall"`
}

func (params *GenMapParams) IsWall(x, y int) bool {
	if params.IsEnclosingWall {
		if x == 0 || x == params.Width-1 {
			return true
		}

		if y == 0 || y == params.Height-1 {
			return true
		}
	}

	return false
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

	return params, nil
}
