package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/functionx/fx-core/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"os"
	"path/filepath"
)

func UpdateValidatorKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-validator [mnemonic]",
		Short: "update validator node private key json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			rootDir := serverCtx.Viper.GetString(flags.FlagHome)
			if len(rootDir) <= 0 {
				return fmt.Errorf("home dir can't be empty")
			}
			valMnemonic := args[0]
			pvKeyFile := serverCtx.Config.PrivValidatorKeyFile()
			_ = os.Remove(pvKeyFile)
			if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
				return err
			}
			pvStateFile := serverCtx.Config.PrivValidatorStateFile()
			if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
				return err
			}
			valPrivKey := privval.NewFilePV(ed25519.GenPrivKeyFromSecret([]byte(valMnemonic)), pvKeyFile, pvStateFile)
			valPrivKey.Save()
			return nil
		},
	}
	return cmd
}

func UpdateNodeKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-node-key",
		Short: "update node key",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			rootDir := serverCtx.Viper.GetString(flags.FlagHome)
			if len(rootDir) <= 0 {
				return fmt.Errorf("home dir can't be empty")
			}
			consMnemonic := app.NewMnemonic()
			if len(args) > 0 {
				consMnemonic = args[0]
			}
			nodeKeyFile := serverCtx.Config.NodeKeyFile()
			if err := tmos.EnsureDir(filepath.Dir(nodeKeyFile), 0777); err != nil {
				return err
			}

			nodeKey := &p2p.NodeKey{PrivKey: ed25519.GenPrivKeyFromSecret([]byte(consMnemonic))}
			return nodeKey.SaveAs(nodeKeyFile)
		},
	}
	return cmd
}
