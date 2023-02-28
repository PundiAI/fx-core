// nolint:staticcheck
package keeper_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	"github.com/functionx/fx-core/v3/x/gravity/keeper"
	v3 "github.com/functionx/fx-core/v3/x/gravity/migrations/v3"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

type MigrationTestSuite struct {
	suite.Suite

	app          *app.App
	ctx          sdk.Context
	migrator     keeper.Migrator
	msgServer    crosschaintypes.MsgServer
	bridgerAddrs []sdk.AccAddress
	externals    []*ecdsa.PrivateKey
	valAddrs     []sdk.ValAddress

	genesisState types.GenesisState
}

func TestMigrationTestSuite(t *testing.T) {
	fxtypes.SetConfig(false)
	suite.Run(t, new(MigrationTestSuite))
}

func (suite *MigrationTestSuite) SetupTest() {
	valNumber := tmrand.Intn(10) + 1

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight(),
		ChainID:         fxtypes.TestnetChainId,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})

	suite.migrator = keeper.NewMigrator(
		suite.app.AppCodec(),
		suite.app.LegacyAmino(),
		suite.app.GetKey(paramstypes.ModuleName),
		suite.app.GetKey(types.ModuleName),
		suite.app.GetKey(ethtypes.ModuleName),
		suite.app.StakingKeeper,
		suite.app.AccountKeeper,
		suite.app.BankKeeper,
	)
	suite.msgServer = crosschainkeeper.NewMsgServerImpl(suite.app.EthKeeper)

	for _, addr := range v3.GetEthOracleAddrs(suite.ctx.ChainID()) {
		helpers.AddTestAddr(suite.app, suite.ctx, sdk.MustAccAddressFromBech32(addr), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18))))
	}

	suite.bridgerAddrs = helpers.AddTestAddrs(suite.app, suite.ctx, valNumber, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100).MulRaw(1e18))))
	suite.externals = helpers.CreateMultiECDSA(valNumber)
	suite.Equal(len(suite.bridgerAddrs), len(valAccounts), valNumber)
	suite.valAddrs = make([]sdk.ValAddress, len(suite.bridgerAddrs))
	for i, account := range valAccounts {
		suite.valAddrs[i] = account.GetAddress().Bytes()
	}
}

func (suite *MigrationTestSuite) InitGravityStore() {
	paramsStore := suite.ctx.MultiStore().GetKVStore(suite.app.GetKey(paramstypes.ModuleName))
	gravityStore := suite.ctx.MultiStore().GetKVStore(suite.app.GetKey(types.ModuleName))
	v3.InitTestGravityDB(suite.app.AppCodec(), suite.app.LegacyAmino(), suite.genesisState, paramsStore, gravityStore)
}

func (suite *MigrationTestSuite) createDefGenesisState() {
	suite.genesisState = types.GenesisState{
		Params:            v3.TestParams(),
		LastObservedNonce: tmrand.Uint64(),
		Erc20ToDenoms: []types.ERC20ToDenom{
			{
				Erc20: helpers.GenerateAddress().Hex(),
				Denom: fxtypes.DefaultDenom,
			},
		},
	}
	var votes []string
	for i, addr := range suite.bridgerAddrs {
		_, found := suite.app.StakingKeeper.GetValidator(suite.ctx, suite.valAddrs[i])
		suite.True(found)
		suite.genesisState.DelegateKeys = append(suite.genesisState.DelegateKeys, types.MsgSetOrchestratorAddress{
			Validator:    suite.valAddrs[i].String(),
			Orchestrator: addr.String(),
			EthAddress:   crypto.PubkeyToAddress(suite.externals[i].PublicKey).String(),
		})
		votes = append(votes, suite.valAddrs[i].String())
	}

	suite.genesisState.Attestations = []types.Attestation{
		{
			Observed: true,
			Votes:    votes,
			Height:   tmrand.Uint64(),
			Claim: v3.AttClaimToAny(&types.MsgDepositClaim{
				EventNonce:    suite.genesisState.LastObservedNonce,
				BlockHeight:   tmrand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				Amount:        sdkmath.NewInt(tmrand.Int63() + 1),
				EthSender:     helpers.GenerateAddress().Hex(),
				FxReceiver:    sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				TargetIbc:     "",
				Orchestrator:  suite.bridgerAddrs[0].String(),
			}),
		},
	}
}

func (suite *MigrationTestSuite) checkBridgeToken(tokenContract string, bridgeTokenLen int) {
	response, err := suite.app.EthKeeper.BridgeTokens(sdk.WrapSDKContext(suite.ctx),
		&crosschaintypes.QueryBridgeTokensRequest{})
	suite.NoError(err)
	suite.Equal(len(response.BridgeTokens), bridgeTokenLen)
	suite.Contains(response.BridgeTokens, &crosschaintypes.BridgeToken{
		Token: tokenContract,
		Denom: fmt.Sprintf("%s%s", ethtypes.ModuleName, tokenContract),
	})
}

func (suite *MigrationTestSuite) TestBridgeTokenClaim() {
	// MsgBridgeTokenClaim
	suite.createDefGenesisState()

	tokenContract := helpers.GenerateAddress().Hex()
	metadata := fxtypes.GetCrossChainMetadata("Test Token", "TEST", uint32(tmrand.Intn(18)),
		fmt.Sprintf("%s%s", ethtypes.ModuleName, tokenContract))
	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)

	suite.InitGravityStore()
	suite.Equal(suite.app.EthKeeper.GetAllOracles(suite.ctx, false).Len(), 0)
	suite.NoError(suite.migrator.Migrate1to2(suite.ctx))

	suite.Equal(suite.app.EthKeeper.GetAllOracles(suite.ctx, false).Len(), len(suite.bridgerAddrs))

	power := suite.app.EthKeeper.GetLastTotalPower(suite.ctx)
	onlineOracles := suite.app.EthKeeper.GetAllOracles(suite.ctx, true)
	suite.True(onlineOracles.Len() > 0 && onlineOracles.Len() <= 20)
	suite.Equal(power.String(), sdkmath.NewInt(int64(100*len(onlineOracles))).String())

	proposalOracle, found := suite.app.EthKeeper.GetProposalOracle(suite.ctx)
	suite.True(found)
	suite.Equal(len(proposalOracle.Oracles), len(suite.bridgerAddrs))

	suite.checkBridgeToken(tokenContract, len(suite.genesisState.Erc20ToDenoms)+1)

	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:    suite.genesisState.LastObservedNonce + 1,
		BlockHeight:   tmrand.Uint64(),
		TokenContract: helpers.GenerateAddress().Hex(),
		Name:          "Test token 2",
		Symbol:        "TEST2",
		Decimals:      uint64(tmrand.Intn(18) + 1),
		ChannelIbc:    "",
		ChainName:     ethtypes.ModuleName,
	}

	for _, onlineOracle := range onlineOracles {
		lastEventNonce := suite.app.EthKeeper.GetLastEventNonceByOracle(suite.ctx, onlineOracle.GetOracle())
		suite.Require().Equal(lastEventNonce, suite.genesisState.LastObservedNonce)

		msg.BridgerAddress = onlineOracle.BridgerAddress
		_, err := suite.msgServer.BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), msg)
		suite.NoError(err)
	}

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()

	suite.checkBridgeToken(tokenContract, len(suite.genesisState.Erc20ToDenoms)+2)
}

func (suite *MigrationTestSuite) TestSendToFxClaim() {
	// MsgSendToFxClaim
	suite.createDefGenesisState()

	suite.InitGravityStore()
	suite.NoError(suite.migrator.Migrate1to2(suite.ctx))

	onlineOracles := suite.app.EthKeeper.GetAllOracles(suite.ctx, true)
	suite.True(onlineOracles.Len() > 0 && onlineOracles.Len() <= 20)

	bridgeToken := suite.app.EthKeeper.GetDenomByBridgeToken(suite.ctx, suite.genesisState.Erc20ToDenoms[0].Denom)
	suite.Equal(bridgeToken.Token, suite.genesisState.Erc20ToDenoms[0].Erc20)
	suite.Equal(bridgeToken.Denom, fxtypes.DefaultDenom)

	msg := &crosschaintypes.MsgSendToFxClaim{
		EventNonce:    suite.genesisState.LastObservedNonce,
		BlockHeight:   tmrand.Uint64(),
		TokenContract: bridgeToken.Token,
		Amount:        sdkmath.NewInt(1),
		Sender:        helpers.GenerateAddress().Hex(),
		Receiver:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		TargetIbc:     "",
		ChainName:     ethtypes.ModuleName,
	}

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(msg.Receiver))
	suite.True(balances.IsZero())

	for _, oracle := range onlineOracles {
		msg.EventNonce = suite.genesisState.LastObservedNonce + 1
		msg.BridgerAddress = oracle.BridgerAddress
		_, err := suite.msgServer.SendToFxClaim(sdk.WrapSDKContext(suite.ctx), msg)
		suite.NoError(err)
	}

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(msg.Receiver))
	suite.Equal(balances.String(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, msg.Amount)).String())
}

func (suite *MigrationTestSuite) TestSendToExternal() {
	/*
		1.MsgSendToExternal
			1.1.MsgCancelSendToExternal
		2.MsgRequestBatch
		3.MsgConfirmBatch
		4.MsgSendToExternalClaim
	*/
	suite.createDefGenesisState()
	suite.genesisState.LastObservedBlockHeight = types.LastObservedEthereumBlockHeight{
		FxBlockHeight:  uint64(suite.ctx.BlockHeight()),
		EthBlockHeight: tmrand.Uint64(),
	}

	tokenContract := helpers.GenerateAddress().Hex()
	denom := fmt.Sprintf("%s%s", ethtypes.ModuleName, tokenContract)
	metadata := fxtypes.GetCrossChainMetadata("Test Token", "TEST", uint32(tmrand.Intn(18)), denom)
	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)

	err := suite.app.BankKeeper.MintCoins(suite.ctx,
		ethtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(200))))
	suite.NoError(err)

	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx,
		ethtypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(200))))
	suite.NoError(err)

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx,
		types.ModuleName, suite.valAddrs[0].Bytes(), sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(100))))
	suite.NoError(err)

	suite.InitGravityStore()
	suite.NoError(suite.migrator.Migrate1to2(suite.ctx))

	bridgeToken := suite.app.EthKeeper.GetBridgeTokenDenom(suite.ctx, tokenContract)

	_, err = suite.msgServer.SendToExternal(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.MsgSendToExternal{
		Sender:    sdk.AccAddress(suite.valAddrs[0]).String(),
		Dest:      helpers.GenerateAddress().Hex(),
		Amount:    sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(1)),
		BridgeFee: sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(1)),
		ChainName: ethtypes.ModuleName,
	})
	suite.NoError(err)

	balances := suite.app.BankKeeper.GetBalance(suite.ctx, suite.valAddrs[0].Bytes(), denom)
	suite.Equal(balances, sdk.NewCoin(denom, sdkmath.NewInt(98)))

	response, err := suite.app.EthKeeper.BatchFees(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.QueryBatchFeeRequest{})
	suite.NoError(err)
	suite.Equal(len(response.BatchFees), 1)
	suite.Equal(response.BatchFees[0].TokenContract, tokenContract)

	_, err = suite.msgServer.RequestBatch(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.MsgRequestBatch{
		Sender:     suite.bridgerAddrs[0].String(),
		Denom:      denom,
		MinimumFee: sdkmath.ZeroInt(),
		FeeReceive: helpers.GenerateAddress().Hex(),
		ChainName:  ethtypes.ModuleName,
		BaseFee:    sdkmath.ZeroInt(),
	})
	suite.NoError(err)

	txBatchesResponse, err := suite.app.EthKeeper.OutgoingTxBatches(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.QueryOutgoingTxBatchesRequest{})
	suite.NoError(err)
	suite.Equal(len(txBatchesResponse.Batches), 1)
	suite.Equal(txBatchesResponse.Batches[0].TokenContract, tokenContract)

	checkpoint, err := txBatchesResponse.Batches[0].GetCheckpoint(suite.genesisState.Params.GravityId)
	suite.NoError(err)

	onlineOracles := suite.app.EthKeeper.GetAllOracles(suite.ctx, true)
	suite.True(onlineOracles.Len() > 0 && onlineOracles.Len() <= 20)

	for i, bridger := range suite.bridgerAddrs {
		var onlineOracle crosschaintypes.Oracle
		for _, oracle := range onlineOracles {
			if oracle.BridgerAddress == bridger.String() {
				onlineOracle = oracle
				break
			}
		}
		if !onlineOracle.Online {
			continue
		}
		externalAddress := crypto.PubkeyToAddress(suite.externals[i].PublicKey).String()
		suite.Equal(externalAddress, onlineOracle.ExternalAddress)
		signature, err := crosschaintypes.NewEthereumSignature(checkpoint, suite.externals[i])
		suite.NoError(err)

		_, err = suite.msgServer.ConfirmBatch(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.MsgConfirmBatch{
			Nonce:           txBatchesResponse.Batches[0].BatchNonce,
			TokenContract:   txBatchesResponse.Batches[0].TokenContract,
			BridgerAddress:  bridger.String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       ethtypes.ModuleName,
		})
		suite.NoError(err)
	}

	msg := &crosschaintypes.MsgSendToExternalClaim{
		EventNonce:    suite.genesisState.LastObservedNonce + 1,
		BlockHeight:   suite.genesisState.LastObservedBlockHeight.EthBlockHeight + 1,
		BatchNonce:    txBatchesResponse.Batches[0].BatchNonce,
		TokenContract: txBatchesResponse.Batches[0].TokenContract,
		ChainName:     ethtypes.ModuleName,
	}

	for _, onlineOracle := range onlineOracles {
		msg.BridgerAddress = onlineOracle.BridgerAddress
		_, err := suite.msgServer.SendToExternalClaim(sdk.WrapSDKContext(suite.ctx), msg)
		suite.NoError(err)
	}
}

func (suite *MigrationTestSuite) TestOracleSetConfirm() {
	/*
		1. MsgOracleSetConfirm
		2. MsgOracleSetUpdatedClaim
	*/

	suite.createDefGenesisState()
	suite.genesisState.LastObservedBlockHeight = types.LastObservedEthereumBlockHeight{
		FxBlockHeight:  uint64(suite.ctx.BlockHeight()),
		EthBlockHeight: tmrand.Uint64(),
	}

	suite.InitGravityStore()
	suite.NoError(suite.migrator.Migrate1to2(suite.ctx))

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.app.Commit()

	onlineOracles := suite.app.EthKeeper.GetAllOracles(suite.ctx, true)
	suite.True(onlineOracles.Len() > 0 && onlineOracles.Len() <= 20)

	oracleSet := suite.app.EthKeeper.GetLatestOracleSet(suite.ctx)
	suite.True(oracleSet.Height > 0)
	suite.True(oracleSet.Nonce > 0)
	suite.True(len(oracleSet.Members) > 0)
	checkpoint, err := oracleSet.GetCheckpoint(suite.genesisState.Params.GravityId)
	suite.NoError(err)

	for i, bridger := range suite.bridgerAddrs {
		var onlineOracle crosschaintypes.Oracle
		for _, oracle := range onlineOracles {
			if oracle.BridgerAddress == bridger.String() {
				onlineOracle = oracle
				break
			}
		}
		if !onlineOracle.Online {
			continue
		}
		externalAddress := crypto.PubkeyToAddress(suite.externals[i].PublicKey).String()
		suite.Equal(externalAddress, onlineOracle.ExternalAddress)
		signature, err := crosschaintypes.NewEthereumSignature(checkpoint, suite.externals[i])
		suite.NoError(err)
		_, err = suite.msgServer.OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.MsgOracleSetConfirm{
			Nonce:           oracleSet.Nonce,
			BridgerAddress:  onlineOracle.BridgerAddress,
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       ethtypes.ModuleName,
		})
		suite.NoError(err)
	}

	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     suite.genesisState.LastObservedNonce + 1,
		BlockHeight:    suite.genesisState.LastObservedBlockHeight.EthBlockHeight + 1,
		OracleSetNonce: oracleSet.Nonce,
		Members:        oracleSet.Members,
		ChainName:      ethtypes.ModuleName,
	}
	for _, onlineOracle := range onlineOracles {
		msg.BridgerAddress = onlineOracle.BridgerAddress
		_, err := suite.msgServer.OracleSetUpdateClaim(sdk.WrapSDKContext(suite.ctx), msg)
		suite.NoError(err)
	}
}
