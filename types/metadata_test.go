package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetadata_Validate(t *testing.T) {
	for _, metadata := range GetMetadata() {
		err := metadata.Validate()
		assert.NoError(t, err)
	}
}
