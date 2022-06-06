package evm_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	erc20types "github.com/functionx/fx-core/x/erc20/types"
)

const (
	_          = "december slow blue fury silly bread friend unknown render resource dry buyer brand final abstract gallery slow since hood shadow neglect travel convince foil"
	privateKey = "845e3bad4fcb74154b491119be75a655f456efae0a8040476a440483ed91de79"
)

func Test_PrivateKey(t *testing.T) {
	priKey, err := mnemonicToFxPrivKey("december slow blue fury silly bread friend unknown render resource dry buyer brand final abstract gallery slow since hood shadow neglect travel convince foil")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("hex", hex.EncodeToString(priKey.Bytes()))
	t.Log("fx address", sdk.AccAddress(priKey.PubKey().Address().Bytes()).String())
}

func Test_FxAddress(t *testing.T) {
	bz, _ := hex.DecodeString(privateKey)
	priKey := secp256k1.PrivKey{Key: bz}

	pubkey, err := legacybech32.MarshalPubKey(legacybech32.AccPK, priKey.PubKey())
	require.NoError(t, err)
	t.Log("fx pubkey", pubkey)
	t.Log("fx hex public", hex.EncodeToString(priKey.PubKey().Bytes()))
	uncompressedPubKey, err := crypto.DecompressPubkey(priKey.PubKey().Bytes())
	require.NoError(t, err)
	t.Log("fx uncompressed public", hex.EncodeToString(crypto.FromECDSAPub(uncompressedPubKey)))
	acc := sdk.AccAddress(priKey.PubKey().Address().Bytes())
	t.Log("fx address", hex.EncodeToString(acc.Bytes()))
	t.Log("fx address", acc.String())

	priKeyEth, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("eth public", hex.EncodeToString(crypto.FromECDSAPub(&priKeyEth.PublicKey)))
	t.Log("eth compress public", hex.EncodeToString(crypto.CompressPubkey(&priKeyEth.PublicKey)))
	t.Log("eth address", hex.EncodeToString(crypto.PubkeyToAddress(priKeyEth.PublicKey).Bytes()))
	t.Log("eth address", crypto.PubkeyToAddress(priKeyEth.PublicKey).Hex())
}

func Test_Erc20ModuleAddress(t *testing.T) {
	address := erc20types.ModuleAddress
	assert.Equal(t, address.Hex(), "0x47EeB2eac350E1923b8CBDfA4396A077b36E62a0")
	assert.Equal(t, sdk.AccAddress(address.Bytes()).String(), "fx1glht96kr2rseywuvhhay894qw7ekuc4qs5z0yh")
}
