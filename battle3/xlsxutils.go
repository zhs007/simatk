package battle3

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

func SaveEvents(fn string, events []*Event) error {
	f := excelize.NewFile()

	sheet := "Sheet1"

	for i, v := range events {
		f.SetCellStr(sheet, goutils.Pos2Cell(i, 0), fmt.Sprintf("第%v种布局", i+1))

		y := 1
		v.ForEach(func(event *Event) bool {
			if IsItem(event.ID) || IsEquipment(event.ID) {
				data, _ := MgrStatic.MgrItem.GetItemData(event.ID)
				f.SetCellStr(sheet, goutils.Pos2Cell(i, y), data.Name)
				y++
			} else if IsMonster(event.ID) {
				data, _ := MgrStatic.MgrCharacter.GetCharacterData(event.ID)
				f.SetCellStr(sheet, goutils.Pos2Cell(i, y), data.Name)
				y++
			}

			return true
		})
	}

	return f.SaveAs(fn)
}
