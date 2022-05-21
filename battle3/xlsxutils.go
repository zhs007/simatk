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
		f.SetCellStr(sheet, goutils.Pos2Cell(i, 0), fmt.Sprintf("第%v种布局", i+1))

		for j, event := range v {
			if IsItem(event.ID) || IsEquipment(event.ID) {
				data, _ := MgrStatic.MgrItem.GetItemData(event.ID)
				f.SetCellStr(sheet, goutils.Pos2Cell(i, j+1), data.Name)
			} else if IsMonster(event.ID) {
				data, _ := MgrStatic.MgrCharacter.GetCharacterData(event.ID)
				f.SetCellStr(sheet, goutils.Pos2Cell(i, j+1), data.Name)
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
