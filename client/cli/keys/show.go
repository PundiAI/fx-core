package keys

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/ledger"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	// FlagAddress is the flag for the user's address on the command line.
	FlagAddress = "address"
	// FlagEip55Address is the flag for the user's EIP address on the command line.
	FlagEip55Address = "eip-address"
	// FlagPublicKey represents the user's public key on the command line.
	FlagPublicKey = "pubkey"
	// FlagBechPrefix defines a desired Bech32 prefix encoding for a key.
	FlagBechPrefix = "bech"
	// FlagDevice indicates that the information should be shown in the device
	FlagDevice = "device"

	flagMultiSigThreshold = "multisig-threshold"
)

// ShowKeysCmd shows key information for a given key name.
func ShowKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [name_or_address [name_or_address...]]",
		Short: "Retrieve key information by name or address",
		Long: `Display keys details. If multiple names or addresses are provided,
then an ephemeral multisig key will be created under the name "multi"
consisting of all the keys provided by name and multisig threshold.`,
		Args: cobra.MinimumNArgs(1),
		RunE: runShowCmd,
	}
	f := cmd.Flags()
	f.String(FlagBechPrefix, sdk.PrefixAccount, "The Bech32 prefix encoding for a key (acc|val|cons)")
	f.BoolP(FlagAddress, "a", false, "Output the bech32 format address only (overrides --output)")
	f.BoolP(FlagEip55Address, "e", false, "Output the EIP55 address only (overrides --output)")
	f.BoolP(FlagPublicKey, "p", false, "Output the public key only (overrides --output)")
	f.BoolP(FlagDevice, "d", false, "Output the address in a ledger device")
	f.Int(flagMultiSigThreshold, 1, "K out of N required signatures")

	return cmd
}

//gocyclo:ignore
func runShowCmd(cmd *cobra.Command, args []string) (err error) {
	k := new(keyring.Record)
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		k, err = fetchKey(clientCtx.Keyring, args[0])
		if err != nil {
			return fmt.Errorf("%s is not a valid name or address: %w", args[0], err)
		}
	} else {
		pks := make([]cryptotypes.PubKey, len(args))
		for i, keyref := range args {
			k, err := fetchKey(clientCtx.Keyring, keyref)
			if err != nil {
				return fmt.Errorf("%s is not a valid name or address: %w", keyref, err)
			}
			key, err := k.GetPubKey()
			if err != nil {
				return err
			}
			pks[i] = key
		}

		multisigThreshold, _ := cmd.Flags().GetInt(flagMultiSigThreshold)
		err = validateMultisigThreshold(multisigThreshold, len(args))
		if err != nil {
			return err
		}

		multikey := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)
		k, err = keyring.NewMultiRecord(k.Name, multikey)
		if err != nil {
			return err
		}
	}

	isShowAddr, _ := cmd.Flags().GetBool(FlagAddress)
	isShowEip55Addr, _ := cmd.Flags().GetBool(FlagEip55Address)
	isShowPubKey, _ := cmd.Flags().GetBool(FlagPublicKey)
	isShowDevice, _ := cmd.Flags().GetBool(FlagDevice)

	isOutputSet := false
	tmp := cmd.Flag(cli.OutputFlag)
	if tmp != nil {
		isOutputSet = tmp.Changed
	}

	if isShowEip55Addr && isShowAddr && isShowPubKey {
		return errors.New("cannot use both --address and --pubkey at once")
	}

	if isOutputSet && (isShowEip55Addr || isShowAddr || isShowPubKey) {
		return errors.New("cannot use --output with --address or --pubkey")
	}

	bechPrefix, _ := cmd.Flags().GetString(FlagBechPrefix)
	bechKeyOut, err := getBechKeyOut(bechPrefix)
	if err != nil {
		return err
	}

	if isOutputSet {
		clientCtx.OutputFormat, _ = cmd.Flags().GetString(cli.OutputFlag)
	}

	switch {
	case isShowEip55Addr, isShowAddr, isShowPubKey:
		ko, err := bechKeyOut(k)
		if err != nil {
			return err
		}
		out := ko.Address
		if isShowEip55Addr {
			out = ko.Eip55Address
		}
		if isShowPubKey {
			out = ko.PubKey
		}
		cmd.Println(out)
	default:
		outInfo, err := bechKeyOut(k)
		if err != nil {
			return err
		}
		printInfo(cmd.OutOrStdout(), outInfo, clientCtx.OutputFormat)
	}

	if isShowDevice {
		if isShowPubKey {
			return fmt.Errorf("the device flag (-d) can only be used for addresses not pubkeys")
		}
		if bechPrefix != "acc" {
			return fmt.Errorf("the device flag (-d) can only be used for accounts")
		}

		// Override and show in the device
		if k.GetType() != keyring.TypeLedger {
			return fmt.Errorf("the device flag (-d) can only be used for accounts stored in devices")
		}

		ledgerItem := k.GetLedger()
		if ledgerItem == nil {
			return errors.New("unable to get ledger item")
		}

		pk, err := k.GetPubKey()
		if err != nil {
			return err
		}

		return ledger.ShowAddress(*ledgerItem.Path, pk, sdk.GetConfig().GetBech32AccountAddrPrefix())
	}

	return nil
}

func fetchKey(kb keyring.Keyring, keyref string) (*keyring.Record, error) {
	// firstly check if the keyref is a key name of a key registered in a keyring.
	k, err := kb.Key(keyref)
	// if the key is not there or if we have a problem with a keyring itself then we move to a
	// fallback: searching for key by address.
	if err == nil || !errorsmod.IsOf(err, errortypes.ErrIO, errortypes.ErrKeyNotFound) {
		return k, err
	}
	accAddr, err := sdk.AccAddressFromBech32(keyref)
	if err != nil {
		return k, err
	}

	k, err = kb.KeyByAddress(accAddr)
	return k, errorsmod.Wrap(err, "Invalid key")
}

func validateMultisigThreshold(k, nKeys int) error {
	if k <= 0 {
		return fmt.Errorf("threshold must be a positive integer")
	}
	if nKeys < k {
		return fmt.Errorf(
			"threshold k of n multisignature: %d < %d", nKeys, k)
	}
	return nil
}

type bechKeyOutFn func(keyInfo *keyring.Record) (KeyOutput, error)

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case sdk.PrefixAccount:
		return MkAccKeyOutput, nil
	case sdk.PrefixValidator:
		return MkValKeyOutput, nil
	case sdk.PrefixConsensus:
		return MkConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}
