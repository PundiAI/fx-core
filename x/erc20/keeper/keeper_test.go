package keeper_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	ante2 "github.com/functionx/fx-core/ante"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"

	"github.com/functionx/fx-core/app/helpers"
	upgradesv2 "github.com/functionx/fx-core/app/upgrades/v2"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	fxtypes "github.com/functionx/fx-core/types"

	ethereumtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/x/evm/statedb"
	evm "github.com/tharsis/ethermint/x/evm/types"

	app "github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/tests"
	"github.com/functionx/fx-core/x/erc20/types"

	ethermint "github.com/tharsis/ethermint/types"
)

const (
	PurseDenom = "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	app              *app.App
	queryClientEvm   evm.QueryClient
	queryClient      types.QueryClient
	address          common.Address
	consAddress      sdk.ConsAddress
	clientCtx        client.Context
	ethSigner        ethereumtypes.Signer
	signer           keyring.Signer
	mintFeeCollector bool
	privateKey       cryptotypes.PrivKey
	checkTx          bool
	purseBalance     sdk.Int
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {

	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = tests.NewSigner(priv)
	suite.privateKey = priv

	// consensus key
	priv, err = ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.consAddress = sdk.ConsAddress(priv.PubKey().Address())

	priv, err = ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(priv.PubKey().Address())

	// init app
	suite.app = helpers.Setup(suite.T(), false)

	suite.ctx = suite.app.BaseApp.NewContext(suite.checkTx, tmproto.Header{Height: 1, ChainID: "fxcore", ProposerAddress: consAddress, Time: time.Now().UTC()})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.Erc20Keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	acc := &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(suite.address.Bytes(), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, priv.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)
	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	amount := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, suite.address.Bytes(), sdk.NewCoins(amount))
	suite.Require().NoError(err)

	if !suite.purseBalance.IsNil() {
		amount := sdk.NewCoin(PurseDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))
		err = suite.app.BankKeeper.MintCoins(suite.ctx, bsctypes.ModuleName, sdk.NewCoins(amount))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, bsctypes.ModuleName, suite.address.Bytes(), sdk.NewCoins(amount))
		suite.Require().NoError(err)
	}

	encodingConfig := app.MakeEncodingConfig()
	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	suite.ethSigner = ethereumtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())

	//update metadata
	err = upgradesv2.UpdateFXMetadata(suite.ctx, suite.app.BankKeeper, suite.app.GetKey(banktypes.StoreKey))
	require.NoError(t, err)
	//migrate account to eth
	upgradesv2.MigrateAccountToEth(suite.ctx, suite.app.AccountKeeper)
	// init logic contract
	for _, contract := range fxtypes.GetInitContracts() {
		require.True(t, len(contract.Code) > 0)
		require.True(t, contract.Address != common.HexToAddress(fxtypes.EmptyEvmAddress))
		err := suite.app.Erc20Keeper.CreateContractWithCode(suite.ctx, contract.Address, contract.Code)
		require.NoError(t, err)
	}

	// register coin
	for _, metadata := range fxtypes.GetMetadata() {
		_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
		require.NoError(t, err)
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *KeeperTestSuite) StateDB() *statedb.StateDB {
	return statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))
}

func (suite *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

// DeployContract deploys the ERC20MinterBurnerDecimalsContract.
func (suite *KeeperTestSuite) DeployContract(from common.Address, name, symbol string, decimals uint8) (common.Address, error) {
	return suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, from, name, symbol, decimals, false)
}

func (suite *KeeperTestSuite) DeployContractDirectBalanceManipulation(name string, symbol string) common.Address {
	ctx := sdk.WrapSDKContext(suite.ctx)
	chainID := suite.app.EvmKeeper.ChainID()

	erc20Config := fxtypes.GetERC20()

	ctorArgs, err := erc20Config.ABI.Pack("", big.NewInt(1000000000000000000))
	suite.Require().NoError(err)

	data := append(erc20Config.Bin, ctorArgs...)
	args, err := json.Marshal(&evm.TransactionArgs{
		From: &suite.address,
		Data: (*hexutil.Bytes)(&data),
	})
	suite.Require().NoError(err)

	res, err := suite.queryClientEvm.EstimateGas(ctx, &evm.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	suite.Require().NoError(err)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, suite.address)

	erc20DeployTx := evm.NewTxContract(
		chainID,
		nonce,
		nil,     // amount
		res.Gas, // gasLimit
		nil,     // gasPrice
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		big.NewInt(1),
		data,                        // input
		&ethereumtypes.AccessList{}, // accesses
	)

	erc20DeployTx.From = suite.address.Hex()
	err = erc20DeployTx.Sign(ethereumtypes.LatestSignerForChainID(chainID), suite.signer)
	suite.Require().NoError(err)
	rsp, err := suite.app.EvmKeeper.EthereumTx(ctx, erc20DeployTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return crypto.CreateAddress(suite.address, nonce)
}

func (suite *KeeperTestSuite) Commit() {
	_ = suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) MintERC20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20Config := fxtypes.GetERC20()
	transferData, err := erc20Config.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) BurnERC20Token(contractAddr, from common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20Config := fxtypes.GetERC20()
	transferData, err := erc20Config.ABI.Pack("transfer", types.ModuleAddress, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) sendTx(contractAddr, from common.Address, transferData []byte) *evm.MsgEthereumTx {
	ctx := sdk.WrapSDKContext(suite.ctx)
	chainID := suite.app.EvmKeeper.ChainID()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contractAddr, From: &from, Data: (*hexutil.Bytes)(&transferData)})
	suite.Require().NoError(err)
	res, err := suite.queryClientEvm.EstimateGas(ctx, &evm.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	suite.Require().NoError(err)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, suite.address)

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64()*int64(res.Gas)))))

	ercTransferTx := evm.NewTx(
		chainID,
		nonce,
		&contractAddr,
		nil,
		res.Gas,
		nil,
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		big.NewInt(1),
		transferData,
		&ethereumtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = suite.address.Hex()
	err = ercTransferTx.Sign(ethereumtypes.LatestSignerForChainID(chainID), suite.signer)
	suite.Require().NoError(err)
	rsp, err := suite.app.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return ercTransferTx
}

func (suite *KeeperTestSuite) BalanceOf(contract, account common.Address) interface{} {
	erc20Config := fxtypes.GetERC20()
	res, err := suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20Config.ABI, types.ModuleAddress, contract, false, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := erc20Config.ABI.Unpack("balanceOf", res.Ret)
	if len(unpacked) == 0 {
		return nil
	}

	return unpacked[0]
}

func (suite *KeeperTestSuite) NameOf(contract common.Address) string {
	erc20Config := fxtypes.GetERC20()

	res, err := suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20Config.ABI, types.ModuleAddress, contract, false, "name")
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	unpacked, err := erc20Config.ABI.Unpack("name", res.Ret)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(unpacked)

	return fmt.Sprintf("%v", unpacked[0])
}

func (suite *KeeperTestSuite) TransferERC20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20Config := fxtypes.GetERC20()
	transferData, err := erc20Config.ABI.Pack("transfer", to, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func NewPriKey() cryptotypes.PrivKey {
	return secp256k1.GenPrivKey()
}

func sendEthTx(t *testing.T, ctx sdk.Context, myApp *app.App,
	signer keyring.Signer, from, contract common.Address, data []byte) {

	chainID := myApp.EvmKeeper.ChainID()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contract, From: &from, Data: (*hexutil.Bytes)(&data)})
	require.NoError(t, err)
	res, err := myApp.EvmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evm.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	require.NoError(t, err)

	nonce, err := myApp.AccountKeeper.GetSequence(ctx, from.Bytes())
	require.NoError(t, err)

	var feeCap, tipCap sdk.Int
	baseFee := myApp.FeeMarketKeeper.GetBaseFee(ctx)
	params := myApp.FeeMarketKeeper.GetParams(ctx)
	if sdk.NewDecFromBigInt(baseFee).LT(params.MinGasPrice) {
		feeCap = params.MinGasPrice.TruncateInt()
		tipCap = feeCap.Sub(sdk.NewIntFromBigInt(baseFee))
	}

	ercTransferTx := evm.NewTx(
		chainID,
		nonce,
		&contract,
		nil,
		res.Gas,
		nil,
		feeCap.BigInt(),
		tipCap.BigInt(),
		data,
		&ethereumtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = from.String()
	err = ercTransferTx.Sign(ethereumtypes.LatestSignerForChainID(chainID), signer)
	require.NoError(t, err)

	options := ante2.HandlerOptions{
		AccountKeeper:   myApp.AccountKeeper,
		BankKeeper:      myApp.BankKeeper,
		FeeMarketKeeper: myApp.FeeMarketKeeper,
		EvmKeeper:       myApp.EvmKeeper,
		SignModeHandler: app.MakeEncodingConfig().TxConfig.SignModeHandler(),
		SigGasConsumer:  ante2.DefaultSigVerificationGasConsumer,
	}
	require.NoError(t, options.Validate())
	handler := ante2.NewAnteHandler(options)

	clientCtx := client.Context{}.WithTxConfig(app.MakeEncodingConfig().TxConfig)
	tx, err := ercTransferTx.BuildTx(clientCtx.TxConfig.NewTxBuilder(), fxtypes.DefaultDenom)
	require.NoError(t, err)
	ctx, err = handler(ctx, tx, false)
	require.NoError(t, err)

	rsp, err := myApp.EvmKeeper.EthereumTx(sdk.WrapSDKContext(ctx), ercTransferTx)
	require.NoError(t, err)
	require.Empty(t, rsp.VmError)
}
