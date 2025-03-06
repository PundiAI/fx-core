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
	require *require.Assertions
	err     error

	keeper contract.StakingPrecompileKeeper
}

func NewStakingPrecompileSuite(require *require.Assertions, caller contract.Caller) StakingPrecompileSuite {
	address := common.HexToAddress(contract.StakingAddress)
	return StakingPrecompileSuite{
		require: require,
		keeper:  contract.NewStakingPrecompileKeeper(caller, address),
	}
}

func (s StakingPrecompileSuite) WithContract(addr common.Address) StakingPrecompileSuite {
	suite := s
	suite.keeper = suite.keeper.WithContract(addr)
	return suite
}

func (s StakingPrecompileSuite) WithError(err error) StakingPrecompileSuite {
	suite := s
	suite.err = err
	return suite
}

func (s StakingPrecompileSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s StakingPrecompileSuite) AllowanceShares(ctx context.Context, args contract.AllowanceSharesArgs) *big.Int {
	shares, err := s.keeper.AllowanceShares(ctx, args)
	s.requireError(err)
	return shares
}

func (s StakingPrecompileSuite) Delegation(ctx context.Context, args contract.DelegationArgs) (*big.Int, *big.Int) {
	amount, shares, err := s.keeper.Delegation(ctx, args)
	s.requireError(err)
	return amount, shares
}

func (s StakingPrecompileSuite) DelegationRewards(ctx context.Context, args contract.DelegationRewardsArgs) *big.Int {
	rewards, err := s.keeper.DelegationRewards(ctx, args)
	s.requireError(err)
	return rewards
}

func (s StakingPrecompileSuite) ValidatorList(ctx context.Context, args contract.ValidatorListArgs) []string {
	rewards, err := s.keeper.ValidatorList(ctx, args)
	s.requireError(err)
	return rewards
}

func (s StakingPrecompileSuite) SlashingInfo(ctx context.Context, args contract.SlashingInfoArgs) (bool, *big.Int) {
	jailed, missed, err := s.keeper.SlashingInfo(ctx, args)
	s.requireError(err)
	return jailed, missed
}

func (s StakingPrecompileSuite) ApproveShares(ctx context.Context, from common.Address, args contract.ApproveSharesArgs) *evmtypes.MsgEthereumTxResponse {
	res, err := s.keeper.ApproveShares(ctx, from, args)
	s.requireError(err)
	return res
}

func (s StakingPrecompileSuite) TransferShares(ctx context.Context, from common.Address, args contract.TransferSharesArgs) (*evmtypes.MsgEthereumTxResponse, *contract.TransferSharesRet) {
	res, ret, err := s.keeper.TransferShares(ctx, from, args)
	s.requireError(err)
	return res, ret
}

func (s StakingPrecompileSuite) TransferFromShares(ctx context.Context, from common.Address, args contract.TransferFromSharesArgs) (*evmtypes.MsgEthereumTxResponse, *contract.TransferFromSharesRet) {
	res, ret, err := s.keeper.TransferFromShares(ctx, from, args)
	s.requireError(err)
	return res, ret
}

func (s StakingPrecompileSuite) Withdraw(ctx context.Context, from common.Address, args contract.WithdrawArgs) (*evmtypes.MsgEthereumTxResponse, *big.Int) {
	res, reward, err := s.keeper.Withdraw(ctx, from, args)
	s.requireError(err)
	return res, reward
}

func (s StakingPrecompileSuite) DelegateV2(ctx context.Context, from common.Address, args contract.DelegateV2Args, value ...*big.Int) *evmtypes.MsgEthereumTxResponse {
	if len(value) == 0 {
		value = []*big.Int{big.NewInt(0)}
	}
	res, err := s.keeper.DelegateV2(ctx, from, value[0], args)
	s.requireError(err)
	return res
}

func (s StakingPrecompileSuite) RedelegateV2(ctx context.Context, from common.Address, args contract.RedelegateV2Args) *evmtypes.MsgEthereumTxResponse {
	res, err := s.keeper.RedelegateV2(ctx, from, args)
	s.requireError(err)
	return res
}

func (s StakingPrecompileSuite) UndelegateV2(ctx context.Context, from common.Address, args contract.UndelegateV2Args) *evmtypes.MsgEthereumTxResponse {
	res, err := s.keeper.UndelegateV2(ctx, from, args)
	s.requireError(err)
	return res
}
