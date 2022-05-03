package battle

// 只考虑HP和DPS的情况，所以只有2种倾向，一种是肉，一种是输出

const (
	PropTypeHP  = 1 // 初始HP
	PropTypeDPS = 2 // 初始dps

	PropTypeCurHP  = 100 // 当前HP
	PropTypeCurDPS = 101 // 当前dps
)

const (
	UnitTypeUnknow = 0
	UnitTypeHP     = 1
	UnitTypeDPS    = 2
)

type BattleResult struct {
	Units      []*Unit // 战斗单位，2个
	WinIndex   int     // 胜利索引
	FirstIndex int     // 先手索引
}
