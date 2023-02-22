package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	LPTokenOwnerModuleName = "lp_token"

	LPTokenSymbol   = "LP Token"
	LPTokenDecimals = uint8(36)

	MethodMintToken    = "mint"
	MethodBurnToken    = "burn"
	MethodSelfDestruct = "selfDestruct"

	LPTokenTransferEventName = "__FXLPTokenTransfer"
)

type FXLPTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}
