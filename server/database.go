package server

import (
	"errors"
	"path/filepath"

	"cosmossdk.io/log"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtdbm "github.com/cometbft/cometbft-db"
	tmcfg "github.com/cometbft/cometbft/config"
	sm "github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/store"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type Database struct {
	genesisFile string
	blockStore  *store.BlockStore
	stateStore  sm.Store
	appDB       dbm.DB
	appStore    *rootmulti.Store
	storeKeys   map[string]*storetypes.KVStoreKey
}

func NewDatabase(cfg *tmcfg.Config, modules ...string) (*Database, error) {
	dataDir := filepath.Join(cfg.RootDir, "data")

	blockStoreDB, err := cmtdbm.NewDB(BlockDBName, cmtdbm.BackendType(cfg.DBBackend), dataDir)
	if err != nil {
		return nil, err
	}
	blockStore := store.NewBlockStore(blockStoreDB)

	stateDB, err := cmtdbm.NewDB(StateDBName, cmtdbm.BackendType(cfg.DBBackend), dataDir)
	if err != nil {
		return nil, err
	}
	stateStore := sm.NewStore(stateDB, sm.StoreOptions{DiscardABCIResponses: false})

	appDB, err := dbm.NewDB(AppDBName, dbm.BackendType(cfg.DBBackend), dataDir)
	if err != nil {
		return nil, err
	}
	storeKeys := storetypes.NewKVStoreKeys(append(modules, upgradetypes.StoreKey)...)
	appStore := rootmulti.NewStore(appDB, log.NewNopLogger(), metrics.NewNoOpMetrics())

	for _, storeKey := range storeKeys {
		appStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	}
	if err = appStore.LoadLatestVersion(); err != nil {
		return nil, err
	}
	return &Database{
		genesisFile: cfg.GenesisFile(),
		blockStore:  blockStore,
		appDB:       appDB,
		stateStore:  stateStore,
		appStore:    appStore,
		storeKeys:   storeKeys,
	}, err
}

func (d *Database) AppDB() dbm.DB {
	return d.appDB
}

func (d *Database) AppStore() *rootmulti.Store {
	return d.appStore
}

func (d *Database) StateStore() sm.Store {
	return d.stateStore
}

func (d *Database) BlockStore() *store.BlockStore {
	return d.blockStore
}

func (d *Database) Close() {
	_ = d.appDB.Close()
	_ = d.blockStore.Close()
	_ = d.stateStore.Close()
}

func (d *Database) GetChainId() (string, error) {
	meta := d.blockStore.LoadBaseMeta()
	if meta == nil {
		return "", errors.New("not found chain id")
	}
	return meta.Header.ChainID, nil
}

func (d *Database) GetBlockHeight() (int64, error) {
	return d.blockStore.Height(), nil
}

func (d *Database) GetSyncing() (bool, error) {
	return true, nil
}

func (d *Database) GetNodeInfo() (*cmtservice.VersionInfo, error) {
	state, err := d.GetState()
	if err != nil {
		return nil, err
	}
	return &cmtservice.VersionInfo{
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
	if err := proto.Unmarshal(bz, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (d *Database) GetConsensusValidators() ([]*cmtservice.Validator, error) {
	state, err := d.GetState()
	if err != nil {
		return nil, err
	}
	if len(state.Validators.Validators) == 0 {
		return nil, nil
	}
	validators := make([]*cmtservice.Validator, 0, len(state.Validators.Validators))
	for _, val := range state.Validators.Validators {
		validator, err := toValidator(val)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validator)
	}
	return validators, nil
}

func (d *Database) GetState() (sm.State, error) {
	state, err := d.stateStore.LoadFromDBOrGenesisFile(d.genesisFile)
	if err != nil {
		return sm.State{}, err
	}
	return state, nil
}

func toValidator(validator *tmtypes.Validator) (*cmtservice.Validator, error) {
	pk, err := cryptocodec.FromCmtPubKeyInterface(validator.PubKey)
	if err != nil {
		return nil, err
	}
	anyPub, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return nil, err
	}
	return &cmtservice.Validator{
		Address:          sdk.ConsAddress(validator.Address.Bytes()).String(),
		ProposerPriority: validator.ProposerPriority,
		PubKey:           anyPub,
		VotingPower:      validator.VotingPower,
	}, nil
}
