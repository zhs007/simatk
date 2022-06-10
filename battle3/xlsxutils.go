package battle3

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

func SaveEvents(fn string, events [][]*Event) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	for i, v := range events {
		avglasthpper := CountAvgLastHPPer(v)

		f.SetCellStr(sheet, goutils.Pos2Cell(i, 0), fmt.Sprintf("第%v种布局 %v", i+1, avglasthpper))

		off := 0
		for j, event := range v {
			if IsItem(event.ID) || IsEquipment(event.ID) {
				data, _ := MgrStatic.MgrItem.GetItemData(event.ID)
				f.SetCellStr(sheet, goutils.Pos2Cell(i, j+1+off),
					fmt.Sprintf("%v %v%%->%v%%", data.Name, event.StartHP*100/event.MaxHP, event.EndHP*100/event.MaxHP))
			} else if IsMonster(event.ID) {
				data, _ := MgrStatic.MgrCharacter.GetCharacterData(event.ID)
				f.SetCellStr(sheet, goutils.Pos2Cell(i, j+1+off),
					fmt.Sprintf("%v %v%%->%v%%", data.Name, event.StartHP*100/event.MaxHP, event.EndHP*100/event.MaxHP))
			}

			for _, ee := range event.Awards {
				off++

				if IsItem(ee.ID) || IsEquipment(ee.ID) {
					data, _ := MgrStatic.MgrItem.GetItemData(ee.ID)
					f.SetCellStr(sheet, goutils.Pos2Cell(i, j+1+off),
						fmt.Sprintf("%v %v%%->%v%%", data.Name, ee.StartHP*100/ee.MaxHP, ee.EndHP*100/ee.MaxHP))
				} else if IsMonster(ee.ID) {
					data, _ := MgrStatic.MgrCharacter.GetCharacterData(ee.ID)
					f.SetCellStr(sheet, goutils.Pos2Cell(i, j+1+off),
						fmt.Sprintf("%v %v%%->%v%%", data.Name, ee.StartHP*100/ee.MaxHP, ee.EndHP*100/ee.MaxHP))
				}
			}
		}
	}

	return f.SaveAs(fn)
}

func SaveEvents2(fn string, events []*Event) error {
	f := excelize.NewFile()

	for i, v := range events {
		sheet := fmt.Sprintf("%d", i+1)
		f.NewSheet(sheet)

		v.OutputExcel(f, sheet)
	}

	f.DeleteSheet("Sheet1")

	return f.SaveAs(fn)
}
