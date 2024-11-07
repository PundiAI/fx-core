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

type BridgeFeeQuoteKeeper struct {
	Caller
	abi      abi.ABI
	from     common.Address
	contract common.Address
}

func NewBridgeFeeQuoteKeeper(caller Caller, contract string) BridgeFeeQuoteKeeper {
	return BridgeFeeQuoteKeeper{
		Caller: caller,
		abi:    GetBridgeFeeQuote().ABI,
		// evm module address
		from:     common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes()),
		contract: common.HexToAddress(contract),
	}
}

func (k BridgeFeeQuoteKeeper) GetQuotesByToken(ctx context.Context, chainName, tokenName string) ([]IBridgeFeeQuoteQuoteInfo, error) {
	var res struct{ Quotes []IBridgeFeeQuoteQuoteInfo }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getQuotesByToken", &res, chainName, tokenName); err != nil {
		return nil, err
	}
	return res.Quotes, nil
}

func (k BridgeFeeQuoteKeeper) GetQuoteById(ctx context.Context, id *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	var res struct{ Quote IBridgeFeeQuoteQuoteInfo }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getQuoteById", &res, id); err != nil {
		return IBridgeFeeQuoteQuoteInfo{}, err
	}
	return res.Quote, nil
}

func (k BridgeFeeQuoteKeeper) Initialize(ctx context.Context, oracle common.Address, maxQuoteIndex *big.Int) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "initialize", oracle, maxQuoteIndex)
}

func (k BridgeFeeQuoteKeeper) GetOwnerRole(ctx context.Context) (common.Hash, error) {
	var res struct{ Role common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "OWNER_ROLE", &res); err != nil {
		return common.Hash{}, err
	}
	return res.Role, nil
}

func (k BridgeFeeQuoteKeeper) GetUpgradeRole(ctx context.Context) (common.Hash, error) {
	var res struct{ Role common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "UPGRADE_ROLE", &res); err != nil {
		return common.Hash{}, err
	}
	return res.Role, nil
}

func (k BridgeFeeQuoteKeeper) GrantRole(ctx context.Context, role common.Hash, account common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "grantRole", role, account)
}

func (k BridgeFeeQuoteKeeper) RegisterChain(ctx context.Context, chainName string, tokenNames ...string) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "registerChain", chainName, tokenNames)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "registerChain", res)
}

func (k BridgeFeeQuoteKeeper) RegisterTokenName(ctx context.Context, chainName string, tokenNames []string) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "registerTokenName", chainName, tokenNames)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "registerTokenName", res)
}
