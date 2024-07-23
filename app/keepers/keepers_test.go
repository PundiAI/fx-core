package keepers_test

import (
	"reflect"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	tmlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/app/keepers"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func TestNewAppKeeper(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino

	baseApp := baseapp.NewBaseApp(
		fxtypes.Name,
		tmlog.NewNopLogger(),
		dbm.NewMemDB(),
		encodingConfig.TxConfig.TxDecoder(),
	)

	appKeeper := keepers.NewAppKeeper(
		appCodec,
		baseApp,
		legacyAmino,
		app.GetMaccPerms(),
		nil,
		nil,
		fxtypes.GetDefaultNodeHome(),
		0,
		app.EmptyAppOptions{},
	)
	assert.NotNil(t, appKeeper)
	typeOf := reflect.TypeOf(appKeeper)
	valueOf := reflect.ValueOf(appKeeper)
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
			// gov hooks not used
			if valueOfField.Type().String() == "types.GovHooks" {
				continue
			}
			assert.Falsef(t, valueOfField.IsNil(), "%s-%s-%s", valueOf.Type().PkgPath(), typeOfField.Name, name)
		}

		switch valueOfField.Kind() {
		case reflect.Pointer, reflect.Interface:
			if typeOfField.Name == "QueryServer" ||
				(name == "EvidenceKeeper" && typeOfField.Name == "router") {
				return
			}
			assert.Falsef(t, valueOfField.IsNil(), "%s-%s-%s", valueOf.Type().PkgPath(), typeOfField.Name, name)
		}
		checkStructField(t, valueOfField, typeOfField.Name)
	}
}
