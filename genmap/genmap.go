package genmap

import "math/rand"

func GenMap(w, h int) (*Map, error) {
	m := NewMap(w, h)

	lstow := m.Foreach(func(cm *Map, x, y int) bool {
		return cm.IsValidStartOrExit(x, y)
	})

	ci := rand.Int() % len(lstow)
	m.Tile[lstow[ci].X][lstow[ci].Y] = TileStart

	lstow = m.Foreach(func(cm *Map, x, y int) bool {
		return cm.IsValidStartOrExit(x, y)
	})

	ci = rand.Int() % len(lstow)
	m.Tile[lstow[ci].X][lstow[ci].Y] = TileExit

	return m, nil
}
