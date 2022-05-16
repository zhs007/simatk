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
)
