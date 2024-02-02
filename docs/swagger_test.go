package docs_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/api"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestSwaggerConfig(t *testing.T) {
	data, err := os.ReadFile("config.json")
	assert.NoError(t, err)
	var c config
	assert.NoError(t, json.Unmarshal(data, &c))
	assert.Equal(t, "2.0", c.Swagger)
	assert.Equal(t, "0.4.0", c.Info.Version)
	assert.Equal(t, 23, len(c.Apis))
	app := helpers.Setup(true, false)
	clientCtx := client.Context{
		InterfaceRegistry: app.InterfaceRegistry(),
	}
	apiSrv := api.New(clientCtx, app.Logger())
	app.RegisterAPIRoutes(apiSrv, srvconfig.APIConfig{Swagger: true})
	assert.NotNil(t, apiSrv.Router.Path("/swagger/"))
	handler := reflect.Indirect(reflect.ValueOf(apiSrv.GRPCGatewayRouter)).Field(0).MapRange()
	route := make(map[string]int)
	for handler.Next() {
		for i := 0; i < handler.Value().Len(); i++ {
			field := handler.Value().Index(i).Field(0)
			pat := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface().(runtime.Pattern)
			split := strings.Split(pat.String(), "/")
			assert.True(t, len(split) > 3)
			if len(split) > 4 && split[3] != "v1" && split[3] != "v1beta1" && (split[4] == "v1" || split[4] == "v1beta1") {
				split[3] = fmt.Sprintf("%s/%s", split[3], split[4])
			}
			key := fmt.Sprintf("%s/%s/%s", split[1], split[2], split[3])
			if key == "ibc/apps/transfer/v1" {
				key = "ibc/applications/transfer/v1"
			}
			route[key] = route[key] + 1
		}
		if handler.Key().String() == "POST" {
			assert.Equal(t, 2, handler.Value().Len())
		}
		if handler.Key().String() == "GET" {
			assert.Equal(t, 205, handler.Value().Len())
		}
	}
	assert.Equal(t, 32, len(route))
	ignoreLen := len(route) - len(c.Apis)
	for _, v := range c.Apis {
		for key := range route {
			if strings.HasPrefix(v.Url, "./tmp-swagger-gen/"+key) {
				delete(route, key)
			}
		}
	}
	for k := range route {
		t.Log("ignore", k)
		// ignore routes:
		// 1. other/v1/gas_price
		// 2. fx/gravity/v1
		// 3. fx/other/gas_price
		// 4. fx/base/v1
		// 5. fx/ibc/applications
		// 6. ibc/core/channel/v1
		// 7. ibc/core/client/v1
		// 8. ibc/core/connection/v1
		// 9. cosmos/gov/v1beta1
	}
	assert.Equal(t, ignoreLen, len(route))
	assert.Equal(t, 9, len(route))
}

type config struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
	} `json:"info"`
	Apis []struct {
		Url string `json:"url"`
	} `json:"apis"`
}
