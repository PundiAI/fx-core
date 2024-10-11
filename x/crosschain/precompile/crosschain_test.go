package precompile_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testcontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func TestCrossChainABI(t *testing.T) {
	crossChain := precompile.NewCrossChainMethod(nil)

	require.Equal(t, 6, len(crossChain.Method.Inputs))
	require.Equal(t, 1, len(crossChain.Method.Outputs))

	require.Equal(t, 8, len(crossChain.Event.Inputs))
}

//nolint:gocyclo
func (suite *PrecompileTestSuite) TestCrossChainIBC() {
	testCases := []struct {
		name     string
		malleate func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - origin token",
			malleate: func(_ *types.TokenPair, _ Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string) {
				suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint)))

				sourcePort, sourceChannel := suite.RandTransferChannel()

				prefix, recipient := suite.RandPrefixAndAddress()
				data, err := crosschaintypes.GetABI().Pack(
					"crossChain",
					common.Address{},
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(fmt.Sprintf("%s/%s/%s", prefix, sourcePort, sourceChannel)),
					"ibc memo",
				)
				suite.Require().NoError(err)

				return data, randMint, sourcePort, sourceChannel, nil
			},
			result: true,
		},
		{
			name: "failed - not zero fee",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string) {
				sourcePort, sourceChannel := suite.RandTransferChannel()

				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				suite.MintToken(signer.AccAddress(), coin)
				_, err := suite.App.Erc20Keeper.ConvertCoin(suite.Ctx,
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), randMint)

				prefix, recipient := suite.RandPrefixAndAddress()

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := crosschaintypes.GetABI().Pack(
					"crossChain",
					pair.GetERC20Contract(),
					recipient,
					amount,
					fee,
					fxtypes.MustStrToByte32(fmt.Sprintf("%s/%s/%s", prefix, sourcePort, sourceChannel)),
					"ibc memo",
				)
				suite.Require().NoError(err)

				return data, big.NewInt(0), sourcePort, sourceChannel, []string{sdk.Coin{Amount: sdkmath.NewIntFromBigInt(fee), Denom: md.metadata.Base}.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("ibc transfer fee must be zero: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - not zero fee - origin token",
			malleate: func(_ *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string) {
				suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint)))

				sourcePort, sourceChannel := suite.RandTransferChannel()

				prefix, recipient := suite.RandPrefixAndAddress()
				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := crosschaintypes.GetABI().Pack(
					"crossChain",
					common.Address{},
					recipient,
					amount,
					fee,
					fxtypes.MustStrToByte32(fmt.Sprintf("%s/%s/%s", prefix, sourcePort, sourceChannel)),
					"ibc memo",
				)
				suite.Require().NoError(err)

				return data, randMint, sourcePort, sourceChannel, []string{sdk.Coin{Amount: sdkmath.NewIntFromBigInt(fee), Denom: fxtypes.DefaultDenom}.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("ibc transfer fee must be zero: %s", args[0])
			},
			result: false,
		},
		{
			name: "contract - ok - origin token",
			malleate: func(_ *types.TokenPair, _ Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string) {
				suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint)))

				sourcePort, sourceChannel := suite.RandTransferChannel()

				prefix, recipient := suite.RandPrefixAndAddress()
				data, err := contract.MustABIJson(testcontract.CrossChainTestMetaData.ABI).Pack(
					"crossChain",
					common.Address{},
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(fmt.Sprintf("%s/%s/%s", prefix, sourcePort, sourceChannel)),
					"ibc memo",
				)
				suite.Require().NoError(err)

				return data, randMint, sourcePort, sourceChannel, nil
			},
			result: true,
		},
		{
			name: "contract - failed - not zero fee",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string) {
				sourcePort, sourceChannel := suite.RandTransferChannel()

				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				suite.MintToken(signer.AccAddress(), coin)
				_, err := suite.App.Erc20Keeper.ConvertCoin(suite.Ctx,
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), suite.crosschain, randMint)

				prefix, recipient := suite.RandPrefixAndAddress()
				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := contract.MustABIJson(testcontract.CrossChainTestMetaData.ABI).Pack(
					"crossChain",
					pair.GetERC20Contract(),
					recipient,
					amount,
					fee,
					fxtypes.MustStrToByte32(fmt.Sprintf("%s/%s/%s", prefix, sourcePort, sourceChannel)),
					"ibc memo",
				)
				suite.Require().NoError(err)

				return data, big.NewInt(0), sourcePort, sourceChannel, []string{sdk.Coin{Amount: sdkmath.NewIntFromBigInt(fee), Denom: md.metadata.Base}.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("ibc transfer fee must be zero: %s", args[0])
			},
			result: false,
		},
		{
			name: "contract - failed - not zero fee - origin token",
			malleate: func(_ *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *big.Int, string, string, []string) {
				suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint)))

				sourcePort, sourceChannel := suite.RandTransferChannel()

				prefix, recipient := suite.RandPrefixAndAddress()
				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := contract.MustABIJson(testcontract.CrossChainTestMetaData.ABI).Pack(
					"crossChain",
					common.Address{},
					recipient,
					amount,
					fee,
					fxtypes.MustStrToByte32(fmt.Sprintf("%s/%s/%s", prefix, sourcePort, sourceChannel)),
					"ibc memo",
				)
				suite.Require().NoError(err)

				return data, randMint, sourcePort, sourceChannel, []string{sdk.Coin{Amount: sdkmath.NewIntFromBigInt(fee), Denom: fxtypes.DefaultDenom}.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("ibc transfer fee must be zero: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			signer := suite.RandSigner()
			// token pair
			md := suite.GenerateCrossChainDenoms()
			pair, err := suite.App.Erc20Keeper.RegisterNativeCoin(suite.Ctx, md.GetMetadata())
			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			packData, value, portId, channelId, errArgs := tc.malleate(pair, md, signer, randMint)

			crosschainContract := crosschaintypes.GetAddress()
			addrQuery := signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				crosschainContract = suite.crosschain
				addrQuery = suite.crosschain
			}

			commitments := suite.App.IBCKeeper.ChannelKeeper.GetAllPacketCommitmentsAtChannel(suite.Ctx, portId, channelId)
			ibcTxs := make(map[string]bool, len(commitments))
			for _, commitment := range commitments {
				ibcTxs[fmt.Sprintf("%s/%s/%d", commitment.PortId, commitment.ChannelId, commitment.Sequence)] = true
			}

			totalBefore, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			res := suite.EthereumTx(signer, crosschainContract, value, packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, sdk.AccAddress(addrQuery.Bytes()))
				suite.Require().True(chainBalances.IsZero(), chainBalances.String())
				balance := suite.BalanceOf(pair.GetERC20Contract(), addrQuery)
				suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

				manyToOne := make(map[string]bool)
				suite.App.BankKeeper.IterateAllDenomMetaData(suite.Ctx, func(md banktypes.Metadata) bool {
					if len(md.DenomUnits) > 0 && len(md.DenomUnits[0].Aliases) > 0 {
						manyToOne[md.Base] = true
					}
					return false
				})
				totalAfter, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)

				for _, coin := range totalBefore.Supply {
					if manyToOne[coin.Denom] {
						continue
					}
					expect := totalAfter.Supply.AmountOf(coin.Denom)
					if strings.HasPrefix(coin.GetDenom(), "ibc/") {
						expect = expect.Add(sdkmath.NewIntFromBigInt(randMint))
					}
					suite.Require().Equal(coin.Amount.String(), expect.String(), coin.Denom)
				}

				for _, event := range suite.Ctx.EventManager().Events() {
					if event.Type != ibcchanneltypes.EventTypeSendPacket {
						continue
					}
					var eventPortId, eventChannelId string
					var sequence string
					var data []byte

					for _, attr := range event.Attributes {
						attrKey, attrValue := attr.Key, attr.Value
						if attrKey == ibcchanneltypes.AttributeKeyDataHex {
							data, err = hex.DecodeString(attrValue)
							suite.Require().NoError(err)
						}
						if attrKey == ibcchanneltypes.AttributeKeySequence {
							sequence = attrValue
						}
						if attrKey == ibcchanneltypes.AttributeKeySrcPort {
							eventPortId = attrValue
						}
						if attrKey == ibcchanneltypes.AttributeKeySrcChannel {
							eventChannelId = attrValue
						}
					}
					if eventPortId != portId || eventChannelId != channelId {
						continue
					}
					txKey := fmt.Sprintf("%s/%s/%s", portId, channelId, sequence)
					if ibcTxs[txKey] {
						continue
					}
					var packet ibctransfertypes.FungibleTokenPacketData
					err = suite.App.LegacyAmino().UnmarshalJSON(data, &packet)
					suite.Require().NoError(err)
					suite.Require().Equal(sdk.AccAddress(addrQuery.Bytes()).String(), packet.Sender)
					suite.Require().Equal(randMint.String(), packet.Amount)
				}
			} else {
				suite.Error(res, errors.New(tc.error(errArgs)))
			}
		})
	}
}
