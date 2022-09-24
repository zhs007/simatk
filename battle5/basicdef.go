package battle5

type PropType int

const (
	// 这一组是基本属性
	PropTypeHP    PropType = 1 // 初始HP
	PropTypeAtk   PropType = 2 // 初始atk
	PropTypeDef   PropType = 3 // 初始def
	PropTypeMagic PropType = 4 // 初始magic
	PropTypeSpeed PropType = 5 // 初始speed

	// 这一组是实际参与战斗的基本属性，一般是经过培养系统以后的值
	PropTypeMaxHP    PropType = 100 // 最大HP
	PropTypeCurHP    PropType = 101 // 当前HP
	PropTypeCurAtk   PropType = 102 // 当前atk
	PropTypeCurDef   PropType = 103 // 当前def
	PropTypeCurMagic PropType = 104 // 当前magic
	PropTypeCurSpeed PropType = 105 // 当前speed
)
