package types_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	_ "github.com/functionx/fx-core/v3/app"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestMsgBondedOracle_ValidateBasic(t *testing.T) {
	randomAddrPrefix := strings.ToLower(tmrand.Str(5))
	errPrefixAddress, err := bech32.ConvertAndEncode(randomAddrPrefix, tmrand.Bytes(20))
	require.NoError(t, err)
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	chainName := "bsc"
	invalidChainName := fmt.Sprintf("a%sb", tmrand.Str(5))
	externalAddr := common.BytesToAddress(tmrand.Bytes(20)).Hex()

	tests := []struct {
		name    string
		msg     *types.MsgBondedOracle
		wantErr string
	}{
		{
			"invalid chain name", &types.MsgBondedOracle{
				ChainName: invalidChainName,
			},
			"unrecognized cross chain name: invalid request",
		},
		{
			"invalid oracle address", &types.MsgBondedOracle{
				ChainName:     chainName,
				OracleAddress: errPrefixAddress,
			},
			fmt.Sprintf("invalid oracle address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			"invalid bridger address", &types.MsgBondedOracle{
				ChainName:      chainName,
				OracleAddress:  addr.String(),
				BridgerAddress: errPrefixAddress,
			},
			fmt.Sprintf("invalid bridger address: invalid Bech32 prefix; expected %s, got %s: invalid address", sdk.Bech32MainPrefix, randomAddrPrefix),
		},
		{
			"invalid external address", &types.MsgBondedOracle{
				ChainName:       chainName,
				OracleAddress:   addr.String(),
				BridgerAddress:  addr.String(),
				ExternalAddress: strings.ToUpper(externalAddr),
				DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1)},
			},
			fmt.Sprintf("invalid external address: invalid address (%s) doesn't pass regex: invalid address", strings.ToUpper(externalAddr)),
		},
		{
			"invalid delegate amount", &types.MsgBondedOracle{
				ChainName:       chainName,
				OracleAddress:   addr.String(),
				BridgerAddress:  addr.String(),
				ExternalAddress: "0x312469f0a5782Ab0f2b3aD223C6798bFd630d61D",
				DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(-1)},
			},
			"invalid delegation amount: invalid request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestEthereumMsgValidate_MsgOracleSetUpdatedClaimValidate(t *testing.T) {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	key, err := crypto.GenerateKey()
	assert.NoError(t, err)
	tests := []struct {
		name    string
		msg     *types.MsgOracleSetUpdatedClaim
		wantErr string
	}{
		{
			"err - zero event nonce", &types.MsgOracleSetUpdatedClaim{
				EventNonce:     0,
				BlockHeight:    100,
				OracleSetNonce: 10,
				Members: []types.BridgeValidator{{
					Power:           100,
					ExternalAddress: crypto.PubkeyToAddress(key.PublicKey).Hex(),
				}},
				BridgerAddress: addr.String(),
				ChainName:      "bsc",
			},
			"zero event nonce: invalid request",
		}, {
			"err - empty members", &types.MsgOracleSetUpdatedClaim{
				EventNonce:     100,
				BlockHeight:    100,
				OracleSetNonce: 10,
				Members:        []types.BridgeValidator{},
				BridgerAddress: addr.String(),
				ChainName:      "bsc",
			},
			"empty members: invalid request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestValidateModuleName(t *testing.T) {
	for _, name := range []string{
		ethtypes.ModuleName,
		bsctypes.ModuleName,
		polygontypes.ModuleName,
		trontypes.ModuleName,
		avalanchetypes.ModuleName,
	} {
		assert.NoError(t, types.ValidateModuleName(name))
	}
}
