package battle3

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type ForEachEventFunc func(*Event) bool

type GenEventFuncData struct {
	PreFunc          string   `yaml:"preFunc"`
	PreFuncParams    []int    `yaml:"preFuncParams"`
	PreFuncStrParams []string `yaml:"preFuncStrParams"`
}

type GenEventData struct {
	ID        int                 `yaml:"id"`
	EventFunc []*GenEventFuncData `yaml:"eventFunc"`
	IsEnding  bool                `yaml:"isEnding"`
	TotalNum  int                 `yaml:"totalNum"`
}

type Event struct {
	ID         int      `yaml:"id"` // 0 是入口
	Children   []*Event `yaml:"children"`
	IsEnding   bool     `yaml:"isEnding"`
	x, y       int      `yaml:"-"`       // 坐标，主要用于输出用，目前只用于Excel输出
	Index      int      `yaml:"index"`   // 索引，单线流程的顺序索引
	StartHP    int      `yaml:"startHP"` // 初始HP
	EndHP      int      `yaml:"endHP"`   // 结束HP
	MaxHP      int      `yaml:"maxHP"`   // 最大HP
	isFinished bool     `yaml:"-"`       // 跑流程用，表示通过了，或者道具被拿走了
	TotalNum   int      `yaml:"-"`       // AI跑了多少次
	WinNum     int      `yaml:"-"`       // AI通关多少次
	Awards     []*Event `yaml:"awards"`
}

func (event *Event) Clone() *Event {
	ne := &Event{
		ID:       event.ID,
		IsEnding: event.IsEnding,
		Index:    event.Index,
	}

	for _, v := range event.Children {
		nv := v.Clone()
		ne.Children = append(ne.Children, nv)
	}

	return ne
}

func (event *Event) CloneOnlyMe() *Event {
	return &Event{
		ID:       event.ID,
		IsEnding: event.IsEnding,
		Index:    event.Index,
	}
}

// 获得可以达到的全部event
func (event *Event) BuildNextEvents() []*Event {
	if event.ID == 0 || event.isFinished {
		lst := []*Event{}

		for _, v := range event.Children {
			nlst := v.BuildNextEvents()
			if len(nlst) > 0 {
				lst = append(lst, nlst...)
			}
		}

		return lst
	}

	return []*Event{event}
}

func (event *Event) CountID(id int) int {
	num := 0

	if event.ID == id {
		num++
	}

	for _, v := range event.Children {
		num += v.CountID(id)
	}

	return num
}

func (event *Event) AddChild(e *Event) {
	ne := e.CloneOnlyMe()
	event.Children = append(event.Children, ne)

	if len(e.Awards) > 0 {
		for _, v := range e.Awards {
			ne.AddChild(v)
		}
	}
}

// 加到最深的子节点
func (event *Event) Add2Last(id int) {
	if event.ID == 0 {
		event.ID = id

		return
	}

	if len(event.Children) == 0 {
		event.Children = []*Event{
			{
				ID: id,
			},
		}

		return
	}

	event.Children[0].Add2Last(id)
}

func (event *Event) ForEach(funcForEach ForEachEventFunc) bool {
	if event.ID >= 0 {
		if !funcForEach(event) {
			return false
		}

		for _, v := range event.Children {
			if !v.ForEach(funcForEach) {
				return false
			}
		}

		return true
	}

	return true
}

func (event *Event) GetLeafNodes() []*Event {
	if len(event.Children) == 0 {
		return []*Event{event}
	}

	nodes := []*Event{}

	for _, v := range event.Children {
		lst := v.GetLeafNodes()
		if len(lst) != 0 {
			nodes = append(nodes, lst...)
		}
	}

	return nodes
}

func (event *Event) GetNotValidLeafNodes() []*Event {
	if len(event.Children) == 0 && IsLeafNode(event) {
		return nil
	}

	nodes := []*Event{event}

	for _, v := range event.Children {
		lst := v.GetNotValidLeafNodes()
		if len(lst) != 0 {
			nodes = append(nodes, lst...)
		}
	}

	return nodes
}

func (event *Event) GetValidLeafNodes() []*Event {
	if len(event.Children) == 0 && !IsLeafNode(event) {
		return []*Event{event}
	}

	nodes := []*Event{}

	for _, v := range event.Children {
		lst := v.GetValidLeafNodes()
		if len(lst) != 0 {
			nodes = append(nodes, lst...)
		}
	}

	return nodes
}

// func (event *Event) rebuildY(parenty int) int {
// 	event.y = parenty + 1
// 	maxy := event.y
// 	for _, v := range event.Children {
// 		cy := v.rebuildY(event.y)
// 		if cy > maxy {
// 			maxy = cy
// 		}
// 	}

// 	return maxy
// }

// func (event *Event) countY(y int) int {
// 	if event.y == y {
// 		return 1
// 	}

// 	cn := 0

// 	for _, v := range event.Children {
// 		cn += v.countY(y)
// 	}

// 	return cn
// }

func (event *Event) countXOff() int {
	if len(event.Children) == 0 {
		return 1
	}

	num := 0

	for _, v := range event.Children {
		num += v.countXOff()
	}

	return num
}

func (event *Event) rebuildPos(x, y int) {
	event.x = x
	event.y = y

	off := 0
	for _, v := range event.Children {
		v.rebuildPos(x+off, y+1)
		off += v.countXOff()
	}
}

func (event *Event) GetName() string {
	if event.ID == 0 {
		if event.TotalNum > 0 {
			return fmt.Sprintf("%d-root %v", event.Index, event.WinNum*100/event.TotalNum)
		}

		return fmt.Sprintf("%d-root", event.Index)
	}

	if IsItem(event.ID) || IsEquipment(event.ID) {
		data, _ := MgrStatic.MgrItem.GetItemData(event.ID)
		return fmt.Sprintf("%d-%v", event.Index, data.Name)
	} else if IsMonster(event.ID) {
		data, _ := MgrStatic.MgrCharacter.GetCharacterData(event.ID)
		return fmt.Sprintf("%d-%v", event.Index, data.Name)
	}

	return fmt.Sprintf("%d-error", event.Index)
}

func (event *Event) OutputExcel(f *excelize.File, sheet string) error {
	event.rebuildPos(0, 0)

	event.ForEach(func(e *Event) bool {
		f.SetCellStr(sheet, goutils.Pos2Cell(e.x, e.y), e.GetName())

		return true
	})

	return nil
}

func (event *Event) CountNodes() int {
	num := 0

	if event.ID > 0 {
		num++
	}

	for _, v := range event.Children {
		num += v.CountNodes()
	}

	return num
}

// 获得可以达到的全部event，增加一个外部接口来判断
func (event *Event) BuildNextEventsEx(eachEvent ForEachEventFunc) []*Event {
	if event.ID == 0 || event.isFinished {
		lst := []*Event{}

		for _, v := range event.Children {
			nlst := v.BuildNextEventsEx(eachEvent)
			if len(nlst) > 0 {
				lst = append(lst, nlst...)
			}
		}

		return lst
	}

	if eachEvent(event) {
		return []*Event{event}
	}

	return nil
}
