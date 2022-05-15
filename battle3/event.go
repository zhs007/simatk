package battle3

type Event struct {
	ID            int      `yaml:"id"`
	PreFunc       string   `yaml:"preFunc"`
	PreFuncParams []string `yaml:"preFuncParams"`
	IsEnding      bool     `yaml:"isEnding"`
	Children      []*Event `yaml:"children"`
}
