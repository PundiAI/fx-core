package keeper_test

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookCrossChainChain() {
	testCases := []struct {
		name   string
		relays func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string)
		error  func(args []string) string
		result bool
	}{
		{
			name: "ok - module",
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				moduleName := md.RandModule()
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(moduleName),
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(moduleName),
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
			name: "failed - from address insufficient funds",
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				moduleName := md.RandModule()
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(moduleName),
						Amount:    randMint,
						Fee:       big.NewInt(1),
						Target:    fxtypes.MustStrToByte32(moduleName),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, []string{
					fmt.Sprintf("%s%s", randMint.String(), pair.GetDenom()),
					fmt.Sprintf("%s%s", big.NewInt(0).Add(randMint, big.NewInt(1)).String(), pair.GetDenom()),
				}
			},
			error: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
			result: false,
		},
		{
			name: "failed - module insufficient funds",
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				addAmount := big.NewInt(0).Add(randMint, big.NewInt(1))
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(addAmount))))
				moduleName := md.RandModule()
				expectedModuleName := moduleName
				if moduleName == "gravity" {
					expectedModuleName = "eth"
				}

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(moduleName),
						Amount:    addAmount,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(moduleName),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, []string{
					fmt.Sprintf("%s%s", big.NewInt(0).Sub(addAmount, big.NewInt(1)).String(), md.GetDenom(expectedModuleName)),
					fmt.Sprintf("%s%s", addAmount.String(), md.GetDenom(expectedModuleName)),
				}
			},
			error: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
			result: false,
		},
		{
			name: "failed - target not support",
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))

				unknownChain := "chainabc"
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(unknownChain),
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(unknownChain),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, []string{"chainabc"}
			},
			error: func(args []string) string {
				return fmt.Sprintf("target %s not support: invalid target", args[0])
			},
			result: false,
		},
		{
			name: "failed - bridge token is not exist",
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))

				randDenom := fmt.Sprintf("t%st", strings.ToLower(tmrand.Str(5)))
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(randDenom, sdk.NewIntFromBigInt(randMint))))

				moduleName := md.RandModule()
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: suite.RandAddress(moduleName),
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(moduleName),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         randDenom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, []string{}
			},
			error: func(args []string) string {
				return "bridge token is not exist: invalid"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			md := suite.GenerateCrossChainDenoms()
			pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
			suite.NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdk.NewIntFromBigInt(randMint))
			// relay event
			relays, errArgs := tc.relays(*pair, md, signer.Address(), randMint)
			// hook transfer cross chain
			err = suite.app.Erc20Keeper.EVMHooks().HookTransferCrossChainEvent(suite.ctx, relays)
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
		relays func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string)
		error  func(args []string) string
		result bool
	}{
		{
			name: "ok - ibc token",
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				prefix := "px" // "evmos"

				recipient, _ := bech32.ConvertAndEncode(prefix, suite.RandSigner().AccAddress().Bytes())

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    randMint,
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
			relays: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				prefix := "px" // "evmos"
				sourcePort1, sourceChannel1 := suite.RandTransferChannel()
				recipient, _ := bech32.ConvertAndEncode(prefix, suite.RandSigner().AccAddress().Bytes())

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    randMint,
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
			md := suite.GenerateCrossChainDenoms(ibcToken)
			pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
			suite.NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdk.NewIntFromBigInt(randMint))
			// relay event
			relays, errArgs := tc.relays(*pair, md, signer.Address(), randMint, sourcePort, sourceChannel)
			// hook transfer cross chain
			err = suite.app.Erc20Keeper.EVMHooks().HookTransferCrossChainEvent(suite.ctx, relays)
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
