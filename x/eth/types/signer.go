package types

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"

	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

const signaturePrefix = "\x19Ethereum Signed Message:\n32"

// NewEthereumSignature creates a new signuature over a given byte array
func NewEthereumSignature(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, crosschaintypes.ErrInvalid.Wrapf("private key is nil")
	}
	protectedHash := crypto.Keccak256Hash(append([]uint8(signaturePrefix), hash...))
	return crypto.Sign(protectedHash.Bytes(), privateKey)
}

func EthAddressFromSignature(hash, signature []byte) (string, error) {
	if len(signature) < 65 {
		return "", crosschaintypes.ErrInvalid.Wrapf("signature too short")
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

	protectedHash := crypto.Keccak256Hash(append([]uint8(signaturePrefix), hash...))
	pubkey, err := crypto.SigToPub(protectedHash.Bytes(), signature)
	if err != nil {
		return "", crosschaintypes.ErrInvalid.Wrapf("signature verification failed: %s", err.Error())
	}

	addr := crypto.PubkeyToAddress(*pubkey)
	return addr.Hex(), nil
}

// ValidateEthereumSignature takes a message, an associated signature and public key and
// returns an error if the signature isn't valid
func ValidateEthereumSignature(hash, signature []byte, ethAddress string) error {
	addr, err := EthAddressFromSignature(hash, signature)
	if err != nil {
		return err
	}
	if addr != ethAddress {
		return crosschaintypes.ErrInvalid.Wrapf("signature not matching")
	}
	return nil
}
