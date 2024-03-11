package state

import "github.com/kmg7/fson/internal/state/codes"

type ErrInternal struct {
	While string
	Err   error
}

func (e *ErrInternal) Error() string {
	return e.Err.Error()
}

func (e *ErrInternal) Meta() map[string][]string {
	return map[string][]string{
		codes.ErrInternal: {e.While, e.Err.Error()},
	}
}
