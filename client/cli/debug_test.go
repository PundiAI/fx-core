package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/client/cli"
)

func TestToStringCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		flags   map[string]string
		output  string
		wantErr bool
	}{
		{
			name:    "decode base64",
			args:    []string{"base64", "OQ=="},
			output:  "9\n",
			wantErr: false,
		},
		{
			name:    "decode hex",
			args:    []string{"hex", "666f6f"},
			output:  "foo\n",
			wantErr: false,
		},
		{
			name:    "decode hex",
			args:    []string{"hex", "0x666f6f"},
			output:  "foo\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cli.ToStringCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetArgs(tt.args)
			for k, v := range tt.flags {
				assert.NoError(t, cmd.Flags().Set(k, v))
			}
			if err := cmd.Execute(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.output, buf.String())
			}
		})
	}
}
