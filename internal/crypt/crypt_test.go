package crypt_test

import (
	"testing"

	"github.com/kmg7/fson/internal/crypt"
	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	t.Run("FineDefaultCost", func(t *testing.T) {
		c := crypt.Instance(crypt.Options{})
		assert.NotNil(t, c)
	})
}
func TestBcrypt(t *testing.T) {
	c := crypt.Instance(crypt.Options{
		BcryptCost: 8,
	})

	t.Run("FineSubject", func(t *testing.T) {
		subject := "Password1_."

		got, err := c.Bcrypt([]byte(subject))
		assert.Nil(t, err)
		assert.NotNil(t, got)
	})
}

func TestBcryptCompare(t *testing.T) {
	c := crypt.Instance(crypt.Options{
		BcryptCost: 8,
	})
	fineSubject1 := []byte("Passowrd_.1")
	fineSubject2 := []byte("Passowrd_.2")
	badSubject := []byte(";i$$")
	bs1Hash := badSubject

	fs1Hash, err := c.Bcrypt(fineSubject1)
	assert.Nil(t, err)
	assert.NotNil(t, fs1Hash)

	fs2Hash, err := c.Bcrypt(fineSubject2)
	assert.Nil(t, err)
	assert.NotNil(t, fs2Hash)

	t.Run("FineSubject", func(t *testing.T) {
		t.Run("Matches", func(t *testing.T) {
			err = c.BcryptCompare(*fs1Hash, fineSubject1)
			assert.Nil(t, err)
		})

		t.Run("NotMathces", func(t *testing.T) {
			err = c.BcryptCompare(*fs1Hash, fineSubject2)
			assert.NotNil(t, err)
		})
	})
	t.Run("BadHash", func(t *testing.T) {
		err = c.BcryptCompare(bs1Hash, fineSubject1)
		assert.NotNil(t, err)
	})

	t.Run("BadSubject", func(t *testing.T) {
		err = c.BcryptCompare(*fs1Hash, badSubject)
		assert.NotNil(t, err)
	})
}
