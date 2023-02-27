package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookCrossChainChain() {
	testCases := []struct {
		name     string
		malleate func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - module",
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				moduleName := md.RandModule()
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: helpers.GenerateAddressByModule(moduleName),
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
			name: "ok - bsc ibc token",
			malleate: func(_ types.TokenPair, _ Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				sourcePort, sourceChannel := suite.RandTransferChannel()
				tokenAddress := helpers.GenerateAddress()
				denom, err := suite.CrossChainKeepers()[bsctypes.ModuleName].SetIbcDenomTrace(suite.ctx, tokenAddress.Hex(), hex.EncodeToString([]byte(fmt.Sprintf("%s/%s", sourcePort, sourceChannel))))
				suite.Require().NoError(err)
				suite.CrossChainKeepers()[bsctypes.ModuleName].AddBridgeToken(suite.ctx, tokenAddress.Hex(), denom)

				symbol := helpers.NewRandSymbol()
				ibcMD := banktypes.Metadata{
					Description: "The cross chain token of the Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    denom,
							Exponent: 0,
						},
						{
							Denom:    symbol,
							Exponent: 18,
						},
					},
					Base:    denom,
					Display: denom,
					Name:    fmt.Sprintf("%s Token", symbol),
					Symbol:  symbol,
				}
				pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, ibcMD)
				suite.Require().NoError(err)
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(randMint))))

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: helpers.GenerateAddressByModule(bsctypes.ModuleName),
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32("chain/" + bsctypes.ModuleName),
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
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				moduleName := md.RandModule()
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: helpers.GenerateAddressByModule(moduleName),
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
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
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
						Recipient: helpers.GenerateAddressByModule(moduleName),
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
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))

				unknownChain := "chainabc"
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: helpers.GenerateAddressByModule(unknownChain),
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
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))

				randDenom := helpers.NewRandDenom()
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(randDenom, sdk.NewIntFromBigInt(randMint))))

				moduleName := md.RandModule()
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: helpers.GenerateAddressByModule(moduleName),
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
			relays, errArgs := tc.malleate(*pair, md, signer.Address(), randMint)
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
		name     string
		malleate func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - ibc token",
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				prefix, recipient := suite.RandPrefixAndAddress()

				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)),
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
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				// add relay token
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				sourcePort1, sourceChannel1 := suite.RandTransferChannel()
				prefix, recipient := suite.RandPrefixAndAddress()

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
		{
			name: "failed - no zero fee",
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				prefix, recipient := suite.RandPrefixAndAddress()

				fee := big.NewInt(int64(tmrand.Intn(1000) + 10))
				relayAmount := big.NewInt(0).Sub(randMint, fee)

				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    relayAmount,
						Fee:       fee,
						Target:    fxtypes.MustStrToByte32(ibcTarget),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransferCrossChain{relay}, []string{fee.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("ibc transfer fee must be zero: %s: invalid coins", args[0])
			},
			result: false,
		},
		{
			name: "failed - invalid recipient address - hex",
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				prefix, recipient := "0x", suite.RandSigner().AccAddress().String()
				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(ibcTarget),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}

				return []types.RelayTransferCrossChain{relay}, []string{recipient, "wrong length: invalid address"}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid to address %s, error %s", args[0], args[1])
			},
			result: false,
		},
		{
			name: "failed - invalid recipient address - bench32",
			malleate: func(pair types.TokenPair, md Metadata, singerAddr common.Address, randMint *big.Int, sourcePort, sourceChannel string) ([]types.RelayTransferCrossChain, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, singerAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdk.NewIntFromBigInt(randMint))))
				prefix, recipient := "px", helpers.GenerateAddress().Hex()
				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)
				relay := types.RelayTransferCrossChain{
					TransferCrossChainEvent: &types.TransferCrossChainEvent{
						From:      singerAddr,
						Recipient: recipient,
						Amount:    randMint,
						Fee:       big.NewInt(0),
						Target:    fxtypes.MustStrToByte32(ibcTarget),
					},
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.GetDenom(),
					ContractOwner: pair.ContractOwner,
				}

				return []types.RelayTransferCrossChain{relay}, []string{recipient, "decoding bech32 failed: string not all lowercase or all uppercase: invalid address"}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid to address %s, error %s", args[0], args[1])
			},
			result: false,
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
			randMint := big.NewInt(int64(tmrand.Uint32() + 100000))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdk.NewIntFromBigInt(randMint))
			// relay event
			relays, errArgs := tc.malleate(*pair, md, signer.Address(), randMint, sourcePort, sourceChannel)
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
