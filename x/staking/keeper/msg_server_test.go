package keeper_test

import (
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v5/testutil/helpers"
	"github.com/functionx/fx-core/v5/x/staking/types"
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
		malleate   func(val sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey
		expectPass bool
		errMsg     string
	}{
		{
			name: "success",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           tmAny,
				}
			},
			expectPass: true,
		},
		{
			name: "success - granted",
			malleate: func(val sdk.ValAddress, _ sdk.AccAddress) *types.MsgEditConsensusPubKey {
				newFrom := helpers.NewEthPrivKey()
				suite.app.StakingKeeper.UpdateValidatorOperator(suite.ctx, val, newFrom.PubKey().Address().Bytes())
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             sdk.AccAddress(newFrom.PubKey().Address()).String(),
					Pubkey:           tmAny,
				}
			},
			expectPass: true,
		},
		{
			name: "fail - invalid validator",
			malleate: func(_ sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: "",
					From:             from.String(),
					Pubkey:           tmAny,
				}
			},
			expectPass: false,
			errMsg:     "empty address string is not allowed: invalid address",
		},
		{
			name: "fail - validator not found",
			malleate: func(_ sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: sdk.ValAddress(unKnownAddr).String(),
					From:             from.String(),
					Pubkey:           tmAny,
				}
			},
			expectPass: false,
			errMsg:     fmt.Sprintf("validator %s not found: unknown address", sdk.ValAddress(unKnownAddr).String()),
		},
		{
			name: "fail - from not authorized",
			malleate: func(val sdk.ValAddress, _ sdk.AccAddress) *types.MsgEditConsensusPubKey {
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Pubkey:           tmAny,
				}
			},
			expectPass: false,
			errMsg:     "from address not authorized: unauthorized",
		},
		{
			name: "fail - invalid pubkey",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey {
				pkAny, err := codectypes.NewAnyWithValue(&banktypes.MsgSend{})
				if err != nil {
					panic(err)
				}
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           pkAny,
				}
			},
			expectPass: false,
			errMsg:     "Expecting cryptotypes.PubKey, got *types.MsgSend: invalid type",
		},
		{
			name: "fail - pubkey exist",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey {
				addr := suite.valAccounts[1].GetAddress()
				validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(addr))
				suite.Require().True(found)

				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           validator.ConsensusPubkey,
				}
			},
			expectPass: false,
			errMsg:     "validator already exist for this pubkey; must use new validator pubkey",
		},
		{
			name: "fail - invalid pubkey type",
			malleate: func(val sdk.ValAddress, from sdk.AccAddress) *types.MsgEditConsensusPubKey {
				pkAny, err := codectypes.NewAnyWithValue(helpers.NewEthPrivKey().PubKey())
				if err != nil {
					panic(err)
				}
				return &types.MsgEditConsensusPubKey{
					ValidatorAddress: val.String(),
					From:             from.String(),
					Pubkey:           pkAny,
				}
			},
			expectPass: false,
			errMsg:     "got: eth_secp256k1, expected: [ed25519]: validator pubkey type is not supported",
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			from := suite.valAccounts[0].GetAddress()
			msg := tc.malleate(sdk.ValAddress(from), from)
			_, err := suite.app.StakingKeeper.EditConsensusPubKey(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(tc.errMsg, err.Error())
			}
		})
	}
}
