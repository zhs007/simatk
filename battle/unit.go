package battle

type Unit struct {
	Props    map[int]int `json:"props"`
	UnitType int         `json:"unitType"`
}

func NewUnit(hp int, dps int) *Unit {
	unit := &Unit{
		Props: make(map[int]int),
	}

	unit.Props[PropTypeHP] = hp
	unit.Props[PropTypeDPS] = dps

	unit.Props[PropTypeCurHP] = hp
	unit.Props[PropTypeCurDPS] = dps

	unit.UnitType = unit.analyzeUnitType()

	return unit
}

func (unit *Unit) analyzeUnitType() int {
	if unit.Props[PropTypeHP] > unit.Props[PropTypeDPS] {
		return UnitTypeHP
	}

	if unit.Props[PropTypeHP] < unit.Props[PropTypeDPS] {
		return UnitTypeDPS
	}

	return UnitTypeUnknow
}

func (unit *Unit) GetLastHPPer() int {
	if unit.Props[PropTypeCurHP] <= 0 {
		return 0
	}

	p := unit.Props[PropTypeCurHP] * 100 / unit.Props[PropTypeHP]

	return p
}

func (unit *Unit) IsAlive() bool {
	return unit.Props[PropTypeCurHP] > 0
}

func (unit *Unit) Attack(target *Unit) {
	target.Props[PropTypeCurHP] -= unit.Props[PropTypeCurDPS]
}

func (unit *Unit) Reset() {
	unit.Props[PropTypeCurHP] = unit.Props[PropTypeHP]
	unit.Props[PropTypeCurDPS] = unit.Props[PropTypeDPS]
}

func (unit *Unit) ResetAndClone() *Unit {
	return NewUnit(unit.Props[PropTypeHP], unit.Props[PropTypeDPS])
}
