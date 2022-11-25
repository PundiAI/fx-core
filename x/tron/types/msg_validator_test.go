package types_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"

	_ "github.com/functionx/fx-core/v3/app"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestMsgBondedOracle_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	normalOrchestratorAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgBondedOracle
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgBondedOracle{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err chain name",
			msg: &types.MsgBondedOracle{
				ChainName: "111",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err oracle address",
			msg: &types.MsgBondedOracle{
				ChainName:     trontypes.ModuleName,
				OracleAddress: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "oracle address", types.ErrInvalid),
		},
		{
			testName: "err oracle address",
			msg: &types.MsgBondedOracle{
				ChainName:     trontypes.ModuleName,
				OracleAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "oracle address", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgBondedOracle{
				ChainName:      trontypes.ModuleName,
				OracleAddress:  normalOracleAddress,
				BridgerAddress: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgBondedOracle{
				ChainName:      trontypes.ModuleName,
				OracleAddress:  normalOracleAddress,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalOrchestratorAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalOrchestratorAddress,
				ExternalAddress: "err external address",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err deposit amount",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalOrchestratorAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "delegate amount", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgBondedOracle{
				ChainName:       trontypes.ModuleName,
				OracleAddress:   normalOracleAddress,
				BridgerAddress:  normalOrchestratorAddress,
				ExternalAddress: normalExternalAddress,
				DelegateAmount:  sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1)),
			},
			expectPass: true,
			err:        nil,
			errReason:  "",
		},
	}

	for _, testCase := range testCases {
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgAddDelegate_ValidateBasic(t *testing.T) {
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgAddDelegate
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgAddDelegate{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err oracle address",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: errPrefixAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1)},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "oracle address", types.ErrInvalid),
		},
		{
			testName: "err amount",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: normalOracleAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(0)},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "amount", types.ErrInvalid),
		},
		{
			testName: "err amount",
			msg: &types.MsgAddDelegate{
				ChainName:     trontypes.ModuleName,
				OracleAddress: normalOracleAddress,
				Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(-1)},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "amount", types.ErrInvalid),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgOracleSetConfirm_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgOracleSetConfirm
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgOracleSetConfirm{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: strings.ToLower(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err signature",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				Signature:       "",
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("signature: %s", types.ErrEmpty),
		},
		{
			testName: "err signature: hex.decode error",
			msg: &types.MsgOracleSetConfirm{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				Signature:       "kkkkk",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not hex decode signature: %s: %s", "kkkkk", types.ErrInvalid),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgOracleSetUpdatedClaim_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgOracleSetUpdatedClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err members",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Members:        nil,
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("%s: %s", "members", types.ErrEmpty),
		},
		{
			testName: "err members: member external error",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: "",
					},
				},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err members: member external power",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           0,
						ExternalAddress: normalExternalAddress,
					},
				},
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("%s: %s", "member power", types.ErrEmpty),
		},
		{
			testName: "err members: member external error",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: strings.ToLower(normalExternalAddress),
					},
				},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err event nonce",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: normalExternalAddress,
					},
				},
				EventNonce: 0,
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "event nonce", types.ErrUnknown),
		},
		{
			testName: "err block height",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
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
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "block height", types.ErrUnknown),
		},
		{
			testName: "success",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgBridgeTokenClaim_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	zeroTronAddress := make([]byte, 0)
	zeroTronAddress = append(zeroTronAddress, tronAddress.TronBytePrefix)
	zeroTronAddress = append(zeroTronAddress, gethCommon.HexToAddress("0x0000000000000000000000000000000000000000").Bytes()...)
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgBridgeTokenClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err channelIBC: not hex.decode channelIBC",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "kkkkk",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not decode hex channelIbc string: %s: %s", "kkkkk", types.ErrInvalid),
		},
		{
			testName: "err name",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "",
				Name:           "",
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("token name: %s", types.ErrEmpty),
		},
		{
			testName: "err symbol",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				ChannelIbc:     "",
				Name:           "DEMO TOKEN",
				Symbol:         "",
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("token symbol: %s", types.ErrEmpty),
		},
		{
			testName: "err event nonce",
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
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "event nonce", types.ErrUnknown),
		},
		{
			testName: "err block height",
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
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "block height", types.ErrUnknown),
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
				TokenContract:  tronAddress.Address(zeroTronAddress).String(),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgSendToFxClaim_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSendToFxClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgSendToFxClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err sender address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         "",
				TokenContract:  "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "sender address", types.ErrInvalid),
		},
		{
			testName: "err sender address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "sender address", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err receiver address",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "receiver address", types.ErrInvalid),
		},
		{
			testName: "err amount",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalFxAddress,
				Amount:         sdk.Int{},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "amount cannot be negative", types.ErrInvalid),
		},
		{
			testName: "err amount",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalFxAddress,
				Amount:         sdk.ZeroInt().Sub(sdk.NewInt(10000)),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "amount cannot be negative", types.ErrInvalid),
		},
		{
			testName: "err channelIBC: not hex.decode channelIBC",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalFxAddress,
				Amount:         sdk.ZeroInt(),
				TargetIbc:      "kkkkk",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not decode hex targetIbc string: %s: %s", "kkkkk", types.ErrInvalid),
		},
		{
			testName: "err event nonce",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalFxAddress,
				Amount:         sdk.ZeroInt(),
				TargetIbc:      "",
				EventNonce:     0,
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "event nonce", types.ErrUnknown),
		},
		{
			testName: "err block height",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalFxAddress,
				Amount:         sdk.ZeroInt(),
				EventNonce:     1,
				BlockHeight:    0,
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "block height", types.ErrUnknown),
		},
		{
			testName: "success",
			msg: &types.MsgSendToFxClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				Sender:         normalExternalAddress,
				TokenContract:  normalExternalAddress,
				Receiver:       normalFxAddress,
				Amount:         sdk.ZeroInt(),
				TargetIbc:      hex.EncodeToString([]byte("bc/transfer/channel-0")),
				EventNonce:     1,
				BlockHeight:    1,
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgSendToExternal_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSendToExternal
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgSendToExternal{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err sender address",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "sender address", types.ErrInvalid),
		},
		{
			testName: "err dest address",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "dest", types.ErrInvalid),
		},
		{
			testName: "err dest address",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "dest", types.ErrInvalid),
		},
		{
			testName: "err amount",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "amount", types.ErrInvalid),
		},
		{
			testName: "err bridge fee: amount coin name != bridgeFee coin name",
			msg: &types.MsgSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo1", sdk.NewInt(1)),
				BridgeFee: sdk.NewCoin("demo2", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("fee and amount must be the same type %s != %s: %s", "demo1", "demo2", types.ErrInvalid),
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
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridge fee", types.ErrInvalid),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgCancelSendToExternal_ValidateBasic(t *testing.T) {
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgCancelSendToExternal
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgCancelSendToExternal{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err sender address",
			msg: &types.MsgCancelSendToExternal{
				ChainName: trontypes.ModuleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "sender address", types.ErrInvalid),
		},
		{
			testName: "err transaction id",
			msg: &types.MsgCancelSendToExternal{
				ChainName:     trontypes.ModuleName,
				Sender:        normalFxAddress,
				TransactionId: 0,
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "transaction id", types.ErrUnknown),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgSendToExternalClaim_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSendToExternalClaim
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgSendToExternalClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err event nonce",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				EventNonce:     0,
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "event nonce", types.ErrUnknown),
		},
		{
			testName: "err block height",
			msg: &types.MsgSendToExternalClaim{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: normalFxAddress,
				TokenContract:  normalExternalAddress,
				EventNonce:     1,
				BlockHeight:    0,
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "block height", types.ErrUnknown),
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
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("%s: %s", "batch nonce", types.ErrUnknown),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgRequestBatch_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalFxAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgRequestBatch
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgRequestBatch{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err sender address",
			msg: &types.MsgRequestBatch{
				ChainName: trontypes.ModuleName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "sender address", types.ErrInvalid),
		},
		{
			testName: "err denom",
			msg: &types.MsgRequestBatch{
				ChainName: trontypes.ModuleName,
				Sender:    normalFxAddress,
				Denom:     "",
			},
			expectPass: false,
			err:        types.ErrUnknown,
			errReason:  fmt.Sprintf("denom: %s", types.ErrUnknown),
		},
		{
			testName: "err tokenContract address",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(-1),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("minimum fee: %s", types.ErrInvalid),
		},
		{
			testName: "err fee receive",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "fee receive address", types.ErrInvalid),
		},
		{
			testName: "err fee receive",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "fee receive address", types.ErrInvalid),
		},

		{
			testName: "success",
			msg: &types.MsgRequestBatch{
				ChainName:  trontypes.ModuleName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: normalExternalAddress,
				BaseFee:    sdk.ZeroInt(),
			},
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestMsgConfirmBatch_ValidateBasic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := tronAddress.PubkeyToAddress(key.PublicKey).String()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgConfirmBatch
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.MsgConfirmBatch{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err bridger address",
			msg: &types.MsgConfirmBatch{
				ChainName:      trontypes.ModuleName,
				BridgerAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "bridger address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: strings.ToLower(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "external address", types.ErrInvalid),
		},
		{
			testName: "err token contract address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   "",
				Nonce:           0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err token contract address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   strings.ToUpper(normalExternalAddress),
				Nonce:           0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err external address",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   strings.ToLower(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "token contract", types.ErrInvalid),
		},
		{
			testName: "err signature",
			msg: &types.MsgConfirmBatch{
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   normalExternalAddress,
				Signature:       "",
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("signature: %s", types.ErrEmpty),
		},
		{
			testName: "err signature: hex.decode error",
			msg: &types.MsgConfirmBatch{
				Nonce:           0,
				ChainName:       trontypes.ModuleName,
				BridgerAddress:  normalOracleAddress,
				ExternalAddress: normalExternalAddress,
				TokenContract:   normalExternalAddress,
				Signature:       "gggg",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not hex decode signature: %s: %s", "gggg", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgConfirmBatch{
				Nonce:           0,
				BridgerAddress:  normalOracleAddress,
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}

func TestUpdateChainOraclesProposal_ValidateBasic(t *testing.T) {
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.UpdateChainOraclesProposal
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name",
			msg: &types.UpdateChainOraclesProposal{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "chain name", types.ErrInvalid),
		},
		{
			testName: "err oracle",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   trontypes.ModuleName,
				Title:       "test title",
				Description: "test description",
			},
			expectPass: false,
			err:        types.ErrEmpty,
			errReason:  fmt.Sprintf("%s: %s", "oracles", types.ErrEmpty),
		},
		{
			testName: "err external address",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   trontypes.ModuleName,
				Title:       "test title",
				Description: "test description",
				Oracles: []string{
					strings.ToUpper(errPrefixAddress),
				},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "oracle address", types.ErrInvalid),
		},
		{
			testName: "err oracle",
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
			errReason:  fmt.Sprintf("oracle address: %s", types.ErrDuplicate),
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
		err = testCase.msg.ValidateBasic()
		if testCase.expectPass {
			require.NoError(t, err)
		} else {
			require.NotNil(t, err, testCase.testName)
			require.ErrorIs(t, err, testCase.err, testCase.testName)
			require.EqualValues(t, testCase.errReason, err.Error(), testCase.testName)
		}
	}
}
