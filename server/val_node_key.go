package server

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

const (
	flagResetNodeKey = "reset-node-key"
	flagResetPrivKey = "reset-priv-key"
	flagUnsafe       = "unsafe"
	flagKeyType      = "key-type"
)

// UnsafeRestPrivValidatorCmd Reset validator node consensus private key file
//
//gocyclo:ignore
func UnsafeRestPrivValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsafe-reset-priv-validator [secret]",
		Short: "(unsafe) Reset validator node consensus private key file (priv_validator_key.json)",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg.EnsureRoot(serverCtx.Config.RootDir)

			unsafe := serverCtx.Viper.GetBool(flagUnsafe)
			resetKey := serverCtx.Viper.GetBool(flagResetPrivKey)
			if (unsafe && !resetKey) || (!unsafe && resetKey) {
				return fmt.Errorf("the flags %s and %s must be used together", flagUnsafe, flagResetPrivKey)
			}
			if !unsafe && !resetKey {
				buf := bufio.NewReader(cmd.InOrStdin())
				if yes, err := input.GetConfirmation("WARNING: The consensus private key of the node will be replaced. Ensure that the backup is complete. USE AT YOUR OWN RISK. Continue?", buf, cmd.ErrOrStderr()); err != nil {
					return err
				} else if !yes {
					return nil
				}
			}

			secret := newMnemonic()
			if len(args) > 0 {
				secret = args[0]
			} else if !unsafe || !resetKey {
				return fmt.Errorf("invalid params, please use reset-priv-key --unsafe")
			}
			if len(secret) < 32 {
				return fmt.Errorf("secret contains less than 32 characters")
			}
			pvKeyFile := serverCtx.Config.PrivValidatorKeyFile()
			if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0o777); err != nil {
				return err
			}
			pvStateFile := serverCtx.Config.PrivValidatorStateFile()
			if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0o777); err != nil {
				return err
			}
			keyType := serverCtx.Viper.GetString(flagKeyType)
			switch keyType {
			case "ed25519":
				_ = os.Remove(pvKeyFile)
				valPrivKey := privval.NewFilePV(ed25519.GenPrivKeyFromSecret([]byte(secret)), pvKeyFile, pvStateFile)
				valPrivKey.Save()
			case "secp256k1":
				_ = os.Remove(pvKeyFile)
				pk := secp256k1.GenPrivKeySecp256k1([]byte(secret))
				valPrivKey := privval.NewFilePV(pk, pvKeyFile, pvStateFile)
				valPrivKey.Save()
			default:
				return fmt.Errorf("invalid key type: %s", keyType)
			}
			return nil
		},
	}
	cmd.Flags().Bool(flagResetPrivKey, false, "Reset ed25519 private key. Requires --unsafe.")
	cmd.Flags().Bool(flagUnsafe, false, "Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.")
	cmd.Flags().String(flagKeyType, "ed25519", "Private key type, ed25519 or secp256k1.")
	return cmd
}

func UnsafeResetNodeKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsafe-reset-node-key [secret]",
		Short: "(unsafe) reset node key file (node_key.json)",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg.EnsureRoot(serverCtx.Config.RootDir)

			unsafe := serverCtx.Viper.GetBool(flagUnsafe)
			resetKey := serverCtx.Viper.GetBool(flagResetNodeKey)
			if (unsafe && !resetKey) || (!unsafe && resetKey) {
				return fmt.Errorf("the flags %s and %s must be used together", flagUnsafe, flagResetNodeKey)
			}
			if !unsafe && !resetKey {
				buf := bufio.NewReader(cmd.InOrStdin())
				if yes, err := input.GetConfirmation("WARNING: The node key of the node will be replaced. Ensure that the backup is complete. USE AT YOUR OWN RISK. Continue?", buf, cmd.ErrOrStderr()); err != nil {
					return err
				} else if !yes {
					return nil
				}
			}
			secret := newMnemonic()
			if len(args) > 0 {
				secret = args[0]
			} else if !unsafe || !resetKey {
				return fmt.Errorf("invalid params. please use --reset-node-key --unsafe")
			}
			if len(secret) < 32 {
				return fmt.Errorf("secret contains less than 32 characters")
			}
			nodeKeyFile := serverCtx.Config.NodeKeyFile()
			_ = os.Remove(nodeKeyFile)
			if err := tmos.EnsureDir(filepath.Dir(nodeKeyFile), 0o777); err != nil {
				return err
			}

			nodeKey := &p2p.NodeKey{PrivKey: ed25519.GenPrivKeyFromSecret([]byte(secret))}
			return nodeKey.SaveAs(nodeKeyFile)
		},
	}
	cmd.Flags().Bool(flagResetNodeKey, false, "Reset node key. Requires --unsafe.")
	cmd.Flags().Bool(flagUnsafe, false, "Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.")
	return cmd
}

func newMnemonic() string {
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
