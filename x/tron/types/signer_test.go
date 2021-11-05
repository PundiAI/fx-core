package types

import (
	"crypto/sha256"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
)

func TestEthAddressFromSignature(t *testing.T) {
	dataHash := sha256.Sum256([]byte("tron - test - data"))
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	ethereumSignature, err := NewTronSignature(dataHash[:], key)
	if err != nil {
		t.Fatal(err)
	}

	tronAddr := tronAddress.PubkeyToAddress(key.PublicKey)
	recoverySignAddress, err := TronAddressFromSignature(dataHash[:], ethereumSignature)
	if err != nil {
		t.Fatal(err)
	}
	require.EqualValues(t, tronAddr.String(), recoverySignAddress)
}

func TestValidateEthereumSignature(t *testing.T) {
	dataHash := sha256.Sum256([]byte("tron - test - data"))
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	ethereumSignature, err := NewTronSignature(dataHash[:], key)
	require.NoError(t, err)

	err = ValidateTronSignature(dataHash[:], ethereumSignature, tronAddress.PubkeyToAddress(key.PublicKey).String())
	require.NoError(t, err)
}
