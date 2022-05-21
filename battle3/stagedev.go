package battle3

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// 根据关卡设计和玩家当前状态，来分析玩家成长数据
// 这里就是关卡设计数据

type StageDevData struct {
	Name        string // 名字，易读而已，没用
	Monsters    []int  // 怪物，id列表
	MonsterNums []int  // 怪物数量，和Monsters对应
	HP          int    // 玩家HP
	DPS         int    // 玩家DPS
	Equipments  []int  // 玩家装备
}

type StageDevDataMgr struct {
	lst []*StageDevData
}

func (mgr *StageDevDataMgr) GetData(index int) *StageDevData {
	if index >= 0 && index < len(mgr.lst) {
		return mgr.lst[index]
	}

	goutils.Error("StageDevDataMgr.GetData",
		zap.Int("index", index),
		zap.Int("len", len(mgr.lst)),
		zap.Error(ErrInvalidStageDevIndex))

	return nil
}

func LoadStageDevData(fn string) (*StageDevDataMgr, error) {
	mgr := &StageDevDataMgr{}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadStageDevData:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadStageDevData:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadStageDevData:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			stagedev := &StageDevData{}

			for x, colCell := range row {
				switch header[x] {
				case "name":
					stagedev.Name = colCell
				case "monsters":
					arr := strings.Split(colCell, "|")
					if len(arr)%2 != 0 {
						goutils.Error("LoadStageDevData:monsters",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(ErrInvalidData))

						return nil, ErrInvalidData
					}

					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadStageDevData:monsters",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						if vi%2 == 0 {
							stagedev.Monsters = append(stagedev.Monsters, int(i64))
						} else {
							stagedev.MonsterNums = append(stagedev.MonsterNums, int(i64))
						}
					}
				case "hp":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadStageDevData:hp",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					stagedev.HP = int(i64)
				case "dps":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadStageDevData:dps",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					stagedev.DPS = int(i64)
				case "equip":
					arr := strings.Split(colCell, "|")
					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadStageDevData:equip",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						stagedev.Equipments = append(stagedev.Equipments, int(i64))
					}
				}
			}

			mgr.lst = append(mgr.lst, stagedev)
		}
	}

	return mgr, nil
}
