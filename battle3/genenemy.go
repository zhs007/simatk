package battle3

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GenEnemyParam struct {
	Title           string   `json:"title"`
	UnitType        UnitType `json:"unitType"`
	MinTurns        int      `json:"minTurns"`
	MaxTurns        int      `json:"maxTurns"`
	MinLastHP       int      `json:"minLastHP"`
	MaxLastHP       int      `json:"maxLastHP"`
	StartTotalVal   int      `json:"startTotalVal"`
	EndTotalVal     int      `json:"endTotalVal"`
	IsWinner        bool     `json:"isWinner"`
	DetailTurnOff   int      `json:"detailTurnOff"`
	DetailLastHPOff int      `json:"detailLastHPOff"`
}

type GenEnemyBattleResult struct {
	LastHP int `json:"lastHP"` // 获胜者剩下的HP
	MaxHP  int `json:"maxHP"`  // 获胜者的总HP
	Turns  int `json:"turns"`  // 战斗回合数
}

type GenEnemyResultNode struct {
	TotalVal        int                     `json:"totalVal"`        // 总值
	LstHP           []int                   `json:"lstHP"`           // HP
	LstBattleResult []*GenEnemyBattleResult `json:"lstBattleResult"` // 战斗结果结算
}

func (node *GenEnemyResultNode) AddBattleResult(param *GenEnemyParam, ret *BattleResult) {
	// 玩家胜利
	if param.IsWinner {
		node.LstHP = append(node.LstHP, ret.Units[1].Props[PropTypeHP])

		node.LstBattleResult = append(node.LstBattleResult, &GenEnemyBattleResult{
			LastHP: ret.Units[0].Props[PropTypeCurHP],
			MaxHP:  ret.Units[0].Props[PropTypeHP],
			Turns:  ret.Turns,
		})
	} else {
		node.LstHP = append(node.LstHP, ret.Units[1].Props[PropTypeHP])

		node.LstBattleResult = append(node.LstBattleResult, &GenEnemyBattleResult{
			LastHP: ret.Units[1].Props[PropTypeCurHP],
			MaxHP:  ret.Units[1].Props[PropTypeHP],
			Turns:  ret.Turns,
		})
	}
}

type GenEnemyResultDetail struct {
	MinTurns  int                   `json:"minTurns"`
	MaxTurns  int                   `json:"maxTurns"`
	MinLastHP int                   `json:"minLastHP"`
	MaxLastHP int                   `json:"maxLastHP"`
	Nodes     []*GenEnemyResultNode `json:"nodes"` //
}

func (detail *GenEnemyResultDetail) addDetail(totalval int, param *GenEnemyParam, ret *BattleResult) {
	for _, v := range detail.Nodes {
		if v.TotalVal == totalval {
			v.AddBattleResult(param, ret)

			return
		}
	}

	node := &GenEnemyResultNode{
		TotalVal: totalval,
	}

	node.AddBattleResult(param, ret)

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
		if result.Param.DetailTurnOff > 0 && result.Param.DetailLastHPOff > 0 {

			for t := result.Param.MinTurns + result.Param.DetailTurnOff; t <= result.Param.MaxTurns; t += result.Param.DetailTurnOff {
				for hp := result.Param.MinLastHP + result.Param.DetailLastHPOff; hp <= result.Param.MaxLastHP; hp += result.Param.DetailLastHPOff {

					detail := &GenEnemyResultDetail{
						MinTurns:  t - result.Param.DetailTurnOff,
						MaxTurns:  t,
						MinLastHP: hp - result.Param.DetailLastHPOff,
						MaxLastHP: hp,
					}

					result.DetailNodes = append(result.DetailNodes, detail)
				}

				continue
			}
		} else if result.Param.DetailTurnOff > 0 {

			for t := result.Param.MinTurns + result.Param.DetailTurnOff; t <= result.Param.MaxTurns; t += result.Param.DetailTurnOff {
				detail := &GenEnemyResultDetail{
					MinTurns: t - result.Param.DetailTurnOff,
					MaxTurns: t,
				}

				result.DetailNodes = append(result.DetailNodes, detail)
			}
		} else if result.Param.DetailLastHPOff > 0 {

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

func (result *GenEnemyResult) AddNodeForDetail(param *GenEnemyParam, ret *BattleResult) {
	for _, v := range result.DetailNodes {
		if v.MinTurns > 0 && v.MaxTurns >= v.MinTurns && v.MinLastHP > 0 && v.MaxLastHP >= v.MinLastHP {
			if ret.Turns >= v.MinTurns && ret.Turns < v.MaxTurns {

				curhp := ret.Units[0].Props[PropTypeCurHP] * 100 / ret.Units[0].Props[PropTypeHP]

				if curhp >= v.MinLastHP &&
					curhp <= v.MaxLastHP {

					v.addDetail(ret.Units[1].Props[PropTypeHP]+ret.Units[1].Props[PropTypeDPS], param, ret)
				}
			}
		} else if v.MinTurns > 0 && v.MaxTurns >= v.MinTurns {
			if ret.Turns >= v.MinTurns && ret.Turns < v.MaxTurns {

				v.addDetail(ret.Units[1].Props[PropTypeHP]+ret.Units[1].Props[PropTypeDPS], param, ret)
			}
		} else if v.MinLastHP > 0 && v.MaxLastHP >= v.MinLastHP {
			curhp := ret.Units[0].Props[PropTypeCurHP] * 100 / ret.Units[0].Props[PropTypeHP]

			if curhp >= v.MinLastHP &&
				curhp < v.MaxLastHP {

				v.addDetail(ret.Units[1].Props[PropTypeHP]+ret.Units[1].Props[PropTypeDPS], param, ret)
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

func (result *GenEnemyResult) OutputExcel(f *excelize.File) error {
	sheet := result.Param.Title
	f.DeleteSheet(sheet)
	f.NewSheet(sheet)

	lsthead := []string{
		"totalval",
		"hp",
		"dps",
		"lasthpper",
		"turns",
	}

	for x, v := range lsthead {
		f.SetCellStr(sheet, goutils.Pos2Cell(x, 0), v)
	}

	y := 1

	for _, v := range result.Nodes {
		for i, curhp := range v.LstHP {
			f.SetCellInt(sheet, goutils.Pos2Cell(0, y), v.TotalVal)
			f.SetCellInt(sheet, goutils.Pos2Cell(1, y), curhp)
			f.SetCellInt(sheet, goutils.Pos2Cell(2, y), v.TotalVal-curhp)
			f.SetCellFloat(sheet, goutils.Pos2Cell(3, y), float64(v.LstBattleResult[i].LastHP)/float64(v.LstBattleResult[i].MaxHP), 2, 32)
			f.SetCellInt(sheet, goutils.Pos2Cell(4, y), v.LstBattleResult[i].Turns)

			y++
		}
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
		if param.UnitType == UnitTypeMoreDPS || param.UnitType == UnitTypeDPS {
			enemy.Props[PropTypeIsFirst] = 1
		}

		arr := []*Unit{hero.Clone(), enemy}

		ret0 := StartBattle(arr)
		if param.IsWinner {
			if ret0.WinIndex == 0 {
				if ret0.Turns >= param.MinTurns && ret0.Turns <= param.MaxTurns {
					curhp := ret0.Units[0].Props[PropTypeCurHP] * 100 / ret0.Units[0].Props[PropTypeHP]

					if curhp >= param.MinLastHP &&
						curhp <= param.MaxLastHP {
						result.AddBattleResult(param, ret0)

						detailRet.AddNodeForDetail(param, ret0)
					}
				}
			}
		} else { // 如果找战斗失败的，只能判断回合数
			if ret0.WinIndex == 1 {
				if ret0.Turns >= param.MinTurns && ret0.Turns <= param.MaxTurns {

					result.AddBattleResult(param, ret0)

					detailRet.AddNodeForDetail(param, ret0)
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
