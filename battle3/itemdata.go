package battle3

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type ItemData struct {
	ID         int      `yaml:"id"`
	Name       string   `yaml:"name"`
	SPName     string   `yaml:"spname"`
	Type       string   `yaml:"type"`
	ValType    []string `yaml:"valtype"`
	TargetProp []int    `yaml:"-"`
	ValFunc    string   `yaml:"valfunc"`
	Val        []int    `yaml:"val"`
	StrVal     []string `yaml:"strval"`
}

func (item *ItemData) onInit() error {
	for _, v := range item.ValType {
		prop, err := Str2Prop(v)
		if err != nil {
			goutils.Error("ItemData.onInit:Str2Prop",
				zap.String("str", v),
				zap.Error(err))

			return err
		}

		item.TargetProp = append(item.TargetProp, prop)
	}

	return nil
}

type ItemDataMgr struct {
	MapItem map[int]*ItemData
}

func LoadItem(fn string) (*ItemDataMgr, error) {
	mgr := &ItemDataMgr{
		MapItem: make(map[int]*ItemData),
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadItem:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadItem:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		goutils.Error("LoadItem:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			item := &ItemData{}

			for x, colCell := range row {
				switch header[x] {
				case "id":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadItem:id",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					item.ID = int(i64)
				case "name":
					item.Name = colCell
				case "spname":
					item.SPName = colCell
				case "type":
					item.Type = colCell
				case "valtype":
					item.ValType = strings.Split(colCell, "|")
				case "valfunc":
					item.ValFunc = colCell
				case "val":
					arr := strings.Split(colCell, "|")
					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadItem:val",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						item.Val = append(item.Val, int(i64))
					}
				case "strval":
					item.StrVal = strings.Split(colCell, "|")
				}
			}

			err = item.onInit()
			if err != nil {
				goutils.Error("LoadItem:onInit",
					zap.Int("y", y),
					zap.Error(err))
			}

			mgr.MapItem[item.ID] = item
		}
	}

	return mgr, nil
}
