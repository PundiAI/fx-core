package helpers

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	hd2 "github.com/evmos/ethermint/crypto/hd"
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
	if algo == hd.Secp256k1Type {
		hdPath = hd.CreateHDPath(118, account, index)
	} else if algo == hd2.EthSecp256k1Type {
		hdPath = hd.CreateHDPath(60, account, index)
	} else {
		return nil, fmt.Errorf("invalid algo")
	}
	signAlgo, err := keyring.NewSigningAlgoFromString(string(algo), hd2.SupportedAlgorithms)
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

func CreateMultiECDSA(count int) []*ecdsa.PrivateKey {
	var ethKeys []*ecdsa.PrivateKey
	for i := 0; i < count; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		ethKeys = append(ethKeys, key)
	}
	return ethKeys
}

// NewPriKey generates cosmos-sdk accAddress private key.
func NewPriKey() cryptotypes.PrivKey {
	return secp256k1.GenPrivKey()
}

// NewEthPrivKey generates an Ethereum address and its corresponding private key.
func NewEthPrivKey() cryptotypes.PrivKey {
	privkey, err := ethsecp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}
	key, err := privkey.ToECDSA()
	if err != nil {
		panic(err)
	}
	addr1 := crypto.PubkeyToAddress(key.PublicKey)
	addr2 := common.BytesToAddress(privkey.PubKey().Address())
	if addr1 != addr2 {
		panic("invalid private key")
	}
	return privkey
}

// GenerateAddress generates an Ethereum address.
func GenerateAddress() common.Address {
	return common.BytesToAddress(NewEthPrivKey().PubKey().Address())
}

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

func (s Signer) Address() common.Address {
	return common.BytesToAddress(s.privKey.PubKey().Address())
}

func (s Signer) AccAddress() sdk.AccAddress {
	return s.privKey.PubKey().Address().Bytes()
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
