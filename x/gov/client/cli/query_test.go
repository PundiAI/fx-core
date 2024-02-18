package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/x/gov/client/cli"
)

func TestCmdParams(t *testing.T) {
	testCases := []struct {
		name         string
		args         []string
		expCmdOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", "output")},
			"--output=json",
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=text", "output")},
			"--output=text",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cmd := cli.GetCmdQueryParams()
			cmd.SetArgs(tc.args)
			assert.Contains(t, fmt.Sprint(cmd), strings.TrimSpace(tc.expCmdOutput))
		})
	}
}

func TestCmdEGFParams(t *testing.T) {
	testCases := []struct {
		name         string
		args         []string
		expCmdOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", "output")},
			"--output=json",
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=text", "output")},
			"--output=text",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cmd := cli.GetCmdQueryEGFParams()
			cmd.SetArgs(tc.args)
			assert.Contains(t, fmt.Sprint(cmd), strings.TrimSpace(tc.expCmdOutput))
		})
	}
}
