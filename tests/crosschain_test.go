package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
)

func (suite *IntegrationTest) CrossChainTest() {
	for _, chain := range suite.crosschain {
		chain.Init()

		tokenAddress := helpers.GenerateAddressByModule(chain.chainName)
		metadata := fxtypes.GetCrossChainMetadata("test token", helpers.NewRandSymbol(), 18)

		bridgeDenom := fmt.Sprintf("%s%s", chain.chainName, tokenAddress)
		channelIBCHex := ""
		if chain.chainName == bsctypes.ModuleName {
			channelIBCHex = hex.EncodeToString([]byte("transfer/channel-0"))
			trace, err := fxtypes.GetIbcDenomTrace(bridgeDenom, channelIBCHex)
			suite.NoError(err)
			bridgeDenom = trace.IBCDenom()
			metadata = fxtypes.GetCrossChainMetadata("ibc token", bridgeDenom, 18)

			suite.erc20.RegisterCoinProposal(metadata)
		}
		chain.SendUpdateChainOraclesProposal()

		chain.BondedOracle()
		chain.SendOracleSetConfirm()

		chain.AddBridgeTokenClaim(metadata.Name, metadata.Symbol,
			uint64(metadata.DenomUnits[1].Exponent), tokenAddress, channelIBCHex)

		if len(channelIBCHex) > 0 {
			channelIbc, err := hex.DecodeString(channelIBCHex)
			suite.NoError(err)
			target := fmt.Sprintf("px/%s", string(channelIbc))
			chain.SendToFxClaim(tokenAddress, sdkmath.NewInt(100), target)
			chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdkmath.NewInt(0)))

			ibcTransferAddr := authtypes.NewModuleAddress(ibctransfertypes.ModuleName)
			chain.CheckBalance(ibcTransferAddr, sdk.NewCoin(bridgeDenom, sdkmath.NewInt(0)))
		}
		chain.SendToFxClaim(tokenAddress, sdkmath.NewInt(100), "")
		chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdkmath.NewInt(100)))

		txId := chain.SendToExternal(5, sdk.NewCoin(bridgeDenom, sdkmath.NewInt(10)))
		suite.True(txId > 0)
		chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdkmath.NewInt(50)))

		chain.CheckIncreaseBridgeFee(chain.AccAddress(), txId)

		chain.SendBatchRequest(5)
		chain.SendConfirmBatch()

		chain.SendToExternalAndCancel(sdk.NewCoin(bridgeDenom, sdkmath.NewInt(40)))
		chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdkmath.NewInt(40)))

		if chain.chainName == ethtypes.ModuleName {
			fxTokenAddress := helpers.GenerateAddress().Hex()
			fxMD := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
			chain.AddBridgeTokenClaim(fxMD.Name, fxMD.Symbol, fxtypes.DenomUnit, fxTokenAddress, "")

			// send fx to chain
			balance := suite.QueryBalances(chain.AccAddress())
			chain.SendToFxClaim(fxTokenAddress, sdkmath.NewInt(100), "")
			chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, balance.AmountOf(fxtypes.DefaultDenom).Add(sdkmath.NewInt(100))))

			// send fx to evm
			fxPair := suite.erc20.TokenPair(fxtypes.DefaultDenom)
			erc20Balance := suite.erc20.BalanceOf(fxPair.GetERC20Contract(), chain.HexAddress())
			chain.SendToFxClaim(fxTokenAddress, sdkmath.NewInt(100), fxtypes.ERC20Target)
			chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, balance.AmountOf(fxtypes.DefaultDenom).Add(sdkmath.NewInt(100))))
			suite.Equal(big.NewInt(0).Add(erc20Balance, big.NewInt(100)), suite.erc20.BalanceOf(fxPair.GetERC20Contract(), chain.HexAddress()))
		}
	}
}
