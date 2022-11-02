package domain

import (
	"fmt"
)

var (
	ErrInternalError = fmt.Errorf("failed with internal error")
	ErrBadRequest    = fmt.Errorf("failed with invalid request parameters")
)
