package battle3

const (
	MaxTurns     = 10 // 最大回合数
	LeechVal     = 5  // 吸血率
	FightBackVal = 15 // 反击率

	MinCharacterID = 1000 // 最小 CharacterID
	MaxCharacterID = 9999 // 最大 CharacterID

	MinItemID = 10000 // 最小 ItemID
	MaxItemID = 99999 // 最大 ItemID
)

type PropType int

const (
	// 这一组是基本属性
	PropTypeHP  PropType = 1 // 初始HP
	PropTypeDPS PropType = 2 // 初始dps

	// 这一组是实际参与战斗的基本属性，一般是经过培养系统以后的值
	PropTypeMaxHP  PropType = 100 // 最大HP
	PropTypeCurHP  PropType = 101 // 当前HP
	PropTypeCurDPS PropType = 102 // 当前dps

	// 这一组是状态，0/1
	PropTypeIsFirst        PropType = 1000 // 是否先手
	PropTypeIsDouble       PropType = 1001 // 是否连击
	PropTypeIsFightBack    PropType = 1002 // 是否反击
	PropTypeIsLeech        PropType = 1003 // 是否吸血
	PropTypeIsReturnDamage PropType = 1004 // 是否反伤

	// 这一组是战斗里需要用的其它属性
	PropTypeReturnDamageVal PropType = 2000 // 反伤率
	PropTypeLeechVal        PropType = 2001 // 吸血率
	PropTypeUpAtk           PropType = 2002 // 加攻
	PropTypeDownAtk         PropType = 2003 // 减攻
	PropTypeUpDamage        PropType = 2004 // 加伤
	PropTypeDownDamage      PropType = 2005 // 减伤
)

type UnitType int

const (
	UnitTypeUnknow  UnitType = 0
	UnitTypeMoreHP  UnitType = 1
	UnitTypeHP      UnitType = 2
	UnitTypeNormal  UnitType = 3
	UnitTypeDPS     UnitType = 4
	UnitTypeMoreDPS UnitType = 5
)

type BattleResult struct {
	Units           []*Unit `json:"units"`           // 战斗单位，2个
	WinIndex        int     `json:"winIndex"`        // 胜利索引
	FirstIndex      int     `json:"firstIndex"`      // 先手索引
	ForceFirstIndex int     `json:"forceFirstIndex"` // 强制先手索引，只在速度相等时才有用，如果无意义为-1
	Turns           int     `json:"turns"`           // 回合数
}

var mapProp map[string]PropType

func Str2Prop(str string) (PropType, error) {
	prop, isok := mapProp[str]
	if !isok {
		return 0, ErrInvalidPropStr
	}

	return prop, nil
}

var mapUnitType map[UnitType]string

func UnitType2Str(ut UnitType) (string, error) {
	str, isok := mapUnitType[ut]
	if !isok {
		return "", ErrInvalidUnitType
	}

	return str, nil
}

func init() {
	mapProp = make(map[string]PropType)

	mapProp["hp"] = PropTypeHP
	mapProp["dps"] = PropTypeDPS

	mapProp["maxhp"] = PropTypeMaxHP
	mapProp["curhp"] = PropTypeCurHP
	mapProp["curdps"] = PropTypeCurDPS

	mapProp["isfirst"] = PropTypeIsFirst
	mapProp["isdouble"] = PropTypeIsDouble
	mapProp["isfightback"] = PropTypeIsFightBack
	mapProp["isleech"] = PropTypeIsLeech
	mapProp["isreturndamage"] = PropTypeIsReturnDamage

	mapProp["returndamageval"] = PropTypeReturnDamageVal
	mapProp["leechval"] = PropTypeLeechVal
	mapProp["upatk"] = PropTypeUpAtk
	mapProp["downatk"] = PropTypeDownAtk
	mapProp["updamage"] = PropTypeUpDamage
	mapProp["downdamage"] = PropTypeDownDamage

	mapUnitType = make(map[UnitType]string)

	mapUnitType[UnitTypeUnknow] = "Unknow"
	mapUnitType[UnitTypeMoreHP] = "More HP"
	mapUnitType[UnitTypeHP] = "HP"
	mapUnitType[UnitTypeNormal] = "Normal"
	mapUnitType[UnitTypeDPS] = "DPS"
	mapUnitType[UnitTypeMoreDPS] = "More DPS"
}
