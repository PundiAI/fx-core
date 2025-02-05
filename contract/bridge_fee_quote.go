package contract

import (
	"context"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	for i := 0; i < len(res.Quotes); i++ {
		if res.Quotes[i].Id.Sign() <= 0 {
			res.Quotes = append(res.Quotes[:i], res.Quotes[i+1:]...)
			i--
		}
	}
	return res.Quotes, nil
}

func (k BridgeFeeQuoteKeeper) GetQuoteById(ctx context.Context, id *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	var res struct{ Quote IBridgeFeeQuoteQuoteInfo }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "getQuoteById", &res, id); err != nil {
		return IBridgeFeeQuoteQuoteInfo{}, err
	}
	if res.Quote.Id.Sign() <= 0 {
		return IBridgeFeeQuoteQuoteInfo{}, sdkerrors.ErrInvalidRequest.Wrapf("quote not found")
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

func (k BridgeFeeQuoteKeeper) Quote(ctx context.Context, from common.Address, inputs []IBridgeFeeQuoteQuoteInput) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contract, nil, k.abi, "quote", inputs)
	if err != nil {
		return k.unpackError(res, err)
	}
	var ret struct{ QuoteIds []*big.Int }
	if err = k.abi.UnpackIntoInterface(&ret, "quote", res.Ret); err != nil {
		return res, sdkerrors.ErrInvalidType.Wrapf("failed to unpack %s: %s", "quote", err.Error())
	}
	if len(ret.QuoteIds) == 0 {
		return res, sdkerrors.ErrInvalidRequest.Wrapf("quote ids not found")
	}
	return res, nil
}

func (k BridgeFeeQuoteKeeper) GrantRole(ctx context.Context, role common.Hash, account common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "grantRole", role, account)
}

func (k BridgeFeeQuoteKeeper) RegisterChain(ctx context.Context, chainName common.Hash, tokenNames ...common.Hash) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "registerChain", chainName, tokenNames)
	if err != nil {
		return k.unpackError(res, err)
	}
	return unpackRetIsOk(k.abi, "registerChain", res)
}

func (k BridgeFeeQuoteKeeper) AddToken(ctx context.Context, chainName common.Hash, tokenNames []common.Hash) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "addToken", chainName, tokenNames)
	if err != nil {
		return k.unpackError(res, err)
	}
	return unpackRetIsOk(k.abi, "addToken", res)
}

func (k BridgeFeeQuoteKeeper) UpgradeTo(ctx context.Context, newLogic common.Address) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "upgradeTo", newLogic)
	if err != nil {
		return k.unpackError(res, err)
	}
	return res, nil
}

func (k BridgeFeeQuoteKeeper) UpgradeToAndCall(ctx context.Context, newLogic common.Address, data []byte) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "upgradeToAndCall", newLogic, data)
	if err != nil {
		return k.unpackError(res, err)
	}
	return unpackRetIsOk(k.abi, "upgradeToAndCall", res)
}

func (k BridgeFeeQuoteKeeper) unpackError(res *types.MsgEthereumTxResponse, err error) (*types.MsgEthereumTxResponse, error) {
	if err == nil {
		return res, nil
	}
	revertInfo, unpackErr := UnpackRevertError(k.abi, res.Ret)
	if unpackErr != nil {
		return res, err
	}
	return res, fmt.Errorf("reverted: %s, vmErr: %s", revertInfo, err.Error())
}
