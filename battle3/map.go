package battle3

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type RoomData struct {
	Width  int `yaml:"width" json:"width"`
	Height int `yaml:"height" json:"height"`
	X      int `yaml:"x" json:"x"`
	Y      int `yaml:"y" json:"y"`
}

func (rd *RoomData) IsInRoom(x, y int) bool {
	return x >= rd.X && x <= rd.X+rd.Width+1 && y >= rd.Y && y <= rd.Y+rd.Height+1
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
	Data  [][]int     `yaml:"data" json:"data"`
	Rooms []*RoomData `yaml:"rooms" json:"rooms"`
	Start []int       `yaml:"start" json:"start"`
	Exit  []int       `yaml:"exit" json:"exit"`
}

func (md *MapData) GetWidth() int {
	if md.Data == nil {
		return 0
	}

	return len(md.Data[0])
}

func (md *MapData) GetHeight() int {
	if md.Data == nil {
		return 0
	}

	return len(md.Data)
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

	nmd.Start = goutils.CloneIntArr(md.Start)
	nmd.Exit = goutils.CloneIntArr(md.Exit)

	return nmd
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
	return y >= 0 && y < md.GetHeight() && x >= 0 && x < md.GetWidth()
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

func (md *MapData) hasRoomDoor(sx, sy int, w, h int, dir int) bool {
	// 0 - up
	// 1 - right
	// 2 - down
	// 3 - left

	if dir == 0 {
		for tx := 1; tx <= w; tx++ {
			if MgrStatic.StaticGenMap.IsDoor(md.Data[sy][sx+tx]) {
				return true
			}
		}
	} else if dir == 2 {
		for tx := 1; tx <= w; tx++ {
			if MgrStatic.StaticGenMap.IsDoor(md.Data[sy+h+1][sx+tx]) {
				return true
			}
		}
	} else if dir == 1 {
		for ty := 1; ty <= h; ty++ {
			if MgrStatic.StaticGenMap.IsDoor(md.Data[sy+ty][sx+w+1]) {
				return true
			}
		}
	} else if dir == 3 {
		for ty := 1; ty <= h; ty++ {
			if MgrStatic.StaticGenMap.IsDoor(md.Data[sy+ty][sx]) {
				return true
			}
		}
	}

	return false
}

func (md *MapData) isValidRoomDoorPos(x, y int) bool {
	if x == 0 || y == 0 || x == md.GetWidth()-1 || y == md.GetHeight()-1 {
		return false
	}

	if !MgrStatic.StaticGenMap.IsWall(md.Data[y][x]) {
		return false
	}

	wallnum := 0
	if !MgrStatic.StaticGenMap.IsWall(md.Data[y-1][x]) {
		wallnum++
	}
	if !MgrStatic.StaticGenMap.IsWall(md.Data[y+1][x]) {
		wallnum++
	}
	if !MgrStatic.StaticGenMap.IsWall(md.Data[y][x-1]) {
		wallnum++
	}
	if !MgrStatic.StaticGenMap.IsWall(md.Data[y][x+1]) {
		wallnum++
	}
	if wallnum != 2 {
		return false
	}

	if !MgrStatic.StaticGenMap.IsNotNearWall(md.Data[y-1][x]) {
		return false
	}
	if !MgrStatic.StaticGenMap.IsNotNearWall(md.Data[y+1][x]) {
		return false
	}
	if !MgrStatic.StaticGenMap.IsNotNearWall(md.Data[y][x-1]) {
		return false
	}
	if !MgrStatic.StaticGenMap.IsNotNearWall(md.Data[y][x+1]) {
		return false
	}

	return true
}

func (md *MapData) genRoomDoorPos(sx, sy int, w, h int) (int, int) {
	// 一面墙就一扇门
	// 不能在外墙上
	// 在墙的中间区域
	// 门四周不能超过3面墙
	// 门四周不能有门、起点、终点

	arr := []int{}

	if sy > 0 && !md.hasRoomDoor(sx, sy, w, h, 0) {
		ss := w/4 + 1
		es := 3 * w / 4
		if ss == es {
			es++
		}

		for tx := ss; tx < es; tx++ {
			if md.isValidRoomDoorPos(sx+tx, sy) {
				arr = append(arr, sx+tx, sy)
			}
		}
	}

	if sy < md.GetHeight()-h-2 && !md.hasRoomDoor(sx, sy, w, h, 2) {
		ss := w/4 + 1
		es := 3 * w / 4
		if ss == es {
			es++
		}

		for tx := ss; tx < es; tx++ {
			if md.isValidRoomDoorPos(sx+tx, sy) {
				arr = append(arr, sx+tx, sy)
			}
		}
	}

	if sx > 0 && !md.hasRoomDoor(sx, sy, w, h, 3) {
		ss := h/4 + 1
		es := 3 * h / 4
		if ss == es {
			es++
		}

		for ty := ss; ty < es; ty++ {
			arr = append(arr, sx, sy+ty)
		}
	}

	if sx < md.GetWidth()-w-2 && !md.hasRoomDoor(sx, sy, w, h, 1) {
		ss := h/4 + 1
		es := 3 * h / 4
		if ss == es {
			es++
		}

		for ty := ss; ty < es; ty++ {
			arr = append(arr, sx+w+1, sy+ty)
		}
	}

	if len(arr) >= 2 {
		r := rand.Int() % (len(arr) / 2)
		return arr[r*2], arr[r*2+1]
	}

	return -1, -1
}

func (md *MapData) genRoomDoor(sx, sy int, w, h int) {
	tx, ty := md.genRoomDoorPos(sx, sy, w, h)
	if tx != -1 && ty != -1 {
		md.Data[ty][tx] = MgrStatic.StaticGenMap.Door.GenVal()
	}
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
				md.Data[ty+sy][tx+sx] = MgrStatic.StaticGenMap.Wall.GenVal()
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
		if sy+h+1 >= md.GetHeight() {
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
		if sx+w+1 >= md.GetWidth() {
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
	// 且房间的墙不允许是入口或出口（地图外墙除外）

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

	if sx != 0 {
		for ty := 0; ty <= h+1; ty++ {
			if MgrStatic.StaticGenMap.IsNonRoomWall(md.Data[sy+ty][sx]) {
				return false
			}
		}
	}

	if sx+w+1 != md.GetWidth()-1 {
		for ty := 0; ty <= h+1; ty++ {
			if MgrStatic.StaticGenMap.IsNonRoomWall(md.Data[sy+ty][sx+w+1]) {
				return false
			}
		}
	}

	if sy != 0 {
		for tx := 0; tx <= w+1; tx++ {
			if MgrStatic.StaticGenMap.IsNonRoomWall(md.Data[sy][sx+tx]) {
				return false
			}
		}
	}

	if sy+h+1 != md.GetHeight()-1 {
		for tx := 0; tx <= w+1; tx++ {
			if MgrStatic.StaticGenMap.IsNonRoomWall(md.Data[sy+h+1][sx+tx]) {
				return false
			}
		}
	}

	return true
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

func (md *MapData) ToXlsx(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	for y, arr := range md.Data {
		for x, v := range arr {
			f.SetCellStr(sheet, goutils.Pos2Cell(x, y), fmt.Sprintf("%v", v))
		}
	}

	return f.SaveAs(fn)
}

func (md *MapData) ToYaml(fn string) error {
	data, err := yaml.Marshal(md)
	if err != nil {
		goutils.Error("MapData:ToYaml:Marshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	err = ioutil.WriteFile(fn, data, 0644)
	if err != nil {
		goutils.Error("MapData:ToYaml:WriteFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return nil
}

func (md *MapData) ToJson(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := json.Marshal(md)
	if err != nil {
		goutils.Error("MapData:ToJson:Marshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	err = ioutil.WriteFile(fn, data, 0644)
	if err != nil {
		goutils.Error("MapData:ToJson:WriteFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return nil
}

func (md *MapData) IsInRoom(x, y int) bool {
	for _, v := range md.Rooms {
		if v.IsInRoom(x, y) {
			return true
		}
	}

	return false
}

func (md *MapData) getMinRoom(x, y int) *RoomData {
	if x == 0 || y == 0 || x == md.GetWidth()-1 || y == md.GetHeight()-1 {
		return nil
	}

	// 因为传入参数是最初发现的可用点，所以不需要太复杂的判断就能得到sx和sy
	sx := x
	for tx := x; tx > 0 && MgrStatic.StaticGenMap.IsRoomFloor(md.Data[y][tx]); tx-- {
		sx = tx
	}
	sx--

	sy := y
	for ty := y; ty > 0 && MgrStatic.StaticGenMap.IsRoomFloor(md.Data[ty][x]); ty-- {
		sy = ty
	}
	sy--

	// 取到ex和ey
	ex := x
	ey := y
	xend := false
	yend := false
	ti := 1

	for !xend || !yend {
		if !xend && !yend {
			if ex+1 < md.GetWidth()-1 && ey+1 < md.GetHeight()-1 && MgrStatic.StaticGenMap.IsRoomFloor(md.Data[ey+1][ex+1]) {
				for i := 0; i < ti; i++ {
					if !MgrStatic.StaticGenMap.IsRoomFloor(md.Data[ey+1][x+i]) {
						yend = true
					}

					if !MgrStatic.StaticGenMap.IsRoomFloor(md.Data[y+i][ex+1]) {
						xend = true
					}
				}
			} else {
				if ex+1 >= md.GetWidth()-1 {
					xend = true
				} else {
					for i := 0; i < ti; i++ {
						if !MgrStatic.StaticGenMap.IsRoomFloor(md.Data[y+i][ex+1]) {
							xend = true
						}
					}
				}

				if ey+1 >= md.GetHeight()-1 {
					yend = true
				} else {
					for i := 0; i < ti; i++ {
						if !MgrStatic.StaticGenMap.IsRoomFloor(md.Data[ey+1][x+i]) {
							yend = true
						}
					}
				}
			}
		} else if !xend {
			if ex+1 >= md.GetWidth()-1 {
				xend = true
			}

			for i := y; i <= ey; i++ {
				if !MgrStatic.StaticGenMap.IsRoomFloor(md.Data[i][ex+1]) {
					xend = true
				}
			}
		} else if !yend {
			if ey+1 >= md.GetHeight()-1 {
				yend = true
			}

			for i := x; i <= ex; i++ {
				if !MgrStatic.StaticGenMap.IsRoomFloor(md.Data[ey+1][i]) {
					yend = true
				}
			}
		}

		if !xend {
			ex++
		}

		if !yend {
			ey++
		}

		ti++
	}

	if ex-sx == 1 || ey-sy == 1 {
		return nil
	}

	return &RoomData{
		X:      sx,
		Y:      sy,
		Width:  ex - sx,
		Height: ey - sy,
	}
}

// 判断这个地图是否合法
func (md *MapData) IsValidMap() bool {
	// 统计所有的非墙壁地块数量
	// 从start开始统计可到达的地块数量
	// 如果2个值相等，表示地图全部位置可到达，则可以看作地图合法

	wallnum := 0
	for _, arr := range md.Data {
		for _, v := range arr {
			if !MgrStatic.StaticGenMap.IsWall(v) {
				wallnum++
			}
		}
	}

	validnum := md.CalcValidTileNum(md.Start[0], md.Start[1])

	return wallnum == validnum
}

func (md *MapData) CalcValidTileNum(sx, sy int) int {
	lst := []int{sx, sy}

	lst = md.calcNearValidTile(sx, sy, lst)

	return len(lst) / 2
}

func (md *MapData) calcNearValidTile(sx, sy int, lst []int) []int {
	if sx > 0 {
		if !MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx-1]) {
			if goutils.IndexOfInt2Slice(lst, sx-1, sy, 0) < 0 {
				lst = append(lst, sx-1, sy)

				lst = md.calcNearValidTile(sx-1, sy, lst)
			}
		}
	}

	if sx < md.GetWidth()-1 {
		if !MgrStatic.StaticGenMap.IsWall(md.Data[sy][sx+1]) {
			if goutils.IndexOfInt2Slice(lst, sx+1, sy, 0) < 0 {
				lst = append(lst, sx+1, sy)

				lst = md.calcNearValidTile(sx+1, sy, lst)
			}
		}
	}

	if sy > 0 {
		if !MgrStatic.StaticGenMap.IsWall(md.Data[sy-1][sx]) {
			if goutils.IndexOfInt2Slice(lst, sx, sy-1, 0) < 0 {
				lst = append(lst, sx, sy-1)

				lst = md.calcNearValidTile(sx, sy-1, lst)
			}
		}
	}

	if sy < md.GetHeight()-1 {
		if !MgrStatic.StaticGenMap.IsWall(md.Data[sy+1][sx]) {
			if goutils.IndexOfInt2Slice(lst, sx, sy+1, 0) < 0 {
				lst = append(lst, sx, sy+1)

				lst = md.calcNearValidTile(sx, sy+1, lst)
			}
		}
	}

	return lst
}
