package battle3

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// 根据关卡配置里的主角属性，生成怪物（敌人）
func GenEnemy2(index int) error {
	data := MgrStatic.MgrStageDev.GetData(index)
	if data != nil {
		f := excelize.NewFile()

		unit := NewUnit(data.HP, data.DPS)

		for ut := UnitTypeMoreHP; ut <= UnitTypeMoreDPS; ut++ {
			title, err := UnitType2Str(ut)
			if err != nil {
				goutils.Error("GenEnemy2:UnitType2Str",
					zap.Int("unittype", int(ut)),
					zap.Error(err))

				return err
			}

			ret := GenEnemy(unit,
				&GenEnemyParam{
					Title:           title,
					UnitType:        ut,
					MinTurns:        data.MinTurn[ut-1],
					MaxTurns:        data.MaxTurn[ut-1],
					MinLastHP:       data.MinLastHP[ut-1],
					MaxLastHP:       data.MaxLastHP[ut-1],
					StartTotalVal:   data.MinTotalVal,
					EndTotalVal:     data.MaxTotalVal,
					IsWinner:        true,
					DetailTurnOff:   0,
					DetailLastHPOff: 0,
				})

			ret.OutputExcel(f)
		}

		f.DeleteSheet("Sheet1")

		f.SaveAs(fmt.Sprintf("genenemy-%v.xlsx", data.Name))

		return nil
	}

	goutils.Error("GenEnemy2",
		zap.Int("index", index),
		zap.Error(ErrInvalidStageDevIndex))

	return ErrInvalidStageDevIndex
}
