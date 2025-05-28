package similarity

import "errors"

var (
	ErrIncompatibleVector = errors.New("incompatible vector dimensions: row count differs")
	ErrIncompatibleRow    = errors.New("incompatible vector row: column count differs at some vector row")
)
