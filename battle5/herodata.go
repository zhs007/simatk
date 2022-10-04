package battle5

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type HeroData struct {
	ID          HeroID
	Name        string
	HP          int
	Atk         int
	Def         int
	Magic       int
	Speed       int
	MovDistance int
	AtkDistance int
	Place       int
	Info        string
	X, Y        int
	Skills      []SkillID
}

func (hd *HeroData) Clone() *HeroData {
	return &HeroData{
		ID:          hd.ID,
		Name:        hd.Name,
		HP:          hd.HP,
		Atk:         hd.Atk,
		Def:         hd.Def,
		Magic:       hd.Magic,
		Speed:       hd.Speed,
		MovDistance: hd.MovDistance,
		AtkDistance: hd.AtkDistance,
		Place:       hd.Place,
		Info:        hd.Info,
		X:           hd.X,
		Y:           hd.Y,
		Skills:      append([]SkillID{}, hd.Skills...),
	}
}

type HeroDataMgr struct {
	mapHeros map[HeroID]*HeroData
}

func (mgr *HeroDataMgr) GetHeroData(id HeroID) *HeroData {
	return mgr.mapHeros[id]
}

func LoadHeroData(fn string) (*HeroDataMgr, error) {
	mgr := &HeroDataMgr{
		mapHeros: make(map[HeroID]*HeroData),
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadHeroData:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadHeroData:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadHeroData:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			hd := &HeroData{
				X: -1,
				Y: -1,
			}

			for x, colCell := range row {
				switch header[x] {
				case "id":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:id",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.ID = HeroID(i64)
				case "name":
					hd.Name = colCell
				case "hp":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:hp",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.HP = int(i64)
				case "atk":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:atk",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.Atk = int(i64)
				case "def":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:def",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.Def = int(i64)
				case "magic":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:magic",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.Magic = int(i64)
				case "speed":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:speed",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.Speed = int(i64)
				case "movdistance":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:movdistance",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.MovDistance = int(i64)
				case "atkdistance":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:atkdistance",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.AtkDistance = int(i64)
				case "place":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadHeroData:place",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					hd.Place = int(i64)
				case "info":
					hd.Info = colCell
				case "skills":
					arr := strings.Split(colCell, "|")
					for _, v := range arr {
						v = strings.TrimSpace(v)
						if v != "" {
							i64, err := goutils.String2Int64(v)
							if err != nil {
								goutils.Error("LoadHeroData:skills",
									zap.Int("x", x),
									zap.Int("y", y),
									zap.String("cell", colCell),
									zap.Error(err))

								return nil, err
							}

							hd.Skills = append(hd.Skills, SkillID(i64))
						}
					}
				}
			}

			mgr.mapHeros[hd.ID] = hd
		}
	}

	return mgr, nil
}
