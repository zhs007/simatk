package battle5

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type BuffEffect int

const (
	BuffEffectUnknow       = 0
	BuffEffectTaunt        = 1
	BuffEffectReduceDamage = 2
)

func Str2BuffEffect(str string) (BuffEffect, error) {
	switch str {
	case "TAUNT":
		return BuffEffectTaunt, nil
	case "REDUCE_DAMAGE":
		return BuffEffectReduceDamage, nil
	}

	return BuffEffectUnknow, ErrInvalidBuffEffectString
}

type BuffData struct {
	ID      BuffID
	Name    string
	Effect  BuffEffect
	Level   int
	Turns   int
	Attach  *FuncData
	Trigger *FuncData
	Find    *FuncData
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

			// attachfunc := &FuncData{}
			// findfunc := &FuncData{}
			// triggerfunc := &FuncData{}

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
				case "effect":
					effect, err := Str2BuffEffect(colCell)
					if err != nil {
						goutils.Error("LoadBuffData:effect",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
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
					// case "attachfunc":
					// 	if colCell == "" {
					// 		attachfunc = nil

					// 		break
					// 	}

					// 	attachfunc.FuncName = colCell

					// 	bd.Attach = attachfunc
					// case "attachvals":
					// 	if attachfunc == nil {
					// 		break
					// 	}

					// 	arr := strings.Split(colCell, "|")
					// 	for _, v := range arr {
					// 		v = strings.TrimSpace(v)
					// 		if v != "" {
					// 			i64, err := goutils.String2Int64(v)
					// 			if err != nil {
					// 				goutils.Error("LoadBuffData:attachvals",
					// 					zap.Int("x", x),
					// 					zap.Int("y", y),
					// 					zap.String("cell", colCell),
					// 					zap.Error(err))

					// 				return nil, err
					// 			}

					// 			attachfunc.InVals = append(attachfunc.InVals, int(i64))
					// 		}
					// 	}

					// 	bd.Attach = attachfunc
					// case "attachstrvals":
					// 	if attachfunc == nil {
					// 		break
					// 	}

					// 	arr := strings.Split(colCell, "|")
					// 	for _, v := range arr {
					// 		v = strings.TrimSpace(v)
					// 		if v != "" {
					// 			attachfunc.InStrVals = append(attachfunc.InStrVals, v)
					// 		}
					// 	}

					// 	bd.Attach = attachfunc
					// case "findfunc":
					// 	if colCell == "" {
					// 		findfunc = nil

					// 		break
					// 	}

					// 	findfunc.FuncName = colCell

					// 	bd.Find = findfunc
					// case "findvals":
					// 	if findfunc == nil {
					// 		break
					// 	}

					// 	arr := strings.Split(colCell, "|")
					// 	for _, v := range arr {
					// 		v = strings.TrimSpace(v)
					// 		if v != "" {
					// 			i64, err := goutils.String2Int64(v)
					// 			if err != nil {
					// 				goutils.Error("LoadBuffData:findvals",
					// 					zap.Int("x", x),
					// 					zap.Int("y", y),
					// 					zap.String("cell", colCell),
					// 					zap.Error(err))

					// 				return nil, err
					// 			}

					// 			findfunc.InVals = append(findfunc.InVals, int(i64))
					// 		}
					// 	}

					// 	bd.Find = findfunc
					// case "findstrvals":
					// 	if findfunc == nil {
					// 		break
					// 	}

					// 	arr := strings.Split(colCell, "|")
					// 	for _, v := range arr {
					// 		v = strings.TrimSpace(v)
					// 		if v != "" {
					// 			findfunc.InStrVals = append(findfunc.InStrVals, v)
					// 		}
					// 	}

					// 	bd.Find = findfunc
					// case "triggerfunc":
					// 	if colCell == "" {
					// 		triggerfunc = nil

					// 		break
					// 	}

					// 	triggerfunc.FuncName = colCell

					// 	bd.Trigger = triggerfunc
					// case "triggetvals":
					// 	if triggerfunc == nil {
					// 		break
					// 	}

					// 	arr := strings.Split(colCell, "|")
					// 	for _, v := range arr {
					// 		v = strings.TrimSpace(v)
					// 		if v != "" {
					// 			i64, err := goutils.String2Int64(v)
					// 			if err != nil {
					// 				goutils.Error("LoadBuffData:triggetvals",
					// 					zap.Int("x", x),
					// 					zap.Int("y", y),
					// 					zap.String("cell", colCell),
					// 					zap.Error(err))

					// 				return nil, err
					// 			}

					// 			triggerfunc.InVals = append(triggerfunc.InVals, int(i64))
					// 		}
					// 	}

					// 	bd.Trigger = triggerfunc
					// case "triggerstrvals":
					// 	if triggerfunc == nil {
					// 		break
					// 	}

					// 	arr := strings.Split(colCell, "|")
					// 	for _, v := range arr {
					// 		v = strings.TrimSpace(v)
					// 		if v != "" {
					// 			triggerfunc.InStrVals = append(triggerfunc.InStrVals, v)
					// 		}
					// 	}

					// 	bd.Trigger = triggerfunc
				}
			}

			// MgrStatic.MgrFunc.InitFuncData(findfunc)
			// MgrStatic.MgrFunc.InitFuncData(attachfunc)
			// MgrStatic.MgrFunc.InitFuncData(triggerfunc)

			mgr.mapBuffs[bd.ID] = bd
		}
	}

	return mgr, nil
}
