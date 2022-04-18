package migrate

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	fxtypes "github.com/functionx/fx-core/types"
	migratetypes "github.com/functionx/fx-core/x/migrate/types"
)

func TestStakingParams(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)

	params, err := cli.StakingQuery().Params(context.Background(), &stakingtypes.QueryParamsRequest{})
	require.NoError(t, err)
	t.Log("params", params)
}

func TestDelegate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()
	//default account balance
	cli.testQueryBalance(ctx, cli.FxAddress())
	//query
	vals := cli.testQueryValidator(ctx)
	val := vals[0]
	//default account delegate
	cli.testDelegate(ctx, val)
	acc, _ := sdk.AccAddressFromBech32("fx1968jve3k63a3u9whswlu2gsns4p0fqn0acxzgg")
	cli.testSetWithdrawAddr(ctx, acc)
	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())
	//send
	cli.testSend(ctx, toAddress, 1000000)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	//cli.testDelegate(ctx, val, toPrivateKey) //if not comments, can not migrate, to address has delegated

	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
	// to delegate
	cli.testDelegate(ctx, val, toPrivateKey)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestWithdrawReward(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()
	//default account balance
	cli.testQueryBalance(ctx, cli.FxAddress())
	//query
	vals := cli.testQueryValidator(ctx)
	val := vals[0]
	//default account delegate
	cli.testDelegate(ctx, val)

	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())
	//send
	cli.testSend(ctx, toAddress, 1000000)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	cli.testWithdrawReward(ctx, val, toPrivateKey)

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestUnDelegate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()
	//query
	vals := cli.testQueryValidator(ctx)
	val := vals[0]
	//default account delegate
	cli.testDelegate(ctx, val)
	//undelegate
	cli.testUndelegate(ctx, val, false)
	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())
	//send
	cli.testSend(ctx, toAddress, 1000000)
	//delegate undelegate
	//cli.testDelegate(ctx, val, toPrivateKey)         //if not comments, can not migrate, to address has delegated
	//cli.testUndelegate(ctx, val, true, toPrivateKey) //if not comments, can not migrate, to address has undelegated
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	time.Sleep(30 * time.Second)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestReDelegate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()
	//create validator
	cli.testCreateValidator(ctx, val2Mnemonic)
	//query
	vals := cli.testQueryValidator(ctx)
	val, val2 := vals[0], vals[1]
	//default account delegate
	cli.testDelegate(ctx, val)
	//redelegate
	cli.testRedelegate(ctx, val, val2, false)
	cli.testRedelegate(ctx, val, val2, false)
	cli.testRedelegate(ctx, val, val2, false)
	cli.testRedelegate(ctx, val, val2, false)
	cli.testRedelegate(ctx, val, val2, false)
	cli.testRedelegate(ctx, val, val2, false)
	cli.testRedelegate(ctx, val, val2, false)
	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())
	//send
	cli.testSend(ctx, toAddress, 1000000)
	//delegate redelegate
	//cli.testDelegate(ctx, val, toPrivateKey)               //if not comments, can not migrate, to address has delegated
	//cli.testRedelegate(ctx, val, val2, true, toPrivateKey) //if not comments, can not migrate, to address has delegated
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	time.Sleep(30 * time.Second)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestProposalDeposit(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()

	//cli.testProposalSubmit(ctx)       //deposit period
	cli.testProposalSubmit(ctx, true) //vote period

	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())

	cli.testSend(ctx, toAddress, 1000000)

	cli.testProposalDeposit(ctx, toPrivateKey)

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestProposalVote(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()

	//cli.testProposalSubmit(ctx)       //deposit period
	cli.testProposalSubmit(ctx, true) //vote period

	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())

	cli.testSend(ctx, toAddress, 1000000)

	cli.testProposalVote(ctx)
	//cli.testProposalVote(ctx, toPrivateKey) //can not migrate, to address has voted

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestAll(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cli := NewClient(t)
	ctx := context.Background()

	//create validator
	cli.testCreateValidator(ctx, val2Mnemonic)
	//query
	vals := cli.testQueryValidator(ctx)
	val, val2 := vals[0], vals[1]

	//default account delegate
	cli.testDelegate(ctx, val)
	acc, _ := sdk.AccAddressFromBech32("fx1968jve3k63a3u9whswlu2gsns4p0fqn0acxzgg")
	cli.testSetWithdrawAddr(ctx, acc)

	//to address
	toPrivateKey, err := mnemonicToEthSecp256k1(toAddressMnemonic)
	require.NoError(cli.t, err)
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())
	//send
	cli.testSend(ctx, toAddress, 1000000)

	//undelegate
	cli.testUndelegate(ctx, val, false)

	//redelegate(max=7)
	cli.testRedelegate(ctx, val, val2, false)

	//proposal
	cli.testProposalSubmit(ctx, true)

	//vote
	cli.testProposalVote(ctx)

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	//migrate to
	cli.testMigrateAccount(ctx, toPrivateKey)

	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)

	time.Sleep(30 * time.Second)
	//query
	cli.testQueryAccount(ctx, cli.FxAddress())
	cli.testQueryAccount(ctx, toAddress)
}

func TestSignature(t *testing.T) {
	bz, err := hex.DecodeString("3741e28e26d1df113bffff063d4121d1559f9efa87cf0110aa3d0be1cf742018")
	require.NoError(t, err)
	pri := &ethsecp256k1.PrivKey{Key: bz}
	t.Log("private key address", common.BytesToAddress(pri.PubKey().Address()).String())

	sig, err := hex.DecodeString("a010cf5b836eb934203ce5cc79544c79c7abca116dc9181c600d69d4163574120d1f1d5fd18225288dc9b8386a98f35af2a34cec36ae67f73cf70726819a9e8001")
	require.NoError(t, err)

	from, err := sdk.AccAddressFromBech32("fx1as048peq0hzfz4d64ew68ulwfx90m27lmckwaq")
	require.NoError(t, err)

	to := common.HexToAddress("0x77F2022532009c5EB4c6C70f395DEAaA793481Bc")

	var b []byte
	b = append(b, []byte(migratetypes.MigrateAccountSignaturePrefix)...)
	b = append(b, from.Bytes()...)
	b = append(b, to.Bytes()...)

	t.Log("data", hex.EncodeToString(b))

	hash := migratetypes.MigrateAccountSignatureHash(from, to.Bytes())
	t.Log("hash", hex.EncodeToString(hash))

	pubKey, err := crypto.SigToPub(hash, sig)
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(*pubKey)
	t.Log("address", address.String())
}

func (cli *Client) testQueryBalance(ctx context.Context, acc sdk.AccAddress) {
	balances, err := cli.BankQuery().AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: acc.String()})
	require.NoError(cli.t, err)
	cli.t.Log("address", acc.String(), "balance", balances.Balances.String())
}
func (cli *Client) testQueryValidator(ctx context.Context) []sdk.ValAddress {
	validators, err := cli.StakingQuery().Validators(ctx, &stakingtypes.QueryValidatorsRequest{Status: stakingtypes.Bonded.String()})
	require.NoError(cli.t, err)
	require.True(cli.t, len(validators.Validators) > 0)

	vals := make([]sdk.ValAddress, 0, len(validators.Validators))
	for _, v := range validators.Validators {
		valAddr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		require.NoError(cli.t, err)
		vals = append(vals, valAddr)
		cli.t.Log("query validator", "val", v.OperatorAddress, "token", v.Tokens.String(), "share", v.DelegatorShares.String())
	}

	return vals
}
func (cli *Client) testQueryAccount(ctx context.Context, acc sdk.AccAddress) {
	balances, err := cli.BankQuery().AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: acc.String()})
	require.NoError(cli.t, err)
	cli.t.Log("all balance", balances.Balances.String())

	validators, err := cli.StakingQuery().Validators(ctx, &stakingtypes.QueryValidatorsRequest{Status: stakingtypes.Bonded.String()})
	require.NoError(cli.t, err)
	require.True(cli.t, len(validators.Validators) > 0)

	withdrawAddr, err := cli.DistrQuery().DelegatorWithdrawAddress(ctx, &distritypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: acc.String()})
	require.NoError(cli.t, err)
	cli.t.Log("withdraw address", withdrawAddr.WithdrawAddress)

	for _, v := range validators.Validators {
		resp, err := cli.StakingQuery().Delegation(ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: acc.String(),
			ValidatorAddr: v.OperatorAddress,
		})
		if err != nil {
			continue
		}
		cli.t.Log("delegate validator", v.OperatorAddress, "balance", resp.DelegationResponse.Balance.String())
	}

	for _, v := range validators.Validators {
		undelegationResp, err := cli.StakingQuery().UnbondingDelegation(ctx, &stakingtypes.QueryUnbondingDelegationRequest{
			DelegatorAddr: acc.String(),
			ValidatorAddr: v.OperatorAddress,
		})
		if err != nil {
			continue
		}
		for _, e := range undelegationResp.Unbond.Entries {
			cli.t.Log("undelegate validator", v.OperatorAddress, "balance", e.Balance.String(), "time", e.CompletionTime.String())
		}
	}

	redelegationResp, err := cli.StakingQuery().Redelegations(ctx, &stakingtypes.QueryRedelegationsRequest{DelegatorAddr: acc.String()})
	require.NoError(cli.t, err)
	for _, r := range redelegationResp.RedelegationResponses {
		for _, e := range r.Entries {
			cli.t.Log("redelegate validator", r.Redelegation.ValidatorSrcAddress, "to", r.Redelegation.ValidatorDstAddress, "balance", e.Balance.String(), "time", e.RedelegationEntry.CompletionTime.String())
		}
	}

	proposalsResp, err := cli.GovQuery().Proposals(ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusDepositPeriod,
		Depositor:      acc.String(),
	})
	require.NoError(cli.t, err)
	for _, p := range proposalsResp.Proposals {
		depositResp, err := cli.GovQuery().Deposit(ctx, &govtypes.QueryDepositRequest{ProposalId: p.ProposalId, Depositor: acc.String()})
		if err == nil {
			cli.t.Log("proposal deposit", "id", p.ProposalId, "title", p.GetTitle(), "status", p.Status.String(), "amount", depositResp.Deposit.Amount.String())
		}
	}
	proposalsResp, err = cli.GovQuery().Proposals(ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusVotingPeriod,
		Depositor:      acc.String(),
	})
	require.NoError(cli.t, err)
	for _, p := range proposalsResp.Proposals {
		depositResp, err := cli.GovQuery().Deposit(ctx, &govtypes.QueryDepositRequest{ProposalId: p.ProposalId, Depositor: acc.String()})
		if err == nil {
			cli.t.Log("proposal vote-deposit", "id", p.ProposalId, "title", p.GetTitle(), "status", p.Status.String(), "amount", depositResp.Deposit.Amount.String())
		}
	}

	proposalsResp, err = cli.GovQuery().Proposals(ctx, &govtypes.QueryProposalsRequest{
		ProposalStatus: govtypes.StatusVotingPeriod,
		Voter:          acc.String(),
	})
	require.NoError(cli.t, err)
	for _, p := range proposalsResp.Proposals {
		depositResp, err := cli.GovQuery().Deposit(ctx, &govtypes.QueryDepositRequest{ProposalId: p.ProposalId, Depositor: acc.String()})
		if err == nil {
			cli.t.Log("proposal vote", "id", p.ProposalId, "title", p.GetTitle(), "status", p.Status.String(), "amount", depositResp.Deposit.Amount.String())
		}
	}
}

func (cli *Client) testCreateValidator(ctx context.Context, mnemonic string) sdk.ValAddress {
	valPriv, err := mnemonicToEthSecp256k1(mnemonic)
	require.NoError(cli.t, err)
	valAddr := sdk.ValAddress(valPriv.PubKey().Address())

	var msg sdk.Msg
	amt := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000))))
	msg = banktypes.NewMsgSend(cli.FxAddress(), sdk.AccAddress(valAddr), amt)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("send to validator txHash", txHash)

	oldKey := cli.privateKey
	defer cli.SetPrivateKey(oldKey)
	cli.SetPrivateKey(valPriv)

	ed25519, err := mnemonicToEd25519(mnemonic)
	require.NoError(cli.t, err)

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
	msg, err = stakingtypes.NewMsgCreateValidator(valAddr, ed25519.PubKey(), selfDelegate, description, rates, minSelfDelegate)
	require.NoError(cli.t, err)
	txHash = cli.BroadcastTx(msg)
	cli.t.Log("create validator txHash", txHash)

	return valAddr
}
func (cli *Client) testSend(ctx context.Context, acc sdk.AccAddress, amt int64, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)

		cli.SetPrivateKey(privateKey[0])
	}
	var msg sdk.Msg
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(amt))))
	msg = banktypes.NewMsgSend(cli.FxAddress(), acc, amount)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("send txHash", txHash)
}
func (cli *Client) testSetWithdrawAddr(ctx context.Context, acc sdk.AccAddress, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)

		cli.SetPrivateKey(privateKey[0])
	}
	msg := distritypes.NewMsgSetWithdrawAddress(cli.FxAddress(), acc)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("set withdraw txHash", txHash)
}
func (cli *Client) testDelegate(ctx context.Context, val sdk.ValAddress, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)

		cli.SetPrivateKey(privateKey[0])
	}
	var msg sdk.Msg
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100000))
	msg = stakingtypes.NewMsgDelegate(cli.FxAddress(), val, sdk.NewCoin(fxtypes.DefaultDenom, amt))
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("delegate txHash", txHash)
}
func (cli *Client) testWithdrawReward(ctx context.Context, val sdk.ValAddress, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)

		cli.SetPrivateKey(privateKey[0])
	}
	msg := distritypes.NewMsgWithdrawDelegatorReward(cli.FxAddress(), val)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("withdraw reward txHash", txHash)
}

func (cli *Client) testUndelegate(ctx context.Context, val sdk.ValAddress, all bool, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)

		cli.SetPrivateKey(privateKey[0])
	}
	var msg sdk.Msg
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))
	if all {
		delegation, err := cli.StakingQuery().Delegation(ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: cli.FxAddress().String(),
			ValidatorAddr: val.String(),
		})
		require.NoError(cli.t, err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg = stakingtypes.NewMsgUndelegate(cli.FxAddress(), val, sdk.NewCoin(fxtypes.DefaultDenom, amt))
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("undelegate txHash", txHash)
}
func (cli *Client) testRedelegate(ctx context.Context, valSrc, valDest sdk.ValAddress, all bool, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)

		cli.SetPrivateKey(privateKey[0])
	}
	var msg sdk.Msg
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))
	if all {
		delegation, err := cli.StakingQuery().Delegation(ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: cli.FxAddress().String(),
			ValidatorAddr: valSrc.String(),
		})
		require.NoError(cli.t, err)
		amt = delegation.DelegationResponse.Balance.Amount
	}
	msg = stakingtypes.NewMsgBeginRedelegate(cli.FxAddress(), valSrc, valDest, sdk.NewCoin(fxtypes.DefaultDenom, amt))
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("redelegate txHash", txHash)
}
func (cli *Client) testProposalSubmit(ctx context.Context, satisfyVote ...bool) {
	content := govtypes.ContentFromProposalType("title", "description", "Text")
	amount := int64(5000)
	if len(satisfyVote) > 0 && satisfyVote[0] {
		amount = 10000
	}
	initDeposit := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(amount))))
	msg, err := govtypes.NewMsgSubmitProposal(content, initDeposit, cli.FxAddress())
	require.NoError(cli.t, err)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("proposal submit txHash", txHash)
}
func (cli *Client) testProposalDeposit(ctx context.Context, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)
		cli.SetPrivateKey(privateKey[0])
	}
	depositAmt := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1))))
	msg := govtypes.NewMsgDeposit(cli.FxAddress(), 1, depositAmt)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("proposal deposit txHash", txHash)
}
func (cli *Client) testProposalVote(ctx context.Context, privateKey ...cryptotypes.PrivKey) {
	if len(privateKey) > 0 {
		oldKey := cli.privateKey
		defer cli.SetPrivateKey(oldKey)
		cli.SetPrivateKey(privateKey[0])
	}
	msg := govtypes.NewMsgVote(cli.FxAddress(), 1, govtypes.OptionYes)
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("proposal vote txHash", txHash)
}
func (cli *Client) testMigrateAccount(ctx context.Context, toPrivateKey cryptotypes.PrivKey) {
	toAddress := sdk.AccAddress(toPrivateKey.PubKey().Address())
	cli.t.Log("migrate from", cli.FxAddress().String(), "migrate to", toAddress.String())
	migrateSign, err := toPrivateKey.Sign(migratetypes.MigrateAccountSignatureHash(cli.FxAddress(), toAddress))
	require.NoError(cli.t, err)

	msg := migratetypes.NewMsgMigrateAccount(cli.FxAddress(), toAddress, hex.EncodeToString(migrateSign))
	txHash := cli.BroadcastTx(msg)
	cli.t.Log("migrate account txHash", txHash)
}
