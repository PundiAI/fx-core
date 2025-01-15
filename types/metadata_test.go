package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetadata(t *testing.T) {
	metadata := NewMetadata("ABC Token", "ABC", 18)
	assert.NoError(t, metadata.Validate())

	assert.NoError(t, NewDefaultMetadata().Validate())
}
