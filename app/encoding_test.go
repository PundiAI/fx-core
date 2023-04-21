package app_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v4/app"
)

func TestMakeEncodingConfig_RegisterInterfaces(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()
	interfaceRegistry := reflect.ValueOf(encodingConfig.Codec).Elem().Field(0).Elem().Elem()

	interfaceNames := interfaceRegistry.Field(0).MapRange()
	var count1 int
	for interfaceNames.Next() {
		count1++
		t.Log(interfaceNames.Key())
	}
	assert.Equal(t, 32, count1)

	interfaceImpls := interfaceRegistry.Field(1).MapRange()
	var count2 int
	for interfaceImpls.Next() {
		count2++
		t.Log(interfaceImpls.Value())
	}
	assert.Equal(t, 32, count2)

	typeURLMap := interfaceRegistry.Field(2).MapRange()
	var count3 int
	for typeURLMap.Next() {
		count3++
		t.Log(typeURLMap.Key())
	}
	assert.Equal(t, 258, count3)
}
