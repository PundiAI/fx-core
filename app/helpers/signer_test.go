package helpers_test

import (
	"testing"

	hd2 "github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v2/app/helpers"
)

func TestPrivKeyFromMnemonic(t *testing.T) {
	mnemonic := "afraid bottom consider camera dragon seek text other addict shift tent zebra feel swim approve prize promote secret process method quick engine return pitch"
	privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, hd2.Secp256k1Type, 0, 0)
	require.NoError(t, err)

	require.Equal(t, "fx18e2hlvcupj2dqhpty99044h2j7jp5d5wh670n8", sdk.AccAddress(privKey.PubKey().Address().Bytes()).String())
}
