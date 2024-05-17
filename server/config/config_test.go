package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBypassMinFee_Validate(t *testing.T) {
	type fields struct {
		MsgTypes       []string
		MsgMaxGasUsage uint64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "bypass-min-fee.msg-types empty",
			fields: fields{
				MsgTypes:       []string{},
				MsgMaxGasUsage: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "bypass-min-fee.msg-types invalid 1",
			fields: fields{
				MsgTypes:       []string{"invalid"},
				MsgMaxGasUsage: 0,
			},
			wantErr: assert.Error,
		},
		{
			name: "bypass-min-fee.msg-types invalid 2",
			fields: fields{
				MsgTypes:       []string{" invalid"},
				MsgMaxGasUsage: 0,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := BypassMinFee{
				MsgTypes:       tt.fields.MsgTypes,
				MsgMaxGasUsage: tt.fields.MsgMaxGasUsage,
			}
			tt.wantErr(t, f.Validate(), "Validate()")
		})
	}
}
