package battle3

const (
	MaxTurns     = 10 // 最大回合数
	LeechVal     = 5  // 吸血率
	FightBackVal = 15 // 反击率
)

// 只考虑HP和DPS的情况，所以只有2种倾向，一种是肉，一种是输出

const (
	// 这一组是基本属性
	PropTypeHP  = 1 // 初始HP
	PropTypeDPS = 2 // 初始dps

	// 这一组是实际参与战斗的基本属性，一般是经过培养系统以后的值
	PropTypeMaxHP  = 100 // 最大HP
	PropTypeCurHP  = 101 // 当前HP
	PropTypeCurDPS = 102 // 当前dps

	// 这一组是状态，0/1
	PropTypeIsFirst        = 1000 // 是否先手
	PropTypeIsDouble       = 1001 // 是否连击
	PropTypeIsFightBack    = 1002 // 是否反击
	PropTypeIsLeech        = 1003 // 是否吸血
	PropTypeIsReturnDamage = 1004 // 是否反伤

	// 这一组是战斗里需要用的其它属性
	PropTypeReturnDamageVal = 2000 // 反伤率
	PropTypeLeechVal        = 2001 // 吸血率
	PropTypeUpAtk           = 2002 // 加攻
	PropTypeDownAtk         = 2003 // 减攻
	PropTypeUpDamage        = 2004 // 加伤
	PropTypeDownDamage      = 2005 // 减伤
)

const (
	UnitTypeUnknow  = 0
	UnitTypeMoreHP  = 1
	UnitTypeHP      = 2
	UnitTypeNormal  = 3
	UnitTypeDPS     = 4
	UnitTypeMoreDPS = 5
)

type BattleResult struct {
	Units           []*Unit `json:"units"`           // 战斗单位，2个
	WinIndex        int     `json:"winIndex"`        // 胜利索引
	FirstIndex      int     `json:"firstIndex"`      // 先手索引
	ForceFirstIndex int     `json:"forceFirstIndex"` // 强制先手索引，只在速度相等时才有用，如果无意义为-1
}
