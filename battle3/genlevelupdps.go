package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenDPSLevelUpResult struct {
	DPS    []int `json:"dps"`    // dps的值
	LoseHP []int `json:"losehp"` // 该dps下，全部战斗完，损失的HP
}

func (result *GenDPSLevelUpResult) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(result)
	if err != nil {
		goutils.Warn("GenDPSLevelUpResult.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("GenDPSLevelUpResult.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}

// 判断dps成长回报，其实只需要考虑损失HP之和即可
func GenDPSLevelUp(index int, off int) (*GenDPSLevelUpResult, error) {
	data := MgrStatic.MgrStageDev.GetData(index)
	if data != nil {
		ret := &GenDPSLevelUpResult{}

		for dpsoff := 0; dpsoff <= off; dpsoff++ {
			unit := NewUnit(data.HP, data.DPS+dpsoff)

			losehp := 0
			for i, v := range data.Monsters {
				for n := 0; n < data.MonsterNums[i]; n++ {
					monster, err := MgrStatic.MgrCharacter.NewUnit(v)
					if err != nil {
						goutils.Error("GenDPSLevelUp:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", v),
							zap.Error(err))

						return nil, err
					}

					ret := startBattle([]*Unit{unit.Clone(), monster})
					if ret == nil {
						goutils.Error("GenDPSLevelUp:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", v),
							zap.Error(ErrBattle))

						return nil, ErrBattle
					}

					losehp += ret.Units[0].Props[PropTypeMaxHP] - ret.Units[0].Props[PropTypeCurHP]
				}
			}

			ret.DPS = append(ret.DPS, data.DPS+dpsoff)
			ret.LoseHP = append(ret.LoseHP, losehp)
		}

		return ret, nil
	}

	goutils.Error("GenDPSLevelUp",
		zap.Int("index", index),
		zap.Error(ErrInvalidStageDevIndex))

	return nil, ErrInvalidStageDevIndex
}
