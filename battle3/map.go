package battle3

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type MapData struct {
	Data [][]int `yaml:"data"`
}

func NewMap(w, h int, params *GenMapParams) *MapData {
	md := &MapData{}

	for y := 0; y < h; y++ {
		arr := []int{}
		for x := 0; x < w; x++ {
			if y == 0 || y == h-1 {
				arr = append(arr, params.GenWall())
			} else if x == 0 || x == w-1 {
				arr = append(arr, params.GenWall())
			} else {
				arr = append(arr, params.GenFloor())
			}
		}

		md.Data = append(md.Data, arr)
	}

	return md
}

func (md *MapData) Save(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	for y, arr := range md.Data {
		for x, v := range arr {
			f.SetCellStr(sheet, goutils.Pos2Cell(x, y), fmt.Sprintf("%v", v))
		}
	}

	return f.SaveAs(fn)
}
