package server

import (
	"errors"
	"path/filepath"

	cmtdbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
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
	dbm "github.com/tendermint/tm-db"
)

var genesisDocNotFoundError = errors.New("genesis doc not found")

type Blockchain interface {
	GetChainId() (string, error)
	GetBlockHeight() (int64, error)
	GetSyncing() (bool, error)
	GetNodeInfo() (*tmservice.VersionInfo, error)
	CurrentPlan() (*upgradetypes.Plan, error)
	GetValidators() ([]stakingtypes.Validator, error)
}

type Database struct {
	blockStore *store.BlockStore
	stateDB    cmtdbm.DB
	stateStore sm.Store
	appDB      dbm.DB
	appStore   *rootmulti.Store
	storeKeys  map[string]*storetypes.KVStoreKey
	codec      codec.Codec
}

func NewDatabase(rootDir string, dbType string, cdc codec.Codec) (Blockchain, error) {
	dataDir := filepath.Join(rootDir, "data")

	blockStoreDB, err := cmtdbm.NewDB(BlockDBName, cmtdbm.BackendType(dbType), dataDir)
	if err != nil {
		return nil, err
	}
	blockStore := store.NewBlockStore(blockStoreDB)

	stateDB, err := cmtdbm.NewDB(StateDBName, cmtdbm.BackendType(dbType), dataDir)
	if err != nil {
		return nil, err
	}
	stateStore := sm.NewStore(stateDB, sm.StoreOptions{DiscardABCIResponses: false})

	appDB, err := dbm.NewDB(AppDBName, dbm.BackendType(dbType), dataDir)
	if err != nil {
		return nil, err
	}
	storeKeys := sdk.NewKVStoreKeys(upgradetypes.StoreKey)
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
		stateStore: stateStore,
		appStore:   appStore,
		storeKeys:  storeKeys,
		codec:      cdc,
	}, err
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
	_ = d.blockStore.Close()
	_ = d.stateDB.Close()
	_ = d.appDB.Close()
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

func (d *Database) GetValidators() ([]stakingtypes.Validator, error) {
	state, err := d.GetState()
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

func (d *Database) GetState() (sm.State, error) {
	genDoc, err := d.genesisDoc()
	if err != nil {
		return sm.State{}, err
	}
	state, err := d.stateStore.LoadFromDBOrGenesisDoc(genDoc)
	if err != nil {
		return sm.State{}, err
	}
	return state, nil
}

func (d *Database) SetGenesis(genesis []byte) error {
	_, err := d.genesisDoc()
	if !errors.Is(err, genesisDocNotFoundError) {
		return nil
	}
	return d.stateDB.Set([]byte("genesisDoc"), genesis)
}

func (d *Database) genesisDoc() (*tmtypes.GenesisDoc, error) {
	genesisDocKey := []byte("genesisDoc")
	b, err := d.stateDB.Get(genesisDocKey)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, genesisDocNotFoundError
	}
	var genDoc *tmtypes.GenesisDoc
	if err = tmjson.Unmarshal(b, &genDoc); err != nil {
		return nil, err
	}
	return genDoc, nil
}
