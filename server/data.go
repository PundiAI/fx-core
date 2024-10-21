package server

import (
	"fmt"
	"sync"
	"time"

	cmtdbm "github.com/cometbft/cometbft-db"
	sm "github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/store"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	flagHeight    = "height"
	flagPruning   = "enable_pruning"
	flagDBBackend = "db_backend"

	BlockDBName = "blockstore"
	StateDBName = "state"
	AppDBName   = "application"
)

func DataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Modify data or query data in database",
	}

	cmd.AddCommand(
		dataQueryBlockCmd(),
		dataPruningCmd(),
	)

	cmd.PersistentFlags().String(flagDBBackend, "goleveldb", "Database backend: goleveldb")
	return cmd
}

func dataQueryBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block",
		Short: "Query blocks heights in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			ctx.Config.DBBackend = ctx.Viper.GetString(flagDBBackend)
			blockStoreDB, err := cmtdbm.NewDB(BlockDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}
			blockStore := store.NewBlockStore(blockStoreDB)
			fmt.Printf("[%d ~ %d]\n", blockStore.Base(), blockStore.Height())

			return nil
		},
	}
}

func dataPruningCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune-compact",
		Short: "Prune and Compact blocks and application states",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}

			backend, err := cmd.Flags().GetString(flagDBBackend)
			if err != nil {
				return err
			}
			if backend != string(cmtdbm.GoLevelDBBackend) {
				return fmt.Errorf("nonsupport db_backend %s, expected <%s>", backend, cmtdbm.GoLevelDBBackend)
			}
			return nil
		},
	}

	cmd.AddCommand(
		pruneAllCmd(),
		pruneAppCmd(),
		pruneBlockCmd(),
	)

	cmd.PersistentFlags().Int64P(flagHeight, "r", 100800, "Removes block or state up to (but not including) a height")
	cmd.PersistentFlags().BoolP(flagPruning, "p", false, "Enable pruning")
	return cmd
}

func pruneAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Compact both application states and blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)

			blockStoreDB, err := cmtdbm.NewDB(BlockDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}
			stateDB, err := cmtdbm.NewDB(StateDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}
			appDB, err := cmtdbm.NewDB(AppDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}

			if err = pruneBlockData(cmd, blockStoreDB, stateDB); err != nil {
				return err
			}

			fmt.Println("--------- compact start... ---------")
			var wg sync.WaitGroup
			wg.Add(3)
			go compactDB(blockStoreDB, BlockDBName, &wg)
			go compactDB(stateDB, StateDBName, &wg)
			go compactDB(appDB, AppDBName, &wg)
			wg.Wait()
			fmt.Println("--------- compact end!!!   ---------")

			return nil
		},
	}

	return cmd
}

func pruneAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Compact application state",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)

			appDB, err := cmtdbm.NewDB(AppDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}
			fmt.Println("--------- compact start ---------")
			var wg sync.WaitGroup
			wg.Add(1)
			compactDB(appDB, AppDBName, &wg)
			fmt.Println("--------- compact end ---------")

			return nil
		},
	}

	return cmd
}

func pruneBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block",
		Short: "Compact while pruning blocks and states",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)

			blockStoreDB, err := cmtdbm.NewDB(BlockDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}
			stateDB, err := cmtdbm.NewDB(StateDBName, cmtdbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
			if err != nil {
				return err
			}

			if err = pruneBlockData(cmd, blockStoreDB, stateDB); err != nil {
				return err
			}

			fmt.Println("--------- compact start... ---------")
			var wg sync.WaitGroup
			wg.Add(2)
			go compactDB(blockStoreDB, BlockDBName, &wg)
			go compactDB(stateDB, StateDBName, &wg)
			wg.Wait()
			fmt.Println("--------- compact end!!!   ---------")

			return nil
		},
	}

	return cmd
}

func pruneBlockData(cmd *cobra.Command, blockStoreDB, stateDB cmtdbm.DB) error {
	enablePruning, err := cmd.Flags().GetBool(flagPruning)
	if err != nil {
		return err
	}
	if !enablePruning {
		return nil
	}

	fullHeight, err := cmd.Flags().GetInt64(flagHeight)
	if err != nil {
		return err
	}
	if fullHeight == 0 {
		return nil
	}

	baseHeight, currentHeight := getBlockInfo(blockStoreDB)

	toHeight := currentHeight - fullHeight
	if toHeight <= baseHeight {
		fmt.Printf("base height greater than equal to full height, skip pruning!baseHeight:%d,currentHeight:%d,toHeight:%d\n", baseHeight, currentHeight, toHeight)
		return nil
	}
	fmt.Printf("--------- pruning start... from:%d,to:%d,size:%d---------\n", baseHeight, toHeight, fullHeight)
	var wg sync.WaitGroup
	wg.Add(2)
	start := time.Now()
	go pruneBlocks(blockStoreDB, baseHeight, toHeight, &wg)
	go pruneStates(stateDB, baseHeight, toHeight, &wg)
	wg.Wait()
	fmt.Printf("--------- pruning end!!!   cast:%v---------\n", time.Since(start))
	return nil
}

// pruneBlocks deletes blocks between the given heights (including from, excluding to).
func pruneBlocks(blockStoreDB cmtdbm.DB, baseHeight, toHeight int64, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Prune blocks start from %d to %d, pruneSize: %d\n", baseHeight, toHeight, toHeight-baseHeight)
	start := time.Now()
	_, _, err := store.NewBlockStore(blockStoreDB).PruneBlocks(toHeight, sm.State{})
	if err != nil {
		panic(fmt.Errorf("failed to prune block store: %w", err))
	}

	baseHeightAfter, currentHeightAfter := getBlockInfo(blockStoreDB)
	fmt.Printf("Prune blocks done new base: %d, height:%d, cost:%v\n", baseHeightAfter, currentHeightAfter, time.Since(start))
}

// pruneStates deletes states between the given heights (including from, excluding to).
func pruneStates(stateDB cmtdbm.DB, from, to int64, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Prune states start from %d to %d, pruneSize: %d\n", from, to, to-from)

	start := time.Now()
	stateStore := sm.NewStore(stateDB, sm.StoreOptions{
		DiscardABCIResponses: false,
	})
	if err := stateStore.PruneStates(from, to, 0); err != nil {
		panic(fmt.Errorf("failed to prune state database: %w", err))
	}

	fmt.Printf("Prune states done in %v \n", time.Since(start))
}

func compactDB(db cmtdbm.DB, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Compact %s... \n", name)
	start := time.Now()

	if ldb, ok := db.(*cmtdbm.GoLevelDB); ok {
		if err := ldb.DB().CompactRange(util.Range{}); err != nil {
			panic(err)
		}
	}

	fmt.Printf("Compact %s done in %v \n", name, time.Since(start))
}

func getBlockInfo(blockStoreDB cmtdbm.DB) (baseHeight, height int64) {
	blockStore := store.NewBlockStore(blockStoreDB)
	return blockStore.Base(), blockStore.Height()
}
