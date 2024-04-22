package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

func TestMsgGrantPrivilegeRoute(t *testing.T) {
	val := sdk.ValAddress("val")
	addr1 := sdk.AccAddress("from")
	toKey := helpers.NewPriKey()

	msg, err := types.NewMsgGrantPrivilege(val, addr1, toKey.PubKey())
	require.NoError(t, err)
	require.Equal(t, msg.Route(), stakingtypes.RouterKey)
	require.Equal(t, msg.Type(), "grant_privilege")
}

func TestMsgGrantPrivilegeValidation(t *testing.T) {
	key1 := helpers.NewPriKey()
	key2 := helpers.NewPriKey()
	eth3 := helpers.NewEthPrivKey()
	eth4 := helpers.NewEthPrivKey()

	any1, _ := codectypes.NewAnyWithValue(key1.PubKey())
	any2, _ := codectypes.NewAnyWithValue(key2.PubKey())
	any3, _ := codectypes.NewAnyWithValue(eth3.PubKey())
	any4, _ := codectypes.NewAnyWithValue(eth4.PubKey())

	val := sdk.ValAddress(key1.PubKey().Address())
	addr1 := sdk.AccAddress(key1.PubKey().Address())
	addr3 := sdk.AccAddress(eth3.PubKey().Address())

	invalidAddr := "xxxxxxxxxxxx"
	invalidPK, _ := codectypes.NewAnyWithValue(&banktypes.MsgSend{})
	invalidPKType, _ := codectypes.NewAnyWithValue(ed25519.GenPrivKey().PubKey())

	cases := []struct {
		name        string
		expectedErr string // empty means no error expected
		msg         *types.MsgGrantPrivilege
	}{
		{"valid grant to acc address", "", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr1.String(), ToPubkey: any2}},
		{"valid grant to eth address", "", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr1.String(), ToPubkey: any3}},
		{"valid grant eth to eth address", "", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr3.String(), ToPubkey: any4}},

		{"empty validator address", "invalid validator address: empty address string is not allowed: invalid address", &types.MsgGrantPrivilege{ValidatorAddress: "", FromAddress: addr1.String(), ToPubkey: any2}},
		{"invalid validator address", "invalid validator address: decoding bech32 failed: invalid separator index -1: invalid address", &types.MsgGrantPrivilege{FromAddress: addr1.String(), ToPubkey: any2, ValidatorAddress: invalidAddr}},

		{"empty from address", "invalid from address: empty address string is not allowed: invalid address", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: "", ToPubkey: any2}},
		{"invalid from address", "invalid from address: decoding bech32 failed: invalid separator index -1: invalid address", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: invalidAddr, ToPubkey: any2}},

		{"empty pubkey", "empty pubkey: invalid pubkey", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr1.String(), ToPubkey: nil}},
		{"invalid pubkey type", "expecting cryptotypes.PubKey, got *types.MsgSend: invalid pubkey", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr1.String(), ToPubkey: invalidPK}},
		{"invalid pubkey key type", "expecting *secp256k1.PubKey or *ethsecp256k1.PubKey, got *ed25519.PubKey: invalid pubkey", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr1.String(), ToPubkey: invalidPKType}},

		{"same from and to address", "same account: invalid request", &types.MsgGrantPrivilege{ValidatorAddress: val.String(), FromAddress: addr1.String(), ToPubkey: any1}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgGrantPrivilegeGetSignBytes(t *testing.T) {
	val := sdk.ValAddress("val")
	addr1 := sdk.AccAddress("input")
	bz, _ := hex.DecodeString("0370c5fe92d864015703aff6d2f3a5608c3740e368370f0c25f090abc2e368b0be")
	toPk := &ethsecp256k1.PubKey{Key: bz}

	msg, err := types.NewMsgGrantPrivilege(val, addr1, toPk)
	require.NoError(t, err)

	res := msg.GetSignBytes()
	expected := `{"type":"staking/MsgGrantPrivilege","value":{"from_address":"cosmos1d9h8qat57ljhcm","to_pubkey":"A3DF/pLYZAFXA6/20vOlYIw3QONoNw8MJfCQq8LjaLC+","validator_address":"cosmosvaloper1weskc8rudzy"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgGrantPrivilegeGetSigners(t *testing.T) {
	toKey := helpers.NewEthPrivKey()
	msg, err := types.NewMsgGrantPrivilege(sdk.ValAddress{}, []byte("input111111111111111"), toKey.PubKey())
	require.NoError(t, err)

	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")
}

func consPubKey() cryptotypes.PubKey {
	var pk cryptotypes.PubKey
	pkStr := `{"@type":"/cosmos.crypto.ed25519.PubKey","key":"ua9OcG6txvyZ2wSxeVR1+NDzUqO8TzZdoYxCyA48qAM="}`
	if err := app.MakeEncodingConfig().Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		panic(err)
	}
	return pk
}

func TestMsgEditConsensusKey_Route(t *testing.T) {
	val := sdk.ValAddress("val")
	from := sdk.AccAddress("from")
	pubKey := consPubKey()

	msg, err := types.NewMsgEditConsensusPubKey(val, from, pubKey)
	require.NoError(t, err)
	require.Equal(t, msg.Route(), stakingtypes.RouterKey)
	require.Equal(t, msg.Type(), "edit_consensus_pubkey")
}

func TestMsgEditConsensusKeyValidation(t *testing.T) {
	key := helpers.NewPriKey()
	val := sdk.ValAddress(key.PubKey().Address().Bytes())
	from := sdk.AccAddress(key.PubKey().Address().Bytes())

	testCase := []struct {
		name        string
		msg         func() *types.MsgEditConsensusPubKey
		expectError string
	}{
		{
			name: "valid",
			msg: func() *types.MsgEditConsensusPubKey {
				msg, err := types.NewMsgEditConsensusPubKey(val, from, consPubKey())
				require.NoError(t, err)
				return msg
			},
		},
		{
			name: "invalid validator address",
			msg: func() *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: "",
					From:             from.String(),
					Pubkey:           nil,
				}
			},
			expectError: "invalid validator address: empty address string is not allowed: invalid address",
		},
		{
			name: "invalid from address",
			msg: func() *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             "",
					Pubkey:           nil,
				}
			},
			expectError: "invalid from address: empty address string is not allowed: invalid address",
		},
		{
			name: "invalid pubkey",
			msg: func() *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           nil,
				}
			},
			expectError: "empty validator public key",
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg().ValidateBasic()
			if len(tc.expectError) == 0 {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectError)
			}
		})
	}
}

func TestMsgEditConsensusKeyGetSignBytes(t *testing.T) {
	val := sdk.ValAddress("val")
	from := sdk.AccAddress("from")
	msg, err := types.NewMsgEditConsensusPubKey(val, from, consPubKey())
	require.NoError(t, err)
	res := msg.GetSignBytes()

	expected := `{"type":"staking/MsgEditConsensusPubKey","value":{"from":"cosmos1veex7mgzt83cu","pubkey":"ua9OcG6txvyZ2wSxeVR1+NDzUqO8TzZdoYxCyA48qAM=","validator_address":"cosmosvaloper1weskc8rudzy"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgEditConsensusKeyGetSigners(t *testing.T) {
	msg, err := types.NewMsgEditConsensusPubKey(sdk.ValAddress{}, []byte("input111111111111111"), consPubKey())
	require.NoError(t, err)
	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")
}
