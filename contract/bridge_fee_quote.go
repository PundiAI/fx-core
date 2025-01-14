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

func NewBridgeFeeQuoteKeeper(caller Caller, contract ...string) BridgeFeeQuoteKeeper {
	if len(contract) == 0 {
		contract = append(contract, BridgeFeeAddress)
	}
	return BridgeFeeQuoteKeeper{
		Caller: caller,
		abi:    GetBridgeFeeQuote().ABI,
		// evm module address
		from:     common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes()),
		contract: common.HexToAddress(contract[0]),
	}
}

func (k BridgeFeeQuoteKeeper) GetQuoteNonce(ctx context.Context) (*big.Int, error) {
	var res struct{ QuoteNonce *big.Int }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "quoteNonce", &res); err != nil {
		return nil, err
	}
	return res.QuoteNonce, nil
}

func (k BridgeFeeQuoteKeeper) GetChainNames(ctx context.Context) ([]common.Hash, error) {
	var res struct{ ChainNames []common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getChainNames", &res); err != nil {
		return nil, err
	}
	return res.ChainNames, nil
}

func (k BridgeFeeQuoteKeeper) GetTokens(ctx context.Context, chainName common.Hash) ([]common.Hash, error) {
	var res struct{ Tokens []common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getTokens", &res, chainName); err != nil {
		return nil, err
	}
	return res.Tokens, nil
}

func (k BridgeFeeQuoteKeeper) GetDefaultOracleQuote(ctx context.Context, chainName, tokenName common.Hash) ([]IBridgeFeeQuoteQuoteInfo, error) {
	var res struct{ Quotes []IBridgeFeeQuoteQuoteInfo }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getDefaultOracleQuote", &res, chainName, tokenName); err != nil {
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

func (k BridgeFeeQuoteKeeper) Initialize(ctx context.Context, oracle common.Address, maxQuoteCap uint8) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "initialize", oracle, maxQuoteCap)
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

func (k BridgeFeeQuoteKeeper) RegisterChain(ctx context.Context, chainName common.Hash, tokenNames ...common.Hash) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "registerChain", chainName, tokenNames)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "registerChain", res)
}

func (k BridgeFeeQuoteKeeper) AddToken(ctx context.Context, chainName common.Hash, tokenNames []common.Hash) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "addToken", chainName, tokenNames)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "addToken", res)
}
