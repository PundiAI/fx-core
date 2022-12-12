package amino_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

func TestAminoEncodeRegisterCoinProposal(t *testing.T) {
	encode := app.MakeEncodingConfig()
	proposal := erc20types.RegisterCoinProposal{
		Title:       "v2",
		Description: "foo",
		Metadata: types.Metadata{
			Description: "test",
			DenomUnits: []*types.DenomUnit{
				{
					Denom:    "test",
					Exponent: 0,
					Aliases: []string{
						"ethtest",
					},
				},
				{
					Denom:    "TEST",
					Exponent: 18,
					Aliases:  []string{},
				},
			},
			Base:    "test",
			Display: "test",
			Name:    "test name",
			Symbol:  "TEST",
		},
	}
	data, err := encode.Amino.MarshalJSON(proposal)
	require.NoError(t, err)
	require.Equal(t, `{"title":"v2","description":"foo","metadata":{"description":"test","denom_units":[{"denom":"test","aliases":["ethtest"]},{"denom":"TEST","exponent":18}],"base":"test","display":"test","name":"test name","symbol":"TEST"}}`, string(data))

	marshal, err := encode.Codec.MarshalJSON(&proposal)
	require.NoError(t, err)
	require.Equal(t, `{"title":"v2","description":"foo","metadata":{"description":"test","denom_units":[{"denom":"test","exponent":0,"aliases":["ethtest"]},{"denom":"TEST","exponent":18,"aliases":[]}],"base":"test","display":"test","name":"test name","symbol":"TEST"}}`, string(marshal))
}
