package state

import (
	"fmt"

	"github.com/kmg7/fson/internal/state/codes"
)

type ErrAlreadyExists struct {
	Resource string
	Subject  string
}

func (e *ErrAlreadyExists) Error() string {
	return fmt.Sprintf("'%v' '%v' already exists", e.Resource, e.Subject)
}

func (e *ErrAlreadyExists) Meta() map[string][]string {
	return map[string][]string{
		codes.ErrAlreadyExists: {e.Resource, e.Subject},
	}
}
