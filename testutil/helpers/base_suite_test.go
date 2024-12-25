package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

func TestBaseSuite(t *testing.T) {
	s := new(helpers.BaseSuite)
	suite.Run(t, s)
	s.SetupTest()
	assert.Equal(t, int64(1), s.Ctx.BlockHeight())
	s.Commit()
	assert.Equal(t, int64(2), s.Ctx.BlockHeight())
}
