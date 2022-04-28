package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/functionx/fx-core/crypto/ethsecp256k1"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
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
		exportKeyCommand(),
		importKeyCommand(),
		keys.ListKeysCmd(),
		keys.ShowKeysCmd(),
		keys.DeleteKeyCommand(),
		keys.MigrateCommand(),
		parseAddressCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")

	return cmd
}

const (
	flagUnarmoredHex = "unarmored-hex"
	flagUnsafe       = "unsafe"
	flagASCIIArmored = "ascii-armor"
)

// exportKeyCommand exports private keys from the key store.
func exportKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export private keys",
		Long: `Export a private key from the local keyring in encrypted format.

When both the --unarmored-hex and --unsafe flags are selected, cryptographic
private key material is exported in an INSECURE fashion that is designed to
allow users to import their keys in hot wallets. This feature is for advanced
users only that are confident about how to handle private keys work and are
FULLY AWARE OF THE RISKS. If you are unsure, you may want to do some research
and export your keys in encrypted format.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			buf := bufio.NewReader(clientCtx.Input)
			unarmored, _ := cmd.Flags().GetBool(flagUnarmoredHex)
			unsafe, _ := cmd.Flags().GetBool(flagUnsafe)

			if unarmored && unsafe {
				return exportUnsafeUnarmored(cmd, args[0], buf, clientCtx.Keyring)
			} else if unarmored || unsafe {
				return fmt.Errorf("the flags %s and %s must be used together", flagUnsafe, flagUnarmoredHex)
			}

			encryptPassword, err := input.GetPassword("Enter passphrase to encrypt the exported key:", buf)
			if err != nil {
				return err
			}

			asciiArmored, _ := cmd.Flags().GetBool(flagASCIIArmored)
			if asciiArmored {
				armored, err := clientCtx.Keyring.ExportPrivKeyArmor(args[0], encryptPassword)
				if err != nil {
					return err
				}

				cmd.Println(armored)

				return nil
			}

			hexPrivKey, err := keyring.NewUnsafe(clientCtx.Keyring).UnsafeExportPrivKeyHex(args[0])
			if err != nil {
				return err
			}
			priv, err := ethcrypto.HexToECDSA(hexPrivKey)
			if err != nil {
				return err
			}
			key := &keystore.Key{
				PrivateKey: priv,
				Address:    ethcrypto.PubkeyToAddress(priv.PublicKey),
			}
			keyjson, err := keystore.EncryptKey(key, encryptPassword, keystore.StandardScryptN, keystore.StandardScryptP)
			if err != nil {
				return err
			}

			cmd.Println(string(keyjson))

			return nil
		},
	}

	cmd.Flags().Bool(flagUnarmoredHex, false, "Export unarmored hex privkey. Requires --unsafe.")
	cmd.Flags().Bool(flagUnsafe, false, "Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.")
	cmd.Flags().Bool(flagASCIIArmored, false, "Enable ASCII-armored encrypted format")
	return cmd
}

func exportUnsafeUnarmored(cmd *cobra.Command, uid string, buf *bufio.Reader, kr keyring.Keyring) error {
	// confirm deletion, unless -y is passed
	if yes, err := input.GetConfirmation("WARNING: The private key will be exported as an unarmored hexadecimal string. USE AT YOUR OWN RISK. Continue?", buf, cmd.ErrOrStderr()); err != nil {
		return err
	} else if !yes {
		return nil
	}

	hexPrivKey, err := keyring.NewUnsafe(kr).UnsafeExportPrivKeyHex(uid)
	if err != nil {
		return err
	}

	cmd.Print(hexPrivKey)

	return nil
}

func importKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <name> <keyfile>",
		Short: "Import private keys into the local keybase",
		Long:  "Import a ASCII armored or ethereum keystore or unencrypted private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			buf := bufio.NewReader(clientCtx.Input)

			bz, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}
			if len(bz) == 64 {
				priv, err := ethcrypto.HexToECDSA(string(bz))
				if err != nil {
					return err
				}
				armoPrivKey := crypto.EncryptArmorPrivKey(&ethsecp256k1.PrivKey{Key: ethcrypto.FromECDSA(priv)}, "", ethsecp256k1.KeyType)
				return clientCtx.Keyring.ImportPrivKey(args[0], armoPrivKey, "")
			}
			passphrase, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
			if err != nil {
				return err
			}
			key, err := keystore.DecryptKey(bz, passphrase)
			if err == nil {
				armoPrivKey := crypto.EncryptArmorPrivKey(&ethsecp256k1.PrivKey{Key: ethcrypto.FromECDSA(key.PrivateKey)}, "", ethsecp256k1.KeyType)
				return clientCtx.Keyring.ImportPrivKey(args[0], armoPrivKey, "")
			}

			return clientCtx.Keyring.ImportPrivKey(args[0], string(bz), passphrase)
		},
	}
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
