package battle5

import "errors"

var (
	// ErrInvalidFuncName - invalid func name
	ErrInvalidFuncName = errors.New("invalid func name")
	// ErrInvalidBuffID - invalid buffid
	ErrInvalidBuffID = errors.New("invalid buffid")
	// ErrInvalidFuncValsLength - invalid funcvals length
	ErrInvalidFuncValsLength = errors.New("invalid funcvals length")
	// ErrHeroIsNull - null hero
	ErrHeroIsNull = errors.New("null hero")
	// ErrDuplicateHero - duplicate hero
	ErrDuplicateHero = errors.New("duplicate hero")
	// ErrInvalidBuffEffectString - invalid Buff Effect string
	ErrInvalidBuffEffectString = errors.New("invalid Buff Effect string")
)
