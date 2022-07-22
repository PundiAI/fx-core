package tests

import (
	"encoding/hex"
	"testing"
	"time"

	hd2 "github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v2/app/helpers"
	fxtypes "github.com/functionx/fx-core/v2/types"
	migratetypes "github.com/functionx/fx-core/v2/x/migrate/types"
)

type MigrateTestSuite struct {
	*TestSuite
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, &MigrateTestSuite{TestSuite: NewTestSuite()})
}

func (suite *MigrateTestSuite) MigrateAccount(fromPrivateKey, toPrivateKey cryptotypes.PrivKey) {
	fromAddr := sdk.AccAddress(fromPrivateKey.PubKey().Address().Bytes())
	toAddress := common.BytesToAddress(toPrivateKey.PubKey().Address())
	suite.T().Log("migrate from", fromAddr.String(), "migrate to", toAddress.String())

	migrateSign, err := toPrivateKey.Sign(migratetypes.MigrateAccountSignatureHash(fromAddr, toAddress.Bytes()))
	suite.NoError(err)

	msg := migratetypes.NewMsgMigrateAccount(fromAddr, toAddress, hex.EncodeToString(migrateSign))
	txResp := suite.BroadcastTx(fromPrivateKey, msg)
	suite.T().Log("migrate account txHash", txResp.TxHash)
}

func (suite *MigrateTestSuite) TestDelegate() {
	fromPrivKey, err := helpers.PrivKeyFromMnemonic(helpers.NewMnemonic(), hd2.Secp256k1Type, 0, 0)
	suite.NoError(err)
	fromAccAddress := fromPrivKey.PubKey().Address().Bytes()
	amount := sdk.NewInt(20).MulRaw(1e18)
	suite.Send(fromAccAddress, helpers.NewCoin(amount))
	suite.CheckBalance(fromAccAddress, helpers.NewCoin(amount))

	valAddress := suite.QueryValidatorByToken()
	delegateAmount := helpers.NewCoin(sdk.NewInt(1).MulRaw(1e18))
	suite.Delegate(fromPrivKey, valAddress, delegateAmount)
	amount = amount.Sub(sdk.NewInt(3).MulRaw(1e18))
	suite.CheckBalance(fromAccAddress, helpers.NewCoin(amount))
	suite.CheckDelegate(fromAccAddress, valAddress, delegateAmount)

	withdrawAddr := sdk.AccAddress(helpers.NewPriKey().PubKey().Address().Bytes())
	suite.SetWithdrawAddr(fromPrivKey, withdrawAddr)
	amount = amount.Sub(sdk.NewInt(2).MulRaw(1e18))
	suite.CheckBalance(fromAccAddress, helpers.NewCoin(amount))
	suite.CheckWithdrawAddr(fromAccAddress, withdrawAddr)

	// ===> migration

	toPrivKey, err := helpers.PrivKeyFromMnemonic(helpers.NewMnemonic(), hd.EthSecp256k1Type, 0, 0)
	suite.NoError(err)
	toAccAddress := sdk.AccAddress(toPrivKey.PubKey().Address().Bytes())
	suite.CheckBalance(toAccAddress, helpers.NewCoin(sdk.ZeroInt()))

	suite.MigrateAccount(fromPrivKey, toPrivKey)
	amount = amount.Sub(sdk.NewInt(2).MulRaw(1e18))

	suite.CheckBalance(fromAccAddress, helpers.NewCoin(sdk.ZeroInt()))
	suite.CheckDelegate(fromAccAddress, valAddress, helpers.NewCoin(sdk.ZeroInt()))

	suite.CheckBalance(toAccAddress, helpers.NewCoin(amount))
	suite.CheckDelegate(toAccAddress, valAddress, delegateAmount)
	suite.CheckWithdrawAddr(toAccAddress, toAccAddress)

	suite.Delegate(toPrivKey, valAddress, helpers.NewCoin(sdk.NewInt(1).MulRaw(1e18)))
	amount = amount.Sub(sdk.NewInt(3).MulRaw(1e18))
	balances := suite.QueryBalances(toAccAddress)
	suite.True(balances.AmountOf(fxtypes.DefaultDenom).GT(amount))

	delegateAmount = delegateAmount.Add(helpers.NewCoin(sdk.NewInt(1).MulRaw(1e18)))
	suite.CheckDelegate(toAccAddress, valAddress, delegateAmount)

	suite.WithdrawReward(toPrivKey, valAddress)
	amount = amount.Sub(sdk.NewInt(2).MulRaw(1e18))
	balances2 := suite.QueryBalances(toAccAddress)
	suite.True(balances2.AmountOf(fxtypes.DefaultDenom).GT(amount))
}

func (suite *MigrateTestSuite) TestUnDelegate() {
	fromPrivKey, err := helpers.PrivKeyFromMnemonic(helpers.NewMnemonic(), hd2.Secp256k1Type, 0, 0)
	suite.NoError(err)
	fromAccAddress := fromPrivKey.PubKey().Address().Bytes()
	amount := sdk.NewInt(20).MulRaw(1e18)
	suite.Send(fromAccAddress, helpers.NewCoin(amount))

	valAddress := suite.QueryValidatorByToken()
	delegateAmount := helpers.NewCoin(sdk.NewInt(1).MulRaw(1e18))
	suite.Delegate(fromPrivKey, valAddress, delegateAmount)
	amount = amount.Sub(sdk.NewInt(3).MulRaw(1e18))

	txHash := suite.Undelegate(fromPrivKey, valAddress, helpers.NewCoin(sdk.ZeroInt()))
	delegateAmount = delegateAmount.Sub(helpers.NewCoin(sdk.NewInt(1).MulRaw(1e18)))
	amount = amount.Sub(sdk.NewInt(2).MulRaw(1e18))
	block := suite.QueryBlockByTxHash(txHash)
	unbondingDelegationEntry := stakingtypes.UnbondingDelegationEntry{
		CreationHeight: block.Header.Height,
		CompletionTime: block.Header.Time.Add(21 * 24 * time.Hour),
		InitialBalance: delegateAmount.Amount,
		Balance:        delegateAmount.Amount,
	}
	suite.CheckUndelegate(fromAccAddress, valAddress, unbondingDelegationEntry)

	// ===> migration

	toPrivKey, err := helpers.PrivKeyFromMnemonic(helpers.NewMnemonic(), hd.EthSecp256k1Type, 0, 0)
	suite.NoError(err)
	toAccAddress := sdk.AccAddress(toPrivKey.PubKey().Address().Bytes())

	suite.MigrateAccount(fromPrivKey, toPrivKey)
	amount = amount.Sub(sdk.NewInt(2).MulRaw(1e18))

	balances2 := suite.QueryBalances(toAccAddress)
	suite.True(balances2.AmountOf(fxtypes.DefaultDenom).GT(amount))
	suite.CheckDelegate(toAccAddress, valAddress, delegateAmount)
	suite.CheckUndelegate(toAccAddress, valAddress, unbondingDelegationEntry)
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
