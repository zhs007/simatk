package battle2

// 只考虑HP和DPS的情况，所以只有2种倾向，一种是肉，一种是输出

const (
	PropTypeHP    = 1 // 初始HP
	PropTypeDPS   = 2 // 初始dps
	PropTypeSpeed = 3 // 初始speed

	PropTypeCurHP    = 100 // 当前HP
	PropTypeCurDPS   = 101 // 当前dps
	PropTypeCurSpeed = 102 // 当前speed
)

const (
	UnitTypeUnknow     = 0
	UnitTypeHP         = 1
	UnitTypeDPS        = 2
	UnitTypeSpeed      = 3
	UnitTypeHPDPS      = 4
	UnitTypeHPSpeed    = 5
	UnitTypeDPSHP      = 6
	UnitTypeDPSSpeed   = 7
	UnitTypeHPSpeedHP  = 8
	UnitTypeHPSpeedDPS = 9
)

type BattleResult struct {
	Units           []*Unit `json:"units"`           // 战斗单位，2个
	WinIndex        int     `json:"winIndex"`        // 胜利索引
	FirstIndex      int     `json:"firstIndex"`      // 先手索引
	ForceFirstIndex int     `json:"forceFirstIndex"` // 强制先手索引，只在速度相等时才有用，如果无意义为-1
}
