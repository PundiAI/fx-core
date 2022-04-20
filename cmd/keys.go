package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/functionx/fx-core/crypto/ethsecp256k1"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	etherminthd "github.com/functionx/fx-core/crypto/hd"
)

// Commands registers a sub-tree of commands to interact with
// local private key storage.
func keyCommands(defaultNodeHome string) *cobra.Command {
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

	addKeyCmd := keys.AddKeyCommand()
	algoFlag := addKeyCmd.Flag(flags.FlagKeyAlgorithm)
	algoFlag.DefValue = ethsecp256k1.KeyType
	_ = algoFlag.Value.Set(algoFlag.DefValue)

	cmd.AddCommand(
		keys.MnemonicKeyCommand(),
		addKeyCmd,
		keys.ExportKeyCommand(),
		keys.ImportKeyCommand(),
		keys.ListKeysCmd(),
		keys.ShowKeysCmd(),
		keys.DeleteKeyCommand(),
		keys.MigrateCommand(),
		parseAddressCommand(),
		unsafeExportEthKeyCommand(),
		unsafeImportKeyCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")

	return cmd
}

func parseAddressCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse [address]",
		Short: "Parse address from hex to bech32 and vice versa",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			addrString := args[0]

			// try hex, then bech32
			var addr []byte
			addr, err = hex.DecodeString(addrString)
			if err != nil {
				var err2 error
				addr, err2 = sdk.AccAddressFromBech32(addrString)
				if err2 != nil {
					var err3 error
					addr, err3 = sdk.ValAddressFromBech32(addrString)

					if err3 != nil {
						return fmt.Errorf("expected hex or bech32. Got errors: hex: %v, bech32 acc: %v, bech32 val: %v", err, err2, err3)

					}
				}
			}
			prefix, err := cmd.Flags().GetString("prefix")
			if err != nil {
				return err
			}
			accAddress, err := bech32.ConvertAndEncode(prefix, addr)
			if err != nil {
				return err
			}
			valAddress, err := bech32.ConvertAndEncode(prefix+sdk.PrefixValidator+sdk.PrefixOperator, addr)
			if err != nil {
				return err
			}

			data, err := json.Marshal(map[string]interface{}{
				"base64_address": addr,
				"hex_address":    hex.EncodeToString(addr),
				"eip55_address":  common.BytesToAddress(addr).String(),
				"acc_address":    accAddress,
				"val_address":    valAddress,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(data)
		},
	}
	cmd.Flags().String("prefix", "fx", "custom address prefix")
	return cmd
}

// unsafeExportEthKeyCommand exports a key with the given name as a private key in hex format.
func unsafeExportEthKeyCommand() *cobra.Command {
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

// unsafeImportKeyCommand imports private keys from a keyfile.
func unsafeImportKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-import-eth-key <name> <pk>",
		Short: "**UNSAFE** Import Ethereum private keys into the local keybase",
		Long:  "**UNSAFE** Import a hex-encoded Ethereum private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}
}
