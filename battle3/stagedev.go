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
	UpAtk       int    // 加攻
	DownAtk     int    // 减攻
	Equipments  []int  // 玩家装备
	MinTurn     []int  // 允许的回合数
	MaxTurn     []int  // 允许的回合数
	MinLastHP   []int  // 剩余血量百分比
	MaxLastHP   []int  // 剩余血量百分比
	MinTotalVal int    // 怪物总属性
	MaxTotalVal int    // 怪物总属性
}

func (data *StageDevData) RemoveMonster(monster int) {
	for i, v := range data.Monsters {
		if v == monster {
			data.MonsterNums[i]--

			if data.MonsterNums[i] <= 0 {
				data.Monsters = append(data.Monsters[:i], data.Monsters[i+1:]...)
				data.MonsterNums = append(data.MonsterNums[:i], data.MonsterNums[i+1:]...)
			}

			return
		}
	}
}

func (data *StageDevData) GetTotalMonsterNum() int {
	num := 0

	for _, v := range data.MonsterNums {
		num += v
	}

	return num
}

func (data *StageDevData) Clone() *StageDevData {
	dst := &StageDevData{
		Name:        data.Name,
		Monsters:    make([]int, len(data.Monsters)),
		MonsterNums: make([]int, len(data.MonsterNums)),
		HP:          data.HP,
		DPS:         data.DPS,
		UpAtk:       data.UpAtk,
		DownAtk:     data.DownAtk,
		Equipments:  make([]int, len(data.Equipments)),
	}

	copy(dst.Monsters, data.Monsters)
	copy(dst.MonsterNums, data.MonsterNums)
	copy(dst.Equipments, data.Equipments)

	return dst
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
				colCell = strings.TrimSpace(colCell)

				switch header[x] {
				case "name":
					stagedev.Name = colCell
				case "monsters":
					if colCell == "" {
						continue
					}

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
				case "upatk":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadStageDevData:upatk",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					stagedev.UpAtk = int(i64)
				case "downatk":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadStageDevData:downatk",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					stagedev.DownAtk = int(i64)
				case "equip":
					if colCell == "" {
						continue
					}

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
				case "turn":
					if colCell == "" {
						continue
					}

					arr := strings.Split(colCell, "|")
					if len(arr)%2 != 0 {
						goutils.Error("LoadStageDevData:turn",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(ErrInvalidData))

						return nil, ErrInvalidData
					}

					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadStageDevData:turn",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						if vi%2 == 0 {
							stagedev.MinTurn = append(stagedev.MinTurn, int(i64))
						} else {
							stagedev.MaxTurn = append(stagedev.MaxTurn, int(i64))
						}
					}
				case "lasthp":
					if colCell == "" {
						continue
					}

					arr := strings.Split(colCell, "|")
					if len(arr)%2 != 0 {
						goutils.Error("LoadStageDevData:lasthp",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(ErrInvalidData))

						return nil, ErrInvalidData
					}

					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadStageDevData:lasthp",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						if vi%2 == 0 {
							stagedev.MinLastHP = append(stagedev.MinLastHP, int(i64))
						} else {
							stagedev.MaxLastHP = append(stagedev.MaxLastHP, int(i64))
						}
					}
				case "totalval":
					if colCell == "" {
						continue
					}

					arr := strings.Split(colCell, "|")
					if len(arr) != 2 {
						goutils.Error("LoadStageDevData:totalval",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(ErrInvalidData))

						return nil, ErrInvalidData
					}

					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadStageDevData:totalval",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						if vi%2 == 0 {
							stagedev.MinTotalVal = int(i64)
						} else {
							stagedev.MaxTotalVal = int(i64)
						}
					}
				}
			}

			mgr.lst = append(mgr.lst, stagedev)
		}
	}

	return mgr, nil
}

type FuncForEachStageDevData func(arr []int)

func ForEach(data *StageDevData, arr []int, each FuncForEachStageDevData) {
	monsterNum := data.GetTotalMonsterNum()
	if monsterNum == 1 {
		arr = append(arr, data.Monsters[0])

		each(arr)

		return
	}

	for _, v := range data.Monsters {
		nd := data.Clone()
		nd.RemoveMonster(v)

		narr := make([]int, len(arr))
		copy(narr, arr)

		narr = append(narr, v)

		ForEach(nd, narr, each)
	}
}
