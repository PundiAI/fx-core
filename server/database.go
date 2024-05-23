package server

import (
	"errors"
	"path/filepath"

	cmtdbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

type Database struct {
	genesisFile string
	blockStore  *store.BlockStore
	stateStore  sm.Store
	appDB       dbm.DB
	appStore    *rootmulti.Store
	storeKeys   map[string]*storetypes.KVStoreKey
	codec       codec.Codec
}

func NewDatabase(cfg *tmcfg.Config, cdc codec.Codec, modules ...string) (*Database, error) {
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
	storeKeys := sdk.NewKVStoreKeys(append(modules, upgradetypes.StoreKey)...)
	appStore := rootmulti.NewStore(appDB, log.NewNopLogger())

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
		codec:       cdc,
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

func (d *Database) GetNodeInfo() (*tmservice.VersionInfo, error) {
	state, err := d.GetState()
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
	if err := d.codec.Unmarshal(bz, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (d *Database) GetConsensusValidators() ([]*tmservice.Validator, error) {
	state, err := d.GetState()
	if err != nil {
		return nil, err
	}
	if len(state.Validators.Validators) == 0 {
		return nil, nil
	}
	validators := make([]*tmservice.Validator, len(state.Validators.Validators))
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

func toValidator(validator *tmtypes.Validator) (*tmservice.Validator, error) {
	pk, err := cryptocodec.FromTmPubKeyInterface(validator.PubKey)
	if err != nil {
		return nil, err
	}
	anyPub, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return nil, err
	}
	return &tmservice.Validator{
		Address:          sdk.ConsAddress(validator.Address.Bytes()).String(),
		ProposerPriority: validator.ProposerPriority,
		PubKey:           anyPub,
		VotingPower:      validator.VotingPower,
	}, nil
}
