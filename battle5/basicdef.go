package battle5

const (
	MaxTurn = 10 // 最大回合数为10

	SceneWidth  = 7 // 宽度
	SceneHeight = 3 // 高度
)

type HeroID int
type SkillID int
type BuffID int
type PropType int

const (
	// 这一组是基本属性
	PropTypeHP    PropType = 1 // 初始HP
	PropTypeAtk   PropType = 2 // 初始atk
	PropTypeDef   PropType = 3 // 初始def
	PropTypeMagic PropType = 4 // 初始magic
	PropTypeSpeed PropType = 5 // 初始speed
	PropTypeDodge PropType = 6 // 初始闪避
	PropTypeCrit  PropType = 7 // 初始暴击

	PropTypeMovDistance PropType = 20 // 移动距离
	PropTypeAtkDistance PropType = 21 // 攻击距离
	PropTypePlace       PropType = 22 // 位置，1是前排、2是中排、3是后排

	PropTypeAttackType PropType = 50 // 默认是物理攻击还是魔法攻击，0 表示物理，1 表示魔法

	// 这一组是实际参与战斗的基本属性，一般是经过培养系统以后的值
	PropTypeMaxHP    PropType = 100 // 最大HP
	PropTypeCurHP    PropType = 101 // 当前HP
	PropTypeCurAtk   PropType = 102 // 当前atk
	PropTypeCurDef   PropType = 103 // 当前def
	PropTypeCurMagic PropType = 104 // 当前magic
	PropTypeCurSpeed PropType = 105 // 当前speed

	PropTypeCurMovDistance PropType = 120 // 当前移动距离
	PropTypeCurAtkDistance PropType = 121 // 当前攻击距离

	PropTypeTeamSpeedVal PropType = 150 // 队伍速度值，快的队里这个值大于慢的队

	// 状态值，是 或者 否
	PropTypeTaunt PropType = 200 // 嘲讽

	// 百分比附加属性
	PropTypeReduceDamage PropType = 300 // 减少伤害
)

type TargetType int

const (
	TargetTypeEnemy  TargetType = 1 // 敌人
	TargetTypeFriend TargetType = 2 // 友方
	TargetTypeAll    TargetType = 3 // 不分敌我
)

type SkillType int

const (
	SkillTypeBasicAtk SkillType = 1 // 普攻
	SkillTypeNatural  SkillType = 2 // 天赋，天赋技能
	SkillTypeUltimate SkillType = 3 // 必杀，终极技能
	SkillTypeNormal   SkillType = 4 // 普通
)

type ReleaseSkillType int

const (
	ReleaseSkillTypeNormal  ReleaseSkillType = 1 // 主动技能
	ReleaseSkillTypePassive ReleaseSkillType = 2 // 被动
)

type TriggerType int

const (
	TriggerTypeNone      TriggerType = 0 // none
	TriggerTypeAttach    TriggerType = 1 // attach buff时，buff在被附加时
	TriggerTypeTurnStart TriggerType = 2 // 回合开始
	TriggerTypeTurnEnd   TriggerType = 3 // 回合结束
	TriggerTypeFind      TriggerType = 4 // 寻找目标时
)

func Str2TriggerType(str string) TriggerType {

	switch str {
	case "find":
		return TriggerTypeFind
	case "attach":
		return TriggerTypeAttach
	case "turnstart":
		return TriggerTypeTurnStart
	case "turnend":
		return TriggerTypeTurnEnd
	}

	return TriggerTypeNone
}

// 触发时数据，这里面的数据可以完整的表达任意时点
type TriggerData struct {
	Parent   *TriggerData // 父节点，某些情况下，是A打B，B反击然后触发一个判断，所以这里需要有溯源机制
	Trigger  TriggerType  // 类型
	Turn     int          // 回合
	Src      *Hero        // 源英雄，发起方，譬如攻击、寻敌等
	Target   *Hero        // 目标英雄，受击方，被确定为目标等
	SrcSkill *Skill       // 源技能
	SrcBuff  *Buff        // 源buff
}

type BuffInstanceID int
type HeroInstanceID int

// 关于位置
// (x3,y1) (x2,y1) (x1,y1)
// (x3,y2) (x2,y2) (x1,y2)
// (x3,y3) (x2,y3) (x1,y3)
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (pos *Pos) SetXY(x, y int) {
	pos.X = x
	pos.Y = y
}

func (pos *Pos) Set(p *Pos) {
	pos.X = p.X
	pos.Y = p.Y
}

func (pos *Pos) CalcDistance(target *Pos) int {
	ox := target.X - pos.X
	oy := target.Y - pos.Y

	if ox < 0 {
		ox = -ox
	}

	if oy < 0 {
		oy = -oy
	}

	return ox + oy
}

func (pos *Pos) Equal(target *Pos) bool {
	return pos.X == target.X && pos.Y == target.Y
}

func (pos *Pos) Clone() *Pos {
	return &Pos{
		X: pos.X,
		Y: pos.Y,
	}
}

func NewPos(x, y int) *Pos {
	return &Pos{
		X: x,
		Y: y,
	}
}

type BattleActionFromData struct {
	Parent *BattleLogNode
	Hero   *Hero
	Skill  *Skill
}

func NewBattleActionFromData(parent *BattleLogNode, hero *Hero, skill *Skill) *BattleActionFromData {
	return &BattleActionFromData{
		Parent: parent,
		Hero:   hero,
		Skill:  skill,
	}
}
