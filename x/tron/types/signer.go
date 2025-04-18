package types

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

const tronSignaturePrefix = "\x19TRON Signed Message:\n32"

// NewTronSignature creates a new signuature over a given byte array
func NewTronSignature(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, types.ErrInvalid.Wrapf("private key")
	}
	protectedHash := crypto.Keccak256Hash(append([]uint8(tronSignaturePrefix), hash...))
	return crypto.Sign(protectedHash.Bytes(), privateKey)
}

func TronAddressFromSignature(hash, signature []byte) (string, error) {
	if len(signature) < 65 {
		return "", types.ErrInvalid.Wrapf("signature too short")
	}
	// To verify signature
	// - use crypto.SigToPub to get the public key
	// - use crypto.PubkeyToAddress to get the address
	// - compare this to the address given.

	// for backwards compatibility reasons  the V value of an Ethereum sig is presented
	// as 27 or 28, internally though it should be a 0-3 value due to changed formats.
	// It seems that go-ethereum expects this to be done before sigs actually reach it's
	// internal validation functions. In order to comply with this requirement we check
	// the sig an dif it's in standard format we correct it. If it's in go-ethereum's expected
	// format already we make no changes.
	//
	// We could attempt to break or otherwise exit early on obviously invalid values for this
	// byte, but that's a task best left to go-ethereum
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	protectedHash := crypto.Keccak256Hash(append([]uint8(tronSignaturePrefix), hash...))
	pubkey, err := crypto.SigToPub(protectedHash.Bytes(), signature)
	if err != nil {
		return "", types.ErrInvalid.Wrapf("signature verification failed: %s", err.Error())
	}

	addr := address.PubkeyToAddress(*pubkey)
	return addr.String(), nil
}

// ValidateTronSignature takes a message, an associated signature and public key and
// returns an error if the signature isn't valid
func ValidateTronSignature(hash, signature []byte, ethAddress string) error {
	addr, err := TronAddressFromSignature(hash, signature)
	if err != nil {
		return err
	}
	if addr != ethAddress {
		return types.ErrInvalid.Wrapf("signature not matching")
	}
	return nil
}
