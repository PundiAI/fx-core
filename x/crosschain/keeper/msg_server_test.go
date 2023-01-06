package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	tmrand "github.com/tendermint/tendermint/libs/rand"

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
	testCases := []struct {
		name   string
		pass   bool
		err    string
		preRun func(msg *types.MsgBondedOracle)
	}{
		{
			name: "error - sender not oracle",
			preRun: func(msg *types.MsgBondedOracle) {
				msg.OracleAddress = msg.BridgerAddress
			},
			pass: false,
			err:  types.ErrNoFoundOracle.Error(),
		},
		{
			name: "error - oracle existed",
			preRun: func(msg *types.MsgBondedOracle) {
				suite.Keeper().SetOracle(suite.ctx, types.Oracle{OracleAddress: msg.OracleAddress})
			},
			pass: false,
			err:  "oracle existed bridger address: invalid",
		},
		{
			name: "error - bridger address is bound",
			preRun: func(msg *types.MsgBondedOracle) {
				suite.Keeper().SetOracleByBridger(suite.ctx, sdk.MustAccAddressFromBech32(msg.BridgerAddress), sdk.MustAccAddressFromBech32(msg.OracleAddress))
			},
			pass: false,
			err:  "bridger address is bound to oracle: invalid",
		},
		{
			name: "error - external address is bound",
			preRun: func(msg *types.MsgBondedOracle) {
				suite.Keeper().SetOracleByExternalAddress(suite.ctx, msg.ExternalAddress, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			},
			pass: false,
			err:  "external address is bound to oracle: invalid",
		},
		{
			name: "error - stake denom not match chain params stake denom",
			preRun: func(msg *types.MsgBondedOracle) {
				msg.DelegateAmount.Denom = "stake"
			},
			pass: false,
			err:  fmt.Sprintf("delegate denom got %s, expected %s: invalid", "stake", "FX"),
		},
		{
			name: "error - delegate amount less than threshold amount",
			preRun: func(msg *types.MsgBondedOracle) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
				msg.DelegateAmount.Amount = delegateThreshold.Amount.Sub(sdk.NewInt(rand.Int63() - 1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgBondedOracle) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
				delegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.ctx)
				maxDelegateAmount := delegateThreshold.Amount.Mul(sdk.NewInt(delegateMultiple))
				msg.DelegateAmount.Amount = maxDelegateAmount.Add(sdk.NewInt(rand.Int63() - 1))
			},
			pass: false,
			err:  types.ErrDelegateAmountAboveMaximum.Error(),
		},
		{
			name: "pass",
			preRun: func(msg *types.MsgBondedOracle) {
			},
			pass: true,
		},
	}
	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			suite.SetupTest()
			oracleIndex := rand.Intn(len(suite.oracles))
			msg := &types.MsgBondedOracle{
				OracleAddress:    suite.oracles[oracleIndex].String(),
				BridgerAddress:   suite.bridgers[oracleIndex].String(),
				ExternalAddress:  crypto.PubkeyToAddress(suite.externals[oracleIndex].PublicKey).Hex(),
				ValidatorAddress: suite.validator[oracleIndex].String(),
				DelegateAmount: sdk.Coin{
					Denom:  fxtypes.DefaultDenom,
					Amount: sdk.NewInt((rand.Int63n(3) + 1) * 10_000).MulRaw(1e18),
				},
				ChainName: suite.chainName,
			}

			testCase.preRun(msg)

			_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msg)
			if !testCase.pass {
				require.Error(t, err)
				require.EqualValues(suite.T(), testCase.err, err.Error())
				return
			}

			// success check
			require.NoError(t, err)

			// check oracle
			oracle, found := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			require.True(t, found)
			require.NotNil(t, oracle)
			require.EqualValues(t, msg.OracleAddress, oracle.OracleAddress)
			require.EqualValues(t, msg.BridgerAddress, oracle.BridgerAddress)
			require.EqualValues(t, msg.ExternalAddress, oracle.ExternalAddress)
			require.True(t, oracle.Online)
			require.EqualValues(t, msg.ValidatorAddress, oracle.DelegateValidator)
			require.EqualValues(t, int64(0), oracle.SlashTimes)

			// check relationship
			oracleAddr, found := suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, sdk.MustAccAddressFromBech32(msg.BridgerAddress))
			suite.True(found)
			suite.EqualValues(msg.OracleAddress, oracleAddr.String())

			oracleAddr, found = suite.Keeper().GetOracleByExternalAddress(suite.ctx, msg.ExternalAddress)
			suite.True(found)
			suite.EqualValues(msg.OracleAddress, oracleAddr.String())

			// check power
			totalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
			suite.EqualValues(msg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())

			// check delegate
			oracleDelegateAddr := oracle.GetDelegateAddress(suite.chainName)
			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, oracleDelegateAddr, suite.validator[oracleIndex])
			suite.True(found)
			suite.EqualValues(oracleDelegateAddr.String(), delegation.DelegatorAddress)
			suite.EqualValues(msg.ValidatorAddress, delegation.ValidatorAddress)
			suite.Truef(msg.DelegateAmount.Amount.Equal(delegation.GetShares().TruncateInt()), "expect:%s,actual:%s", msg.DelegateAmount.Amount.String(), delegation.GetShares().TruncateInt().String())
		})
	}
}

func (suite *KeeperTestSuite) TestMsgAddDelegate() {
	initDelegateAmount := suite.Keeper().GetOracleDelegateThreshold(suite.ctx).Amount
	testCases := []struct {
		name                 string
		pass                 bool
		err                  string
		preRun               func(msg *types.MsgAddDelegate)
		expectDelegateAmount func(msg *types.MsgAddDelegate) sdk.Int
	}{
		{
			name: "error - sender not oracle",
			preRun: func(msg *types.MsgAddDelegate) {
				msg.OracleAddress = sdk.AccAddress(tmrand.Bytes(20)).String()
			},
			pass: false,
			err:  types.ErrNoFoundOracle.Error(),
		},
		{
			name: "error - stake denom not match chain params stake denom",
			preRun: func(msg *types.MsgAddDelegate) {
				msg.Amount.Denom = "stake"
			},
			pass: false,
			err:  fmt.Sprintf("delegate denom got %s, expected %s: invalid", "stake", "FX"),
		},
		{
			name: "error - not sufficient slash amount",
			preRun: func(msg *types.MsgAddDelegate) {
				oracle, _ := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
				oracle.SlashTimes = 1
				suite.Keeper().SetOracle(suite.ctx, oracle)
				slashFraction := suite.Keeper().GetSlashFraction(suite.ctx)
				slashAmount := initDelegateAmount.ToDec().Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				randomAmount := rand.Int63n(slashAmount.QuoRaw(1e18).Int64()) + 1
				msg.Amount.Amount = sdk.NewInt(randomAmount).MulRaw(1e18).Sub(sdk.NewInt(1))
			},
			pass: false,
			err:  "not sufficient slash amount: invalid",
		},
		{
			name: "error - delegate amount less than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				params := suite.Keeper().GetParams(suite.ctx)
				addDelegateThreshold := rand.Int63n(100000) + 1
				params.DelegateThreshold.Amount = initDelegateAmount.Add(sdk.NewInt(addDelegateThreshold).MulRaw(1e18))
				suite.Keeper().SetParams(suite.ctx, &params)
				msg.Amount.Amount = sdk.NewInt(rand.Int63n(addDelegateThreshold) + 1).MulRaw(1e18).Sub(sdk.NewInt(1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
				delegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.ctx)
				maxDelegateAmount := delegateThreshold.Amount.Mul(sdk.NewInt(delegateMultiple))
				msg.Amount.Amount = maxDelegateAmount.Add(sdk.NewInt(rand.Int63() - 1))
			},
			pass: false,
			err:  types.ErrDelegateAmountAboveMaximum.Error(),
		},
		{
			name: "pass",
			preRun: func(msg *types.MsgAddDelegate) {
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdk.Int {
				return initDelegateAmount.Add(msg.Amount.Amount)
			},
		},
		{
			name: "pass - add slash amount",
			preRun: func(msg *types.MsgAddDelegate) {
				oracle, _ := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
				oracle.SlashTimes = 1
				oracle.Online = false
				suite.Keeper().SetOracle(suite.ctx, oracle)

				slashFraction := suite.Keeper().GetSlashFraction(suite.ctx)
				slashAmount := initDelegateAmount.ToDec().Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				msg.Amount.Amount = slashAmount
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdk.Int {
				return initDelegateAmount
			},
		},
		{
			name: "pass - add more slash amount",
			preRun: func(msg *types.MsgAddDelegate) {
				oracle, _ := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
				oracle.SlashTimes = 1
				oracle.Online = false
				suite.Keeper().SetOracle(suite.ctx, oracle)

				slashFraction := suite.Keeper().GetSlashFraction(suite.ctx)
				slashAmount := initDelegateAmount.ToDec().Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				msg.Amount.Amount = slashAmount.Add(sdk.NewInt(1000).MulRaw(1e18))
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdk.Int {
				return initDelegateAmount.Add(sdk.NewInt(1000).MulRaw(1e18))
			},
		},
	}
	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			suite.SetupTest()
			oracleIndex := rand.Intn(len(suite.oracles))

			// init bonded oracle
			_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
				OracleAddress:    suite.oracles[oracleIndex].String(),
				BridgerAddress:   suite.bridgers[oracleIndex].String(),
				ExternalAddress:  crypto.PubkeyToAddress(suite.externals[oracleIndex].PublicKey).Hex(),
				ValidatorAddress: suite.validator[oracleIndex].String(),
				DelegateAmount: sdk.Coin{
					Denom:  fxtypes.DefaultDenom,
					Amount: initDelegateAmount,
				},
				ChainName: suite.chainName,
			})
			require.NoError(t, err)

			oracleDelegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
			oracleDelegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.ctx)
			maxDelegateAmount := oracleDelegateThreshold.Amount.Mul(sdk.NewInt(oracleDelegateMultiple)).Sub(initDelegateAmount)
			msg := &types.MsgAddDelegate{
				ChainName:     suite.chainName,
				OracleAddress: suite.oracles[oracleIndex].String(),
				Amount: sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(
					rand.Int63n(maxDelegateAmount.QuoRaw(1e18).Int64())+1,
				).
					MulRaw(1e18).
					Sub(sdk.NewInt(1))),
			}
			testCase.preRun(msg)

			_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), msg)
			if !testCase.pass {
				require.Error(t, err)
				require.EqualValues(suite.T(), testCase.err, err.Error())
				return
			}

			// success check
			require.NoError(t, err)

			// check oracle
			oracle, found := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			require.True(t, found)
			require.NotNil(t, oracle)
			require.EqualValues(t, msg.OracleAddress, oracle.OracleAddress)
			require.True(t, oracle.Online)
			require.EqualValues(t, 0, oracle.SlashTimes)

			// check power
			totalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
			expectDelegateAmount := testCase.expectDelegateAmount(msg)
			suite.EqualValues(expectDelegateAmount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())

			// check delegate
			oracleDelegateAddr := oracle.GetDelegateAddress(suite.chainName)
			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, oracleDelegateAddr, suite.validator[oracleIndex])
			suite.True(found)
			suite.EqualValues(oracleDelegateAddr.String(), delegation.DelegatorAddress)
			suite.EqualValues(oracle.DelegateValidator, delegation.ValidatorAddress)
			suite.Truef(expectDelegateAmount.Equal(delegation.GetShares().TruncateInt()), "expect:%s,actual:%s", expectDelegateAmount.String(), delegation.GetShares().TruncateInt().String())
		})
	}
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
