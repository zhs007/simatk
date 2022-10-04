package battle5

const (
	MaxTurn = 10 // 最大回合数为10

	SceneWidth  = 7 // 宽度
	SceneHeight = 3 // 高度
)

type HeroID int
type SkillID int
type PropType int

const (
	// 这一组是基本属性
	PropTypeHP    PropType = 1 // 初始HP
	PropTypeAtk   PropType = 2 // 初始atk
	PropTypeDef   PropType = 3 // 初始def
	PropTypeMagic PropType = 4 // 初始magic
	PropTypeSpeed PropType = 5 // 初始speed

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
)

// 关于位置
// (x3,y1) (x2,y1) (x1,y1)
// (x3,y2) (x2,y2) (x1,y2)
// (x3,y3) (x2,y3) (x1,y3)
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}
