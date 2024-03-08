package app_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/app"
)

func TestMakeEncodingConfig_RegisterInterfaces(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()

	interfaceRegistry := reflect.ValueOf(encodingConfig.Codec).Elem().Field(0).Elem().Elem()

	interfaceNames := interfaceRegistry.Field(0).MapRange()
	var count1 int
	for interfaceNames.Next() {
		count1++
	}
	assert.Equal(t, 33, count1)

	interfaceImpls := interfaceRegistry.Field(1).MapRange()
	var count2 int
	for interfaceImpls.Next() {
		count2++
	}
	assert.Equal(t, 33, count2)

	typeURLMap := interfaceRegistry.Field(2).MapRange()
	var count3 int
	for typeURLMap.Next() {
		count3++
	}
	assert.Equal(t, 267, count3)

	govContent := encodingConfig.InterfaceRegistry.ListImplementations("cosmos.gov.v1beta1.Content")
	assert.Equal(t, 14, len(govContent))

	msgImplementations := encodingConfig.InterfaceRegistry.ListImplementations(sdk.MsgInterfaceProtoName)
	assert.Equal(t, 103, len(msgImplementations))

	type govProposalMsg interface {
		GetAuthority() string
	}
	var govMsg []string
	for _, implementation := range msgImplementations {
		resolvedMsg, err := encodingConfig.InterfaceRegistry.Resolve(implementation)
		assert.NoError(t, err)

		if _, ok := resolvedMsg.(govProposalMsg); ok {
			govMsg = append(govMsg, implementation)
		}
	}
	assert.Equal(t, 15, len(govMsg))
}
