package integration

import (
	"sort"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pundiai/fx-core/v8/client/grpc"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

type FxCoreSuite struct {
	*EthSuite

	codec             codec.Codec
	validators        []*helpers.Signer
	grpcCli           *grpc.Client
	gasPrices         sdk.Coins
	defDenom          string
	timeoutCommit     time.Duration
	waitForHeightFunc func(height int64) (int64, error)

	proposalId     uint64
	proposalStatus govv1.ProposalStatus
}

func (suite *FxCoreSuite) WithGasPrices(gasPrices ...sdk.Coin) *FxCoreSuite {
	newSuite := *suite
	newSuite.gasPrices = gasPrices
	return &newSuite
}

func (suite *FxCoreSuite) GetValSigner() *helpers.Signer {
	return suite.validators[0]
}

func (suite *FxCoreSuite) GetValAddr() sdk.ValAddress {
	return suite.GetValSigner().AccAddress().Bytes()
}

func (suite *FxCoreSuite) FindValSigner(addr sdk.AccAddress) *helpers.Signer {
	for _, validator := range suite.validators {
		if addr.Equals(validator.AccAddress()) {
			return validator
		}
	}
	return nil
}

func (suite *FxCoreSuite) getNextProposalId() uint64 {
	suite.proposalId = suite.proposalId + 1
	return suite.proposalId
}

// Deprecated: Use NewStakingCoin
func (suite *FxCoreSuite) NewCoin(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(suite.defDenom, amount)
}

func (suite *FxCoreSuite) NewStakingCoin(amount, power int64) sdk.Coin {
	coin := helpers.NewStakingCoin(amount, power)
	suite.Require().Equal(coin.Denom, suite.defDenom)
	return coin
}

func (suite *FxCoreSuite) BlockNumber() int64 {
	height, err := suite.grpcCli.GetBlockHeight()
	suite.Require().NoError(err)
	return height
}

func (suite *FxCoreSuite) QueryTx(txHash string) *sdk.TxResponse {
	txResponse, err := suite.grpcCli.TxByHash(txHash)
	suite.Require().NoError(err)
	return txResponse
}

func (suite *FxCoreSuite) QueryBlock(blockHeight int64) *cmtservice.Block {
	txResponse, err := suite.grpcCli.GetBlockByHeight(blockHeight)
	suite.Require().NoError(err)
	return txResponse
}

func (suite *FxCoreSuite) QueryBlockByTxHash(txHash string) *cmtservice.Block {
	txResponse := suite.QueryTx(txHash)
	return suite.QueryBlock(txResponse.Height)
}

func (suite *FxCoreSuite) BroadcastTx(signer *helpers.Signer, msgList ...sdk.Msg) *sdk.TxResponse {
	balances, err := suite.grpcCli.QueryBalances(signer.AccAddress().String())
	suite.Require().NoError(err)
	suite.Require().True(balances.AmountOf(suite.defDenom).GT(sdkmath.NewInt(2).MulRaw(1e18)))

	txRaw, err := suite.grpcCli.WithGasPrices(suite.gasPrices).
		BuildTxRaw(signer.PrivKey(), msgList, 500000, 0, "")
	suite.Require().NoError(err)

	txResponse, err := suite.grpcCli.BroadcastTx(txRaw)
	suite.Require().NoError(err)
	suite.Require().NotNil(txResponse) // txResponse might be nil, but error is also nil
	suite.Require().EqualValues(0, txResponse.Code)
	txResponse, err = suite.grpcCli.WaitMined(txResponse.TxHash, 200*suite.timeoutCommit, 10*suite.timeoutCommit)
	suite.Require().NoError(err)
	msgTypeURL := sdk.MsgTypeURL(msgList[0])
	suite.T().Log("broadcast tx", "msg:", msgTypeURL, "height:", txResponse.Height, "txHash:", txResponse.TxHash)
	waitBlock := txResponse.Height + 1
	if msgTypeURL == sdk.MsgTypeURL(&govv1beta1.MsgSubmitProposal{}) ||
		msgTypeURL == sdk.MsgTypeURL(&govv1.MsgSubmitProposal{}) {
		waitBlock += 2
	}
	_, err = suite.waitForHeightFunc(waitBlock)
	suite.Require().NoError(err)
	return txResponse
}

func (suite *FxCoreSuite) Send(toAddress sdk.AccAddress, amount ...sdk.Coin) *sdk.TxResponse {
	signer := suite.GetValSigner()
	txResponse := suite.BroadcastTx(signer, banktypes.NewMsgSend(signer.AccAddress(), toAddress, amount))
	suite.Require().Less(txResponse.GasUsed, int64(100_000))
	return txResponse
}

// --- Bank Module

func (suite *FxCoreSuite) GetAllBalances(accAddress sdk.AccAddress) sdk.Coins {
	balances, err := suite.grpcCli.QueryBalances(accAddress.String())
	suite.Require().NoError(err)
	return balances
}

func (suite *FxCoreSuite) EqualBalance(accAddress sdk.AccAddress, balance sdk.Coin) {
	queryBalance, err := suite.grpcCli.QueryBalance(accAddress.String(), balance.Denom)
	suite.Require().NoError(err)
	suite.Require().Equal(balance.String(), queryBalance.String())
}

func (suite *FxCoreSuite) GetDenomsMetadata() []banktypes.Metadata {
	resp, err := suite.grpcCli.BankQuery().DenomsMetadata(suite.ctx, &banktypes.QueryDenomsMetadataRequest{})
	suite.Require().NoError(err)
	return resp.Metadatas
}

func (suite *FxCoreSuite) GetMetadata(denom string) banktypes.Metadata {
	response, err := suite.grpcCli.BankQuery().DenomMetadata(suite.ctx, &banktypes.QueryDenomMetadataRequest{Denom: denom})
	suite.Require().NoError(err)
	return response.Metadata
}

// --- Auth Module

func (suite *FxCoreSuite) QueryModuleAccountByName(moduleName string) sdk.AccAddress {
	moduleAccount, err := suite.grpcCli.GetModuleAccountByName(moduleName)
	suite.Require().NoError(err)
	return moduleAccount.GetAddress()
}

// --- Gov Module

func (suite *FxCoreSuite) WithProposalStatus(status govv1.ProposalStatus) *FxCoreSuite {
	newSuite := *suite
	newSuite.proposalStatus = status
	return &newSuite
}

func (suite *FxCoreSuite) BroadcastProposalTxV1(msgs ...sdk.Msg) (*sdk.TxResponse, uint64) {
	proposalMsg, err := govv1.NewMsgSubmitProposal(
		msgs,
		sdk.NewCoins(suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18))),
		suite.GetValSigner().AccAddress().String(),
		"",
		sdk.MsgTypeURL(msgs[0]),
		sdk.MsgTypeURL(msgs[0]),
		false,
	)
	suite.Require().NoError(err)
	proposalId := suite.getNextProposalId()
	voteMsg := govv1.NewMsgVote(suite.GetValSigner().AccAddress(), proposalId, govv1.OptionYes, "")
	txResponse := suite.BroadcastTx(suite.GetValSigner(), proposalMsg, voteMsg)
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
				suite.Require().NoError(err)
				suite.Require().Equal(proposalId, id)
				break
			}
		}
	}
	suite.Require().NoError(err)

	if suite.proposalStatus > 0 {
		suite.EqualProposal(proposalId, suite.proposalStatus)
	}
	return txResponse, proposalId
}

func (suite *FxCoreSuite) GetProposals(depositor sdk.AccAddress) govv1.Proposals {
	proposalsResp, err := suite.grpcCli.GovQuery().Proposals(suite.ctx, &govv1.QueryProposalsRequest{
		ProposalStatus: govv1.StatusDepositPeriod,
		Depositor:      depositor.String(),
	})
	suite.Require().NoError(err)
	return proposalsResp.Proposals
}

func (suite *FxCoreSuite) ProposalDeposit(signer *helpers.Signer, proposalID uint64, amount sdk.Coin) *sdk.TxResponse {
	coins := sdk.NewCoins(amount)
	txResponse := suite.BroadcastTx(signer, govv1beta1.NewMsgDeposit(signer.AccAddress(), proposalID, coins))

	depositResp, err := suite.grpcCli.GovQuery().Deposit(suite.ctx, &govv1.QueryDepositRequest{
		ProposalId: proposalID,
		Depositor:  signer.AccAddress().String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(depositResp.Deposit.Amount, coins)
	return txResponse
}

func (suite *FxCoreSuite) ProposalVote(signer *helpers.Signer, proposalID uint64, option govv1beta1.VoteOption) *sdk.TxResponse {
	return suite.BroadcastTx(signer, govv1beta1.NewMsgVote(signer.AccAddress(), proposalID, option))
}

func (suite *FxCoreSuite) EqualProposal(proposalId uint64, _ govv1.ProposalStatus) govv1.Proposal {
	proposalResp, err := suite.grpcCli.GovQuery().Proposal(suite.ctx, &govv1.QueryProposalRequest{
		ProposalId: proposalId,
	})
	suite.Require().NoError(err)

	suite.Require().Greater(proposalResp.Proposal.Status, govv1.StatusDepositPeriod)
	return *proposalResp.Proposal
}

// --- Distribution Module

func (suite *FxCoreSuite) DistrQuery() distrtypes.QueryClient {
	return suite.grpcCli.DistrQuery()
}

func (suite *FxCoreSuite) EqualWithdrawAddr(delegatorAddr, withdrawAddr sdk.AccAddress) {
	withdrawAddressResp, err := suite.DistrQuery().DelegatorWithdrawAddress(suite.ctx,
		&distrtypes.QueryDelegatorWithdrawAddressRequest{
			DelegatorAddress: delegatorAddr.String(),
		})
	suite.Require().NoError(err)
	suite.Require().Equal(withdrawAddressResp.WithdrawAddress, withdrawAddr.String())
}

func (suite *FxCoreSuite) WithdrawReward(signer *helpers.Signer, valAddress sdk.ValAddress) *sdk.TxResponse {
	return suite.BroadcastTx(signer, distrtypes.NewMsgWithdrawDelegatorReward(signer.AccAddress().String(), valAddress.String()))
}

func (suite *FxCoreSuite) DelegationRewards(delAddr, valAddr string) sdk.DecCoins {
	response, err := suite.DistrQuery().DelegationRewards(suite.ctx,
		&distrtypes.QueryDelegationRewardsRequest{
			DelegatorAddress: delAddr, ValidatorAddress: valAddr,
		})
	suite.Require().NoError(err)
	return response.Rewards
}

func (suite *FxCoreSuite) SetWithdrawAddress(signer *helpers.Signer, withdrawAddr sdk.AccAddress) *sdk.TxResponse {
	delAddr := signer.AccAddress()
	setWithdrawAddress := distrtypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)
	txResponse := suite.BroadcastTx(signer, setWithdrawAddress)

	suite.EqualWithdrawAddr(delAddr, withdrawAddr)
	return txResponse
}

// --- Staking Module

func (suite *FxCoreSuite) StakingQuery() stakingtypes.QueryClient {
	return suite.grpcCli.StakingQuery()
}

func (suite *FxCoreSuite) GetDelegation(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) *stakingtypes.DelegationResponse {
	delegationResp, err := suite.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	if status.Code(err) == codes.NotFound {
		return nil
	}
	suite.Require().NoError(err)
	return delegationResp.DelegationResponse
}

func (suite *FxCoreSuite) CreateValidator(signer *helpers.Signer, toBondedVal bool) *sdk.TxResponse {
	valAddr := sdk.ValAddress(signer.AccAddress().Bytes())
	minSelfDelegate := sdkmath.NewInt(1)
	stakingDenom := suite.defDenom
	selfDelegate := sdk.NewCoin(stakingDenom, minSelfDelegate)
	if toBondedVal {
		selfDelegate = sdk.NewCoin(stakingDenom, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(100)))
	}
	description := stakingtypes.Description{
		Moniker:         "val2",
		Identity:        "",
		Website:         "",
		SecurityContact: "",
		Details:         "",
	}
	rates := stakingtypes.CommissionRates{
		Rate:          sdkmath.LegacyNewDecWithPrec(2, 2),  // 2%
		MaxRate:       sdkmath.LegacyNewDecWithPrec(50, 2), // 5%
		MaxChangeRate: sdkmath.LegacyNewDecWithPrec(2, 2),  // 2%
	}
	ed25519PrivKey := ed25519.GenPrivKeyFromSecret(valAddr.Bytes())
	msg, err := stakingtypes.NewMsgCreateValidator(valAddr.String(), ed25519PrivKey.PubKey(), selfDelegate, description, rates, minSelfDelegate)
	suite.Require().NoError(err)
	return suite.BroadcastTx(signer, msg)
}

func (suite *FxCoreSuite) GetValSortByToken() sdk.ValAddress {
	response, err := suite.StakingQuery().Validators(suite.ctx,
		&stakingtypes.QueryValidatorsRequest{Status: stakingtypes.Bonded.String()})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(response.Validators)
	validators := response.Validators
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Tokens.LT(validators[j].Tokens)
	})
	valAddr, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
	suite.Require().NoError(err)
	return valAddr
}

func (suite *FxCoreSuite) Delegate(signer *helpers.Signer, valAddress sdk.ValAddress, amount sdk.Coin) *sdk.TxResponse {
	delegation := suite.GetDelegation(signer.AccAddress(), valAddress)

	txResponse := suite.BroadcastTx(signer, stakingtypes.NewMsgDelegate(signer.AccAddress().String(), valAddress.String(), amount))

	if delegation != nil {
		amount = amount.Add(delegation.Balance)
	}
	suite.EqualDelegate(signer.AccAddress(), valAddress, amount)
	return txResponse
}

func (suite *FxCoreSuite) EqualDelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, delegation sdk.Coin) {
	delegationResp, err := suite.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	if delegation.IsZero() {
		suite.Require().Error(sdkerrors.ErrNotFound)
	} else {
		suite.Require().NoError(err)
		suite.Require().Equal(delegation.String(), delegationResp.DelegationResponse.Balance.String())
	}
}

func (suite *FxCoreSuite) Undelegate(signer *helpers.Signer, valAddress sdk.ValAddress, amount sdk.Coin) *sdk.TxResponse {
	if amount.IsZero() {
		delegation, err := suite.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: signer.AccAddress().String(),
			ValidatorAddr: valAddress.String(),
		})
		suite.Require().NoError(err)
		amount = delegation.DelegationResponse.Balance
	}
	txResponse := suite.BroadcastTx(signer, stakingtypes.NewMsgUndelegate(signer.AccAddress().String(), valAddress.String(), amount))
	return txResponse
}

func (suite *FxCoreSuite) EqualUndelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, entries ...stakingtypes.UnbondingDelegationEntry) {
	response, err := suite.StakingQuery().UnbondingDelegation(suite.ctx, &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(len(response.Unbond.Entries), len(entries))
	for i, entry := range response.Unbond.Entries {
		entry.UnbondingId = 0
		suite.Require().Equal(entry.String(), entries[i].String())
	}
}

func (suite *FxCoreSuite) Redelegate(signer *helpers.Signer, valSrc, valDest sdk.ValAddress, all bool) *sdk.TxResponse {
	amt := sdkmath.NewInt(1)
	if all {
		delegation, err := suite.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: signer.AccAddress().String(),
			ValidatorAddr: valSrc.String(),
		})
		suite.Require().NoError(err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg := stakingtypes.NewMsgBeginRedelegate(signer.AccAddress().String(), valSrc.String(), valDest.String(), sdk.NewCoin(suite.defDenom, amt))
	return suite.BroadcastTx(signer, msg)
}

func (suite *FxCoreSuite) EqualRedelegate(delegatorAddr sdk.AccAddress, redelegationResponses stakingtypes.RedelegationResponses) {
	redelegationResp, err := suite.StakingQuery().Redelegations(suite.ctx, &stakingtypes.QueryRedelegationsRequest{DelegatorAddr: delegatorAddr.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(len(redelegationResp.RedelegationResponses), len(redelegationResponses))
	for i, item := range redelegationResp.RedelegationResponses {
		suite.Require().Equal(item.Redelegation.String(), redelegationResponses[i].Redelegation.String())
		for j, entry := range item.Entries {
			suite.Require().Equal(entry.RedelegationEntry.String(), redelegationResponses[i].Entries[j].RedelegationEntry.String())
			suite.Require().Equal(entry.Balance.String(), redelegationResponses[i].Entries[j].Balance.String())
		}
	}
}

// --- ERC20 Module

func (suite *FxCoreSuite) ERC20Query() erc20types.QueryClient {
	return suite.grpcCli.ERC20Query()
}

func (suite *FxCoreSuite) GetErc20TokenAddress(denom string) common.Address {
	pair, err := suite.ERC20Query().TokenPair(suite.ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	suite.Require().NoError(err)
	return common.HexToAddress(pair.Erc20Token.Erc20Address)
}

func (suite *FxCoreSuite) ToggleTokenConversionProposal(denom string) (*sdk.TxResponse, uint64) {
	msg := &erc20types.MsgToggleTokenConversion{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Token:     denom,
	}
	return suite.BroadcastProposalTxV1(msg)
}

func (suite *FxCoreSuite) ConvertCoin(signer *helpers.Signer, recipient common.Address, coin sdk.Coin) *sdk.TxResponse {
	fromAddress := signer.AccAddress()

	beforeBalance := suite.GetAllBalances(fromAddress).AmountOf(coin.Denom)
	erc20TokenAddress := suite.GetErc20TokenAddress(coin.Denom)

	erc20TokenSuite := NewERC20TokenSuite(suite.EthSuite, erc20TokenAddress, signer)
	beforeBalanceOf := erc20TokenSuite.BalanceOf(recipient)

	msg := erc20types.NewMsgConvertCoin(coin, recipient, fromAddress)
	txResponse := suite.BroadcastTx(signer, msg)

	afterBalance := suite.GetAllBalances(fromAddress).AmountOf(coin.Denom)
	afterBalanceOf := erc20TokenSuite.BalanceOf(recipient)
	suite.Require().Equal(beforeBalance.Sub(afterBalance).String(), coin.Amount.String())
	suite.Require().Equal(afterBalanceOf.String(), beforeBalanceOf.String())
	return txResponse
}
