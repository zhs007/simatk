package genmap

type Map struct {
	Tile [][]TileType
}

func NewMap(w, h int) *Map {
	dat := [][]TileType{}

	for x := 0; x < w; x++ {
		arr := []TileType{}

		for y := 0; y < h; y++ {
			arr = append(arr, TileWall)
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
