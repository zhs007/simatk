package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenDPSLevelUp2Result struct {
	DPS    []int `json:"dps"`    // dps的值
	WinNum []int `json:"winnum"` // 该dps下，可能的胜利次数
}

func (result *GenDPSLevelUp2Result) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(result)
	if err != nil {
		goutils.Warn("GenDPSLevelUp2Result.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("GenDPSLevelUp2Result.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}

func GenDPSLevelUp2(index int, off int) (*GenDPSLevelUp2Result, error) {
	data := MgrStatic.MgrStageDev.GetData(index)
	if data != nil {
		ret := &GenDPSLevelUp2Result{}

		for dpsoff := 0; dpsoff <= off; dpsoff++ {
			unit := NewUnit(data.HP, data.DPS+dpsoff)
			winnum := 0

			ForEach(data.Clone(), []int{}, func(arr []int) {
				nu := unit.Clone()
				for _, m := range arr {
					monster, err := MgrStatic.MgrCharacter.NewUnit(m)
					if err != nil {
						goutils.Error("GenDPSLevelUp2:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", m),
							zap.Error(err))

						return
					}

					ret := startBattle([]*Unit{nu, monster})
					if ret == nil {
						goutils.Error("GenDPSLevelUp2:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", m),
							zap.Error(ErrBattle))

						return
					}

					if ret.WinIndex == 1 {
						return
					}
				}

				winnum++
			})

			ret.DPS = append(ret.DPS, data.DPS+dpsoff)
			ret.WinNum = append(ret.WinNum, winnum)
		}

		return ret, nil
	}

	goutils.Error("GenDPSLevelUp2",
		zap.Int("index", index),
		zap.Error(ErrInvalidStageDevIndex))

	return nil, ErrInvalidStageDevIndex
}
