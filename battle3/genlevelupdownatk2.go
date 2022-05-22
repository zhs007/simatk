package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenDownAtkLevelUp2Result struct {
	DownAtk []int `json:"downatk"` // downatk的值
	WinNum  []int `json:"winnum"`  // 可能的胜利次数
}

func (result *GenDownAtkLevelUp2Result) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(result)
	if err != nil {
		goutils.Warn("GenDownAtkLevelUp2Result.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("GenDownAtkLevelUp2Result.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}

// 判断DownAtk成长回报，其实只需要考虑损失HP之和即可
func GenDownAtkLevelUp2(index int, off int) (*GenDownAtkLevelUp2Result, error) {
	data := MgrStatic.MgrStageDev.GetData(index)
	if data != nil {
		ret := &GenDownAtkLevelUp2Result{}

		for downatk := 0; downatk <= off; downatk++ {
			unit := NewUnit(data.HP, data.DPS)
			unit.Props[PropTypeDownAtk] = data.DownAtk + downatk

			winnum := 0

			ForEach(data.Clone(), []int{}, func(arr []int) {
				nu := unit.Clone()
				for _, m := range arr {
					monster, err := MgrStatic.MgrCharacter.NewUnit(m)
					if err != nil {
						goutils.Error("GenDownAtkLevelUp2:NewUnit",
							zap.Int("index", index),
							zap.Int("monster", m),
							zap.Error(err))

						return
					}

					ret := startBattle([]*Unit{nu, monster})
					if ret == nil {
						goutils.Error("GenDownAtkLevelUp2:NewUnit",
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

			ret.DownAtk = append(ret.DownAtk, data.DownAtk+downatk)
			ret.WinNum = append(ret.WinNum, winnum)
		}

		return ret, nil
	}

	goutils.Error("GenDownAtkLevelUp2",
		zap.Int("index", index),
		zap.Error(ErrInvalidStageDevIndex))

	return nil, ErrInvalidStageDevIndex
}
