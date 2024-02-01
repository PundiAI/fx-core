package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v6/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v6/types"
	crosschaintypes "github.com/functionx/fx-core/v6/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v6/x/erc20/types"
)

func (suite *IntegrationTest) ERC20TokenOriginTest() {
	suite.Send(suite.erc20.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	decimals := 18
	metadata := fxtypes.GetCrossChainMetadataManyToOne("test token", helpers.NewRandSymbol(), uint32(decimals))
	aliases := make([]string, 0, len(suite.crosschain))
	bridgeTokens := make([]crosschaintypes.BridgeToken, 0, len(suite.crosschain))
	for _, chain := range suite.crosschain {
		denom, bridgeToken := chain.AddBridgeToken(metadata)
		aliases = append(aliases, denom)
		bridgeTokens = append(bridgeTokens, bridgeToken)
	}
	metadata.DenomUnits[0].Aliases = aliases

	suite.erc20.RegisterCoinProposal(metadata)
	suite.erc20.CheckRegisterCoin(metadata.Base)

	tokenPair := suite.erc20.TokenPair(metadata.Base)
	suite.Equal(tokenPair.Denom, metadata.Base)
	suite.Equal(tokenPair.Enabled, true)
	suite.Equal(tokenPair.ContractOwner, erc20types.OWNER_MODULE)

	for i, chain := range suite.crosschain {
		bridgeToken := bridgeTokens[i]
		chain.SendToTxClaimWithReceiver(suite.erc20.AccAddress(), bridgeToken.Token, sdkmath.NewInt(100), fxtypes.LegacyERC20Target)

		// covert chain.address erc20 token to native token: metadata.base
		suite.erc20.ConvertERC20(tokenPair.GetERC20Contract(), sdkmath.NewInt(50), suite.erc20.AccAddress())
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(50)))

		// covert erc20.addr metadata.base
		suite.erc20.ConvertDenom(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(50)), chain.chainName)
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(50)))

		// send to chain.address
		baseTokenBalanceAmount := suite.QueryBalances(chain.AccAddress()).AmountOf(metadata.Base)
		chain.SendToTxClaimWithReceiver(suite.erc20.AccAddress(), bridgeToken.Token, sdkmath.NewInt(100), "")
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, baseTokenBalanceAmount.Add(sdkmath.NewInt(100))))

		// convert native token(metadata base) to erc20 token
		suite.erc20.ConvertCoin(suite.erc20.HexAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(100)))

		if i < len(suite.crosschain)-1 {
			// remove proposal
			suite.erc20.UpdateDenomAliasProposal(metadata.Base, bridgeToken.Denom)

			// check remove
			response, err := suite.erc20.ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: metadata.Base})
			suite.NoError(err)
			suite.Equal(len(suite.crosschain)-i-1, len(response.Aliases))

			_, err = suite.erc20.ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: bridgeToken.Denom})
			suite.Error(err)
		}
	}

	suite.erc20.ToggleTokenConversionProposal(metadata.Base)

	suite.False(suite.erc20.TokenPair(metadata.Base).Enabled)
}

func (suite *IntegrationTest) ERC20IBCChainTokenOriginTest() {
	suite.Send(suite.erc20.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	portID := "transfer"
	channelID := "channel-0"

	for _, chain := range suite.crosschain {
		tokenAddress := helpers.GenerateAddressByModule(chain.chainName)
		bridgeDenom := fmt.Sprintf("%s%s", chain.chainName, tokenAddress)
		channelIBCHex := hex.EncodeToString([]byte(fmt.Sprintf("%s/%s", portID, channelID)))
		trace, err := fxtypes.GetIbcDenomTrace(bridgeDenom, channelIBCHex)
		suite.NoError(err)
		bridgeDenom = trace.IBCDenom()

		symbol := helpers.NewRandSymbol()
		metadata := banktypes.Metadata{
			Description: "The cross chain token of the Function X",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    bridgeDenom,
					Exponent: 0,
				}, {
					Denom:    symbol,
					Exponent: 18,
				},
			},
			Base:    bridgeDenom,
			Display: bridgeDenom,
			Name:    fmt.Sprintf("Token %s", symbol),
			Symbol:  symbol,
		}

		chain.AddBridgeTokenClaim(metadata.Name, metadata.Symbol, uint64(metadata.DenomUnits[1].Exponent), tokenAddress, channelIBCHex)
		suite.erc20.RegisterCoinProposal(metadata)
		suite.erc20.CheckRegisterCoin(metadata.Base)

		tokenPair := suite.erc20.TokenPair(metadata.Base)
		suite.Equal(tokenPair.Denom, metadata.Base)
		suite.Equal(tokenPair.Enabled, true)
		suite.Equal(tokenPair.ContractOwner, erc20types.OWNER_MODULE)
	}
}

func (suite *IntegrationTest) ERC20TokenERC20Test() {
	suite.Send(suite.erc20.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	proxy := suite.evm.DeployERC20Contract(suite.erc20.privKey)
	suite.evm.MintERC20(suite.erc20.privKey, proxy, common.BytesToAddress(suite.erc20.privKey.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18)))
	suite.True(suite.evm.CheckBalanceOf(proxy, common.BytesToAddress(suite.erc20.privKey.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))))

	metadataBrdige := fxtypes.GetCrossChainMetadataManyToOne("test token", helpers.NewRandSymbol(), uint32(18))
	aliases := make([]string, 0, len(suite.crosschain))
	bridgeTokens := make([]crosschaintypes.BridgeToken, 0, len(suite.crosschain))
	for _, chain := range suite.crosschain {
		denom, bridgeToken := chain.AddBridgeToken(metadataBrdige)
		aliases = append(aliases, denom)
		bridgeTokens = append(bridgeTokens, bridgeToken)
	}
	metadataBrdige.DenomUnits[0].Aliases = aliases
	suite.erc20.RegisterErc20Proposal(proxy.String(), aliases)

	symbol := suite.evm.Symbol(proxy)
	suite.erc20.CheckRegisterCoin(strings.ToLower(symbol))
	metadata := suite.GetMetadata(strings.ToLower(symbol))
	suite.T().Log("metadata", metadata.String())

	tokenPair := suite.erc20.TokenPair(metadata.Base)
	suite.Equal(tokenPair.Denom, metadata.Base)
	suite.Equal(tokenPair.Enabled, true)
	suite.Equal(tokenPair.ContractOwner, erc20types.OWNER_EXTERNAL)

	for i, chain := range suite.crosschain {
		bridgeToken := bridgeTokens[i]
		suite.erc20.ConvertERC20(tokenPair.GetERC20Contract(), sdkmath.NewInt(100), suite.erc20.AccAddress())
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(100)))

		suite.erc20.ConvertDenom(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(100)), chain.chainName)
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(100)))

		suite.erc20.ConvertERC20(tokenPair.GetERC20Contract(), sdkmath.NewInt(100), suite.erc20.AccAddress())
		// convert native token(metadata base) to erc20 token
		suite.erc20.ConvertCoin(suite.erc20.HexAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(100)))

		if i < len(suite.crosschain)-1 {
			// remove proposal
			suite.erc20.UpdateDenomAliasProposal(metadata.Base, bridgeToken.Denom)

			// check remove
			response, err := suite.erc20.ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: metadata.Base})
			suite.NoError(err)
			suite.Equal(len(suite.crosschain)-i-1, len(response.Aliases))

			_, err = suite.erc20.ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: bridgeToken.Denom})
			suite.Error(err)
		}
	}
	suite.erc20.ToggleTokenConversionProposal(metadata.Base)
	suite.False(suite.erc20.TokenPair(metadata.Base).Enabled)
}

func (suite *IntegrationTest) ERC20IBCChainTokenERC20Test() {
	suite.Send(suite.erc20.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	proxy := suite.evm.DeployERC20Contract(suite.erc20.privKey)
	portID := "transfer"
	channelID := "channel-0"

	aliases := make([]string, 0, len(suite.crosschain))
	for _, chain := range suite.crosschain {
		tokenAddress := helpers.GenerateAddressByModule(chain.chainName)
		bridgeDenom := fmt.Sprintf("%s%s", chain.chainName, tokenAddress)
		channelIBCHex := hex.EncodeToString([]byte(fmt.Sprintf("%s/%s", portID, channelID)))
		trace, err := fxtypes.GetIbcDenomTrace(bridgeDenom, channelIBCHex)
		suite.NoError(err)
		bridgeDenom = trace.IBCDenom()
		aliases = append(aliases, bridgeDenom)
		chain.AddBridgeTokenClaim("Test ERC20", "ERC20IBC", uint64(18), tokenAddress, channelIBCHex)
	}

	suite.erc20.RegisterErc20Proposal(proxy.String(), aliases)
	symbol := suite.evm.Symbol(proxy)
	suite.erc20.CheckRegisterCoin(strings.ToLower(symbol))
	metadata := suite.GetMetadata(strings.ToLower(symbol))

	tokenPair := suite.erc20.TokenPair(metadata.Base)
	suite.Equal(tokenPair.Denom, metadata.Base)
	suite.Equal(tokenPair.Enabled, true)
	suite.Equal(tokenPair.ContractOwner, erc20types.OWNER_EXTERNAL)
	suite.evm.MintERC20(suite.erc20.privKey, proxy, common.BytesToAddress(suite.erc20.privKey.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18)))
	suite.True(suite.evm.CheckBalanceOf(proxy, common.BytesToAddress(suite.erc20.privKey.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))))
}
