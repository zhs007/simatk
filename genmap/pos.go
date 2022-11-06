package genmap

type Pos struct {
	X int
	Y int
}

func NewPos(x, y int) *Pos {
	return &Pos{
		X: x,
		Y: y,
	}
}
