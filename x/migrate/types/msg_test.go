package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	_ "github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/migrate/types"
)

func TestMsgMigrateAccountRoute(t *testing.T) {
	addr1 := sdk.AccAddress("from")
	addr2 := sdk.AccAddress("to")
	var msg = types.NewMsgMigrateAccount(addr1, addr2, "empty string")

	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "migrate_account")
}

func TestMsgMigrateAccountValidation(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	toAddressByte := crypto.PubkeyToAddress(privateKey.PublicKey)
	addr1 := sdk.AccAddress("from________________")
	addr2 := sdk.AccAddress(toAddressByte.Bytes())
	addrTo := sdk.AccAddress("to__________________")
	addrEmpty := sdk.AccAddress("")
	addrTooLong := sdk.AccAddress("Accidentally used 268 bytes pubkey test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content")

	sign, err := crypto.Sign(types.MigrateAccountSignatureHash(addr1, addr2), privateKey)
	require.NoError(t, err)
	validSignHex := hex.EncodeToString(sign)

	emptySign := ""
	invalidSign := "xx xxx"
	errSignStr := hex.EncodeToString(append([]byte("----------------------------------length 64 and latest less than"), []byte{0x1}...))

	cases := []struct {
		name        string
		expectedErr string // empty means no error expected
		msg         *types.MsgMigrateAccount
	}{
		{"valid migrate", "", types.NewMsgMigrateAccount(addr1, addr2, validSignHex)},

		{"empty from address", "invalid sender address (empty address string is not allowed): invalid address", types.NewMsgMigrateAccount(addrEmpty, addr2, emptySign)},
		{"invalid from address", "invalid sender address (address max length is 255, got 268: unknown address): invalid address", types.NewMsgMigrateAccount(addrTooLong, addr2, emptySign)},

		{"empty to address", "invalid to address (empty address string is not allowed): invalid address", types.NewMsgMigrateAccount(addr1, addrEmpty, emptySign)},
		{"invalid to address", "invalid to address (address max length is 255, got 268: unknown address): invalid address", types.NewMsgMigrateAccount(addr1, addrTooLong, emptySign)},

		{"same from address to address", fmt.Sprintf("%s: same account", addr1.String()), types.NewMsgMigrateAccount(addr1, addr1, emptySign)},

		{"empty sign", "signature is empty: invalid signature", types.NewMsgMigrateAccount(addr1, addr2, emptySign)},
		{"invalid sign", "could not hex decode signature: xx xxx: invalid signature", types.NewMsgMigrateAccount(addr1, addr2, invalidSign)},
		{"signature key not equal to address", fmt.Sprintf("signature key not equal to address, expected %s, got fx1n3g5mp5y08at0wdkw8uktqgmvj3yzunmyt8qsf: invalid signature", addrTo.String()), types.NewMsgMigrateAccount(addr1, addrTo, errSignStr)},
	}

	for _, tc := range cases {
		err = tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgMigrateAccountGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("input"))
	addr2 := sdk.AccAddress([]byte("output"))
	sign := "0x1"
	var msg = types.NewMsgMigrateAccount(addr1, addr2, sign)
	res := msg.GetSignBytes()

	expected := `{"type":"migrate/MsgMigrateAccount","value":{"from":"fx1d9h8qat556rm7q","signature":"0x1","to":"fx1da6hgur4ws9u7g24"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgMigrateAccountGetSigners(t *testing.T) {
	var msg = types.NewMsgMigrateAccount(sdk.AccAddress([]byte("input111111111111111")), sdk.AccAddress{}, "")
	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")
}
