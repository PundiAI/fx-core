package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestOracleUpdate() {
	if len(suite.oracleAddrs) < 10 {
		return
	}
	for i := 0; i < 10; i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.ValPrivs[i].ValAddress().String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(100).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		suite.Require().NoError(msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, msgBondedOracle)

		suite.Require().NoError(err)
		suite.EndBlocker()
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.Ctx)
		suite.Require().NotNil(oracleSets)
		suite.Require().EqualValues(i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.Ctx)
		expectPower := sdkmath.NewInt(100).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		suite.Require().True(expectPower.Equal(power))
	}

	bridgeToken := helpers.GenExternalAddr(suite.chainName)

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}
		err := suite.SendClaimReturnErr(addBridgeTokenClaim)
		suite.Require().NoError(err)
		endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.Ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(endBlockBeforeAttestation)
		suite.Require().False(endBlockBeforeAttestation.Observed)
		suite.Require().NotNil(endBlockBeforeAttestation.Votes)
		suite.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))

		suite.EndBlocker()
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.Ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(endBlockAfterAttestation)
		suite.Require().False(endBlockAfterAttestation.Observed)
	}

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[6].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	err := suite.SendClaimReturnErr(addBridgeTokenClaim)
	suite.Require().NoError(err)
	suite.EndBlocker()
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	attestation := suite.Keeper().GetAttestation(suite.Ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())

	suite.Require().NotNil(attestation)
	suite.Require().True(attestation.Observed)

	var newOracleList []string
	for i := 0; i < 7; i++ {
		newOracleList = append(newOracleList, suite.oracleAddrs[i].String())
	}
	_, err = suite.MsgServer().UpdateChainOracles(suite.Ctx, &types.MsgUpdateChainOracles{
		ChainName: suite.chainName,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Oracles:   newOracleList,
	})

	suite.Require().ErrorIs(types.ErrInvalid, err)

	expectTotalPower := sdkmath.NewInt(100).MulRaw(1e18).Mul(sdkmath.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := suite.Keeper().GetLastTotalPower(suite.Ctx)
	suite.Require().True(expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalOracleChangePowerThreshold.Mul(expectTotalPower).Quo(sdkmath.NewInt(100))

	expectDeletePower := sdkmath.NewInt(100).MulRaw(1e18).Mul(sdkmath.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	suite.Require().EqualValues(fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newOracleList2 []string
	for i := 0; i < 8; i++ {
		newOracleList2 = append(newOracleList2, suite.oracleAddrs[i].String())
	}
	_, err = suite.MsgServer().UpdateChainOracles(suite.Ctx, &types.MsgUpdateChainOracles{
		ChainName: suite.chainName,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Oracles:   newOracleList2,
	})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestAttestationAfterOracleUpdate() {
	if len(suite.bridgerAddrs) < 20 {
		return
	}
	for i := 0; i < 20; i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.ValPrivs[i].ValAddress().String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(100).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, msgBondedOracle)
		suite.Require().NoError(err)
		suite.EndBlocker()
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.Ctx)
		suite.Require().NotNil(oracleSets)
		suite.Require().EqualValues(i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.Ctx)
		expectPower := sdkmath.NewInt(100).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		suite.Require().True(expectPower.Equal(power))
	}

	{
		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  helpers.GenExternalAddr(suite.chainName),
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 13; i++ {
			firstBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[i].String()
			err := suite.SendClaimReturnErr(firstBridgeTokenClaim)
			suite.Require().NoError(err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			suite.Require().NotNil(endBlockBeforeAttestation)
			suite.Require().False(endBlockBeforeAttestation.Observed)
			suite.Require().NotNil(endBlockBeforeAttestation.Votes)
			suite.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))

			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			suite.Require().NotNil(endBlockAfterAttestation)
			suite.Require().False(endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[13].String()
		err := suite.SendClaimReturnErr(firstBridgeTokenClaim)
		suite.Require().NoError(err)
		suite.EndBlocker()
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		attestation := suite.Keeper().GetAttestation(suite.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())

		suite.Require().NotNil(attestation)
		suite.Require().True(attestation.Observed)
	}

	{
		secondBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     2,
			BlockHeight:    1001,
			TokenContract:  helpers.GenExternalAddr(suite.chainName),
			Name:           "Test Token2",
			Symbol:         "TEST2",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 6; i++ {
			secondBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[i].String()
			err := suite.SendClaimReturnErr(secondBridgeTokenClaim)
			suite.Require().NoError(err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			suite.Require().NotNil(endBlockBeforeAttestation)
			suite.Require().False(endBlockBeforeAttestation.Observed)
			suite.Require().NotNil(endBlockBeforeAttestation.Votes)
			suite.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))

			suite.EndBlocker()
			suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			suite.Require().NotNil(endBlockAfterAttestation)
			suite.Require().False(endBlockAfterAttestation.Observed)
		}

		secondClaimAttestation := suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(secondClaimAttestation)
		suite.Require().False(secondClaimAttestation.Observed)
		suite.Require().NotNil(secondClaimAttestation.Votes)
		suite.Require().EqualValues(6, len(secondClaimAttestation.Votes))

		var newOracleList []string
		for i := 0; i < 15; i++ {
			newOracleList = append(newOracleList, suite.oracleAddrs[i].String())
		}
		_, err := suite.MsgServer().UpdateChainOracles(suite.Ctx, &types.MsgUpdateChainOracles{
			Oracles:   newOracleList,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		suite.Require().NoError(err)
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		suite.EndBlocker()

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(secondClaimAttestation)
		suite.Require().False(secondClaimAttestation.Observed)
		suite.Require().NotNil(secondClaimAttestation.Votes)
		suite.Require().EqualValues(6, len(secondClaimAttestation.Votes))

		activeOracles := suite.Keeper().GetAllOracles(suite.Ctx, true)
		suite.Require().NotNil(activeOracles)
		suite.Require().EqualValues(15, len(activeOracles))
		for i := 0; i < 15; i++ {
			suite.Require().NotNil(newOracleList[i], activeOracles[i].OracleAddress)
		}

		var newOracleList2 []string
		for i := 0; i < 11; i++ {
			newOracleList2 = append(newOracleList2, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainOracles(suite.Ctx, &types.MsgUpdateChainOracles{
			Oracles:   newOracleList2,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		suite.Require().NoError(err)
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		suite.EndBlocker()

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(secondClaimAttestation)
		suite.Require().False(secondClaimAttestation.Observed)
		suite.Require().NotNil(secondClaimAttestation.Votes)
		suite.Require().EqualValues(6, len(secondClaimAttestation.Votes))

		activeOracles = suite.Keeper().GetAllOracles(suite.Ctx, true)
		suite.Require().NotNil(activeOracles)
		suite.Require().EqualValues(11, len(activeOracles))
		for i := 0; i < 11; i++ {
			suite.Require().NotNil(newOracleList2[i], activeOracles[i].OracleAddress)
		}

		var newOracleList3 []string
		for i := 0; i < 10; i++ {
			newOracleList3 = append(newOracleList3, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainOracles(suite.Ctx, &types.MsgUpdateChainOracles{
			Oracles:   newOracleList3,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		suite.Require().NoError(err)
		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		suite.EndBlocker()

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(secondClaimAttestation)
		suite.Require().False(secondClaimAttestation.Observed)
		suite.Require().NotNil(secondClaimAttestation.Votes)
		suite.Require().EqualValues(6, len(secondClaimAttestation.Votes))

		activeOracles = suite.Keeper().GetAllOracles(suite.Ctx, true)
		suite.Require().NotNil(activeOracles)
		suite.Require().EqualValues(10, len(activeOracles))
		for i := 0; i < 10; i++ {
			suite.Require().NotNil(newOracleList3[i], activeOracles[i].OracleAddress)
		}

		secondBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[6].String()
		err = suite.SendClaimReturnErr(secondBridgeTokenClaim)
		suite.Require().NoError(err)

		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
		suite.EndBlocker()

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		suite.Require().NotNil(secondClaimAttestation)
		suite.Require().True(secondClaimAttestation.Observed)
		suite.Require().NotNil(secondClaimAttestation.Votes)
		suite.Require().EqualValues(7, len(secondClaimAttestation.Votes))
	}
}

func (suite *KeeperTestSuite) TestOracleDelete() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.ValPrivs[i].ValAddress().String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(100).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		suite.Require().NoError(msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, msgBondedOracle)
		suite.Require().NoError(err)
	}
	suite.EndBlocker()
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	allOracles := suite.Keeper().GetAllOracles(suite.Ctx, false)
	suite.Require().NotNil(allOracles)
	suite.Require().EqualValues(len(suite.oracleAddrs), len(allOracles))

	oracle := suite.oracleAddrs[0]
	bridger := suite.bridgerAddrs[0]
	externalAddress := suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey)

	oracleAddr, found := suite.Keeper().GetOracleAddrByBridgerAddr(suite.Ctx, bridger)
	suite.Require().True(found)
	suite.Require().EqualValues(oracle.String(), oracleAddr.String())

	oracleAddr, found = suite.Keeper().GetOracleAddrByExternalAddr(suite.Ctx, externalAddress)
	suite.Require().True(found)
	suite.Require().EqualValues(oracle.String(), oracleAddr.String())

	oracleData, found := suite.Keeper().GetOracle(suite.Ctx, oracle)
	suite.Require().True(found)
	suite.Require().NotNil(oracleData)
	suite.Require().EqualValues(oracle.String(), oracleData.OracleAddress)
	suite.Require().EqualValues(bridger.String(), oracleData.BridgerAddress)
	suite.Require().EqualValues(externalAddress, oracleData.ExternalAddress)

	suite.Require().True(sdkmath.NewInt(100).MulRaw(1e18).Equal(oracleData.DelegateAmount))

	newOracleAddressList := make([]string, 0, len(suite.oracleAddrs)-1)
	for _, address := range suite.oracleAddrs[1:] {
		newOracleAddressList = append(newOracleAddressList, address.String())
	}

	_, err := suite.MsgServer().UpdateChainOracles(suite.Ctx, &types.MsgUpdateChainOracles{
		Oracles:   newOracleAddressList,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ChainName: suite.chainName,
	})
	suite.Require().NoError(err)
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.EndBlocker()

	oracleAddr, found = suite.Keeper().GetOracleAddrByBridgerAddr(suite.Ctx, bridger)
	suite.Require().True(found)
	suite.Require().Equal(oracleAddr, oracle)

	oracleAddr, found = suite.Keeper().GetOracleAddrByExternalAddr(suite.Ctx, externalAddress)
	suite.Require().True(found)
	suite.Require().Equal(oracleAddr, oracle)

	_, found = suite.Keeper().GetOracle(suite.Ctx, oracle)
	suite.Require().True(found)
}

func (suite *KeeperTestSuite) TestOracleSetSlash() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.ValPrivs[i].ValAddress().String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(100).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		suite.Require().NoError(msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, msgBondedOracle)
		suite.Require().NoError(err)
	}
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.Keeper().EndBlocker(suite.Ctx)

	allOracles := suite.Keeper().GetAllOracles(suite.Ctx, false)
	suite.Require().NotNil(allOracles)
	suite.Require().Equal(len(suite.oracleAddrs), len(allOracles))

	oracleSets := suite.Keeper().GetOracleSets(suite.Ctx)
	suite.Require().NotNil(oracleSets)
	suite.Require().EqualValues(1, len(oracleSets))

	for i := 0; i < len(suite.oracleAddrs)-1; i++ {
		externalAddress, signature := suite.SignOracleSetConfirm(suite.externalPris[i], oracleSets[0])
		oracleSetConfirm := &types.MsgOracleSetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  suite.bridgerAddrs[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		}
		suite.Require().NoError(oracleSetConfirm.ValidateBasic())
		_, err := suite.MsgServer().OracleSetConfirm(suite.Ctx, oracleSetConfirm)
		suite.Require().NoError(err)
	}

	suite.Keeper().EndBlocker(suite.Ctx)
	oracleSetHeight := int64(oracleSets[0].Height)
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.EndBlocker()

	oracle, found := suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	suite.Require().True(found)
	suite.Require().True(oracle.Online)
	suite.Require().Equal(int64(0), oracle.SlashTimes)

	suite.Ctx = suite.Ctx.WithBlockHeight(oracleSetHeight + int64(suite.Keeper().GetParams(suite.Ctx).SignedWindow) + 1)
	suite.Keeper().EndBlocker(suite.Ctx)

	oracle, found = suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	suite.Require().True(found)
	suite.Require().False(oracle.Online)
	suite.Require().Equal(int64(1), oracle.SlashTimes)
}

func (suite *KeeperTestSuite) TestSlashOracle() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.ValPrivs[i].ValAddress().String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(100).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		suite.Require().NoError(msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, msgBondedOracle)
		suite.Require().NoError(err)
	}

	params := suite.Keeper().GetParams(suite.Ctx)
	err := suite.Keeper().SetParams(suite.Ctx, &params)
	suite.Require().NoError(err)
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[i])
		suite.Require().True(found)
		suite.Require().True(oracle.Online)
		suite.Require().Equal(int64(0), oracle.SlashTimes)

		suite.Keeper().SlashOracle(suite.Ctx, oracle.OracleAddress)

		oracle, found = suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[i])
		suite.Require().True(found)
		suite.Require().False(oracle.Online)
		suite.Require().Equal(int64(1), oracle.SlashTimes)
	}

	// repeat slash test.
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[i])
		suite.Require().True(found)
		suite.Require().False(oracle.Online)
		suite.Require().Equal(int64(1), oracle.SlashTimes)

		suite.Keeper().SlashOracle(suite.Ctx, oracle.OracleAddress)

		oracle, found = suite.Keeper().GetOracle(suite.Ctx, suite.oracleAddrs[i])
		suite.Require().True(found)
		suite.Require().False(oracle.Online)
		suite.Require().Equal(int64(1), oracle.SlashTimes)
	}
}
