package helpers

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	hd2 "github.com/evmos/ethermint/crypto/hd"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

func NewMnemonic() string {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		panic(err)
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		panic(err)
	}
	return mnemonic
}

func PrivKeyFromMnemonic(mnemonic string, algo hd.PubKeyType, account, index uint32) (cryptotypes.PrivKey, error) {
	var hdPath *hd.BIP44Params
	var signAlgo keyring.SignatureAlgo
	var err error
	if algo == hd.Secp256k1Type {
		hdPath = hd.CreateHDPath(118, account, index)
		signAlgo, err = keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyring.SigningAlgoList{hd.Secp256k1})
	} else if algo == hd2.EthSecp256k1Type {
		hdPath = hd.CreateHDPath(60, account, index)
		signAlgo, err = keyring.NewSigningAlgoFromString(string(hd2.EthSecp256k1Type), keyring.SigningAlgoList{hd2.EthSecp256k1})
	}
	if err != nil {
		return nil, err
	}
	// create master key and derive first key for keyring
	derivedPriv, err := signAlgo.Derive()(mnemonic, "", hdPath.String())
	if err != nil {
		return nil, err
	}
	privKey := signAlgo.Generate()(derivedPriv)
	return privKey, nil
}

func CreateMultiEthKey(count int) []*ecdsa.PrivateKey {
	var ethKeys []*ecdsa.PrivateKey
	for i := 0; i < count; i++ {
		ethKeys = append(ethKeys, GenerateEthKey())
	}
	return ethKeys
}

func GenerateEthKey() *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return key
}

func GenerateEthAddress() common.Address {
	return crypto.PubkeyToAddress(GenerateEthKey().PublicKey)
}

// NewPriKey generates cosmos-sdk accAddress private key.
func NewPriKey() cryptotypes.PrivKey {
	return secp256k1.GenPrivKey()
}

// NewAddrKey generates an Ethereum address and its corresponding private key.
func NewAddrKey() (common.Address, cryptotypes.PrivKey) {
	privkey, err := ethsecp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}
	key, err := privkey.ToECDSA()
	if err != nil {
		panic(err)
	}

	addr := crypto.PubkeyToAddress(key.PublicKey)

	return addr, privkey
}

// GenerateAddress generates an Ethereum address.
func GenerateAddress() common.Address {
	addr, _ := NewAddrKey()
	return addr
}

var _ keyring.Signer = &Signer{}

// Signer defines a type that is used on testing for signing MsgEthereumTx
type Signer struct {
	privKey cryptotypes.PrivKey
}

func NewSigner(sk cryptotypes.PrivKey) keyring.Signer {
	return &Signer{
		privKey: sk,
	}
}

// Sign signs the message using the underlying private key
func (s Signer) Sign(_ string, msg []byte) ([]byte, cryptotypes.PubKey, error) {
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
func (s Signer) SignByAddress(address sdk.Address, msg []byte) ([]byte, cryptotypes.PubKey, error) {
	signer := sdk.AccAddress(s.privKey.PubKey().Address())
	if !signer.Equals(address) {
		return nil, nil, fmt.Errorf("address mismatch: signer %s ≠ given address %s", signer, address)
	}

	return s.Sign("", msg)
}
