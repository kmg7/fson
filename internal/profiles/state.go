package profiles

import "github.com/kmg7/fson/internal/state"

type State = state.AppState

const (
	StatusAlreadyExists = "already-exists"
	StatusNotFound      = "not-found"
	StatusInternalErr   = "internal-error"
)

func stateAlreadyExists(thing string) *State {
	return &state.AppState{
		Status: StatusAlreadyExists,
		Meta:   thing,
	}
}

func stateNotFound() *State {
	return &state.AppState{
		Status: StatusNotFound,
	}
}

func stateInternalErr(err error) *State {
	return &state.AppState{
		Status:      StatusInternalErr,
		Err:         err,
		ErrInternal: true,
	}
}
