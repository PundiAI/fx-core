package cli

import (
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

// NewPruneCmd creates a command to prune cosmos-sdk and tendermint state.
func NewPruneCmd(defaultNodeHome string) *cobra.Command {
	pruneCmd := &cobra.Command{
		Use:                "prune",
		Short:              "Prune cosmos-sdk and tendermint state",
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	pruneCmd.AddCommand(
		LastBlockState(),
	)

	pruneCmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	return pruneCmd
}

func LastBlockState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-block-state",
		Short: "prune last block state",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			cfg := serverCtx.Config
			blockState, err := LoadBlockState(cfg.RootDir, GetAppDBBackend(serverCtx.Viper))
			if err != nil {
				return err
			}
			pruneLastHeight, err := blockState.PruneLastBlock()
			if err != nil {
				return fmt.Errorf("failed to prune last block state: %w", err)
			}
			fmt.Printf("prune block state to height %d", pruneLastHeight-1)
			return nil
		},
	}
	return cmd
}

func LoadBlockState(rootDir string, backendType dbm.BackendType) (*store.BlockStore, error) {
	dataDir := filepath.Join(rootDir, "data")
	blockStoreDB, err := dbm.NewDB("blockstore", backendType, dataDir)
	if err != nil {
		return nil, err
	}
	return store.NewBlockStore(blockStoreDB), nil
}

// GetAppDBBackend gets the backend type to use for the application DBs.
func GetAppDBBackend(opts types.AppOptions) dbm.BackendType {
	rv := cast.ToString(opts.Get("app-db-backend"))
	if len(rv) == 0 {
		rv = sdk.DBBackend
	}
	if len(rv) == 0 {
		rv = cast.ToString(opts.Get("db-backend"))
	}
	if len(rv) != 0 {
		return dbm.BackendType(rv)
	}
	return dbm.GoLevelDBBackend
}
