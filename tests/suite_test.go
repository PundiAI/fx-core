package tests

import (
	"context"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	"github.com/functionx/fx-core/client/jsonrpc"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app/helpers"
	"github.com/functionx/fx-core/client/grpc"
	fxtypes "github.com/functionx/fx-core/types"
)

type TestSuite struct {
	suite.Suite
	sync.Mutex
	ctx        context.Context
	grpcClient *grpc.Client
	tmClient   jsonrpc.CustomRPC
	privateKey cryptotypes.PrivKey
}

func NewTestSuite() TestSuite {
	client, err := grpc.NewGRPCClient(GetGrpcUrl())
	if err != nil {
		panic(err)
	}
	mnemonic, keyType := GetAdminMnemonic()
	privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, keyType, 0, 0)
	if err != nil {
		panic(err)
	}
	return TestSuite{
		Suite:      suite.Suite{},
		ctx:        context.Background(),
		grpcClient: client,
		privateKey: privKey,
		tmClient:   jsonrpc.NewCustomRPC(jsonrpc.NewFastClient(GetNodeJsonRpcUrl())),
	}
}

func (suite *TestSuite) SetPrivateKey(key cryptotypes.PrivKey) {
	suite.privateKey = key
}

func (suite *TestSuite) AdminAddress() sdk.AccAddress {
	return suite.privateKey.PubKey().Address().Bytes()
}

func (suite *TestSuite) BroadcastTx(msgList ...sdk.Msg) string {
	buildTx, err := suite.grpcClient.BuildTx(suite.privateKey, msgList)
	suite.Require().NoError(err)

	txResponse, err := suite.grpcClient.BroadcastTxOk(buildTx)
	suite.Require().NoError(err)
	return txResponse.TxHash
}

func (suite *TestSuite) QueryBalance(accAddress sdk.AccAddress) sdk.Coins {
	balances, err := suite.grpcClient.QueryBalances(accAddress.String())
	suite.Require().NoError(err)
	suite.Require().True(balances.IsAllPositive())
	return balances
}

func (suite *TestSuite) QueryValidator() []sdk.ValAddress {
	validators, err := suite.grpcClient.StakingQuery().Validators(suite.ctx, &stakingtypes.QueryValidatorsRequest{Status: stakingtypes.Bonded.String()})
	suite.Require().NoError(err)
	suite.Require().True(len(validators.Validators) > 0)

	vals := make([]sdk.ValAddress, 0, len(validators.Validators))
	for _, v := range validators.Validators {
		valAddr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		suite.Require().NoError(err)
		vals = append(vals, valAddr)
		suite.T().Log("query validator", "val", v.OperatorAddress, "token", v.Tokens.String(), "share", v.DelegatorShares.String())
	}

	return vals
}

func (suite *TestSuite) CreateValidator(valPriv cryptotypes.PrivKey) {
	valAddr := sdk.ValAddress(valPriv.PubKey().Address())

	amt := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000))))
	txHash := suite.BroadcastTx(banktypes.NewMsgSend(suite.AdminAddress(), sdk.AccAddress(valAddr), amt))
	suite.T().Log("send to validator txHash", txHash)

	oldKey := suite.privateKey
	defer suite.SetPrivateKey(oldKey)
	suite.SetPrivateKey(valPriv)

	selfDelegate := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100)))
	minSelfDelegate := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1))
	description := stakingtypes.Description{
		Moniker:         "val2",
		Identity:        "",
		Website:         "",
		SecurityContact: "",
		Details:         "",
	}
	rates := stakingtypes.CommissionRates{
		Rate:          sdk.NewDecWithPrec(2, 2),
		MaxRate:       sdk.NewDecWithPrec(50, 2),
		MaxChangeRate: sdk.NewDecWithPrec(2, 2),
	}
	ed25519PrivKey := ed25519.GenPrivKeyFromSecret(valAddr.Bytes())
	msg, err := stakingtypes.NewMsgCreateValidator(valAddr, ed25519PrivKey.PubKey(), selfDelegate, description, rates, minSelfDelegate)
	suite.Require().NoError(err)
	txHash = suite.BroadcastTx(msg)
	suite.T().Log("create validator txHash", txHash)
}

func (suite *TestSuite) Send(accAddress sdk.AccAddress, amount sdk.Coin) {
	txHash := suite.BroadcastTx(banktypes.NewMsgSend(suite.AdminAddress(), accAddress, sdk.NewCoins(amount)))
	suite.T().Log("send txHash", txHash)
}

func (suite *TestSuite) SetWithdrawAddr(accAddress sdk.AccAddress) {
	txHash := suite.BroadcastTx(distritypes.NewMsgSetWithdrawAddress(suite.AdminAddress(), accAddress))
	suite.T().Log("set withdraw txHash", txHash)
}

func (suite *TestSuite) Delegate(valAddress sdk.ValAddress, amount sdk.Coin) {
	txHash := suite.BroadcastTx(stakingtypes.NewMsgDelegate(suite.AdminAddress(), valAddress, amount))
	suite.T().Log("delegate txHash", txHash)
}

func (suite *TestSuite) WithdrawReward(valAddress sdk.ValAddress) {
	txHash := suite.BroadcastTx(distritypes.NewMsgWithdrawDelegatorReward(suite.AdminAddress(), valAddress))
	suite.T().Log("withdraw reward txHash", txHash)
}

func (suite *TestSuite) Undelegate(valAddress sdk.ValAddress, all bool) {
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))
	if all {
		delegation, err := suite.grpcClient.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: suite.AdminAddress().String(),
			ValidatorAddr: valAddress.String(),
		})
		suite.Require().NoError(err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg := stakingtypes.NewMsgUndelegate(suite.AdminAddress(), valAddress, sdk.NewCoin(fxtypes.DefaultDenom, amt))
	txHash := suite.BroadcastTx(msg)
	suite.T().Log("undelegate txHash", txHash)
}

func (suite *TestSuite) Redelegate(valSrc, valDest sdk.ValAddress, all bool) {
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))
	if all {
		delegation, err := suite.grpcClient.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: suite.AdminAddress().String(),
			ValidatorAddr: valSrc.String(),
		})
		suite.Require().NoError(err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg := stakingtypes.NewMsgBeginRedelegate(suite.AdminAddress(), valSrc, valDest, sdk.NewCoin(fxtypes.DefaultDenom, amt))
	txHash := suite.BroadcastTx(msg)
	suite.T().Log("redelegate txHash", txHash)
}

func (suite *TestSuite) ProposalSubmit(deposit sdk.Coin) {
	content := govtypes.ContentFromProposalType("title", "description", "Text")
	msg, err := govtypes.NewMsgSubmitProposal(content, sdk.NewCoins(deposit), suite.AdminAddress())
	suite.Require().NoError(err)
	txHash := suite.BroadcastTx(msg)
	suite.T().Log("proposal submit txHash", txHash)
}

func (suite *TestSuite) ProposalDeposit(proposalID uint64, amount sdk.Coin) {
	txHash := suite.BroadcastTx(govtypes.NewMsgDeposit(suite.AdminAddress(), proposalID, sdk.NewCoins(amount)))
	suite.T().Log("proposal deposit txHash", txHash)
}

func (suite *TestSuite) ProposalVote(proposalID uint64, option govtypes.VoteOption) {
	txHash := suite.BroadcastTx(govtypes.NewMsgVote(suite.AdminAddress(), proposalID, option))
	suite.T().Log("proposal vote txHash", txHash)
}
