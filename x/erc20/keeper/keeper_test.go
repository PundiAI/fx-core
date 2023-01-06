package keeper_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/server/config"
	"github.com/evmos/ethermint/x/evm/statedb"
	evm "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.App
	queryClient types.QueryClient
	signer      *helpers.Signer
	randSigners map[common.Address]*helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) SetupTest() {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(suite.T(), err)
	suite.signer = helpers.NewSigner(priv)

	suite.randSigners = make(map[common.Address]*helpers.Signer)
	suite.randSigners[suite.signer.Address()] = suite.signer
	for i := 0; i < 10; i++ {
		privKey := helpers.NewEthPrivKey()
		suite.randSigners[common.BytesToAddress(privKey.PubKey().Address())] = helpers.NewSigner(privKey)
	}

	set, accs, balances := helpers.GenerateGenesisValidator(100, nil)
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), set, accs, balances...)

	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight(),
		ChainID:         fxtypes.ChainId(),
		ProposerAddress: set.Proposer.Address,
		Time:            time.Now().UTC(),
	})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.Erc20Keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))))
	for addr := range suite.randSigners {
		helpers.AddTestAddr(suite.app, suite.ctx, addr.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))))
	}
}

func (suite *KeeperTestSuite) Commit() {
	suite.app.EndBlock(abci.RequestEndBlock{
		Height: suite.ctx.BlockHeight(),
	})
	suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	header.Time = time.Now().UTC()
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})
	suite.ctx = suite.ctx.WithBlockHeight(header.Height)
}

func (suite *KeeperTestSuite) StateDB() *statedb.StateDB {
	return statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))
}

func (suite *KeeperTestSuite) RandSigner() *helpers.Signer {
	idx := rand.Intn(len(suite.randSigners) - 1)
	i := 0
	for _, signer := range suite.randSigners {
		if i == idx {
			return signer
		}
		i++
	}
	return suite.signer
}

func (suite *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) DeployContract(from common.Address) (common.Address, error) {
	contract, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), "Test token", "TEST", 18)
	suite.Require().NoError(err)

	_, err = suite.app.Erc20Keeper.CallEVM(suite.ctx, fxtypes.GetERC20().ABI, suite.app.Erc20Keeper.ModuleAddress(), contract, true, "transferOwnership", from)
	if err != nil {
		return common.Address{}, err
	}
	return contract, nil
}

func (suite *KeeperTestSuite) DeployFXRelayToken() (types.TokenPair, banktypes.Metadata) {
	fxToken := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)

	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, fxToken)
	suite.Require().NoError(err)
	return *pair, fxToken
}

func (suite *KeeperTestSuite) GenerateCrossChainDenoms() []string {
	modules := []string{"eth", "bsc", "polygon", "tron", "avalanche"}

	count := rand.Intn(len(modules)-1) + 1

	denoms := make([]string, len(modules))
	for index, m := range modules {
		address := helpers.GenerateAddress().String()
		if m == "tron" {
			address = trontypes.AddressFromHex(address)
		}
		denom := fmt.Sprintf("%s%s", m, address)
		denoms[index] = denom
	}
	if count >= len(modules) {
		return denoms
	}
	return denoms[:count]
}

func (suite *KeeperTestSuite) DeployNativeRelayToken(symbol string, denom ...string) (types.TokenPair, banktypes.Metadata) {
	testToken := fxtypes.GetCrossChainMetadata("Test Token", symbol, 18, denom...)

	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, testToken)
	suite.Require().NoError(err)
	return *pair, testToken
}

func (suite *KeeperTestSuite) DeployERC20RelayToken(contract common.Address) (types.TokenPair, banktypes.Metadata) {
	pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contract)
	suite.Require().NoError(err)
	md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.Denom)
	suite.Require().True(found)
	return *pair, md
}

func (suite *KeeperTestSuite) MintLockNativeTokenToModule(md banktypes.Metadata, amt sdk.Int) *big.Int {
	generateAddress := helpers.GenerateAddress()

	count := 1
	if len(md.DenomUnits) > 0 && len(md.DenomUnits[0].Aliases) > 0 {
		// add alias to erc20 module
		for _, alias := range md.DenomUnits[0].Aliases {
			// add alias for erc20 module
			coins := sdk.NewCoins(sdk.NewCoin(alias, amt))
			helpers.AddTestAddr(suite.app, suite.ctx, generateAddress.Bytes(), coins)
			err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, generateAddress.Bytes(), types.ModuleName, coins)
			suite.Require().NoError(err)
		}
		count = len(md.DenomUnits[0].Aliases)
	}

	// add denom to erc20 module
	coin := sdk.NewCoin(md.Base, amt.Mul(sdk.NewInt(int64(count))))
	helpers.AddTestAddr(suite.app, suite.ctx, generateAddress.Bytes(), sdk.NewCoins(coin))
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, generateAddress.Bytes(), types.ModuleName, sdk.NewCoins(coin))
	suite.Require().NoError(err)

	return coin.Amount.BigInt()
}

func (suite *KeeperTestSuite) BalanceOf(contract, account common.Address) *big.Int {
	balance, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, contract, account)
	suite.NoError(err)
	return balance
}

func (suite *KeeperTestSuite) MintERC20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := fxtypes.GetERC20()
	transferData, err := erc20.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) ModuleMintERC20Token(contractAddr, to common.Address, amount *big.Int) {
	erc20 := fxtypes.GetERC20()
	transferData, err := erc20.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	suite.moduleSendEvmTx(contractAddr, transferData)
}

func (suite *KeeperTestSuite) TransferERC20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := fxtypes.GetERC20()
	transferData, err := erc20.ABI.Pack("transfer", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) TransferERC20TokenToModule(contractAddr, from common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := fxtypes.GetERC20()
	moduleAddress := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	transferData, err := erc20.ABI.Pack("transfer", common.BytesToAddress(moduleAddress.Bytes()), amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) TransferERC20TokenToModuleWithoutHook(contractAddr, from common.Address, amount *big.Int) {
	erc20 := fxtypes.GetERC20()
	moduleAddress := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	transferData, err := erc20.ABI.Pack("transfer", common.BytesToAddress(moduleAddress.Bytes()), amount)
	suite.Require().NoError(err)
	_, err = suite.app.EvmKeeper.CallEVMWithData(suite.ctx, from, &contractAddr, transferData, true)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) sendEvmTx(contractAddr, from common.Address, data []byte) *evm.MsgEthereumTx {
	chainID := suite.app.EvmKeeper.ChainID()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contractAddr, From: &from, Data: (*hexutil.Bytes)(&data)})
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	res, err := evm.NewQueryClient(queryHelper).EstimateGas(sdk.WrapSDKContext(suite.ctx),
		&evm.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		},
	)
	suite.Require().NoError(err)

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	// suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64()*int64(res.Gas)))))

	ercTransferTx := evm.NewTx(
		chainID,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		&contractAddr,
		nil,
		res.Gas,
		nil,
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		big.NewInt(1),
		data,
		&ethereumtypes.AccessList{}, // accesses
	)
	signer, ok := suite.randSigners[from]
	suite.Require().True(ok)

	ercTransferTx.From = signer.Address().Hex()
	err = ercTransferTx.Sign(ethereumtypes.LatestSignerForChainID(chainID), signer)
	suite.Require().NoError(err)

	rsp, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return ercTransferTx
}

func (suite *KeeperTestSuite) moduleSendEvmTx(contractAddr common.Address, data []byte) {
	rsp, err := suite.app.EvmKeeper.CallEVMWithData(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), &contractAddr, data, true)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
}

func newMetadata() banktypes.Metadata {
	return banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases: []string{
					fmt.Sprintf("%s%s", bsctypes.ModuleName, helpers.GenerateAddress().String()),
					fmt.Sprintf("%s%s", ethtypes.ModuleName, helpers.GenerateAddress().String()),
					// fmt.Sprintf("%s%s", "ibc/", helpers.GenerateAddress().String()),
				},
			}, {
				Denom:    "USDT",
				Exponent: uint32(18),
			},
		},
		Base:    "usdt",
		Display: "display usdt",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}
}
