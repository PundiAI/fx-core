package keepers_test

import (
	"reflect"
	"testing"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/app/keepers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func TestNewAppKeeper(t *testing.T) {
	interfaceRegistry := app.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(appCodec, authtx.DefaultSignModes)
	baseApp := baseapp.NewBaseApp(
		fxtypes.Name,
		log.NewNopLogger(),
		dbm.NewMemDB(),
		txConfig.TxDecoder(),
	)

	appKeeper := keepers.NewAppKeeper(
		appCodec,
		baseApp,
		codec.NewLegacyAmino(),
		app.GetMaccPerms(),
		nil,
		nil,
		fxtypes.GetDefaultNodeHome(),
		0,
		log.NewNopLogger(),
		app.EmptyAppOptions{},
	)
	assert.NotNil(t, appKeeper)
	typeOf := reflect.TypeOf(appKeeper)
	valueOf := reflect.ValueOf(appKeeper)
	checkStructField(t, valueOf, typeOf.Name())
}

//nolint:gocyclo
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
			// gov hooks not used
			if valueOfField.Type().String() == "types.GovHooks" {
				continue
			}
			assert.Falsef(t, valueOfField.IsNil(), "%s-%s-%s", valueOf.Type().PkgPath(), typeOfField.Name, name)
		}

		if valueOfField.Kind() == reflect.Pointer || valueOfField.Kind() == reflect.Interface {
			if typeOfField.Name == "QueryServer" ||
				(name == "EvidenceKeeper" && typeOfField.Name == "router") ||
				(name == "FeeGrantKeeper" && typeOfField.Name == "bankKeeper") || // deprecated in v0.50
				(name == "AuthzKeeper" && typeOfField.Name == "bankKeeper") { // deprecated in v0.50
				return
			}
			assert.Falsef(t, valueOfField.IsNil(), "%s-%s-%s", valueOf.Type().PkgPath(), typeOfField.Name, name)
		}
		checkStructField(t, valueOfField, typeOfField.Name)
	}
}
