package types_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/pundiai/fx-core/v8/app"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

const (
	tronAddressErr = ": doesn't pass format validation: invalid address"
)

func TestMsgBondedOracle_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgBondedOracle{
				ChainName: helpers.NewRandSymbol(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty oracle address",
			msg: &types.MsgBondedOracle{
				ChainName:     moduleName,
				OracleAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid oracle address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix oracle address",
			msg: &types.MsgBondedOracle{
				ChainName:     moduleName,
				OracleAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty bridger address",
			msg: &types.MsgBondedOracle{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				ChainName:       moduleName,
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
				ChainName:       moduleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid external address: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "error - oracle address is same bridge address",
			msg: &types.MsgBondedOracle{
				ChainName:       moduleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin(helpers.NewRandDenom(), sdkmath.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "same address: invalid request",
		},
		{
			testName: "success - zero delegate amount",
			msg: &types.MsgBondedOracle{
				ChainName:       moduleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin(helpers.NewRandDenom(), sdkmath.NewInt(0)),
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
		{
			testName: "success",
			msg: &types.MsgBondedOracle{
				ChainName:       moduleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  types.NewDelegateAmount(sdkmath.NewInt(1)),
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgAddDelegate_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := helpers.NewRandDenom()
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgAddDelegate{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty oracle address",
			msg: &types.MsgAddDelegate{
				ChainName:     moduleName,
				OracleAddress: "",
				Amount:        types.NewDelegateAmount(sdkmath.NewInt(1)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid oracle address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address oracle address",
			msg: &types.MsgAddDelegate{
				ChainName:     moduleName,
				OracleAddress: errPrefixAddress,
				Amount:        types.NewDelegateAmount(sdkmath.NewInt(1)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - zero delegate amount",
			msg: &types.MsgAddDelegate{
				ChainName:     moduleName,
				OracleAddress: normalOracleAddress,
				Amount:        types.NewDelegateAmount(sdkmath.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgAddDelegate{
				ChainName:     moduleName,
				OracleAddress: normalOracleAddress,
				Amount:        types.NewDelegateAmount(sdkmath.NewInt(1)),
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgOracleSetConfirm_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)

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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgOracleSetConfirm{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty bridger address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:      moduleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - err address prefix bridger address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:      moduleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty external address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       moduleName,
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
				ChainName:       moduleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid external address: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - empty signature",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       moduleName,
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
				ChainName:       moduleName,
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
				ChainName:       moduleName,
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgOracleSetUpdatedClaim_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)

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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      moduleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address oracle address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      moduleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty members",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: externalAddressToUpper(normalExternalAddress),
					},
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid external address: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err event nonce",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgBridgeTokenClaim_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
	addressBytes := helpers.GenAccAddress()
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      moduleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address oracle address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      moduleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty tokenContract address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid token contract: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - invalid memo",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				Memo:           tmrand.Str(100),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "could not decode hex memo string: invalid request",
		},
		{
			testName: "err - empty name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				Memo:           "",
				Name:           "empty token name: invalid request",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "empty token symbol: invalid request",
		},
		{
			testName: "err - empty symbol",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				Memo:           "",
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
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				Memo:           "",
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
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				Memo:           "",
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
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				Memo:           hex.EncodeToString([]byte("transfer/channel-0")),
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
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  helpers.GenZeroExternalAddr(moduleName),
				Memo:           hex.EncodeToString([]byte("transfer/channel-0")),
				Name:           tmrand.Str(5),
				Symbol:         tmrand.Str(10),
				EventNonce:     uint64(tmrand.Int63n(1000000000)),
				BlockHeight:    uint64(tmrand.Int63n(1000000000)),
				Decimals:       uint64(tmrand.Int63n(18)),
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgSendToFxClaim_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgSendToFxClaim{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address bridger address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty sender address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid sender address: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - tokenContract address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid token contract: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - invalid receiver address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdkmath.Int{},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - negative amount",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdkmath.ZeroInt().Sub(sdkmath.NewInt(10000)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - invalid channel ibc",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdkmath.ZeroInt(),
				TargetIbc:      tmrand.Str(100),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "could not decode hex targetIbc: invalid request",
		},
		{
			testName: "err - zero event nonce",
			msg: &types.MsgSendToFxClaim{
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdkmath.ZeroInt(),
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
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdkmath.ZeroInt(),
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
				ChainName:      moduleName,
				BridgerAddress: normalBridgeAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalBridgeAddress,
				Amount:         sdkmath.ZeroInt(),
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgSendToExternal_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
	normalFxAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgSendToExternal{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - error prefix sender address",
			msg: &types.MsgSendToExternal{
				ChainName: moduleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid sender address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty dest address",
			msg: &types.MsgSendToExternal{
				ChainName: moduleName,
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
				ChainName: moduleName,
				Sender:    normalFxAddress,
				Dest:      externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid dest address: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - empty amount",
			msg: &types.MsgSendToExternal{
				ChainName: moduleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdkmath.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid amount: invalid request",
		},
		{
			testName: "err - bridge fee denom not match amount denom",
			msg: &types.MsgSendToExternal{
				ChainName: moduleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin(helpers.NewRandDenom(), sdkmath.NewInt(1)),
				BridgeFee: sdk.NewCoin(helpers.NewRandDenom(), sdkmath.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "bridge fee denom not equal amount denom: invalid request",
		},
		{
			testName: "err bridge fee",
			msg: &types.MsgSendToExternal{
				ChainName: moduleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdkmath.NewInt(1)),
				BridgeFee: sdk.NewCoin("demo", sdkmath.NewInt(0)),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "invalid bridge fee: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgSendToExternal{
				ChainName: moduleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdkmath.NewInt(1)),
				BridgeFee: sdk.NewCoin("demo", sdkmath.NewInt(1)),
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgSendToExternalClaim_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgSendToExternalClaim{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - error prefix bridger address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      moduleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty tokenContract address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid token contract: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - zero event nonce",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				ChainName:      moduleName,
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgConfirmBatch_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - invalid chain name",
			msg: &types.MsgConfirmBatch{
				ChainName: helpers.NewRandDenom(),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty bridge address",
			msg: &types.MsgConfirmBatch{
				ChainName:      moduleName,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			testName: "err - error prefix address bridger address",
			msg: &types.MsgConfirmBatch{
				ChainName:      moduleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			testName: "err - empty external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       moduleName,
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
				ChainName:       moduleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid external address: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - empty token contract address",
			msg: &types.MsgConfirmBatch{
				ChainName:       moduleName,
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
				ChainName:       moduleName,
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   externalAddressToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid token contract: mismatch expected: %s, got: %s: %s",
				normalExternalAddress, externalAddressToUpper(normalExternalAddress), sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err signature",
			msg: &types.MsgConfirmBatch{
				ChainName:       moduleName,
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
				ChainName:       moduleName,
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
				ChainName:       moduleName,
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestUpdateChainOraclesProposal_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalOracleAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
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
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "err - empty oracle",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   moduleName,
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
				ChainName:   moduleName,
				Title:       tmrand.Str(20),
				Description: tmrand.Str(20),
				Oracles: []string{
					strings.ToUpper(errPrefixAddress),
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason: fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: %s",
				sdk.Bech32MainPrefix, randomAddrPrefix, sdkerrors.ErrInvalidAddress.Error()),
		},
		{
			testName: "err - duplicate oracle",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   moduleName,
				Title:       "test title",
				Description: "test description",
				Oracles: []string{
					normalOracleAddress,
					normalOracleAddress,
				},
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("duplicate oracle address: %s: invalid address", normalOracleAddress),
		},
		{
			testName: "success",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   moduleName,
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
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				if moduleName == trontypes.ModuleName && strings.Contains(testCase.errReason, "mismatch expected") {
					testCase.errReason = strings.Split(testCase.errReason, ":")[0] + tronAddressErr
				}
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgBridgeCallConfirm_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	normalBridgeAddress := sdk.AccAddress(tmrand.Bytes(20)).String()
	normalExternalAddress := helpers.GenExternalAddr(moduleName)
	testCases := []struct {
		testName   string
		msg        *types.MsgBridgeCallConfirm
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err - empty chain name",
			msg: &types.MsgBridgeCallConfirm{
				ChainName: "",
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidRequest,
			errReason:  "unrecognized crosschain name: invalid request",
		},
		{
			testName: "success",
			msg: &types.MsgBridgeCallConfirm{
				Nonce:           uint64(tmrand.Int63n(100000)),
				BridgerAddress:  normalBridgeAddress,
				ExternalAddress: normalExternalAddress,
				Signature:       hex.EncodeToString(tmrand.Bytes(100)),
				ChainName:       moduleName,
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err := testCase.msg.ValidateBasic()
			if testCase.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.err, "%+v", testCase.msg)
				require.EqualValuesf(t, testCase.errReason, err.Error(), "%+v", testCase.msg)
			}
		})
	}
}

func TestMsgBridgeCallClaim_ValidateBasic(t *testing.T) {
	moduleName := getRandModule()
	tests := []struct {
		name          string
		msg           *types.MsgBridgeCallClaim
		expectedError string
	}{
		{
			name: "success",
			msg: &types.MsgBridgeCallClaim{
				ChainName:      moduleName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(moduleName),
				Refund:         helpers.GenExternalAddr(moduleName),
				TokenContracts: []string{helpers.GenExternalAddr(moduleName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(100)},
				To:             helpers.GenExternalAddr(moduleName),
				Data:           "",
				QuoteId:        sdkmath.NewInt(1),
				GasLimit:       sdkmath.ZeroInt(),
				Memo:           "",
				TxOrigin:       helpers.GenExternalAddr(moduleName),
			},
		},
		{
			name: "error - empty bridger address",
			msg: &types.MsgBridgeCallClaim{
				ChainName:      moduleName,
				BridgerAddress: "",
				EventNonce:     1,
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(moduleName),
				Refund:         helpers.GenExternalAddr(moduleName),
				TokenContracts: []string{helpers.GenExternalAddr(moduleName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(100)},
				To:             helpers.GenExternalAddr(moduleName),
				Data:           "",
				QuoteId:        sdkmath.NewInt(1),
				GasLimit:       sdkmath.ZeroInt(),
				Memo:           "",
				TxOrigin:       helpers.GenExternalAddr(moduleName),
			},
			expectedError: "invalid bridger address: empty address string is not allowed: invalid address",
		},
		{
			name: "error - mismatched tokenContracts and amounts",
			msg: &types.MsgBridgeCallClaim{
				ChainName:      moduleName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(moduleName),
				Refund:         helpers.GenExternalAddr(moduleName),
				TokenContracts: []string{helpers.GenExternalAddr(moduleName), helpers.GenExternalAddr(moduleName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(100)}, // Only one amount provided
				To:             helpers.GenExternalAddr(moduleName),
				Data:           "",
				QuoteId:        sdkmath.NewInt(1),
				GasLimit:       sdkmath.ZeroInt(),
				Memo:           "",
				TxOrigin:       helpers.GenExternalAddr(moduleName),
			},
			expectedError: "mismatched token contracts and amounts: invalid request",
		},
		{
			name: "error - zero event nonce",
			msg: &types.MsgBridgeCallClaim{
				ChainName:      moduleName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     0, // Zero event nonce
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(moduleName),
				Refund:         helpers.GenExternalAddr(moduleName),
				TokenContracts: []string{helpers.GenExternalAddr(moduleName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(100)},
				To:             helpers.GenExternalAddr(moduleName),
				Data:           "",
				QuoteId:        sdkmath.NewInt(1),
				GasLimit:       sdkmath.ZeroInt(),
				Memo:           "",
				TxOrigin:       helpers.GenExternalAddr(moduleName),
			},
			expectedError: "zero event nonce: invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if len(tt.expectedError) > 0 {
				require.Error(t, err, "ValidateBasic()")
				assert.Equal(t, tt.expectedError, err.Error(), "ValidateBasic()")
			} else {
				require.NoError(t, err, "ValidateBasic()")
			}
		})
	}
}

// externalAddressToUpper for test case address to upper
func externalAddressToUpper(address string) string {
	if strings.HasPrefix(address, "0x") {
		result := fmt.Sprintf("%s%s", address[0:2], strings.ToLower(address[2:]))
		if result == address {
			result = fmt.Sprintf("%s%s", address[0:2], strings.ToUpper(address[2:]))
		}
		return result
	} else if strings.HasPrefix(address, "T") {
		return fmt.Sprintf("%s%s", address[0:1], strings.ToLower(address[1:]))
	}
	panic(fmt.Sprintf("not support address prefix: %s", address))
}

func getRandModule() string {
	modules := fxtypes.GetSupportChains()
	return modules[tmrand.Intn(len(modules))]
}
