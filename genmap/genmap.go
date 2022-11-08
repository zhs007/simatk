package genmap

import "math/rand"

func genRoom(minw, minh, maxw, maxh int) (int, int) {
	return minw + rand.Int()%(maxw-minw), minh + rand.Int()%(maxh-minh)
}

func genValidRoomPos(m *Map, rw, rh int) []*Pos {
	// lst := []*Pos{}

	// 最好的位置是四周全部都有东西，刚好只够放下
	lst := m.Foreach(func(cm *Map, x, y int) bool {
		// 先判断外边是否合适
		if !cm.checkAreaOutside(x, y, rw, rh, func(curmap *Map, cx, cy int) bool {
			return curmap.Tile[cx][cy] != TileNone
		}) {
			return false
		}

		// 再判断里面
		return cm.checkArea(x, y, rw, rh, func(curmap *Map, cx, cy int) bool {
			return curmap.Tile[cx][cy] == TileNone
		})
	})

	if len(lst) > 0 {
		return lst
	}

	// 再判断3边
	lst = m.Foreach(func(cm *Map, x, y int) bool {
		// 先判断外边是否合适
		if cm.checkAreaOutsideEx(x, y, rw, rh, func(curmap *Map, cx, cy int) bool {
			return curmap.Tile[cx][cy] != TileNone
		}) != 3 {
			return false
		}

		// 再判断里面
		return cm.checkArea(x, y, rw, rh, func(curmap *Map, cx, cy int) bool {
			return curmap.Tile[cx][cy] == TileNone
		})
	})

	if len(lst) > 0 {
		return lst
	}

	// 再判断2边
	lst = m.Foreach(func(cm *Map, x, y int) bool {
		// 先判断外边是否合适
		if cm.checkAreaOutsideEx(x, y, rw, rh, func(curmap *Map, cx, cy int) bool {
			return curmap.Tile[cx][cy] != TileNone
		}) != 2 {
			return false
		}

		// 再判断里面
		return cm.checkArea(x, y, rw, rh, func(curmap *Map, cx, cy int) bool {
			return curmap.Tile[cx][cy] == TileNone
		})
	})

	if len(lst) > 0 {
		return lst
	}

	return nil
}

func genRoomEx(m *Map, minw, minh, maxw, maxh int) bool {
	var lst []*Pos

	rw, rh := genRoom(minw, minh, maxw, maxh)
	lst0 := genValidRoomPos(m, rw, rh)
	if len(lst0) > 0 {
		lst = append(lst, lst0...)
	}

	if rw != rh {
		lst1 := genValidRoomPos(m, rh, rw)
		if len(lst1) > 0 {
			lst = append(lst, lst1...)
		}
	}

	if len(lst) > 0 {
		ci := rand.Int() % len(lst)

		if ci >= len(lst0) {
			m.SetRoom(lst[ci].X, lst[ci].Y, rh, rw)
		} else {
			m.SetRoom(lst[ci].X, lst[ci].Y, rw, rh)
		}

		return true
	}

	return false
}

func GenMap(w, h int, minw, minh, maxw, maxh int) (*Map, error) {
	m := NewMap(w, h)

	// 决定入口
	lstow := m.Foreach(func(cm *Map, x, y int) bool {
		return cm.IsValidStartOrExit(x, y)
	})

	ci := rand.Int() % len(lstow)
	m.Tile[lstow[ci].X][lstow[ci].Y] = TileStart

	// lstow = m.Foreach(func(cm *Map, x, y int) bool {
	// 	return cm.IsValidStartOrExit(x, y)
	// })

	// ci = rand.Int() % len(lstow)
	// m.Tile[lstow[ci].X][lstow[ci].Y] = TileExit

	// 房间
	failnum := 0
	for {
		if !genRoomEx(m, minw, minh, maxw, maxh) {
			failnum++

			if failnum >= 10 {
				return m, nil
			}
		} else {
			failnum = 0
		}
	}
	// for i := 0; i < 10; i++ {
	// 	rw, rh := genRoom(minw, minh, maxw, maxh)
	// 	lstow = genValidRoomPos(m, rw, rh)
	// 	// if rw != rh {
	// 	// 	lstow1 := genValidRoomPos(m, rh, rw)
	// 	// 	if len(lstow1) > 0 {
	// 	// 		lstow = append(lstow, lstow1...)
	// 	// 	}
	// 	// }

	// 	if len(lstow) > 0 {
	// 		ci := rand.Int() % len(lstow)

	// 		m.SetRoom(lstow[ci].X, lstow[ci].Y, rw, rh)
	// 		// m.Tile[lstow[ci].X][lstow[ci].Y] = TileStart
	// 	}
	// }

	return m, nil
}
