package battle5

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

// 模拟战斗
func simHeroBattles(hero *Hero, maxpropval int, minpropval int) (int, int, int, int) {
	var totalnum, winnum, losenum, drawnum int

	if maxpropval > minpropval {
		for hp := minpropval; hp < maxpropval; hp++ {
			for atk := minpropval; atk < maxpropval; atk++ {
				for def := minpropval; def < maxpropval; def++ {
					for magic := minpropval; magic < maxpropval; magic++ {
						for speed := minpropval; speed < maxpropval; speed++ {
							target := NewHero(hp, atk, def, magic, speed)

							ret := SimBattle(hero.Clone(), target)
							if ret == 1 {
								winnum++
							} else if ret == -1 {
								losenum++
							} else {
								drawnum++
							}

							totalnum++
						}
					}
				}
			}
		}
	}

	return totalnum, winnum, drawnum, losenum
}

func SimAllBattles(fn string, maxpropval int, minpropval int) error {
	if maxpropval > minpropval {
		f := excelize.NewFile()

		sheet := f.GetSheetName(0)
		f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "hp")
		f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "atk")
		f.SetCellStr(sheet, goutils.Pos2Cell(2, 0), "def")
		f.SetCellStr(sheet, goutils.Pos2Cell(3, 0), "magic")
		f.SetCellStr(sheet, goutils.Pos2Cell(4, 0), "speed")
		f.SetCellStr(sheet, goutils.Pos2Cell(5, 0), "total")
		f.SetCellStr(sheet, goutils.Pos2Cell(6, 0), "win")
		f.SetCellStr(sheet, goutils.Pos2Cell(7, 0), "draw")
		f.SetCellStr(sheet, goutils.Pos2Cell(8, 0), "lose")

		num := 0

		for hp := minpropval; hp < maxpropval; hp++ {
			for atk := minpropval; atk < maxpropval; atk++ {
				for def := minpropval; def < maxpropval; def++ {
					for magic := minpropval; magic < maxpropval; magic++ {
						for speed := minpropval; speed < maxpropval; speed++ {
							num++

							hero := NewHero(hp, atk, def, magic, speed)

							total, win, draw, lose := simHeroBattles(hero, maxpropval, minpropval)

							f.SetCellInt(sheet, goutils.Pos2Cell(0, num), hp)
							f.SetCellInt(sheet, goutils.Pos2Cell(1, num), atk)
							f.SetCellInt(sheet, goutils.Pos2Cell(2, num), def)
							f.SetCellInt(sheet, goutils.Pos2Cell(3, num), magic)
							f.SetCellInt(sheet, goutils.Pos2Cell(4, num), speed)
							f.SetCellInt(sheet, goutils.Pos2Cell(5, num), total)
							f.SetCellInt(sheet, goutils.Pos2Cell(6, num), win)
							f.SetCellInt(sheet, goutils.Pos2Cell(7, num), draw)
							f.SetCellInt(sheet, goutils.Pos2Cell(8, num), lose)

							fmt.Printf("%v %v %v %v %v\n", hp, atk, def, magic, speed)
						}
					}
				}
			}
		}

		f.SaveAs(fn)
	}

	return nil
}
