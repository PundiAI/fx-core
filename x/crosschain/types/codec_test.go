package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/app"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func Test_Codec_MsgConfirm(t *testing.T) {
	interfaceRegistry := app.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	crosschaintypes.RegisterInterfaces(interfaceRegistry)

	bridgeCallConfirm := mockBridgeCallConfirm()
	bridgeCallConfirmAny, err := types.NewAnyWithValue(bridgeCallConfirm)
	require.NoError(t, err)
	confirm, ok := bridgeCallConfirmAny.GetCachedValue().(crosschaintypes.Confirm)
	require.True(t, ok)
	require.Equal(t, crosschaintypes.Confirm(bridgeCallConfirm), confirm)

	msgConfirm := mockMsgConfirm(bridgeCallConfirmAny)
	result := appCodec.MustMarshal(msgConfirm)

	msgConfirm2 := &crosschaintypes.MsgConfirm{}
	appCodec.MustUnmarshal(result, msgConfirm2)
	require.Equal(t, msgConfirm, msgConfirm2)
	require.Equal(t, msgConfirm.String(), msgConfirm2.String())

	confirm, ok = msgConfirm2.Confirm.GetCachedValue().(crosschaintypes.Confirm)
	require.True(t, ok)
	require.Equal(t, crosschaintypes.Confirm(bridgeCallConfirm), confirm)
}

func mockMsgConfirm(bridgeCallConfirmAny *types.Any) *crosschaintypes.MsgConfirm {
	return &crosschaintypes.MsgConfirm{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: helpers.GenHexAddress().String(),
		Confirm:        bridgeCallConfirmAny,
	}
}

func mockBridgeCallConfirm() *crosschaintypes.MsgBridgeCallConfirm {
	return &crosschaintypes.MsgBridgeCallConfirm{
		ChainName:       ethtypes.ModuleName,
		BridgerAddress:  helpers.GenHexAddress().String(),
		ExternalAddress: helpers.GenHexAddress().String(),
		Nonce:           1,
		Signature:       helpers.GenHexAddress().String(),
	}
}

func Test_Codec_MsgClaim(t *testing.T) {
	interfaceRegistry := app.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	crosschaintypes.RegisterInterfaces(interfaceRegistry)

	bridgeCallClaim := &crosschaintypes.MsgBridgeCallClaim{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: helpers.GenHexAddress().String(),
		EventNonce:     1,
		BlockHeight:    1,
		Sender:         helpers.GenHexAddress().String(),
		Refund:         helpers.GenHexAddress().String(),
		TokenContracts: []string{helpers.GenHexAddress().String()},
		Amounts:        []sdkmath.Int{sdkmath.NewInt(1)},
		To:             helpers.GenHexAddress().String(),
		Data:           "data",
		QuoteId:        sdkmath.NewInt(1),
		GasLimit:       sdkmath.NewInt(1),
		Memo:           "memo",
		TxOrigin:       helpers.GenHexAddress().String(),
	}
	bridgeCallClaimAny, err := types.NewAnyWithValue(bridgeCallClaim)
	require.NoError(t, err)
	claim, ok := bridgeCallClaimAny.GetCachedValue().(crosschaintypes.ExternalClaim)
	require.True(t, ok)
	require.Equal(t, crosschaintypes.ExternalClaim(bridgeCallClaim), claim)

	msgClaim := &crosschaintypes.MsgClaim{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: helpers.GenHexAddress().String(),
		Claim:          bridgeCallClaimAny,
	}
	result := appCodec.MustMarshal(msgClaim)

	msgClaim2 := &crosschaintypes.MsgClaim{}
	appCodec.MustUnmarshal(result, msgClaim2)
	require.Equal(t, msgClaim, msgClaim2)
	require.Equal(t, msgClaim.String(), msgClaim2.String())

	claim, ok = msgClaim2.Claim.GetCachedValue().(crosschaintypes.ExternalClaim)
	require.True(t, ok)
	require.Equal(t, crosschaintypes.ExternalClaim(bridgeCallClaim), claim)
}
