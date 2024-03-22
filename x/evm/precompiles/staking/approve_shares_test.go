package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
)

func TestStakingApproveSharesABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.ApproveSharesMethodName]
	require.Equal(t, method, staking.ApproveSharesMethod)
	require.Equal(t, 3, len(staking.ApproveSharesMethod.Inputs))
	require.Equal(t, 1, len(staking.ApproveSharesMethod.Outputs))

	event := stakingABI.Events[staking.ApproveSharesEventName]
	require.Equal(t, event, staking.ApproveSharesEvent)
	require.Equal(t, 4, len(staking.ApproveSharesEvent.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestApproveShares() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.ApproveSharesMethodName, val.String(), spender.Address(), allowance.BigInt())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "ok - approve zero",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.ApproveSharesMethodName, val.String(), spender.Address(), big.NewInt(0))
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string) {
				valStr := val.String() + "1"
				pack, err := staking.GetABI().Pack(staking.ApproveSharesMethodName, valStr, spender.Address(), allowance.BigInt())
				suite.Require().NoError(err)
				return pack, []string{valStr}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string) {
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestApproveSharesName, val.String(), spender.Address(), allowance.BigInt())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - ok - approve zero",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string) {
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestApproveSharesName, val.String(), spender.Address(), big.NewInt(0))
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) ([]byte, []string) {
				valStr := val.String() + "1"
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestApproveSharesName, valStr, spender.Address(), allowance.BigInt())
				suite.Require().NoError(err)
				return pack, []string{valStr}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: approve shares failed: invalid validator address: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val := vals[0]
			owner := suite.RandSigner()
			spender := suite.RandSigner()
			allowanceAmt := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))

			contract := staking.GetAddress()
			sender := owner.Address()
			if strings.HasPrefix(tc.name, "contract") {
				contract = suite.staking
				sender = suite.staking
			}

			allowance := suite.app.StakingKeeper.GetAllowance(suite.ctx, val.GetOperator(), owner.AccAddress(), spender.AccAddress())
			suite.Require().Equal(0, allowance.Cmp(big.NewInt(0)))

			pack, errArgs := tc.malleate(val.GetOperator(), spender, allowanceAmt)

			tx, err := suite.PackEthereumTx(owner, contract, big.NewInt(0), pack)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				allowance = suite.app.StakingKeeper.GetAllowance(suite.ctx, val.GetOperator(), sender.Bytes(), spender.AccAddress())
				if allowance.Cmp(big.NewInt(0)) != 0 {
					suite.Require().Equal(0, allowance.Cmp(allowanceAmt.BigInt()))
				}

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == staking.ApproveSharesEvent.ID.String() {
						suite.Require().Equal(log.Address, staking.GetAddress().String())
						suite.Require().Equal(log.Topics[1], sender.Hash().String())
						suite.Require().Equal(log.Topics[2], spender.Address().Hash().String())
						unpack, err := staking.ApproveSharesEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						unpackValidator := unpack[0].(string)
						suite.Require().Equal(unpackValidator, val.GetOperator().String())
						shares := unpack[1].(*big.Int)
						if allowance.Cmp(big.NewInt(0)) != 0 {
							suite.Require().Equal(shares.String(), allowanceAmt.BigInt().String())
						}
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type == fxstakingtypes.EventTypeApproveShares {
						for _, attr := range event.Attributes {
							if string(attr.Key) == stakingtypes.AttributeKeyValidator {
								suite.Require().Equal(string(attr.Value), val.GetOperator().String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeyOwner {
								suite.Require().Equal(string(attr.Value), sdk.AccAddress(sender.Bytes()).String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeySpender {
								suite.Require().Equal(string(attr.Value), spender.AccAddress().String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeyShares {
								if strings.Contains(tc.name, "zero") {
									suite.Require().Equal(string(attr.Value), "0")
								} else {
									suite.Require().Equal(string(attr.Value), allowanceAmt.String())
								}
							}
						}
						existEvent = true
					}
				}
				suite.Require().True(existEvent)
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				} else {
					if res.VmError != vm.ErrExecutionReverted.Error() {
						suite.Require().Equal(tc.error(errArgs), res.VmError)
					} else {
						if len(res.Ret) > 0 {
							reason, err := abi.UnpackRevert(common.CopyBytes(res.Ret))
							suite.Require().NoError(err)

							suite.Require().Equal(tc.error(errArgs), reason)
						} else {
							suite.Require().Equal(tc.error(errArgs), vm.ErrExecutionReverted.Error())
						}
					}
				}
			}
		})
	}
}
