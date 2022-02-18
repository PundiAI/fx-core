package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewERC20Data(t *testing.T) {
	data := NewFIP20Data("test", "ERC20", uint8(18))
	exp := FIP20Data(FIP20Data{Name: "test", Symbol: "ERC20", Decimals: 0x12})
	require.Equal(t, exp, data)
}
