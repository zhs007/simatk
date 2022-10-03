package battle5

type BattleLogNodeID int
type BattleLogNodeType int

const (
	BLNTBattleStart BattleLogNodeType = 1 // 战斗开始，根节点，不需要额外的数据
	BLNTHeroComeIn  BattleLogNodeType = 2 // 角色进场
	BLNTTurnStart   BattleLogNodeType = 3 // 回合开始
	BLNTTurnEnd     BattleLogNodeType = 4 // 回合结束
	BLNTHeroMove    BattleLogNodeType = 5 // 角色主动移动
	BLNTChgProp     BattleLogNodeType = 6 // 改变属性
	BLNTBattleReady BattleLogNodeType = 7 // 战斗开始后的准备阶段，也不需要额外的数据，处理角色上场及角色特性技能、被动技能释放、装备等
)

type BattleLogNode struct {
	NodeID        BattleLogNodeID   `json:"nodeid"`
	Type          BattleLogNodeType `json:"type"`
	Turn          int               `json:"turn,omitempty"` // turn从1开始，turnindex从0开始
	SrcHeroID     HeroID            `json:"srcheroid,omitempty"`
	SrcRealHeroID int               `json:"srcrealheroid,omitempty"`
	SrcPos        *Pos              `json:"srcpos,omitempty"`
	TargetPos     *Pos              `json:"targetpos,omitempty"`
	Children      []*BattleLogNode  `json:"children,omitempty"`
	Props         []PropType        `json:"props,omitempty"`
	PropVals      []int             `json:"propvals,omitempty"`
}

func (bln *BattleLogNode) SetSrcPos(hero *Hero) {
	bln.SrcPos = &Pos{
		X: hero.X,
		Y: hero.Y,
	}
}

func (bln *BattleLogNode) SetTargetPos(hero *Hero) {
	bln.TargetPos = &Pos{
		X: hero.X,
		Y: hero.Y,
	}
}

type BattleLog struct {
	Root      *BattleLogNode  `json:"root"`
	HashCode  string          `json:"hashcode"`
	curNodeID BattleLogNodeID `json:"-"`
}

func (bl *BattleLog) GenNodeID() BattleLogNodeID {
	nid := bl.curNodeID

	bl.curNodeID++

	return nid
}

func (bl *BattleLog) StartBattle() *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTBattleStart,
	}

	bl.Root = node

	return node
}

func (bl *BattleLog) BattleReady(parent *BattleLogNode) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTBattleReady,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) HeroComeIn(parent *BattleLogNode, hero *Hero) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTHeroComeIn,
	}

	for proptype, v := range hero.Props {
		node.Props = append(node.Props, proptype)
		node.PropVals = append(node.PropVals, v)
	}

	node.SetTargetPos(hero)

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) StartTurn(parent *BattleLogNode, turn int) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTTurnStart,
		Turn:   turn,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func NewBattleLog() *BattleLog {
	return &BattleLog{
		curNodeID: 1,
	}
}
