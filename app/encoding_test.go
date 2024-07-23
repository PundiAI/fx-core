package app_test

import (
	"reflect"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestMakeEncodingConfig_RegisterInterfaces(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()

	// github.com/cosmos/cosmos/codec/types.interfaceRegistry
	interfaceRegistry := reflect.ValueOf(encodingConfig.Codec).Elem().Field(0).Elem().Elem()

	result := struct {
		InterfaceNames []string
		TypeURLMap     []string
		GovContent     []string
		Msgs           []string
		ProposalMsgs   []string
	}{}

	interfaceNames := interfaceRegistry.FieldByName("interfaceNames").MapRange()
	for interfaceNames.Next() {
		result.InterfaceNames = append(result.InterfaceNames, interfaceNames.Key().String())
	}
	sort.Strings(result.InterfaceNames)

	interfaceImpls := interfaceRegistry.FieldByName("interfaceImpls").MapRange()
	var count1 int
	for interfaceImpls.Next() {
		count1++
	}
	assert.Equal(t, 32, count1)

	implInterfaces := interfaceRegistry.FieldByName("implInterfaces").MapRange()
	var count2 int
	for implInterfaces.Next() {
		count2++
	}
	assert.Equal(t, 296, count2)

	typeURLMap := interfaceRegistry.FieldByName("typeURLMap").MapRange()
	for typeURLMap.Next() {
		result.TypeURLMap = append(result.TypeURLMap, typeURLMap.Key().String())
	}
	sort.Strings(result.TypeURLMap)

	result.GovContent = encodingConfig.InterfaceRegistry.ListImplementations("cosmos.gov.v1beta1.Content")
	sort.Strings(result.GovContent)

	result.Msgs = encodingConfig.InterfaceRegistry.ListImplementations(sdk.MsgInterfaceProtoName)
	sort.Strings(result.Msgs)

	type govProposalMsg interface {
		GetAuthority() string
	}
	for _, implementation := range result.Msgs {
		resolvedMsg, err := encodingConfig.InterfaceRegistry.Resolve(implementation)
		assert.NoError(t, err)

		if _, ok := resolvedMsg.(govProposalMsg); ok {
			result.ProposalMsgs = append(result.ProposalMsgs, implementation)
		}
	}
	sort.Strings(result.ProposalMsgs)

	helpers.AssertJsonFile(t, "./interface_registry.json", result)
}
