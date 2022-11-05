package battle5

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type SkillData struct {
	ID          SkillID
	Name        string
	Info        string
	Type        SkillType
	ReleaseType ReleaseSkillType
	Atk         *FuncData
	Find        *FuncData
	OnStart     *FuncData
	OnEnd       *FuncData
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

			atkfunc, err := BuildFuncData(header, row, "atk")
			if err != nil {
				goutils.Error("LoadSkillData:atk",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if atkfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(atkfunc)
				sd.Atk = atkfunc
			}

			findfunc, err := BuildFuncData(header, row, "find")
			if err != nil {
				goutils.Error("LoadSkillData:find",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if findfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(findfunc)
				sd.Find = findfunc
			}

			onstartfunc, err := BuildFuncData(header, row, "onstart")
			if err != nil {
				goutils.Error("LoadSkillData:onstart",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if onstartfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(onstartfunc)
				sd.OnStart = onstartfunc
			}

			onendfunc, err := BuildFuncData(header, row, "onend")
			if err != nil {
				goutils.Error("LoadSkillData:onend",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if onendfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(onendfunc)
				sd.OnEnd = onendfunc
			}

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
				case "type":
					sd.Type = Str2SkillType(colCell)
				case "releasetype":
					sd.ReleaseType = Str2ReleaseSkillType(colCell)
				}
			}

			mgr.mapSkills[sd.ID] = sd
		}
	}

	return mgr, nil
}

func Str2SkillType(str string) SkillType {
	switch str {
	case "basicatk":
		return SkillTypeBasicAtk
	case "natural":
		return SkillTypeNatural
	case "ultimate":
		return SkillTypeUltimate
	case "normal":
		return SkillTypeNormal
	}

	return SkillTypeNormal
}

func Str2ReleaseSkillType(str string) ReleaseSkillType {
	switch str {
	case "normal":
		return ReleaseSkillTypeNormal
	case "passive":
		return ReleaseSkillTypePassive
	}

	return ReleaseSkillTypeNormal
}
