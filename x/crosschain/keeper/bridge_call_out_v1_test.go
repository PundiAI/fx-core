package keeper_test

/*
func (suite *KeeperTestSuite) TestKeeper_BridgeCallRefund() {
	suite.bondedOracle()
	oracleAddr, found := suite.Keeper().GetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0])
	suite.Require().True(found)
	suite.Require().EqualValues(oracleAddr.Bytes(), suite.oracleAddrs[0].Bytes())
	suite.Commit()

	oracleAddr, found = suite.Keeper().GetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0])
	suite.Require().True(found)
	suite.Require().EqualValues(oracleAddr.Bytes(), suite.oracleAddrs[0].Bytes())

	bridgeToken := helpers.GenHexAddress()
	bridgeTokenStr := types.ExternalAddrToStr(suite.chainName, bridgeToken.Bytes())
	suite.addBridgeToken(bridgeTokenStr, fxtypes.GetCrossChainMetadataManyToOne("test token", "TTT", 18))

	suite.registerCoin(types.NewBridgeDenom(suite.chainName, bridgeTokenStr))

	fxAddr1 := helpers.GenHexAddress()
	randomBlock := tmrand.Int63n(1000000000)
	randomAmount := tmrand.Int63n(1000000000)
	claim := &types.MsgSendToFxClaim{
		EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, suite.oracleAddrs[0]) + 1,
		BlockHeight:    uint64(randomBlock),
		TokenContract:  bridgeTokenStr,
		Amount:         sdkmath.NewInt(randomAmount),
		Sender:         types.ExternalAddrToStr(suite.chainName, helpers.GenHexAddress().Bytes()),
		Receiver:       sdk.AccAddress(fxAddr1.Bytes()).String(),
		TargetIbc:      "",
		BridgerAddress: suite.bridgerAddrs[0].String(),
	}
	suite.SendClaim(claim)

	pair, b := suite.App.Erc20Keeper.GetTokenPair(suite.Ctx, "ttt")
	suite.True(b)
	suite.Equal(sdkmath.NewInt(randomAmount), suite.App.BankKeeper.GetBalance(suite.Ctx, fxAddr1.Bytes(), pair.Denom).Amount)

	// todo remove after bridge call refactor
	suite.MintTokenToModule(erc20types.ModuleName, sdk.NewCoin(types.NewBridgeDenom(suite.chainName, bridgeTokenStr), sdkmath.NewInt(randomAmount)))

	bridgeCallRefundAddr := helpers.GenAccAddress()
	_, err := suite.MsgServer().BridgeCall(suite.Ctx, &types.MsgBridgeCall{
		ChainName: suite.chainName,
		Sender:    sdk.AccAddress(fxAddr1.Bytes()).String(),
		Refund:    bridgeCallRefundAddr.String(),
		Coins:     sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdkmath.NewInt(randomAmount))),
		Value:     sdkmath.ZeroInt(),
	})
	suite.NoError(err)

	suite.Equal(sdkmath.NewInt(0), suite.App.BankKeeper.GetBalance(suite.Ctx, fxAddr1.Bytes(), pair.Denom).Amount)

	var outgoingBridgeCall *types.OutgoingBridgeCall
	suite.Keeper().IterateOutgoingBridgeCallsByAddress(suite.Ctx, types.ExternalAddrToStr(suite.chainName, fxAddr1.Bytes()), func(value *types.OutgoingBridgeCall) bool {
		outgoingBridgeCall = value
		return true
	})
	suite.NotNil(outgoingBridgeCall)

	// Triggering the SendtoFx claim once is just to trigger timeout
	sendToFxClaim := &types.MsgSendToFxClaim{
		EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.Ctx, suite.oracleAddrs[0]) + 1,
		BlockHeight:    outgoingBridgeCall.Timeout,
		TokenContract:  bridgeTokenStr,
		Amount:         sdkmath.NewInt(randomAmount),
		Sender:         types.ExternalAddrToStr(suite.chainName, helpers.GenHexAddress().Bytes()),
		Receiver:       sdk.AccAddress(fxAddr1.Bytes()).String(),
		TargetIbc:      hex.EncodeToString([]byte(fxtypes.ERC20Target)),
		BridgerAddress: suite.bridgerAddrs[0].String(),
	}
	suite.SendClaim(sendToFxClaim)
	// expect balance = sendToFx value + outgointBridgeCallRefund value
	suite.CheckBalanceOf(pair.GetERC20Contract(), fxAddr1, big.NewInt(randomAmount))
	suite.Equal(sdkmath.NewInt(0), suite.App.BankKeeper.GetBalance(suite.Ctx, fxAddr1.Bytes(), pair.Denom).Amount)
	suite.Equal(sdkmath.NewInt(randomAmount), suite.App.BankKeeper.GetBalance(suite.Ctx, bridgeCallRefundAddr, pair.Denom).Amount)
}
*/
