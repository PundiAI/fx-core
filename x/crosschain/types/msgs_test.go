package types_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	"github.com/functionx/fx-core/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
)

func init() {
	types.InitMsgValidatorBasicRouter()
	types.RegisterValidatorBasic(bsctypes.ModuleName, types.EthereumMsgValidate{})
	types.RegisterValidatorBasic(polygontypes.ModuleName, types.EthereumMsgValidate{})
}

func TestMsgBondedOracle_ValidateBasic(t *testing.T) {
	addrTooLong := sdk.AccAddress("Accidentally used 268 bytes pubkey test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content test content")
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	chainName := "bsc"
	invalidChainName := "tron"
	tests := []struct {
		name    string
		msg     *types.MsgBondedOracle
		wantErr string
	}{
		{"invalid chain name", &types.MsgBondedOracle{
			OracleAddress:   addrTooLong.String(),
			BridgerAddress:  addr.String(),
			ExternalAddress: "0x312469f0a5782Ab0f2b3aD223C6798bFd630d61D",
			DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1)},
			ChainName:       invalidChainName,
		},
			fmt.Sprintf("Unrecognized cross chain type: %v: unknown request", invalidChainName),
		},
		{"invalid oracle address", &types.MsgBondedOracle{
			OracleAddress:   addrTooLong.String(),
			BridgerAddress:  addr.String(),
			ExternalAddress: "0x312469f0a5782Ab0f2b3aD223C6798bFd630d61D",
			DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1)},
			ChainName:       chainName,
		},
			"oracle address: invalid",
		},
		{"invalid bridger address", &types.MsgBondedOracle{
			OracleAddress:   addr.String(),
			BridgerAddress:  addrTooLong.String(),
			ExternalAddress: "0x312469f0a5782Ab0f2b3aD223C6798bFd630d61D",
			DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1)},
			ChainName:       chainName,
		},
			"bridger address: invalid",
		},
		{"invalid external address", &types.MsgBondedOracle{
			OracleAddress:   addr.String(),
			BridgerAddress:  addr.String(),
			ExternalAddress: "0x312469f0a5782Ab0f2b3aD223C6798bF",
			DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1)},
			ChainName:       chainName,
		},
			"external address: invalid",
		},
		{"invalid delegate amount", &types.MsgBondedOracle{
			OracleAddress:   addr.String(),
			BridgerAddress:  addr.String(),
			ExternalAddress: "0x312469f0a5782Ab0f2b3aD223C6798bFd630d61D",
			DelegateAmount:  sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(-1)},
			ChainName:       chainName,
		},
			"delegate amount: invalid",
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
		{"invalid event nonce", &types.MsgOracleSetUpdatedClaim{
			EventNonce:     0,
			BlockHeight:    100,
			OracleSetNonce: 10,
			Members: []types.BridgeValidator{{Power: 100,
				ExternalAddress: crypto.PubkeyToAddress(key.PublicKey).Hex()}},
			BridgerAddress: addr.String(),
			ChainName:      "bsc",
		},
			"event nonce: unknown",
		}, {"empty members", &types.MsgOracleSetUpdatedClaim{
			EventNonce:     100,
			BlockHeight:    100,
			OracleSetNonce: 10,
			Members:        []types.BridgeValidator{},
			BridgerAddress: addr.String(),
			ChainName:      "bsc",
		},
			"members: empty",
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
