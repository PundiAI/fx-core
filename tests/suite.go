package tests

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v2/app/helpers"
	"github.com/functionx/fx-core/v2/client/grpc"
	"github.com/functionx/fx-core/v2/client/jsonrpc"
	"github.com/functionx/fx-core/v2/testutil"
	"github.com/functionx/fx-core/v2/testutil/network"
)

type TestSuite struct {
	suite.Suite
	useNetwork bool
	network    *network.Network
	ctx        context.Context
	sync.Mutex
}

func NewTestSuite() *TestSuite {
	testSuite := &TestSuite{
		Suite:      suite.Suite{},
		useNetwork: true,
		Mutex:      sync.Mutex{},
		ctx:        context.Background(),
	}
	if os.Getenv("USE_NETWORK") == "false" {
		testSuite.useNetwork = false
	}
	// nolint
	return testSuite
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

func (suite *TestSuite) GetUseNetwork() bool {
	return suite.useNetwork
}

func (suite *TestSuite) SetupSuite() {
	if !suite.useNetwork {
		return
	}
	suite.T().Log("setting up integration test suite")

	cfg := testutil.DefaultNetworkConfig()
	cfg.Mnemonics = append(cfg.Mnemonics, helpers.NewMnemonic())
	cfg.NumValidators = 1
	cfg.VotingPeriod = 5 * time.Second

	baseDir, err := os.MkdirTemp(suite.T().TempDir(), cfg.ChainID)
	suite.Require().NoError(err)
	suite.network, err = network.New(suite.T(), baseDir, cfg)
	suite.Require().NoError(err)

	_, err = suite.network.WaitForHeight(1)
	suite.Require().NoError(err)
}

func (suite *TestSuite) TearDownSuite() {
	if !suite.useNetwork {
		return
	}
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *TestSuite) GetFirstValidtor() *network.Validator {
	return suite.network.Validators[0]
}

func (suite *TestSuite) AdminPrivateKey() cryptotypes.PrivKey {
	privKey, err := helpers.PrivKeyFromMnemonic(suite.network.Config.Mnemonics[0], hd.Secp256k1Type, 0, 0)
	suite.NoError(err)
	return privKey
}

func (suite *TestSuite) GRPCClient() *grpc.Client {
	grpcUrl := "http://localhost:9090"
	if suite.useNetwork {
		grpcUrl = fmt.Sprintf("http://%s", suite.GetFirstValidtor().AppConfig.GRPC.Address)
	}
	client, err := grpc.NewClient(grpcUrl)
	suite.NoError(err)
	client.WithContext(suite.ctx)
	return client
}

func (suite *TestSuite) NodeClient() *jsonrpc.NodeRPC {
	nodeUrl := "http://localhost:26657"
	if suite.useNetwork {
		nodeUrl = suite.GetFirstValidtor().RPCAddress
	}
	rpc := jsonrpc.NewNodeRPC(jsonrpc.NewFastClient(nodeUrl))
	rpc.WithContext(suite.ctx)
	return rpc
}

func (suite *TestSuite) ValAddress() sdk.ValAddress {
	return suite.AdminPrivateKey().PubKey().Address().Bytes()
}

func (suite *TestSuite) GetStakingDenom() string {
	return suite.network.Config.BondDenom
}

func (suite *TestSuite) NewCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(suite.GetStakingDenom(), amount)
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

func (suite *TestSuite) QueryBlock(blockHeight int64) *types.Block {
	txResponse, err := suite.GRPCClient().GetBlockByHeight(blockHeight)
	suite.NoError(err)
	return txResponse
}

func (suite *TestSuite) QueryBlockByTxHash(txHash string) *types.Block {
	txResponse := suite.QueryTx(txHash)
	return suite.QueryBlock(txResponse.Height)
}

func (suite *TestSuite) BroadcastTx(privKey cryptotypes.PrivKey, msgList ...sdk.Msg) *sdk.TxResponse {
	suite.Lock()
	defer suite.Unlock()
	grpcClient := suite.GRPCClient()
	balances, err := grpcClient.QueryBalances(sdk.AccAddress(privKey.PubKey().Address().Bytes()).String())
	suite.NoError(err)
	suite.True(balances.AmountOf(suite.GetStakingDenom()).GT(sdk.NewInt(2).MulRaw(1e18)))

	grpcClient.WithGasPrices(sdk.NewCoins(sdk.NewCoin(suite.GetStakingDenom(), sdk.NewInt(4_000).MulRaw(1e9))))
	txRaw, err := grpcClient.BuildTxV2(privKey, msgList, 500000, "", 0)
	suite.NoError(err)

	txResponse, err := grpcClient.BroadcastTxOk(txRaw)
	suite.NoError(err)
	// txResponse might be nil, but error is also nil
	if txResponse != nil {
		suite.Equal(uint32(0), txResponse.Code)
	}

	if suite.useNetwork {
		suite.NoError(suite.network.WaitForNextBlock())
	}
	return txResponse
}

func (suite *TestSuite) BroadcastProposalTx(privKey cryptotypes.PrivKey, msgList ...sdk.Msg) (proposalId uint64) {
	txResponse := suite.BroadcastTx(privKey, msgList...)
	suite.T().Log("proposal submit txHash", txResponse.TxHash)
	for _, log := range txResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "proposal_deposit" {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != "proposal_id" {
					continue
				}
				proposalId, err := strconv.ParseUint(attribute.Value, 10, 64)
				suite.NoError(err)
				suite.CheckProposal(proposalId, govtypes.StatusVotingPeriod)
				return proposalId
			}
		}
	}
	return 0
}

func (suite *TestSuite) CreateValidator(valPriv cryptotypes.PrivKey) {
	valAddr := sdk.ValAddress(valPriv.PubKey().Address())
	selfDelegate := sdk.NewCoin(suite.GetStakingDenom(), sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100)))
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
	suite.NoError(err)
	txResponse := suite.BroadcastTx(valPriv, msg)
	suite.T().Log("create validator txHash", txResponse.TxHash)
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

func (suite *TestSuite) Send(toAddress sdk.AccAddress, amount sdk.Coin) {
	suite.SendFrom(suite.AdminPrivateKey(), toAddress, amount)
}

func (suite *TestSuite) SendFrom(priv cryptotypes.PrivKey, toAddress sdk.AccAddress, amount sdk.Coin) {
	txResponse := suite.BroadcastTx(priv, banktypes.NewMsgSend(priv.PubKey().Address().Bytes(), toAddress, sdk.NewCoins(amount)))
	suite.T().Log("send txHash", txResponse.TxHash)
}

func (suite *TestSuite) QueryBalances(accAddress sdk.AccAddress) sdk.Coins {
	balances, err := suite.GRPCClient().QueryBalances(accAddress.String())
	suite.NoError(err)
	return balances
}

func (suite *TestSuite) CheckBalance(accAddress sdk.AccAddress, balance sdk.Coin) {
	queryBalance, err := suite.GRPCClient().QueryBalance(accAddress.String(), balance.Denom)
	suite.NoError(err)
	suite.Equal(queryBalance, balance)
}

func (suite *TestSuite) SetWithdrawAddr(priv cryptotypes.PrivKey, withdrawAddr sdk.AccAddress) {
	fromAddr := sdk.AccAddress(priv.PubKey().Address().Bytes())
	txResponse := suite.BroadcastTx(priv, distritypes.NewMsgSetWithdrawAddress(fromAddr, withdrawAddr))
	suite.T().Log("set withdraw txHash", txResponse.TxHash)
}

func (suite *TestSuite) CheckWithdrawAddr(delegatorAddr, withdrawAddr sdk.AccAddress) {
	withdrawAddressResp, err := suite.GRPCClient().DistrQuery().DelegatorWithdrawAddress(suite.ctx, &distritypes.QueryDelegatorWithdrawAddressRequest{
		DelegatorAddress: delegatorAddr.String(),
	})
	suite.NoError(err)
	suite.Equal(withdrawAddressResp.WithdrawAddress, withdrawAddr.String())
}

func (suite *TestSuite) Delegate(priv cryptotypes.PrivKey, valAddress sdk.ValAddress, amount sdk.Coin) {
	txResponse := suite.BroadcastTx(priv, stakingtypes.NewMsgDelegate(priv.PubKey().Address().Bytes(), valAddress, amount))
	suite.T().Log("delegate txHash", txResponse.TxHash)
}

func (suite *TestSuite) CheckDelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, delegation sdk.Coin) {
	delegationResp, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	if delegation.IsZero() {
		suite.Error(sdkerrors.ErrNotFound)
	} else {
		suite.NoError(err)
		suite.Equal(delegation, delegationResp.DelegationResponse.Balance)
	}
}

func (suite *TestSuite) WithdrawReward(priv cryptotypes.PrivKey, valAddress sdk.ValAddress) {
	txResponse := suite.BroadcastTx(priv, distritypes.NewMsgWithdrawDelegatorReward(priv.PubKey().Address().Bytes(), valAddress))
	suite.T().Log("withdraw reward txHash", txResponse.TxHash)
}

func (suite *TestSuite) Undelegate(priv cryptotypes.PrivKey, valAddress sdk.ValAddress, amount sdk.Coin) string {
	if amount.IsZero() {
		delegation, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: sdk.AccAddress(priv.PubKey().Address().Bytes()).String(),
			ValidatorAddr: valAddress.String(),
		})
		suite.NoError(err)
		amount = delegation.DelegationResponse.Balance
	}
	txResponse := suite.BroadcastTx(priv, stakingtypes.NewMsgUndelegate(priv.PubKey().Address().Bytes(), valAddress, amount))
	suite.T().Log("undelegate txHash", txResponse.TxHash)
	return txResponse.TxHash
}

func (suite *TestSuite) CheckUndelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, entries ...stakingtypes.UnbondingDelegationEntry) {
	undelegationResp, err := suite.GRPCClient().StakingQuery().UnbondingDelegation(suite.ctx, &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	suite.NoError(err)
	suite.T().Log(undelegationResp.Unbond.Entries)
}

func (suite *TestSuite) Redelegate(priv cryptotypes.PrivKey, valSrc, valDest sdk.ValAddress, all bool) {
	amt := sdk.NewInt(1)
	if all {
		delegation, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: sdk.AccAddress(priv.PubKey().Address().Bytes()).String(),
			ValidatorAddr: valSrc.String(),
		})
		suite.NoError(err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg := stakingtypes.NewMsgBeginRedelegate(priv.PubKey().Address().Bytes(), valSrc, valDest, sdk.NewCoin(suite.GetStakingDenom(), amt))
	txResponse := suite.BroadcastTx(priv, msg)
	suite.T().Log("redelegate txHash", txResponse.TxHash)
}

func (suite *TestSuite) CheckRedelegate(delegatorAddr sdk.AccAddress, entries []stakingtypes.RedelegationResponses) {
	redelegationResp, err := suite.GRPCClient().StakingQuery().Redelegations(suite.ctx, &stakingtypes.QueryRedelegationsRequest{DelegatorAddr: delegatorAddr.String()})
	suite.NoError(err)
	suite.T().Log(redelegationResp.RedelegationResponses)
}

func (suite *TestSuite) ProposalSubmit(priv cryptotypes.PrivKey, deposit sdk.Coin) (proposalId uint64) {
	content := govtypes.ContentFromProposalType("title", "description", "Text")
	msg, err := govtypes.NewMsgSubmitProposal(content, sdk.NewCoins(deposit), priv.PubKey().Address().Bytes())
	suite.NoError(err)
	return suite.BroadcastProposalTx(priv, msg)
}

func (suite *TestSuite) CheckProposals(depositor sdk.AccAddress) govtypes.Proposals {
	proposalsResp, err := suite.GRPCClient().GovQuery().Proposals(suite.ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusDepositPeriod,
		Depositor:      depositor.String(),
	})
	suite.NoError(err)
	return proposalsResp.Proposals
}

func (suite *TestSuite) ProposalDeposit(priv cryptotypes.PrivKey, proposalID uint64, amount sdk.Coin) {
	txResponse := suite.BroadcastTx(priv, govtypes.NewMsgDeposit(priv.PubKey().Address().Bytes(), proposalID, sdk.NewCoins(amount)))
	suite.T().Log("proposal deposit txHash", txResponse.TxHash)
}

func (suite *TestSuite) CheckDeposit(proposalId uint64, depositor sdk.AccAddress, amount sdk.Coin) {
	depositResp, err := suite.GRPCClient().GovQuery().Deposit(suite.ctx, &govtypes.QueryDepositRequest{
		ProposalId: proposalId,
		Depositor:  depositor.String(),
	})
	suite.NoError(err)
	suite.Equal(depositResp.Deposit.Amount, amount)
}

func (suite *TestSuite) ProposalVote(priv cryptotypes.PrivKey, proposalID uint64, option govtypes.VoteOption) {
	txResponse := suite.BroadcastTx(priv, govtypes.NewMsgVote(priv.PubKey().Address().Bytes(), proposalID, option))
	suite.T().Log("proposal vote txHash", txResponse.TxHash)
}

func (suite *TestSuite) CheckProposal(proposalId uint64, status govtypes.ProposalStatus) govtypes.Proposal {
	timeoutCtx, cancel := context.WithTimeout(suite.ctx, suite.network.Config.VotingPeriod)
	defer cancel()
	for {
		proposalResp, err := suite.GRPCClient().GovQuery().Proposal(timeoutCtx, &govtypes.QueryProposalRequest{
			ProposalId: proposalId,
		})
		suite.NoError(err)
		if proposalResp.Proposal.Status == status {
			return proposalResp.Proposal
		} else {
			suite.T().Log("proposal status", proposalId, proposalResp.Proposal.Status.String())
		}
		time.Sleep(500 * time.Millisecond)
	}
}
