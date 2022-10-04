package battle5

var MapPropTypeStr map[PropType]string

func init() {
	MapPropTypeStr = make(map[PropType]string)

	MapPropTypeStr[PropTypeHP] = "血量"
	MapPropTypeStr[PropTypeAtk] = "攻"
	MapPropTypeStr[PropTypeDef] = "防"
	MapPropTypeStr[PropTypeMagic] = "法"
	MapPropTypeStr[PropTypeSpeed] = "速"

	MapPropTypeStr[PropTypeMovDistance] = "移动距离"
	MapPropTypeStr[PropTypeAtkDistance] = "攻击距离"
	MapPropTypeStr[PropTypePlace] = "站位"

	MapPropTypeStr[PropTypeAttackType] = "攻击类型>物理|魔法"

	MapPropTypeStr[PropTypeMaxHP] = "当前最大血量"
	MapPropTypeStr[PropTypeCurHP] = "当前血量"
	MapPropTypeStr[PropTypeCurAtk] = "当前攻"
	MapPropTypeStr[PropTypeCurDef] = "当前防"
	MapPropTypeStr[PropTypeCurMagic] = "当前法"
	MapPropTypeStr[PropTypeCurSpeed] = "当前速"

	MapPropTypeStr[PropTypeCurMovDistance] = "当前移动距离"
	MapPropTypeStr[PropTypeCurAtkDistance] = "当前攻击距离"

	MapPropTypeStr[PropTypeTeamSpeedVal] = "队伍速度"
}
