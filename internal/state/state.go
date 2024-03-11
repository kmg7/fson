package state

import (
	"encoding/json"
)

// Application state is being used for state transfers.
// It is used for communication between app modules and
// Outer clients.
type AppState struct {
	Status      string `json:"status"`         // Explains the state in short.
	Data        any    `json:"data,omitempty"` // Data could be anything.
	Meta        any    `json:"meta,omitempty"` // Additional explanations of the state.
	Err         error  `json:"-"`              // Omitted from outer application
	ErrInternal bool   `json:"-"`              // Omitted from outer application
}

func (s *AppState) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

type Error interface {
	Meta() map[string][]string
	Error() string
}
