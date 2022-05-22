package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenHPLevelUpResult struct {
	HP     []int `json:"hp"`     // hp的值
	WinNum []int `json:"winnum"` // 该hp下，可能的胜利次数
}

func (result *GenHPLevelUpResult) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(result)
	if err != nil {
		goutils.Warn("GenHPLevelUpResult.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("GenHPLevelUpResult.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}

func GenHPLevelUp(index int, off int) (*GenHPLevelUpResult, error) {
	data := MgrStatic.MgrStageDev.GetData(index)
	if data != nil {
		ret := &GenHPLevelUpResult{}

		for hpoff := 0; hpoff <= off; hpoff++ {
			unit := NewUnit(data.HP+hpoff, data.DPS)
			winnum := 0

			ForEach(data.Clone(), []int{}, func(arr []int) {
				nu := unit.Clone()
				for _, m := range arr {
					monster, err := MgrStatic.MgrCharacter.NewUnit(m)
					if err != nil {
						goutils.Error("GenHPLevelUp:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", m),
							zap.Error(err))

						return
					}

					ret := startBattle([]*Unit{nu, monster})
					if ret == nil {
						goutils.Error("GenHPLevelUp:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", m),
							zap.Error(ErrBattle))

						return
					}

					if ret.WinIndex == 1 {
						return
					}

					winnum++
				}
			})

			ret.HP = append(ret.HP, data.HP+hpoff)
			ret.WinNum = append(ret.WinNum, winnum)
		}

		return ret, nil
	}

	goutils.Error("GenHPLevelUp",
		zap.Int("index", index),
		zap.Error(ErrInvalidStageDevIndex))

	return nil, ErrInvalidStageDevIndex
}
