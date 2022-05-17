package battle3

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
	ID       int      `yaml:"id"`
	Children []*Event `yaml:"children"`
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
