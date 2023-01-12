package types_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gethcommon "github.com/ethereum/go-ethereum/common"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	_ "github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestMsgBondedOracle_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))

	externalAddressLenth := trontypes.TronContractAddressLen
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgBondedOracle
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgBondedOracle{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgBondedOracle{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty oracle address",
			msg: &types.MsgBondedOracle{
				ChainName:     trontypes.ModuleName,
				OracleAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid oracle address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix oracle address",
			msg: &types.MsgBondedOracle{
				ChainName:     trontypes.ModuleName,
				OracleAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty bridger address",
			msg: &types.MsgBondedOracle{
				ChainName:      trontypes.ModuleName,
				OracleAddress:  normalOracleAddress,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - invalid bridger address",
			msg: &types.MsgBondedOracle{
				ChainName:      trontypes.ModuleName,
				OracleAddress:  normalOracleAddress,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty external address",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid external address: empty: invalid address",
		},
		{
			testName: "err - invalid external address",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: "error external address",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid external address: invalid address (%s) of the wrong length exp (%d) actual (%d): invalid address", "error external address", externalAddressLenth, len("error external address")),
		},
		{
			testName: "error - oracle address is same bridge address",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin(fmt.Sprintf("a%sb", tmrand.Str(5)), sdk.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "same address: invalid request",
		},
		{
			testName: "success - zero delegate amount",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin(fmt.Sprintf("a%sb", tmrand.Str(5)), sdk.NewInt(0)),
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
		{
			testName: "success",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1)),
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgAddDelegate_ValidateBasic(t *testing.T) {
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5)))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))

	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgAddDelegate
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgAddDelegate{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgAddDelegate{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty oracle address",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: "",
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1)},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid oracle address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address oracle address",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: errPrefixAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1)},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - zero delegate amount",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: normalOracleAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(0)},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - negative delegate amount",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: normalOracleAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(-1)},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: normalOracleAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1)},
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgOracleSetConfirm_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))

	testCases := []struct {
		testName   string
		msg        *types.MsgOracleSetConfirm
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgOracleSetConfirm{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgOracleSetConfirm{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty bridger address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - err address prefix bridger address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty external address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid external address: empty: invalid address",
		},
		{
			testName: "err - error external address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid external address: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - empty signature",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				Signature:       "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty signature: invalid request",
		},
		{
			testName: "err - signature: hex.decode error",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				Signature:       tmrand.Str(100),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "could not hex decode signature: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           0,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				Signature:       hex.EncodeToString([]byte("kkkkk")),
				ChainName:       trontypes.ModuleName,
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgOracleSetUpdatedClaim_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))

	testCases := []struct {
		testName   string
		msg        *types.MsgOracleSetUpdatedClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address oracle address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty members",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members:        nil,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty members: invalid request",
		},
		{
			testName: "err - empty member external address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: "",
					},
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid external address: empty: invalid address",
		},
		{
			testName: "err - zero member external power",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           0,
						ExternalAddress: normalExternalAddress,
					},
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero power: invalid request",
		},
		{
			testName: "err - invalid member external error",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: strings.ToLower(normalExternalAddress),
					},
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid external address: invalid address: %s: invalid address", strings.ToLower(normalExternalAddress)),
		},
		{
			testName: "err event nonce",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: normalExternalAddress,
					},
				},
				EventNonce: 0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero event nonce: invalid request",
		},
		{
			testName: "err block height",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: normalExternalAddress,
					},
				},
				EventNonce:  1,
				BlockHeight: 0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero block height: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: normalExternalAddress,
					},
				},
				EventNonce:     1,
				BlockHeight:    1,
				OracleSetNonce: 0,
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgBridgeTokenClaim_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	zeroTronAddress := make([]byte, 0)
	zeroTronAddress = append(zeroTronAddress, tronaddress.TronBytePrefix)
	zeroTronAddress = append(zeroTronAddress, gethcommon.HexToAddress("0x0000000000000000000000000000000000000000").Bytes()...)
	normalFxAddress := addressBytes.String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgBridgeTokenClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address oracle address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty tokenContract address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid token contract: empty: invalid address",
		},
		{
			testName: "err - invalid tokenContract address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid token contract: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - invalid channelIBC",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     tmrand.Str(100),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "could not decode hex channelIbc string: invalid request",
		},
		{
			testName: "err - empty name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "",
				Name:           "empty token name: invalid request",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty token symbol: invalid request",
		},
		{
			testName: "err - empty symbol",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "",
				Name:           "DEMO TOKEN",
				Symbol:         "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty token symbol: invalid request",
		},
		{
			testName: "err - zero event nonce",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "",
				Name:           "DEMO TOKEN",
				Symbol:         "DEMO",
				EventNonce:     0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero event nonce: invalid request",
		},
		{
			testName: "err - zero block height",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "",
				Name:           "DEMO TOKEN",
				Symbol:         "DEMO",
				EventNonce:     1,
				BlockHeight:    0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero block height: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				Name:           "DEMO TOKEN",
				Symbol:         "DEMO",
				EventNonce:     1,
				BlockHeight:    1,
				Decimals:       0,
			},
			expectPass: true,
		},
		{
			testName: "success-0x0000000000000000000000000000000000000000",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  tronaddress.Address(zeroTronAddress).String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				Name:           "TRX",
				Symbol:         "TRX",
				EventNonce:     1,
				BlockHeight:    1,
				Decimals:       0,
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgSendToFxClaim_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSendToFxClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgSendToFxClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgSendToFxClaim{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address bridger address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty sender address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid sender address: empty: invalid address",
		},
		{
			testName: "err - invalid sender address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid sender address: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - tokenContract address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid token contract: empty: invalid address",
		},
		{
			testName: "err - invalid tokenContract address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid token contract: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - invalid receiver address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid receiver address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty amount",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdk.Int{},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - negative amount",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdk.ZeroInt().Sub(sdk.NewInt(10000)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - invalid channel ibc",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdk.ZeroInt(),
				TargetIbc:      tmrand.Str(100),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "could not decode hex targetIbc: invalid request",
		},
		{
			testName: "err - zero event nonce",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdk.ZeroInt(),
				TargetIbc:      "",
				EventNonce:     0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero event nonce: invalid request",
		},
		{
			testName: "err block height",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdk.ZeroInt(),
				EventNonce:     1,
				BlockHeight:    0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero block height: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdk.ZeroInt(),
				TargetIbc:      hex.EncodeToString([]byte("bc/transfer/channel-0")),
				EventNonce:     1,
				BlockHeight:    1,
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgSendToExternal_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalFxAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSendToExternal
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgSendToExternal{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgSendToExternal{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - error prefix sender address",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid sender address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty dest address",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid dest address: empty: invalid address",
		},
		{
			testName: "err - invalid dest address",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid dest address: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - empty amount",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - bridge fee denom not match amount denom",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin(fmt.Sprintf("a%sb", strings.ToLower(tmrand.Str(4))), sdk.NewInt(1)),
				BridgeFee: sdk.NewCoin(fmt.Sprintf("a%sb", strings.ToLower(tmrand.Str(5))), sdk.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "bridge fee denom not equal amount denom: invalid request",
		},
		{
			testName: "err bridge fee",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdk.NewInt(1)),
				BridgeFee: sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid bridge fee: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdk.NewInt(1)),
				BridgeFee: sdk.NewCoin("demo", sdk.NewInt(1)),
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgCancelSendToExternal_ValidateBasic(t *testing.T) {
	normalFxAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgCancelSendToExternal
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgCancelSendToExternal{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid sender address",
			msg: &types.MsgCancelSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid sender address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - zero transaction id",
			msg: &types.MsgCancelSendToExternal{
				ChainName:     trontypes.ModuleName,
				Sender:        normalFxAddress,
				TransactionId: 0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero transaction id: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgCancelSendToExternal{
				ChainName:     trontypes.ModuleName,
				Sender:        normalFxAddress,
				TransactionId: 1,
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgSendToExternalClaim_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalFxAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSendToExternalClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgSendToExternalClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgSendToExternalClaim{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - error prefix bridger address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty tokenContract address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid token contract: empty: invalid address",
		},
		{
			testName: "err - invalid tokenContract address toUpper",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid token contract: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - zero event nonce",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				EventNonce:     0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero event nonce: invalid request",
		},
		{
			testName: "err - zero block height",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				EventNonce:     1,
				BlockHeight:    0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero block height: invalid request",
		},
		{
			testName: "err batch nonce",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				EventNonce:     1,
				BlockHeight:    1,
				BatchNonce:     0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "zero batch nonce: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				EventNonce:     1,
				BlockHeight:    1,
				BatchNonce:     1,
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgRequestBatch_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgRequestBatch
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgRequestBatch{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgRequestBatch{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - error prefix sender address",
			msg: &types.MsgRequestBatch{
				ChainName: trontypes.ModuleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid sender address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty denom",
			msg: &types.MsgRequestBatch{
				ChainName: trontypes.ModuleName,
				Sender:    normalBridgeAddress,
				Denom:     "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty denom: invalid request",
		},
		{
			testName: "err - invalid minimum fee - negative",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalBridgeAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(-1),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid minimum fee: invalid request",
		},
		{
			testName: "err fee receive",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalBridgeAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid fee receive address: empty: invalid address",
		},
		{
			testName: "err fee receive",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalBridgeAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid fee receive address: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},

		{
			testName: "success",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalBridgeAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: normalExternalAddress,
				BaseFee:    sdk.ZeroInt(),
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestMsgConfirmBatch_ValidateBasic(t *testing.T) {
	normalExternalAddress := helpers.GenerateAddressByModule(trontypes.ModuleName)
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)

	testCases := []struct {
		testName   string
		msg        *types.MsgConfirmBatch
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgConfirmBatch{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid chain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgConfirmBatch{
				ChainName: strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5))),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized cross chain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgConfirmBatch{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address bridger address",
			msg: &types.MsgConfirmBatch{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid external address: empty: invalid address",
		},
		{
			testName: "err - invalid external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid external address: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err - empty token contract address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   "",
				Nonce:           0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid token contract: empty: invalid address",
		},
		{
			testName: "err - invalid token contract address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   strings.ToUpper(normalExternalAddress),
				Nonce:           0,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid token contract: invalid address: %s: invalid address", strings.ToUpper(normalExternalAddress)),
		},
		{
			testName: "err signature",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   normalExternalAddress,
				Signature:       "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty signature: invalid request",
		},
		{
			testName: "err signature: hex.decode error",
			msg: &types.MsgConfirmBatch{
				Nonce:           0,
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   normalExternalAddress,
				Signature:       tmrand.Str(100),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "could not hex decode signature: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgConfirmBatch{
				Nonce:           0,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   normalExternalAddress,
				Signature:       hex.EncodeToString([]byte("abcd")),
				ChainName:       trontypes.ModuleName,
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}

func TestUpdateChainOraclesProposal_ValidateBasic(t *testing.T) {
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	var err error

	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.UpdateChainOraclesProposal
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.UpdateChainOraclesProposal{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  fmt.Sprintf("invalid chain name: %s", sdkerrors.ErrInvalidRequest),
		},
		{
			testName: "err - empty oracle",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   trontypes.ModuleName,
				Title:       tmrand.Str(20),
				Description: tmrand.Str(20),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  fmt.Sprintf("empty oracles: %s", sdkerrors.ErrInvalidRequest),
		},
		{
			testName: "err external address",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   trontypes.ModuleName,
				Title:       tmrand.Str(20),
				Description: tmrand.Str(20),
				Oracles: []string{
					strings.ToUpper(errPrefixAddress),
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: %s",
				sdk.Bech32MainPrefix, randomAddrPrefix, sdkerrors.ErrInvalidAddress),
		},
		{
			testName: "err - duplicate oracle",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   trontypes.ModuleName,
				Title:       "test title",
				Description: "test description",
				Oracles: []string{
					normalOracleAddress,
					normalOracleAddress,
				},
			},
			expectPass: false,
			err:        types.ErrDuplicate,
			errReason:  fmt.Sprintf("oracle address: %s: %s", normalOracleAddress, types.ErrDuplicate),
		},
		{
			testName: "success",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   trontypes.ModuleName,
				Title:       "test title",
				Description: "test description",
				Oracles: []string{
					normalOracleAddress,
				},
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err = testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err, testCase.testName)
				require.Error(t, err, testCase.err, testCase.testName)
				require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
			}
		})
	}
}
