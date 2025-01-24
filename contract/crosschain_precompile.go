package contract

import (
	"context"
	"math/big"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type CrosschainPrecompileKeeper struct {
	Caller
	abi          abi.ABI
	contractAddr common.Address
}

func NewCrosschainPrecompileKeeper(caller Caller, contractAddr common.Address) CrosschainPrecompileKeeper {
	if IsZeroEthAddress(contractAddr) {
		contractAddr = common.HexToAddress(CrosschainAddress)
	}
	return CrosschainPrecompileKeeper{
		Caller:       caller,
		abi:          MustABIJson(ICrosschainMetaData.ABI),
		contractAddr: contractAddr,
	}
}

func (k CrosschainPrecompileKeeper) BridgeCoinAmount(ctx context.Context, args BridgeCoinAmountArgs) (*big.Int, error) {
	res := struct{ Amount *big.Int }{}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "bridgeCoinAmount", &res, args.Token, args.Target)
	if err != nil {
		return nil, err
	}
	return res.Amount, nil
}

func (k CrosschainPrecompileKeeper) HasOracle(ctx context.Context, args HasOracleArgs) (bool, error) {
	res := struct{ Has bool }{}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "hasOracle", &res, args.Chain, args.ExternalAddress)
	if err != nil {
		return false, err
	}
	return res.Has, nil
}

func (k CrosschainPrecompileKeeper) IsOracleOnline(ctx context.Context, args IsOracleOnlineArgs) (bool, error) {
	res := struct{ Online bool }{}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "isOracleOnline", &res, args.Chain, args.ExternalAddress)
	if err != nil {
		return false, err
	}
	return res.Online, nil
}

func (k CrosschainPrecompileKeeper) GetERC20Token(ctx context.Context, args GetERC20TokenArgs) (common.Address, bool, error) {
	res := struct {
		Token  common.Address
		Enable bool
	}{}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "getERC20Token", &res, args.Denom)
	if err != nil {
		return common.Address{}, false, err
	}
	return res.Token, res.Enable, nil
}

func (k CrosschainPrecompileKeeper) BridgeCall(ctx context.Context, value *big.Int, from common.Address, args BridgeCallArgs) (*evmtypes.MsgEthereumTxResponse, *big.Int, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, value, k.abi, "bridgeCall",
		args.DstChain, args.Refund, args.Tokens, args.Amounts, args.To, args.Data, args.QuoteId, args.GasLimit, args.Memo)
	if err != nil {
		return nil, nil, err
	}
	ret := struct{ EventNonce *big.Int }{}
	if err = k.abi.UnpackIntoInterface(&ret, "bridgeCall", res.Ret); err != nil {
		return res, nil, sdkerrors.ErrInvalidType.Wrapf("failed to unpack bridgeCall: %s", err.Error())
	}
	return res, ret.EventNonce, nil
}

func (k CrosschainPrecompileKeeper) ExecuteClaim(ctx context.Context, from common.Address, args ExecuteClaimArgs) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "executeClaim",
		args.Chain, args.EventNonce)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "executeClaim", res)
}

func (k CrosschainPrecompileKeeper) Crosschain(ctx context.Context, value *big.Int, from common.Address, args CrosschainArgs) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, value, k.abi, "crossChain",
		args.Token, args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "crossChain", res)
}
