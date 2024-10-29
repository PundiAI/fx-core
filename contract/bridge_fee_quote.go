package contract

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type BrideFeeQuoteKeeper struct {
	Caller
	abi      abi.ABI
	from     common.Address
	contract common.Address
}

func NewBridgeFeeQuoteKeeper(caller Caller, contract string) BrideFeeQuoteKeeper {
	return BrideFeeQuoteKeeper{
		Caller: caller,
		abi:    bridgeFeeQuoteABI,
		// evm module address
		from:     common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes()),
		contract: common.HexToAddress(contract),
	}
}

func (k BrideFeeQuoteKeeper) GetQuotesByToken(ctx context.Context, chainName, tokenName string) ([]IBridgeFeeQuoteQuoteInfo, error) {
	var res struct{ Quotes []IBridgeFeeQuoteQuoteInfo }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getQuotesByToken", &res, chainName, tokenName); err != nil {
		return nil, err
	}
	return res.Quotes, nil
}

func (k BrideFeeQuoteKeeper) GetQuoteById(ctx context.Context, id *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	var res struct{ Quote IBridgeFeeQuoteQuoteInfo }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getQuoteById", &res, id); err != nil {
		return IBridgeFeeQuoteQuoteInfo{}, err
	}
	return res.Quote, nil
}
