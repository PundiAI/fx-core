package types_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/x/crosschain/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const (
	chainName    = "demo"
	depositDenom = "FX"
)

var (
	depositAmount = sdk.NewInt(1)
)

func init() {
	types.InitMsgValidatorBasicRouter()
	types.RegisterValidatorBasic(chainName, types.EthereumMsgValidateBasic{})
}

func TestMsgSetOrchestrator(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	normalOrchestratorAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgSetOrchestratorAddress
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name - empty",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err chain name - 111",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName: "111",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "111", types.ErrInvalidChainName),
		},
		{
			testName: "err oracle address - empty",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName: chainName,
				Oracle:    "",
			},
			expectPass: false,
			err:        types.ErrOracleAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrOracleAddress),
		},
		{
			testName: "err oracle address - err prefix",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName: chainName,
				Oracle:    errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOracleAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOracleAddress),
		},
		{
			testName: "err orchestrator address - empty",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName:    chainName,
				Oracle:       normalOracleAddress,
				Orchestrator: "",
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrOrchestratorAddress),
		},
		{
			testName: "err orchestrator address - err prefix",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName:    chainName,
				Oracle:       normalOracleAddress,
				Orchestrator: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err externalAddress address - empty",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName:       chainName,
				Oracle:          normalOracleAddress,
				Orchestrator:    normalOrchestratorAddress,
				ExternalAddress: "",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err externalAddress address - err external address",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName:       chainName,
				Oracle:          normalOracleAddress,
				Orchestrator:    normalOrchestratorAddress,
				ExternalAddress: "err external address",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "err external address", types.ErrExternalAddress),
		},
		{
			testName: "err deposit amount - amount is not positive",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName:       chainName,
				Oracle:          normalOracleAddress,
				Orchestrator:    normalOrchestratorAddress,
				ExternalAddress: normalExternalAddress,
				Deposit:         sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        types.ErrInvalidCoin,
			errReason:  fmt.Sprintf("%s: %s", "0demo", types.ErrInvalidCoin),
		},
		{
			testName: "success",
			msg: &types.MsgSetOrchestratorAddress{
				ChainName:       chainName,
				Oracle:          normalOracleAddress,
				Orchestrator:    normalOrchestratorAddress,
				ExternalAddress: normalExternalAddress,
				Deposit:         sdk.NewCoin(depositDenom, depositAmount),
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

func TestMsgAddOracleDeposit(t *testing.T) {
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	normalOracleAddress := addressBytes.String()
	var err error
	errPrefixAddress, err := bech32.ConvertAndEncode("demo", addressBytes)
	require.NoError(t, err)
	testCases := []struct {
		testName   string
		msg        *types.MsgAddOracleDeposit
		expectPass bool
		err        error
		errReason  string
	}{
		{
			testName: "err chain name - empty",
			msg: &types.MsgAddOracleDeposit{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err oracle address - err prefix",
			msg: &types.MsgAddOracleDeposit{
				ChainName: chainName,
				Oracle:    errPrefixAddress,
				Amount:    sdk.Coin{Denom: depositDenom, Amount: depositAmount},
			},
			expectPass: false,
			err:        types.ErrOracleAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOracleAddress),
		},
		{
			testName: "err amount - value:0",
			msg: &types.MsgAddOracleDeposit{
				ChainName: chainName,
				Oracle:    normalOracleAddress,
				Amount:    sdk.Coin{Denom: depositDenom, Amount: sdk.NewInt(0)},
			},
			expectPass: false,
			err:        types.ErrInvalidCoin,
			errReason:  fmt.Sprintf("0%s: %s", depositDenom, types.ErrInvalidCoin),
		},
		{
			testName: "err amount - value:-1",
			msg: &types.MsgAddOracleDeposit{
				ChainName: chainName,
				Oracle:    normalOracleAddress,
				Amount:    sdk.Coin{Denom: depositDenom, Amount: sdk.NewInt(-1)},
			},
			expectPass: false,
			err:        types.ErrInvalidCoin,
			errReason:  fmt.Sprintf("-1%s: %s", depositDenom, types.ErrInvalidCoin),
		},
		{
			testName: "success",
			msg: &types.MsgAddOracleDeposit{
				ChainName: chainName,
				Oracle:    normalOracleAddress,
				Amount:    sdk.Coin{Denom: depositDenom, Amount: depositAmount},
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

func TestMsgOracleSetConfirm(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgOracleSetConfirm{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err orchestrator address",
			msg: &types.MsgOracleSetConfirm{
				ChainName:           chainName,
				OrchestratorAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err external address: empty",
			msg: &types.MsgOracleSetConfirm{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     "",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err external address: ToUpper",
			msg: &types.MsgOracleSetConfirm{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err external address: ToLower",
			msg: &types.MsgOracleSetConfirm{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     strings.ToLower(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToLower(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err signature: empty",
			msg: &types.MsgOracleSetConfirm{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				Signature:           "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("signature is empty: %s", types.ErrInvalid),
		},
		{
			testName: "err signature: hex.decode error",
			msg: &types.MsgOracleSetConfirm{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				Signature:           "kkkkk",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not hex decode signature: %s: %s", "kkkkk", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgOracleSetConfirm{
				Nonce:               0,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				Signature:           hex.EncodeToString([]byte("kkkkk")),
				ChainName:           chainName,
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

func TestMsgOracleSetUpdatedClaim(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err orchestrator address - err prefix address",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err members: members len == 0",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
				Members:      nil,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "members len == 0", types.ErrInvalid),
		},
		{
			testName: "err members: member external error: empty",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: "",
					},
				},
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err members: member external power is 0 ",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           0,
						ExternalAddress: normalExternalAddress,
					},
				},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "member power == 0", types.ErrInvalid),
		},
		{
			testName: "err members: member external error: case not match",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: strings.ToLower(normalExternalAddress),
					},
				},
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToLower(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err event nonce: event nonce == 0",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
				Members: types.BridgeValidators{
					{
						Power:           1,
						ExternalAddress: normalExternalAddress,
					},
				},
				EventNonce: 0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "event nonce == 0", types.ErrInvalid),
		},
		{
			testName: "err block height: block height == 0",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
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
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "block height == 0", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgOracleSetUpdatedClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
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

func TestMsgBridgeTokenClaim(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
	addressBytes := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
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
			testName: "err chain name - empty",
			msg: &types.MsgBridgeTokenClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err orchestrator address - err prefix address",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:    chainName,
				Orchestrator: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err tokenContract address - empty",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: "",
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrTokenContractAddress),
		},
		{
			testName: "err tokenContract address - ToUpper",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrTokenContractAddress),
		},
		{
			testName: "err channelIBC: not hex.decode channelIBC",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				ChannelIbc:    "kkkkk",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not decode hex channelIbc string: %s: %s", "kkkkk", types.ErrInvalid),
		},
		{
			testName: "err name: empty",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				ChannelIbc:    "",
				Name:          "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("token name is empty: %s", types.ErrInvalid),
		},
		{
			testName: "err symbol: empty",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				ChannelIbc:    "",
				Name:          "DEMO TOKEN",
				Symbol:        "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("token symbol is empty: %s", types.ErrInvalid),
		},
		{
			testName: "err block height: event nonce == 0",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				ChannelIbc:    "",
				Name:          "DEMO TOKEN",
				Symbol:        "DEMO",
				EventNonce:    0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "event nonce == 0", types.ErrInvalid),
		},
		{
			testName: "err block height: block height == 0",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				ChannelIbc:    "",
				Name:          "DEMO TOKEN",
				Symbol:        "DEMO",
				EventNonce:    1,
				BlockHeight:   0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "block height == 0", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgBridgeTokenClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				ChannelIbc:    hex.EncodeToString([]byte("transfer/channel-0")),
				Name:          "DEMO TOKEN",
				Symbol:        "DEMO",
				EventNonce:    1,
				BlockHeight:   1,
				Decimals:      0,
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

func TestMsgSendToFxClaim(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgSendToFxClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err orchestrator address - err prefix address",
			msg: &types.MsgSendToFxClaim{
				ChainName:    chainName,
				Orchestrator: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err sender address - empty",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        "",
				TokenContract: "",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err sender address - ToUpper",
			msg: &types.MsgSendToFxClaim{
				ChainName:    chainName,
				Orchestrator: normalFxAddress,
				Sender:       strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err tokenContract address - empty",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: "",
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrTokenContractAddress),
		},
		{
			testName: "err tokenContract address - ToUpper",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrTokenContractAddress),
		},
		{
			testName: "err receiver address - err prefix address",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: normalExternalAddress,
				Receiver:      errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, sdkerrors.ErrInvalidAddress),
		},
		{
			testName: "err channelIBC: not hex.decode channelIBC",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: normalExternalAddress,
				Receiver:      normalFxAddress,
				TargetIbc:     "kkkkk",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not decode hex targetIbc string: %s: %s", "kkkkk", types.ErrInvalid),
		},
		{
			testName: "err event nonce: event nonce == 0",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: normalExternalAddress,
				Receiver:      normalFxAddress,
				TargetIbc:     "",
				EventNonce:    0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "event nonce == 0", types.ErrInvalid),
		},
		{
			testName: "err block height: block height == 0",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: normalExternalAddress,
				Receiver:      normalFxAddress,
				EventNonce:    1,
				BlockHeight:   0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "block height == 0", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgSendToFxClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				Sender:        normalExternalAddress,
				TokenContract: normalExternalAddress,
				Receiver:      normalFxAddress,
				TargetIbc:     hex.EncodeToString([]byte("bc/transfer/channel-0")),
				EventNonce:    1,
				BlockHeight:   1,
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

func TestMsgSendToExternal(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgSendToExternal{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err sender address - err prefix address",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, sdkerrors.ErrInvalidAddress),
		},
		{
			testName: "err dest address - empty",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
				Sender:    normalFxAddress,
				Dest:      "",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err dest address - ToUpper",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
				Sender:    normalFxAddress,
				Dest:      strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err amount: amount is zero",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        types.ErrInvalidCoin,
			errReason:  fmt.Sprintf("%s: %s", "0demo", types.ErrInvalidCoin),
		},
		{
			testName: "err bridgeFee: amount coin name != bridgeFee coin name",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
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
			testName: "err bridgeFee: fee is zero",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
				Sender:    normalFxAddress,
				Dest:      normalExternalAddress,
				Amount:    sdk.NewCoin("demo", sdk.NewInt(1)),
				BridgeFee: sdk.NewCoin("demo", sdk.NewInt(0)),
			},
			expectPass: false,
			err:        types.ErrInvalidCoin,
			errReason:  fmt.Sprintf("%s: %s", "0demo", types.ErrInvalidCoin),
		},
		{
			testName: "success",
			msg: &types.MsgSendToExternal{
				ChainName: chainName,
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

func TestMsgCancelSendToExternal(t *testing.T) {
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
			testName: "err chain name - empty",
			msg: &types.MsgCancelSendToExternal{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err sender address - err prefix address",
			msg: &types.MsgCancelSendToExternal{
				ChainName: chainName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, sdkerrors.ErrInvalidAddress),
		},
		{
			testName: "err transactionId - transactionId == 0",
			msg: &types.MsgCancelSendToExternal{
				ChainName:     chainName,
				Sender:        normalFxAddress,
				TransactionId: 0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "transaction id == 0", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgCancelSendToExternal{
				ChainName:     chainName,
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

func TestMsgSendToExternalClaim(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgSendToExternalClaim{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err orchestrator address - err prefix address",
			msg: &types.MsgSendToExternalClaim{
				ChainName:    chainName,
				Orchestrator: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err tokenContract address - empty",
			msg: &types.MsgSendToExternalClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: "",
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrTokenContractAddress),
		},
		{
			testName: "err tokenContract address - ToUpper",
			msg: &types.MsgSendToExternalClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrTokenContractAddress),
		},
		{
			testName: "err event nonce: event nonce == 0",
			msg: &types.MsgSendToExternalClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				EventNonce:    0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "event nonce == 0", types.ErrInvalid),
		},
		{
			testName: "err block height: block height == 0",
			msg: &types.MsgSendToExternalClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				EventNonce:    1,
				BlockHeight:   0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "block height == 0", types.ErrInvalid),
		},
		{
			testName: "err batch nonce : batch nonce == 0",
			msg: &types.MsgSendToExternalClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				EventNonce:    1,
				BlockHeight:   1,
				BatchNonce:    0,
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "batch nonce == 0", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgSendToExternalClaim{
				ChainName:     chainName,
				Orchestrator:  normalFxAddress,
				TokenContract: normalExternalAddress,
				EventNonce:    1,
				BlockHeight:   1,
				BatchNonce:    1,
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

func TestMsgRequestBatch(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgRequestBatch{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err sender address - err prefix address",
			msg: &types.MsgRequestBatch{
				ChainName: chainName,
				Sender:    errPrefixAddress,
			},
			expectPass: false,
			err:        sdkerrors.ErrInvalidAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, sdkerrors.ErrInvalidAddress),
		},
		{
			testName: "err denom - empty",
			msg: &types.MsgRequestBatch{
				ChainName: chainName,
				Sender:    normalFxAddress,
				Denom:     "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("denom is empty:%s: %s", "", types.ErrInvalid),
		},
		{
			testName: "err tokenContract address - ToUpper",
			msg: &types.MsgRequestBatch{
				ChainName:  chainName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(-1),
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("minimum fee is positive:%s: %s", sdk.NewInt(-1).String(), types.ErrInvalid),
		},
		{
			testName: "err fee receive: empty",
			msg: &types.MsgRequestBatch{
				ChainName:  chainName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: "",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err fee receive: ToUpper",
			msg: &types.MsgRequestBatch{
				ChainName:  chainName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrExternalAddress),
		},

		{
			testName: "success",
			msg: &types.MsgRequestBatch{
				ChainName:  chainName,
				Sender:     normalFxAddress,
				Denom:      "demo",
				MinimumFee: sdk.NewInt(1),
				FeeReceive: normalExternalAddress,
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

func TestMsgConfirmBatch(t *testing.T) {
	key, _ := crypto.GenerateKey()
	normalExternalAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
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
			testName: "err chain name - empty",
			msg: &types.MsgConfirmBatch{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err orchestrator address",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: errPrefixAddress,
			},
			expectPass: false,
			err:        types.ErrOrchestratorAddress,
			errReason:  fmt.Sprintf("%s: %s", errPrefixAddress, types.ErrOrchestratorAddress),
		},
		{
			testName: "err external address: empty",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     "",
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrExternalAddress),
		},
		{
			testName: "err external address: ToUpper",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     strings.ToUpper(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err external address: ToLower",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     strings.ToLower(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrExternalAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToLower(normalExternalAddress), types.ErrExternalAddress),
		},
		{
			testName: "err token contract address: empty",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				TokenContract:       "",
				Nonce:               0,
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrTokenContractAddress),
		},
		{
			testName: "err token contract address: ToUpper",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				TokenContract:       strings.ToUpper(normalExternalAddress),
				Nonce:               0,
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(normalExternalAddress), types.ErrTokenContractAddress),
		},
		{
			testName: "err external address: ToLower",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				TokenContract:       strings.ToLower(normalExternalAddress),
			},
			expectPass: false,
			err:        types.ErrTokenContractAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToLower(normalExternalAddress), types.ErrTokenContractAddress),
		},
		{
			testName: "err signature: empty",
			msg: &types.MsgConfirmBatch{
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				TokenContract:       normalExternalAddress,
				Signature:           "",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("signature is empty: %s", types.ErrInvalid),
		},
		{
			testName: "err signature: hex.decode error",
			msg: &types.MsgConfirmBatch{
				Nonce:               0,
				ChainName:           chainName,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				TokenContract:       normalExternalAddress,
				Signature:           "gggg",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("could not hex decode signature: %s: %s", "gggg", types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.MsgConfirmBatch{
				Nonce:               0,
				OrchestratorAddress: normalOracleAddress,
				ExternalAddress:     normalExternalAddress,
				TokenContract:       normalExternalAddress,
				Signature:           hex.EncodeToString([]byte("abcd")),
				ChainName:           chainName,
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

func TestUpdateChainOraclesProposal(t *testing.T) {
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
			testName: "err chain name - empty",
			msg: &types.UpdateChainOraclesProposal{
				ChainName: "",
			},
			expectPass: false,
			err:        types.ErrInvalidChainName,
			errReason:  fmt.Sprintf("%s: %s", "", types.ErrInvalidChainName),
		},
		{
			testName: "err oracle: empty",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   chainName,
				Title:       "test title",
				Description: "test description",
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("%s: %s", "oracles cannot be empty", types.ErrInvalid),
		},
		{
			testName: "err external address: err prefix address",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   chainName,
				Title:       "test title",
				Description: "test description",
				Oracles: []string{
					strings.ToUpper(errPrefixAddress),
				},
			},
			expectPass: false,
			err:        types.ErrOracleAddress,
			errReason:  fmt.Sprintf("%s: %s", strings.ToUpper(errPrefixAddress), types.ErrOracleAddress),
		},
		{
			testName: "err oracle: duplicate oracle",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   chainName,
				Title:       "test title",
				Description: "test description",
				Oracles: []string{
					normalOracleAddress,
					normalOracleAddress,
				},
			},
			expectPass: false,
			err:        types.ErrInvalid,
			errReason:  fmt.Sprintf("duplicate oracle %s: %s", normalOracleAddress, types.ErrInvalid),
		},
		{
			testName: "success",
			msg: &types.UpdateChainOraclesProposal{
				ChainName:   chainName,
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
