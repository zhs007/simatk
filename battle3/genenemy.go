package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenEnemyParam struct {
	UnitType      int  `json:"unitType"`
	MinTurns      int  `json:"minTurns"`
	MaxTurns      int  `json:"maxTurns"`
	MinLastHP     int  `json:"minLastHP"`
	MaxLastHP     int  `json:"maxLastHP"`
	StartTotalVal int  `json:"startTotalVal"`
	EndTotalVal   int  `json:"endTotalVal"`
	IsWinner      bool `json:"isWinner"`
}

type GenEnemyResultNode struct {
	TotalVal int   `json:"totalVal"` // 总值
	LstHP    []int `json:"lstHP"`    // HP
}

type GenEnemyResult struct {
	Param *GenEnemyParam        `json:"param"` //
	Nodes []*GenEnemyResultNode `json:"nodes"` //
}

func (result *GenEnemyResult) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(result)
	if err != nil {
		goutils.Warn("GenEnemyResult.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("GenEnemyResult.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}

func genEnemyWithTotaVal(hero *Unit, param *GenEnemyParam, totalval int) *GenEnemyResultNode {
	result := &GenEnemyResultNode{
		TotalVal: totalval,
	}

	minhp, maxhp := GetHPAreaForUnitType(param.UnitType, totalval)
	for hp := minhp; hp < maxhp; hp++ {
		if hp == 0 || hp == totalval {
			continue
		}

		dps := totalval - hp

		enemy := NewUnit(hp, dps)
		arr := []*Unit{hero.ResetAndClone(), enemy}

		ret0 := StartBattle(arr)
		if param.IsWinner {
			if ret0.WinIndex == 0 {
				if ret0.Turns >= param.MinTurns && ret0.Turns <= param.MaxTurns {
					if ret0.Units[0].Props[PropTypeCurHP] >= param.MinLastHP &&
						ret0.Units[0].Props[PropTypeCurHP] <= param.MaxLastHP {

						result.LstHP = append(result.LstHP, ret0.Units[1].Props[PropTypeHP])
					}
				}
			}
		} else { // 如果找战斗失败的，只能判断回合数
			if ret0.WinIndex == 1 {
				if ret0.Turns >= param.MinTurns && ret0.Turns <= param.MaxTurns {

					result.LstHP = append(result.LstHP, ret0.Units[1].Props[PropTypeHP])
				}
			}
		}
	}

	if len(result.LstHP) > 0 {
		return result
	}

	return nil
}

func GenEnemy(hero *Unit, param *GenEnemyParam) *GenEnemyResult {
	result := &GenEnemyResult{
		Param: param,
	}

	for tv := param.StartTotalVal; tv <= param.EndTotalVal; tv++ {
		ret := genEnemyWithTotaVal(hero, param, tv)
		if ret != nil {
			result.Nodes = append(result.Nodes, ret)
		}
	}

	return result
}
