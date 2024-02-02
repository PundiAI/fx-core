package keeper_test

import (
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

func (suite *KeeperTestSuite) TestGrantPrivilege() {
	addrNotExist := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	testCase := []struct {
		name       string
		malleate   func() *types.MsgGrantPrivilege
		expectPass bool
		errMsg     string
	}{
		{
			name: "success - secp256k1",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				valAddr := sdk.ValAddress(acc)
				_, pkAny := suite.GenerateGrantPubkey()
				return &types.MsgGrantPrivilege{
					ValidatorAddress: valAddr.String(),
					FromAddress:      acc.String(),
					ToPubkey:         pkAny,
				}
			},
			expectPass: true,
		},
		{
			name: "success - eth_secp256k1",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				valAddr := sdk.ValAddress(acc)
				_, pkAny := suite.GenerateGrantPubkey()
				return &types.MsgGrantPrivilege{
					ValidatorAddress: valAddr.String(),
					FromAddress:      acc.String(),
					ToPubkey:         pkAny,
				}
			},
			expectPass: true,
		},
		{
			name: "invalid validator address",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				return &types.MsgGrantPrivilege{
					ValidatorAddress: acc.String(),
					FromAddress:      acc.String(),
				}
			},
			expectPass: false,
			errMsg:     "invalid Bech32 prefix; expected fxvaloper, got fx: invalid address",
		},
		{
			name: "validator not found",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				_, pkAny := suite.GenerateGrantPubkey()
				return &types.MsgGrantPrivilege{
					ValidatorAddress: sdk.ValAddress(addrNotExist).String(),
					FromAddress:      acc.String(),
					ToPubkey:         pkAny,
				}
			},
			expectPass: false,
			errMsg:     fmt.Sprintf("validator %s not found: unknown address", sdk.ValAddress(addrNotExist).String()),
		},
		{
			name: "from address not authorized",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				_, pkAny := suite.GenerateGrantPubkey()
				return &types.MsgGrantPrivilege{
					ValidatorAddress: sdk.ValAddress(acc).String(),
					FromAddress:      addrNotExist.String(),
					ToPubkey:         pkAny,
				}
			},
			expectPass: false,
			errMsg:     "from address not authorized: unauthorized",
		},
		{
			name: "val not authorized",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				_, pkAny := suite.GenerateGrantPubkey()
				msg := &types.MsgGrantPrivilege{
					ValidatorAddress: sdk.ValAddress(acc).String(),
					FromAddress:      acc.String(),
					ToPubkey:         pkAny,
				}
				_, err := suite.app.StakingKeeper.GrantPrivilege(sdk.WrapSDKContext(suite.ctx), msg)
				suite.Require().NoError(err)

				return &types.MsgGrantPrivilege{
					ValidatorAddress: sdk.ValAddress(acc).String(),
					FromAddress:      acc.String(),
					ToPubkey:         pkAny,
				}
			},
			expectPass: false,
			errMsg:     "from address not authorized: unauthorized",
		},
		{
			name: "invalid pubkey - empty",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				valAddr := sdk.ValAddress(acc)
				return &types.MsgGrantPrivilege{
					ValidatorAddress: valAddr.String(),
					FromAddress:      acc.String(),
					ToPubkey:         nil,
				}
			},
			expectPass: false,
			errMsg:     "empty pubkey: invalid pubkey",
		},
		{
			name: "invalid pubkey - invalid msg type",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				valAddr := sdk.ValAddress(acc)
				pk, _ := codectypes.NewAnyWithValue(&banktypes.MsgSend{})
				return &types.MsgGrantPrivilege{
					ValidatorAddress: valAddr.String(),
					FromAddress:      acc.String(),
					ToPubkey:         pk,
				}
			},
			expectPass: false,
			errMsg:     "expecting cryptotypes.PubKey, got *types.MsgSend: invalid pubkey",
		},
		{
			name: "invalid pubkey - invalid key",
			malleate: func() *types.MsgGrantPrivilege {
				acc := suite.valAccounts[0].GetAddress()
				valAddr := sdk.ValAddress(acc)
				pk, _ := codectypes.NewAnyWithValue(ed25519.GenPrivKey().PubKey())
				return &types.MsgGrantPrivilege{
					ValidatorAddress: valAddr.String(),
					FromAddress:      acc.String(),
					ToPubkey:         pk,
				}
			},
			expectPass: false,
			errMsg:     "expecting *secp256k1.PubKey or *ethsecp256k1.PubKey, got *ed25519.PubKey: invalid pubkey",
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			msg := tc.malleate()
			_, err := suite.app.StakingKeeper.GrantPrivilege(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(tc.errMsg, err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrantAccount() {
	key1 := helpers.NewPriKey()
	key2 := helpers.NewPriKey()
	eth3 := helpers.NewEthPrivKey()
	eth4 := helpers.NewEthPrivKey()

	any1, _ := codectypes.NewAnyWithValue(key1.PubKey())
	any2, _ := codectypes.NewAnyWithValue(key2.PubKey())
	any3, _ := codectypes.NewAnyWithValue(eth3.PubKey())
	any4, _ := codectypes.NewAnyWithValue(eth4.PubKey())

	addr1 := sdk.AccAddress(key1.PubKey().Address())
	addr2 := sdk.AccAddress(key2.PubKey().Address())
	addr3 := sdk.AccAddress(eth3.PubKey().Address())
	addr4 := sdk.AccAddress(eth4.PubKey().Address())

	acc1 := suite.valAccounts[0].GetAddress()
	valAddr1 := sdk.ValAddress(acc1)

	found := suite.app.StakingKeeper.HasValidatorOperator(suite.ctx, valAddr1)
	suite.Require().False(found)

	// val1 grant to key1
	sign, err := key1.Sign(types.GrantPrivilegeSignatureData(valAddr1, acc1, addr1))
	suite.Require().NoError(err)
	msgGrant := &types.MsgGrantPrivilege{FromAddress: acc1.String(), ToPubkey: any1, Signature: hex.EncodeToString(sign), ValidatorAddress: valAddr1.String()}
	_, err = suite.app.StakingKeeper.GrantPrivilege(sdk.WrapSDKContext(suite.ctx), msgGrant)
	suite.Require().NoError(err)

	operator, found := suite.app.StakingKeeper.GetValidatorOperator(suite.ctx, valAddr1)
	suite.Require().True(found)
	suite.Require().Equal(operator.String(), addr1.String())

	// val1 grant key1 to key2
	sign, err = key2.Sign(types.GrantPrivilegeSignatureData(valAddr1, addr1, addr2))
	suite.Require().NoError(err)
	msgGrant = &types.MsgGrantPrivilege{FromAddress: addr1.String(), ToPubkey: any2, Signature: hex.EncodeToString(sign), ValidatorAddress: valAddr1.String()}
	_, err = suite.app.StakingKeeper.GrantPrivilege(sdk.WrapSDKContext(suite.ctx), msgGrant)
	suite.Require().NoError(err)

	operator, found = suite.app.StakingKeeper.GetValidatorOperator(suite.ctx, valAddr1)
	suite.Require().True(found)
	suite.Require().Equal(operator.String(), addr2.String())

	// val1 grant key2 to eth3
	sign, err = eth3.Sign(types.GrantPrivilegeSignatureData(valAddr1, addr2, addr3))
	suite.Require().NoError(err)
	msgGrant = &types.MsgGrantPrivilege{FromAddress: addr2.String(), ToPubkey: any3, Signature: hex.EncodeToString(sign), ValidatorAddress: valAddr1.String()}
	_, err = suite.app.StakingKeeper.GrantPrivilege(sdk.WrapSDKContext(suite.ctx), msgGrant)
	suite.Require().NoError(err)

	operator, found = suite.app.StakingKeeper.GetValidatorOperator(suite.ctx, valAddr1)
	suite.Require().True(found)
	suite.Require().Equal(operator.String(), addr3.String())

	// val1 grant eth3 to eth4
	sign, err = eth4.Sign(types.GrantPrivilegeSignatureData(valAddr1, addr3, addr4))
	suite.Require().NoError(err)
	msgGrant = &types.MsgGrantPrivilege{FromAddress: addr3.String(), ToPubkey: any4, Signature: hex.EncodeToString(sign), ValidatorAddress: valAddr1.String()}
	_, err = suite.app.StakingKeeper.GrantPrivilege(sdk.WrapSDKContext(suite.ctx), msgGrant)
	suite.Require().NoError(err)

	operator, found = suite.app.StakingKeeper.GetValidatorOperator(suite.ctx, valAddr1)
	suite.Require().True(found)
	suite.Require().Equal(operator.String(), addr4.String())
}

func (suite *KeeperTestSuite) TestEditConsensusPubKey() {
	unKnownAddr := helpers.NewPriKey().PubKey().Address()
	_, tmAny := suite.GenerateConsKey()
	testCase := []struct {
		name       string
		malleate   func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string)
		expectPass bool
	}{
		{
			name: "success",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           tmAny,
				}, ""
			},
			expectPass: true,
		},
		{
			name: "success - granted",
			malleate: func(val sdk.ValAddress, _ sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				newFrom := helpers.NewEthPrivKey()
				suite.app.StakingKeeper.UpdateValidatorOperator(suite.ctx, val, newFrom.PubKey().Address().Bytes())
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             sdk.AccAddress(newFrom.PubKey().Address()).String(),
					Pubkey:           tmAny,
				}, ""
			},
			expectPass: true,
		},
		{
			name: "fail - invalid validator",
			malleate: func(_ sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: "",
					From:             from.String(),
					Pubkey:           tmAny,
				}, "empty address string is not allowed: invalid address"
			},
			expectPass: false,
		},
		{
			name: "fail - slashing double sign",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val)
				suite.Require().True(found)
				consAddr, err := validator.GetConsAddr()
				suite.Require().NoError(err)
				info, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, consAddr)
				suite.Require().True(found)
				info.JailedUntil = evidencetypes.DoubleSignJailEndTime
				suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, consAddr, info)

				validator.Jailed = true
				suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           tmAny,
				}, fmt.Sprintf("validator %s is jailed for double sign: invalid request", val.String())
			},
			expectPass: false,
		},
		{
			name: "fail - validator not found",
			malleate: func(_ sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: sdk.ValAddress(unKnownAddr).String(),
					From:             from.String(),
					Pubkey:           tmAny,
				}, fmt.Sprintf("validator %s not found: unknown address", sdk.ValAddress(unKnownAddr).String())
			},
			expectPass: false,
		},
		{
			name: "fail - from not authorized",
			malleate: func(val sdk.ValAddress, _ sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Pubkey:           tmAny,
				}, "from address not authorized: unauthorized"
			},
			expectPass: false,
		},
		{
			name: "fail - validator updating",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				err := suite.app.StakingKeeper.SetConsensusProcess(suite.ctx, val, tmAny.GetCachedValue().(cryptotypes.PubKey), types.ProcessStart)
				suite.Require().NoError(err)
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           tmAny,
				}, fmt.Sprintf("validator %s is updating consensus pubkey: invalid request", val.String())
			},
			expectPass: false,
		},
		{
			name: "fail - invalid pubkey",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				pkAny, err := codectypes.NewAnyWithValue(&banktypes.MsgSend{})
				if err != nil {
					panic(err)
				}
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           pkAny,
				}, "Expecting cryptotypes.PubKey, got *types.MsgSend: invalid type"
			},
			expectPass: false,
		},
		{
			name: "fail - pubkey exist",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				addr := suite.valAccounts[1].GetAddress()
				validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(addr))
				suite.Require().True(found)

				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           validator.ConsensusPubkey,
				}, "validator already exist for this pubkey; must use new validator pubkey"
			},
			expectPass: false,
		},
		{
			name: "fail - invalid pubkey type",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				pkAny, err := codectypes.NewAnyWithValue(helpers.NewEthPrivKey().PubKey())
				if err != nil {
					panic(err)
				}
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           pkAny,
				}, "got: eth_secp256k1, expected: [ed25519]: validator pubkey type is not supported"
			},
			expectPass: false,
		},
		{
			name: "fail - new pubkey updated by other validator",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				_, pkAny := suite.GenerateConsKey()
				pk := pkAny.GetCachedValue().(cryptotypes.PubKey)

				otherVal := sdk.ValAddress(suite.valAccounts[1].GetAddress())
				err := suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, otherVal, pk)
				suite.Require().NoError(err)

				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           pkAny,
				}, fmt.Sprintf("new consensus pubkey %s already exists: validator already exist for this pubkey; must use new validator pubkey", sdk.ConsAddress(pk.Address()))
			},
			expectPass: false,
		},
		{
			name: "fail - validator power more than 1/3",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				amount := math.NewInt(10000).Mul(math.NewInt(1e18))
				helpers.AddTestAddr(suite.app, suite.ctx, from, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)))

				validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val)
				suite.Require().True(found)
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, from, amount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				suite.Commit()

				valPower := suite.app.StakingKeeper.GetLastValidatorPower(suite.ctx, val)
				totalPower := suite.app.StakingKeeper.GetLastTotalPower(suite.ctx)

				return &types.MsgEditConsensusPubKey{
						ValidatorAddress: val.String(),
						From:             from.String(),
						Pubkey:           tmAny,
					}, fmt.Sprintf("total update power %d more than 1/3 total power %s: invalid request",
						valPower, totalPower.Quo(sdk.NewInt(3)).String())
			},
			expectPass: false,
		},
		{
			name: "fail - update power more than 1/3",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) (*types.MsgEditConsensusPubKey, string) {
				amount := math.NewInt(10000).Mul(math.NewInt(1e18))
				helpers.AddTestAddr(suite.app, suite.ctx, from, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)))

				validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(suite.valAccounts[1].GetAddress()))
				suite.Require().True(found)
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, from, amount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				suite.Commit()

				err = suite.app.StakingKeeper.SetConsensusProcess(suite.ctx, validator.GetOperator(), ed25519.GenPrivKey().PubKey(), types.ProcessStart)
				suite.Require().NoError(err)
				val1Power := suite.app.StakingKeeper.GetLastValidatorPower(suite.ctx, validator.GetOperator())

				valPower := suite.app.StakingKeeper.GetLastValidatorPower(suite.ctx, val)
				totalPower := suite.app.StakingKeeper.GetLastTotalPower(suite.ctx)

				return &types.MsgEditConsensusPubKey{
						ValidatorAddress: val.String(),
						From:             from.String(),
						Pubkey:           tmAny,
					}, fmt.Sprintf("total update power %d more than 1/3 total power %s: invalid request",
						valPower+val1Power, totalPower.Quo(sdk.NewInt(3)).String())
			},
			expectPass: false,
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			from := suite.valAccounts[0].GetAddress()
			msg, errMsg := tc.malleate(sdk.ValAddress(from), from)
			_, err := suite.app.StakingKeeper.EditConsensusPubKey(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errMsg, err.Error())
			}
		})
	}
}
