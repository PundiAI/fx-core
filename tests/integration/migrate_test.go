package integration

import (
	"encoding/hex"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	migratetypes "github.com/functionx/fx-core/v8/x/migrate/types"
)

func (suite *IntegrationTest) migrateAccount(fromSigner, toSigner *helpers.Signer) {
	fromAddr := fromSigner.AccAddress()
	toAddress := toSigner.Address()

	migrateSign, err := toSigner.PrivKey().Sign(migratetypes.MigrateAccountSignatureHash(fromAddr, toAddress.Bytes()))
	suite.Require().NoError(err)

	msg := migratetypes.NewMsgMigrateAccount(fromAddr, toAddress, hex.EncodeToString(migrateSign))
	suite.BroadcastTx(fromSigner, msg)
}

func (suite *IntegrationTest) MigrateTestDelegate() {
	fromSigner := helpers.NewSigner(helpers.NewPriKey())
	fromAccAddress := fromSigner.AccAddress()
	amount := sdkmath.NewInt(20).MulRaw(1e18)
	suite.Send(fromAccAddress, suite.NewCoin(amount))
	suite.EqualBalance(fromAccAddress, suite.NewCoin(amount))

	valAddress := suite.GetValSortByToken()
	delegateAmount := suite.NewCoin(sdkmath.NewInt(1).MulRaw(1e18))
	suite.Delegate(fromSigner, valAddress, delegateAmount)
	amount = amount.Sub(sdkmath.NewInt(3).MulRaw(1e18))
	suite.EqualBalance(fromAccAddress, suite.NewCoin(amount))
	suite.EqualDelegate(fromAccAddress, valAddress, delegateAmount)

	withdrawAddr := sdk.AccAddress(helpers.NewPriKey().PubKey().Address().Bytes())
	suite.SetWithdrawAddress(fromSigner, withdrawAddr)
	amount = amount.Sub(sdkmath.NewInt(2).MulRaw(1e18))
	suite.EqualBalance(fromAccAddress, suite.NewCoin(amount))

	// ===> migration

	toSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	toAccAddress := toSigner.AccAddress()
	suite.EqualBalance(toAccAddress, suite.NewCoin(sdkmath.ZeroInt()))

	suite.migrateAccount(fromSigner, toSigner)
	amount = amount.Sub(sdkmath.NewInt(2).MulRaw(1e18))

	suite.EqualBalance(fromAccAddress, suite.NewCoin(sdkmath.ZeroInt()))
	suite.EqualDelegate(fromAccAddress, valAddress, suite.NewCoin(sdkmath.ZeroInt()))

	suite.EqualBalance(toAccAddress, suite.NewCoin(amount))
	suite.EqualDelegate(toAccAddress, valAddress, delegateAmount)
	suite.EqualWithdrawAddr(toAccAddress, toAccAddress)

	suite.Delegate(toSigner, valAddress, suite.NewCoin(sdkmath.NewInt(1).MulRaw(1e18)))
	amount = amount.Sub(sdkmath.NewInt(3).MulRaw(1e18))
	balances := suite.GetAllBalances(toAccAddress)
	suite.True(balances.AmountOf(fxtypes.DefaultDenom).GT(amount))

	delegateAmount = delegateAmount.Add(suite.NewCoin(sdkmath.NewInt(1).MulRaw(1e18)))
	suite.EqualDelegate(toAccAddress, valAddress, delegateAmount)

	suite.WithdrawReward(toSigner, valAddress)
	amount = amount.Sub(sdkmath.NewInt(2).MulRaw(1e18))
	balances2 := suite.GetAllBalances(toAccAddress)
	suite.True(balances2.AmountOf(fxtypes.DefaultDenom).GT(amount))
}

func (suite *IntegrationTest) MigrateTestUnDelegate() {
	fromSigner := helpers.NewSigner(helpers.NewPriKey())
	fromAccAddress := fromSigner.AccAddress()
	amount := sdkmath.NewInt(20).MulRaw(1e18)
	suite.Send(fromAccAddress, suite.NewCoin(amount))

	valAddress := suite.GetValSortByToken()
	delegateAmount := suite.NewCoin(sdkmath.NewInt(2).MulRaw(1e18))
	suite.Delegate(fromSigner, valAddress, delegateAmount)
	amount = amount.Sub(sdkmath.NewInt(2 + 2).MulRaw(1e18))

	delegateAmount = delegateAmount.Sub(suite.NewCoin(sdkmath.NewInt(1).MulRaw(1e18)))
	txResponse := suite.Undelegate(fromSigner, valAddress, delegateAmount)
	amount = amount.Sub(sdkmath.NewInt(2).MulRaw(1e18))

	block := suite.QueryBlockByTxHash(txResponse.TxHash)
	unbondingDelegationEntry := stakingtypes.UnbondingDelegationEntry{
		CreationHeight: block.Header.Height,
		CompletionTime: block.Header.Time.Add(21 * 24 * time.Hour),
		InitialBalance: delegateAmount.Amount,
		Balance:        delegateAmount.Amount,
	}
	suite.EqualUndelegate(fromAccAddress, valAddress, unbondingDelegationEntry)

	// ===> migration

	toSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	toAccAddress := toSigner.AccAddress()

	suite.migrateAccount(fromSigner, toSigner)
	amount = amount.Sub(sdkmath.NewInt(2).MulRaw(1e18))

	balances2 := suite.GetAllBalances(toAccAddress)
	suite.True(balances2.AmountOf(fxtypes.DefaultDenom).GT(amount))
	suite.EqualDelegate(toAccAddress, valAddress, delegateAmount)
	suite.EqualUndelegate(toAccAddress, valAddress, unbondingDelegationEntry)
}

func TestSignature(t *testing.T) {
	bz, err := hex.DecodeString("3741e28e26d1df113bffff063d4121d1559f9efa87cf0110aa3d0be1cf742018")
	require.NoError(t, err)

	pri := &ethsecp256k1.PrivKey{Key: bz}
	require.Equal(t, "0x77F2022532009c5EB4c6C70f395DEAaA793481Bc", common.BytesToAddress(pri.PubKey().Address()).String())

	sig, err := hex.DecodeString("a010cf5b836eb934203ce5cc79544c79c7abca116dc9181c600d69d4163574120d1f1d5fd18225288dc9b8386a98f35af2a34cec36ae67f73cf70726819a9e8001")
	require.NoError(t, err)

	from, err := sdk.AccAddressFromHexUnsafe("ec1f5387207dc49155baae5da3f3ee498afdabdf")
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
