package battle5

import "errors"

var (
	// ErrInvalidFuncName - invalid func name
	ErrInvalidFuncName = errors.New("invalid func name")
	// ErrInvalidFuncValsLength - invalid funcvals length
	ErrInvalidFuncValsLength = errors.New("invalid funcvals length")
)
