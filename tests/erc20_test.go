package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *IntegrationTest) ERC20Test() {
	suite.Send(suite.erc20.AccAddress(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))

	decimals := 18
	metadata := fxtypes.GetCrossChainMetadata("test token", helpers.NewRandSymbol(), uint32(decimals))

	var aliases []string
	var bridgeTokens []crosschaintypes.BridgeToken
	for _, chain := range suite.crosschain {
		bridgeTokenAddr := helpers.GenerateAddressByModule(chain.chainName)
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
		chain.SendToFxClaim(bridgeToken.Token, sdkmath.NewInt(200), fxtypes.LegacyERC20Target)
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
		suite.erc20.ConvertERC20(chain.privKey, tokenPair.GetERC20Contract(), sdkmath.NewInt(50), suite.erc20.AccAddress())
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(50)))

		// covert erc20.addr metadata.base
		suite.erc20.ConvertDenom(suite.erc20.privKey, suite.erc20.AccAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(50)), chain.chainName)
		suite.CheckBalance(suite.erc20.AccAddress(), sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(50)))

		// send to chain.address
		baseTokenBalanceAmount := suite.QueryBalances(chain.AccAddress()).AmountOf(metadata.Base)
		chain.SendToFxClaim(bridgeToken.Token, sdkmath.NewInt(100), "")
		suite.CheckBalance(chain.AccAddress(), sdk.NewCoin(metadata.Base, baseTokenBalanceAmount.Add(sdkmath.NewInt(100))))

		// convert native token(metadata base) to erc20 token
		balance = suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), suite.erc20.HexAddress())
		suite.erc20.ConvertCoin(chain.privKey, suite.erc20.HexAddress(), sdk.NewCoin(metadata.Base, sdkmath.NewInt(50)))
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

//gocyclo:ignore
func (suite *IntegrationTest) ERC20IBCChainTokenTest() {
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

		chain.AddBridgeTokenClaim(metadata.Name, metadata.Symbol, uint64(18), tokenAddress, channelIBCHex)
		suite.erc20.RegisterCoinProposal(metadata)
		suite.erc20.CheckRegisterCoin(metadata.Base)

		tokenPair := suite.erc20.TokenPair(metadata.Base)
		suite.Equal(tokenPair.Denom, metadata.Base)
		suite.Equal(tokenPair.Enabled, true)
		suite.Equal(tokenPair.ContractOwner, erc20types.OWNER_MODULE)

		balance := suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), chain.HexAddress())
		chain.SendToFxClaim(tokenAddress, sdkmath.NewInt(200), fxtypes.LegacyERC20Target)
		suite.Equal(big.NewInt(0).Add(balance, big.NewInt(200)), suite.erc20.BalanceOf(tokenPair.GetERC20Contract(), chain.HexAddress()))

		// ibc token transfer to chain
		receive := suite.erc20.HexAddress().String()
		if chain.chainName == trontypes.ModuleName {
			receive = trontypes.AddressFromHex(receive)
		}
		suite.erc20.TransferCrossChain(chain.privKey, tokenPair.GetERC20Contract(), receive,
			big.NewInt(50), big.NewInt(50), fmt.Sprintf("chain/%s", chain.chainName))

		resp, err := chain.CrosschainQuery().GetPendingSendToExternal(suite.ctx,
			&crosschaintypes.QueryPendingSendToExternalRequest{
				ChainName:     chain.chainName,
				SenderAddress: chain.AccAddress().String(),
			})
		suite.NoError(err)
		suite.Equal(1, len(resp.UnbatchedTransfers))
		suite.Equal(int64(50), resp.UnbatchedTransfers[0].Token.Amount.Int64())
		suite.Equal(int64(50), resp.UnbatchedTransfers[0].Fee.Amount.Int64())
		suite.Equal(tokenAddress, resp.UnbatchedTransfers[0].Token.Contract)
		suite.Equal(chain.AccAddress().String(), resp.UnbatchedTransfers[0].Sender)
		if chain.chainName == trontypes.ModuleName {
			suite.Equal(trontypes.AddressFromHex(suite.erc20.HexAddress().String()), resp.UnbatchedTransfers[0].DestAddress)
		} else {
			suite.Equal(suite.erc20.HexAddress().String(), resp.UnbatchedTransfers[0].DestAddress)
		}

		// ibc token transfer to other cosmos chain
		receive, err = sdk.Bech32ifyAddressBytes("px", suite.erc20.AccAddress().Bytes())
		suite.Require().NoError(err)

		respTX := suite.erc20.TransferCrossChain(chain.privKey, tokenPair.GetERC20Contract(), receive, big.NewInt(50), big.NewInt(0), "ibc/0/px")

		// "send_packet.packet_src_channel='channel-0' AND send_packet.packet_sequence='1'"
		search, err := suite.NodeClient().TxSearch("send_packet.packet_src_channel='channel-0'", 1, 100, "")
		suite.NoError(err)
		for _, tx := range search.Txs {
			find := false
			for _, event := range tx.TxResult.Events {
				if event.Type == "ethereum_tx" {
					for _, attr := range event.Attributes {
						if string(attr.Key) == "ethereumTxHash" {
							if string(attr.Value) == respTX.Hash().String() {
								find = true
							}
						}
					}
				}
			}
			if find {
				for _, event := range tx.TxResult.Events {
					if event.Type == "relay_transfer_cross_chain" {
						for _, attr := range event.Attributes {
							if string(attr.Key) == "from" {
								suite.Equal(string(attr.Value), chain.HexAddress().String())
							}
							if string(attr.Key) == "recipient" {
								suite.Equal(string(attr.Value), receive)
							}
							if string(attr.Key) == "amount" {
								suite.Equal(string(attr.Value), "50")
							}
							if string(attr.Key) == "fee" {
								suite.Equal(string(attr.Value), "0")
							}
							if string(attr.Key) == "target" {
								suite.Equal(string(attr.Value), "ibc/0/px")
							}
							if string(attr.Key) == "token_address" {
								suite.Equal(string(attr.Value), tokenPair.GetErc20Address())
							}
							if string(attr.Key) == "coin" {
								suite.Equal(string(attr.Value), metadata.Base)
							}
						}
					}
				}
			}
		}
	}
}
