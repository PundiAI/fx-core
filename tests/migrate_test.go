package tests

import (
	"encoding/hex"
	"testing"
	"time"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	hd2 "github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/stretchr/testify/suite"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/evmos/ethermint/crypto/hd"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/app/helpers"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	migratetypes "github.com/functionx/fx-core/x/migrate/types"
)

type MigrateTestSuite struct {
	TestSuite
	toPrivateKey cryptotypes.PrivKey
	newValPriv   cryptotypes.PrivKey
}

func TestMigrateTestSuite(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	privKey, err := helpers.PrivKeyFromMnemonic(helpers.NewMnemonic(), hd.EthSecp256k1Type, 0, 0)
	require.NoError(t, err)

	valPrivKey, err := helpers.PrivKeyFromMnemonic(helpers.NewMnemonic(), hd2.Secp256k1Type, 0, 0)
	require.NoError(t, err)

	suite.Run(t, &MigrateTestSuite{
		TestSuite:    NewTestSuite(),
		toPrivateKey: privKey,
		newValPriv:   valPrivKey,
	})
}

func (suite *MigrateTestSuite) ToAccAddress() sdk.AccAddress {
	return suite.toPrivateKey.PubKey().Address().Bytes()
}

func (suite *MigrateTestSuite) NewValAddress() sdk.AccAddress {
	return suite.newValPriv.PubKey().Address().Bytes()
}

func (suite *MigrateTestSuite) MigrateAccount(toPrivateKey cryptotypes.PrivKey) {
	toAddress := common.BytesToAddress(toPrivateKey.PubKey().Address())
	suite.T().Log("migrate from", suite.AdminAddress().String(), "migrate to", toAddress.String())
	migrateSign, err := toPrivateKey.Sign(migratetypes.MigrateAccountSignatureHash(suite.AdminAddress(), toAddress.Bytes()))
	suite.Require().NoError(err)

	msg := migratetypes.NewMsgMigrateAccount(suite.AdminAddress(), toAddress, hex.EncodeToString(migrateSign))
	txHash := suite.BroadcastTx(msg)
	suite.T().Log("migrate account txHash", txHash)
}

func (suite *MigrateTestSuite) SetupAllSuite() {
	balances := suite.QueryBalance(suite.AdminAddress())
	amount := balances.AmountOf(fxtypes.DefaultDenom).QuoRaw(3)

	suite.Send(suite.ToAccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, amount))

	suite.Send(suite.NewValAddress(), sdk.NewCoin(fxtypes.DefaultDenom, amount))
}

func (suite *MigrateTestSuite) TestDelegate() {
	//default account balance
	suite.QueryBalance(suite.AdminAddress())
	//query
	vals := suite.QueryValidator()
	val := vals[0]
	//default account delegate
	suite.Delegate(val, sdk.NewCoin(
		fxtypes.DefaultDenom,
		sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100000)),
	))
	acc, _ := sdk.AccAddressFromBech32("fx1968jve3k63a3u9whswlu2gsns4p0fqn0acxzgg")
	suite.SetWithdrawAddr(acc)
	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	//suite.Delegate( val, toPrivateKey) //if not comments, can not migrate, to address has delegated

	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
	// to delegate
	suite.privateKey = suite.toPrivateKey
	suite.Delegate(val, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100000))))
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) TestWithdrawReward() {

	//default account balance
	suite.QueryBalance(suite.AdminAddress())
	//query
	vals := suite.QueryValidator()
	val := vals[0]
	//default account delegate
	suite.Delegate(val, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100000))))

	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	suite.privateKey = suite.toPrivateKey
	suite.WithdrawReward(val)

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) TestUnDelegate() {
	//query
	vals := suite.QueryValidator()
	val := vals[0]
	//default account delegate
	suite.Delegate(val, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100000))))
	//undelegate
	suite.Undelegate(val, false)
	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())
	//delegate undelegate
	//suite.Delegate( val, toPrivateKey)         //if not comments, can not migrate, to address has delegated
	//suite.Undelegate( val, true, toPrivateKey) //if not comments, can not migrate, to address has undelegated
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	time.Sleep(30 * time.Second)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) TestReDelegate() {
	//create validator
	suite.CreateValidator(suite.newValPriv)
	//query
	vals := suite.QueryValidator()
	val, val2 := vals[0], vals[1]
	//default account delegate
	suite.Delegate(val, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100000))))
	//redelegate
	suite.Redelegate(val, val2, false)
	suite.Redelegate(val, val2, false)
	suite.Redelegate(val, val2, false)
	suite.Redelegate(val, val2, false)
	suite.Redelegate(val, val2, false)
	suite.Redelegate(val, val2, false)
	suite.Redelegate(val, val2, false)
	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())
	//delegate redelegate
	//suite.Delegate( val, toPrivateKey)               //if not comments, can not migrate, to address has delegated
	//suite.Redelegate( val, val2, true, toPrivateKey) //if not comments, can not migrate, to address has delegated
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	time.Sleep(30 * time.Second)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) TestProposalDeposit() {

	//suite.ProposalSubmit()       //deposit period
	suite.ProposalSubmit(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10_000)))) //vote period

	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())

	suite.privateKey = suite.toPrivateKey
	suite.ProposalDeposit(1, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1))))

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) TestProposalVote() {
	//suite.ProposalSubmit()       //deposit period
	suite.ProposalSubmit(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10_000)))) //vote period

	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())

	suite.ProposalVote(1, govtypes.OptionYes)
	//suite.ProposalVote( toPrivateKey) //can not migrate, to address has voted

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) TestAll() {
	//create validator
	suite.CreateValidator(suite.newValPriv)
	//query
	vals := suite.QueryValidator()
	val, val2 := vals[0], vals[1]

	//default account delegate
	suite.Delegate(val, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10_000))))
	acc, _ := sdk.AccAddressFromBech32("fx1968jve3k63a3u9whswlu2gsns4p0fqn0acxzgg")
	suite.SetWithdrawAddr(acc)

	//to address
	toAddress := sdk.AccAddress(suite.toPrivateKey.PubKey().Address())

	//undelegate
	suite.Undelegate(val, false)

	//redelegate(max=7)
	suite.Redelegate(val, val2, false)

	//proposal
	suite.ProposalSubmit(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10_100))))

	//vote
	suite.ProposalVote(1, govtypes.OptionYes)

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	//migrate to
	suite.MigrateAccount(suite.toPrivateKey)

	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)

	time.Sleep(30 * time.Second)
	//query
	suite.checkAccount(suite.AdminAddress())
	suite.checkAccount(toAddress)
}

func (suite *MigrateTestSuite) checkAccount(acc sdk.AccAddress) {
	balances, err := suite.grpcClient.BankQuery().AllBalances(suite.ctx, &banktypes.QueryAllBalancesRequest{Address: acc.String()})
	suite.Require().NoError(err)
	suite.T().Log("all balance", balances.Balances.String())

	validators, err := suite.grpcClient.StakingQuery().Validators(suite.ctx, &stakingtypes.QueryValidatorsRequest{Status: stakingtypes.Bonded.String()})
	suite.Require().NoError(err)
	suite.Require().True(len(validators.Validators) > 0)

	withdrawAddr, err := suite.grpcClient.DistrQuery().DelegatorWithdrawAddress(suite.ctx, &distritypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: acc.String()})
	suite.Require().NoError(err)
	suite.T().Log("withdraw address", withdrawAddr.WithdrawAddress)

	for _, v := range validators.Validators {
		resp, err := suite.grpcClient.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: acc.String(),
			ValidatorAddr: v.OperatorAddress,
		})
		if err != nil {
			continue
		}
		suite.T().Log("delegate validator", v.OperatorAddress, "balance", resp.DelegationResponse.Balance.String())
	}

	for _, v := range validators.Validators {
		undelegationResp, err := suite.grpcClient.StakingQuery().UnbondingDelegation(suite.ctx, &stakingtypes.QueryUnbondingDelegationRequest{
			DelegatorAddr: acc.String(),
			ValidatorAddr: v.OperatorAddress,
		})
		if err != nil {
			continue
		}
		for _, e := range undelegationResp.Unbond.Entries {
			suite.T().Log("undelegate validator", v.OperatorAddress, "balance", e.Balance.String(), "time", e.CompletionTime.String())
		}
	}

	redelegationResp, err := suite.grpcClient.StakingQuery().Redelegations(suite.ctx, &stakingtypes.QueryRedelegationsRequest{DelegatorAddr: acc.String()})
	suite.Require().NoError(err)
	for _, r := range redelegationResp.RedelegationResponses {
		for _, e := range r.Entries {
			suite.T().Log("redelegate validator", r.Redelegation.ValidatorSrcAddress, "to", r.Redelegation.ValidatorDstAddress, "balance", e.Balance.String(), "time", e.RedelegationEntry.CompletionTime.String())
		}
	}

	proposalsResp, err := suite.grpcClient.GovQuery().Proposals(suite.ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusDepositPeriod,
		Depositor:      acc.String(),
	})
	suite.Require().NoError(err)
	for _, p := range proposalsResp.Proposals {
		depositResp, err := suite.grpcClient.GovQuery().Deposit(suite.ctx, &govtypes.QueryDepositRequest{ProposalId: p.ProposalId, Depositor: acc.String()})
		if err == nil {
			suite.T().Log("proposal deposit", "id", p.ProposalId, "title", p.GetTitle(), "status", p.Status.String(), "amount", depositResp.Deposit.Amount.String())
		}
	}
	proposalsResp, err = suite.grpcClient.GovQuery().Proposals(suite.ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusVotingPeriod,
		Depositor:      acc.String(),
	})
	suite.Require().NoError(err)
	for _, p := range proposalsResp.Proposals {
		depositResp, err := suite.grpcClient.GovQuery().Deposit(suite.ctx, &govtypes.QueryDepositRequest{ProposalId: p.ProposalId, Depositor: acc.String()})
		if err == nil {
			suite.T().Log("proposal vote-deposit", "id", p.ProposalId, "title", p.GetTitle(), "status", p.Status.String(), "amount", depositResp.Deposit.Amount.String())
		}
	}

	proposalsResp, err = suite.grpcClient.GovQuery().Proposals(suite.ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusVotingPeriod,
		Voter:          acc.String(),
	})
	suite.Require().NoError(err)
	for _, p := range proposalsResp.Proposals {
		depositResp, err := suite.grpcClient.GovQuery().Deposit(suite.ctx, &govtypes.QueryDepositRequest{ProposalId: p.ProposalId, Depositor: acc.String()})
		if err == nil {
			suite.T().Log("proposal vote", "id", p.ProposalId, "title", p.GetTitle(), "status", p.Status.String(), "amount", depositResp.Deposit.Amount.String())
		}
	}
}

func TestSignature(t *testing.T) {
	bz, err := hex.DecodeString("3741e28e26d1df113bffff063d4121d1559f9efa87cf0110aa3d0be1cf742018")
	require.NoError(t, err)

	pri := &ethsecp256k1.PrivKey{Key: bz}
	require.Equal(t, "0x77F2022532009c5EB4c6C70f395DEAaA793481Bc", common.BytesToAddress(pri.PubKey().Address()).String())

	sig, err := hex.DecodeString("a010cf5b836eb934203ce5cc79544c79c7abca116dc9181c600d69d4163574120d1f1d5fd18225288dc9b8386a98f35af2a34cec36ae67f73cf70726819a9e8001")
	require.NoError(t, err)

	from, err := sdk.AccAddressFromBech32("fx1as048peq0hzfz4d64ew68ulwfx90m27lmckwaq")
	require.NoError(t, err)

	to := common.HexToAddress("0x77F2022532009c5EB4c6C70f395DEAaA793481Bc")

	var bt []byte
	bt = append(bt, []byte(migratetypes.MigrateAccountSignaturePrefix)...)
	bt = append(bt, from.Bytes()...)
	bt = append(bt, to.Bytes()...)
	require.Equal(t, "4d6967726174654163636f756e743aec1f5387207dc49155baae5da3f3ee498afdabdf77f2022532009c5eb4c6c70f395deaaa793481bc", hex.EncodeToString(bt))

	hash := migratetypes.MigrateAccountSignatureHash(from, to.Bytes())
	require.Equal(t, "fc88a1e6d3bebe443be968e6c88fc7646bfc6fe31eea1d357dd347be6579dc80", hex.EncodeToString(hash))

	pubKey, err := crypto.SigToPub(hash, sig)
	require.NoError(t, err)

	address := crypto.PubkeyToAddress(*pubKey)
	require.Equal(t, "0x77F2022532009c5EB4c6C70f395DEAaA793481Bc", address.String())
}
