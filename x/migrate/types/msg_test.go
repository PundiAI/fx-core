package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	_ "github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/x/migrate/types"
)

func TestMsgMigrateAccountRoute(t *testing.T) {
	addr1 := sdk.AccAddress("from")
	addr2 := common.BytesToAddress([]byte("to"))
	msg := types.NewMsgMigrateAccount(addr1, addr2, "empty string")

	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "migrate_account")
}

func TestMsgMigrateAccountValidation(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	toAddressByte := crypto.PubkeyToAddress(privateKey.PublicKey)
	addr1 := sdk.AccAddress("from________________")
	addr2 := toAddressByte
	addrEmpty := sdk.AccAddress("")
	addrTooLong := sdk.AccAddress("Accidentally used 268 bytes pubkey test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content")

	sign, err := crypto.Sign(types.MigrateAccountSignatureHash(addr1, addr2.Bytes()), privateKey)
	require.NoError(t, err)
	validSignHex := hex.EncodeToString(sign)

	emptySign := ""
	invalidSign := "xx xxx"
	data := append([]byte("----------------------------------length 64 and latest less than"), []byte{0x1}...)
	errSignStr := hex.EncodeToString(data)

	cases := []struct {
		name        string
		expectedErr string // empty means no error expected
		msg         *types.MsgMigrateAccount
	}{
		{"valid migrate", "", types.NewMsgMigrateAccount(addr1, addr2, validSignHex)},

		{"empty from address", "invalid from address: empty address string is not allowed: invalid address", types.NewMsgMigrateAccount(addrEmpty, addr2, emptySign)},
		{"invalid from address", "invalid from address: address max length is 255, got 268: unknown address: invalid address", types.NewMsgMigrateAccount(addrTooLong, addr2, emptySign)},

		{"empty to address", "invalid to address: empty: invalid address", &types.MsgMigrateAccount{From: addr1.String(), To: "", Signature: emptySign}},
		{"invalid to address", "invalid to address: wrong length: invalid address", &types.MsgMigrateAccount{From: addr1.String(), To: "1234567890", Signature: emptySign}},

		{"same from address to address", "same account: invalid request", types.NewMsgMigrateAccount(addr1, common.BytesToAddress(addr1.Bytes()), emptySign)},

		{"empty sign", "empty signature: invalid request", types.NewMsgMigrateAccount(addr1, addr2, emptySign)},
		{"invalid sign", "could not hex decode signature: invalid request", types.NewMsgMigrateAccount(addr1, addr2, invalidSign)},
		{"signature key not equal to address", "signature key not equal to address: invalid request", types.NewMsgMigrateAccount(addr1, addr2, errSignStr)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err = tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgMigrateAccountGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress("input")
	addr2 := common.BytesToAddress([]byte("output"))
	sign := "0x1"
	msg := types.NewMsgMigrateAccount(addr1, addr2, sign)
	res := msg.GetSignBytes()

	expected := `{"type":"migrate/MsgMigrateAccount","value":{"from":"cosmos1d9h8qat57ljhcm","signature":"0x1","to":"0x00000000000000000000000000006F7574707574"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgMigrateAccountGetSigners(t *testing.T) {
	msg := types.NewMsgMigrateAccount([]byte("input111111111111111"), common.Address{}, "")
	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")
}
