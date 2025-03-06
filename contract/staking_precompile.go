package contract

import (
	"context"
	"math/big"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type StakingPrecompileKeeper struct {
	Caller
	abi          abi.ABI
	contractAddr common.Address
}

func NewStakingPrecompileKeeper(caller Caller, contractAddr common.Address) StakingPrecompileKeeper {
	if IsZeroEthAddress(contractAddr) {
		contractAddr = common.HexToAddress(StakingAddress)
	}
	return StakingPrecompileKeeper{
		Caller:       caller,
		abi:          MustABIJson(IStakingMetaData.ABI),
		contractAddr: contractAddr,
	}
}

func (k StakingPrecompileKeeper) WithContract(addr common.Address) StakingPrecompileKeeper {
	keeper := k
	keeper.contractAddr = addr
	return keeper
}

func (k StakingPrecompileKeeper) AllowanceShares(ctx context.Context, args AllowanceSharesArgs) (*big.Int, error) {
	var output struct {
		Shares *big.Int
	}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "allowanceShares", &output, args.Validator, args.Owner, args.Spender)
	if err != nil {
		return nil, err
	}
	return output.Shares, nil
}

func (k StakingPrecompileKeeper) Delegation(ctx context.Context, args DelegationArgs) (*big.Int, *big.Int, error) {
	var output struct {
		Shares         *big.Int
		DelegateAmount *big.Int
	}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "delegation", &output, args.Validator, args.Delegator)
	if err != nil {
		return nil, nil, err
	}
	return output.Shares, output.DelegateAmount, nil
}

func (k StakingPrecompileKeeper) DelegationRewards(ctx context.Context, args DelegationRewardsArgs) (*big.Int, error) {
	var output struct {
		Rewards *big.Int
	}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "delegationRewards", &output, args.Validator, args.Delegator)
	if err != nil {
		return nil, err
	}
	return output.Rewards, nil
}

func (k StakingPrecompileKeeper) ValidatorList(ctx context.Context, args ValidatorListArgs) ([]string, error) {
	var valList []string
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "validatorList", &valList, args.SortBy)
	if err != nil {
		return nil, err
	}
	return valList, nil
}

func (k StakingPrecompileKeeper) SlashingInfo(ctx context.Context, args SlashingInfoArgs) (bool, *big.Int, error) {
	var output struct {
		Jailed bool
		Missed *big.Int
	}
	err := k.QueryContract(ctx, common.Address{}, k.contractAddr, k.abi, "slashingInfo", &output, args.Validator)
	if err != nil {
		return false, nil, err
	}
	return output.Jailed, output.Missed, nil
}

func (k StakingPrecompileKeeper) ApproveShares(ctx context.Context, from common.Address, args ApproveSharesArgs) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "approveShares", args.Validator, args.Spender, args.Shares)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "approveShares", res)
}

func (k StakingPrecompileKeeper) TransferShares(ctx context.Context, from common.Address, args TransferSharesArgs) (*evmtypes.MsgEthereumTxResponse, *TransferSharesRet, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "transferShares", args.Validator, args.To, args.Shares)
	if err != nil {
		return nil, nil, err
	}
	ret := new(TransferSharesRet)
	if err = k.abi.UnpackIntoInterface(ret, "transferShares", res.Ret); err != nil {
		return res, nil, sdkerrors.ErrInvalidType.Wrapf("failed to unpack transferShares: %s", err.Error())
	}
	return res, ret, nil
}

func (k StakingPrecompileKeeper) TransferFromShares(ctx context.Context, from common.Address, args TransferFromSharesArgs) (*evmtypes.MsgEthereumTxResponse, *TransferFromSharesRet, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "transferFromShares", args.Validator, args.From, args.To, args.Shares)
	if err != nil {
		return nil, nil, err
	}
	ret := new(TransferFromSharesRet)
	if err = k.abi.UnpackIntoInterface(ret, "transferFromShares", res.Ret); err != nil {
		return res, nil, sdkerrors.ErrInvalidType.Wrapf("failed to unpack transferFromShares: %s", err.Error())
	}
	return res, ret, nil
}

func (k StakingPrecompileKeeper) Withdraw(ctx context.Context, from common.Address, args WithdrawArgs) (*evmtypes.MsgEthereumTxResponse, *big.Int, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "withdraw", args.Validator)
	if err != nil {
		return nil, nil, err
	}
	ret := struct{ Reward *big.Int }{}
	if err = k.abi.UnpackIntoInterface(&ret, "withdraw", res.Ret); err != nil {
		return res, nil, sdkerrors.ErrInvalidType.Wrapf("failed to unpack withdraw: %s", err.Error())
	}
	return res, ret.Reward, nil
}

func (k StakingPrecompileKeeper) DelegateV2(ctx context.Context, from common.Address, value *big.Int, args DelegateV2Args) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, value, k.abi, "delegateV2", args.Validator, args.Amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "delegateV2", res)
}

func (k StakingPrecompileKeeper) RedelegateV2(ctx context.Context, from common.Address, args RedelegateV2Args) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "redelegateV2", args.ValidatorSrc, args.ValidatorDst, args.Amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "redelegateV2", res)
}

func (k StakingPrecompileKeeper) UndelegateV2(ctx context.Context, from common.Address, args UndelegateV2Args) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "undelegateV2", args.Validator, args.Amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "undelegateV2", res)
}
