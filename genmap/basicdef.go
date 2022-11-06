package genmap

type TileType uint8

const (
	TileWall        TileType = 0 // 墙
	TileRoad        TileType = 1
	TileRoom        TileType = 2
	TileDoor        TileType = 3
	TileUpStair     TileType = 4
	TileDownStair   TileType = 5
	TileOutsideWall TileType = 6 // 外墙，外墙上不可以有门
)
