package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/input"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

const (
	flagUpdateNodeKey = "update-node-key"
	flagUpdatePrivKey = "update-priv-key"
	flagUnsafe        = "unsafe"
)

func UpdateValidatorKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-validator [secret]",
		Short: "update validator node consensus private key file (.fxcore/config/priv_validator_key.json)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			rootDir := serverCtx.Viper.GetString(flags.FlagHome)
			if len(rootDir) <= 0 {
				return errors.New("please set the home directory")
			}
			unsafe := serverCtx.Viper.GetBool(flagUnsafe)
			updateKey := serverCtx.Viper.GetBool(flagUpdatePrivKey)
			if (unsafe && !updateKey) || (!unsafe && updateKey) {
				return fmt.Errorf("the flags %s and %s must be used together", flagUnsafe, flagUpdatePrivKey)
			}
			if !unsafe && !updateKey {
				buf := bufio.NewReader(cmd.InOrStdin())
				if yes, err := input.GetConfirmation("WARNING: The consensus private key of the node will be replaced. Ensure that the backup is complete. USE AT YOUR OWN RISK. Continue?", buf, cmd.ErrOrStderr()); err != nil {
					return err
				} else if !yes {
					return nil
				}
			}

			secret := args[0]
			if len(secret) < 32 {
				return fmt.Errorf("secret contains less than 32 characters")
			}
			pvKeyFile := serverCtx.Config.PrivValidatorKeyFile()
			_ = os.Remove(pvKeyFile)
			if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
				return err
			}
			pvStateFile := serverCtx.Config.PrivValidatorStateFile()
			if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
				return err
			}
			valPrivKey := privval.NewFilePV(ed25519.GenPrivKeyFromSecret([]byte(secret)), pvKeyFile, pvStateFile)
			valPrivKey.Save()
			return nil
		},
	}
	cmd.Flags().Bool(flagUpdatePrivKey, false, "Update ed25519 private key. Requires --unsafe.")
	cmd.Flags().Bool(flagUnsafe, false, "Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.")
	return cmd
}

func UpdateNodeKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-node-key [secret]",
		Short: "update node key file (fxcore/config/node_key.json)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			rootDir := serverCtx.Viper.GetString(flags.FlagHome)
			if len(rootDir) <= 0 {
				return errors.New("please set the home directory")
			}
			unsafe := serverCtx.Viper.GetBool(flagUnsafe)
			updateKey := serverCtx.Viper.GetBool(flagUpdateNodeKey)
			if (unsafe && !updateKey) || (!unsafe && updateKey) {
				return fmt.Errorf("the flags %s and %s must be used together", flagUnsafe, flagUpdateNodeKey)
			}
			if !unsafe && !updateKey {
				buf := bufio.NewReader(cmd.InOrStdin())
				if yes, err := input.GetConfirmation("WARNING: The node key of the node will be replaced. Ensure that the backup is complete. USE AT YOUR OWN RISK. Continue?", buf, cmd.ErrOrStderr()); err != nil {
					return err
				} else if !yes {
					return nil
				}
			}
			secret := args[0]
			if len(secret) < 32 {
				return fmt.Errorf("secret contains less than 32 characters")
			}
			nodeKeyFile := serverCtx.Config.NodeKeyFile()
			_ = os.Remove(nodeKeyFile)
			if err := tmos.EnsureDir(filepath.Dir(nodeKeyFile), 0777); err != nil {
				return err
			}

			nodeKey := &p2p.NodeKey{PrivKey: ed25519.GenPrivKeyFromSecret([]byte(secret))}
			return nodeKey.SaveAs(nodeKeyFile)
		},
	}
	cmd.Flags().Bool(flagUpdateNodeKey, false, "Update node key. Requires --unsafe.")
	cmd.Flags().Bool(flagUnsafe, false, "Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.")
	return cmd
}
