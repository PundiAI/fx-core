// nolint
package v2_test

import (
	"math/rand"
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

func (suite *TestSuite) TestCodec() {

	testCases := []struct {
		name     string
		oldValue codec.ProtoMarshaler
		newValue codec.ProtoMarshaler
	}{
		{
			name: "OracleSet",
			oldValue: &types.Valset{
				Nonce: rand.Uint64(),
				Members: []*types.BridgeValidator{
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().String(),
					},
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().String(),
					},
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().String(),
					},
				},
				Height: rand.Uint64(),
			},
			newValue: &crosschaintypes.OracleSet{},
		},
		{
			name: "OutgoingTransferTx",
			oldValue: &types.OutgoingTransferTx{
				Id:          rand.Uint64(),
				Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				DestAddress: helpers.GenerateAddress().String(),
				Erc20Token: &types.ERC20Token{
					Contract: helpers.GenerateAddress().String(),
					Amount:   sdk.NewInt(rand.Int63() + 1),
				},
				Erc20Fee: &types.ERC20Token{
					Contract: helpers.GenerateAddress().String(),
					Amount:   sdk.NewInt(rand.Int63() + 1),
				},
			},
			newValue: &crosschaintypes.OutgoingTransferTx{},
		},
		{
			name: "OutgoingTxBatch",
			oldValue: &types.OutgoingTxBatch{
				BatchNonce:   rand.Uint64(),
				BatchTimeout: rand.Uint64(),
				Transactions: []*types.OutgoingTransferTx{
					{
						Id:          rand.Uint64(),
						Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
						DestAddress: helpers.GenerateAddress().String(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers.GenerateAddress().String(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers.GenerateAddress().String(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
					},
					{
						Id:          rand.Uint64(),
						Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
						DestAddress: helpers.GenerateAddress().String(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers.GenerateAddress().String(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers.GenerateAddress().String(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
					},
				},
				TokenContract: helpers.GenerateAddress().String(),
				Block:         rand.Uint64(),
				FeeReceive:    helpers.GenerateAddress().String(),
			},
			newValue: &crosschaintypes.OutgoingTxBatch{},
		},
	}

	for _, test := range testCases {
		suite.Run(test.name, func() {

			suite.cdc.MustUnmarshal(suite.cdc.MustMarshal(test.oldValue), test.newValue)

			valueOf1 := reflect.Indirect(reflect.ValueOf(test.oldValue))
			valueOf2 := reflect.Indirect(reflect.ValueOf(test.newValue))
			for i := 0; i < valueOf1.NumField(); i++ {
				if valueOf1.Field(i).Kind() == reflect.Slice {
					for j := 0; j < valueOf1.Field(i).Len(); j++ {
						object := reflect.Indirect(valueOf1.Field(i).Index(j))
						for n := 0; n < object.NumField(); n++ {
							if object.Field(n).Kind() == reflect.Pointer || valueOf1.Field(i).Kind() == reflect.Struct {
								subObject := reflect.Indirect(object.Field(n))
								for m := 0; m < subObject.NumField(); m++ {
									suite.T().Log(subObject.Field(m).Interface(), reflect.Indirect(valueOf2.Field(i).Index(j)).Field(n).Field(m).Interface())
									suite.Equal(subObject.Field(m).Interface(), reflect.Indirect(valueOf2.Field(i).Index(j)).Field(n).Field(m).Interface())
								}
							} else {
								suite.T().Log(object.Field(n).Interface(), reflect.Indirect(valueOf2.Field(i).Index(j)).Field(n).Interface())
								suite.Equal(object.Field(n).Interface(), reflect.Indirect(valueOf2.Field(i).Index(j)).Field(n).Interface())
							}
						}
					}
				} else if valueOf1.Field(i).Kind() == reflect.Pointer || valueOf1.Field(i).Kind() == reflect.Struct {
					object := reflect.Indirect(valueOf1.Field(i))
					for n := 0; n < object.NumField(); n++ {
						suite.T().Log(object.Field(n).Interface(), valueOf2.Field(i).Field(n).Interface())
						suite.Equal(object.Field(n).Interface(), valueOf2.Field(i).Field(n).Interface())
					}
				} else {
					suite.T().Log(valueOf1.Field(i).Interface(), valueOf2.Field(i).Interface())
					suite.Equal(valueOf1.Field(i).Interface(), valueOf2.Field(i).Interface())
				}
			}
		})
	}
}
