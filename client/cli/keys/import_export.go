package keys

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/spf13/cobra"
)

const (
	flagUnarmoredHex = "unarmored-hex"
	flagUnsafe       = "unsafe"
	flagASCIIArmored = "ascii-armor"
)

// ExportKeyCommand exports private keys from the key store.
func ExportKeyCommand() *cobra.Command {
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

			privKey, err := clientCtx.Keyring.(unsafeExporter).ExportPrivateKeyObject(args[0])
			if err != nil {
				return err
			}

			hexPrivKey := hex.EncodeToString(privKey.Bytes())
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

	priv, err := kr.(unsafeExporter).ExportPrivateKeyObject(uid)
	if err != nil {
		return err
	}
	hexPrivKey := hex.EncodeToString(priv.Bytes())

	cmd.Print(hexPrivKey)

	return nil
}

func ImportKeyCommand() *cobra.Command {
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

			bz, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			// os.ReadFile read all data, contains line break at the end of a file
			bz = bytes.TrimPrefix(bytes.TrimSuffix(bz, []byte{'\n'}), []byte("0x"))

			if len(bz) == 64 {
				priv, err := ethcrypto.HexToECDSA(string(bz))
				if err != nil {
					return err
				}
				algoStr, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
				var armoPrivKey string
				switch hd.PubKeyType(algoStr) {
				case hd.Secp256k1Type:
					armoPrivKey = crypto.EncryptArmorPrivKey(&secp256k1.PrivKey{Key: ethcrypto.FromECDSA(priv)}, "", string(hd.Secp256k1Type))
				case hd2.EthSecp256k1Type:
					armoPrivKey = crypto.EncryptArmorPrivKey(&ethsecp256k1.PrivKey{Key: ethcrypto.FromECDSA(priv)}, "", string(hd2.EthSecp256k1Type))
				default:
					return fmt.Errorf("provided algorithm %q is not supported", algoStr)
				}

				return clientCtx.Keyring.ImportPrivKey(args[0], armoPrivKey, "")
			}
			passphrase, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
			if err != nil && !strings.Contains(err.Error(), "password must be at least") {
				return err
			}
			key, err := keystore.DecryptKey(bz, passphrase)
			if err == nil {
				algoStr, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
				var armoPrivKey string
				if hd.PubKeyType(algoStr) == hd.Secp256k1Type {
					armoPrivKey = crypto.EncryptArmorPrivKey(&secp256k1.PrivKey{Key: ethcrypto.FromECDSA(key.PrivateKey)}, "", string(hd.Secp256k1Type))
				} else if hd.PubKeyType(algoStr) == hd2.EthSecp256k1Type {
					armoPrivKey = crypto.EncryptArmorPrivKey(&ethsecp256k1.PrivKey{Key: ethcrypto.FromECDSA(key.PrivateKey)}, "", string(hd2.EthSecp256k1Type))
				} else {
					return fmt.Errorf("provided algorithm %q is not supported", algoStr)
				}
				return clientCtx.Keyring.ImportPrivKey(args[0], armoPrivKey, "")
			}

			return clientCtx.Keyring.ImportPrivKey(args[0], string(bz), passphrase)
		},
	}
	cmd.Flags().String(flags.FlagKeyAlgorithm, ethsecp256k1.KeyType, "Key signing algorithm to generate keys for")
	return cmd
}

// unsafeExporter is implemented by key stores that support unsafe export
// of private keys' material.
type unsafeExporter interface {
	// ExportPrivateKeyObject returns a private key in unarmored format.
	ExportPrivateKeyObject(uid string) (cryptotypes.PrivKey, error)
}
