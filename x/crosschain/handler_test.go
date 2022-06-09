package crosschain_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/functionx/fx-core/app/helpers"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/crosschain"
	"github.com/functionx/fx-core/x/crosschain/types"
)

var (
	minStakeAmount = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))
)

const (
	chainName      = "xxx"
	chainGravityId = "local-test-xxx"
)

// 1. Test MsgCreateOracleBridger
func TestHandlerMsgSetOrchestratorAddress(t *testing.T) {
	// get test env
	_, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 4)
	// 1. sender not in chain oracle
	notOracleMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   orchestratorAddressList[0].String(),
		BridgeAddress:   orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:       chainName,
	}
	var err error
	_, err = h(ctx, notOracleMsg)
	require.ErrorIs(t, types.ErrNotOracle, err)
	require.EqualValues(t, types.ErrNotOracle.Error(), err.Error())

	// 2. stake denom not match chain params stake denom
	notMatchStakeDenomMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: "abctoken", Amount: sdk.NewInt(100000)},
		ChainName:       chainName,
	}
	_, err = h(ctx, notMatchStakeDenomMsg)
	require.ErrorIs(t, err, types.ErrBadStakeDenom)
	require.EqualValues(t, fmt.Sprintf("got %s, expected %s: %s", notMatchStakeDenomMsg.DelegateAmount.Denom, fxtypes.DefaultDenom, types.ErrBadStakeDenom), err.Error())

	// 3. insufficient stake amount msg.
	belowMinimumStakeAmountMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:       chainName,
	}
	_, err = h(ctx, belowMinimumStakeAmountMsg)
	require.ErrorIs(t, types.ErrStakeAmountBelowMinimum, err)
	require.EqualValues(t, types.ErrStakeAmountBelowMinimum.Error(), err.Error())

	// 4. success msg
	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	// 5. oracle duplicate set orchestrator
	oracleDuplicateSetOrchestratorMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:       chainName,
	}
	_, err = h(ctx, oracleDuplicateSetOrchestratorMsg)
	require.ErrorIs(t, types.ErrInvalid, err)
	require.EqualValues(t, fmt.Sprintf("oracle existed orchestrator address: %s", types.ErrInvalid.Error()), err.Error())

	// 6. Set the same orchestrator address for different Oracle databases
	duplicateSetOrchestratorMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[1].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:       chainName,
	}
	_, err = h(ctx, duplicateSetOrchestratorMsg)
	require.ErrorIs(t, types.ErrInvalid, err)
	require.EqualValues(t, fmt.Sprintf("orchestrator address is bound to oracle: %s", types.ErrInvalid.Error()), err.Error())

	// 7. Set the same external address for different Oracle databases
	duplicateSetExternalAddressMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[1].String(),
		BridgerAddress:  orchestratorAddressList[1].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:       chainName,
	}
	_, err = h(ctx, duplicateSetExternalAddressMsg)
	require.ErrorIs(t, types.ErrInvalid, err)
	require.EqualValues(t, fmt.Sprintf("external address is bound to oracle: %s", types.ErrInvalid.Error()), err.Error())

	// 8. Margin is not allowed to be submitted more than ten times the threshold
	depositAmountBelowMaximumMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[1].String(),
		BridgerAddress:  orchestratorAddressList[1].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[1].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount.Mul(sdk.NewInt(10).Add(sdk.NewInt(1)))},
		ChainName:       chainName,
	}
	_, err = h(ctx, depositAmountBelowMaximumMsg)
	require.ErrorIs(t, types.ErrStakeAmountBelowMaximum, err)
	require.EqualValues(t, types.ErrStakeAmountBelowMaximum.Error(), err.Error())

	// 9. success msg
	normalMsgOracle2 := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[1].String(),
		BridgerAddress:  orchestratorAddressList[1].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[1].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsgOracle2)
	require.NoError(t, err)
}

// 2. Test MsgAddOracleStake
func TestMsgAddOracleStake(t *testing.T) {
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 4)
	keep := myApp.BscKeeper
	var err error

	// Query the status before the configuration
	totalStakeBefore := keep.GetTotalStake(ctx)
	require.EqualValues(t, sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()), totalStakeBefore)

	// 1. First sets up a valid validator
	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	// Query the totalStake after the address is set
	totalStakeAfter := keep.GetTotalStake(ctx)
	require.True(t, normalMsg.DelegateAmount.IsEqual(totalStakeAfter))

	denomNotMatchMsg := &types.MsgAddOracleStake{
		OracleAddress: oracleAddressList[0].String(),
		Amount: sdk.Coin{
			Denom:  "abc",
			Amount: minStakeAmount,
		},
		ChainName: chainName,
	}
	_, err = h(ctx, denomNotMatchMsg)
	require.ErrorIs(t, err, types.ErrBadStakeDenom)
	require.EqualValues(t, fmt.Sprintf("got %s, expected %s: %s", denomNotMatchMsg.Amount.Denom, fxtypes.DefaultDenom, types.ErrBadStakeDenom), err.Error())

	notOracleMsg := &types.MsgAddOracleStake{
		OracleAddress: orchestratorAddressList[0].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: minStakeAmount,
		},
		ChainName: chainName,
	}
	_, err = h(ctx, notOracleMsg)
	require.ErrorIs(t, types.ErrNotOracle, err)
	require.EqualValues(t, types.ErrNotOracle.Error(), err.Error())

	notSetOrchestratorOracleMsg := &types.MsgAddOracleStake{
		OracleAddress: oracleAddressList[1].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: minStakeAmount,
		},
		ChainName: chainName,
	}
	_, err = h(ctx, notSetOrchestratorOracleMsg)
	require.ErrorIs(t, types.ErrNoOracleFound, err)
	require.EqualValues(t, types.ErrNoOracleFound.Error(), err.Error())

	depositAmountBelowMaximumMsg := &types.MsgAddOracleStake{
		OracleAddress: oracleAddressList[0].String(),
		Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount.Mul(sdk.NewInt(9)).Add(sdk.NewInt(1))},
		ChainName:     chainName,
	}
	_, err = h(ctx, depositAmountBelowMaximumMsg)
	require.ErrorIs(t, types.ErrStakeAmountBelowMaximum, err)
	require.EqualValues(t, types.ErrStakeAmountBelowMaximum.Error(), err.Error())

	normalAddStakeMsg := &types.MsgAddOracleStake{
		OracleAddress: oracleAddressList[0].String(),
		Amount:        sdk.NewCoin(fxtypes.DefaultDenom, minStakeAmount),
		ChainName:     chainName,
	}

	addStake1Before := keep.GetTotalStake(ctx)
	_, err = h(ctx, normalAddStakeMsg)
	require.NoError(t, err)
	addStake1After := keep.GetTotalStake(ctx)
	require.True(t, addStake1Before.Add(normalAddStakeMsg.Amount).IsEqual(addStake1After))
}

func TestMsgSetOracleSetConfirm(t *testing.T) {
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 4)
	keep := myApp.BscKeeper
	var err error

	totalStakeBefore := keep.GetTotalStake(ctx)
	require.EqualValues(t, sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()), totalStakeBefore)

	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	latestOracleSetNonce := keep.GetLatestOracleSetNonce(ctx)
	require.EqualValues(t, 0, latestOracleSetNonce)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
	latestOracleSetNonce = keep.GetLatestOracleSetNonce(ctx)
	require.EqualValues(t, 1, latestOracleSetNonce)

	require.True(t, keep.HasOracleSetRequest(ctx, 1))

	require.False(t, keep.HasOracleSetRequest(ctx, 2))

	nonce1OracleSet := keep.GetOracleSet(ctx, 1)
	require.EqualValues(t, uint64(1), nonce1OracleSet.Nonce)
	require.EqualValues(t, uint64(2), nonce1OracleSet.Height)
	require.EqualValues(t, 1, len(nonce1OracleSet.Members))
	require.EqualValues(t, normalMsg.ExternalAddress, nonce1OracleSet.Members[0].ExternalAddress)
	require.EqualValues(t, math.MaxUint32, nonce1OracleSet.Members[0].Power)

	var gravityId string
	require.NotPanics(t, func() {
		gravityId = keep.GetGravityID(ctx)
	})
	require.EqualValues(t, chainGravityId, gravityId)
	checkpoint := nonce1OracleSet.GetCheckpoint(gravityId)

	external1Signature, err := types.NewEthereumSignature(checkpoint, ethKeys[0])
	require.NoError(t, err)
	external2Signature, err := types.NewEthereumSignature(checkpoint, ethKeys[1])
	require.NoError(t, err)
	errMsgDatas := []struct {
		name      string
		msg       *types.MsgOracleSetConfirm
		err       error
		errReason string
	}{
		{
			name: "Error oracleSet nonce",
			msg: &types.MsgOracleSetConfirm{
				Nonce:               0,
				OrchestratorAddress: orchestratorAddressList[0].String(),
				ExternalAddress:     normalMsg.ExternalAddress,
				Signature:           hex.EncodeToString(external1Signature),
				ChainName:           chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("couldn't find oracleSet: %s", types.ErrInvalid),
		},
		{
			name: "not oracle msg",
			msg: &types.MsgOracleSetConfirm{
				Nonce:               nonce1OracleSet.Nonce,
				OrchestratorAddress: orchestratorAddressList[0].String(),
				ExternalAddress:     crypto.PubkeyToAddress(ethKeys[1].PublicKey).Hex(),
				Signature:           hex.EncodeToString(external1Signature),
				ChainName:           chainName,
			},
			err:       types.ErrNotOracle,
			errReason: fmt.Sprintf("%s", types.ErrNotOracle),
		},
		{
			name: "sign not match external-1  external-sign-2",
			msg: &types.MsgOracleSetConfirm{
				Nonce:               nonce1OracleSet.Nonce,
				OrchestratorAddress: orchestratorAddressList[0].String(),
				ExternalAddress:     crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
				Signature:           hex.EncodeToString(external2Signature),
				ChainName:           chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s: %s", crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(), hex.EncodeToString(checkpoint), hex.EncodeToString(external2Signature), types.ErrInvalid),
		},
		{
			name: "orchestrator address not match",
			msg: &types.MsgOracleSetConfirm{
				Nonce:               nonce1OracleSet.Nonce,
				OrchestratorAddress: orchestratorAddressList[1].String(),
				ExternalAddress:     crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
				Signature:           hex.EncodeToString(external1Signature),
				ChainName:           chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("got %s, expected %s: %s", orchestratorAddressList[1].String(), orchestratorAddressList[0].String(), types.ErrInvalid),
		},
	}

	for _, testData := range errMsgDatas {
		_, err = h(ctx, testData.msg)
		require.ErrorIs(t, err, testData.err, testData.name)
		require.EqualValues(t, err.Error(), testData.errReason, testData.name)
	}

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:               nonce1OracleSet.Nonce,
		OrchestratorAddress: orchestratorAddressList[0].String(),
		ExternalAddress:     crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		Signature:           hex.EncodeToString(external1Signature),
		ChainName:           chainName,
	}
	_, err = h(ctx, normalOracleSetConfirmMsg)
	require.NoError(t, err)

	endBlockBeforeLatestOracleSet := keep.GetLatestOracleSet(ctx)
	require.NotNil(t, endBlockBeforeLatestOracleSet)
}

func TestClaimWithOracleJailed(t *testing.T) {
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 10)
	keeper := myApp.BscKeeper
	var err error

	totalStakeBefore := keeper.GetTotalStake(ctx)
	require.EqualValues(t, sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()), totalStakeBefore)

	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	latestOracleSetNonce := keeper.GetLatestOracleSetNonce(ctx)
	require.EqualValues(t, 1, latestOracleSetNonce)

	nonce1OracleSet := keeper.GetOracleSet(ctx, latestOracleSetNonce)
	require.EqualValues(t, uint64(1), nonce1OracleSet.Nonce)
	require.EqualValues(t, uint64(2), nonce1OracleSet.Height)

	var gravityId string
	require.NotPanics(t, func() {
		gravityId = keeper.GetGravityID(ctx)
	})
	require.EqualValues(t, chainGravityId, gravityId)
	checkpoint := nonce1OracleSet.GetCheckpoint(gravityId)

	// oracle jailed!!!
	oracle, found := keeper.GetOracle(ctx, oracleAddressList[0])
	require.True(t, found)
	oracle.Jailed = true
	keeper.SetOracle(ctx, oracle)

	external1Signature, err := types.NewEthereumSignature(checkpoint, ethKeys[0])
	require.NoError(t, err)

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:               latestOracleSetNonce,
		OrchestratorAddress: orchestratorAddressList[0].String(),
		ExternalAddress:     crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		Signature:           hex.EncodeToString(external1Signature),
		ChainName:           chainName,
	}
	_, err = h(ctx, normalOracleSetConfirmMsg)
	require.Nil(t, err)
}

func TestClaimTest(t *testing.T) {
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 10)
	var err error

	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgerAddress:  orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	oracleLastEventNonce := myApp.BscKeeper.GetLastEventNonceByOracle(ctx, oracleAddressList[0])
	require.EqualValues(t, 0, oracleLastEventNonce)

	errMsgDatas := []struct {
		name      string
		msg       *types.MsgBridgeTokenClaim
		err       error
		errReason string
	}{
		{
			name: "error oracleSet nonce: 2",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     2,
				BlockHeight:    1,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: orchestratorAddressList[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 2, 1, types.ErrNonContiguousEventNonce),
		},
		{
			name: "error oracleSet nonce: 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    1,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: orchestratorAddressList[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 3, 1, types.ErrNonContiguousEventNonce),
		},
		{
			name: "Normal claim msg: 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    1,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: orchestratorAddressList[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      chainName,
			},
			err:       nil,
			errReason: "",
		},
		{
			name: "error oracleSet nonce: 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    2,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: orchestratorAddressList[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 1, 2, types.ErrNonContiguousEventNonce),
		},
		{
			name: "error oracleSet nonce: 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    2,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: orchestratorAddressList[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 3, 2, types.ErrNonContiguousEventNonce),
		},
	}

	for _, testData := range errMsgDatas {
		_, err = h(ctx, testData.msg)
		require.ErrorIs(t, err, testData.err, testData.name)
		if err == nil {
			continue
		}
		require.EqualValues(t, testData.errReason, err.Error(), testData.name)
	}

}

// Test Support RequestBatch baseFee
func TestSupportRequestBatchBaseFee(t *testing.T) {
	//myApp.SetAppLog(server.ZeroLogWrapper{Logger: log.Logger.Level(zerolog.DebugLevel)})
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 10)
	keep := myApp.BscKeeper
	var err error

	// Query the status before the configuration
	totalStakeBefore := keep.GetTotalStake(ctx)
	require.EqualValues(t, sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()), totalStakeBefore)

	endBlock := func() {
		//ctx = ctx.WithBlockHeight(fxtypes.CrossChainSupportBscBlock() + 1)
		crosschain.EndBlocker(ctx, keep)
	}

	// 1. First sets up a valid validator
	for i, oracle := range oracleAddressList {
		normalMsg := &types.MsgCreateOracleBridger{
			OracleAddress:   oracle.String(),
			BridgeAddress:   orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
			ChainName:       chainName,
		}
		_, err = h(ctx, normalMsg)
		require.NoError(t, err)
	}

	endBlock()

	var externalOracleMembers []*types.BridgeValidator
	for _, key := range ethKeys {
		externalOracleMembers = append(externalOracleMembers, &types.BridgeValidator{
			Power:           100,
			ExternalAddress: crypto.PubkeyToAddress(key.PublicKey).Hex(),
		})
	}

	// 2. oracle update claim
	for i := range oracleAddressList {
		normalMsg := &types.MsgOracleSetUpdatedClaim{
			EventNonce:     1,
			BlockHeight:    1,
			OracleSetNonce: 1,
			Members:        externalOracleMembers,
			BridgerAddress: orchestratorAddressList[i].String(),
			ChainName:      chainName,
		}
		_, err = h(ctx, normalMsg)
		require.NoError(t, err)
	}

	endBlock()

	// 3. add bridge token.
	sendToFxSendAddr := crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex()
	sendToFxReceiveAddr := orchestratorAddressList[0]
	sendToFxAmount := sdk.NewIntWithDecimal(1000, 18)
	sendToFxToken := "0x0000000000000000000000000000000000001000"

	for i, oracle := range oracleAddressList {
		normalMsg := &types.MsgBridgeTokenClaim{
			EventNonce:     keep.GetLastEventNonceByOracle(ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Name:           "BSC USDT",
			Symbol:         "USDT",
			Decimals:       18,
			BridgerAddress: orchestratorAddressList[i].String(),
			ChannelIbc:     "",
			ChainName:      chainName,
		}
		_, err = h(ctx, normalMsg)
		require.NoError(t, err)
	}

	endBlock()

	bridgeDenomData := keep.GetBridgeTokenDenom(ctx, sendToFxToken)
	require.NotNil(t, bridgeDenomData)
	tokenDenom := fmt.Sprintf("%s%s", chainName, sendToFxToken)
	require.EqualValues(t, tokenDenom, bridgeDenomData.Denom)
	bridgeTokenData := keep.GetDenomByBridgeToken(ctx, tokenDenom)
	require.NotNil(t, bridgeTokenData)
	require.EqualValues(t, sendToFxToken, bridgeTokenData.Token)

	// 4. sendToFx.
	for i, oracle := range oracleAddressList {
		normalMsg := &types.MsgSendToFxClaim{
			EventNonce:     keep.GetLastEventNonceByOracle(ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Amount:         sendToFxAmount,
			Sender:         sendToFxSendAddr,
			Receiver:       sendToFxReceiveAddr.String(),
			TargetIbc:      "",
			BridgerAddress: orchestratorAddressList[i].String(),
			ChainName:      chainName,
		}
		_, err = h(ctx, normalMsg)
		require.NoError(t, err)
	}

	endBlock()

	balance := myApp.BankKeeper.GetBalance(ctx, sendToFxReceiveAddr, tokenDenom)
	require.NotNil(t, balance)
	require.EqualValues(t, balance.Denom, tokenDenom)
	require.True(t, balance.Amount.Equal(sendToFxAmount))

	sendToExternal := func(bridgeFees []sdk.Int) {
		for _, bridgeFee := range bridgeFees {
			sendToExternal := &types.MsgSendToExternal{
				Sender:    sendToFxReceiveAddr.String(),
				Dest:      sendToFxSendAddr,
				Amount:    sdk.NewCoin(tokenDenom, sdk.NewInt(3)),
				BridgeFee: sdk.NewCoin(tokenDenom, bridgeFee),
				ChainName: chainName,
			}
			_, err = h(ctx, sendToExternal)
			require.NoError(t, err)
		}
	}

	sendToExternal([]sdk.Int{sdk.NewInt(1), sdk.NewInt(2), sdk.NewInt(3)})
	usdtBatchFee := keep.GetBatchFeesByTokenType(ctx, sendToFxToken, 100, sdk.NewInt(0))
	require.EqualValues(t, sendToFxToken, usdtBatchFee.TokenContract)
	require.EqualValues(t, 3, usdtBatchFee.TotalTxs)
	require.EqualValues(t, sdk.NewInt(6), usdtBatchFee.TotalFees)

	fn := func(i sdk.Int) *sdk.Int {
		return &i
	}

	testCases := []struct {
		testName       string
		baseFee        *sdk.Int
		pass           bool
		expectTotalTxs uint64
		err            error
	}{
		{
			testName:       "Support - baseFee 1000",
			baseFee:        fn(sdk.NewInt(1000)),
			pass:           false,
			expectTotalTxs: 3,
			err:            types.ErrEmpty,
		},
		{
			testName:       "Support - baseFee 2",
			baseFee:        fn(sdk.NewInt(2)),
			pass:           true,
			expectTotalTxs: 1,
			err:            nil,
		},
		{
			testName:       "Support - baseFee 0",
			baseFee:        fn(sdk.NewInt(0)),
			pass:           true,
			expectTotalTxs: 0,
			err:            nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			cacheCtx, _ := ctx.CacheContext()
			_, err = h(cacheCtx, &types.MsgRequestBatch{
				Sender:     orchestratorAddressList[0].String(),
				Denom:      tokenDenom,
				MinimumFee: sdk.NewInt(1),
				FeeReceive: "0x0000000000000000000000000000000000000000",
				ChainName:  chainName,
				BaseFee:    testCase.baseFee,
			})
			if testCase.pass {
				require.NoError(t, err)
				usdtBatchFee = keep.GetBatchFeesByTokenType(cacheCtx, sendToFxToken, 100, sdk.NewInt(0))
				require.EqualValues(t, testCase.expectTotalTxs, usdtBatchFee.TotalTxs)
				return
			}

			require.NotNil(t, err)
			require.True(t, errors.As(err, &testCase.err))
			require.Equal(t, err, testCase.err)
		})
	}
}

func createTestEnv(t *testing.T, generateAccountNum int) (myApp *app.App, ctx sdk.Context, oracleAddressList, orchestratorAddressList []sdk.AccAddress, ethKeys []*ecdsa.PrivateKey, handler sdk.Handler) {
	fxtypes.ChangeNetworkForTest(fxtypes.NetworkDevnet())

	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := helpers.GenerateGenesisValidator(t, 2, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initBalances)))
	myApp = helpers.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx = myApp.BaseApp.NewContext(false, tmproto.Header{Height: 1})
	oracleAddressList = helpers.AddTestAddrsIncremental(myApp, ctx, generateAccountNum, minStakeAmount.Mul(sdk.NewInt(1000)))
	orchestratorAddressList = helpers.AddTestAddrs(myApp, ctx, generateAccountNum, sdk.ZeroInt())
	ethKeys = genEthKey(generateAccountNum)
	// chain module oracle list
	var oracles []string
	for _, account := range oracleAddressList {
		oracles = append(oracles, account.String())
	}

	var err error
	// init bsc params by proposal
	proposalHandler := crosschain.NewCrossChainProposalHandler(myApp.CrosschainKeeper)
	err = proposalHandler(ctx, &types.InitCrossChainParamsProposal{
		Title:       "init bsc chain params",
		Description: "init fx chain <-> bsc chain params",
		Params:      defaultModuleParams(oracles),
		ChainName:   chainName,
	})
	require.NoError(t, err)

	crosschianHandler := crosschain.NewHandler(myApp.CrosschainKeeper)

	proxyHandler := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		require.NoError(t, msg.ValidateBasic(), fmt.Sprintf("msg %s validate basic error", sdk.MsgTypeURL(msg)))
		return crosschianHandler(ctx, msg)
	}
	return myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, proxyHandler
}

func defaultModuleParams(oracles []string) *types.Params {
	return &types.Params{
		GravityId:                         chainGravityId,
		SignedWindow:                      20000,
		ExternalBatchTimeout:              43200000,
		AverageBlockTime:                  5000,
		AverageExternalBlockTime:          3000,
		SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		IbcTransferTimeoutHeight:          10000,
		OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		Oracles:                           oracles,
		StakeThreshold:                    sdk.NewCoin(fxtypes.DefaultDenom, minStakeAmount),
	}
}

func genEthKey(count int) []*ecdsa.PrivateKey {
	var ethKeys []*ecdsa.PrivateKey
	for i := 0; i < count; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		ethKeys = append(ethKeys, key)
	}
	return ethKeys
}
