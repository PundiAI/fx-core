package tests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

func (suite *IntegrationTest) PrecompileTransferCrossChainTest() {
	suite.Send(suite.precompile.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	tokenPair, bridgeTokens := suite.precompileInit()
	for i, chain := range suite.crosschain {
		bridgeToken := bridgeTokens[i]
		suite.erc20.HandleWithCheckBalance(tokenPair.GetERC20Contract(), suite.precompile.HexAddress(), big.NewInt(200), func() {
			chain.SendToTxClaimWithReceiver(suite.precompile.AccAddress(), bridgeToken.Token, sdkmath.NewInt(200), fxtypes.LegacyERC20Target)
		})

		receive := chain.FormatAddress(suite.precompile.HexAddress())
		suite.precompile.TransferCrossChainAndCheckPendingTx(tokenPair.GetERC20Contract(), receive,
			big.NewInt(20), big.NewInt(30), chain.chainName)

		chain.CancelAllSendToExternal()
	}
}

func (suite *IntegrationTest) PrecompileCrossChainTest() {
	suite.Send(suite.precompile.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	tokenPair, bridgeTokens := suite.precompileInit()
	for i, chain := range suite.crosschain {
		bridgeToken := bridgeTokens[i]
		suite.erc20.HandleWithCheckBalance(tokenPair.GetERC20Contract(), suite.precompile.HexAddress(), big.NewInt(200), func() {
			chain.SendToTxClaimWithReceiver(suite.precompile.AccAddress(), bridgeToken.Token, sdkmath.NewInt(200), fxtypes.LegacyERC20Target)
		})

		receive := chain.FormatAddress(suite.precompile.HexAddress())
		suite.precompile.CrossChain(tokenPair.GetERC20Contract(), receive,
			big.NewInt(20), big.NewInt(30), chain.chainName)

		chain.CancelAllSendToExternal()
	}
}

func (suite *IntegrationTest) PrecompileCancelSendToExternalTest() {
	suite.Send(suite.precompile.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	tokenPair, bridgeTokens := suite.precompileInit()
	for i, chain := range suite.crosschain {
		bridgeToken := bridgeTokens[i]
		suite.erc20.HandleWithCheckBalance(tokenPair.GetERC20Contract(), suite.precompile.HexAddress(), big.NewInt(200), func() {
			chain.SendToTxClaimWithReceiver(suite.precompile.AccAddress(), bridgeToken.Token, sdkmath.NewInt(200), fxtypes.LegacyERC20Target)
		})

		receive := chain.FormatAddress(suite.precompile.HexAddress())
		txId := suite.precompile.CrossChain(tokenPair.GetERC20Contract(), receive,
			big.NewInt(20), big.NewInt(30), chain.chainName)

		suite.precompile.CancelSendToExternalAndCheckPendingTx(chain.chainName, txId)
	}
}

func (suite *IntegrationTest) PrecompileIncreaseBridgeFeeTest() {
	suite.Send(suite.precompile.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	tokenPair, bridgeTokens := suite.precompileInit()
	for i, chain := range suite.crosschain {
		bridgeToken := bridgeTokens[i]
		suite.erc20.HandleWithCheckBalance(tokenPair.GetERC20Contract(), suite.precompile.HexAddress(), big.NewInt(200), func() {
			chain.SendToTxClaimWithReceiver(suite.precompile.AccAddress(), bridgeToken.Token, sdkmath.NewInt(200), fxtypes.LegacyERC20Target)
		})

		receive := chain.FormatAddress(suite.precompile.HexAddress())
		txId := suite.precompile.CrossChain(tokenPair.GetERC20Contract(), receive,
			big.NewInt(20), big.NewInt(30), chain.chainName)

		suite.precompile.IncreaseBridgeFeeCheckPendingTx(chain.chainName, txId, tokenPair.GetERC20Contract(), big.NewInt(50))

		chain.CancelAllSendToExternal()
	}
}

func (suite *IntegrationTest) PrecompileCrossChainConvertedDenomTest() {
	ethChain := suite.GetCrossChainByName(ethtypes.ModuleName)
	bscChain := suite.GetCrossChainByName(bsctypes.ModuleName)
	tronChain := suite.GetCrossChainByName(trontypes.ModuleName)

	receiver := suite.precompile.AccAddress()
	receiverHex := common.BytesToAddress(receiver.Bytes())

	fxMd := suite.GetMetadata(fxtypes.DefaultDenom)
	fxPair := suite.erc20.TokenPair(fxMd.Base)
	// fx
	suite.erc20.HandleWithCheckBalance(fxPair.GetERC20Contract(), receiverHex, big.NewInt(100), func() {
		ethFXTokenAddress := ethChain.GetBridgeTokenByDenom(fxtypes.DefaultDenom)
		ethChain.SendToTxClaimWithReceiver(receiver, ethFXTokenAddress, sdkmath.NewInt(100), "erc20")
	})
	suite.erc20.HandleWithCheckBalance(fxPair.GetERC20Contract(), receiverHex, big.NewInt(100), func() {
		bscFXTokenAddress := bscChain.GetBridgeTokenByDenom(fxtypes.DefaultDenom)
		bscChain.SendToTxClaimWithReceiver(receiver, bscFXTokenAddress, sdkmath.NewInt(100), "erc20")
	})

	// fx cross chain
	suite.precompile.CrossChain(fxPair.GetERC20Contract(),
		ethChain.HexAddressString(), big.NewInt(10), big.NewInt(1), ethtypes.ModuleName)
	suite.precompile.CrossChain(fxPair.GetERC20Contract(),
		bscChain.HexAddressString(), big.NewInt(11), big.NewInt(1), bsctypes.ModuleName)

	// purse
	purseMd := bscChain.SelectTokenMetadata("ibc/")
	pursePair := suite.erc20.TokenPair(purseMd.Base)

	suite.erc20.HandleWithCheckBalance(pursePair.GetERC20Contract(), receiverHex, big.NewInt(200), func() {
		ethPURSETokenAddress := ethChain.GetBridgeTokenByDenom(purseMd.DenomUnits[0].Aliases[0])
		ethChain.SendToTxClaimWithReceiver(receiver, ethPURSETokenAddress, sdkmath.NewInt(200), "erc20")
	})
	suite.erc20.HandleWithCheckBalance(pursePair.GetERC20Contract(), receiverHex, big.NewInt(200), func() {
		bscPurseTokenAddress := bscChain.GetBridgeTokenByDenom(purseMd.Base)
		bscChain.SendToTxClaimWithReceiver(receiver, bscPurseTokenAddress, sdkmath.NewInt(200), "erc20")
	})

	// purse cross chain
	suite.precompile.CrossChain(pursePair.GetERC20Contract(),
		ethChain.HexAddressString(), big.NewInt(20), big.NewInt(1), ethtypes.ModuleName)
	suite.precompile.CrossChain(pursePair.GetERC20Contract(),
		bscChain.HexAddressString(), big.NewInt(21), big.NewInt(1), bsctypes.ModuleName)

	// pundix
	pundixMd := ethChain.SelectTokenMetadata("eth")
	pundixPair := suite.erc20.TokenPair(pundixMd.Base)

	suite.erc20.HandleWithCheckBalance(pundixPair.GetERC20Contract(), receiverHex, big.NewInt(300), func() {
		ethPundixTokenAddress := ethChain.GetBridgeTokenByDenom(pundixMd.Base)
		ethChain.SendToTxClaimWithReceiver(receiver, ethPundixTokenAddress, sdkmath.NewInt(300), "erc20")
	})
	suite.erc20.HandleWithCheckBalance(pundixPair.GetERC20Contract(), receiverHex, big.NewInt(300), func() {
		tronPundixTokenAddress := tronChain.GetBridgeTokenByDenom(pundixMd.DenomUnits[0].Aliases[0])
		tronChain.SendToTxClaimWithReceiver(receiver, tronPundixTokenAddress, sdkmath.NewInt(300), "erc20")
	})

	// pundix cross chain
	suite.precompile.CrossChain(pundixPair.GetERC20Contract(),
		ethChain.HexAddressString(), big.NewInt(30), big.NewInt(1), ethtypes.ModuleName)
	suite.precompile.CrossChain(pundixPair.GetERC20Contract(),
		tronChain.HexAddressString(), big.NewInt(31), big.NewInt(1), trontypes.ModuleName)

	// clear send to external tx
	ethChain.CancelAllSendToExternal()
	bscChain.CancelAllSendToExternal()
	tronChain.CancelAllSendToExternal()
}

func (suite *IntegrationTest) precompileInit() (*erc20types.TokenPair, []crosschaintypes.BridgeToken) {
	metadata := fxtypes.GetCrossChainMetadataManyToOne("test token", helpers.NewRandSymbol(), uint32(18))
	aliases := make([]string, 0, len(suite.crosschain))
	bridgeTokens := make([]crosschaintypes.BridgeToken, 0, len(suite.crosschain))
	for _, chain := range suite.crosschain {
		denom, bridgeToken := chain.AddBridgeToken(metadata)
		aliases = append(aliases, denom)
		bridgeTokens = append(bridgeTokens, bridgeToken)
	}
	metadata.DenomUnits[0].Aliases = aliases

	suite.erc20.RegisterCoinProposal(metadata)
	tokenPair := suite.erc20.TokenPair(metadata.Base)
	return tokenPair, bridgeTokens
}
