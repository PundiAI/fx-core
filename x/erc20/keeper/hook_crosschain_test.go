package keeper_test

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookCrossChainChain() {
	pairTokenFn := func() (types.TokenPair, banktypes.Metadata, *big.Int) {
		denoms := suite.GenerateCrossChainDenoms()
		pair, md := suite.DeployNativeRelayToken("TEST", denoms...)
		totalMint := suite.MintLockNativeTokenToModule(md, sdk.NewIntFromBigInt(big.NewInt(int64(tmrand.Uint32()+1))))
		return pair, md, totalMint
	}

	testCases := []struct {
		name      string
		pairToken func() (types.TokenPair, banktypes.Metadata, *big.Int)
		relays    func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransferCrossChain, []string)
		error     func(args []string) string
		result    bool
	}{
		{
			name:      "ok - chain/module",
			pairToken: pairTokenFn,
			relays: func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(totalCanMint))))
				relayAmount := big.NewInt(0).Div(totalCanMint, big.NewInt(int64(len(md.GetDenomUnits()[0].GetAliases()))))

				modules, _ := moduleDenom(md.GetDenomUnits()[0].GetAliases(), suite.CrossChainKeepers())

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(modules[0]),
						Amount:    relayAmount,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(fmt.Sprintf("chain/%s", modules[0])),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, nil
			},
			result: true,
		},
		{
			name:      "ok - module",
			pairToken: pairTokenFn,
			relays: func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(totalCanMint))))
				relayAmount := big.NewInt(0).Div(totalCanMint, big.NewInt(int64(len(md.GetDenomUnits()[0].GetAliases()))))
				modules, _ := moduleDenom(md.GetDenomUnits()[0].GetAliases(), suite.CrossChainKeepers())
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(modules[0]),
						Amount:    relayAmount,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(modules[0]),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, nil
			},
			result: true,
		},
		{
			name:      "ok - all",
			pairToken: pairTokenFn,
			relays: func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(totalCanMint))))
				relayAmount := big.NewInt(0).Div(totalCanMint, big.NewInt(int64(len(md.GetDenomUnits()[0].GetAliases()))))
				modules, _ := moduleDenom(md.GetDenomUnits()[0].GetAliases(), suite.CrossChainKeepers())
				relays := make([]types.RelayTransferCrossChain, 0, len(modules))
				for _, m := range modules {
					relay := types.RelayTransferCrossChain{
						TransferCrossChainEvent: &types.TransferCrossChainEvent{
							From:      singerAddr,
							Recipient: suite.RandAddress(m),
							Amount:    relayAmount,
							Fee:       big.NewInt(0),
							Target:    fxtypes.MustStrToByte32(m),
						},
						TokenContract: pair.GetERC20Contract(),
						Denom:         pair.GetDenom(),
						ContractOwner: pair.ContractOwner,
					}
					relays = append(relays, relay)
				}

				return relays, nil
			},
			result: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			pair, md, totalMint := tc.pairToken()
			// relay event
			relays, errArgs := tc.relays(pair, md, signer.Address(), totalMint)
			// hook transfer cross chain
			err := suite.app.Erc20Keeper.EVMHooks().HookTransferCrossChainEvent(suite.ctx, relays)
			// check result
			if tc.result {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHookCrossChainIBC() {
	testCases := []struct {
		name   string
		relays func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string)
		error  func(args []string) string
		result bool
	}{
		{
			name: "ok - ibc token",
			relays: func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(totalCanMint))))
				prefix := "px" // "evmos"

				recipient, _ := bech32.ConvertAndEncode(prefix, suite.RandSigner().AccAddress().Bytes())
				relayAmount := big.NewInt(0).Div(totalCanMint, big.NewInt(int64(len(md.GetDenomUnits()[0].GetAliases()))))

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    relayAmount,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s/%s", prefix, sourcePort, sourceChannel)),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, nil
			},
			result: true,
		},
		{
			name: "ok - base token",
			relays: func(pair types.TokenPair, md banktypes.Metadata, singerAddr common.Address, totalCanMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(totalCanMint))))
				prefix := "px" // "evmos"
				sourcePort1, sourceChannel1 := suite.RandTransferChannel()

				recipient, _ := bech32.ConvertAndEncode(prefix, suite.RandSigner().AccAddress().Bytes())
				relayAmount := big.NewInt(0).Div(totalCanMint, big.NewInt(int64(len(md.GetDenomUnits()[0].GetAliases()))))

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    relayAmount,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s/%s", prefix, sourcePort1, sourceChannel1)),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, nil
			},
			result: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// set port channel
			sourcePort, sourceChannel := suite.RandTransferChannel()
			// add ibc token
			ibcToken := suite.AddIBCToken(sourcePort, sourceChannel)
			// token pair
			denoms := suite.GenerateCrossChainDenoms()
			pair, md := suite.DeployNativeRelayToken("TEST", append(denoms, ibcToken)...)
			// mint and lock token
			totalMint := suite.MintLockNativeTokenToModule(md, sdk.NewIntFromBigInt(big.NewInt(int64(tmrand.Uint32()+1))))
			// relay event
			relays, errArgs := tc.relays(pair, md, signer.Address(), totalMint, sourcePort, sourceChannel)
			// hook transfer cross chain
			err := suite.app.Erc20Keeper.EVMHooks().HookTransferCrossChainEvent(suite.ctx, relays)
			// check result
			if tc.result {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func moduleDenom(denom []string, module map[string]CrossChainKeeper) ([]string, []string) {
	modules, denoms := make([]string, 0, len(module)), make([]string, 0, len(denom))
	for m := range module {
		for _, d := range denom {
			if strings.HasPrefix(d, m) {
				modules = append(modules, m)
				denoms = append(denoms, d)
			}
		}
	}
	return modules, denoms
}
