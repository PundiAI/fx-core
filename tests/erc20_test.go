package tests

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *IntegrationTest) ERC20Test() {
	suite.Send(suite.erc20.AccAddress(), suite.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))

	decimals := 18
	metadata := fxtypes.GetCrossChainMetadata("test token", helpers.NewRandSymbol(), uint32(decimals))

	var aliases []string
	var bridgeTokens []crosschaintypes.BridgeToken
	for _, chain := range suite.crosschain {
		bridgeTokenAddr := helpers.GenerateAddress().Hex()
		if chain.chainName == trontypes.ModuleName {
			bridgeTokenAddr = trontypes.AddressFromHex(bridgeTokenAddr)
		}
		chain.AddBridgeTokenClaim(metadata.Name, metadata.Symbol, uint64(decimals), bridgeTokenAddr, "")
		bridgeTokenDenom := chain.GetBridgeDenomByToken(bridgeTokenAddr)
		aliases = append(aliases, bridgeTokenDenom)
		bridgeTokens = append(bridgeTokens, crosschaintypes.BridgeToken{
			Token: bridgeTokenAddr,
			Denom: bridgeTokenDenom,
		})
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
		chain.SendToFxClaim(bridgeToken.Token, sdk.NewInt(200), fxtypes.LegacyERC20Target)
		balance := suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), chain.HexAddress())
		suite.Equal(balance, big.NewInt(200))

		balance = suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), suite.erc20.HexAddress())
		suite.erc20.TransferERC20(chain.privKey, tokenPair.GetERC20Contract(), suite.erc20.HexAddress(), big.NewInt(100))
		suite.Equal(big.NewInt(0).Add(balance, big.NewInt(100)), suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), suite.erc20.HexAddress()))

		receive := suite.erc20.HexAddress().String()
		if chain.chainName == trontypes.ModuleName {
			receive = trontypes.AddressFromHex(receive)
		}
		suite.erc20.TransferCrossChain(suite.erc20.privKey, tokenPair.GetERC20Contract(), receive,
			big.NewInt(50), big.NewInt(50), fmt.Sprintf("chain/%s", chain.chainName))

		resp, err := chain.CrosschainQuery().GetPendingSendToExternal(suite.ctx,
			&crosschaintypes.QueryPendingSendToExternalRequest{
				ChainName:     chain.chainName,
				SenderAddress: suite.erc20.AccAddress().String(),
			})
		suite.NoError(err)
		suite.Equal(1, len(resp.UnbatchedTransfers))
		suite.Equal(int64(50), resp.UnbatchedTransfers[0].Token.Amount.Int64())
		suite.Equal(int64(50), resp.UnbatchedTransfers[0].Fee.Amount.Int64())
		suite.Equal(suite.erc20.AccAddress().String(), resp.UnbatchedTransfers[0].Sender)
		if chain.chainName == trontypes.ModuleName {
			suite.Equal(trontypes.AddressFromHex(suite.erc20.HexAddress().String()), resp.UnbatchedTransfers[0].DestAddress)
		} else {
			suite.Equal(suite.erc20.HexAddress().String(), resp.UnbatchedTransfers[0].DestAddress)
		}

		// covert chain.address erc20 token to native token: metadata.base
		suite.erc20.ConvertERC20(chain.privKey, tokenPair.GetERC20Contract(), sdk.NewInt(50), suite.erc20.AccAddress())
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdk.NewInt(50)))

		// covert erc20.addr metadata.base
		suite.erc20.ConvertDenom(suite.erc20.privKey, suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdk.NewInt(50)), chain.chainName)
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(bridgeToken.Denom, sdk.NewInt(50)))

		// send to chain.address
		baseTokenBalanceAmount := suite.QueryBalances(chain.AccAddress()).AmountOf(metadata.Base)
		chain.SendToFxClaim(bridgeToken.Token, sdk.NewInt(100), "")
		suite.CheckBalance(chain.AccAddress(), sdk.NewCoin(metadata.Base, baseTokenBalanceAmount.Add(sdk.NewInt(100))))

		// convert native token(metadata base) to erc20 token
		balance = suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), suite.erc20.HexAddress())
		suite.erc20.ConvertCoin(chain.privKey, suite.erc20.HexAddress(), sdk.NewCoin(metadata.Base, sdk.NewInt(50)))
		suite.Equal(suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), suite.erc20.HexAddress()), big.NewInt(0).Add(balance, big.NewInt(50)))

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
