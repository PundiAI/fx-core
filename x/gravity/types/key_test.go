package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetFeeSecondIndexKey(t *testing.T) {
	t.Log(GetFeeSecondIndexKey(ERC20Token{Amount: sdk.NewInt(100), Contract: "0x0412C7c846bb6b7DC462CF6B453f76D8440b2609"}))
}
