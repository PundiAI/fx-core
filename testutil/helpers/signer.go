package helpers

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

var _ keyring.Signer = &Signer{}

// Signer defines a type that is used on testing for signing MsgEthereumTx
type Signer struct {
	privKey cryptotypes.PrivKey
}

func NewSigner(sk cryptotypes.PrivKey) *Signer {
	return &Signer{
		privKey: sk,
	}
}

func (s Signer) PrivKey() cryptotypes.PrivKey {
	return s.privKey
}

func (s Signer) PubKey() cryptotypes.PubKey {
	return s.privKey.PubKey()
}

func (s Signer) Address() common.Address {
	return common.BytesToAddress(s.privKey.PubKey().Address())
}

func (s Signer) AccAddress() sdk.AccAddress {
	return s.privKey.PubKey().Address().Bytes()
}

func (s Signer) ValAddress() sdk.ValAddress {
	return s.privKey.PubKey().Address().Bytes()
}

func (s Signer) ExternalAddr(module string) string {
	return fxtypes.ExternalAddrToStr(module, s.AccAddress())
}

// Sign signs the message using the underlying private key
func (s Signer) Sign(_ string, msg []byte, _ signing.SignMode) ([]byte, cryptotypes.PubKey, error) {
	if s.privKey.Type() != ethsecp256k1.KeyType {
		return nil, nil, fmt.Errorf(
			"invalid private key type for signing ethereum tx; expected %s, got %s",
			ethsecp256k1.KeyType,
			s.privKey.Type(),
		)
	}

	sig, err := s.privKey.Sign(msg)
	if err != nil {
		return nil, nil, err
	}

	return sig, s.privKey.PubKey(), nil
}

// SignByAddress sign byte messages with a user key providing the address.
func (s Signer) SignByAddress(address sdk.Address, msg []byte, signMode signing.SignMode) ([]byte, cryptotypes.PubKey, error) {
	signer := sdk.AccAddress(s.privKey.PubKey().Address())
	if !signer.Equals(address) {
		return nil, nil, fmt.Errorf("address mismatch: signer %s ≠ given address %s", signer, address)
	}

	return s.Sign("", msg, signMode)
}
