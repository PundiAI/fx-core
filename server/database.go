package server

import (
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type Database struct{}

func (d *Database) GetChainId() (string, error) {
	// TODO implement me
	panic("implement me")
}

func (d *Database) GetBlockHeight() (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (d *Database) GetSyncing() (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (d *Database) GetNodeInfo() (*tmservice.VersionInfo, error) {
	// TODO implement me
	panic("implement me")
}

func (d *Database) CurrentPlan() (*upgradetypes.Plan, error) {
	// TODO implement me
	panic("implement me")
}
