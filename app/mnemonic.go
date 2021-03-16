package app

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
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

func NewPrivKeyFromMnemonic(mnemonic string) *secp256k1.PrivKey {
	hdPath := hd.CreateHDPath(sdk.CoinType, 0, 0).String()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyring.SigningAlgoList{hd.Secp256k1})
	if err != nil {
		panic(err.Error())
	}
	// create master key and derive first key for keyring
	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath)
	if err != nil {
		panic(err.Error())
	}
	privKey := algo.Generate()(derivedPriv)
	accPriKey, ok := privKey.(*secp256k1.PrivKey)
	if !ok {
		panic("not secp256k1.PrivKey")
	}
	return accPriKey
}
