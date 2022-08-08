package amino

import (
	"testing"
	"time"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v2/app"
)

func TestAminoEncodeSoftwareUpgradeProposal(t *testing.T) {
	encode := app.MakeEncodingConfig()
	proposal := upgradetypes.SoftwareUpgradeProposal{
		Title:       "v2",
		Description: "foo",
		Plan: upgradetypes.Plan{
			Name:   "foo",
			Time:   time.Time{},
			Height: 123,
			Info:   "foo",
		},
	}
	data, err := encode.Amino.MarshalJSON(proposal)
	require.NoError(t, err)
	require.Equal(t, `{"type":"cosmos-sdk/SoftwareUpgradeProposal","value":{"title":"v2","description":"foo","plan":{"name":"foo","time":"0001-01-01T00:00:00Z","height":"123","info":"foo"}}}`, string(data))

	marshal, err := encode.Marshaler.MarshalJSON(&proposal)
	require.NoError(t, err)
	require.Equal(t, `{"title":"v2","description":"foo","plan":{"name":"foo","time":"0001-01-01T00:00:00Z","height":"123","info":"foo","upgraded_client_state":null}}`, string(marshal))
}
