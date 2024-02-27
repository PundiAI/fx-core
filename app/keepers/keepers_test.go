package keepers_test

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/app/keepers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	arbitrumtypes "github.com/functionx/fx-core/v7/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v7/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	layer2types "github.com/functionx/fx-core/v7/x/layer2/types"
	optimismtypes "github.com/functionx/fx-core/v7/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v7/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

func TestNewAppKeeper(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino

	bApp := baseapp.NewBaseApp(
		fxtypes.Name,
		tmlog.NewNopLogger(),
		dbm.NewMemDB(),
		encodingConfig.TxConfig.TxDecoder(),
	)
	maccPerms := map[string][]string{
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		bsctypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
		polygontypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		avalanchetypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		ethtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
		trontypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		arbitrumtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		optimismtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		layer2types.ModuleName:         {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
		erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
	}

	keeper := keepers.NewAppKeeper(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		nil,
		nil,
		fxtypes.GetDefaultNodeHome(),
		0,
		app.EmptyAppOptions{},
	)
	assert.NotNil(t, keeper)
	typeOf := reflect.TypeOf(keeper)
	valueOf := reflect.ValueOf(keeper)
	checkStructField(t, valueOf, typeOf.Name())
}

func checkStructField(t *testing.T, valueOf reflect.Value, name string) {
	valueOf = reflect.Indirect(valueOf)
	if valueOf.Kind() != reflect.Struct ||
		valueOf.Type().String() == "baseapp.MsgServiceRouter" {
		return
	}

	numberField := valueOf.NumField()
	for i := 0; i < numberField; i++ {
		valueOfField := valueOf.Field(i)
		typeOfField := valueOf.Type().Field(i)
		switch typeOfField.Name {
		case "storeKey":
			assert.False(t, valueOfField.IsNil(), typeOfField.Name)
		case "hooks":
			// evm hooks deprecated
			if valueOfField.Type().String() == "types.EvmHooks" {
				continue
			}
			assert.False(t, valueOfField.IsNil(), typeOfField.Name)
		}

		switch valueOfField.Kind() {
		case reflect.Pointer, reflect.Interface:
			if typeOfField.Name == "QueryServer" ||
				(name == "EvidenceKeeper" && typeOfField.Name == "router") {
				return
			}
			assert.False(t, valueOfField.IsNil(), typeOfField.Name)
		}
		checkStructField(t, valueOfField, typeOfField.Name)
	}
}
