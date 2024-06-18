package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrPtr(t *testing.T) {

	t.Run("converts a string to a pointer", func(t *testing.T) {
		ptr := StrPtr("test")
		str := "test"
		assert.IsType(t, &str, ptr)
	})
}
