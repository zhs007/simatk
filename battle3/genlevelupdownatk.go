package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenDownAtkLevelUpResult struct {
	DownAtk []int `json:"downatk"` // downatk的值
	LoseHP  []int `json:"losehp"`  // 该downatk下，全部战斗完，损失的HP
}

func (result *GenDownAtkLevelUpResult) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(result)
	if err != nil {
		goutils.Warn("GenDownAtkLevelUpResult.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("GenDownAtkLevelUpResult.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}

// 判断DownAtk成长回报，其实只需要考虑损失HP之和即可
func GenDownAtkLevelUp(index int, off int) (*GenDownAtkLevelUpResult, error) {
	data := MgrStatic.MgrStageDev.GetData(index)
	if data != nil {
		ret := &GenDownAtkLevelUpResult{}

		for downatk := 0; downatk <= off; downatk++ {
			unit := NewUnit(data.HP, data.DPS)
			unit.Props[PropTypeDownAtk] = data.DownAtk + downatk

			losehp := 0
			for i, v := range data.Monsters {
				for n := 0; n < data.MonsterNums[i]; n++ {
					monster, err := MgrStatic.MgrCharacter.NewUnit(v)
					if err != nil {
						goutils.Error("GenDownAtkLevelUp:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", v),
							zap.Error(err))

						return nil, err
					}

					ret := startBattle([]*Unit{unit.Clone(), monster})
					if ret == nil {
						goutils.Error("GenDownAtkLevelUp:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", v),
							zap.Error(ErrBattle))

						return nil, ErrBattle
					}

					losehp += ret.Units[0].Props[PropTypeMaxHP] - ret.Units[0].Props[PropTypeCurHP]
				}
			}

			ret.DownAtk = append(ret.DownAtk, data.DownAtk+downatk)
			ret.LoseHP = append(ret.LoseHP, losehp)
		}

		return ret, nil
	}

	goutils.Error("GenDownAtkLevelUp",
		zap.Int("index", index),
		zap.Error(ErrInvalidStageDevIndex))

	return nil, ErrInvalidStageDevIndex
}
