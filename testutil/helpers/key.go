package helpers

import (
	"crypto/ecdsa"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	hd2 "github.com/evmos/ethermint/crypto/hd"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
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
	switch algo {
	case hd.Secp256k1Type:
		hdPath = hd.CreateHDPath(118, account, index)
	case hd2.EthSecp256k1Type:
		hdPath = hd.CreateHDPath(60, account, index)
	default:
		return nil, sdkerrors.ErrInvalidPubKey
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

// GenHexAddress generates an Ethereum address.
func GenHexAddress() common.Address {
	return common.BytesToAddress(NewEthPrivKey().PubKey().Address())
}

// GenAccAddress generates an cosmos-sdk accAddress
func GenAccAddress() sdk.AccAddress {
	return NewPriKey().PubKey().Address().Bytes()
}

// GenExternalAddr generates an Ethereum or Tron address.
func GenExternalAddr(module string) string {
	addr := GenHexAddress()
	return types.ExternalAddrToStr(module, addr.Bytes())
}

// GenZeroExternalAddr generates an Ethereum or Tron zero address.
func GenZeroExternalAddr(module string) string {
	addr := common.Address{}
	return types.ExternalAddrToStr(module, addr.Bytes())
}

// HexAddrToTronAddr returns a Tron address from an hex string.
func HexAddrToTronAddr(str string) string {
	return types.ExternalAddrToStr(trontypes.ModuleName, common.FromHex(str))
}
