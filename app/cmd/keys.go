package cmd

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/ledger"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	yaml "gopkg.in/yaml.v2"
	"io"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	etherminthd "github.com/functionx/fx-core/crypto/hd"

	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
)

const (
	flagInteractive       = "interactive"
	flagRecover           = "recover"
	flagNoBackup          = "no-backup"
	flagCoinType          = "coin-type"
	flagAccount           = "account"
	flagIndex             = "index"
	flagMultisig          = "multisig"
	flagMultiSigThreshold = "multisig-threshold"
	flagNoSort            = "nosort"
	flagHDPath            = "hd-path"

	flagMultiSigKeyName = "multi"
	flagShowMore        = "more"

	flagListNames = "list-names"

	mnemonicEntropySize = 256
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"

	// defaultKeyDBName is the client's subdirectory where keys are stored.
	defaultKeyDBName = "keys"
)

// KeyCommands registers a sub-tree of commands to interact with
// local private key storage.
func KeyCommands(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage your application's keys",
		Long: `Keyring management commands. These keys may be in any format supported by the
Tendermint crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be use only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     https://github.com/KDE/kwallet
    pass        https://www.passwordstore.org/

The pass backend requires GnuPG: https://gnupg.org/
`,
	}

	// support adding Ethereum supported keys
	addCmd := keys.AddKeyCommand()
	addCmd.RunE = runAddCmdPrepare

	showCmd := keys.ShowKeysCmd()
	showCmd.RunE = runShowCmd

	listCmd := keys.ListKeysCmd()
	listCmd.RunE = runListCmd

	cmd.AddCommand(
		keys.MnemonicKeyCommand(),
		addCmd,
		keys.ExportKeyCommand(),
		keys.ImportKeyCommand(),
		listCmd,
		showCmd,
		keys.DeleteKeyCommand(),
		keys.ParseKeyStringCommand(),
		keys.MigrateCommand(),
		UnsafeExportEthKeyCommand(),
		UnsafeImportKeyCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, keyring.BackendOS, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")
	cmd.PersistentFlags().BoolP(flagShowMore, "m", false, "Show more info of account")
	return cmd
}

func runAddCmdPrepare(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(cmd.InOrStdin())
	return runAddCmd(clientCtx, cmd, args, buf)
}

// UnsafeExportEthKeyCommand exports a key with the given name as a private key in hex format.
func UnsafeExportEthKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-export-eth-key [name]",
		Short: "**UNSAFE** Export an Ethereum private key",
		Long:  `**UNSAFE** Export an Ethereum private key unencrypted to use in dev tooling`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())

			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			rootDir, _ := cmd.Flags().GetString(flags.FlagHome)

			kr, err := keyring.New(
				sdk.KeyringServiceName(),
				keyringBackend,
				rootDir,
				inBuf,
				etherminthd.EthSecp256k1Option(),
			)
			if err != nil {
				return err
			}

			decryptPassword := ""
			conf := true

			switch keyringBackend {
			case keyring.BackendFile:
				decryptPassword, err = input.GetPassword(
					"**WARNING this is an unsafe way to export your unencrypted private key**\nEnter key password:",
					inBuf)
			case keyring.BackendOS:
				conf, err = input.GetConfirmation(
					"**WARNING** this is an unsafe way to export your unencrypted private key, are you sure?",
					inBuf, cmd.ErrOrStderr())
			}
			if err != nil || !conf {
				return err
			}

			// Exports private key from keybase using password
			armor, err := kr.ExportPrivKeyArmor(args[0], decryptPassword)
			if err != nil {
				return err
			}

			privKey, algo, err := crypto.UnarmorDecryptPrivKey(armor, decryptPassword)
			if err != nil {
				return err
			}

			if algo != string(hd.Secp256k1Type) {
				return fmt.Errorf("invalid key algorithm, got %s, expected %s", algo, string(hd.Secp256k1Type))
			}

			// Converts key to cosmos secp256k1 implementation
			secp256k1PrivKey, ok := privKey.(*secp256k1.PrivKey)
			if !ok {
				return fmt.Errorf("invalid private key type %T, expected %T", privKey, &secp256k1.PrivKey{})
			}

			key, err := ethcrypto.ToECDSA(secp256k1PrivKey.Bytes())
			if err != nil {
				return err
			}

			// Formats key for output
			privB := ethcrypto.FromECDSA(key)
			keyS := strings.ToUpper(hexutil.Encode(privB)[2:])

			fmt.Println(keyS)

			return nil
		},
	}
}

// UnsafeImportKeyCommand imports private keys from a keyfile.
func UnsafeImportKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-import-eth-key <name> <pk>",
		Short: "**UNSAFE** Import Ethereum private keys into the local keybase",
		Long:  "**UNSAFE** Import a hex-encoded Ethereum private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE:  runImportCmd,
	}
}

func runImportCmd(cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
	rootDir, _ := cmd.Flags().GetString(flags.FlagHome)

	kb, err := keyring.New(
		sdk.KeyringServiceName(),
		keyringBackend,
		rootDir,
		inBuf,
		etherminthd.EthSecp256k1Option(),
	)
	if err != nil {
		return err
	}

	privKey := &secp256k1.PrivKey{
		Key: common.FromHex(args[1]),
	}

	armor := crypto.EncryptArmorPrivKey(privKey, "", "secp256k1")

	return kb.ImportPrivKey(args[0], armor, "")
}

/*
input
	- bip39 mnemonic
	- bip39 passphrase
	- bip44 path
	- local encryption password
output
	- armor encrypted private key (saved to file)
*/
func runAddCmd(ctx client.Context, cmd *cobra.Command, args []string, inBuf *bufio.Reader) error {
	var err error

	name := args[0]
	interactive, _ := cmd.Flags().GetBool(flagInteractive)
	noBackup, _ := cmd.Flags().GetBool(flagNoBackup)
	showMnemonic := !noBackup
	kb := ctx.Keyring
	outputFormat := ctx.OutputFormat

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algoStr, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
	algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
	if err != nil {
		return err
	}

	if dryRun, _ := cmd.Flags().GetBool(flags.FlagDryRun); dryRun {
		// use in memory keybase
		kb = keyring.NewInMemory(etherminthd.EthSecp256k1Option())
	} else {
		_, err = kb.Key(name)
		if err == nil {
			// account exists, ask for user confirmation
			response, err2 := input.GetConfirmation(fmt.Sprintf("override the existing name %s", name), inBuf, cmd.ErrOrStderr())
			if err2 != nil {
				return err2
			}

			if !response {
				return errors.New("aborted")
			}

			err2 = kb.Delete(name)
			if err2 != nil {
				return err2
			}
		}

		multisigKeys, _ := cmd.Flags().GetStringSlice(flagMultisig)
		if len(multisigKeys) != 0 {
			pks := make([]cryptotypes.PubKey, len(multisigKeys))
			multisigThreshold, _ := cmd.Flags().GetInt(flagMultiSigThreshold)
			if err := validateMultisigThreshold(multisigThreshold, len(multisigKeys)); err != nil {
				return err
			}

			for i, keyname := range multisigKeys {
				k, err := kb.Key(keyname)
				if err != nil {
					return err
				}

				pks[i] = k.GetPubKey()
			}

			if noSort, _ := cmd.Flags().GetBool(flagNoSort); !noSort {
				sort.Slice(pks, func(i, j int) bool {
					return bytes.Compare(pks[i].Address(), pks[j].Address()) < 0
				})
			}

			pk := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)
			info, err := kb.SaveMultisig(name, pk)
			if err != nil {
				return err
			}

			return printCreate(cmd, info, false, "", outputFormat)
		}
	}

	pubKeyStr, _ := cmd.Flags().GetString(keys.FlagPublicKey)
	if pubKeyStr != "" {
		var pubkey cryptotypes.PubKey
		pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pubKeyStr)
		if err != nil {
			if err = ctx.Codec.UnmarshalInterfaceJSON([]byte(pubKeyStr), &pubkey); err != nil {
				return err
			}
		}
		info, err := kb.SavePubKey(name, pubkey, algo.Name())
		if err != nil {
			return err
		}

		return printCreate(cmd, info, false, "", outputFormat)
	}

	coinType, _ := cmd.Flags().GetUint32(flagCoinType)
	account, _ := cmd.Flags().GetUint32(flagAccount)
	index, _ := cmd.Flags().GetUint32(flagIndex)
	hdPath, _ := cmd.Flags().GetString(flagHDPath)
	useLedger, _ := cmd.Flags().GetBool(flags.FlagUseLedger)

	if len(hdPath) == 0 {
		hdPath = hd.CreateHDPath(coinType, account, index).String()
	} else if useLedger {
		return errors.New("cannot set custom bip32 path with ledger")
	}

	// If we're using ledger, only thing we need is the path and the bech32 prefix.
	if useLedger {
		bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

		info, err := kb.SaveLedgerKey(name, algo, bech32PrefixAccAddr, coinType, account, index)
		if err != nil {
			return err
		}

		return printCreate(cmd, info, false, "", outputFormat)
	}

	// Get bip39 mnemonic
	var mnemonic, bip39Passphrase string

	recover, _ := cmd.Flags().GetBool(flagRecover)
	if recover {
		mnemonic, err = input.GetString("Enter your bip39 mnemonic", inBuf)
		if err != nil {
			return err
		}

		if !bip39.IsMnemonicValid(mnemonic) {
			return errors.New("invalid mnemonic")
		}
	} else if interactive {
		mnemonic, err = input.GetString("Enter your bip39 mnemonic, or hit enter to generate one.", inBuf)
		if err != nil {
			return err
		}

		if !bip39.IsMnemonicValid(mnemonic) && mnemonic != "" {
			return errors.New("invalid mnemonic")
		}
	}

	if len(mnemonic) == 0 {
		// read entropy seed straight from tmcrypto.Rand and convert to mnemonic
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			return err
		}

		mnemonic, err = bip39.NewMnemonic(entropySeed)
		if err != nil {
			return err
		}
	}

	// override bip39 passphrase
	if interactive {
		bip39Passphrase, err = input.GetString(
			"Enter your bip39 passphrase. This is combined with the mnemonic to derive the seed. "+
				"Most users should just hit enter to use the default, \"\"", inBuf)
		if err != nil {
			return err
		}

		// if they use one, make them re-enter it
		if len(bip39Passphrase) != 0 {
			p2, err := input.GetString("Repeat the passphrase:", inBuf)
			if err != nil {
				return err
			}

			if bip39Passphrase != p2 {
				return errors.New("passphrases don't match")
			}
		}
	}

	info, err := kb.NewAccount(name, mnemonic, bip39Passphrase, hdPath, algo)
	if err != nil {
		return err
	}

	// Recover key from seed passphrase
	if recover {
		// Hide mnemonic from output
		showMnemonic = false
		mnemonic = ""
	}

	return printCreate(cmd, info, showMnemonic, mnemonic, outputFormat)
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

type bechKeyOutFn func(keyInfo cryptokeyring.Info) (cryptokeyring.KeyOutput, error)

// NewLegacyKeyBaseFromDir initializes a legacy keybase at the rootDir directory. Keybase
// options can be applied when generating this new Keybase.
func NewLegacyKeyBaseFromDir(rootDir string, opts ...cryptokeyring.KeybaseOption) (cryptokeyring.LegacyKeybase, error) {
	return getLegacyKeyBaseFromDir(rootDir, opts...)
}

func getLegacyKeyBaseFromDir(rootDir string, opts ...cryptokeyring.KeybaseOption) (cryptokeyring.LegacyKeybase, error) {
	return cryptokeyring.NewLegacy(defaultKeyDBName, filepath.Join(rootDir, "keys"), opts...)
}

func printCreate(cmd *cobra.Command, info keyring.Info, showMnemonic bool, mnemonic, outputFormat string) error {
	switch outputFormat {
	case OutputFormatText:
		cmd.PrintErrln()
		printKeyInfo(cmd.OutOrStdout(), info, keyring.Bech32KeyOutput, outputFormat, false)

		// print mnemonic unless requested not to.
		if showMnemonic {
			fmt.Fprintln(cmd.ErrOrStderr(), "\n**Important** write this mnemonic phrase in a safe place.")
			fmt.Fprintln(cmd.ErrOrStderr(), "It is the only way to recover your account if you ever forget your password.")
			fmt.Fprintln(cmd.ErrOrStderr(), "")
			fmt.Fprintln(cmd.ErrOrStderr(), mnemonic)
		}
	case OutputFormatJSON:
		out, err := keyring.Bech32KeyOutput(info)
		if err != nil {
			return err
		}

		if showMnemonic {
			out.Mnemonic = mnemonic
		}

		jsonString, err := keys.KeysCdc.MarshalJSON(out)
		if err != nil {
			return err
		}

		cmd.Println(string(jsonString))

	default:
		return fmt.Errorf("invalid output format %s", outputFormat)
	}

	return nil
}

func printKeyAddress(w io.Writer, info cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ko.Address)
}

func printPubKey(w io.Writer, info cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ko.PubKey)
}

type KeyOutputV2 struct {
	Name              string                 `json:"name" yaml:"name"`
	Type              string                 `json:"type" yaml:"type"`
	Algo              hd.PubKeyType          `json:"algo" yaml:"algo"`
	Address           string                 `json:"address" yaml:"address"`
	AddressByte       string                 `json:"address_byte" yaml:"address_byte"`
	AddressHex        string                 `json:"address_hex" yaml:"address_hex"`
	PubKey            string                 `json:"pubkey" yaml:"pubkey"`
	ValAddress        string                 `json:"val_address" yaml:"val_address"`
	ValPubKey         string                 `json:"val_pubkey" yaml:"val_pubkey"`
	DecompressAddress string                 `json:"decompress_address" yaml:"decompress_address"`
	DecompressPubkey  string                 `json:"decompress_pubkey" yaml:"decompress_pubkey"`
	Mnemonic          string                 `json:"mnemonic,omitempty" yaml:"mnemonic"`
	Threshold         uint                   `json:"threshold,omitempty" yaml:"threshold"`
	PubKeys           []multisigPubKeyOutput `json:"pubkeys,omitempty" yaml:"pubkeys"`
}

type multisigPubKeyOutput struct {
	Address string `json:"address" yaml:"address"`
	PubKey  string `json:"pubkey" yaml:"pubkey"`
	Weight  uint   `json:"weight" yaml:"weight"`
}

func KeyOutputToV2(v1 cryptokeyring.KeyOutput, info keyring.Info) KeyOutputV2 {
	v2 := KeyOutputV2{
		Name:      v1.Name,
		Type:      v1.Type,
		Address:   v1.Address,
		PubKey:    v1.PubKey,
		Mnemonic:  v1.Mnemonic,
		Threshold: v1.Threshold,
		Algo:      info.GetAlgo(),
	}
	info.GetPubKey().Bytes()
	for _, p := range v1.PubKeys {
		v2.PubKeys = append(v2.PubKeys, multisigPubKeyOutput{
			Address: p.Address,
			PubKey:  p.PubKey,
			Weight:  p.Weight,
		})
	}
	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, v1.PubKey)
	if err != nil {
		panic(err)
	}
	v2.AddressByte = fmt.Sprintf("%v", pubkey.Address().Bytes())
	v2.AddressHex = common.BytesToAddress(pubkey.Address()).Hex()
	valPub, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeValPub, pubkey)
	if err != nil {
		panic(err)
	}
	v2.ValPubKey = valPub
	v2.ValAddress = sdk.ValAddress(pubkey.Address()).String()
	decompressPubkey, err := ethcrypto.DecompressPubkey(pubkey.Bytes())
	if err != nil {
		panic(err)
	}
	v2.DecompressPubkey = hex.EncodeToString(ethcrypto.FromECDSAPub(decompressPubkey))
	v2.DecompressAddress = ethcrypto.PubkeyToAddress(*decompressPubkey).Hex()
	return v2
}

func printKeyInfo(w io.Writer, keyInfo cryptokeyring.Info, bechKeyOut bechKeyOutFn, output string, isShowMore bool) {
	ko, err := bechKeyOut(keyInfo)
	if err != nil {
		panic(err)
	}
	var keyOutput interface{}
	keyOutput = ko
	if isShowMore {
		keyOutput = KeyOutputToV2(ko, keyInfo)
	}
	switch output {
	case OutputFormatText:
		outputTextUnlimitedWidth(w, []interface{}{keyOutput})
	case OutputFormatJSON:
		outputJSON(w, keyOutput)
	}
}

func runShowCmd(cmd *cobra.Command, args []string) (err error) {
	var info keyring.Info
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		info, err = fetchKey(clientCtx.Keyring, args[0])
		if err != nil {
			return fmt.Errorf("%s is not a valid name or address: %v", args[0], err)
		}
	} else {
		pks := make([]cryptotypes.PubKey, len(args))
		for i, keyref := range args {
			info, err := fetchKey(clientCtx.Keyring, keyref)
			if err != nil {
				return fmt.Errorf("%s is not a valid name or address: %v", keyref, err)
			}

			pks[i] = info.GetPubKey()
		}

		multisigThreshold, _ := cmd.Flags().GetInt(flagMultiSigThreshold)
		err = validateMultisigThreshold(multisigThreshold, len(args))
		if err != nil {
			return err
		}

		multikey := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)
		info, err = keyring.NewMultiInfo(flagMultiSigKeyName, multikey)
		if err != nil {
			return err
		}
	}

	isShowAddr, _ := cmd.Flags().GetBool(keys.FlagAddress)
	isShowPubKey, _ := cmd.Flags().GetBool(keys.FlagPublicKey)
	isShowDevice, _ := cmd.Flags().GetBool(keys.FlagDevice)
	showMore, _ := cmd.Flags().GetBool(flagShowMore)

	isOutputSet := false
	tmp := cmd.Flag(cli.OutputFlag)
	if tmp != nil {
		isOutputSet = tmp.Changed
	}

	if isShowAddr && isShowPubKey {
		return errors.New("cannot use both --address and --pubkey at once")
	}

	if isOutputSet && (isShowAddr || isShowPubKey) {
		return errors.New("cannot use --output with --address or --pubkey")
	}

	bechPrefix, _ := cmd.Flags().GetString(keys.FlagBechPrefix)
	bechKeyOut, err := getBechKeyOut(bechPrefix)
	if err != nil {
		return err
	}

	output, _ := cmd.Flags().GetString(cli.OutputFlag)

	switch {
	case isShowAddr:
		printKeyAddress(cmd.OutOrStdout(), info, bechKeyOut)
	case isShowPubKey:
		printPubKey(cmd.OutOrStdout(), info, bechKeyOut)
	default:
		printKeyInfo(cmd.OutOrStdout(), info, bechKeyOut, output, showMore)
	}

	if isShowDevice {
		if isShowPubKey {
			return fmt.Errorf("the device flag (-d) can only be used for addresses not pubkeys")
		}
		if bechPrefix != "acc" {
			return fmt.Errorf("the device flag (-d) can only be used for accounts")
		}

		// Override and show in the device
		if info.GetType() != keyring.TypeLedger {
			return fmt.Errorf("the device flag (-d) can only be used for accounts stored in devices")
		}

		hdpath, err := info.GetPath()
		if err != nil {
			return nil
		}

		return ledger.ShowAddress(*hdpath, info.GetPubKey(), sdk.GetConfig().GetBech32AccountAddrPrefix())
	}

	return nil
}

func fetchKey(kb keyring.Keyring, keyref string) (keyring.Info, error) {
	info, err := kb.Key(keyref)
	if err != nil {
		accAddr, err := sdk.AccAddressFromBech32(keyref)
		if err != nil {
			return info, err
		}

		info, err = kb.KeyByAddress(accAddr)
		if err != nil {
			return info, errors.New("key not found")
		}
	}
	return info, nil
}

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case sdk.PrefixAccount:
		return keyring.Bech32KeyOutput, nil
	case sdk.PrefixValidator:
		return keyring.Bech32ValKeyOutput, nil
	case sdk.PrefixConsensus:
		return keyring.Bech32ConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}

func runListCmd(cmd *cobra.Command, _ []string) error {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}

	infos, err := clientCtx.Keyring.List()
	if err != nil {
		return err
	}

	cmd.SetOut(cmd.OutOrStdout())

	showMore, _ := cmd.Flags().GetBool(flagShowMore)

	if ok, _ := cmd.Flags().GetBool(flagListNames); !ok {
		output, _ := cmd.Flags().GetString(cli.OutputFlag)
		printInfos(cmd.OutOrStdout(), infos, output, showMore)
		return nil
	}
	for _, info := range infos {
		cmd.Println(info.GetName())
	}
	return nil
}

func printInfos(w io.Writer, infos []cryptokeyring.Info, output string, showMore bool) {
	var op interface{}
	if showMore {
		kos, err := Bech32KeysOutputV2(infos)
		if err != nil {
			panic(err)
		}
		op = kos
	} else {
		kos, err := cryptokeyring.Bech32KeysOutput(infos)
		if err != nil {
			panic(err)
		}
		op = kos
	}

	switch output {
	case OutputFormatText:
		outputTextUnlimitedWidth(w, &op)
	case OutputFormatJSON:
		outputJSON(w, op)
	}
}

func outputText(w io.Writer, info interface{}) {
	out, err := yaml.Marshal(&info)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(out))
}

//yamlOutputUnlimitedWidth unlimited emitter best_width, unsafe link https://github.com/go-yaml/yaml/pull/455
func outputTextUnlimitedWidth(w io.Writer, info interface{}) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	v := reflect.ValueOf(enc).Elem().FieldByName("encoder").Elem()
	width := v.FieldByName("emitter").FieldByName("best_width").UnsafeAddr()
	widthPtr := (*int)(unsafe.Pointer(width))
	// don't limit the width
	*widthPtr = -1
	// encode contents
	if err := enc.Encode(info); err != nil {
		panic(err)
	}
	if err := enc.Close(); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, buf.String())
}
func outputJSON(w io.Writer, info interface{}) {
	out, err := keys.KeysCdc.MarshalJSON(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%s", out)
}

// Bech32KeysOutputV2 returns a slice of KeyOutput objects, each with the "acc"
// Bech32 prefixes, given a slice of Info objects. It returns an error if any
// call to Bech32KeyOutput fails.
func Bech32KeysOutputV2(infos []cryptokeyring.Info) ([]KeyOutputV2, error) {
	kos := make([]KeyOutputV2, len(infos))
	for i, info := range infos {
		ko, err := cryptokeyring.Bech32KeyOutput(info)
		if err != nil {
			return nil, err
		}
		kos[i] = KeyOutputToV2(ko, info)
	}

	return kos, nil
}
