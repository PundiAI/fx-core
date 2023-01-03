package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	types2 "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *KeeperTestSuite) TestMsgBondedOracle() {

	// 1. sender not in chain oracle
	notOracleMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.bridgers[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(rand.Int63()),
		},
		ChainName: suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), notOracleMsg)
	require.ErrorIs(suite.T(), err, types.ErrNoFoundOracle)

	// 2. stake denom not match chain params stake denom
	notMatchStakeDenomMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  "abctoken",
			Amount: sdk.NewInt(rand.Int63()),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), notMatchStakeDenomMsg)
	require.ErrorIs(suite.T(), err, types.ErrInvalid)
	require.EqualValues(suite.T(), fmt.Sprintf(
		"delegate denom got %s, expected %s: %s",
		notMatchStakeDenomMsg.DelegateAmount.Denom, fxtypes.DefaultDenom, types.ErrInvalid), err.Error())

	// 3. insufficient stake amount msg.
	belowMinimumStakeAmountMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(rand.Int63()),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), belowMinimumStakeAmountMsg)
	require.ErrorIs(suite.T(), types.ErrDelegateAmountBelowMinimum, err)
	require.EqualValues(suite.T(), types.ErrDelegateAmountBelowMinimum.Error(), err.Error())

	// 4. success msg
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(3) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	// 5. oracle duplicate set bridger
	oracleDuplicateBondedMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(rand.Int63()),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), oracleDuplicateBondedMsg)
	require.ErrorIs(suite.T(), types.ErrInvalid, err)
	require.EqualValues(suite.T(),
		fmt.Sprintf("oracle existed bridger address: %s", types.ErrInvalid.Error()), err.Error())

	// 6. Set the same bridger tronAddress for different Oracle databases
	duplicateSetBridgerMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(rand.Int63()),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), duplicateSetBridgerMsg)
	require.ErrorIs(suite.T(), types.ErrInvalid, err)
	require.EqualValues(suite.T(),
		fmt.Sprintf("bridger address is bound to oracle: %s", types.ErrInvalid.Error()), err.Error())

	// 7. Set the same external tronAddress for different Oracle databases
	duplicateSetExternalAddressMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[1].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(rand.Int63()),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), duplicateSetExternalAddressMsg)
	require.ErrorIs(suite.T(), types.ErrInvalid, err)
	require.EqualValues(suite.T(),
		fmt.Sprintf("external address is bound to oracle: %s", types.ErrInvalid.Error()), err.Error())

	// 8. Margin is not allowed to be submitted more than ten times the threshold
	depositAmountBelowMaximumMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[1].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[1].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(10_100).MulRaw(1e18).Mul(sdk.NewInt(10).Add(sdk.NewInt(1))),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), depositAmountBelowMaximumMsg)
	require.ErrorIs(suite.T(), types.ErrDelegateAmountAboveMaximum, err)
	require.EqualValues(suite.T(), types.ErrDelegateAmountAboveMaximum.Error(), err.Error())

	// 9. success msg
	normalMsgOracle2 := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[1].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[1].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsgOracle2)
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestMsgAddDelegate() {

	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	denomNotMatchMsg := &types.MsgAddDelegate{
		OracleAddress: suite.oracles[0].String(),
		Amount: sdk.Coin{
			Denom:  "abc",
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), denomNotMatchMsg)
	require.ErrorIs(suite.T(), err, types.ErrInvalid)
	require.EqualValues(suite.T(), fmt.Sprintf("delegate denom got %s, expected %s: %s",
		denomNotMatchMsg.Amount.Denom, fxtypes.DefaultDenom, types.ErrInvalid), err.Error())

	notOracleMsg := &types.MsgAddDelegate{
		OracleAddress: suite.bridgers[0].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), notOracleMsg)
	require.ErrorIs(suite.T(), types.ErrNoFoundOracle, err)
	require.EqualValues(suite.T(), types.ErrNoFoundOracle.Error(), err.Error())

	notSetBridgerOracleMsg := &types.MsgAddDelegate{
		OracleAddress: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), notSetBridgerOracleMsg)
	require.ErrorIs(suite.T(), types.ErrNoFoundOracle, err)
	require.EqualValues(suite.T(), types.ErrNoFoundOracle.Error(), err.Error())

	depositAmountBelowMaximumMsg := &types.MsgAddDelegate{
		OracleAddress: suite.oracles[0].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(10_000).MulRaw(1e18).Mul(sdk.NewInt(9)).Add(sdk.NewInt(1)),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), depositAmountBelowMaximumMsg)
	require.ErrorIs(suite.T(), types.ErrDelegateAmountAboveMaximum, err)
	require.EqualValues(suite.T(), types.ErrDelegateAmountAboveMaximum.Error(), err.Error())

	normalAddStakeMsg := &types.MsgAddDelegate{
		OracleAddress: suite.oracles[0].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt(1),
		},
		ChainName: suite.chainName,
	}
	_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), normalAddStakeMsg)
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestMsgEditBridger() {
	for i := range suite.oracles {
		bondedMsg := &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount: sdk.Coin{
				Denom:  fxtypes.DefaultDenom,
				Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
			},
			ChainName: suite.chainName,
		}
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), bondedMsg)
		suite.NoError(err)
	}

	token := helpers.GenerateAddress().Hex()
	denom := fmt.Sprintf("%s%s", suite.chainName, token)
	suite.Keeper().AddBridgeToken(suite.ctx, token, denom)

	sendToMsg := &types.MsgSendToFxClaim{
		EventNonce:    1,
		BlockHeight:   100,
		TokenContract: token,
		Amount:        sdk.NewInt(int64(rand.Uint32())),
		Sender:        helpers.GenerateAddress().Hex(),
		Receiver:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		TargetIbc:     "",
		ChainName:     suite.chainName,
	}
	for i := 0; i < len(suite.bridgers)/2; i++ {
		sendToMsg.BridgerAddress = suite.bridgers[i].String()
		_, err := suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToMsg)
		suite.NoError(err)
	}

	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: types2.Header{ChainID: suite.ctx.ChainID(), Height: suite.ctx.BlockHeight()}})

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToMsg.Receiver))
	suite.Equal(balances.String(), sdk.NewCoins().String(), len(suite.bridgers))

	for i := 0; i < len(suite.oracles); i++ {
		_, err := suite.MsgServer().EditBridger(sdk.WrapSDKContext(suite.ctx), &types.MsgEditBridger{
			ChainName:      suite.chainName,
			OracleAddress:  suite.oracles[i].String(),
			BridgerAddress: suite.bridgers[i].String(),
		})
		suite.Require().Error(err)

		_, err = suite.MsgServer().EditBridger(sdk.WrapSDKContext(suite.ctx), &types.MsgEditBridger{
			ChainName:      suite.chainName,
			OracleAddress:  suite.oracles[i].String(),
			BridgerAddress: sdk.AccAddress(suite.validator[i]).String(),
		})
		suite.NoError(err)

		sendToMsg.BridgerAddress = sdk.AccAddress(suite.validator[i]).String()
		_, err = suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToMsg)
		if i < len(suite.oracles)/2 {
			suite.ErrorContains(err, types.ErrNonContiguousEventNonce.Error())
		} else {
			suite.NoError(err)
		}
	}
	for _, bridger := range suite.bridgers {
		_, found := suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
		suite.False(found)
	}

	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToMsg.Receiver))
	suite.Equal(balances.String(), sdk.NewCoins(sdk.NewCoin(denom, sendToMsg.Amount)).String())
}

func (suite *KeeperTestSuite) TestMsgSetOracleSetConfirm() {

	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	if trontypes.ModuleName == suite.chainName {
		normalMsg.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[0].PublicKey).String()
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	require.EqualValues(suite.T(), 0, latestOracleSetNonce)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	latestOracleSetNonce = suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	require.EqualValues(suite.T(), 1, latestOracleSetNonce)

	require.True(suite.T(), suite.Keeper().HasOracleSetRequest(suite.ctx, 1))

	require.False(suite.T(), suite.Keeper().HasOracleSetRequest(suite.ctx, 2))

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.ctx, 1)
	require.EqualValues(suite.T(), uint64(1), nonce1OracleSet.Nonce)
	require.EqualValues(suite.T(), uint64(2), nonce1OracleSet.Height)
	require.EqualValues(suite.T(), 1, len(nonce1OracleSet.Members))
	require.EqualValues(suite.T(), normalMsg.ExternalAddress, nonce1OracleSet.Members[0].ExternalAddress)
	require.EqualValues(suite.T(), math.MaxUint32, nonce1OracleSet.Members[0].Power)

	gravityId := suite.Keeper().GetGravityID(suite.ctx)
	checkpoint, err := nonce1OracleSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == suite.chainName {
		checkpoint, err = trontypes.GetCheckpointOracleSet(nonce1OracleSet, gravityId)
	}
	require.NoError(suite.T(), err)

	external1Signature, err := types.NewEthereumSignature(checkpoint, suite.externals[0])
	if trontypes.ModuleName == suite.chainName {
		external1Signature, err = trontypes.NewTronSignature(checkpoint, suite.externals[0])
	}
	require.NoError(suite.T(), err)
	external2Signature, err := types.NewEthereumSignature(checkpoint, suite.externals[1])
	if trontypes.ModuleName == suite.chainName {
		external2Signature, err = trontypes.NewTronSignature(checkpoint, suite.externals[1])
	}
	require.NoError(suite.T(), err)

	errMsgData := []struct {
		name      string
		msg       *types.MsgOracleSetConfirm
		err       error
		errReason string
	}{
		{
			name: "Error oracleSet nonce",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           0,
				BridgerAddress:  suite.bridgers[0].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("couldn't find oracleSet: %s", types.ErrInvalid),
		},
		{
			name: "not oracle msg",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgers[0].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[1].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external2Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrNoFoundOracle,
			errReason: fmt.Sprintf("%s", types.ErrNoFoundOracle),
		},
		{
			name: "sign not match external-1  external-sign-2",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgers[0].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external2Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s: %s", normalMsg.ExternalAddress, hex.EncodeToString(checkpoint), hex.EncodeToString(external2Signature), types.ErrInvalid),
		},
		{
			name: "bridger tronAddress not match",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgers[1].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("got %s, expected %s: %s", suite.bridgers[1].String(), suite.bridgers[0].String(), types.ErrInvalid),
		},
	}

	for _, testData := range errMsgData {
		if trontypes.ModuleName == suite.chainName {
			testData.msg.ExternalAddress = trontypes.AddressFromHex(testData.msg.ExternalAddress)
		}
		_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), testData.msg)
		require.ErrorIs(suite.T(), err, testData.err, testData.name)
		require.EqualValues(suite.T(), err.Error(), testData.errReason, testData.name)
	}

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:           nonce1OracleSet.Nonce,
		BridgerAddress:  suite.bridgers[0].String(),
		ExternalAddress: normalMsg.ExternalAddress,
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       suite.chainName,
	}
	_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), normalOracleSetConfirmMsg)
	require.NoError(suite.T(), err)

	endBlockBeforeLatestOracleSet := suite.Keeper().GetLatestOracleSet(suite.ctx)
	require.NotNil(suite.T(), endBlockBeforeLatestOracleSet)
}

func (suite *KeeperTestSuite) TestClaimWithOracleOnline() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	if trontypes.ModuleName == suite.chainName {
		normalMsg.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[0].PublicKey).String()
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	require.EqualValues(suite.T(), 1, latestOracleSetNonce)

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.ctx, latestOracleSetNonce)
	require.EqualValues(suite.T(), uint64(1), nonce1OracleSet.Nonce)
	require.EqualValues(suite.T(), uint64(2), nonce1OracleSet.Height)

	var gravityId string
	require.NotPanics(suite.T(), func() {
		gravityId = suite.Keeper().GetGravityID(suite.ctx)
	})
	if suite.chainName == ethtypes.ModuleName {
		require.EqualValues(suite.T(), fmt.Sprintf("fx-bridge-%s", suite.chainName), gravityId)
	} else {
		require.EqualValues(suite.T(), fmt.Sprintf("fx-%s-bridge", suite.chainName), gravityId)
	}
	checkpoint, err := nonce1OracleSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == suite.chainName {
		checkpoint, err = trontypes.GetCheckpointOracleSet(nonce1OracleSet, gravityId)
	}
	require.NoError(suite.T(), err)

	// oracle Online!!!
	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[0])
	require.True(suite.T(), found)
	oracle.Online = true
	suite.Keeper().SetOracle(suite.ctx, oracle)

	external1Signature, err := types.NewEthereumSignature(checkpoint, suite.externals[0])
	if trontypes.ModuleName == suite.chainName {
		external1Signature, err = trontypes.NewTronSignature(checkpoint, suite.externals[0])
	}
	require.NoError(suite.T(), err)

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:           latestOracleSetNonce,
		BridgerAddress:  suite.bridgers[0].String(),
		ExternalAddress: normalMsg.ExternalAddress,
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       suite.chainName,
	}
	_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), normalOracleSetConfirmMsg)
	require.Nil(suite.T(), err)
}

func (suite *KeeperTestSuite) TestClaimTest() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
		},
		ChainName: suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracles[0])
	require.EqualValues(suite.T(), 0, oracleLastEventNonce)

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
				TokenContract:  helpers.GenerateAddress().String(),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 2, 1, types.ErrNonContiguousEventNonce),
		},
		{
			name: "error oracleSet nonce: 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    1,
				TokenContract:  helpers.GenerateAddress().String(),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 3, 1, types.ErrNonContiguousEventNonce),
		},
		{
			name: "Normal claim msg: 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    1,
				TokenContract:  helpers.GenerateAddress().String(),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       nil,
			errReason: "",
		},
		{
			name: "error oracleSet nonce: 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    2,
				TokenContract:  helpers.GenerateAddress().String(),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 1, 2, types.ErrNonContiguousEventNonce),
		},
		{
			name: "error oracleSet nonce: 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    2,
				TokenContract:  helpers.GenerateAddress().String(),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 3, 2, types.ErrNonContiguousEventNonce),
		},
	}

	for _, testData := range errMsgDatas {
		if testData.msg.ChainName == trontypes.ModuleName {
			testData.msg.TokenContract = trontypes.AddressFromHex(testData.msg.TokenContract)
		}
		err := testData.msg.ValidateBasic()
		require.NoError(suite.T(), err)
		_, err = suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), testData.msg)
		require.ErrorIs(suite.T(), err, testData.err, testData.name)
		if err == nil {
			continue
		}
		require.EqualValues(suite.T(), testData.errReason, err.Error(), testData.name)
	}

}

func (suite *KeeperTestSuite) TestRequestBatchBaseFee() {

	// 1. First sets up a valid validator
	var totalPower = sdk.ZeroInt()
	var delegateAmounts []sdk.Int
	for i, oracle := range suite.oracles {
		normalMsg := &types.MsgBondedOracle{
			OracleAddress:    oracle.String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[0].String(),
			DelegateAmount: sdk.Coin{
				Denom:  fxtypes.DefaultDenom,
				Amount: sdk.NewInt((rand.Int63n(5) + 1) * 10_000).MulRaw(1e18),
			},
			ChainName: suite.chainName,
		}
		if len(suite.validator) > i {
			normalMsg.ValidatorAddress = suite.validator[i].String()
		}
		delegateAmounts = append(delegateAmounts, normalMsg.DelegateAmount.Amount)
		totalPower = totalPower.Add(normalMsg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction))
		if trontypes.ModuleName == suite.chainName {
			normalMsg.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()
		}
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	var externalOracleMembers types.BridgeValidators
	for i, key := range suite.externals {
		power := delegateAmounts[i].Quo(sdk.DefaultPowerReduction).MulRaw(math.MaxUint32).Quo(totalPower)
		bridgeVal := types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: crypto.PubkeyToAddress(key.PublicKey).Hex(),
		}
		if trontypes.ModuleName == suite.chainName {
			bridgeVal.ExternalAddress = tronAddress.PubkeyToAddress(key.PublicKey).String()
		}
		externalOracleMembers = append(externalOracleMembers, bridgeVal)
	}
	sort.Sort(externalOracleMembers)

	// 2. oracle update claim
	for i := range suite.oracles {
		normalMsg := &types.MsgOracleSetUpdatedClaim{
			EventNonce:     1,
			BlockHeight:    1,
			OracleSetNonce: 1,
			Members:        externalOracleMembers,
			BridgerAddress: suite.bridgers[i].String(),
			ChainName:      suite.chainName,
		}
		_, err := suite.MsgServer().OracleSetUpdateClaim(sdk.WrapSDKContext(suite.ctx), normalMsg)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	// 3. add bridge token.
	sendToFxSendAddr := crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex()
	sendToFxReceiveAddr := suite.bridgers[0]
	sendToFxAmount := sdk.NewIntWithDecimal(1000, 18)
	sendToFxToken := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		sendToFxToken = trontypes.AddressFromHex(sendToFxToken)
		sendToFxSendAddr = tronAddress.PubkeyToAddress(suite.externals[0].PublicKey).String()
	}

	for i, oracle := range suite.oracles {
		normalMsg := &types.MsgBridgeTokenClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Name:           "Test USDT",
			Symbol:         "USDT",
			Decimals:       18,
			BridgerAddress: suite.bridgers[i].String(),
			ChannelIbc:     "",
			ChainName:      suite.chainName,
		}
		_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), normalMsg)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	bridgeDenomData := suite.Keeper().GetBridgeTokenDenom(suite.ctx, sendToFxToken)
	require.NotNil(suite.T(), bridgeDenomData)
	tokenDenom := fmt.Sprintf("%s%s", suite.chainName, sendToFxToken)
	require.EqualValues(suite.T(), tokenDenom, bridgeDenomData.Denom)
	bridgeTokenData := suite.Keeper().GetDenomByBridgeToken(suite.ctx, tokenDenom)
	require.NotNil(suite.T(), bridgeTokenData)
	require.EqualValues(suite.T(), sendToFxToken, bridgeTokenData.Token)

	// 4. sendToFx.
	for i, oracle := range suite.oracles {
		normalMsg := &types.MsgSendToFxClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Amount:         sendToFxAmount,
			Sender:         sendToFxSendAddr,
			Receiver:       sendToFxReceiveAddr.String(),
			TargetIbc:      "",
			BridgerAddress: suite.bridgers[i].String(),
			ChainName:      suite.chainName,
		}
		_, err := suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), normalMsg)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	balance := suite.app.BankKeeper.GetBalance(suite.ctx, sendToFxReceiveAddr, tokenDenom)
	require.NotNil(suite.T(), balance)
	require.EqualValues(suite.T(), balance.Denom, tokenDenom)
	require.True(suite.T(), balance.Amount.Equal(sendToFxAmount))

	sendToExternal := func(bridgeFees []sdk.Int) {
		for _, bridgeFee := range bridgeFees {
			sendToExternal := &types.MsgSendToExternal{
				Sender:    sendToFxReceiveAddr.String(),
				Dest:      sendToFxSendAddr,
				Amount:    sdk.NewCoin(tokenDenom, sdk.NewInt(3)),
				BridgeFee: sdk.NewCoin(tokenDenom, bridgeFee),
				ChainName: suite.chainName,
			}
			_, err := suite.MsgServer().SendToExternal(sdk.WrapSDKContext(suite.ctx), sendToExternal)
			require.NoError(suite.T(), err)
		}
	}

	sendToExternal([]sdk.Int{sdk.NewInt(1), sdk.NewInt(2), sdk.NewInt(3)})
	usdtBatchFee := suite.Keeper().GetBatchFeesByTokenType(suite.ctx, sendToFxToken, 100, sdk.NewInt(0))
	require.EqualValues(suite.T(), sendToFxToken, usdtBatchFee.TokenContract)
	require.EqualValues(suite.T(), 3, usdtBatchFee.TotalTxs)
	require.EqualValues(suite.T(), sdk.NewInt(6), usdtBatchFee.TotalFees)

	testCases := []struct {
		testName       string
		baseFee        sdk.Int
		pass           bool
		expectTotalTxs uint64
		err            error
	}{
		{
			testName:       "Support - baseFee 1000",
			baseFee:        sdk.NewInt(1000),
			pass:           false,
			expectTotalTxs: 3,
			err:            sdkerrors.Wrap(types.ErrEmpty, "no batch tx"),
		},
		{
			testName:       "Support - baseFee 2",
			baseFee:        sdk.NewInt(2),
			pass:           true,
			expectTotalTxs: 1,
			err:            nil,
		},
		{
			testName:       "Support - baseFee 0",
			baseFee:        sdk.NewInt(0),
			pass:           false,
			expectTotalTxs: 0,
			err:            sdkerrors.Wrap(types.ErrInvalid, "new batch would not be more profitable"),
		},
	}

	for _, testCase := range testCases {
		_, err := suite.MsgServer().RequestBatch(sdk.WrapSDKContext(suite.ctx), &types.MsgRequestBatch{
			Sender:     suite.bridgers[0].String(),
			Denom:      tokenDenom,
			MinimumFee: sdk.NewInt(1),
			FeeReceive: "0x0000000000000000000000000000000000000000",
			ChainName:  suite.chainName,
			BaseFee:    testCase.baseFee,
		})
		if testCase.pass {
			require.NoError(suite.T(), err)
			usdtBatchFee = suite.Keeper().GetBatchFeesByTokenType(suite.ctx, sendToFxToken, 100, sdk.NewInt(0))
			require.EqualValues(suite.T(), testCase.expectTotalTxs, usdtBatchFee.TotalTxs)
		} else {
			require.NotNil(suite.T(), err)
			require.Equal(suite.T(), err.Error(), testCase.err.Error())
		}
	}
}
