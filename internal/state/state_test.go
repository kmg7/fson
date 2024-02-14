package state_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/kmg7/fson/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	s := &state.AppState{
		Status:      "test",
		Data:        "test",
		Meta:        "test",
		Err:         errors.New("Err field must be omitted"),
		ErrInternal: true,
	}
	gjson, err := s.ToJSON()

	t.Run("JsonMarshalling", func(t *testing.T) {
		assert.Nil(t, err)
		assert.NotNil(t, gjson)
	})

	gs := state.AppState{}

	t.Run("JsonUnmarshalling", func(t *testing.T) {
		err = json.Unmarshal(gjson, &gs)
		assert.Nil(t, err)
	})

	t.Run("PrivateFields", func(t *testing.T) {
		assert.False(t, gs.ErrInternal)
		assert.Equal(t, s.Status, gs.Status)
		assert.Equal(t, s.Data, gs.Data)
		assert.Equal(t, s.Meta, gs.Meta)

	})
}
