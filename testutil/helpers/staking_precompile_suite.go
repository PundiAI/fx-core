package helpers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
)

type StakingPrecompileSuite struct {
	*ContractBaseSuite
	err error

	contract.StakingPrecompileKeeper
}

func NewStakingPrecompileSuite(require *require.Assertions, signer *Signer, caller contract.Caller, contractAddr common.Address) StakingPrecompileSuite {
	contractBaseSuite := NewContractBaseSuite(require, signer).WithContract(contractAddr)
	return StakingPrecompileSuite{
		ContractBaseSuite:       contractBaseSuite,
		StakingPrecompileKeeper: contract.NewStakingPrecompileKeeper(caller, contractAddr),
	}
}

func (s StakingPrecompileSuite) WithError(err error) StakingPrecompileSuite {
	stakingPrecompileKeeper := s
	stakingPrecompileKeeper.err = err
	return stakingPrecompileKeeper
}

func (s StakingPrecompileSuite) Error(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s StakingPrecompileSuite) AllowanceShares(ctx context.Context, args contract.AllowanceSharesArgs) *big.Int {
	shares, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).AllowanceShares(ctx, args)
	s.Error(err)
	return shares
}

func (s StakingPrecompileSuite) Delegation(ctx context.Context, args contract.DelegationArgs) (*big.Int, *big.Int) {
	amount, shares, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).Delegation(ctx, args)
	s.Error(err)
	return amount, shares
}

func (s StakingPrecompileSuite) DelegationRewards(ctx context.Context, args contract.DelegationRewardsArgs) *big.Int {
	rewards, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).DelegationRewards(ctx, args)
	s.Error(err)
	return rewards
}

func (s StakingPrecompileSuite) ValidatorList(ctx context.Context, args contract.ValidatorListArgs) []string {
	rewards, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).ValidatorList(ctx, args)
	s.Error(err)
	return rewards
}

func (s StakingPrecompileSuite) SlashingInfo(ctx context.Context, args contract.SlashingInfoArgs) (bool, *big.Int) {
	jailed, missed, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).SlashingInfo(ctx, args)
	s.Error(err)
	return jailed, missed
}

func (s StakingPrecompileSuite) ApproveShares(ctx context.Context, args contract.ApproveSharesArgs) *evmtypes.MsgEthereumTxResponse {
	res, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).ApproveShares(ctx, s.HexAddress(), args)
	s.Error(err)
	return res
}

func (s StakingPrecompileSuite) TransferShares(ctx context.Context, args contract.TransferSharesArgs) (*evmtypes.MsgEthereumTxResponse, *contract.TransferSharesRet) {
	res, ret, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).TransferShares(ctx, s.HexAddress(), args)
	s.Error(err)
	return res, ret
}

func (s StakingPrecompileSuite) TransferFromShares(ctx context.Context, args contract.TransferFromSharesArgs) (*evmtypes.MsgEthereumTxResponse, *contract.TransferFromSharesRet) {
	res, ret, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).TransferFromShares(ctx, s.HexAddress(), args)
	s.Error(err)
	return res, ret
}

func (s StakingPrecompileSuite) Withdraw(ctx context.Context, args contract.WithdrawArgs) (*evmtypes.MsgEthereumTxResponse, *big.Int) {
	res, reward, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).Withdraw(ctx, s.HexAddress(), args)
	s.Error(err)
	return res, reward
}

func (s StakingPrecompileSuite) DelegateV2(ctx context.Context, args contract.DelegateV2Args) *evmtypes.MsgEthereumTxResponse {
	res, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).DelegateV2(ctx, s.HexAddress(), args)
	s.Error(err)
	return res
}

func (s StakingPrecompileSuite) RedelegateV2(ctx context.Context, args contract.RedelegateV2Args) *evmtypes.MsgEthereumTxResponse {
	res, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).RedelegateV2(ctx, s.HexAddress(), args)
	s.Error(err)
	return res
}

func (s StakingPrecompileSuite) UndelegateV2(ctx context.Context, args contract.UndelegateV2Args) *evmtypes.MsgEthereumTxResponse {
	res, err := s.StakingPrecompileKeeper.WithContractAddr(s.contract).UndelegateV2(ctx, s.HexAddress(), args)
	s.Error(err)
	return res
}
