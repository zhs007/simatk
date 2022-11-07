package genmap

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type FuncCheck func(m *Map, x, y int) bool
type FuncCheckCross func(t0, t1 TileType) bool

type Map struct {
	Tile [][]TileType
}

func (m *Map) Save(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	for y, arr := range m.Tile {
		for x, v := range arr {
			f.SetCellStr(sheet, goutils.Pos2Cell(x, y), fmt.Sprintf("%v", v))
		}
	}

	return f.SaveAs(fn)
}

func (m *Map) IsValidPos(x, y int) bool {
	return x >= 0 && x < len(m.Tile) && y >= 0 && y < len(m.Tile[0])
}

// 遍历周围8个位置，只要有1个返回true，就返回true
func (m *Map) IsNear(x, y int, check FuncCheck) bool {
	sx := x
	if x > 0 {
		sx = x - 1
	}

	ex := x
	if x < len(m.Tile)-1 {
		ex = x + 1
	}

	sy := y
	if y > 0 {
		sy = y - 1
	}

	ey := y
	if y < len(m.Tile[0])-1 {
		ey = y + 1
	}

	for tx := sx; tx <= ex; tx++ {
		for ty := sy; ty <= ey; ty++ {
			if !(tx == x && ty == y) {
				if check(m, tx, ty) {
					return true
				}
			}
		}
	}

	return false
}

// 是否穿过
func (m *Map) IsCross(x, y int, check0 FuncCheck, check1 FuncCheck) bool {
	if x == 0 {
		if y == 0 || y == len(m.Tile[0])-1 {
			return false
		}

		if check0(m, x+1, y) {
			if check1(m, x, y-1) && check1(m, x, y+1) {
				return true
			}
		} else if check1(m, x+1, y) {
			if check0(m, x, y-1) && check0(m, x, y+1) {
				return true
			}
		}

		return false
	} else if x == len(m.Tile)-1 {
		if y == 0 || y == len(m.Tile[0])-1 {
			return false
		}

		if check0(m, x-1, y) {
			if check1(m, x, y-1) && check1(m, x, y+1) {
				return true
			}
		} else if check1(m, x-1, y) {
			if check0(m, x, y-1) && check0(m, x, y+1) {
				return true
			}
		}

		return false
	}

	if y == 0 {
		if check0(m, x, y+1) {
			if check1(m, x-1, y) && check1(m, x+1, y) {
				return true
			}
		} else if check1(m, x, y+1) {
			if check0(m, x-1, y) && check0(m, x+1, y) {
				return true
			}
		}

		return false
	} else if y == len(m.Tile[0])-1 {
		if check0(m, x, y-1) {
			if check1(m, x-1, y) && check1(m, x+1, y) {
				return true
			}
		} else if check1(m, x, y-1) {
			if check0(m, x-1, y) && check0(m, x+1, y) {
				return true
			}
		}

		return false
	}

	if check0(m, x, y+1) && check0(m, x, y-1) {
		if check1(m, x-1, y) && check1(m, x+1, y) {
			return true
		}
	} else if check1(m, x, y+1) && check1(m, x, y-11) {
		if check0(m, x-1, y) && check0(m, x+1, y) {
			return true
		}
	}

	return false
}

// 是否穿过
func (m *Map) IsCross2(x, y int, check0 FuncCheckCross, check1 FuncCheckCross) bool {
	if x == 0 {
		if y == 0 || y == len(m.Tile[0])-1 {
			return false
		}

		if check0(m.Tile[x+1][y], TileNone) {
			if check1(m.Tile[x][y-1], m.Tile[x][y+1]) {
				return true
			}
		} else if check1(m.Tile[x+1][y], TileNone) {
			if check0(m.Tile[x][y-1], m.Tile[x][y+1]) {
				return true
			}
		}

		return false
	} else if x == len(m.Tile)-1 {
		if y == 0 || y == len(m.Tile[0])-1 {
			return false
		}

		if check0(m.Tile[x-1][y], TileNone) {
			if check1(m.Tile[x][y-1], m.Tile[x][y+1]) {
				return true
			}
		} else if check1(m.Tile[x-1][y], TileNone) {
			if check0(m.Tile[x][y-1], m.Tile[x][y+1]) {
				return true
			}
		}

		return false
	}

	if y == 0 {
		if check0(m.Tile[x][y+1], TileNone) {
			if check1(m.Tile[x-1][y], m.Tile[x+1][y]) {
				return true
			}
		} else if check1(m.Tile[x][y+1], TileNone) {
			if check0(m.Tile[x-1][y], m.Tile[x+1][y]) {
				return true
			}
		}

		return false
	} else if y == len(m.Tile[0])-1 {
		if check0(m.Tile[x][y-1], TileNone) {
			if check1(m.Tile[x-1][y], m.Tile[x+1][y]) {
				return true
			}
		} else if check1(m.Tile[x][y-1], TileNone) {
			if check0(m.Tile[x-1][y], m.Tile[x+1][y]) {
				return true
			}
		}

		return false
	}

	if check0(m.Tile[x][y-1], m.Tile[x][y+1]) {
		if check1(m.Tile[x-1][y], m.Tile[x+1][y]) {
			return true
		}
	} else if check1(m.Tile[x][y-1], m.Tile[x][y+1]) {
		if check0(m.Tile[x-1][y], m.Tile[x+1][y]) {
			return true
		}
	}

	return false
}

// 该坐标是否合适做起点或终点
func (m *Map) IsValidStartOrExit(x, y int) bool {
	if m.IsValidPos(x, y) {
		// 是外墙
		if m.Tile[x][y] == TileOutsideWall {
			// 不在角落
			if x == 0 {
				if y == 0 || y == len(m.Tile[0])-1 {
					return false
				}
			} else if x == len(m.Tile)-1 {
				if y == 0 || y == len(m.Tile[0])-1 {
					return false
				}
			}

			// 周围没有起点或终点
			if m.IsNear(x, y, func(cm *Map, cx, cy int) bool {
				return cm.Tile[cx][cy] == TileStart || cm.Tile[cx][cy] == TileExit
			}) {
				return false
			}

			if m.IsCross2(x, y, func(t0, t1 TileType) bool {
				return t0 == TileOutsideWall && t1 == TileOutsideWall
			}, func(t0, t1 TileType) bool {
				if t0 == TileNone {
					return t1 != TileOutsideWall
				}

				if t1 == TileNone {
					return t0 != TileOutsideWall
				}

				return t0 != TileOutsideWall && t1 != TileOutsideWall
			}) {
				return true
			}
		}
	}

	return false
}

func (m *Map) Foreach(check FuncCheck) []*Pos {
	lst := []*Pos{}

	for x := 0; x < len(m.Tile); x++ {
		for y := 0; y < len(m.Tile[0]); y++ {
			if check(m, x, y) {
				lst = append(lst, NewPos(x, y))
			}
		}
	}

	if len(lst) > 0 {
		return lst
	}

	return nil
}

// 检查4个外边，只要有一个位置返回false，就直接返回false，如果正好在边缘，则不会检查
func (m *Map) checkAreaOutside(x, y int, w, h int, check FuncCheck) bool {
	if x+w >= len(m.Tile) {
		return false
	}

	if y+h >= len(m.Tile[0]) {
		return false
	}

	bx := x
	by := y
	ex := x + w
	ey := y + h

	if x > 0 {
		bx--
	}

	if y > 0 {
		by--
	}

	if x+w < len(m.Tile)-1 {
		ex++
	}

	if y+h < len(m.Tile[0])-1 {
		ey++
	}

	if x > 0 {
		for ty := y; ty <= ey; ty++ {
			if !check(m, bx, ty) {
				return false
			}
		}
	}

	if x+w < len(m.Tile)-1 {
		for ty := y; ty <= ey; ty++ {
			if !check(m, ex, ty) {
				return false
			}
		}
	}

	if y > 0 {
		for tx := x; tx <= ex; tx++ {
			if !check(m, tx, by) {
				return false
			}
		}
	}

	if y+h < len(m.Tile[0])-1 {
		for tx := x; tx <= ex; tx++ {
			if !check(m, tx, ey) {
				return false
			}
		}
	}

	return true
}

// 检查4个外边，有几条边彻底符合check
func (m *Map) checkAreaOutsideEx(x, y int, w, h int, check FuncCheck) int {
	if x+w >= len(m.Tile) {
		return 0
	}

	if y+h >= len(m.Tile[0]) {
		return 0
	}

	bx := x
	by := y
	ex := x + w
	ey := y + h

	if x > 0 {
		bx--
	}

	if y > 0 {
		by--
	}

	if x+w < len(m.Tile)-1 {
		ex++
	}

	if y+h < len(m.Tile[0])-1 {
		ey++
	}

	falsenum := 0

	if x > 0 {
		for ty := y; ty <= ey; ty++ {
			if !check(m, bx, ty) {
				falsenum++

				break
			}
		}
	}

	if x+w < len(m.Tile)-1 {
		for ty := y; ty <= ey; ty++ {
			if !check(m, ex, ty) {
				falsenum++

				break
			}
		}
	}

	if y > 0 {
		for tx := x; tx <= ex; tx++ {
			if !check(m, tx, by) {
				falsenum++

				break
			}
		}
	}

	if y+h < len(m.Tile[0])-1 {
		for tx := x; tx <= ex; tx++ {
			if !check(m, tx, ey) {
				falsenum++

				break
			}
		}
	}

	return 4 - falsenum
}

// 检查区域内，只要有一个位置返回false，就返回false
func (m *Map) checkArea(x, y int, w, h int, check FuncCheck) bool {
	if x+w >= len(m.Tile) {
		return false
	}

	if y+h >= len(m.Tile[0]) {
		return false
	}

	for tx := x; tx <= x+w; tx++ {
		for ty := y; ty <= y+h; ty++ {
			if !check(m, tx, ty) {
				return false
			}
		}
	}

	return true
}

// 检查4个外边，有几条边彻底符合check
func (m *Map) SetRoom(x, y int, w, h int) {
	bx := x
	by := y
	ex := x + w
	ey := y + h

	if x > 0 {
		bx--
	}

	if y > 0 {
		by--
	}

	if x+w < len(m.Tile) {
		ex++
	}

	if y+h < len(m.Tile[0]) {
		ey++
	}

	if x > 0 {
		for ty := y; ty <= ey; ty++ {
			if m.Tile[bx][ty] == TileNone {
				m.Tile[bx][ty] = TileWall
			}
		}
	}

	if x+w < len(m.Tile) {
		for ty := y; ty <= ey; ty++ {
			if m.Tile[ex][ty] == TileNone {
				m.Tile[ex][ty] = TileWall
			}
		}
	}

	if y > 0 {
		for tx := x; tx <= ex; tx++ {
			if m.Tile[tx][by] == TileNone {
				m.Tile[tx][by] = TileWall
			}
		}
	}

	if y+h < len(m.Tile[0]) {
		for tx := x; tx <= ex; tx++ {
			if m.Tile[tx][ey] == TileNone {
				m.Tile[tx][ey] = TileWall
			}
		}
	}

	for tx := x; tx <= x+w; tx++ {
		for ty := y; ty <= y+h; ty++ {
			if m.Tile[tx][ty] == TileNone {
				m.Tile[tx][ty] = TileRoom
			}
		}
	}
}

func NewMap(w, h int) *Map {
	dat := [][]TileType{}

	for x := 0; x < w; x++ {
		arr := []TileType{}

		for y := 0; y < h; y++ {
			arr = append(arr, TileNone)
		}

		dat = append(dat, arr)
	}

	for x := 0; x < w; x++ {
		dat[x][0] = TileOutsideWall
		dat[x][h-1] = TileOutsideWall
	}

	for y := 0; y < h; y++ {
		dat[0][y] = TileOutsideWall
		dat[w-1][y] = TileOutsideWall
	}

	return &Map{
		Tile: dat,
	}
}
