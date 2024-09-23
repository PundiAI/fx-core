package testutil

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/testutil"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type StakingPrecompileSuite struct {
	testutil.EVMSuite
}

func (s *StakingPrecompileSuite) EthereumTx(data []byte, value *big.Int, gasLimit uint64, success bool) *evmtypes.MsgEthereumTxResponse {
	toAddr := s.EVMSuite.GetContractAddr()
	if gasLimit == 0 {
		gasLimit = fxcontract.DefaultGasCap
	}
	if gasLimit < 10_000 {
		gasLimit += 21_000 + gasLimit
	}
	tx, err := s.EVMSuite.EthereumTx(toAddr, data, value, gasLimit)
	if success {
		s.NoError(err)
	} else {
		s.Error(err)
	}
	return tx
}

func (s *StakingPrecompileSuite) Allowance(validator sdk.ValAddress, owner, spender common.Address) *big.Int {
	method := precompile.NewAllowanceSharesMethod(nil)
	data, err := method.PackInput(fxstakingtypes.AllowanceSharesArgs{
		Validator: validator.String(),
		Owner:     owner,
		Spender:   spender,
	})
	s.NoError(err)
	tx := s.CallEVM(data, method.RequiredGas())
	output, err := method.UnpackOutput(tx.Ret)
	s.NoError(err)
	return output
}

func (s *StakingPrecompileSuite) Delegation(validator sdk.ValAddress, delegator common.Address) (*big.Int, *big.Int) {
	method := precompile.NewDelegationMethod(nil)
	data, err := method.PackInput(fxstakingtypes.DelegationArgs{
		Validator: validator.String(),
		Delegator: delegator,
	})
	s.NoError(err)
	tx := s.CallEVM(data, method.RequiredGas())
	shares, amount, err := method.UnpackOutput(tx.Ret)
	s.NoError(err)
	return shares, amount
}

func (s *StakingPrecompileSuite) DelegationRewards(validator sdk.ValAddress, delegator common.Address) *big.Int {
	method := precompile.NewDelegationRewardsMethod(nil)
	data, err := method.PackInput(fxstakingtypes.DelegationRewardsArgs{
		Validator: validator.String(),
		Delegator: delegator,
	})
	s.NoError(err)
	tx := s.CallEVM(data, method.RequiredGas())
	rewards, err := method.UnpackOutput(tx.Ret)
	s.NoError(err)
	return rewards
}

func (s *StakingPrecompileSuite) Approve(validator sdk.ValAddress, spender common.Address, shares *big.Int, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewApproveSharesMethod(nil)
	data, err := method.PackInput(fxstakingtypes.ApproveSharesArgs{
		Validator: validator.String(),
		Spender:   spender,
		Shares:    shares,
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}

func (s *StakingPrecompileSuite) TransferShares(validator sdk.ValAddress, to common.Address, shares *big.Int, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewTransferSharesMethod(nil)
	data, err := method.PackInput(fxstakingtypes.TransferSharesArgs{
		Validator: validator.String(),
		To:        to,
		Shares:    shares,
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}

func (s *StakingPrecompileSuite) TransferFromShares(validator sdk.ValAddress, from, to common.Address, shares *big.Int, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewTransferFromSharesMethod(nil)
	data, err := method.PackInput(fxstakingtypes.TransferFromSharesArgs{
		Validator: validator.String(),
		From:      from,
		To:        to,
		Shares:    shares,
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}

func (s *StakingPrecompileSuite) Withdraw(validator sdk.ValAddress, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewWithdrawMethod(nil)
	data, err := method.PackInput(fxstakingtypes.WithdrawArgs{
		Validator: validator.String(),
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}

func (s *StakingPrecompileSuite) DelegateV2(validator sdk.ValAddress, amount *big.Int, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewDelegateV2Method(nil)
	data, err := method.PackInput(fxstakingtypes.DelegateV2Args{
		Validator: validator.String(),
		Amount:    amount,
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}

func (s *StakingPrecompileSuite) RedelegateV2(validatorSrc, validatorDst sdk.ValAddress, amount *big.Int, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewRedelegateV2Method(nil)
	data, err := method.PackInput(fxstakingtypes.RedelegateV2Args{
		ValidatorSrc: validatorSrc.String(),
		ValidatorDst: validatorDst.String(),
		Amount:       amount,
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}

func (s *StakingPrecompileSuite) UndelegateV2(validator sdk.ValAddress, amount *big.Int, success bool) *evmtypes.MsgEthereumTxResponse {
	method := precompile.NewUndelegateV2Method(nil)
	data, err := method.PackInput(fxstakingtypes.UndelegateV2Args{
		Validator: validator.String(),
		Amount:    amount,
	})
	s.NoError(err)
	return s.EthereumTx(data, nil, method.RequiredGas(), success)
}
