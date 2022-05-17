package battle3

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type Unit struct {
	Props        map[int]int  `json:"props"`
	UnitType     int          `json:"unitType"`
	LstEquipment []*Equipment `json:"-"`
}

func NewUnit(hp int, dps int) *Unit {
	unit := &Unit{
		Props: make(map[int]int),
	}

	unit.Props[PropTypeHP] = hp
	unit.Props[PropTypeDPS] = dps

	unit.Props[PropTypeMaxHP] = hp
	unit.Props[PropTypeCurHP] = hp
	unit.Props[PropTypeCurDPS] = dps

	unit.UnitType = unit.analyzeUnitType()

	return unit
}

func (unit *Unit) analyzeUnitType() int {
	total := unit.Props[PropTypeHP] + unit.Props[PropTypeDPS]
	for ut := UnitTypeMoreHP; ut <= UnitTypeMoreDPS; ut++ {
		minhp, maxhp := GetHPAreaForUnitType(ut, total)
		if unit.Props[PropTypeHP] >= minhp && unit.Props[PropTypeHP] < maxhp {
			return ut
		}
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

// 是否活着
func (unit *Unit) IsAlive() bool {
	return unit.Props[PropTypeCurHP] > 0
}

// 受到伤害
func (unit *Unit) BeHurt(damage int) (bool, int) {
	if unit.Props[PropTypeCurHP]-damage < 0 {
		damage = unit.Props[PropTypeCurHP]
	}

	unit.Props[PropTypeCurHP] -= damage

	return unit.Props[PropTypeCurHP] <= 0, damage
}

// 回血
func (unit *Unit) RestoreHP(hp int) int {
	if unit.Props[PropTypeCurHP]+hp > unit.Props[PropTypeMaxHP] {
		hp = unit.Props[PropTypeMaxHP] - unit.Props[PropTypeCurHP]
	}

	unit.Props[PropTypeCurHP] += hp

	return hp
}

// 造成伤害后事件
func (unit *Unit) onDamage(damage int) bool {
	if unit.Props[PropTypeIsLeech] == 1 {
		hp := damage * unit.Props[PropTypeLeechVal] / 100
		if hp > 0 {
			unit.RestoreHP(hp)
		}
	}

	return false
}

// 被攻击后事件
func (unit *Unit) onBeAttacked(from *Unit, damage int) bool {
	if unit.Props[PropTypeIsReturnDamage] == 1 {
		// 暂时没有得到伤害后触发事件
		iskilled, _ := from.BeHurt(damage * unit.Props[PropTypeReturnDamageVal] / 100)
		if iskilled {
			return true
		}
	}

	// 反击
	if unit.Props[PropTypeIsFightBack] == 1 {
		return unit.Attack(from, false)
	}

	return false
}

// 攻击后事件
func (unit *Unit) onAttackEnd(target *Unit, isFirstAttack bool) bool {
	if isFirstAttack {
		if unit.Props[PropTypeIsDouble] == 1 {
			return unit.Attack(target, false)
		}
	}

	return false
}

// 攻击
func (unit *Unit) Attack(target *Unit, isFirstAttack bool) bool {
	damage := (unit.Props[PropTypeCurDPS] + unit.Props[PropTypeUpAtk] - target.Props[PropTypeDownAtk]) * (100 + unit.Props[PropTypeUpDamage] - target.Props[PropTypeDownDamage]) / 100
	if damage <= 0 {
		damage = 1
	}

	// 这里伤害可能溢出，所以要处理实际伤害
	iskilled, realDamage := target.BeHurt(damage)

	// 这里需要注意，就算造成击杀，也要处理造成伤害，这样吸血才不会少算一次
	// 而且，这里应该拿实际伤害来计算
	// 处理造成伤害，吸血在这里处理
	// 这里现在其实不可能造成额外的游戏结束，但还是独立检查一下，省得后面忘记了
	iskilled2 := unit.onDamage(realDamage)

	if iskilled || iskilled2 {
		return true
	}

	// 处理受击事件
	iskilled = target.onBeAttacked(unit, damage)
	if iskilled {
		return true
	}

	// 自己攻击结束的处理
	return unit.onAttackEnd(target, isFirstAttack)
}

func (unit *Unit) Reset() {
	unit.Props[PropTypeCurHP] = unit.Props[PropTypeHP]
	unit.Props[PropTypeCurDPS] = unit.Props[PropTypeDPS]
}

func (unit *Unit) ResetAndClone() *Unit {
	return NewUnit(unit.Props[PropTypeHP], unit.Props[PropTypeDPS])
}

func (unit *Unit) ChgProp(prop int, val int) (int, error) {
	return MgrStatic.MgrPropFunc.ChgProp(unit, prop, val)
}

func (unit *Unit) UseItem(id int) error {
	itemdata, err := MgrStatic.MgrItem.GetItemData(id)
	if err != nil {
		goutils.Error("Unit.UseItem:GetItemData",
			zap.Int("id", id),
			zap.Error(err))

		return err
	}

	err = MgrStatic.MgrPropFunc.Run(itemdata.ValFunc, unit, nil, PropFuncStateOn, itemdata.TargetProp, itemdata.Val, itemdata.StrVal)
	if err != nil {
		goutils.Error("Unit.UseItem:Run",
			zap.Int("id", id),
			zap.Error(err))

		return err
	}

	return nil
}

func (unit *Unit) HasEquipment(id int) bool {
	for _, v := range unit.LstEquipment {
		if v.Data.ID == id {
			return true
		}
	}

	return false
}

func (unit *Unit) TakeOff(id int) error {
	for i, v := range unit.LstEquipment {
		if v.Data.ID == id {
			err := MgrStatic.MgrPropFunc.Run(v.Data.ValFunc, unit, v.AddOn, PropFuncStateOff, v.Data.TargetProp, v.Data.Val, v.Data.StrVal)
			if err != nil {
				goutils.Error("Unit.TakeOff:Run",
					zap.Int("id", id),
					zap.Error(err))

				return err
			}

			unit.LstEquipment = append(unit.LstEquipment[:i], unit.LstEquipment[i+1:]...)

			return nil
		}
	}

	goutils.Error("Unit.TakeOff",
		zap.Int("id", id),
		zap.Error(ErrCantEquip))

	return ErrCantEquip
}

func (unit *Unit) Equip(id int) error {
	if unit.HasEquipment(id) {
		goutils.Error("Unit.Equip:HasEquipment",
			zap.Int("id", id),
			zap.Error(ErrEquiped))

		return ErrEquiped
	}

	itemdata, err := MgrStatic.MgrItem.GetItemData(id)
	if err != nil {
		goutils.Error("Unit.Equip:GetItemData",
			zap.Int("id", id),
			zap.Error(err))

		return err
	}

	equ := &Equipment{
		Data:  itemdata,
		AddOn: &AddOn{},
	}

	err = MgrStatic.MgrPropFunc.Run(itemdata.ValFunc, unit, equ.AddOn, PropFuncStateOn, itemdata.TargetProp, itemdata.Val, itemdata.StrVal)
	if err != nil {
		goutils.Error("Unit.Equip:Run",
			zap.Int("id", id),
			zap.Error(err))

		return err
	}

	unit.LstEquipment = append(unit.LstEquipment, equ)

	return nil
}
