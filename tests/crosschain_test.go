package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

func (suite *IntegrationTest) CrossChainTest() {
	for index := 0; index < len(suite.crosschain); index++ {
		suite.crosschain[index].Init()
		chain := suite.crosschain[index]

		tokenAddress := helpers.GenerateAddressByModule(chain.chainName)
		metadata := fxtypes.GetCrossChainMetadataManyToOne("test token", helpers.NewRandSymbol(), 18)

		bridgeDenom := crosschaintypes.NewBridgeDenom(chain.chainName, tokenAddress)
		channelIBCHex := ""
		if chain.chainName == bsctypes.ModuleName {
			channelIBCHex = hex.EncodeToString([]byte("transfer/channel-0"))
			trace, err := fxtypes.GetIbcDenomTrace(bridgeDenom, channelIBCHex)
			suite.NoError(err)
			bridgeDenom = trace.IBCDenom()
			metadata = fxtypes.GetCrossChainMetadataOneToOne("ibc token", bridgeDenom, "PURSE", 18)

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

			// add pundix token
			pundixAddress := helpers.GenerateAddressByModule(chain.chainName)
			pundixDenom := crosschaintypes.NewBridgeDenom(ethtypes.ModuleName, pundixAddress)
			pundixMetadata := fxtypes.GetCrossChainMetadataOneToOne("test token", pundixDenom, "PUNDIX", 18)
			suite.erc20.RegisterCoinProposal(pundixMetadata)

			chain.AddBridgeTokenClaim(pundixMetadata.Name, pundixMetadata.Symbol,
				uint64(pundixMetadata.DenomUnits[1].Exponent), pundixAddress, "")
		}
	}
}

func (suite *IntegrationTest) OriginalCrossChainTest() {
	ethChain := suite.GetCrossChainByName(ethtypes.ModuleName)
	bscChain := suite.GetCrossChainByName(bsctypes.ModuleName)
	tronChain := suite.GetCrossChainByName(trontypes.ModuleName)

	// eth add purse token
	purseMd := ethChain.SelectTokenMetadata("ibc/")

	newTokenContract := helpers.GenerateAddressByModule(ethtypes.ModuleName)
	purseNewAlias := crosschaintypes.NewBridgeDenom(ethtypes.ModuleName, newTokenContract)
	resp, _ := suite.erc20.UpdateDenomAliasProposal(purseMd.Base, purseNewAlias)
	suite.Equal(uint32(0), resp.Code)

	ethChain.AddBridgeTokenClaim("PURSE Token", "PURSE", 18, newTokenContract, "")
	purseTokenEth := newTokenContract

	// bsc add FX token
	newTokenContract = helpers.GenerateAddressByModule(bsctypes.ModuleName)

	bscChain.AddBridgeTokenClaim("Function X", "FX", 18, newTokenContract, "")
	fxTokenBSC := newTokenContract

	// polygon add pundix token
	pundixMd := tronChain.SelectTokenMetadata("eth")

	newTokenContract = helpers.GenerateAddressByModule(trontypes.ModuleName)
	pundixAlias := crosschaintypes.NewBridgeDenom(trontypes.ModuleName, newTokenContract)
	resp, _ = suite.erc20.UpdateDenomAliasProposal(pundixMd.Base, pundixAlias)
	suite.Equal(uint32(0), resp.Code)

	tronChain.AddBridgeTokenClaim("Pundix Token", "PUNDIX", 18, newTokenContract, "")
	pundixTokenPolygon := newTokenContract

	// init amount
	initAmount := sdkmath.NewInt(1000)

	bscChain.SendToExternalAndConfirm(sdk.NewCoin(fxtypes.DefaultDenom, initAmount))

	bscPurseTokenAddress := bscChain.GetBridgeTokenByDenom(purseMd.Base)
	bscChain.SendToFxClaim(bscPurseTokenAddress, initAmount, "")
	bscChain.Send(ethChain.AccAddress(), sdk.NewCoin(purseMd.Base, initAmount))
	ethChain.SendToExternalAndConfirm(sdk.NewCoin(purseMd.Base, initAmount))

	ethPundixTokenAddress := ethChain.GetBridgeTokenByDenom(pundixMd.Base)
	ethChain.SendToFxClaim(ethPundixTokenAddress, initAmount, "")
	ethChain.Send(tronChain.AccAddress(), sdk.NewCoin(pundixMd.Base, initAmount))
	tronChain.SendToExternalAndConfirm(sdk.NewCoin(pundixMd.Base, initAmount))

	// send to fx

	// fx
	fxTokenAddress := ethChain.GetBridgeTokenByDenom(fxtypes.DefaultDenom)
	ethChain.SendToFxClaimAndCheckBalance(fxTokenAddress, sdkmath.NewInt(100), "", sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100)))
	bscChain.SendToFxClaimAndCheckBalance(fxTokenBSC, sdkmath.NewInt(100), "", sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100)))

	bscBalances := suite.QueryBalances(authtypes.NewModuleAddress(bsctypes.ModuleName))
	suite.Equal(initAmount.Sub(sdkmath.NewInt(100)), bscBalances.AmountOf(fxtypes.DefaultDenom))

	// pundix
	ethChain.SendToFxClaimAndCheckBalance(ethPundixTokenAddress, sdkmath.NewInt(200), "", sdk.NewCoin(pundixMd.Base, sdkmath.NewInt(200)))
	tronChain.SendToFxClaimAndCheckBalance(pundixTokenPolygon, sdkmath.NewInt(200), "", sdk.NewCoin(pundixMd.Base, sdkmath.NewInt(200)))

	tronBalances := suite.QueryBalances(authtypes.NewModuleAddress(trontypes.ModuleName))
	suite.Equal(initAmount.Sub(sdkmath.NewInt(200)), tronBalances.AmountOf(pundixAlias))
	pxAliasSupply, err := suite.GRPCClient().BankQuery().SupplyOf(suite.ctx, &banktypes.QuerySupplyOfRequest{Denom: pundixAlias})
	suite.NoError(err)
	suite.Equal(pxAliasSupply.Amount.Amount, tronBalances.AmountOf(pundixAlias))

	// purse
	bscChain.SendToFxClaimAndCheckBalance(bscPurseTokenAddress, sdkmath.NewInt(300), "", sdk.NewCoin(purseMd.Base, sdkmath.NewInt(300)))
	ethChain.SendToFxClaimAndCheckBalance(purseTokenEth, sdkmath.NewInt(300), "", sdk.NewCoin(purseMd.Base, sdkmath.NewInt(300)))

	ethBalances := suite.QueryBalances(authtypes.NewModuleAddress(ethtypes.ModuleName))
	suite.Equal(initAmount.Sub(sdkmath.NewInt(300)), ethBalances.AmountOf(purseNewAlias))
	purseAliasSupply, err := suite.GRPCClient().BankQuery().SupplyOf(suite.ctx, &banktypes.QuerySupplyOfRequest{Denom: purseNewAlias})
	suite.NoError(err)
	suite.Equal(purseAliasSupply.Amount.Amount, ethBalances.AmountOf(purseNewAlias))

	// send to external
	// fx eth
	ethChain.SendToExternalAndCheckBalance(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100)))
	bscChain.SendToExternalAndCheckBalance(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100)))

	bscBalances = suite.QueryBalances(authtypes.NewModuleAddress(bsctypes.ModuleName))
	suite.Equal(initAmount, bscBalances.AmountOf(fxtypes.DefaultDenom))

	// pundix
	ethChain.SendToExternalAndCheckBalance(sdk.NewCoin(pundixMd.Base, sdkmath.NewInt(100)))
	tronChain.SendToExternalAndCheckBalance(sdk.NewCoin(pundixMd.Base, sdkmath.NewInt(100)))

	tronBalances = suite.QueryBalances(authtypes.NewModuleAddress(trontypes.ModuleName))
	suite.Equal(initAmount.Sub(sdkmath.NewInt(100)), tronBalances.AmountOf(pundixAlias))
	pxAliasSupply, err = suite.GRPCClient().BankQuery().SupplyOf(suite.ctx, &banktypes.QuerySupplyOfRequest{Denom: pundixAlias})
	suite.NoError(err)
	suite.Equal(initAmount.Sub(sdkmath.NewInt(100)), pxAliasSupply.Amount.Amount)

	// purse
	bscChain.SendToExternalAndCheckBalance(sdk.NewCoin(purseMd.Base, sdkmath.NewInt(100)))
	ethChain.SendToExternalAndCheckBalance(sdk.NewCoin(purseMd.Base, sdkmath.NewInt(100)))

	ethBalances = suite.QueryBalances(authtypes.NewModuleAddress(ethtypes.ModuleName))
	suite.Equal(initAmount.Sub(sdkmath.NewInt(200)), ethBalances.AmountOf(purseNewAlias))
	purseAliasSupply, err = suite.GRPCClient().BankQuery().SupplyOf(suite.ctx, &banktypes.QuerySupplyOfRequest{Denom: purseNewAlias})
	suite.NoError(err)
	suite.Equal(initAmount.Sub(sdkmath.NewInt(200)), purseAliasSupply.Amount.Amount)

	// clear send to external tx
	ethChain.CancelAllSendToExternal()
	bscChain.CancelAllSendToExternal()
	tronChain.CancelAllSendToExternal()
}
