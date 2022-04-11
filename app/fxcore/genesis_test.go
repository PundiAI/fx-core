package fxcore

import (
	"encoding/json"
	fxtypes "github.com/functionx/fx-core/types"
	"testing"

	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

func TestNewDefaultGenesisByDenom(t *testing.T) {
	encodingConfig := MakeEncodingConfig()
	genAppState := NewDefAppGenesisByDenom(fxtypes.MintDenom, encodingConfig.Marshaler)

	state := gravitytypes.DefaultGenesisState()
	state.Erc20ToDenoms = []*gravitytypes.ERC20ToDenom{
		{
			Denom: fxtypes.MintDenom,                            // token symbol
			Erc20: "0x0AD5CE837A789423CC6158053CAd5eB75A6183AC", // token contract address
		},
	}
	data, err := json.Marshal(map[string]interface{}{gravitytypes.ModuleName: state})
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(data, &genAppState); err != nil && len(data) > 0 {
		t.Fatal(err)
	}
	genAppStateStr, err := json.Marshal(genAppState)
	if err != nil {
		t.Fatal(err)
	}
	_ = genAppStateStr
	//t.Log(string(genAppStateStr))
}
