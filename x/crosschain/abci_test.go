package crosschain_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/x/crosschain"
	"github.com/functionx/fx-core/x/crosschain/types"
	ibcTransferTypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

func TestABCIEndBlockDepositClaim(t *testing.T) {
	//myApp.SetAppLog(server.ZeroLogWrapper{Logger: log.Logger.Level(zerolog.DebugLevel)})
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 4)
	keep := myApp.BscKeeper
	var err error

	totalDepositBefore := keep.GetTotalStake(ctx)
	require.EqualValues(t, sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()), totalDepositBefore)

	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:   oracleAddressList[0].String(),
		BridgeAddress:   orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Name:           "Pundix Reward Token",
		Symbol:         "PURES",
		Decimals:       18,
		BridgerAddress: orchestratorAddressList[0].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      chainName,
	}
	_, err = h(ctx, addBridgeTokenClaim)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

	sendToFxClaim := &types.MsgSendToFxClaim{
		EventNonce:     2,
		BlockHeight:    1001,
		TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Amount:         sdk.NewInt(1234),
		Sender:         "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Receiver:       "fx16wvwsmpp4y4ttgzknyr6kqla877jud6u04lqey",
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: orchestratorAddressList[0].String(),
		ChainName:      chainName,
	}
	_, err = h(ctx, sendToFxClaim)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

	receiveAddr, err := sdk.AccAddressFromBech32(sendToFxClaim.Receiver)
	require.NoError(t, err)
	allBalances := myApp.BankKeeper.GetAllBalances(ctx, receiveAddr)
	//t.Logf("%s allBalances:%s", receiveAddr.String(), allBalances)
	tokenContract := common.HexToAddress(addBridgeTokenClaim.TokenContract).Hex()
	// transfer/channel-0/bscPURES
	tokenName := fmt.Sprintf("%s%s", chainName, tokenContract)
	if len(addBridgeTokenClaim.ChannelIbc) > 0 {
		channel, err := hex.DecodeString(addBridgeTokenClaim.ChannelIbc)
		require.NoError(t, err)
		tokenName = ibcTransferTypes.DenomTrace{
			Path:      string(channel),
			BaseDenom: fmt.Sprintf("%s%s", chainName, tokenContract),
		}.IBCDenom()
	}
	require.EqualValues(t, fmt.Sprintf("%s%s", sendToFxClaim.Amount.String(), tokenName), allBalances.String())
}

func TestOracleUpdate(t *testing.T) {
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 25)
	keeper := myApp.BscKeeper
	var err error

	for i := 0; i < 10; i++ {
		_, err = h(ctx, &types.MsgCreateOracleBridger{
			OracleAddress:   oracleAddressList[i].String(),
			BridgeAddress:   orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
			ChainName:       chainName,
		})
		require.NoError(t, err)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		oracleSets := keeper.GetOracleSets(ctx)
		require.NotNil(t, oracleSets)
		require.EqualValues(t, i+1, len(oracleSets))

		power := keeper.GetLastTotalPower(ctx)
		expectPower := minStakeAmount.Mul(sdk.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(t, expectPower.Equal(power))
	}

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
			Name:           "Pundix Reward Token",
			Symbol:         "PURES",
			Decimals:       18,
			BridgerAddress: orchestratorAddressList[i].String(),
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      chainName,
		}
		_, err = h(ctx, addBridgeTokenClaim)
		require.NoError(t, err)
		endBlockBeforeAttestation := keeper.GetAttestation(ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(t, endBlockBeforeAttestation)
		require.False(t, endBlockBeforeAttestation.Observed)
		require.NotNil(t, endBlockBeforeAttestation.Votes)
		require.EqualValues(t, i+1, len(endBlockBeforeAttestation.Votes))

		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		endBlockAfterAttestation := keeper.GetAttestation(ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(t, endBlockAfterAttestation)
		require.False(t, endBlockAfterAttestation.Observed)
	}

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Name:           "Pundix Reward Token",
		Symbol:         "PURES",
		Decimals:       18,
		BridgerAddress: orchestratorAddressList[6].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      chainName,
	}
	_, err = h(ctx, addBridgeTokenClaim)
	require.NoError(t, err)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	attestation := keeper.GetAttestation(ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())

	require.NotNil(t, attestation)

	require.True(t, attestation.Observed)
	//t.Logf("attestation votes:%s", attestation.Votes)

	proposalHandler := crosschain.NewCrossChainProposalHandler(myApp.CrosschainKeeper)

	var newOralceList []string
	for i := 0; i < 7; i++ {
		newOralceList = append(newOralceList, oracleAddressList[i].String())
	}
	err = proposalHandler(ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 1: try update chain oracle power >= 30%, expect error",
		Description: "",
		Oracles:     newOralceList,
		ChainName:   chainName,
	})
	require.ErrorIs(t, types.ErrInvalid, err)

	expectTotalPower := minStakeAmount.Mul(sdk.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := keeper.GetLastTotalPower(ctx)
	require.True(t, expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalOracleChangePowerThreshold.Mul(expectTotalPower).Quo(sdk.NewInt(100))

	expectDeletePower := minStakeAmount.Mul(sdk.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	require.EqualValues(t, fmt.Sprintf("max change power!maxChangePower:%s,deletePower:%s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newOracleList2 []string
	for i := 0; i < 8; i++ {
		newOracleList2 = append(newOracleList2, oracleAddressList[i].String())
	}
	err = proposalHandler(ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 2: try update chain oracle power <= 30%, expect success",
		Description: "",
		Oracles:     newOracleList2,
		ChainName:   chainName,
	})
	require.NoError(t, err)
}

func TestAttestationAfterOracleUpdate(t *testing.T) {
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 25)
	keeper := myApp.BscKeeper
	var err error

	for i := 0; i < 20; i++ {
		_, err = h(ctx, &types.MsgCreateOracleBridger{
			OracleAddress:   oracleAddressList[i].String(),
			BridgeAddress:   orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
			ChainName:       chainName,
		})
		require.NoError(t, err)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		oracleSets := keeper.GetOracleSets(ctx)
		require.NotNil(t, oracleSets)
		require.EqualValues(t, i+1, len(oracleSets))

		power := keeper.GetLastTotalPower(ctx)
		expectPower := minStakeAmount.Mul(sdk.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(t, expectPower.Equal(power))
	}

	{
		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
			Name:           "Pundix Reward Token",
			Symbol:         "PURES",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      chainName,
		}

		for i := 0; i < 13; i++ {
			firstBridgeTokenClaim.BridgerAddress = orchestratorAddressList[i].String()
			_, err = h(ctx, firstBridgeTokenClaim)
			require.NoError(t, err)
			endBlockBeforeAttestation := keeper.GetAttestation(ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(t, endBlockBeforeAttestation)
			require.False(t, endBlockBeforeAttestation.Observed)
			require.NotNil(t, endBlockBeforeAttestation.Votes)
			require.EqualValues(t, i+1, len(endBlockBeforeAttestation.Votes))

			myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			endBlockAfterAttestation := keeper.GetAttestation(ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(t, endBlockAfterAttestation)
			require.False(t, endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.BridgerAddress = orchestratorAddressList[13].String()
		_, err = h(ctx, firstBridgeTokenClaim)
		require.NoError(t, err)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		attestation := keeper.GetAttestation(ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())

		require.NotNil(t, attestation)

		require.True(t, attestation.Observed)
		//t.Logf("attestation votes:%s", attestation.Votes)
	}

	{
		secondBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     2,
			BlockHeight:    1001,
			TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
			Name:           "Pundix Reward Token2",
			Symbol:         "PURES2",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      chainName,
		}

		for i := 0; i < 6; i++ {
			secondBridgeTokenClaim.BridgerAddress = orchestratorAddressList[i].String()
			_, err = h(ctx, secondBridgeTokenClaim)
			require.NoError(t, err)
			endBlockBeforeAttestation := keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(t, endBlockBeforeAttestation)
			require.False(t, endBlockBeforeAttestation.Observed)
			require.NotNil(t, endBlockBeforeAttestation.Votes)
			require.EqualValues(t, i+1, len(endBlockBeforeAttestation.Votes))

			myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			endBlockAfterAttestation := keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(t, endBlockAfterAttestation)
			require.False(t, endBlockAfterAttestation.Observed)
		}

		secondClaimAttestation := keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(t, secondClaimAttestation)
		require.False(t, secondClaimAttestation.Observed)
		require.NotNil(t, secondClaimAttestation.Votes)
		require.EqualValues(t, 6, len(secondClaimAttestation.Votes))

		proposalHandler := crosschain.NewCrossChainProposalHandler(myApp.CrosschainKeeper)

		var newOralceList []string
		for i := 0; i < 15; i++ {
			newOralceList = append(newOralceList, oracleAddressList[i].String())
		}
		err = proposalHandler(ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 1: try update chain oracle save top 15 oracle, expect success",
			Description: "",
			Oracles:     newOralceList,
			ChainName:   chainName,
		})
		require.NoError(t, err)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

		secondClaimAttestation = keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(t, secondClaimAttestation)
		require.False(t, secondClaimAttestation.Observed)
		require.NotNil(t, secondClaimAttestation.Votes)
		require.EqualValues(t, 6, len(secondClaimAttestation.Votes))

		activeOracles := keeper.GetAllActiveOracles(ctx)
		require.NotNil(t, activeOracles)
		require.EqualValues(t, 15, len(activeOracles))
		for i := 0; i < 15; i++ {
			require.NotNil(t, newOralceList[i], activeOracles[i].OracleAddress)
		}

		var newOracleList2 []string
		for i := 0; i < 11; i++ {
			newOracleList2 = append(newOracleList2, oracleAddressList[i].String())
		}
		err = proposalHandler(ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 2: try update chain oracle save top 11 oracle, expect success",
			Description: "",
			Oracles:     newOracleList2,
			ChainName:   chainName,
		})
		require.NoError(t, err)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

		secondClaimAttestation = keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(t, secondClaimAttestation)
		require.False(t, secondClaimAttestation.Observed)
		require.NotNil(t, secondClaimAttestation.Votes)
		require.EqualValues(t, 6, len(secondClaimAttestation.Votes))

		activeOracles = keeper.GetAllActiveOracles(ctx)
		require.NotNil(t, activeOracles)
		require.EqualValues(t, 11, len(activeOracles))
		for i := 0; i < 11; i++ {
			require.NotNil(t, newOracleList2[i], activeOracles[i].OracleAddress)
		}

		var newOracleList3 []string
		for i := 0; i < 10; i++ {
			newOracleList3 = append(newOracleList3, oracleAddressList[i].String())
		}
		err = proposalHandler(ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 3: try update chain oracle save top 10 oracle, expect success",
			Description: "",
			Oracles:     newOracleList3,
			ChainName:   chainName,
		})
		require.NoError(t, err)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

		secondClaimAttestation = keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(t, secondClaimAttestation)
		require.False(t, secondClaimAttestation.Observed)
		require.NotNil(t, secondClaimAttestation.Votes)
		require.EqualValues(t, 6, len(secondClaimAttestation.Votes))

		activeOracles = keeper.GetAllActiveOracles(ctx)
		require.NotNil(t, activeOracles)
		require.EqualValues(t, 10, len(activeOracles))
		for i := 0; i < 10; i++ {
			require.NotNil(t, newOracleList3[i], activeOracles[i].OracleAddress)
		}

		secondBridgeTokenClaim.BridgerAddress = orchestratorAddressList[6].String()
		_, err = h(ctx, secondBridgeTokenClaim)
		require.NoError(t, err)

		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

		secondClaimAttestation = keeper.GetAttestation(ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(t, secondClaimAttestation)
		require.True(t, secondClaimAttestation.Observed)
		require.NotNil(t, secondClaimAttestation.Votes)
		require.EqualValues(t, 7, len(secondClaimAttestation.Votes))
	}
}

func TestOracleDelete(t *testing.T) {
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 25)
	keeper := myApp.BscKeeper
	var err error

	for i := 0; i < 10; i++ {
		_, err = h(ctx, &types.MsgCreateOracleBridger{
			OracleAddress:   oracleAddressList[i].String(),
			BridgeAddress:   orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
			ChainName:       chainName,
		})
		require.NoError(t, err)
	}
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	allOracles := keeper.GetAllOracles(ctx)
	require.NotNil(t, allOracles)
	require.EqualValues(t, 10, len(allOracles))

	oracle := oracleAddressList[0]
	orchestrator := orchestratorAddressList[0]
	externalAddress := crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex()

	oracleAddr, found := keeper.GetOracleAddressByOrchestratorKey(ctx, orchestrator)
	require.True(t, found)
	require.EqualValues(t, oracle.String(), oracleAddr.String())

	oracleAddr, found = keeper.GetOracleByExternalAddress(ctx, externalAddress)
	require.True(t, found)
	require.EqualValues(t, oracle.String(), oracleAddr.String())

	oracleData, found := keeper.GetOracle(ctx, oracle)
	require.True(t, found)
	require.NotNil(t, oracleData)
	require.EqualValues(t, oracle.String(), oracleData.OracleAddress)
	require.EqualValues(t, orchestrator.String(), oracleData.BridgerAddress)
	require.EqualValues(t, externalAddress, oracleData.ExternalAddress)

	require.EqualValues(t, fxtypes.DefaultDenom, oracleData.DelegateAmount.Denom)
	require.True(t, minStakeAmount.Equal(oracleData.DelegateAmount.Amount))

	proposalHandler := crosschain.NewCrossChainProposalHandler(myApp.CrosschainKeeper)

	var newOracleAddressList []string
	for _, address := range oracleAddressList[1:] {
		newOracleAddressList = append(newOracleAddressList, address.String())
	}

	err = proposalHandler(ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 1: try update chain oracle remove first oracle, expect success",
		Description: "",
		Oracles:     newOracleAddressList,
		ChainName:   chainName,
	})
	require.NoError(t, err)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})

	oracleAddr, found = keeper.GetOracleAddressByOrchestratorKey(ctx, orchestrator)
	require.False(t, found)
	require.Nil(t, oracleAddr)

	oracleAddr, found = keeper.GetOracleByExternalAddress(ctx, externalAddress)
	require.False(t, found)
	require.Nil(t, oracleAddr)

	oracleData, found = keeper.GetOracle(ctx, oracle)
	require.False(t, found)
	require.EqualValues(t, types.Oracle{}, oracleData)
	require.EqualValues(t, "", oracleData.OracleAddress)
	require.EqualValues(t, "", oracleData.DelegateAmount.Denom)
	require.True(t, oracleData.DelegateAmount.Amount.IsNil())
}

func TestOracleSetSlash(t *testing.T) {
	//myApp.SetAppLog(server.ZeroLogWrapper{Logger: log.Logger.Level(zerolog.DebugLevel)})
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 10)
	keeper := myApp.BscKeeper
	var err error

	for i := 0; i < 10; i++ {
		_, err = h(ctx, &types.MsgCreateOracleBridger{
			OracleAddress:   oracleAddressList[i].String(),
			BridgeAddress:   orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
			ChainName:       chainName,
		})
		require.NoError(t, err)
	}
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	crosschain.EndBlocker(ctx, keeper)
	allOracles := keeper.GetAllOracles(ctx)
	require.NotNil(t, allOracles)
	require.EqualValues(t, 10, len(allOracles))
	oracleSets := keeper.GetOracleSets(ctx)
	require.NotNil(t, oracleSets)
	require.EqualValues(t, 1, len(oracleSets))

	gravityId := keeper.GetGravityID(ctx)
	checkpoint := oracleSets[0].GetCheckpoint(gravityId)
	for i := 0; i < 9; i++ {
		signature, err := types.NewEthereumSignature(checkpoint, ethKeys[i])
		require.NoError(t, err)
		_, err = h(ctx, &types.MsgOracleSetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			Signature:       hex.EncodeToString(signature),
			ChainName:       chainName,
		})
		require.NoError(t, err)
	}
	crosschain.EndBlocker(ctx, keeper)
	oracleSetHeight := int64(oracleSets[0].Height)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	myApp.EndBlock(abci.RequestEndBlock{Height: ctx.BlockHeight()})
	oracle, found := keeper.GetOracle(ctx, oracleAddressList[9])
	require.True(t, found)
	require.False(t, oracle.Jailed)

	ctx = ctx.WithBlockHeight(oracleSetHeight + int64(keeper.GetParams(ctx).SignedWindow) + 1)
	crosschain.EndBlocker(ctx, keeper)
	oracle, found = keeper.GetOracle(ctx, oracleAddressList[9])
	require.True(t, found)
	require.True(t, oracle.Jailed)
}

func TestSlashFactoryGreat1(t *testing.T) {
	//myApp.SetAppLog(server.ZeroLogWrapper{Logger: log.Logger.Level(zerolog.DebugLevel)})
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h := createTestEnv(t, 10)
	keeper := myApp.BscKeeper
	minStakeAmount, _ := sdk.NewIntFromString("11111111111111111111111")
	var err error

	for i := 0; i < 10; i++ {
		_, err = h(ctx, &types.MsgCreateOracleBridger{
			OracleAddress:   oracleAddressList[i].String(),
			BridgeAddress:   orchestratorAddressList[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(ethKeys[i].PublicKey).Hex(),
			DelegateAmount:  sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
			ChainName:       chainName,
		})
		require.NoError(t, err)
	}
	params := keeper.GetParams(ctx)
	params.SlashFraction, _ = sdk.NewDecFromStr("1.1")

	expectSlashAfterStakeAmount := sdk.MaxInt(
		// remainAmount = max (0, (depositAmount - slashAmount))
		minStakeAmount.Sub(
			sdk.MinInt(minStakeAmount, minStakeAmount.ToDec().Mul(params.SlashFraction).TruncateInt()),
		),
		sdk.ZeroInt())
	require.NotPanics(t, func() {
		keeper.SetParams(ctx, params)
	})

	require.NotPanics(t, func() {
		for i := 0; i < 10; i++ {
			oracle, found := keeper.GetOracle(ctx, oracleAddressList[i])
			require.True(t, found)
			require.False(t, oracle.Jailed)
			require.True(t, oracle.DelegateAmount.IsEqual(sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount}))

			keeper.SlashOracle(ctx, oracle, params.SlashFraction)

			oracle, found = keeper.GetOracle(ctx, oracleAddressList[i])
			require.True(t, found)
			require.True(t, oracle.Jailed)
			require.True(t, oracle.DelegateAmount.IsEqual(sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: expectSlashAfterStakeAmount}))
		}

		// repeat slash test.
		for i := 0; i < 10; i++ {
			oracle, found := keeper.GetOracle(ctx, oracleAddressList[i])
			require.True(t, found)
			require.True(t, oracle.Jailed)
			require.True(t, oracle.DelegateAmount.IsEqual(sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: expectSlashAfterStakeAmount}))

			keeper.SlashOracle(ctx, oracle, params.SlashFraction)

			oracle, found = keeper.GetOracle(ctx, oracleAddressList[i])
			require.True(t, found)
			require.True(t, oracle.Jailed)
			require.True(t, oracle.DelegateAmount.IsEqual(sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: expectSlashAfterStakeAmount}))
		}
	})
}
