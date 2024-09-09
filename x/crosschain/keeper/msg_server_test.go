package keeper_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
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
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.ctx, sdk.MustAccAddressFromBech32(msg.BridgerAddress), sdk.MustAccAddressFromBech32(msg.OracleAddress))
			},
			pass: false,
			err:  "bridger address is bound to oracle: invalid",
		},
		{
			name: "error - external address is bound",
			preRun: func(msg *types.MsgBondedOracle) {
				suite.Keeper().SetOracleAddrByExternalAddr(suite.ctx, msg.ExternalAddress, sdk.MustAccAddressFromBech32(msg.OracleAddress))
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
				msg.DelegateAmount.Amount = delegateThreshold.Amount.Sub(sdkmath.NewInt(tmrand.Int63() - 1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgBondedOracle) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
				delegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.ctx)
				maxDelegateAmount := delegateThreshold.Amount.Mul(sdkmath.NewInt(delegateMultiple))
				msg.DelegateAmount.Amount = maxDelegateAmount.Add(sdkmath.NewInt(tmrand.Int63() - 1))
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
		suite.Run(testCase.name, func() {
			oracleIndex := tmrand.Intn(len(suite.oracleAddrs))
			msg := &types.MsgBondedOracle{
				OracleAddress:    suite.oracleAddrs[oracleIndex].String(),
				BridgerAddress:   suite.bridgerAddrs[oracleIndex].String(),
				ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[oracleIndex].PublicKey),
				ValidatorAddress: suite.valAddrs[oracleIndex].String(),
				DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(3) + 1) * 10_000).MulRaw(1e18)),
				ChainName:        suite.chainName,
			}

			testCase.preRun(msg)

			_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msg)
			if !testCase.pass {
				suite.Require().Error(err)
				suite.Require().EqualValues(testCase.err, err.Error())
				return
			}

			// success check
			suite.Require().NoError(err)

			// check oracle
			oracle, found := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			suite.Require().True(found)
			suite.Require().NotNil(oracle)
			suite.Require().EqualValues(msg.OracleAddress, oracle.OracleAddress)
			suite.Require().EqualValues(msg.BridgerAddress, oracle.BridgerAddress)
			suite.Require().EqualValues(msg.ExternalAddress, oracle.ExternalAddress)
			suite.Require().True(oracle.Online)
			suite.Require().EqualValues(msg.ValidatorAddress, oracle.DelegateValidator)
			suite.Require().EqualValues(int64(0), oracle.SlashTimes)

			// check relationship
			oracleAddr, found := suite.Keeper().GetOracleAddrByBridgerAddr(suite.ctx, sdk.MustAccAddressFromBech32(msg.BridgerAddress))
			suite.True(found)
			suite.Require().EqualValues(msg.OracleAddress, oracleAddr.String())

			oracleAddr, found = suite.Keeper().GetOracleAddrByExternalAddr(suite.ctx, msg.ExternalAddress)
			suite.True(found)
			suite.Require().EqualValues(msg.OracleAddress, oracleAddr.String())

			// check power
			totalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
			suite.Require().EqualValues(msg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())

			// check delegate
			oracleDelegateAddr := oracle.GetDelegateAddress(suite.chainName)
			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, oracleDelegateAddr, suite.valAddrs[oracleIndex])
			suite.True(found)
			suite.Require().EqualValues(oracleDelegateAddr.String(), delegation.DelegatorAddress)
			suite.Require().EqualValues(msg.ValidatorAddress, delegation.ValidatorAddress)
			suite.Truef(msg.DelegateAmount.Amount.Equal(delegation.GetShares().TruncateInt()), "expect:%s,actual:%s", msg.DelegateAmount.Amount.String(), delegation.GetShares().TruncateInt().String())

			startingInfo := suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, suite.valAddrs[oracleIndex], oracleDelegateAddr)
			suite.NotNil(startingInfo)
			suite.EqualValues(uint64(suite.ctx.BlockHeight()), startingInfo.Height)
			suite.True(startingInfo.PreviousPeriod > 0)
			suite.EqualValues(sdk.NewDecFromInt(msg.DelegateAmount.Amount).String(), startingInfo.Stake.String())
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
		expectDelegateAmount func(msg *types.MsgAddDelegate) sdkmath.Int
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
				slashAmount := sdk.NewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				randomAmount := tmrand.Int63n(slashAmount.QuoRaw(1e18).Int64()) + 1
				msg.Amount.Amount = sdkmath.NewInt(randomAmount).MulRaw(1e18).Sub(sdkmath.NewInt(1))
			},
			pass: false,
			err:  "not sufficient slash amount: invalid",
		},
		{
			name: "error - delegate amount less than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				params := suite.Keeper().GetParams(suite.ctx)
				addDelegateThreshold := tmrand.Int63n(100000) + 1
				params.DelegateThreshold.Amount = initDelegateAmount.Add(sdkmath.NewInt(addDelegateThreshold).MulRaw(1e18))
				err := suite.Keeper().SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				msg.Amount.Amount = sdkmath.NewInt(tmrand.Int63n(addDelegateThreshold) + 1).MulRaw(1e18).Sub(sdkmath.NewInt(1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
				delegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.ctx)
				maxDelegateAmount := delegateThreshold.Amount.Mul(sdkmath.NewInt(delegateMultiple))
				msg.Amount.Amount = maxDelegateAmount.Add(sdkmath.NewInt(tmrand.Int63() - 1))
			},
			pass: false,
			err:  types.ErrDelegateAmountAboveMaximum.Error(),
		},
		{
			name: "pass",
			preRun: func(msg *types.MsgAddDelegate) {
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdkmath.Int {
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
				slashAmount := sdk.NewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				msg.Amount.Amount = slashAmount
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdkmath.Int {
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
				slashAmount := sdk.NewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				msg.Amount.Amount = slashAmount.Add(sdkmath.NewInt(1000).MulRaw(1e18))
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdkmath.Int {
				return initDelegateAmount.Add(sdkmath.NewInt(1000).MulRaw(1e18))
			},
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			oracleIndex := tmrand.Intn(len(suite.oracleAddrs))

			// init bonded oracle
			_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
				OracleAddress:    suite.oracleAddrs[oracleIndex].String(),
				BridgerAddress:   suite.bridgerAddrs[oracleIndex].String(),
				ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[oracleIndex].PublicKey),
				ValidatorAddress: suite.valAddrs[oracleIndex].String(),
				DelegateAmount:   types.NewDelegateAmount(initDelegateAmount),
				ChainName:        suite.chainName,
			})
			suite.Require().NoError(err)

			oracleDelegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.ctx)
			oracleDelegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.ctx)
			maxDelegateAmount := oracleDelegateThreshold.Amount.Mul(sdkmath.NewInt(oracleDelegateMultiple)).Sub(initDelegateAmount)
			msg := &types.MsgAddDelegate{
				ChainName:     suite.chainName,
				OracleAddress: suite.oracleAddrs[oracleIndex].String(),
				Amount: types.NewDelegateAmount(sdkmath.NewInt(
					tmrand.Int63n(maxDelegateAmount.QuoRaw(1e18).Int64()) + 1,
				).
					MulRaw(1e18).
					Sub(sdkmath.NewInt(1))),
			}
			testCase.preRun(msg)

			_, err = suite.MsgServer().AddDelegate(sdk.WrapSDKContext(suite.ctx), msg)
			if !testCase.pass {
				suite.Require().Error(err)
				suite.Require().EqualValues(testCase.err, err.Error())
				return
			}

			// success check
			suite.Require().NoError(err)

			// check oracle
			oracle, found := suite.Keeper().GetOracle(suite.ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			suite.Require().True(found)
			suite.Require().NotNil(oracle)
			suite.Require().EqualValues(msg.OracleAddress, oracle.OracleAddress)
			suite.Require().True(oracle.Online)
			suite.Require().EqualValues(0, oracle.SlashTimes)

			// check power
			totalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
			expectDelegateAmount := testCase.expectDelegateAmount(msg)
			suite.Require().EqualValues(expectDelegateAmount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())

			// check delegate
			oracleDelegateAddr := oracle.GetDelegateAddress(suite.chainName)
			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, oracleDelegateAddr, suite.valAddrs[oracleIndex])
			suite.True(found)
			suite.Require().EqualValues(oracleDelegateAddr.String(), delegation.DelegatorAddress)
			suite.Require().EqualValues(oracle.DelegateValidator, delegation.ValidatorAddress)
			suite.Truef(expectDelegateAmount.Equal(delegation.GetShares().TruncateInt()), "expect:%s,actual:%s", expectDelegateAmount.String(), delegation.GetShares().TruncateInt().String())
			startingInfo := suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, suite.valAddrs[oracleIndex], oracleDelegateAddr)
			suite.NotNil(startingInfo)
			suite.EqualValues(uint64(suite.ctx.BlockHeight()), startingInfo.Height)
			suite.True(startingInfo.PreviousPeriod > 0)
			suite.EqualValues(sdk.NewDecFromInt(expectDelegateAmount).String(), startingInfo.Stake.String())
		})
	}
}

func (suite *KeeperTestSuite) TestMsgEditBridger() {
	delegateAmount := sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)
	for i := range suite.oracleAddrs {
		bondedMsg := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(delegateAmount),
			ChainName:        suite.chainName,
		}
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), bondedMsg)
		suite.NoError(err)
	}

	token := fmt.Sprintf("0x%s", tmrand.Str(40))
	denom := types.NewBridgeDenom(suite.chainName, token)
	suite.Keeper().AddBridgeToken(suite.ctx, token, denom)

	privateKey, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	sendToMsg := &types.MsgSendToFxClaim{
		EventNonce:    1,
		BlockHeight:   100,
		TokenContract: token,
		Amount:        sdkmath.NewInt(int64(tmrand.Uint32())),
		Sender:        suite.PubKeyToExternalAddr(privateKey.PublicKey),
		Receiver:      sdk.AccAddress(tmrand.Bytes(20)).String(),
		TargetIbc:     "",
		ChainName:     suite.chainName,
	}
	for i := 0; i < len(suite.bridgerAddrs)/2; i++ {
		sendToMsg.BridgerAddress = suite.bridgerAddrs[i].String()
		err = suite.SendClaimReturnErr(sendToMsg)
		suite.NoError(err)
	}

	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{ChainID: suite.ctx.ChainID(), Height: suite.ctx.BlockHeight()}})

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToMsg.Receiver))
	suite.Require().Equal(balances.String(), sdk.NewCoins().String(), len(suite.bridgerAddrs))

	for i := 0; i < len(suite.oracleAddrs); i++ {
		_, err := suite.MsgServer().EditBridger(sdk.WrapSDKContext(suite.ctx), &types.MsgEditBridger{
			ChainName:      suite.chainName,
			OracleAddress:  suite.oracleAddrs[i].String(),
			BridgerAddress: suite.bridgerAddrs[i].String(),
		})
		suite.Require().Error(err)

		_, err = suite.MsgServer().EditBridger(sdk.WrapSDKContext(suite.ctx), &types.MsgEditBridger{
			ChainName:      suite.chainName,
			OracleAddress:  suite.oracleAddrs[i].String(),
			BridgerAddress: sdk.AccAddress(suite.valAddrs[i]).String(),
		})
		suite.NoError(err)

		sendToMsg.BridgerAddress = sdk.AccAddress(suite.valAddrs[i]).String()
		err = suite.SendClaimReturnErr(sendToMsg)
		if i < len(suite.oracleAddrs)/2 {
			suite.Require().ErrorContains(err, types.ErrNonContiguousEventNonce.Error())
		} else {
			suite.Require().NoError(err)
		}
	}
	err = suite.Keeper().ExecuteClaim(suite.ctx, sendToMsg.EventNonce)
	suite.Require().NoError(err)

	for _, bridger := range suite.bridgerAddrs {
		suite.False(suite.Keeper().HasOracleAddrByBridgerAddr(suite.ctx, bridger))
	}

	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToMsg.Receiver))
	suite.Require().Equal(balances.String(), sdk.NewCoins(sdk.NewCoin(denom, sendToMsg.Amount)).String())
}

func (suite *KeeperTestSuite) TestMsgSetOracleSetConfirm() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	suite.Require().NoError(err)

	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	suite.Require().EqualValues(0, latestOracleSetNonce)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	latestOracleSetNonce = suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	suite.Require().EqualValues(1, latestOracleSetNonce)

	suite.Require().True(suite.Keeper().HasOracleSetRequest(suite.ctx, 1))

	suite.Require().False(suite.Keeper().HasOracleSetRequest(suite.ctx, 2))

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.ctx, 1)
	suite.Require().EqualValues(uint64(1), nonce1OracleSet.Nonce)
	suite.Require().EqualValues(uint64(2), nonce1OracleSet.Height)
	suite.Require().EqualValues(1, len(nonce1OracleSet.Members))
	suite.Require().EqualValues(normalMsg.ExternalAddress, nonce1OracleSet.Members[0].ExternalAddress)
	suite.Require().EqualValues(math.MaxUint32, nonce1OracleSet.Members[0].Power)

	gravityId := suite.Keeper().GetGravityID(suite.ctx)
	checkpoint, err := nonce1OracleSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == suite.chainName {
		checkpoint, err = trontypes.GetCheckpointOracleSet(nonce1OracleSet, gravityId)
	}
	suite.Require().NoError(err)

	external1Signature, err := types.NewEthereumSignature(checkpoint, suite.externalPris[0])
	if trontypes.ModuleName == suite.chainName {
		external1Signature, err = trontypes.NewTronSignature(checkpoint, suite.externalPris[0])
	}
	suite.Require().NoError(err)
	external2Signature, err := types.NewEthereumSignature(checkpoint, suite.externalPris[1])
	if trontypes.ModuleName == suite.chainName {
		external2Signature, err = trontypes.NewTronSignature(checkpoint, suite.externalPris[1])
	}
	suite.Require().NoError(err)

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
				BridgerAddress:  suite.bridgerAddrs[0].String(),
				ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
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
				BridgerAddress:  suite.bridgerAddrs[0].String(),
				ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[1].PublicKey),
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
				BridgerAddress:  suite.bridgerAddrs[0].String(),
				ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
				Signature:       hex.EncodeToString(external2Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s: %s", normalMsg.ExternalAddress, hex.EncodeToString(checkpoint), hex.EncodeToString(external2Signature), types.ErrInvalid),
		},
		{
			name: "bridger address not match",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgerAddrs[1].String(),
				ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("got %s, expected %s: %s", suite.bridgerAddrs[1].String(), suite.bridgerAddrs[0].String(), types.ErrInvalid),
		},
	}

	for _, testData := range errMsgData {
		_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), testData.msg)
		suite.Require().ErrorIs(err, testData.err, testData.name)
		suite.Require().EqualValues(err.Error(), testData.errReason, testData.name)
	}

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:           nonce1OracleSet.Nonce,
		BridgerAddress:  suite.bridgerAddrs[0].String(),
		ExternalAddress: normalMsg.ExternalAddress,
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       suite.chainName,
	}
	_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), normalOracleSetConfirmMsg)
	suite.Require().NoError(err)

	endBlockBeforeLatestOracleSet := suite.Keeper().GetLatestOracleSet(suite.ctx)
	suite.Require().NotNil(endBlockBeforeLatestOracleSet)
}

func (suite *KeeperTestSuite) TestClaimWithOracleOnline() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	suite.Require().NoError(err)

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	suite.Require().EqualValues(1, latestOracleSetNonce)

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.ctx, latestOracleSetNonce)
	suite.Require().EqualValues(uint64(1), nonce1OracleSet.Nonce)
	suite.Require().EqualValues(uint64(2), nonce1OracleSet.Height)

	var gravityId string
	suite.Require().NotPanics(func() {
		gravityId = suite.Keeper().GetGravityID(suite.ctx)
	})
	if suite.chainName == ethtypes.ModuleName {
		suite.Require().EqualValues(fmt.Sprintf("fx-bridge-%s", suite.chainName), gravityId)
	} else {
		suite.Require().EqualValues(fmt.Sprintf("fx-%s-bridge", suite.chainName), gravityId)
	}
	checkpoint, err := nonce1OracleSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == suite.chainName {
		checkpoint, err = trontypes.GetCheckpointOracleSet(nonce1OracleSet, gravityId)
	}
	suite.Require().NoError(err)

	// oracle Online!!!
	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[0])
	suite.Require().True(found)
	oracle.Online = true
	suite.Keeper().SetOracle(suite.ctx, oracle)

	external1Signature, err := types.NewEthereumSignature(checkpoint, suite.externalPris[0])
	if trontypes.ModuleName == suite.chainName {
		external1Signature, err = trontypes.NewTronSignature(checkpoint, suite.externalPris[0])
	}
	suite.Require().NoError(err)

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:           latestOracleSetNonce,
		BridgerAddress:  suite.bridgerAddrs[0].String(),
		ExternalAddress: normalMsg.ExternalAddress,
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       suite.chainName,
	}
	_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), normalOracleSetConfirmMsg)
	suite.Require().Nil(err)
}

func (suite *KeeperTestSuite) TestClaimMsgGasConsumed() {
	gasStatics := func(gasConsumed, maxGas uint64, minGas uint64, avgGas uint64) (uint64, uint64, uint64) {
		if gasConsumed > maxGas {
			maxGas = gasConsumed
		}
		if minGas == 0 || gasConsumed < minGas {
			minGas = gasConsumed
		}
		if avgGas == 0 {
			avgGas = gasConsumed
		} else {
			avgGas = (avgGas + gasConsumed) / 2
		}
		return maxGas, minGas, avgGas
	}

	testCases := []struct {
		name     string
		buildMsg func() types.ExternalClaim
		execute  func(claim types.ExternalClaim) (minGas, maxGas, avgGas uint64)
	}{
		{
			name: "MsgSendToFx",
			buildMsg: func() types.ExternalClaim {
				return &types.MsgBridgeTokenClaim{
					BlockHeight:   tmrand.Uint64(),
					TokenContract: helpers.GenHexAddress().String(),
					Name:          tmrand.Str(10),
					Symbol:        tmrand.Str(10),
					Decimals:      uint64(tmrand.Int63n(18) + 1),
					ChannelIbc:    "",
					ChainName:     suite.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgBridgeTokenClaim)
				suite.True(ok)
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					err := suite.SendClaimReturnErr(msg)
					suite.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
		{
			name: "MsgSendToFxClaim",
			buildMsg: func() types.ExternalClaim {
				return &types.MsgSendToFxClaim{
					BlockHeight:   tmrand.Uint64(),
					TokenContract: helpers.GenHexAddress().String(),
					Amount:        sdkmath.NewInt(tmrand.Int63n(100000) + 1).MulRaw(1e18),
					Sender:        helpers.GenExternalAddr(suite.chainName),
					Receiver:      sdk.AccAddress(tmrand.Bytes(20)).String(),
					TargetIbc:     "",
					ChainName:     suite.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgSendToFxClaim)
				suite.True(ok)
				suite.Keeper().AddBridgeToken(suite.ctx, msg.TokenContract, types.NewBridgeDenom(suite.chainName, msg.TokenContract))
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					err := suite.SendClaimReturnErr(msg)
					suite.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
		{
			name: "OracleSetUpdateClaim",
			buildMsg: func() types.ExternalClaim {
				var externalOracleMembers types.BridgeValidators
				for _, key := range suite.externalPris {
					bridgeVal := types.BridgeValidator{
						Power:           tmrand.Uint64(),
						ExternalAddress: suite.PubKeyToExternalAddr(key.PublicKey),
					}
					externalOracleMembers = append(externalOracleMembers, bridgeVal)
				}
				return &types.MsgOracleSetUpdatedClaim{
					BlockHeight:    tmrand.Uint64(),
					OracleSetNonce: tmrand.Uint64(),
					Members:        externalOracleMembers,
					ChainName:      suite.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgOracleSetUpdatedClaim)
				suite.True(ok)
				suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
					Nonce:   msg.OracleSetNonce,
					Height:  msg.BlockHeight,
					Members: msg.Members,
				})
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					err := suite.SendClaimReturnErr(msg)
					suite.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
		{
			name: "MsgSendToExternalClaim",
			buildMsg: func() types.ExternalClaim {
				return &types.MsgSendToExternalClaim{
					BlockHeight:   tmrand.Uint64(),
					BatchNonce:    tmrand.Uint64(),
					TokenContract: helpers.GenHexAddress().String(),
					ChainName:     suite.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgSendToExternalClaim)
				suite.True(ok)
				suite.Require().NoError(suite.Keeper().StoreBatch(suite.ctx, &types.OutgoingTxBatch{
					BatchNonce:    msg.BatchNonce,
					TokenContract: msg.TokenContract,
				}))
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					err := suite.SendClaimReturnErr(msg)
					suite.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("%s-%s", suite.chainName, testCase.name), func() {
			for i, oracle := range suite.oracleAddrs {
				msg := &types.MsgBondedOracle{
					OracleAddress:    oracle.String(),
					BridgerAddress:   suite.bridgerAddrs[i].String(),
					ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
					ValidatorAddress: suite.valAddrs[0].String(),
					DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
					ChainName:        suite.chainName,
				}
				_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msg)
				suite.Require().NoError(err)
			}

			claimMsg := testCase.buildMsg()
			minGas, maxGas, avgGas := testCase.execute(claimMsg)
			suite.Require().EqualValuesf(minGas, maxGas, "expect equal min:%d, max:%d, diff:%d", minGas, maxGas, maxGas-minGas)
			suite.Require().EqualValuesf(minGas, maxGas, "expect equal min:%d, avg:%d, diff:%d", minGas, avgGas, avgGas-minGas)
		})
	}
}

func (suite *KeeperTestSuite) TestClaimTest() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	suite.Require().NoError(err)

	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0])
	suite.Require().EqualValues(0, oracleLastEventNonce)

	randomPrivateKey, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	testMsgs := []struct {
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
				TokenContract:  suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgerAddrs[0].String(),
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
				TokenContract:  suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgerAddrs[0].String(),
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
				TokenContract:  suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgerAddrs[0].String(),
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
				TokenContract:  suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgerAddrs[0].String(),
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
				TokenContract:  suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				BridgerAddress: suite.bridgerAddrs[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 3, 2, types.ErrNonContiguousEventNonce),
		},
	}

	for _, testData := range testMsgs {
		err = testData.msg.ValidateBasic()
		suite.Require().NoError(err)
		err = suite.SendClaimReturnErr(testData.msg)
		suite.Require().ErrorIs(err, testData.err, testData.name)
		if err == nil {
			continue
		}
		suite.Require().EqualValues(testData.errReason, err.Error(), testData.name)
	}
}

func (suite *KeeperTestSuite) TestRequestBatchBaseFee() {
	// 1. First sets up a valid validator
	totalPower := sdkmath.ZeroInt()
	delegateAmounts := make([]sdkmath.Int, 0, len(suite.oracleAddrs))
	for i, oracle := range suite.oracleAddrs {
		normalMsg := &types.MsgBondedOracle{
			OracleAddress:    oracle.String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[0].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		if len(suite.valAddrs) > i {
			normalMsg.ValidatorAddress = suite.valAddrs[i].String()
		}
		delegateAmounts = append(delegateAmounts, normalMsg.DelegateAmount.Amount)
		totalPower = totalPower.Add(normalMsg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction))
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	var externalOracleMembers types.BridgeValidators
	for i, key := range suite.externalPris {
		power := delegateAmounts[i].Quo(sdk.DefaultPowerReduction).MulRaw(math.MaxUint32).Quo(totalPower)
		bridgeVal := types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: suite.PubKeyToExternalAddr(key.PublicKey),
		}
		externalOracleMembers = append(externalOracleMembers, bridgeVal)
	}
	sort.Sort(externalOracleMembers)

	// 2. oracle update claim
	for i := range suite.oracleAddrs {
		normalMsg := &types.MsgOracleSetUpdatedClaim{
			EventNonce:     1,
			BlockHeight:    1,
			OracleSetNonce: 1,
			Members:        externalOracleMembers,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChainName:      suite.chainName,
		}
		err := suite.SendClaimReturnErr(normalMsg)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	// 3. add bridge token.
	sendToFxSendAddr := suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey)
	sendToFxReceiveAddr := suite.bridgerAddrs[0]
	sendToFxAmount := sdkmath.NewIntWithDecimal(1000, 18)
	randomPrivateKey, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	sendToFxToken := suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey)

	for i, oracle := range suite.oracleAddrs {
		normalMsg := &types.MsgBridgeTokenClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Name:           "Test USDT",
			Symbol:         "USDT",
			Decimals:       18,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChannelIbc:     "",
			ChainName:      suite.chainName,
		}
		err = suite.SendClaimReturnErr(normalMsg)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	bridgeDenomData := suite.Keeper().GetBridgeTokenDenom(suite.ctx, sendToFxToken)
	suite.Require().NotNil(bridgeDenomData)
	tokenDenom := types.NewBridgeDenom(suite.chainName, sendToFxToken)
	suite.Require().EqualValues(tokenDenom, bridgeDenomData.Denom)
	bridgeTokenData := suite.Keeper().GetDenomBridgeToken(suite.ctx, tokenDenom)
	suite.Require().NotNil(bridgeTokenData)
	suite.Require().EqualValues(sendToFxToken, bridgeTokenData.Token)

	// 4. sendToFx.
	sendToFxClaim := new(types.MsgSendToFxClaim)
	for i, oracle := range suite.oracleAddrs {
		sendToFxClaim = &types.MsgSendToFxClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Amount:         sendToFxAmount,
			Sender:         sendToFxSendAddr,
			Receiver:       sendToFxReceiveAddr.String(),
			TargetIbc:      "",
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChainName:      suite.chainName,
		}
		err = suite.SendClaimReturnErr(sendToFxClaim)
		suite.Require().NoError(err)
	}

	err = suite.Keeper().ExecuteClaim(suite.ctx, sendToFxClaim.EventNonce)
	suite.Require().NoError(err)

	balance := suite.app.BankKeeper.GetBalance(suite.ctx, sendToFxReceiveAddr, tokenDenom)
	suite.Require().NotNil(balance)
	suite.Require().EqualValues(balance.Denom, tokenDenom)
	suite.Require().True(balance.Amount.Equal(sendToFxAmount))

	sendToExternal := func(bridgeFees []sdkmath.Int) {
		for _, bridgeFee := range bridgeFees {
			sendToExternal := &types.MsgSendToExternal{
				Sender:    sendToFxReceiveAddr.String(),
				Dest:      sendToFxSendAddr,
				Amount:    sdk.NewCoin(tokenDenom, sdkmath.NewInt(3)),
				BridgeFee: sdk.NewCoin(tokenDenom, bridgeFee),
				ChainName: suite.chainName,
			}
			_, err := suite.MsgServer().SendToExternal(sdk.WrapSDKContext(suite.ctx), sendToExternal)
			suite.Require().NoError(err)
		}
	}

	sendToExternal([]sdkmath.Int{sdkmath.NewInt(1), sdkmath.NewInt(2), sdkmath.NewInt(3)})
	usdtBatchFee := suite.Keeper().GetBatchFeesByTokenType(suite.ctx, sendToFxToken, 100, sdkmath.NewInt(0))
	suite.Require().EqualValues(sendToFxToken, usdtBatchFee.TokenContract)
	suite.Require().EqualValues(3, usdtBatchFee.TotalTxs)
	suite.Require().EqualValues(sdkmath.NewInt(6), usdtBatchFee.TotalFees)

	testCases := []struct {
		testName       string
		baseFee        sdkmath.Int
		pass           bool
		expectTotalTxs uint64
		err            error
	}{
		{
			testName:       "Support - baseFee 1000",
			baseFee:        sdkmath.NewInt(1000),
			pass:           false,
			expectTotalTxs: 3,
			err:            errorsmod.Wrap(types.ErrEmpty, "no batch tx"),
		},
		{
			testName:       "Support - baseFee 2",
			baseFee:        sdkmath.NewInt(2),
			pass:           true,
			expectTotalTxs: 1,
			err:            nil,
		},
		{
			testName:       "Support - baseFee 0",
			baseFee:        sdkmath.NewInt(0),
			pass:           false,
			expectTotalTxs: 0,
			err:            errorsmod.Wrap(types.ErrInvalid, "new batch would not be more profitable"),
		},
	}

	for _, testCase := range testCases {
		_, err := suite.MsgServer().RequestBatch(sdk.WrapSDKContext(suite.ctx), &types.MsgRequestBatch{
			Sender:     suite.bridgerAddrs[0].String(),
			Denom:      tokenDenom,
			MinimumFee: sdkmath.NewInt(1),
			FeeReceive: "0x0000000000000000000000000000000000000000",
			ChainName:  suite.chainName,
			BaseFee:    testCase.baseFee,
		})
		if testCase.pass {
			suite.Require().NoError(err)
			usdtBatchFee = suite.Keeper().GetBatchFeesByTokenType(suite.ctx, sendToFxToken, 100, sdkmath.NewInt(0))
			suite.Require().EqualValues(testCase.expectTotalTxs, usdtBatchFee.TotalTxs)
		} else {
			suite.Require().NotNil(err)
			suite.Require().Equal(err.Error(), testCase.err.Error())
		}
	}
}

func (suite *KeeperTestSuite) TestMsgUpdateChainOracles() {
	updateOracle := &types.MsgUpdateChainOracles{
		Oracles:   []string{},
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ChainName: suite.chainName,
	}
	for _, oracle := range suite.oracleAddrs {
		updateOracle.Oracles = append(updateOracle.Oracles, oracle.String())
	}

	_, err := suite.MsgServer().UpdateChainOracles(suite.ctx, updateOracle)
	suite.Require().NoError(err)
	for _, oracle := range suite.oracleAddrs {
		suite.Require().True(suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	updateOracle.Oracles = []string{}
	number := tmrand.Intn(100)
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	_, err = suite.MsgServer().UpdateChainOracles(suite.ctx, updateOracle)
	suite.Require().NoError(err)

	updateOracle.Oracles = []string{}
	number = tmrand.Intn(2) + 101
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	_, err = suite.MsgServer().UpdateChainOracles(suite.ctx, updateOracle)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestBridgeCallClaim() {
	suite.bondedOracle()

	tokenContract := helpers.GenExternalAddr(suite.chainName)

	suite.addBridgeToken(tokenContract, fxtypes.GetCrossChainMetadataManyToOne("test token", "TT", 18))

	suite.registerCoin(types.NewBridgeDenom(suite.chainName, tokenContract))

	fxTokenContract := helpers.GenExternalAddr(suite.chainName)
	suite.addBridgeToken(fxTokenContract, fxtypes.GetFXMetaData())

	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0])

	testMsgs := []struct {
		name      string
		msg       *types.MsgBridgeCallClaim
		err       error
		errReason string
		expect    bool
	}{
		{
			name: "success",
			msg: &types.MsgBridgeCallClaim{
				EventNonce:     oracleLastEventNonce + 1,
				Sender:         helpers.GenExternalAddr(suite.chainName),
				TokenContracts: []string{tokenContract},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(100)},
				Refund:         helpers.GenExternalAddr(suite.chainName),
				To:             helpers.GenExternalAddr(suite.chainName),
				Data:           "",
				Value:          sdkmath.NewInt(0),
				BlockHeight:    1,
				BridgerAddress: suite.bridgerAddrs[0].String(),
				ChainName:      suite.chainName,
				TxOrigin:       helpers.GenExternalAddr(suite.chainName),
			},
			err:       nil,
			errReason: "",
			expect:    true,
		},
	}

	for _, testData := range testMsgs {
		err := testData.msg.ValidateBasic()
		suite.Require().NoError(err)
		suite.ctx = suite.ctx.WithEventManager(sdk.NewEventManager())
		suite.Require().NoError(testData.msg.ValidateBasic())
		err = suite.SendClaimReturnErr(testData.msg)
		suite.Require().ErrorIs(err, testData.err, testData.name)
		if testData.err == nil {
			suite.checkObservationState(suite.ctx, testData.expect)
		}
		if err == nil {
			continue
		}

		suite.Require().EqualValues(testData.errReason, err.Error(), testData.name)
	}
}

func (suite *KeeperTestSuite) TestMsgBridgeCall() {
	// 1. First sets up a valid validator
	totalPower := sdkmath.ZeroInt()
	delegateAmounts := make([]sdkmath.Int, 0, len(suite.oracleAddrs))
	for i, oracle := range suite.oracleAddrs {
		normalMsg := &types.MsgBondedOracle{
			OracleAddress:    oracle.String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[0].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		if len(suite.valAddrs) > i {
			normalMsg.ValidatorAddress = suite.valAddrs[i].String()
		}
		delegateAmounts = append(delegateAmounts, normalMsg.DelegateAmount.Amount)
		totalPower = totalPower.Add(normalMsg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction))
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	var externalOracleMembers types.BridgeValidators
	for i, key := range suite.externalPris {
		power := delegateAmounts[i].Quo(sdk.DefaultPowerReduction).MulRaw(math.MaxUint32).Quo(totalPower)
		bridgeVal := types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: suite.PubKeyToExternalAddr(key.PublicKey),
		}
		externalOracleMembers = append(externalOracleMembers, bridgeVal)
	}
	sort.Sort(externalOracleMembers)

	// 2. oracle update claim
	for i := range suite.oracleAddrs {
		normalMsg := &types.MsgOracleSetUpdatedClaim{
			EventNonce:     1,
			BlockHeight:    1,
			OracleSetNonce: 1,
			Members:        externalOracleMembers,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChainName:      suite.chainName,
		}
		err := suite.SendClaimReturnErr(normalMsg)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	// 3. add bridge token.
	sendToFxSendAddr := suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey)
	sendToFxReceiveAddr := suite.bridgerAddrs[0]
	sendToFxAmount := sdkmath.NewIntWithDecimal(1000, 18)
	randomPrivateKey, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	sendToFxToken := suite.PubKeyToExternalAddr(randomPrivateKey.PublicKey)

	for i, oracle := range suite.oracleAddrs {
		normalMsg := &types.MsgBridgeTokenClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Name:           "Test USDT",
			Symbol:         "USDT",
			Decimals:       18,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChannelIbc:     "",
			ChainName:      suite.chainName,
		}
		err = suite.SendClaimReturnErr(normalMsg)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.ctx)

	bridgeDenomData := suite.Keeper().GetBridgeTokenDenom(suite.ctx, sendToFxToken)
	suite.Require().NotNil(bridgeDenomData)
	tokenDenom := types.NewBridgeDenom(suite.chainName, sendToFxToken)
	suite.Require().EqualValues(tokenDenom, bridgeDenomData.Denom)
	bridgeTokenData := suite.Keeper().GetDenomBridgeToken(suite.ctx, tokenDenom)
	suite.Require().NotNil(bridgeTokenData)
	suite.Require().EqualValues(sendToFxToken, bridgeTokenData.Token)

	// 4. register coin
	tokenPair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, banktypes.Metadata{
		Description: "FunctionX Token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "test",
				Exponent: 0,
				Aliases:  []string{types.NewBridgeDenom(suite.chainName, sendToFxToken)},
			}, {
				Denom:    "TEST",
				Exponent: 18,
			},
		},
		Base:    "test",
		Display: "TEST",
		Name:    "Test Token",
		Symbol:  "TEST",
	})
	suite.NoError(err)

	// 5. sendToFx.
	sendToFxClaim := new(types.MsgSendToFxClaim)
	for i, oracle := range suite.oracleAddrs {
		sendToFxClaim = &types.MsgSendToFxClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Amount:         sendToFxAmount,
			Sender:         sendToFxSendAddr,
			Receiver:       sendToFxReceiveAddr.String(),
			TargetIbc:      "",
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChainName:      suite.chainName,
		}
		err = suite.SendClaimReturnErr(sendToFxClaim)
		suite.Require().NoError(err)
	}
	err = suite.Keeper().ExecuteClaim(suite.ctx, sendToFxClaim.EventNonce)
	suite.Require().NoError(err)

	suite.Equal(sendToFxAmount, suite.app.BankKeeper.GetBalance(suite.ctx, sendToFxReceiveAddr, tokenPair.GetDenom()).Amount)

	testCases := []struct {
		name     string
		malleate func() *types.MsgBridgeCall
		pass     bool
		err      error
	}{
		{
			name: "pass",
			malleate: func() *types.MsgBridgeCall {
				return &types.MsgBridgeCall{
					ChainName: suite.chainName,
					Sender:    sendToFxReceiveAddr.String(),
					Refund:    helpers.GenAccAddress().String(),
					To:        helpers.GenExternalAddr(suite.chainName),
					Coins:     sdk.NewCoins(sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewInt(1e18))),
					Data:      "",
					Value:     sdkmath.ZeroInt(),
				}
			},
			pass: true,
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			msg := testCase.malleate()

			_, err = suite.MsgServer().BridgeCall(sdk.WrapSDKContext(suite.ctx), msg)
			if testCase.pass {
				suite.Require().NoError(err)
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equal(err.Error(), testCase.err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAddPendingPoolRewards() {
	txId := tmrand.Uint64()
	initRewards := sdk.NewCoins()
	addRewards := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1)))
	tx := types.NewPendingOutgoingTx(txId, helpers.GenHexAddress().Bytes(), helpers.GenExternalAddr(suite.chainName),
		tmrand.Str(40), sdk.NewCoin("test", sdkmath.NewInt(100)), sdk.NewCoin("test", sdkmath.NewInt(100)),
		initRewards)
	suite.Keeper().SetPendingTx(suite.ctx, &tx)

	// mint add reward coins to sender.
	sender := helpers.GenAccAddress()
	suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, suite.chainName, addRewards))
	suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, suite.chainName, sender, addRewards))

	testCases := []struct {
		name         string
		malleate     func() *types.MsgAddPendingPoolRewards
		pass         bool
		err          error
		expectReward sdk.Coins
	}{
		{
			name: "pass",
			malleate: func() *types.MsgAddPendingPoolRewards {
				return &types.MsgAddPendingPoolRewards{
					ChainName: suite.chainName,
					Id:        txId,
					Sender:    sender.String(),
					Rewards:   sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1))),
				}
			},
			pass:         true,
			expectReward: initRewards.Add(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1)))...),
		},
		{
			name: "err - rewards not FX denom",
			malleate: func() *types.MsgAddPendingPoolRewards {
				return &types.MsgAddPendingPoolRewards{
					ChainName: suite.chainName,
					Id:        txId,
					Sender:    sender.String(),
					Rewards:   sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1))),
				}
			},
			pass: false,
			err:  errors.ErrInvalidRequest.Wrapf("unsupported denomination %s, only %s is supported", "test", fxtypes.DefaultDenom),
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			msg := testCase.malleate()

			_, err := suite.MsgServer().AddPendingPoolRewards(sdk.WrapSDKContext(suite.ctx), msg)
			if testCase.pass {
				suite.Require().NoError(err)
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equal(err.Error(), testCase.err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) bondedOracle() {
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	})
	suite.Require().NoError(err)

	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0])
	suite.Require().EqualValues(0, oracleLastEventNonce)
}

func (suite *KeeperTestSuite) addBridgeToken(tokenContract string, md banktypes.Metadata) {
	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0])
	suite.ctx = suite.ctx.WithEventManager(sdk.NewEventManager())
	err := suite.SendClaimReturnErr(&types.MsgBridgeTokenClaim{
		EventNonce:     oracleLastEventNonce + 1,
		BlockHeight:    uint64(suite.ctx.BlockHeight()),
		TokenContract:  tokenContract,
		Name:           md.Name,
		Symbol:         md.Symbol,
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChannelIbc:     "",
		ChainName:      suite.chainName,
	})
	suite.Require().NoError(err)

	suite.checkObservationState(suite.ctx, true)

	newOracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0])
	suite.Require().EqualValues(oracleLastEventNonce+1, newOracleLastEventNonce)
}

func (suite *KeeperTestSuite) registerCoin(bridgeDenom string) {
	_, err := suite.app.Erc20Keeper.RegisterCoin(sdk.WrapSDKContext(suite.ctx), &erc20types.MsgRegisterCoin{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Metadata: banktypes.Metadata{
			Description: "Test token",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "ttt",
					Exponent: 0,
					Aliases:  []string{bridgeDenom},
				},
				{
					Denom:    "TTT",
					Exponent: 18,
				},
			},
			Base:    "ttt",
			Display: "TTT",
			Name:    "Test Token",
			Symbol:  "TTT",
		},
	})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) checkObservationState(ctx context.Context, expect bool) {
	foundObservation := false
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, event := range sdkCtx.EventManager().Events() {
		if event.Type != types.EventTypeContractEvent {
			continue
		}
		suite.Require().False(foundObservation, "found multiple observation event")
		for _, attr := range event.Attributes {
			if attr.Key != types.AttributeKeyStateSuccess {
				continue
			}
			suite.Require().EqualValues(fmt.Sprintf("%v", expect), attr.Value)
			foundObservation = true
			break
		}
	}
	suite.Require().True(foundObservation, "not found observation event")
	sdkCtx.WithEventManager(sdk.NewEventManager())
}
