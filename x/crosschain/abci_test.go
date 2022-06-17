package crosschain_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"reflect"
	"strings"
	"sync"
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

type IntegrationTestSuite struct {
	suite.Suite
	sync.Mutex

	app            *app.App
	ctx            sdk.Context
	oracles        []sdk.AccAddress
	bridgers       []sdk.AccAddress
	externals      []*ecdsa.PrivateKey
	validator      []sdk.ValAddress
	chainName      string
	delegateAmount sdk.Int
}

func TestIntegrationTestSuite(t *testing.T) {
	methodFinder := reflect.TypeOf(new(IntegrationTestSuite))
	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)
		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}
		t.Run(method.Name, func(subT *testing.T) {
			mySuite := new(IntegrationTestSuite)
			mySuite.SetT(t)
			mySuite.SetupTest()
			method.Func.Call([]reflect.Value{reflect.ValueOf(mySuite)})
		})
	}
}

func (suite *IntegrationTestSuite) Msg() types.MsgServer {
	return keeper.NewMsgServerImpl(suite.Keeper())
}

func (suite *IntegrationTestSuite) Keeper() keeper.Keeper {
	return suite.app.BscKeeper
}

func (suite *IntegrationTestSuite) SetupTest() {
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(types.MaxOracleSize, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.externals = genEthKey(types.MaxOracleSize)
	suite.delegateAmount = sdk.NewInt(10 * 1e3).MulRaw(1e18)
	for i := 0; i < types.MaxOracleSize; i++ {
		suite.validator = append(suite.validator, valAccounts[i].GetAddress().Bytes())
	}
	suite.chainName = "bsc"

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracles {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.ctx, proposalOracle)
}

func (suite *IntegrationTestSuite) TestABCIEndBlockDepositClaim() {

	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err := suite.Msg().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Name:           "Pundix Reward Token",
		Symbol:         "PURES",
		Decimals:       18,
		BridgerAddress: suite.bridgers[0].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err = suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	sendToFxClaim := &types.MsgSendToFxClaim{
		EventNonce:     2,
		BlockHeight:    1001,
		TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Amount:         sdk.NewInt(1234),
		Sender:         "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Receiver:       "fx16wvwsmpp4y4ttgzknyr6kqla877jud6u04lqey",
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: suite.bridgers[0].String(),
		ChainName:      suite.chainName,
	}
	_, err = suite.Msg().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToFxClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	receiveAddr, err := sdk.AccAddressFromBech32(sendToFxClaim.Receiver)
	require.NoError(suite.T(), err)
	allBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, receiveAddr)
	tokenContract := common.HexToAddress(addBridgeTokenClaim.TokenContract).Hex()
	// transfer/channel-0/bscPURES
	tokenName := fmt.Sprintf("%s%s", suite.chainName, tokenContract)
	if len(addBridgeTokenClaim.ChannelIbc) > 0 {
		channel, err := hex.DecodeString(addBridgeTokenClaim.ChannelIbc)
		require.NoError(suite.T(), err)
		tokenName = ibcTransferTypes.DenomTrace{
			Path:      string(channel),
			BaseDenom: fmt.Sprintf("%s%s", suite.chainName, tokenContract),
		}.IBCDenom()
	}
	require.EqualValues(suite.T(), fmt.Sprintf("%s%s", sendToFxClaim.Amount.String(), tokenName), allBalances.String())
}

func (suite *IntegrationTestSuite) TestOracleUpdate() {
	for i := 0; i < 10; i++ {
		_, err := suite.Msg().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
			ChainName:        suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := suite.delegateAmount.Mul(sdk.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
			Name:           "Pundix Reward Token",
			Symbol:         "PURES",
			Decimals:       18,
			BridgerAddress: suite.bridgers[i].String(),
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}
		_, err := suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
		require.NoError(suite.T(), err)
		endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), endBlockBeforeAttestation)
		require.False(suite.T(), endBlockBeforeAttestation.Observed)
		require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
		require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), endBlockAfterAttestation)
		require.False(suite.T(), endBlockAfterAttestation.Observed)
	}

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
		Name:           "Pundix Reward Token",
		Symbol:         "PURES",
		Decimals:       18,
		BridgerAddress: suite.bridgers[6].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err := suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	attestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())

	require.NotNil(suite.T(), attestation)

	require.True(suite.T(), attestation.Observed)

	proposalHandler := crosschain.NewChainProposalHandler(suite.app.CrosschainKeeper)

	var newOracleList []string
	for i := 0; i < 7; i++ {
		newOracleList = append(newOracleList, suite.oracles[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 1: try update chain oracle power >= 30%, expect error",
		Description: "",
		Oracles:     newOracleList,
		ChainName:   suite.chainName,
	})
	require.ErrorIs(suite.T(), types.ErrInvalid, err)

	expectTotalPower := suite.delegateAmount.Mul(sdk.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
	require.True(suite.T(), expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalOracleChangePowerThreshold.Mul(expectTotalPower).Quo(sdk.NewInt(100))

	expectDeletePower := suite.delegateAmount.Mul(sdk.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	require.EqualValues(suite.T(), fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newOracleList2 []string
	for i := 0; i < 8; i++ {
		newOracleList2 = append(newOracleList2, suite.oracles[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 2: try update chain oracle power <= 30%, expect success",
		Description: "",
		Oracles:     newOracleList2,
		ChainName:   suite.chainName,
	})
	require.NoError(suite.T(), err)
}

func (suite *IntegrationTestSuite) TestAttestationAfterOracleUpdate() {

	for i := 0; i < 20; i++ {
		_, err := suite.Msg().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
			ChainName:        suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := suite.delegateAmount.Mul(sdk.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
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
			ChainName:      suite.chainName,
		}

		for i := 0; i < 13; i++ {
			firstBridgeTokenClaim.BridgerAddress = suite.bridgers[i].String()
			_, err := suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
			suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.BridgerAddress = suite.bridgers[13].String()
		_, err := suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		attestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())

		require.NotNil(suite.T(), attestation)

		require.True(suite.T(), attestation.Observed)
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
			ChainName:      suite.chainName,
		}

		for i := 0; i < 6; i++ {
			secondBridgeTokenClaim.BridgerAddress = suite.bridgers[i].String()
			_, err := suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), secondBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
			suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		secondClaimAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		proposalHandler := crosschain.NewChainProposalHandler(suite.app.CrosschainKeeper)

		var newOracleList []string
		for i := 0; i < 15; i++ {
			newOracleList = append(newOracleList, suite.oracles[i].String())
		}
		err := proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 1: try update chain oracle save top 15 oracle, expect success",
			Description: "",
			Oracles:     newOracleList,
			ChainName:   suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeOracles := suite.Keeper().GetAllOracles(suite.ctx, true)
		require.NotNil(suite.T(), activeOracles)
		require.EqualValues(suite.T(), 15, len(activeOracles))
		for i := 0; i < 15; i++ {
			require.NotNil(suite.T(), newOracleList[i], activeOracles[i].OracleAddress)
		}

		var newOracleList2 []string
		for i := 0; i < 11; i++ {
			newOracleList2 = append(newOracleList2, suite.oracles[i].String())
		}
		err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 2: try update chain oracle save top 11 oracle, expect success",
			Description: "",
			Oracles:     newOracleList2,
			ChainName:   suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeOracles = suite.Keeper().GetAllOracles(suite.ctx, true)
		require.NotNil(suite.T(), activeOracles)
		require.EqualValues(suite.T(), 11, len(activeOracles))
		for i := 0; i < 11; i++ {
			require.NotNil(suite.T(), newOracleList2[i], activeOracles[i].OracleAddress)
		}

		var newOracleList3 []string
		for i := 0; i < 10; i++ {
			newOracleList3 = append(newOracleList3, suite.oracles[i].String())
		}
		err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 3: try update chain oracle save top 10 oracle, expect success",
			Description: "",
			Oracles:     newOracleList3,
			ChainName:   suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeOracles = suite.Keeper().GetAllOracles(suite.ctx, true)
		require.NotNil(suite.T(), activeOracles)
		require.EqualValues(suite.T(), 10, len(activeOracles))
		for i := 0; i < 10; i++ {
			require.NotNil(suite.T(), newOracleList3[i], activeOracles[i].OracleAddress)
		}

		secondBridgeTokenClaim.BridgerAddress = suite.bridgers[6].String()
		_, err = suite.Msg().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), secondBridgeTokenClaim)
		require.NoError(suite.T(), err)

		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.True(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 7, len(secondClaimAttestation.Votes))
	}
}

func (suite *IntegrationTestSuite) TestOracleDelete() {
	for i := 0; i < 10; i++ {
		_, err := suite.Msg().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
			ChainName:        suite.chainName,
		})
		require.NoError(suite.T(), err)
	}
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	allOracles := suite.Keeper().GetAllOracles(suite.ctx, false)
	require.NotNil(suite.T(), allOracles)
	require.EqualValues(suite.T(), 10, len(allOracles))

	oracle := suite.oracles[0]
	bridger := suite.bridgers[0]
	externalAddress := crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex()

	oracleAddr, found := suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleAddr, found = suite.Keeper().GetOracleByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleData, found := suite.Keeper().GetOracle(suite.ctx, oracle)
	require.True(suite.T(), found)
	require.NotNil(suite.T(), oracleData)
	require.EqualValues(suite.T(), oracle.String(), oracleData.OracleAddress)
	require.EqualValues(suite.T(), bridger.String(), oracleData.BridgerAddress)
	require.EqualValues(suite.T(), externalAddress, oracleData.ExternalAddress)

	require.True(suite.T(), suite.delegateAmount.Equal(oracleData.DelegateAmount))

	proposalHandler := crosschain.NewChainProposalHandler(suite.app.CrosschainKeeper)

	var newOracleAddressList []string
	for _, address := range suite.oracles[1:] {
		newOracleAddressList = append(newOracleAddressList, address.String())
	}

	err := proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 1: try update chain oracle remove first oracle, expect success",
		Description: "",
		Oracles:     newOracleAddressList,
		ChainName:   suite.chainName,
	})
	require.NoError(suite.T(), err)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracleAddr, found = suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)

	oracleAddr, found = suite.Keeper().GetOracleByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)

	oracleData, found = suite.Keeper().GetOracle(suite.ctx, oracle)
	require.True(suite.T(), found)
}

func (suite *IntegrationTestSuite) TestOracleSetSlash() {
	for i := 0; i < 10; i++ {
		_, err := suite.Msg().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
			ChainName:        suite.chainName,
		})
		require.NoError(suite.T(), err)
	}
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	crosschain.EndBlocker(suite.ctx, suite.Keeper())
	allOracles := suite.Keeper().GetAllOracles(suite.ctx, false)
	require.NotNil(suite.T(), allOracles)
	require.EqualValues(suite.T(), 10, len(allOracles))
	oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
	require.NotNil(suite.T(), oracleSets)
	require.EqualValues(suite.T(), 1, len(oracleSets))

	gravityId := suite.Keeper().GetGravityID(suite.ctx)
	checkpoint, _ := oracleSets[0].GetCheckpoint(gravityId)
	for i := 0; i < 9; i++ {
		signature, err := types.NewEthereumSignature(checkpoint, suite.externals[i])
		require.NoError(suite.T(), err)
		_, err = suite.Msg().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), &types.MsgOracleSetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  suite.bridgers[i].String(),
			ExternalAddress: crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		})
		require.NoError(suite.T(), err)
	}
	crosschain.EndBlocker(suite.ctx, suite.Keeper())
	oracleSetHeight := int64(oracleSets[0].Height)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[9])
	require.True(suite.T(), found)
	require.True(suite.T(), oracle.Online)

	suite.ctx = suite.ctx.WithBlockHeight(oracleSetHeight + int64(suite.Keeper().GetParams(suite.ctx).SignedWindow) + 1)
	crosschain.EndBlocker(suite.ctx, suite.Keeper())
	oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracles[9])
	require.True(suite.T(), found)
	require.False(suite.T(), oracle.Online)
}

func (suite *IntegrationTestSuite) TestSlashFactoryGreat1() {
	for i := 0; i < 10; i++ {
		_, err := suite.Msg().BondedOracle(sdk.WrapSDKContext(suite.ctx), &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
			ChainName:        suite.chainName,
		})
		require.NoError(suite.T(), err)
	}
	params := suite.Keeper().GetParams(suite.ctx)
	params.SlashFraction, _ = sdk.NewDecFromStr("1.1")

	//expectSlashAfterStakeAmount := sdk.MaxInt(
	//	suite.delegateAmount.Sub(
	//		sdk.MinInt(suite.delegateAmount, suite.delegateAmount.ToDec().Mul(params.SlashFraction).TruncateInt()),
	//	),
	//	sdk.ZeroInt())
	require.NotPanics(suite.T(), func() {
		suite.Keeper().SetParams(suite.ctx, &params)
	})

	require.NotPanics(suite.T(), func() {
		for i := 0; i < 10; i++ {
			oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
			require.True(suite.T(), found)
			require.True(suite.T(), oracle.Online)
			require.True(suite.T(), oracle.DelegateAmount.Equal(suite.delegateAmount))

			suite.Keeper().SlashOracle(suite.ctx, oracle)

			oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
			require.True(suite.T(), found)
			require.False(suite.T(), oracle.Online)
			//require.True(suite.T(), oracle.DelegateAmount.Equal(expectSlashAfterStakeAmount))
		}

		// repeat slash test.
		for i := 0; i < 10; i++ {
			oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
			require.True(suite.T(), found)
			require.False(suite.T(), oracle.Online)
			//require.True(suite.T(), oracle.DelegateAmount.Equal(expectSlashAfterStakeAmount))

			suite.Keeper().SlashOracle(suite.ctx, oracle)

			oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
			require.True(suite.T(), found)
			//require.True(suite.T(), oracle.Online)
			//require.True(suite.T(), oracle.DelegateAmount.Equal(expectSlashAfterStakeAmount))
		}
	})
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
