package battle5

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type SkillData struct {
	ID   SkillID
	Name string
	Info string
	Atk  *FuncData
	Find *FuncData
}

type SkillDataMgr struct {
	mapSkills map[SkillID]*SkillData
}

func (mgr *SkillDataMgr) GetSkillData(id SkillID) *SkillData {
	return mgr.mapSkills[id]
}

func LoadSkillData(fn string) (*SkillDataMgr, error) {
	mgr := &SkillDataMgr{
		mapSkills: make(map[SkillID]*SkillData),
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadSkillData:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadSkillData:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadSkillData:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			sd := &SkillData{}
			atkfunc := &FuncData{}
			findfunc := &FuncData{}

			for x, colCell := range row {
				switch header[x] {
				case "id":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadSkillData:id",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					sd.ID = SkillID(i64)
				case "name":
					sd.Name = colCell
				case "info":
					sd.Info = colCell
				case "atkfunc":
					atkfunc.FuncName = colCell

					sd.Atk = atkfunc
				case "atkvals":
					arr := strings.Split(colCell, "|")
					for _, v := range arr {
						v = strings.TrimSpace(v)
						if v != "" {
							i64, err := goutils.String2Int64(v)
							if err != nil {
								goutils.Error("LoadSkillData:skills",
									zap.Int("x", x),
									zap.Int("y", y),
									zap.String("cell", colCell),
									zap.Error(err))

								return nil, err
							}

							atkfunc.Vals = append(atkfunc.Vals, int(i64))
						}
					}

					sd.Atk = atkfunc
				case "atkstrvals":
					arr := strings.Split(colCell, "|")
					for _, v := range arr {
						v = strings.TrimSpace(v)
						if v != "" {
							atkfunc.StrVals = append(atkfunc.StrVals, v)
						}
					}

					sd.Atk = atkfunc
				case "findfunc":
					findfunc.FuncName = colCell

					sd.Find = findfunc
				case "findvals":
					arr := strings.Split(colCell, "|")
					for _, v := range arr {
						v = strings.TrimSpace(v)
						if v != "" {
							i64, err := goutils.String2Int64(v)
							if err != nil {
								goutils.Error("LoadSkillData:skills",
									zap.Int("x", x),
									zap.Int("y", y),
									zap.String("cell", colCell),
									zap.Error(err))

								return nil, err
							}

							findfunc.Vals = append(findfunc.Vals, int(i64))
						}
					}

					sd.Find = findfunc
				case "findstrvals":
					arr := strings.Split(colCell, "|")
					for _, v := range arr {
						v = strings.TrimSpace(v)
						if v != "" {
							findfunc.StrVals = append(findfunc.StrVals, v)
						}
					}

					sd.Find = findfunc
				}
			}

			mgr.mapSkills[sd.ID] = sd
		}
	}

	return mgr, nil
}
