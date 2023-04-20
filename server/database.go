package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v4/app"
)

type Database struct {
	blockStore *store.BlockStore
	stateDB    dbm.DB
	appDB      tmdb.DB
	appStore   *rootmulti.Store
	storeKeys  map[string]*storetypes.KVStoreKey
}

func NewDatabase(rootDir string, dbType tmdb.BackendType) (*Database, error) {
	dataDir := filepath.Join(rootDir, "data")
	if !Exists(filepath.Join(dataDir, fmt.Sprintf("%s.db", BlockDBName))) ||
		!Exists(filepath.Join(dataDir, fmt.Sprintf("%s.db", StateDBName))) ||
		!Exists(filepath.Join(dataDir, fmt.Sprintf("%s.db", AppDBName))) {
		fmt.Println("\tWarning: Not found data file!")
		return nil, nil
	}
	blockStoreDB, err := dbm.NewDB(BlockDBName, dbm.BackendType(dbType), dataDir)
	if err != nil {
		return nil, err
	}
	blockStore := store.NewBlockStore(blockStoreDB)

	stateDB, err := dbm.NewDB(StateDBName, dbm.BackendType(dbType), dataDir)
	if err != nil {
		return nil, err
	}
	appDB, err := tmdb.NewDB(AppDBName, dbType, dataDir)
	if err != nil {
		return nil, err
	}
	storeKeys := sdk.NewKVStoreKeys(
		upgradetypes.StoreKey,
	)
	appStore := rootmulti.NewStore(appDB, log.NewNopLogger())
	for _, storeKey := range storeKeys {
		appStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	}
	if err = appStore.LoadLatestVersion(); err != nil {
		return nil, err
	}

	return &Database{
		blockStore: blockStore,
		appDB:      appDB,
		stateDB:    stateDB,
		appStore:   appStore,
		storeKeys:  storeKeys,
	}, err
}

func (d *Database) Close() {
	_ = d.blockStore.Close()
	_ = d.stateDB.Close()
	_ = d.appDB.Close()
}

func (d *Database) GetChainId() (string, error) {
	return d.blockStore.LoadBaseMeta().Header.ChainID, nil
}

func (d *Database) GetBlockHeight() (int64, error) {
	return d.blockStore.Height(), nil
}

func (d *Database) GetSyncing() (bool, error) {
	return true, nil
}

func (d *Database) GetNodeInfo() (*tmservice.VersionInfo, error) {
	stateStore := sm.NewStore(d.stateDB, sm.StoreOptions{DiscardABCIResponses: false})
	genesisDocKey := []byte("genesisDoc")
	b, err := d.stateDB.Get(genesisDocKey)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, errors.New("genesis doc not found")
	}
	var genDoc *tmtypes.GenesisDoc
	err = tmjson.Unmarshal(b, &genDoc)
	if err != nil {
		panic(fmt.Sprintf("Failed to load genesis doc due to unmarshaling error: %v (bytes: %X)", err, b))
	}
	state, err := stateStore.LoadFromDBOrGenesisDoc(genDoc)
	if err != nil {
		return nil, err
	}
	return &tmservice.VersionInfo{
		Version: state.Version.GetSoftware(),
	}, nil
}

func (d *Database) CurrentPlan() (*upgradetypes.Plan, error) {
	kvStore := d.appStore.GetKVStore(d.storeKeys[upgradetypes.StoreKey])
	bz := kvStore.Get(upgradetypes.PlanKey())
	if bz == nil {
		return nil, nil
	}
	var plan upgradetypes.Plan
	app.MakeEncodingConfig().Codec.MustUnmarshal(bz, &plan)
	return &plan, nil
}

func (d *Database) GetValidators() ([]stakingtypes.Validator, error) {
	stateStore := sm.NewStore(d.stateDB, sm.StoreOptions{DiscardABCIResponses: false})
	genesisDocKey := []byte("genesisDoc")
	b, err := d.stateDB.Get(genesisDocKey)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, errors.New("genesis doc not found")
	}
	var genDoc *tmtypes.GenesisDoc
	err = tmjson.Unmarshal(b, &genDoc)
	if err != nil {
		panic(fmt.Sprintf("Failed to load genesis doc due to unmarshaling error: %v (bytes: %X)", err, b))
	}
	state, err := stateStore.LoadFromDBOrGenesisDoc(genDoc)
	if err != nil {
		return nil, err
	}
	validators := make([]stakingtypes.Validator, 0)
	for _, validator := range state.Validators.Validators {
		validators = append(validators, stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(validator.Address.Bytes()).String(),
		})
	}
	return validators, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
