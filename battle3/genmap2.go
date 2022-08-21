package battle3

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// 构建房间方式来生成地图
func GenMap2(params *GenMapParams) (*MapData, error) {
	md := &MapData{}

	for y := 0; y < params.Height; y++ {
		arr := []int{}
		for x := 0; x < params.Width; x++ {
			if params.IsWall(x, y) {
				arr = append(arr, MgrStatic.StaticGenMap.GenWall())
			} else {
				arr = append(arr, MgrStatic.StaticGenMap.GenFloor())
			}
		}

		md.Data = append(md.Data, arr)
	}

	err := md.initStartExit(params)
	if err != nil {
		goutils.Error("NewMap:initStartExit",
			zap.Error(err))

		return nil, err
	}

	md.Data[md.Start[1]][md.Start[0]] = MgrStatic.StaticGenMap.GenStart()
	md.Data[md.Exit[1]][md.Exit[0]] = MgrStatic.StaticGenMap.GenExit()

	lst := []int{}
	for ri := range params.Rooms {
		rd := params.Rooms[ri]
		for i := 0; i < rd.Num; i++ {
			lst = append(lst, rd.Width, rd.Height)
		}
	}

	nmd := md.GenRooms(lst)
	if nmd == nil {
		goutils.Error("NewMap:GenRooms",
			zap.Error(ErrCannotGenMap))

		return nil, ErrCannotGenMap
	}

	nmd.Data[md.Start[1]][md.Start[0]] = MgrStatic.StaticGenMap.GenStart()
	nmd.Data[md.Exit[1]][md.Exit[0]] = MgrStatic.StaticGenMap.GenExit()

	return nmd, nil
}
