package battle5

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type BuffEffect int

const (
	BuffEffectUnknow    = 0
	BuffEffectTaunt     = 1
	BuffEffectChgDamage = 2
)

func Str2BuffEffect(str string) BuffEffect {
	switch str {
	case "TAUNT":
		return BuffEffectTaunt
	case "CHG_DAMAGE":
		return BuffEffectChgDamage
	}

	return BuffEffectUnknow
}

type BuffData struct {
	ID           BuffID
	Name         string
	Effect       BuffEffect
	Triggers     []TriggerType
	Level        int
	Turns        int
	TriggerTimes int
	Cover        *FuncData
	Attach       *FuncData
	Trigger      *FuncData
	Find         *FuncData
	Remove       *FuncData
}

type BuffDataMgr struct {
	mapBuffs map[BuffID]*BuffData
}

func (mgr *BuffDataMgr) GetBuffData(id BuffID) *BuffData {
	return mgr.mapBuffs[id]
}

func LoadBuffData(fn string) (*BuffDataMgr, error) {
	mgr := &BuffDataMgr{
		mapBuffs: make(map[BuffID]*BuffData),
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadBuffData:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadBuffData:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadBuffData:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			bd := &BuffData{}

			attachfunc, err := BuildFuncData(header, row, "attach")
			if err != nil {
				goutils.Error("LoadBuffData:attach",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if attachfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(attachfunc)
				bd.Attach = attachfunc
			}

			findfunc, err := BuildFuncData(header, row, "find")
			if err != nil {
				goutils.Error("LoadBuffData:find",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if findfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(findfunc)
				bd.Find = findfunc
			}

			triggerfunc, err := BuildFuncData(header, row, "trigger")
			if err != nil {
				goutils.Error("LoadBuffData:trigger",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if triggerfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(triggerfunc)
				bd.Trigger = triggerfunc
			}

			coverfunc, err := BuildFuncData(header, row, "cover")
			if err != nil {
				goutils.Error("LoadBuffData:cover",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if coverfunc != nil {
				MgrStatic.MgrFunc.InitFuncData(coverfunc)
				bd.Cover = coverfunc
			}

			removefunc, err := BuildFuncData(header, row, "remove")
			if err != nil {
				goutils.Error("LoadBuffData:remove",
					zap.Int("y", y),
					zap.Error(err))

				return nil, err
			}

			if removefunc != nil {
				MgrStatic.MgrFunc.InitFuncData(removefunc)
				bd.Remove = removefunc
			}

			for x, colCell := range row {
				switch header[x] {
				case "id":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadBuffData:id",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					bd.ID = BuffID(i64)
				case "name":
					bd.Name = colCell
				case "trigger":
					bd.Triggers = ParseTriggerList(colCell)
				case "effect":
					effect := Str2BuffEffect(colCell)
					if effect == BuffEffectUnknow {
						goutils.Error("LoadBuffData:effect",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(ErrInvalidBuffEffectString))

						return nil, ErrInvalidBuffEffectString
					}

					bd.Effect = effect
				case "level":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadBuffData:level",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					bd.Level = int(i64)
				case "turns":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadBuffData:turns",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					bd.Turns = int(i64)
				case "triggertimes":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadBuffData:triggertimes",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					bd.TriggerTimes = int(i64)
				}
			}

			mgr.mapBuffs[bd.ID] = bd
		}
	}

	return mgr, nil
}
