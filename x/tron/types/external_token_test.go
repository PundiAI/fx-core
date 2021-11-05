package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateExternalAddress(t *testing.T) {
	testCases := []struct {
		testName   string
		value      string
		expectPass bool
		err        error
	}{
		{
			testName:   "empty address",
			value:      "",
			expectPass: false,
			err:        fmt.Errorf("empty"),
		},
		{
			testName:   "address length not match",
			value:      "abcdddddd",
			expectPass: false,
			err:        fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", "abcdddddd", len("abcdddddd"), ExternalContractAddressLen),
		},
		{
			testName:   "address length great than tron address",
			value:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t6666",
			expectPass: false,
			err:        fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t6666", len("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t6666"), ExternalContractAddressLen),
		},
		{
			testName:   "lowercase address",
			value:      "tr7nhqjekqxgtci8q8zy4pl8otszgjlj6t",
			expectPass: false,
			err:        fmt.Errorf("invalid address: %s", "tr7nhqjekqxgtci8q8zy4pl8otszgjlj6t"),
		},
		{
			testName:   "uppercase address",
			value:      "TR7NHQJEKQXGTCI8Q8ZY4PL8OTSZGJLJ6T",
			expectPass: false,
			err:        fmt.Errorf("invalid address: %s", "TR7NHQJEKQXGTCI8Q8ZY4PL8OTSZGJLJ6T"),
		},
		{
			testName:   "normal address",
			value:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			expectPass: true,
			err:        nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err := ValidateExternalAddress(testCase.value)
			if testCase.expectPass {
				require.NoError(t, err)
				return
			}
			require.EqualValues(t, err.Error(), testCase.err.Error())
		})
	}
}
