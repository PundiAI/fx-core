package tests

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	tmclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/client/grpc"
	"github.com/functionx/fx-core/v3/client/jsonrpc"
	"github.com/functionx/fx-core/v3/testutil"
	"github.com/functionx/fx-core/v3/testutil/helpers"
	"github.com/functionx/fx-core/v3/testutil/network"
	fxgovtypes "github.com/functionx/fx-core/v3/x/gov/types"
)

type TestSuite struct {
	suite.Suite
	useLocalNetwork bool
	network         *network.Network
	ctx             context.Context
	proposalId      uint64
}

func NewTestSuite() *TestSuite {
	testSuite := &TestSuite{
		Suite:           suite.Suite{},
		useLocalNetwork: true,
		proposalId:      0,
		ctx:             context.Background(),
	}
	if os.Getenv("USE_LOCAL_NETWORK") == "false" {
		testSuite.useLocalNetwork = false
	}
	return testSuite
}

func (suite *TestSuite) SetupSuite() {
	if !suite.useLocalNetwork {
		return
	}
	suite.T().Log("setting up integration test suite")

	encCfg := app.MakeEncodingConfig()
	ibcGenesisOpt := func(config *network.Config) {
		config.GenesisState = testutil.IbcGenesisState(encCfg.Codec, config.GenesisState)
	}
	bankGenesisOpt := func(config *network.Config) {
		config.GenesisState = testutil.BankGenesisState(encCfg.Codec, config.GenesisState)
	}
	cfg := testutil.DefaultNetworkConfig(encCfg, ibcGenesisOpt, bankGenesisOpt)
	cfg.TimeoutCommit = time.Millisecond
	cfg.NumValidators = 1
	// cfg.EnableTMLogging = true

	baseDir, err := os.MkdirTemp(suite.T().TempDir(), cfg.ChainID)
	suite.Require().NoError(err)
	suite.network, err = network.New(suite.T(), baseDir, cfg)
	suite.Require().NoError(err)
}

func (suite *TestSuite) TearDownSuite() {
	if !suite.useLocalNetwork {
		return
	}
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *TestSuite) WithNetwork(network *network.Network) {
	suite.network = network
}

func (suite *TestSuite) GetNetwork() *network.Network {
	return suite.network
}

func (suite *TestSuite) Context() context.Context {
	return suite.ctx
}

func (suite *TestSuite) IsUseLocalNetwork() bool {
	return suite.useLocalNetwork
}

func (suite *TestSuite) getNextProposalId() uint64 {
	suite.proposalId = suite.proposalId + 1
	return suite.proposalId
}

func (suite *TestSuite) GetFirstValidator() *network.Validator {
	return suite.network.Validators[0]
}

func (suite *TestSuite) GetFirstValPrivKey() cryptotypes.PrivKey {
	privKey, err := helpers.PrivKeyFromMnemonic(suite.network.Config.Mnemonics[0], hd.Secp256k1Type, 0, 0)
	suite.NoError(err)
	return privKey
}

func (suite *TestSuite) GRPCClient() *grpc.Client {
	grpcUrl := "http://localhost:9090"
	if suite.useLocalNetwork {
		grpcUrl = fmt.Sprintf("http://%s", suite.GetFirstValidator().AppConfig.GRPC.Address)
	}
	client, err := grpc.NewClient(grpcUrl, suite.ctx)
	suite.NoError(err)
	return client
}

func (suite *TestSuite) NodeClient() *jsonrpc.NodeRPC {
	nodeUrl := "http://localhost:26657"
	if suite.useLocalNetwork {
		nodeUrl = suite.GetFirstValidator().RPCAddress
	}
	return jsonrpc.NewNodeRPC(jsonrpc.NewFastClient(nodeUrl), suite.ctx)
}

func (suite *TestSuite) ValNodeClient() tmclient.Client {
	return suite.GetFirstValidator().RPCClient
}

func (suite *TestSuite) GetFirstValAddr() sdk.ValAddress {
	return suite.GetFirstValPrivKey().PubKey().Address().Bytes()
}

func (suite *TestSuite) GetStakingDenom() string {
	return suite.network.Config.BondDenom
}

func (suite *TestSuite) NewCoin(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(suite.GetStakingDenom(), amount)
}

func (suite *TestSuite) GetMetadata(denom string) banktypes.Metadata {
	response, err := suite.GRPCClient().BankQuery().DenomMetadata(suite.ctx, &banktypes.QueryDenomMetadataRequest{Denom: denom})
	suite.NoError(err)
	return response.Metadata
}

func (suite *TestSuite) BlockNumber() int64 {
	height, err := suite.GRPCClient().GetBlockHeight()
	suite.Error(err)
	return height
}

func (suite *TestSuite) QueryTx(txHash string) *sdk.TxResponse {
	txResponse, err := suite.GRPCClient().TxByHash(txHash)
	suite.NoError(err)
	return txResponse
}

func (suite *TestSuite) QueryBlock(blockHeight int64) *tmservice.Block {
	txResponse, err := suite.GRPCClient().GetBlockByHeight(blockHeight)
	suite.NoError(err)
	return txResponse
}

func (suite *TestSuite) QueryBlockByTxHash(txHash string) *tmservice.Block {
	txResponse := suite.QueryTx(txHash)
	return suite.QueryBlock(txResponse.Height)
}

func (suite *TestSuite) BroadcastTx(privKey cryptotypes.PrivKey, msgList ...sdk.Msg) *sdk.TxResponse {
	grpcClient := suite.GRPCClient()
	balances, err := grpcClient.QueryBalances(sdk.AccAddress(privKey.PubKey().Address().Bytes()).String())
	suite.NoError(err)
	suite.True(balances.AmountOf(suite.GetStakingDenom()).GT(sdkmath.NewInt(2).MulRaw(1e18)))

	gasPrices, err := sdk.ParseCoinsNormalized(suite.network.Config.MinGasPrices)
	suite.NoError(err)
	grpcClient.WithGasPrices(gasPrices)
	txRaw, err := grpcClient.BuildTxV1(privKey, msgList, 500000, "", 0)
	suite.NoError(err)

	txResponse, err := grpcClient.BroadcastTxOk(txRaw, tx.BroadcastMode_BROADCAST_MODE_BLOCK)
	suite.NoError(err)
	// txResponse might be nil, but error is also nil
	suite.NotNil(txResponse)
	suite.T().Log("broadcast tx", "msg:", sdk.MsgTypeURL(msgList[0]), "txHash:", txResponse.TxHash)
	return txResponse
}

func (suite *TestSuite) BroadcastProposalTx(content govv1beta1.Content, expectedStatus ...govv1.ProposalStatus) (*sdk.TxResponse, uint64) {
	proposalMsg, err := govv1beta1.NewMsgSubmitProposal(
		content,
		sdk.NewCoins(suite.NewCoin(sdkmath.NewInt(10_000).MulRaw(1e18))),
		suite.GetFirstValAddr().Bytes(),
	)
	suite.NoError(err)
	proposalId := suite.getNextProposalId()
	voteMsg := govv1beta1.NewMsgVote(suite.GetFirstValAddr().Bytes(), proposalId, govv1beta1.OptionYes)
	txResponse := suite.BroadcastTx(suite.GetFirstValPrivKey(), proposalMsg, voteMsg)
	for _, log := range txResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "proposal_deposit" {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != "proposal_id" {
					continue
				}
				id, err := strconv.ParseUint(attribute.Value, 10, 64)
				suite.NoError(err)
				suite.Require().Equal(proposalId, id)
				break
			}
		}
	}
	_, err = suite.network.WaitForHeight(txResponse.Height + 2)
	suite.NoError(err)
	status := govv1.StatusPassed
	if len(expectedStatus) > 0 {
		status = expectedStatus[0]
	}
	suite.CheckProposal(proposalId, status)
	return txResponse, proposalId
}

func (suite *TestSuite) BroadcastProposalTx2(msgs []sdk.Msg, title, summary string, expectedStatus ...govv1.ProposalStatus) (*sdk.TxResponse, uint64) {
	fxMetadata := fxgovtypes.NewFXMetadata(title, summary, "")
	proposalMsg, err := govv1.NewMsgSubmitProposal(
		msgs,
		sdk.NewCoins(suite.NewCoin(sdkmath.NewInt(10_000).MulRaw(1e18))),
		sdk.AccAddress(suite.GetFirstValAddr().Bytes()).String(),
		fxMetadata.String())
	suite.NoError(err)
	proposalId := suite.getNextProposalId()
	voteMsg := govv1.NewMsgVote(suite.GetFirstValAddr().Bytes(), proposalId, govv1.OptionYes, "")
	txResponse := suite.BroadcastTx(suite.GetFirstValPrivKey(), proposalMsg, voteMsg)
	for _, log := range txResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "proposal_deposit" {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != "proposal_id" {
					continue
				}
				id, err := strconv.ParseUint(attribute.Value, 10, 64)
				suite.NoError(err)
				suite.Require().Equal(proposalId, id)
				break
			}
		}
	}
	_, err = suite.network.WaitForHeight(txResponse.Height + 2)
	suite.NoError(err)
	status := govv1.StatusPassed
	if len(expectedStatus) > 0 {
		status = expectedStatus[0]
	}
	suite.CheckProposal(proposalId, status)
	return txResponse, proposalId
}

func (suite *TestSuite) CreateValidator(valPriv cryptotypes.PrivKey) *sdk.TxResponse {
	valAddr := sdk.ValAddress(valPriv.PubKey().Address())
	selfDelegate := sdk.NewCoin(suite.GetStakingDenom(), sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(100)))
	minSelfDelegate := sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1))
	description := stakingtypes.Description{
		Moniker:         "val2",
		Identity:        "",
		Website:         "",
		SecurityContact: "",
		Details:         "",
	}
	rates := stakingtypes.CommissionRates{
		Rate:          sdk.NewDecWithPrec(2, 2),  // 2%
		MaxRate:       sdk.NewDecWithPrec(50, 2), // 5%
		MaxChangeRate: sdk.NewDecWithPrec(2, 2),  // 2%
	}
	ed25519PrivKey := ed25519.GenPrivKeyFromSecret(valAddr.Bytes())
	msg, err := stakingtypes.NewMsgCreateValidator(valAddr, ed25519PrivKey.PubKey(), selfDelegate, description, rates, minSelfDelegate)
	suite.NoError(err)
	return suite.BroadcastTx(valPriv, msg)
}

func (suite *TestSuite) QueryValidatorByToken() sdk.ValAddress {
	response, err := suite.GRPCClient().StakingQuery().Validators(suite.ctx, &stakingtypes.QueryValidatorsRequest{Status: stakingtypes.Bonded.String()})
	suite.NoError(err)
	suite.True(len(response.Validators) > 0)
	validators := response.Validators
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Tokens.LT(validators[j].Tokens)
	})
	valAddr, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
	suite.NoError(err)
	return valAddr
}

func (suite *TestSuite) Send(toAddress sdk.AccAddress, amount sdk.Coin) *sdk.TxResponse {
	return suite.SendFrom(suite.GetFirstValPrivKey(), toAddress, amount)
}

func (suite *TestSuite) SendFrom(priv cryptotypes.PrivKey, toAddress sdk.AccAddress, amount sdk.Coin) *sdk.TxResponse {
	txResponse := suite.BroadcastTx(priv, banktypes.NewMsgSend(priv.PubKey().Address().Bytes(), toAddress, sdk.NewCoins(amount)))
	suite.True(txResponse.GasUsed < 100_000, txResponse.GasUsed)
	return txResponse
}

func (suite *TestSuite) QueryBalances(accAddress sdk.AccAddress) sdk.Coins {
	balances, err := suite.GRPCClient().QueryBalances(accAddress.String())
	suite.NoError(err)
	return balances
}

func (suite *TestSuite) CheckBalance(accAddress sdk.AccAddress, balance sdk.Coin) {
	queryBalance, err := suite.GRPCClient().QueryBalance(accAddress.String(), balance.Denom)
	suite.NoError(err)
	suite.Equal(balance.String(), queryBalance.String())
}

func (suite *TestSuite) SetWithdrawAddr(priv cryptotypes.PrivKey, withdrawAddr sdk.AccAddress) *sdk.TxResponse {
	fromAddr := sdk.AccAddress(priv.PubKey().Address().Bytes())
	return suite.BroadcastTx(priv, distritypes.NewMsgSetWithdrawAddress(fromAddr, withdrawAddr))
}

func (suite *TestSuite) CheckWithdrawAddr(delegatorAddr, withdrawAddr sdk.AccAddress) {
	withdrawAddressResp, err := suite.GRPCClient().DistrQuery().DelegatorWithdrawAddress(suite.ctx, &distritypes.QueryDelegatorWithdrawAddressRequest{
		DelegatorAddress: delegatorAddr.String(),
	})
	suite.NoError(err)
	suite.Equal(withdrawAddressResp.WithdrawAddress, withdrawAddr.String())
}

func (suite *TestSuite) Delegate(priv cryptotypes.PrivKey, valAddress sdk.ValAddress, amount sdk.Coin) *sdk.TxResponse {
	return suite.BroadcastTx(priv, stakingtypes.NewMsgDelegate(priv.PubKey().Address().Bytes(), valAddress, amount))
}

func (suite *TestSuite) CheckDelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, delegation sdk.Coin) {
	delegationResp, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	if delegation.IsZero() {
		suite.Error(errortypes.ErrNotFound)
	} else {
		suite.NoError(err)
		suite.Equal(delegation, delegationResp.DelegationResponse.Balance)
	}
}

func (suite *TestSuite) WithdrawReward(priv cryptotypes.PrivKey, valAddress sdk.ValAddress) *sdk.TxResponse {
	return suite.BroadcastTx(priv, distritypes.NewMsgWithdrawDelegatorReward(priv.PubKey().Address().Bytes(), valAddress))
}

func (suite *TestSuite) Undelegate(priv cryptotypes.PrivKey, valAddress sdk.ValAddress, amount sdk.Coin) *sdk.TxResponse {
	if amount.IsZero() {
		delegation, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: sdk.AccAddress(priv.PubKey().Address().Bytes()).String(),
			ValidatorAddr: valAddress.String(),
		})
		suite.NoError(err)
		amount = delegation.DelegationResponse.Balance
	}
	return suite.BroadcastTx(priv, stakingtypes.NewMsgUndelegate(priv.PubKey().Address().Bytes(), valAddress, amount))
}

func (suite *TestSuite) CheckUndelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, entries ...stakingtypes.UnbondingDelegationEntry) {
	response, err := suite.GRPCClient().StakingQuery().UnbondingDelegation(suite.ctx, &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	suite.NoError(err)
	suite.Equal(len(response.Unbond.Entries), len(entries))
	for i, entry := range response.Unbond.Entries {
		suite.Equal(entry.String(), entries[i].String())
	}
}

func (suite *TestSuite) Redelegate(priv cryptotypes.PrivKey, valSrc, valDest sdk.ValAddress, all bool) *sdk.TxResponse {
	amt := sdkmath.NewInt(1)
	if all {
		delegation, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: sdk.AccAddress(priv.PubKey().Address().Bytes()).String(),
			ValidatorAddr: valSrc.String(),
		})
		suite.NoError(err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg := stakingtypes.NewMsgBeginRedelegate(priv.PubKey().Address().Bytes(), valSrc, valDest, sdk.NewCoin(suite.GetStakingDenom(), amt))
	return suite.BroadcastTx(priv, msg)
}

func (suite *TestSuite) CheckRedelegate(delegatorAddr sdk.AccAddress, redelegationResponses stakingtypes.RedelegationResponses) {
	redelegationResp, err := suite.GRPCClient().StakingQuery().Redelegations(suite.ctx, &stakingtypes.QueryRedelegationsRequest{DelegatorAddr: delegatorAddr.String()})
	suite.NoError(err)
	suite.Equal(len(redelegationResp.RedelegationResponses), len(redelegationResponses))
	for i, item := range redelegationResp.RedelegationResponses {
		suite.Equal(item.Redelegation.String(), redelegationResponses[i].Redelegation.String())
		for j, entry := range item.Entries {
			suite.Equal(entry.RedelegationEntry.String(), redelegationResponses[i].Entries[j].RedelegationEntry.String())
			suite.Equal(entry.Balance.String(), redelegationResponses[i].Entries[j].Balance.String())
		}
	}
}

func (suite *TestSuite) CheckProposals(depositor sdk.AccAddress) govv1.Proposals {
	proposalsResp, err := suite.GRPCClient().GovQuery().Proposals(suite.ctx, &govv1.QueryProposalsRequest{
		ProposalStatus: govv1.StatusDepositPeriod,
		Depositor:      depositor.String(),
	})
	suite.NoError(err)
	return proposalsResp.Proposals
}

func (suite *TestSuite) ProposalDeposit(priv cryptotypes.PrivKey, proposalID uint64, amount sdk.Coin) *sdk.TxResponse {
	return suite.BroadcastTx(priv, govv1beta1.NewMsgDeposit(priv.PubKey().Address().Bytes(), proposalID, sdk.NewCoins(amount)))
}

func (suite *TestSuite) CheckDeposit(proposalId uint64, depositor sdk.AccAddress, amount sdk.Coin) {
	depositResp, err := suite.GRPCClient().GovQuery().Deposit(suite.ctx, &govv1.QueryDepositRequest{
		ProposalId: proposalId,
		Depositor:  depositor.String(),
	})
	suite.NoError(err)
	suite.Equal(depositResp.Deposit.Amount, amount)
}

func (suite *TestSuite) ProposalVote(priv cryptotypes.PrivKey, proposalID uint64, option govv1beta1.VoteOption) *sdk.TxResponse {
	return suite.BroadcastTx(priv, govv1beta1.NewMsgVote(priv.PubKey().Address().Bytes(), proposalID, option))
}

func (suite *TestSuite) CheckProposal(proposalId uint64, _ govv1.ProposalStatus) govv1.Proposal {
	proposalResp, err := suite.GRPCClient().GovQuery().Proposal(suite.ctx, &govv1.QueryProposalRequest{
		ProposalId: proposalId,
	})
	suite.NoError(err)

	suite.Require().True(proposalResp.Proposal.Status > govv1.StatusDepositPeriod)
	return *proposalResp.Proposal
}
