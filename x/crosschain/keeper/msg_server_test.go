package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
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
				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{OracleAddress: msg.OracleAddress})
			},
			pass: false,
			err:  "oracle existed bridger address: invalid",
		},
		{
			name: "error - bridger address is bound",
			preRun: func(msg *types.MsgBondedOracle) {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, sdk.MustAccAddressFromBech32(msg.BridgerAddress), sdk.MustAccAddressFromBech32(msg.OracleAddress))
			},
			pass: false,
			err:  "bridger address is bound to oracle: invalid",
		},
		{
			name: "error - external address is bound",
			preRun: func(msg *types.MsgBondedOracle) {
				suite.Keeper().SetOracleAddrByExternalAddr(suite.Ctx, msg.ExternalAddress, sdk.MustAccAddressFromBech32(msg.OracleAddress))
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
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.Ctx)
				msg.DelegateAmount.Amount = delegateThreshold.Amount.Sub(sdkmath.NewInt(tmrand.Int63() - 1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgBondedOracle) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.Ctx)
				delegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.Ctx)
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
				ValidatorAddress: suite.ValAddr[oracleIndex].String(),
				DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(3) + 1) * 10_000).MulRaw(1e18)),
				ChainName:        suite.chainName,
			}

			testCase.preRun(msg)

			_, err := suite.MsgServer().BondedOracle(suite.Ctx, msg)
			if !testCase.pass {
				suite.Require().Error(err)
				suite.Require().EqualValues(testCase.err, err.Error())
				return
			}

			// success check
			suite.Require().NoError(err)

			// check oracle
			oracle, found := suite.Keeper().GetOracle(suite.Ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			suite.Require().True(found)
			suite.Require().NotNil(oracle)
			suite.Require().EqualValues(msg.OracleAddress, oracle.OracleAddress)
			suite.Require().EqualValues(msg.BridgerAddress, oracle.BridgerAddress)
			suite.Require().EqualValues(msg.ExternalAddress, oracle.ExternalAddress)
			suite.Require().True(oracle.Online)
			suite.Require().EqualValues(msg.ValidatorAddress, oracle.DelegateValidator)
			suite.Require().EqualValues(int64(0), oracle.SlashTimes)

			// check relationship
			oracleAddr, found := suite.Keeper().GetOracleAddrByBridgerAddr(suite.Ctx, sdk.MustAccAddressFromBech32(msg.BridgerAddress))
			suite.True(found)
			suite.Require().EqualValues(msg.OracleAddress, oracleAddr.String())

			oracleAddr, found = suite.Keeper().GetOracleAddrByExternalAddr(suite.Ctx, msg.ExternalAddress)
			suite.True(found)
			suite.Require().EqualValues(msg.OracleAddress, oracleAddr.String())

			// check power
			totalPower := suite.Keeper().GetLastTotalPower(suite.Ctx)
			suite.Require().EqualValues(msg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())

			// check delegate
			oracleDelegateAddr := oracle.GetDelegateAddress(suite.chainName)
			delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, oracleDelegateAddr, suite.ValAddr[oracleIndex])
			suite.NoError(err)
			suite.Require().EqualValues(oracleDelegateAddr.String(), delegation.DelegatorAddress)
			suite.Require().EqualValues(msg.ValidatorAddress, delegation.ValidatorAddress)
			suite.Truef(msg.DelegateAmount.Amount.Equal(delegation.GetShares().TruncateInt()), "expect:%s,actual:%s", msg.DelegateAmount.Amount.String(), delegation.GetShares().TruncateInt().String())

			startingInfo, err := suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, suite.ValAddr[oracleIndex], oracleDelegateAddr)
			suite.Require().NoError(err)
			suite.NotNil(startingInfo)
			suite.EqualValues(uint64(suite.Ctx.BlockHeight()), startingInfo.Height)
			suite.True(startingInfo.PreviousPeriod > 0)
			suite.EqualValues(sdkmath.LegacyNewDecFromInt(msg.DelegateAmount.Amount).String(), startingInfo.Stake.String())
		})
	}
}

func (suite *KeeperTestSuite) TestMsgAddDelegate() {
	initDelegateAmount := suite.Keeper().GetOracleDelegateThreshold(suite.Ctx).Amount
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
				oracle, _ := suite.Keeper().GetOracle(suite.Ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
				oracle.SlashTimes = 1
				suite.Keeper().SetOracle(suite.Ctx, oracle)
				slashFraction := suite.Keeper().GetSlashFraction(suite.Ctx)
				slashAmount := sdkmath.LegacyNewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
				randomAmount := tmrand.Int63n(slashAmount.QuoRaw(1e18).Int64()) + 1
				msg.Amount.Amount = sdkmath.NewInt(randomAmount).MulRaw(1e18).Sub(sdkmath.NewInt(1))
			},
			pass: false,
			err:  "not sufficient slash amount: invalid",
		},
		{
			name: "error - delegate amount less than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				params := suite.Keeper().GetParams(suite.Ctx)
				addDelegateThreshold := tmrand.Int63n(100000) + 1
				params.DelegateThreshold.Amount = initDelegateAmount.Add(sdkmath.NewInt(addDelegateThreshold).MulRaw(1e18))
				err := suite.Keeper().SetParams(suite.Ctx, &params)
				suite.Require().NoError(err)
				msg.Amount.Amount = sdkmath.NewInt(tmrand.Int63n(addDelegateThreshold) + 1).MulRaw(1e18).Sub(sdkmath.NewInt(1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				delegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.Ctx)
				delegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.Ctx)
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
				oracle, _ := suite.Keeper().GetOracle(suite.Ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
				oracle.SlashTimes = 1
				oracle.Online = false
				suite.Keeper().SetOracle(suite.Ctx, oracle)

				slashFraction := suite.Keeper().GetSlashFraction(suite.Ctx)
				slashAmount := sdkmath.LegacyNewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
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
				oracle, _ := suite.Keeper().GetOracle(suite.Ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
				oracle.SlashTimes = 1
				oracle.Online = false
				suite.Keeper().SetOracle(suite.Ctx, oracle)

				slashFraction := suite.Keeper().GetSlashFraction(suite.Ctx)
				slashAmount := sdkmath.LegacyNewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(oracle.SlashTimes).TruncateInt()
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
			_, err := suite.MsgServer().BondedOracle(suite.Ctx, &types.MsgBondedOracle{
				OracleAddress:    suite.oracleAddrs[oracleIndex].String(),
				BridgerAddress:   suite.bridgerAddrs[oracleIndex].String(),
				ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[oracleIndex].PublicKey),
				ValidatorAddress: suite.ValAddr[oracleIndex].String(),
				DelegateAmount:   types.NewDelegateAmount(initDelegateAmount),
				ChainName:        suite.chainName,
			})
			suite.Require().NoError(err)

			oracleDelegateThreshold := suite.Keeper().GetOracleDelegateThreshold(suite.Ctx)
			oracleDelegateMultiple := suite.Keeper().GetOracleDelegateMultiple(suite.Ctx)
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

			_, err = suite.MsgServer().AddDelegate(suite.Ctx, msg)
			if !testCase.pass {
				suite.Require().Error(err)
				suite.Require().EqualValues(testCase.err, err.Error())
				return
			}

			// success check
			suite.Require().NoError(err)

			// check oracle
			oracle, found := suite.Keeper().GetOracle(suite.Ctx, sdk.MustAccAddressFromBech32(msg.OracleAddress))
			suite.Require().True(found)
			suite.Require().NotNil(oracle)
			suite.Require().EqualValues(msg.OracleAddress, oracle.OracleAddress)
			suite.Require().True(oracle.Online)
			suite.Require().EqualValues(0, oracle.SlashTimes)

			// check power
			totalPower := suite.Keeper().GetLastTotalPower(suite.Ctx)
			expectDelegateAmount := testCase.expectDelegateAmount(msg)
			suite.Require().EqualValues(expectDelegateAmount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())

			// check delegate
			oracleDelegateAddr := oracle.GetDelegateAddress(suite.chainName)
			delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, oracleDelegateAddr, suite.ValAddr[oracleIndex])
			suite.NoError(err)
			suite.Require().EqualValues(oracleDelegateAddr.String(), delegation.DelegatorAddress)
			suite.Require().EqualValues(oracle.DelegateValidator, delegation.ValidatorAddress)
			suite.Truef(expectDelegateAmount.Equal(delegation.GetShares().TruncateInt()), "expect:%s,actual:%s", expectDelegateAmount.String(), delegation.GetShares().TruncateInt().String())
			startingInfo, err := suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, suite.ValAddr[oracleIndex], oracleDelegateAddr)
			suite.Require().NoError(err)
			suite.NotNil(startingInfo)
			suite.EqualValues(uint64(suite.Ctx.BlockHeight()), startingInfo.Height)
			suite.True(startingInfo.PreviousPeriod > 0)
			suite.EqualValues(sdkmath.LegacyNewDecFromInt(expectDelegateAmount).String(), startingInfo.Stake.String())
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
			ValidatorAddress: suite.ValAddr[i].String(),
			DelegateAmount:   types.NewDelegateAmount(delegateAmount),
			ChainName:        suite.chainName,
		}
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, bondedMsg)
		suite.NoError(err)
	}

	for i := 0; i < len(suite.oracleAddrs); i++ {
		_, err := suite.MsgServer().EditBridger(suite.Ctx, &types.MsgEditBridger{
			ChainName:      suite.chainName,
			OracleAddress:  suite.oracleAddrs[i].String(),
			BridgerAddress: suite.bridgerAddrs[i].String(),
		})
		suite.Require().Error(err)

		_, err = suite.MsgServer().EditBridger(suite.Ctx, &types.MsgEditBridger{
			ChainName:      suite.chainName,
			OracleAddress:  suite.oracleAddrs[i].String(),
			BridgerAddress: sdk.AccAddress(suite.ValAddr[i]).String(),
		})
		suite.NoError(err)
	}

	for _, bridger := range suite.bridgerAddrs {
		suite.False(suite.Keeper().HasOracleAddrByBridgerAddr(suite.Ctx, bridger))
	}
}

func (suite *KeeperTestSuite) TestMsgSetOracleSetConfirm() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.ValAddr[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(suite.Ctx, normalMsg)
	suite.Require().NoError(err)

	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.Ctx)
	suite.Require().EqualValues(0, latestOracleSetNonce)

	suite.Commit()

	latestOracleSetNonce = suite.Keeper().GetLatestOracleSetNonce(suite.Ctx)
	suite.Require().EqualValues(1, latestOracleSetNonce)

	suite.Require().True(suite.Keeper().HasOracleSetRequest(suite.Ctx, 1))

	suite.Require().False(suite.Keeper().HasOracleSetRequest(suite.Ctx, 2))

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.Ctx, 1)
	suite.Require().EqualValues(uint64(1), nonce1OracleSet.Nonce)
	suite.Require().EqualValues(uint64(suite.Ctx.BlockHeight()-1), nonce1OracleSet.Height)
	suite.Require().EqualValues(1, len(nonce1OracleSet.Members))
	suite.Require().EqualValues(normalMsg.ExternalAddress, nonce1OracleSet.Members[0].ExternalAddress)
	suite.Require().EqualValues(math.MaxUint32, nonce1OracleSet.Members[0].Power)

	gravityId := suite.Keeper().GetGravityID(suite.Ctx)
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
		_, err = suite.MsgServer().OracleSetConfirm(suite.Ctx, testData.msg)
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
	_, err = suite.MsgServer().OracleSetConfirm(suite.Ctx, normalOracleSetConfirmMsg)
	suite.Require().NoError(err)

	endBlockBeforeLatestOracleSet := suite.Keeper().GetLatestOracleSet(suite.Ctx)
	suite.Require().NotNil(endBlockBeforeLatestOracleSet)
}

func (suite *KeeperTestSuite) TestClaimWithOracleOnline() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.ValAddr[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(suite.Ctx, normalMsg)
	suite.Require().NoError(err)

	suite.Commit()

	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.Ctx)
	suite.Require().EqualValues(1, latestOracleSetNonce)

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.Ctx, latestOracleSetNonce)
	suite.Require().EqualValues(uint64(1), nonce1OracleSet.Nonce)
	suite.Require().EqualValues(uint64(suite.Ctx.BlockHeight()-1), nonce1OracleSet.Height)

	var gravityId string
	suite.Require().NotPanics(func() {
		gravityId = suite.Keeper().GetGravityID(suite.Ctx)
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
	oracle, found := suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[0])
	suite.Require().True(found)
	oracle.Online = true
	suite.Keeper().SetOracle(suite.Ctx, oracle)

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
	_, err = suite.MsgServer().OracleSetConfirm(suite.Ctx, normalOracleSetConfirmMsg)
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
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
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
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
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
				suite.Keeper().StoreOracleSet(suite.Ctx, &types.OracleSet{
					Nonce:   msg.OracleSetNonce,
					Height:  msg.BlockHeight,
					Members: msg.Members,
				})
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
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
				suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
					BatchNonce:    msg.BatchNonce,
					TokenContract: msg.TokenContract,
				}))
				for i, oracle := range suite.oracleAddrs {
					eventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, oracle)
					msg.EventNonce = eventNonce + 1
					msg.BridgerAddress = suite.bridgerAddrs[i].String()
					ctxWithGasMeter := suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
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
					ValidatorAddress: suite.ValAddr[0].String(),
					DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
					ChainName:        suite.chainName,
				}
				_, err := suite.MsgServer().BondedOracle(suite.Ctx, msg)
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
		ValidatorAddress: suite.ValAddr[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(suite.Ctx, normalMsg)
	suite.Require().NoError(err)

	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, suite.oracleAddrs[0])
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

func (suite *KeeperTestSuite) TestMsgUpdateChainOracles() {
	updateOracle := &types.MsgUpdateChainOracles{
		Oracles:   []string{},
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ChainName: suite.chainName,
	}
	for _, oracle := range suite.oracleAddrs {
		updateOracle.Oracles = append(updateOracle.Oracles, oracle.String())
	}

	_, err := suite.MsgServer().UpdateChainOracles(suite.Ctx, updateOracle)
	suite.Require().NoError(err)
	for _, oracle := range suite.oracleAddrs {
		suite.Require().True(suite.Keeper().IsProposalOracle(suite.Ctx, oracle.String()))
	}

	updateOracle.Oracles = []string{}
	number := tmrand.Intn(100)
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	_, err = suite.MsgServer().UpdateChainOracles(suite.Ctx, updateOracle)
	suite.Require().NoError(err)

	updateOracle.Oracles = []string{}
	number = tmrand.Intn(2) + 101
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	_, err = suite.MsgServer().UpdateChainOracles(suite.Ctx, updateOracle)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) BondedOracle() {
	_, err := suite.MsgServer().BondedOracle(suite.Ctx, &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.ValAddr[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt((tmrand.Int63n(5) + 1) * 10_000).MulRaw(1e18)),
		ChainName:        suite.chainName,
	})
	suite.Require().NoError(err)

	oracleLastEventNonce := suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, suite.oracleAddrs[0])
	suite.Require().EqualValues(0, oracleLastEventNonce)
}
