package keys

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"gopkg.in/yaml.v2"
)

// Use protobuf interface marshaler rather then generic JSON

// KeyOutput defines a structure wrapping around an Info object used for output
// functionality.
type KeyOutput struct {
	Name         string `json:"name" yaml:"name"`
	Type         string `json:"type" yaml:"type"`
	Eip55Address string `json:"eip55_address,omitempty" yaml:"eip55_address,omitempty"`
	Address      string `json:"address" yaml:"address"`
	PubKey       string `json:"pubkey" yaml:"pubkey"`
	Mnemonic     string `json:"mnemonic,omitempty" yaml:"mnemonic,omitempty"`
}

// NewKeyOutput creates a default KeyOutput instance without Mnemonic, Threshold and PubKeys
func NewKeyOutput(name string, keyType keyring.KeyType, a sdk.Address, pk cryptotypes.PubKey) (KeyOutput, error) {
	apk, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return KeyOutput{}, err
	}
	bz, err := codec.ProtoMarshalJSON(apk, nil)
	if err != nil {
		return KeyOutput{}, err
	}
	keyOutput := KeyOutput{
		Name:    name,
		Type:    keyType.String(),
		Address: a.String(),
		PubKey:  string(bz),
	}
	if strings.EqualFold(pk.Type(), ethsecp256k1.KeyType) {
		keyOutput.Eip55Address = common.BytesToAddress(a.Bytes()).String()
	}
	return keyOutput, nil
}

// MkConsKeyOutput create a KeyOutput in with "cons" Bech32 prefixes.
func MkConsKeyOutput(k *keyring.Record) (KeyOutput, error) {
	pk, err := k.GetPubKey()
	if err != nil {
		return KeyOutput{}, err
	}
	addr := sdk.ConsAddress(pk.Address())
	return NewKeyOutput(k.Name, k.GetType(), addr, pk)
}

// MkValKeyOutput create a KeyOutput in with "val" Bech32 prefixes.
func MkValKeyOutput(k *keyring.Record) (KeyOutput, error) {
	pk, err := k.GetPubKey()
	if err != nil {
		return KeyOutput{}, err
	}

	addr := sdk.ValAddress(pk.Address())
	return NewKeyOutput(k.Name, k.GetType(), addr, pk)
}

// MkAccKeyOutput create a KeyOutput in with "acc" Bech32 prefixes. If the
// public key is a multisig public key, then the threshold and constituent
// public keys will be added.
func MkAccKeyOutput(k *keyring.Record) (KeyOutput, error) {
	pk, err := k.GetPubKey()
	if err != nil {
		return KeyOutput{}, err
	}
	addr := sdk.AccAddress(pk.Address())
	return NewKeyOutput(k.Name, k.GetType(), addr, pk)
}

// MkAccKeysOutput returns a slice of KeyOutput objects, each with the "acc"
// Bech32 prefixes, given a slice of Info objects. It returns an error if any
// call to MkKeyOutput fails.
func MkAccKeysOutput(records []*keyring.Record) ([]KeyOutput, error) {
	kos := make([]KeyOutput, len(records))
	var err error
	for i, r := range records {
		kos[i], err = MkAccKeyOutput(r)
		if err != nil {
			return nil, err
		}
	}

	return kos, nil
}

func printInfo(w io.Writer, kos interface{}, output string) {
	switch output {
	case flags.OutputFormatText:
		out, err := yaml.Marshal(&kos)
		if err != nil {
			panic(err)
		}
		_, _ = fmt.Fprintln(w, string(out))

	case flags.OutputFormatJSON:
		out, err := json.Marshal(kos)
		if err != nil {
			panic(err)
		}

		_, _ = fmt.Fprintln(w, string(out))
	}
}
