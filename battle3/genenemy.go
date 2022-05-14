package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenEnemyParam struct {
	UnitType        int  `json:"unitType"`
	MinTurns        int  `json:"minTurns"`
	MaxTurns        int  `json:"maxTurns"`
	MinLastHP       int  `json:"minLastHP"`
	MaxLastHP       int  `json:"maxLastHP"`
	StartTotalVal   int  `json:"startTotalVal"`
	EndTotalVal     int  `json:"endTotalVal"`
	IsWinner        bool `json:"isWinner"`
	DetailTurnOff   int  `json:"detailTurnOff"`
	DetailLastHPOff int  `json:"detailLastHPOff"`
}

type GenEnemyResultNode struct {
	TotalVal int   `json:"totalVal"` // 总值
	LstHP    []int `json:"lstHP"`    // HP
}

type GenEnemyResultDetail struct {
	MinTurns  int                   `json:"minTurns"`
	MaxTurns  int                   `json:"maxTurns"`
	MinLastHP int                   `json:"minLastHP"`
	MaxLastHP int                   `json:"maxLastHP"`
	Nodes     []*GenEnemyResultNode `json:"nodes"` //
}

func (detail *GenEnemyResultDetail) addDetail(totalval int, hp int) {
	for _, v := range detail.Nodes {
		if v.TotalVal == totalval {
			v.LstHP = append(v.LstHP, hp)

			return
		}
	}

	node := &GenEnemyResultNode{
		TotalVal: totalval,
		LstHP:    []int{hp},
	}

	detail.Nodes = append(detail.Nodes, node)
}

type GenEnemyResult struct {
	Param       *GenEnemyParam          `json:"param"`       //
	Nodes       []*GenEnemyResultNode   `json:"nodes"`       //
	DetailNodes []*GenEnemyResultDetail `json:"detailNodes"` //
}

func (result *GenEnemyResult) RebuildDetail() {
	result.DetailNodes = nil

	if result.Param != nil {
		if result.Param.DetailTurnOff > 0 {
			for t := result.Param.MinTurns + result.Param.DetailTurnOff; t <= result.Param.MaxTurns; t += result.Param.DetailTurnOff {

				detail := &GenEnemyResultDetail{
					MinTurns: t - result.Param.DetailTurnOff,
					MaxTurns: t,
				}

				result.DetailNodes = append(result.DetailNodes, detail)
			}
		}

		if result.Param.DetailLastHPOff > 0 {
			for hp := result.Param.MinLastHP + result.Param.DetailLastHPOff; hp <= result.Param.MaxLastHP; hp += result.Param.DetailLastHPOff {

				detail := &GenEnemyResultDetail{
					MinLastHP: hp - result.Param.DetailLastHPOff,
					MaxLastHP: hp,
				}

				result.DetailNodes = append(result.DetailNodes, detail)
			}
		}
	}
}

func (result *GenEnemyResult) AddNodeForDetail(ret *BattleResult) {
	for _, v := range result.DetailNodes {
		if v.MinTurns > 0 && v.MaxTurns >= v.MinTurns {
			if ret.Turns >= v.MinTurns && ret.Turns < v.MaxTurns {

				// if v.MinLastHP > 0 && v.MaxLastHP >= v.MinLastHP {

				// 	curhp := ret.Units[0].Props[PropTypeCurHP] * 100 / ret.Units[0].Props[PropTypeHP]

				// 	if curhp >= v.MinLastHP &&
				// 		curhp <= v.MaxLastHP {

				// 		v.addDetail(ret.Units[1].Props[PropTypeHP]+ret.Units[1].Props[PropTypeDPS], ret.Units[1].Props[PropTypeHP])

				// 		continue
				// 	}
				// }

				v.addDetail(ret.Units[1].Props[PropTypeHP]+ret.Units[1].Props[PropTypeDPS], ret.Units[1].Props[PropTypeHP])
			}
		}

		if v.MinLastHP > 0 && v.MaxLastHP >= v.MinLastHP {
			curhp := ret.Units[0].Props[PropTypeCurHP] * 100 / ret.Units[0].Props[PropTypeHP]

			if curhp >= v.MinLastHP &&
				curhp < v.MaxLastHP {

				v.addDetail(ret.Units[1].Props[PropTypeHP]+ret.Units[1].Props[PropTypeDPS], ret.Units[1].Props[PropTypeHP])
			}
		}
	}
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

func genEnemyWithTotaVal(hero *Unit, param *GenEnemyParam, totalval int, detailRet *GenEnemyResult) *GenEnemyResultNode {
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

		// 如果是极致DPS，则一定先手
		if param.UnitType == UnitTypeMoreDPS {
			enemy.Props[PropTypeIsFirst] = 1
		}

		arr := []*Unit{hero.ResetAndClone(), enemy}

		ret0 := StartBattle(arr)
		if param.IsWinner {
			if ret0.WinIndex == 0 {
				if ret0.Turns >= param.MinTurns && ret0.Turns <= param.MaxTurns {
					curhp := ret0.Units[0].Props[PropTypeCurHP] * 100 / ret0.Units[0].Props[PropTypeHP]

					if curhp >= param.MinLastHP &&
						curhp <= param.MaxLastHP {
						result.LstHP = append(result.LstHP, ret0.Units[1].Props[PropTypeHP])

						detailRet.AddNodeForDetail(ret0)
					}
				}
			}
		} else { // 如果找战斗失败的，只能判断回合数
			if ret0.WinIndex == 1 {
				if ret0.Turns >= param.MinTurns && ret0.Turns <= param.MaxTurns {

					result.LstHP = append(result.LstHP, ret0.Units[1].Props[PropTypeHP])

					detailRet.AddNodeForDetail(ret0)
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

	result.RebuildDetail()

	for tv := param.StartTotalVal; tv <= param.EndTotalVal; tv++ {
		ret := genEnemyWithTotaVal(hero, param, tv, result)
		if ret != nil {
			result.Nodes = append(result.Nodes, ret)
		}
	}

	return result
}
