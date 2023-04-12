package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuerier_ValAddressToAccAddress(t *testing.T) {
	type args struct {
		in0 context.Context
		req *ConvertAddressRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *ConvertAddressResponse
		wantErr bool
	}{
		{
			name: "use default bech32 prefix",
			args: args{
				in0: nil,
				req: &ConvertAddressRequest{
					Address: "fxvaloper15jrnlw3mrdrptjt93faa8x038na43yhmq7zz0a",
					Prefix:  "",
				},
			},
			want: &ConvertAddressResponse{
				Address: "fx15jrnlw3mrdrptjt93faa8x038na43yhmg0la3a",
			},
			wantErr: false,
		},
		{
			name: "use fx bech32 prefix",
			args: args{
				in0: nil,
				req: &ConvertAddressRequest{
					Address: "fxvaloper15jrnlw3mrdrptjt93faa8x038na43yhmq7zz0a",
					Prefix:  "fx",
				},
			},
			want: &ConvertAddressResponse{
				Address: "fx15jrnlw3mrdrptjt93faa8x038na43yhmg0la3a",
			},
			wantErr: false,
		},
		{
			name: "use fxvaloper bech32 prefix",
			args: args{
				in0: nil,
				req: &ConvertAddressRequest{
					Address: "fx15jrnlw3mrdrptjt93faa8x038na43yhmg0la3a",
					Prefix:  "fxvaloper",
				},
			},
			want: &ConvertAddressResponse{
				Address: "fxvaloper15jrnlw3mrdrptjt93faa8x038na43yhmq7zz0a",
			},
			wantErr: false,
		},
		{
			name: "use cosmos bech32 prefix",
			args: args{
				in0: nil,
				req: &ConvertAddressRequest{
					Address: "fx15jrnlw3mrdrptjt93faa8x038na43yhmg0la3a",
					Prefix:  "cosmos",
				},
			},
			want: &ConvertAddressResponse{
				Address: "cosmos15jrnlw3mrdrptjt93faa8x038na43yhmjlwguv",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := Querier{}
			got, err := qu.ConvertAddress(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValAddressToAccAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValAddressToAccAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertBech32Prefix(t *testing.T) {
	cases := []struct {
		name      string
		address   string
		prefix    string
		converted string
		err       error
	}{
		{
			name:      "Convert valid bech 32 address",
			address:   "akash1a6zlyvpnksx8wr6wz8wemur2xe8zyh0ytz6d88",
			converted: "cosmos1a6zlyvpnksx8wr6wz8wemur2xe8zyh0yxeh27a",
			prefix:    "cosmos",
		},
		{
			name:    "Convert invalid address",
			address: "invalidaddress",
			prefix:  "cosmos",
			err:     errors.New("cannot decode invalidaddress address: decoding bech32 failed: invalid separator index -1"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			convertedAddress, err := ConvertBech32Prefix(tt.address, tt.prefix)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.converted, convertedAddress)
		})
	}
}
