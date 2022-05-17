package battle3

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type CharacterData struct {
	ID      int    `yaml:"id"`
	Name    string `yaml:"name"`
	SPName  string `yaml:"spname"`
	HP      int    `yaml:"hp"`
	DPS     int    `yaml:"dps"`
	IsFirst bool   `yaml:"isfirst"`
}

type CharacterDataMgr struct {
	mapCharacter map[int]*CharacterData
}

func (mgr *CharacterDataMgr) GetCharacterData(id int) (*CharacterData, error) {
	if id >= MinCharacterID && id <= MaxCharacterID {
		data, isok := mgr.mapCharacter[id]
		if isok {
			return data, nil
		}
	}

	goutils.Error("CharacterDataMgr.GetCharacterData",
		zap.Int("id", id),
		zap.Error(ErrInvalidCharacterID))

	return nil, ErrInvalidCharacterID
}

func (mgr *CharacterDataMgr) NewUnit(id int) (*Unit, error) {
	if id >= MinCharacterID && id <= MaxCharacterID {
		data, isok := mgr.mapCharacter[id]
		if isok {
			unit := NewUnit(data.HP, data.DPS)

			if data.IsFirst {
				unit.Props[PropTypeIsFirst] = 1
			}

			unit.Data = data

			return unit, nil
		}
	}

	goutils.Error("CharacterDataMgr.NewUnit",
		zap.Int("id", id),
		zap.Error(ErrInvalidCharacterID))

	return nil, ErrInvalidCharacterID
}

func LoadCharacter(fn string) (*CharacterDataMgr, error) {
	mgr := &CharacterDataMgr{
		mapCharacter: make(map[int]*CharacterData),
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadCharacter:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadCharacter:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		goutils.Error("LoadCharacter:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			character := &CharacterData{}

			for x, colCell := range row {
				switch header[x] {
				case "id":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadCharacter:id",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					character.ID = int(i64)
				case "name":
					character.Name = colCell
				case "spname":
					character.SPName = colCell
				case "hp":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadCharacter:hp",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					character.HP = int(i64)
				case "dps":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadCharacter:dps",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					character.DPS = int(i64)
				case "isfirst":
					s := strings.ToLower(colCell)
					if s == "true" {
						character.IsFirst = true
					}
				}
			}

			mgr.mapCharacter[character.ID] = character
		}
	}

	return mgr, nil
}
