package battle3

import "errors"

var (
	// ErrInvalidPropStr - invalid prop string
	ErrInvalidPropStr = errors.New("invalid prop string")
	// ErrInvalidPropFunc - invalid PropFunc
	ErrInvalidPropFunc = errors.New("invalid PropFunc")
	// ErrInvalidFuncPropParam - invalid FuncItem parameter
	ErrInvalidFuncPropParam = errors.New("invalid FuncItem parameter")
	// ErrInvalidBasicPropFunc - invalid BasicPropFunc
	ErrInvalidBasicPropFunc = errors.New("invalid BasicPropFunc")
	// ErrInvalidItemID - invalid ItemID
	ErrInvalidItemID = errors.New("invalid ItemID")
	// ErrEquiped - equiped
	ErrEquiped = errors.New("equiped")
	// ErrCantEquip - can not equip
	ErrCantEquip = errors.New("can not equip")
	// ErrInvalidFuncEventParam - invalid FuncEvent parameter
	ErrInvalidFuncEventParam = errors.New("invalid FuncEvent parameter")
	// ErrInvalidCharacterID - invalid CharacterID
	ErrInvalidCharacterID = errors.New("invalid CharacterID")
	// ErrInvalidEventFunc - invalid EventFunc
	ErrInvalidEventFunc = errors.New("invalid EventFunc")
	// ErrNoEvent - no event
	ErrNoEvent = errors.New("no event")
	// ErrInvalidData - invalid data
	ErrInvalidData = errors.New("invalid data")
	// ErrInvalidStageDevIndex - invalid StageDev index
	ErrInvalidStageDevIndex = errors.New("invalid StageDev index")
	// ErrBattle - battle error
	ErrBattle = errors.New("battle error")
	// ErrInvalidUnitType - invalid UnitType
	ErrInvalidUnitType = errors.New("invalid UnitType")
	// ErrInvalidStageID - invalid StageID
	ErrInvalidStageID = errors.New("invalid StageID")
	// ErrCannotGenMap - cannot gen map
	ErrCannotGenMap = errors.New("cannot gen map")
)
