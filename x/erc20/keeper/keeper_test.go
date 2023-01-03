package keeper_test

import (
	"encoding/json"
	"fmt"
	"math/big"
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
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.App
	queryClient types.QueryClient
	signer      *helpers.Signer
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

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))
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

func (suite *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) DeployContract(from common.Address) (common.Address, error) {
	return suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, from, "Test token", "TEST", 18)
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
	//suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64()*int64(res.Gas)))))

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

	ercTransferTx.From = suite.signer.Address().Hex()
	err = ercTransferTx.Sign(ethereumtypes.LatestSignerForChainID(chainID), suite.signer)
	suite.Require().NoError(err)

	rsp, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return ercTransferTx
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
					//fmt.Sprintf("%s%s", "ibc/", helpers.GenerateAddress().String()),
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
