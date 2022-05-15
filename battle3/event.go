package battle3

type GenEventData struct {
	ID            int    `yaml:"id"`
	PreFunc       string `yaml:"preFunc"`
	PreFuncParams []int  `yaml:"preFuncParams"`
	IsEnding      bool   `yaml:"isEnding"`
	TotalNum      int    `yaml:"totalNum"`
}

type Event struct {
	ID       int      `yaml:"id"`
	Children []*Event `yaml:"children"`
}
