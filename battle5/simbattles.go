package battle5

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

// 模拟战斗
func simHeroBattles(hero *Hero, totalval int, maxpropval int, minpropval int) (int, int, int, int) {
	var totalnum, winnum, losenum, drawnum int

	lastval := totalval - minpropval*4

	if lastval > 0 {
		for hp := minpropval; hp < maxpropval && hp < lastval; hp++ {
			lastval0 := totalval - minpropval*3 - hp

			for atk := minpropval; atk < maxpropval && atk < lastval0; atk++ {
				lastval1 := totalval - minpropval*2 - hp - atk

				for def := minpropval; def < maxpropval && def < lastval1; def++ {
					lastval2 := totalval - minpropval - hp - atk - def

					for magic := minpropval; magic < maxpropval && magic < lastval2; magic++ {
						speed := totalval - hp - atk - def - magic

						// for speed := minpropval; speed < maxpropval && speed < lastval3; speed++ {

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
						// }
					}
				}
			}
		}
	}

	return totalnum, winnum, drawnum, losenum
}

func SimAllBattles(fn string, totalval int, maxpropval int, minpropval int) error {
	lastval := totalval - minpropval*4

	if lastval > 0 {
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
		xlsxnum := 0

		chanexcel := make(chan []int, 100)
		go func() {
			for {
				lst := <-chanexcel

				xlsxnum++

				for i, v := range lst {
					f.SetCellInt(sheet, goutils.Pos2Cell(i, xlsxnum), v)
				}
			}
		}()

		chanend := make(chan int, 100)

		startnum := 0
		endnum := 0
		for hp := minpropval; hp < maxpropval && hp < lastval; hp++ {
			startnum++
			go func(curhp int) {
				lastval0 := totalval - minpropval*3 - curhp

				for atk := minpropval; atk < maxpropval && atk < lastval0; atk++ {
					lastval1 := totalval - minpropval*2 - curhp - atk

					for def := minpropval; def < maxpropval && def < lastval1; def++ {
						lastval2 := totalval - minpropval - curhp - atk - def

						for magic := minpropval; magic < maxpropval && magic < lastval2; magic++ {
							speed := totalval - curhp - atk - def - magic

							// for speed := minpropval; speed < maxpropval && speed < lastval3; speed++ {
							num++

							hero := NewHero(curhp, atk, def, magic, speed)

							total, win, draw, lose := simHeroBattles(hero, totalval, maxpropval, minpropval)

							chanexcel <- []int{curhp, atk, def, magic, speed, total, win, draw, lose}

							// f.SetCellInt(sheet, goutils.Pos2Cell(0, num), hp)
							// f.SetCellInt(sheet, goutils.Pos2Cell(1, num), atk)
							// f.SetCellInt(sheet, goutils.Pos2Cell(2, num), def)
							// f.SetCellInt(sheet, goutils.Pos2Cell(3, num), magic)
							// f.SetCellInt(sheet, goutils.Pos2Cell(4, num), speed)
							// f.SetCellInt(sheet, goutils.Pos2Cell(5, num), total)
							// f.SetCellInt(sheet, goutils.Pos2Cell(6, num), win)
							// f.SetCellInt(sheet, goutils.Pos2Cell(7, num), draw)
							// f.SetCellInt(sheet, goutils.Pos2Cell(8, num), lose)

							fmt.Printf("%v %v %v %v %v\n", curhp, atk, def, magic, speed)
							// }
						}
					}
				}

				chanend <- 1
			}(hp)
		}

		for {
			<-chanend

			endnum++

			if endnum == startnum && startnum == maxpropval-minpropval {
				break
			}
		}

		for {
			if num == xlsxnum {
				break
			}

			time.Sleep(time.Second)
		}

		f.SaveAs(fn)
	}

	return nil
}
