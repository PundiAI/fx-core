package tests

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *IntegrationTest) CrossChainTest() {
	for _, chain := range suite.crosschain {
		chain.Init()

		tokenAddress := helpers.GenerateAddress().Hex()
		if chain.chainName == trontypes.ModuleName {
			tokenAddress = trontypes.AddressFromHex(tokenAddress)
		}

		bridgeDenom := fmt.Sprintf("%s%s", chain.chainName, tokenAddress)
		tokenChannelIBC := ""
		if chain.chainName == bsctypes.ModuleName {
			tokenChannelIBC = "transfer/channel-0"

			suite.erc20.metadata.DenomUnits[0].Aliases = []string{
				ibctransfertypes.DenomTrace{
					Path:      tokenChannelIBC,
					BaseDenom: bridgeDenom,
				}.IBCDenom(),
				bridgeDenom,
			}
			suite.T().Log(suite.erc20.metadata.DenomUnits[0].Aliases)

			proposalId := suite.erc20.RegisterCoinProposal(suite.erc20.metadata)
			suite.NoError(suite.network.WaitForNextBlock())
			suite.CheckProposal(proposalId, govtypes.StatusPassed)
		}

		proposalId := chain.SendUpdateChainOraclesProposal()
		suite.NoError(suite.network.WaitForNextBlock())
		suite.CheckProposal(proposalId, govtypes.StatusPassed)

		chain.BondedOracle()
		chain.SendOracleSetConfirm()

		denom := chain.AddBridgeTokenClaim(suite.erc20.metadata.Name, suite.erc20.metadata.Symbol,
			uint64(suite.erc20.TokenDecimals()), tokenAddress, tokenChannelIBC)
		suite.Equal(denom, bridgeDenom)

		if len(tokenChannelIBC) > 0 {
			tokenChannelIBC = fmt.Sprintf("px/%s", tokenChannelIBC)
		}
		if len(tokenChannelIBC) > 0 {
			// todo need create ibc channel
			chain.SendToFxClaim(tokenAddress, sdk.NewInt(100), tokenChannelIBC)
			chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdk.NewInt(0)))
			chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(suite.erc20.metadata.DenomUnits[0].Aliases[0], sdk.NewInt(0)))

			ibcTransferAddr := authtypes.NewModuleAddress(ibctransfertypes.ModuleName)
			chain.CheckBalance(ibcTransferAddr, sdk.NewCoin(bridgeDenom, sdk.NewInt(0)))
			chain.CheckBalance(ibcTransferAddr, sdk.NewCoin(suite.erc20.metadata.DenomUnits[0].Aliases[0], sdk.NewInt(0)))
		} else {
			chain.SendToFxClaim(tokenAddress, sdk.NewInt(100), "")
			chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdk.NewInt(100)))
		}

		txId := chain.SendToExternal(5, sdk.NewCoin(bridgeDenom, sdk.NewInt(10)))
		suite.True(txId > 0)
		chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdk.NewInt(50)))

		chain.SendBatchRequest(5)
		chain.SendConfirmBatch()

		chain.SendToExternalAndCancel(sdk.NewCoin(bridgeDenom, sdk.NewInt(50)))
		chain.CheckBalance(chain.AccAddress(), sdk.NewCoin(bridgeDenom, sdk.NewInt(50)))

		chain.SendFrom(chain.privKey, suite.erc20.AccAddress(), sdk.NewCoin(bridgeDenom, sdk.NewInt(50)))
	}
}
