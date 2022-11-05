package battle5funcs

import "errors"

var (
	// ErrInvalidVals - invalid vals
	ErrInvalidVals = errors.New("invalid vals")
	// ErrInvalidStrVals - invalid strvals
	ErrInvalidStrVals = errors.New("invalid strvals")
	// ErrInvalidValsOrStrVals - invalid vals or strvals
	ErrInvalidValsOrStrVals = errors.New("invalid vals or strvals")
)
