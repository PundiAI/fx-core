package amino

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v2/app"
	erc20types "github.com/functionx/fx-core/v2/x/erc20/types"
)

func TestAminoEncodeUpdateDenomAliasProposal(t *testing.T) {
	encode := app.MakeEncodingConfig()
	proposal := erc20types.RegisterCoinProposal{
		Title:       "v2",
		Description: "foo",
		Metadata:    types.Metadata{},
	}
	data, err := encode.Amino.MarshalJSON(proposal)
	require.NoError(t, err)
	require.Equal(t, `{"title":"v2","description":"foo","metadata":{}}`, string(data))

	marshal, err := encode.Marshaler.MarshalJSON(&proposal)
	require.NoError(t, err)
	require.Equal(t, `{"title":"v2","description":"foo","metadata":{"description":"","denom_units":[],"base":"","display":"","name":"","symbol":""}}`, string(marshal))
}
