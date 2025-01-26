package crosschain_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestExecuteClaimMethod_ABI(t *testing.T) {
	executeClaimABI := crosschain.NewExecuteClaimABI()

	methodStr := `function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)`
	assert.Equal(t, methodStr, executeClaimABI.Method.String())

	eventStr := `event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain, string _errReason)`
	assert.Equal(t, eventStr, executeClaimABI.Event.String())
}

func TestExecuteClaimMethod_PackInput(t *testing.T) {
	executeClaimABI := crosschain.NewExecuteClaimABI()
	input, err := executeClaimABI.PackInput(contract.ExecuteClaimArgs{
		Chain:      ethtypes.ModuleName,
		EventNonce: big.NewInt(1),
	})
	require.NoError(t, err)
	expected := "4ac3bdc30000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000036574680000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expected, hex.EncodeToString(input))
}

func (suite *CrosschainPrecompileTestSuite) TestContract_ExecuteClaim_SendToFx() {
	testCases := []struct {
		name            string
		malleate        func() erc20types.BridgeToken
		amount          sdkmath.Int
		erc20Balance    *big.Int
		moduleAmount    sdkmath.Int
		baseDenomAmount sdkmath.Int
		logsLen         int
	}{
		{
			name: "native coin",
			malleate: func() erc20types.BridgeToken {
				symbol := helpers.NewRandSymbol()
				suite.AddBridgeToken(symbol, true)

				bridgeToken := suite.GetBridgeToken(strings.ToLower(symbol))
				return bridgeToken
			},
			amount:          sdkmath.NewInt(100),
			erc20Balance:    big.NewInt(100),
			moduleAmount:    sdkmath.NewInt(100),
			baseDenomAmount: sdkmath.NewInt(0),
			logsLen:         2,
		},
		{
			name: "native erc20",
			malleate: func() erc20types.BridgeToken {
				symbol := helpers.NewRandSymbol()

				erc20TokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, symbol)
				suite.AddBridgeToken(erc20TokenAddr.String(), false)

				// eth module lock some tokens
				bridgeToken := suite.GetBridgeToken(strings.ToLower(symbol))
				suite.MintTokenToModule(suite.chainName, sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(100)))

				// erc20 module lock some tokens
				suite.MintTokenToModule(crosschaintypes.ModuleName, sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(100)))

				erc20ModuelAddr := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName))
				suite.erc20TokenSuite.Mint(suite.Ctx, suite.signer.Address(), erc20ModuelAddr, big.NewInt(100))
				return bridgeToken
			},
			amount:          sdkmath.NewInt(100),
			erc20Balance:    big.NewInt(100),
			moduleAmount:    sdkmath.NewInt(0),
			baseDenomAmount: sdkmath.NewInt(0),
			logsLen:         2,
		},
		{
			name: "original token",
			malleate: func() erc20types.BridgeToken {
				suite.AddBridgeToken(fxtypes.DefaultSymbol, false)

				// eth module lock some tokens
				bridgeToken := suite.GetBridgeToken(fxtypes.DefaultDenom)
				suite.MintTokenToModule(ethtypes.ModuleName, sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(100)))

				return bridgeToken
			},
			amount:          sdkmath.NewInt(100),
			erc20Balance:    big.NewInt(0),
			moduleAmount:    sdkmath.NewInt(0),
			baseDenomAmount: sdkmath.NewInt(100),
			logsLen:         1,
		},
		{
			name: "legacy FX token",
			malleate: func() erc20types.BridgeToken {
				suite.AddBridgeToken(fxtypes.LegacyFXDenom, false)
				suite.AddBridgeToken(fxtypes.DefaultSymbol, false)

				// eth module lock some tokens
				bridgeToken := suite.GetBridgeToken(fxtypes.DefaultDenom)
				suite.MintTokenToModule(ethtypes.ModuleName, sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(1)))

				return suite.GetBridgeToken(fxtypes.FXDenom)
			},
			amount:          sdkmath.NewInt(100),
			erc20Balance:    big.NewInt(0),
			moduleAmount:    sdkmath.NewInt(0),
			baseDenomAmount: sdkmath.NewInt(1),
			logsLen:         1,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			bridgeToken := tc.malleate()

			receiver := helpers.GenAccAddress()
			claim := &crosschaintypes.MsgSendToFxClaim{
				EventNonce:     1,
				BlockHeight:    100,
				TokenContract:  bridgeToken.Contract,
				Amount:         tc.amount,
				Sender:         helpers.GenExternalAddr(suite.chainName),
				Receiver:       receiver.String(),
				BridgerAddress: helpers.GenAccAddress().String(),
				ChainName:      suite.chainName,
			}
			keeper := suite.App.CrosschainKeepers.GetKeeper(suite.chainName)
			err := keeper.SavePendingExecuteClaim(suite.Ctx, claim)
			suite.Require().NoError(err)

			txResponse := suite.ExecuteClaim(suite.Ctx, suite.signer.Address(),
				contract.ExecuteClaimArgs{Chain: suite.chainName, EventNonce: big.NewInt(int64(claim.EventNonce))},
			)
			suite.NotNil(txResponse)
			suite.Len(txResponse.Logs, tc.logsLen)
			event, err := crosschain.NewExecuteClaimABI().
				UnpackEvent(txResponse.Logs[len(txResponse.Logs)-1].ToEthereum())
			suite.Require().NoError(err)
			suite.Equal(suite.GetSender(), event.Sender)
			suite.Equal(claim.EventNonce, event.EventNonce.Uint64())
			suite.Equal(claim.ChainName, event.Chain)
			suite.Empty(event.ErrReason)

			// crosschain module balance
			bridgeTokenBalance := sdk.NewCoin(bridgeToken.BridgeDenom(), tc.moduleAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), bridgeTokenBalance)

			// erc20 module balance
			baseDenomBalance := sdk.NewCoin(bridgeToken.Denom, tc.moduleAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseDenomBalance)

			baseDenom := bridgeToken.Denom
			if baseDenom == fxtypes.FXDenom { // swap base denom to apundiai denom
				bridgeTokenBalance = sdk.NewCoin(bridgeToken.BridgeDenom(), tc.amount)
				suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeTokenBalance)
				baseDenom = fxtypes.DefaultDenom
			}

			// receiver balance
			baseDenomBalance = sdk.NewCoin(baseDenom, tc.baseDenomAmount)
			suite.AssertBalance(receiver, baseDenomBalance)

			erc20Contract := suite.GetERC20Token(baseDenom).GetERC20Contract()
			suite.erc20TokenSuite.WithContract(erc20Contract)

			balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, common.BytesToAddress(receiver))
			suite.Equal(tc.erc20Balance.String(), balance.String())
		})
	}
}
