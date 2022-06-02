package types_test

import (
	"github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDenomTrace(t *testing.T) {
	testCases := []struct {
		name     string
		denom    string
		expTrace types.DenomTrace
	}{
		{"empty denom", "", types.DenomTrace{}},
		{"base denom", "uatom", types.DenomTrace{BaseDenom: "uatom"}},
		{"trace info", "transfer/channelToA/uatom", types.DenomTrace{BaseDenom: "uatom", Path: "transfer/channelToA"}},
		{"incomplete path", "transfer/uatom", types.DenomTrace{BaseDenom: "uatom", Path: "transfer"}},
		{"invalid path (1)", "transfer//uatom", types.DenomTrace{BaseDenom: "uatom", Path: "transfer/"}},
		{"invalid path (2)", "transfer/channelToA/uatom/", types.DenomTrace{BaseDenom: "", Path: "transfer/channelToA/uatom"}},
	}

	for _, tc := range testCases {
		trace := types.ParseDenomTrace(tc.denom)
		require.Equal(t, tc.expTrace, trace, tc.name)
	}
}

func TestDenomTrace_IBCDenom(t *testing.T) {
	testCases := []struct {
		name     string
		trace    types.DenomTrace
		expDenom string
	}{
		{"base denom", types.DenomTrace{BaseDenom: "uatom"}, "uatom"},
		{"trace info", types.DenomTrace{BaseDenom: "uatom", Path: "transfer/channelToA"}, "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2"},
	}

	for _, tc := range testCases {
		denom := tc.trace.IBCDenom()
		require.Equal(t, tc.expDenom, denom, tc.name)
	}
}

func TestDenomTrace_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		trace    types.DenomTrace
		expError bool
	}{
		{"base denom only", types.DenomTrace{BaseDenom: "uatom"}, false},
		{"empty DenomTrace", types.DenomTrace{}, true},
		{"valid single trace info", types.DenomTrace{BaseDenom: "uatom", Path: "transfer/channelToA"}, false},
		{"valid multiple trace info", types.DenomTrace{BaseDenom: "uatom", Path: "transfer/channelToA/transfer/channelToB"}, false},
		{"single trace identifier", types.DenomTrace{BaseDenom: "uatom", Path: "transfer"}, true},
		{"invalid port ID", types.DenomTrace{BaseDenom: "uatom", Path: "(transfer)/channelToA"}, true},
		{"invalid channel ID", types.DenomTrace{BaseDenom: "uatom", Path: "transfer/(channelToA)"}, true},
		{"empty base denom with trace", types.DenomTrace{BaseDenom: "", Path: "transfer/channelToA"}, true},
	}

	for _, tc := range testCases {
		err := tc.trace.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
			continue
		}
		require.NoError(t, err, tc.name)
	}
}

func TestTraces_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		traces   types.Traces
		expError bool
	}{
		{"empty Traces", types.Traces{}, false},
		{"valid multiple trace info", types.Traces{{BaseDenom: "uatom", Path: "transfer/channelToA/transfer/channelToB"}}, false},
		{
			"valid multiple trace info",
			types.Traces{
				{BaseDenom: "uatom", Path: "transfer/channelToA/transfer/channelToB"},
				{BaseDenom: "uatom", Path: "transfer/channelToA/transfer/channelToB"},
			},
			true,
		},
		{"empty base denom with trace", types.Traces{{BaseDenom: "", Path: "transfer/channelToA"}}, true},
	}

	for _, tc := range testCases {
		err := tc.traces.Validate()
		if tc.expError {
			require.Error(t, err, tc.name)
			continue
		}
		require.NoError(t, err, tc.name)
	}
}

func TestValidatePrefixedDenom(t *testing.T) {
	testCases := []struct {
		name     string
		denom    string
		expError bool
	}{
		{"prefixed denom", "transfer/channelToA/uatom", false},
		{"base denom", "uatom", false},
		{"empty denom", "", true},
		{"empty prefix", "/uatom", true},
		{"empty identifiers", "//uatom", true},
		{"single trace identifier", "transfer/", true},
		{"invalid port ID", "(transfer)/channelToA/uatom", true},
		{"invalid channel ID", "transfer/(channelToA)/uatom", true},
	}

	for _, tc := range testCases {
		err := types.ValidatePrefixedDenom(tc.denom)
		if tc.expError {
			require.Error(t, err, tc.name)
			continue
		}
		require.NoError(t, err, tc.name)
	}
}

func TestValidateIBCDenom(t *testing.T) {
	testCases := []struct {
		name     string
		denom    string
		expError bool
	}{
		{"denom with trace hash", "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", false},
		{"base denom", "uatom", false},
		{"base denom with single '/'s", "gamm/pool/1", false},
		{"base denom with double '/'s", "gamm//pool//1", false},
		{"non-ibc prefix with hash", "notibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", false},
		{"empty denom", "", true},
		{"denom 'ibc'", "ibc", true},
		{"denom 'ibc/'", "ibc/", true},
		{"invald hash", "ibc/!@#$!@#", true},
	}

	for _, tc := range testCases {
		err := types.ValidateIBCDenom(tc.denom)
		if tc.expError {
			require.Error(t, err, tc.name)
			continue
		}
		require.NoError(t, err, tc.name)
	}
}
