package battle3

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type StageData struct {
	ID    int
	Award []int
}

type StageDataMgr struct {
	mapStage map[int]*StageData
}

func (mgr *StageDataMgr) GetData(id int) *StageData {
	d, isok := mgr.mapStage[id]
	if isok {
		return d
	}

	goutils.Error("StageDataMgr.GetData",
		zap.Int("id", id),
		zap.Error(ErrInvalidStageID))

	return nil
}

func LoadStageData(fn string) (*StageDataMgr, error) {
	mgr := &StageDataMgr{
		mapStage: make(map[int]*StageData),
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadStageData:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			goutils.Error("LoadStageData:Close",
				zap.String("fn", fn),
				zap.Error(err))
		}
	}()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadStageData:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	var header []string
	for y, row := range rows {
		if y == 0 {
			header = row
		} else {
			stage := &StageData{}

			for x, colCell := range row {
				colCell = strings.TrimSpace(colCell)

				switch header[x] {
				case "id":
					i64, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadStageDevData:id",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("cell", colCell),
							zap.Error(err))

						return nil, err
					}

					stage.ID = int(i64)
				case "award":
					if colCell == "" {
						continue
					}

					arr := strings.Split(colCell, "|")
					for vi, cv := range arr {
						i64, err := goutils.String2Int64(cv)
						if err != nil {
							goutils.Error("LoadStageDevData:award",
								zap.Int("x", x),
								zap.Int("y", y),
								zap.Int("vi", vi),
								zap.String("cell", colCell),
								zap.String("cv", cv),
								zap.Error(err))

							return nil, err
						}

						stage.Award = append(stage.Award, int(i64))
					}
				}
			}

			mgr.mapStage[stage.ID] = stage
		}
	}

	return mgr, nil
}
