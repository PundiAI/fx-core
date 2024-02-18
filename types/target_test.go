package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

func TestParseTargetIBC(t *testing.T) {
	type expect struct {
		target  string
		prefix  string
		port    string
		channel string
		isIBC   bool
	}
	testCases := []struct {
		name      string
		targetStr string
		expect    expect
	}{
		{
			name:      "normal ibc data hex fx/transfer/channel-0 to targetStr ",
			targetStr: "fx/transfer/channel-0",
			expect: expect{
				prefix:  "fx",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "normal ibc data hex 0x/transfer/channel-0 to targetStr ",
			targetStr: "0x/transfer/channel-0",
			expect: expect{
				prefix:  "0x",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "normal ibc data hex upper prefix 0X/transfer/channel-0 to targetStr ",
			targetStr: "0X/transfer/channel-0",
			expect: expect{
				prefix:  "0X",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "no prefix ibc data /transfer/channel-0",
			targetStr: "/transfer/channel-0",
			expect: expect{
				target: "/transfer/channel-0",
				isIBC:  false,
			},
		},
		{
			name:      "no prefix and no port ibc data /channel-0",
			targetStr: "/channel-0",
			expect: expect{
				target: "/channel-0",
				isIBC:  false,
			},
		},
		{
			name:      "empty ibc data ''",
			targetStr: "''",
			expect: expect{
				target: "''",
				isIBC:  false,
			},
		},
		{
			name:      "two slash ibc data //",
			targetStr: "//",
			expect: expect{
				target: "//",
				isIBC:  false,
			},
		},
		{
			name:      "chain prefix",
			targetStr: "chain/gravity",
			expect: expect{
				target: "eth",
				isIBC:  false,
			},
		},
		{
			name:      "chain prefix, empty module",
			targetStr: "chain/",
			expect: expect{
				target: "",
				isIBC:  false,
			},
		},
		{
			name:      "module",
			targetStr: "gravity",
			expect: expect{
				target: "eth",
				isIBC:  false,
			},
		},
		{
			name:      "empty",
			targetStr: "",
			expect: expect{
				target: "",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with channel/prefix",
			targetStr: "ibc/0/px",
			expect: expect{
				prefix:  "px",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with channel/prefix, but empty address prefix",
			targetStr: "ibc/0/",
			expect: expect{
				target: "ibc/0/",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with channel/prefix, but empty channel sequence",
			targetStr: "ibc//px",
			expect: expect{
				target: "ibc//px",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with channel/prefix, but empty channel sequence and address prefix",
			targetStr: "ibc//",
			expect: expect{
				target: "ibc//",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel",
			targetStr: "ibc/px/transfer/channel-0",
			expect: expect{
				prefix:  "px",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty address prefix",
			targetStr: "ibc//transfer/channel-0",
			expect: expect{
				target: "/transfer/channel-0",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port",
			targetStr: "ibc/px//channel-0",
			expect: expect{
				target: "px//channel-0",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty channel",
			targetStr: "ibc/px/transfer/",
			expect: expect{
				target: "px/transfer/",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port and address prefix",
			targetStr: "ibc///channel-0",
			expect: expect{
				target: "//channel-0",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port and channel",
			targetStr: "ibc/px//",
			expect: expect{
				target: "px//",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty prefix and channel",
			targetStr: "ibc//transfer/",
			expect: expect{
				target: "/transfer/",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty all",
			targetStr: "ibc///",
			expect: expect{
				target: "//",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with '/'",
			targetStr: "ibc/",
			expect: expect{
				target: "ibc/",
				isIBC:  false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target := fxtypes.ParseFxTarget(tc.targetStr)
			require.EqualValues(t, tc.expect.isIBC, target.IsIBC(), tc.name)
			if tc.expect.isIBC {
				require.EqualValues(t, tc.expect.prefix, target.Prefix, tc.name)
				require.EqualValues(t, tc.expect.port, target.SourcePort, tc.name)
				require.EqualValues(t, tc.expect.channel, target.SourceChannel, tc.name)
			} else {
				require.EqualValues(t, tc.expect.target, target.GetTarget(), tc.name)
			}
		})
	}
}

func TestGetIbcDenomTrace(t *testing.T) {
	type args struct {
		denom      string
		channelIBC string
	}
	tests := []struct {
		name    string
		args    args
		want    ibctransfertypes.DenomTrace
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: hex.EncodeToString([]byte("transfer/channel-0")),
			},
			want: ibctransfertypes.DenomTrace{
				Path:      "transfer/channel-0",
				BaseDenom: fxtypes.DefaultDenom,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return true
			},
		},
		{
			name: "ok empty",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: "",
			},
			want: ibctransfertypes.DenomTrace{
				Path:      "",
				BaseDenom: fxtypes.DefaultDenom,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return true
			},
		},
		{
			name: "error decode hex",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: "transfer/channel-0",
			},
			want: ibctransfertypes.DenomTrace{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "decode hex channel-ibc err")
				return true
			},
		},
		{
			name: "error split channel-ibc",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: hex.EncodeToString([]byte("channel-0")),
			},
			want: ibctransfertypes.DenomTrace{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, err.Error(), "invalid params channel-ibc")
				return true
			},
		},
		{
			name: "error source port",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: hex.EncodeToString([]byte("tran/channel-0")),
			},
			want: ibctransfertypes.DenomTrace{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, err.Error(), "invalid source port")
				return true
			},
		},
		{
			name: "error source channel",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: hex.EncodeToString([]byte("transfer/chan-0")),
			},
			want: ibctransfertypes.DenomTrace{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, err.Error(), "invalid source channel")
				return true
			},
		},
		{
			name: "error source channel-index",
			args: args{
				denom:      fxtypes.DefaultDenom,
				channelIBC: hex.EncodeToString([]byte("transfer/channel-x")),
			},
			want: ibctransfertypes.DenomTrace{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, err.Error(), "invalid source channel")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fxtypes.GetIbcDenomTrace(tt.args.denom, tt.args.channelIBC)
			if !tt.wantErr(t, err, fmt.Sprintf("GetIbcDenomTrace(%v, %v)", tt.args.denom, tt.args.channelIBC)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetIbcDenomTrace(%v, %v)", tt.args.denom, tt.args.channelIBC)
		})
	}
}
