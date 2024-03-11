package state

import (
	"fmt"

	"github.com/kmg7/fson/internal/state/codes"
)

type ErrNotFound struct {
	Resource string
	Subject  string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("'%v' '%v' not found", e.Resource, e.Subject)
}

func (e *ErrNotFound) Meta() map[string][]string {
	return map[string][]string{
		codes.ErrNotFound: {e.Resource, e.Subject},
	}
}
