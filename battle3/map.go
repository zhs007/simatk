package battle3

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type RoomData struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
	X      int `yaml:"x"`
	Y      int `yaml:"y"`
}

func (rd *RoomData) Clone() *RoomData {
	return &RoomData{
		Width:  rd.Width,
		Height: rd.Height,
		X:      rd.X,
		Y:      rd.Y,
	}
}

type MapData struct {
	Data  [][]int     `yaml:"data"`
	Rooms []*RoomData `yaml:"rooms"`
	Start []int       `yaml:"start"`
	Exit  []int       `yaml:"exit"`
}

func (md *MapData) calcInstance(sx, sy int, cx, cy int) int {
	dx := AbsInt(sx - cx)
	dy := AbsInt(sy - cy)

	return dx + dy
}

func (md *MapData) initStartExit(params *GenMapParams) error {
	// 如果没有起点，那么默认起点随机在屏幕最下面
	if len(params.StartPos) != 2 {
		y := params.Height - 1

		md.Start = []int{rand.Int()%(params.Width-2) + 1, y}
	} else {
		md.Start = params.StartPos
	}

	// 如果没有终点，终点应该距离起点至少宽高中更大的距离（水平位移加垂直位移，不是直线距离）
	if len(params.ExitPos) != 2 {
		arr := []int{}
		for y := 1; y < params.Height-1; y++ {
			for x := 1; x < params.Width-1; x++ {
				d := md.calcInstance(md.Start[0], md.Start[1], x, y)
				if d >= params.Height-2 && d >= params.Width-2 {
					arr = append(arr, x, y)
				}
			}
		}

		if len(arr) <= 0 {
			return ErrInvalidExitPos
		}

		r := rand.Int() % (len(arr) / 2)

		md.Exit = []int{arr[r*2], arr[r*2+1]}
	} else {
		md.Exit = params.ExitPos
	}

	return nil
}

func (md *MapData) IsValidPos(x, y int) bool {
	return y >= 0 && y < len(md.Data) && x >= 0 && x < len(md.Data[y])
}

func (md *MapData) isRoomPos(x, y int) bool {
	for _, v := range md.Rooms {
		if v.X == x && v.Y == y {
			return true
		}
	}

	return false
}

func (md *MapData) GenRoomPos(w, h int) []int {
	arr := []int{}

	for y, arr1 := range md.Data {
		for x := range arr1 {
			// 房间不可能重复位置
			if md.isRoomPos(x, y) {
				continue
			}

			if md.isValidRoomArea(x, y, w, h) {
				if md.isValidRoomPos(x, y, w, h) {
					arr = append(arr, x, y, w, h)
				}
			}
		}
	}

	if w != h {
		for y, arr1 := range md.Data {
			for x := range arr1 {
				// 房间不可能重复位置
				if md.isRoomPos(x, y) {
					continue
				}

				if md.isValidRoomArea(x, y, h, w) {
					if md.isValidRoomPos(x, y, h, w) {
						arr = append(arr, x, y, h, w)
					}
				}
			}
		}
	}

	return arr
}

func (md *MapData) GenRooms(lst []int) *MapData {
	w := lst[0]
	h := lst[1]

	lst = lst[2:]

	arr := md.GenRoomPos(w, h)
	if len(arr) >= 4 {
	retry:
		ri := rand.Int() % (len(arr) / 4)
		x := arr[ri*4]
		y := arr[ri*4+1]
		tw := arr[ri*4+2]
		th := arr[ri*4+3]

		nmd := md.Clone()
		nmd.SetRoom(x, y, tw, th)
		if len(lst) == 0 {
			return nmd
		}

		nnmd := nmd.GenRooms(lst)
		if nnmd != nil {
			return nnmd
		}

		arr = arr[4:]
		if len(arr) >= 4 {
			goto retry
		}
	}

	return nil
}

func (md *MapData) SetRoom(sx, sy int, w, h int) {
	md.Rooms = append(md.Rooms, &RoomData{
		Width:  w,
		Height: h,
		X:      sx,
		Y:      sy,
	})

	for tx := 0; tx <= w+1; tx++ {
		for ty := 0; ty <= h+1; ty++ {
			if tx == 0 || tx == w+1 || ty == 0 || ty == h+1 {
				md.Data[ty+sy][tx+sx] = MgrStatic.StaticGenMap.GenWall()
			}
		}
	}
}

func (md *MapData) isValidRoomArea(sx, sy int, w, h int) bool {
	if !md.IsValidPos(sx, sy) || !md.IsValidPos(sx+w+1, sy+h+1) {
		return false
	}

	return true
}

func (md *MapData) isRoomWall(sx, sy int, w, h int, dir int) bool {
	// 0 - up
	// 1 - right
	// 2 - down
	// 3 - left

	// 如果一面是墙，至少有一个角是墙，且墙的区域不少于一半

	if dir == 0 {
		if MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx]) || MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx+w+1]) {
			n := 0
			for tx := 1; tx <= w; tx++ {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx+tx]) {
					n++
				}
			}

			return n >= w/2
		}
	} else if dir == 2 {
		if MgrStatic.StaticGenMap.IsWall(md.Data[sy+h+1][sx]) || MgrStatic.StaticGenMap.IsWall(md.Data[sy+h+1][sx+w+1]) {
			n := 0
			for tx := 1; tx <= w; tx++ {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy+h+1][sx+tx]) {
					n++
				}
			}

			return n >= w/2
		}
	} else if dir == 1 {
		if MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx+w+1]) || MgrStatic.StaticGenMap.IsWall(md.Data[sy+h+1][sx+w+1]) {
			n := 0
			for ty := 1; ty <= h; ty++ {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy+ty][sx+w+1]) {
					n++
				}
			}

			return n >= h/2
		}
	} else if dir == 3 {
		if MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx]) || MgrStatic.StaticGenMap.IsWall(md.Data[sy+h+1][sx]) {
			n := 0
			for ty := 1; ty <= h; ty++ {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy+ty][sx]) {
					n++
				}
			}

			return n >= h/2
		}
	}

	return false
}

func (md *MapData) isRoomDoubleWall(sx, sy int, w, h int, dir int) bool {
	// 0 - up
	// 1 - right
	// 2 - down
	// 3 - left

	// 如果一面是双墙，至少一块就算
	if dir == 0 {
		if sy == 0 {
			return false
		}

		for tx := 1; tx <= w; tx++ {
			if md.IsValidPos(sx+tx, sy-1) {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy-1][sx+tx]) {
					return true
				}
			}
		}

		return false
	} else if dir == 2 {
		if sy+h+1 >= len(md.Data) {
			return false
		}

		for tx := 1; tx <= w; tx++ {
			if md.IsValidPos(sx+tx, sy+h+2) {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy+h+2][sx+tx]) {
					return true
				}
			}
		}

		return false
	} else if dir == 1 {
		if sx+w+1 >= len(md.Data[0]) {
			return false
		}

		for ty := 1; ty <= h; ty++ {
			if md.IsValidPos(sx+w+2, sy+ty) {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy+ty][sx+w+2]) {
					return true
				}
			}
		}

		return false
	} else if dir == 3 {
		if sx == 0 {
			return false
		}

		for ty := 1; ty <= h; ty++ {
			if md.IsValidPos(sx-1, sy+ty) {
				if MgrStatic.StaticGenMap.IsWall(md.Data[sy+ty][sx-1]) {
					return true
				}
			}
		}

		return false
	}

	return true
}

func (md *MapData) isValidRoomPos(sx, sy int, w, h int) bool {
	// 至少有1面临墙（但可以不用全是墙）
	// 不允许出现双墙连起来的情况

	haswall := false
	for dir := 0; dir < 4; dir++ {
		if md.isRoomWall(sx, sy, w, h, dir) {
			haswall = true

			break
		}
	}

	if !haswall {
		return false
	}

	for dir := 0; dir < 4; dir++ {
		if md.isRoomDoubleWall(sx, sy, w, h, dir) {
			return false
		}
	}

	for tx := 1; tx <= w; tx++ {
		for ty := 1; ty <= h; ty++ {
			if MgrStatic.StaticGenMap.IsWall(md.Data[sy+ty][sx+tx]) {
				return false
			}
		}
	}

	return true
}

func (md *MapData) Clone() *MapData {
	nmd := &MapData{}

	for _, arr := range md.Data {
		narr := make([]int, len(arr))
		copy(narr, arr)

		nmd.Data = append(nmd.Data, narr)
	}

	for _, v := range md.Rooms {
		nmd.Rooms = append(nmd.Rooms, v.Clone())
	}

	return nmd
}

// func NewMap(params *GenMapParams) (*MapData, error) {
// 	md := &MapData{}

// 	for y := 0; y < params.Height; y++ {
// 		arr := []int{}
// 		for x := 0; x < params.Width; x++ {
// 			if params.IsWall(x, y) {
// 				arr = append(arr, MgrStatic.StaticGenMap.GenWall())
// 			} else {
// 				arr = append(arr, MgrStatic.StaticGenMap.GenFloor())
// 			}
// 		}

// 		md.Data = append(md.Data, arr)
// 	}

// 	err := md.initStartExit(params)
// 	if err != nil {
// 		goutils.Error("NewMap:initStartExit",
// 			zap.Error(err))

// 		return nil, err
// 	}

// 	md.Data[md.Start[1]][md.Start[0]] = MgrStatic.StaticGenMap.GenStart()
// 	md.Data[md.Exit[1]][md.Exit[0]] = MgrStatic.StaticGenMap.GenExit()

// 	lst := []int{}
// 	for ri := range params.Rooms {
// 		rd := params.Rooms[ri]
// 		for i := 0; i < rd.Num; i++ {
// 			lst = append(lst, rd.Width, rd.Height)
// 		}
// 	}

// 	nmd := md.GenRooms(lst)
// 	if nmd == nil {
// 		goutils.Error("NewMap:GenRooms",
// 			zap.Error(ErrCannotGenMap))

// 		return nil, ErrCannotGenMap
// 	}

// 	nmd.Data[md.Start[1]][md.Start[0]] = MgrStatic.StaticGenMap.GenStart()
// 	nmd.Data[md.Exit[1]][md.Exit[0]] = MgrStatic.StaticGenMap.GenExit()

// 	return nmd, nil
// }

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
