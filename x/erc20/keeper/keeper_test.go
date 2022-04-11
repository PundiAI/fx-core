package keeper_test

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/functionx/fx-core/contracts"
	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
	"github.com/functionx/fx-core/x/evm/statedb"
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	app "github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/server/config"
	"github.com/functionx/fx-core/tests"
	fxcoretypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/erc20/types"
	evm "github.com/functionx/fx-core/x/evm/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	app              *app.App
	queryClientEvm   evm.QueryClient
	queryClient      types.QueryClient
	address          common.Address
	priKey           cryptotypes.PrivKey
	consAddress      sdk.ConsAddress
	clientCtx        client.Context
	ethSigner        ethtypes.Signer
	signer           keyring.Signer
	mintFeeCollector bool
	dynamicTxFee     bool
}

// Test helpers
func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	checkTx := false
	suite.mintFeeCollector = true
	suite.dynamicTxFee = true

	// account key
	priKey := NewPriKey()
	//ethsecp256k1.GenerateKey()
	ethPriv := &ethsecp256k1.PrivKey{Key: priKey.Bytes()}
	suite.priKey = ethPriv
	suite.signer = tests.NewSigner(ethPriv)
	suite.address = common.BytesToAddress(suite.priKey.PubKey().Address())

	// consensus key
	priv := NewPriKey()
	suite.consAddress = sdk.ConsAddress(priv.PubKey().Address())

	// setup feemarketGenesis params
	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	feemarketGenesis.Params.EnableHeight = 1
	feemarketGenesis.Params.NoBaseFee = false
	suite.app = app.Setup(checkTx, nil)

	if suite.mintFeeCollector {
		// mint some coin to fee collector
		coins := sdk.NewCoins(sdk.NewCoin(evm.DefaultEVMDenom, sdk.NewInt(int64(params.TxGas)-1)))
		genesisState := app.ModuleBasics.DefaultGenesis(suite.app.AppCodec())
		balances := []banktypes.Balance{
			{
				Address: suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName).String(),
				Coins:   coins,
			},
		}
		// update total supply
		bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, sdk.NewCoins(sdk.NewCoin(evm.DefaultEVMDenom, sdk.NewInt((int64(params.TxGas)-1)))), []banktypes.Metadata{})
		bz := suite.app.AppCodec().MustMarshalJSON(bankGenesis)
		require.NotNil(t, bz)
		genesisState[banktypes.ModuleName] = suite.app.AppCodec().MustMarshalJSON(bankGenesis)

		// we marshal the genesisState of all module to a byte array
		stateBytes, err := tmjson.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		//Initialize the chain
		suite.app.InitChain(
			abci.RequestInitChain{
				InitialHeight:   fxcoretypes.EvmSupportBlock(),
				ChainId:         fxcoretypes.EIP155ChainID().String(),
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{
		Height:          fxcoretypes.EvmSupportBlock(),
		ChainID:         fxcoretypes.EIP155ChainID().String(),
		Time:            time.Now().UTC(),
		ProposerAddress: suite.consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	require.NoError(suite.T(), InitEvmModuleParams(suite.ctx, suite.app.Erc20Keeper, suite.dynamicTxFee))
	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.Erc20Keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	acc := authtypes.NewBaseAccount(suite.address.Bytes(), nil, 0, 0)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	suite.app.EvmKeeper.SetAddressCode(suite.ctx, suite.address, common.BytesToHash(crypto.Keccak256(nil)).Bytes())

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, priv.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)
	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	encodingConfig := app.MakeEncodingConfig()
	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	suite.ethSigner = ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
}

func (suite *KeeperTestSuite) SetupTest() {
	fxcoretypes.ChangeNetworkForTest(fxcoretypes.NetworkDevnet())
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

func (suite *KeeperTestSuite) DeployContract(from common.Address, name string, symbol string, decimals uint8) common.Address {
	contractAddress, err := suite.app.Erc20Keeper.DeployTokenUpgrade(suite.ctx, from, name, symbol, decimals, false)
	suite.Require().NoError(err)
	return contractAddress
}

func (suite *KeeperTestSuite) DeployContractDirectBalanceManipulation(name string, symbol string) common.Address {
	ctx := sdk.WrapSDKContext(suite.ctx)
	chainID := suite.app.EvmKeeper.ChainID()

	erc20Config := contracts.GetERC20(suite.ctx.BlockHeight())

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

	fip20DeployTx := evm.NewTxContract(
		chainID,
		nonce,
		nil,     // amount
		res.Gas, // gasLimit
		nil,     // gasPrice
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		big.NewInt(1),
		data,                   // input
		&ethtypes.AccessList{}, // accesses
	)

	fip20DeployTx.From = suite.address.Hex()
	err = fip20DeployTx.Sign(ethtypes.LatestSignerForChainID(chainID), suite.signer)
	suite.Require().NoError(err)
	rsp, err := suite.app.EvmKeeper.EthereumTx(ctx, fip20DeployTx)
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

func (suite *KeeperTestSuite) MintFIP20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := contracts.GetERC20(suite.ctx.BlockHeight())
	transferData, err := erc20.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) BurnFIP20Token(contractAddr, from common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := contracts.GetERC20(suite.ctx.BlockHeight())
	transferData, err := erc20.ABI.Pack("transfer", types.ModuleAddress, amount)
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

	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, suite.address.Bytes())
	suite.Require().NoError(err)

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(evm.DefaultEVMDenom, sdk.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64()*int64(res.Gas)))))

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
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = suite.address.Hex()
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), suite.signer)
	suite.Require().NoError(err)
	rsp, err := suite.app.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return ercTransferTx
}

func (suite *KeeperTestSuite) BalanceOf(contract, account common.Address) interface{} {
	erc20 := contracts.GetERC20(suite.ctx.BlockHeight())

	res, err := suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20.ABI, types.ModuleAddress, contract, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := erc20.ABI.Unpack("balanceOf", res.Ret)
	if len(unpacked) == 0 {
		return nil
	}

	return unpacked[0]
}

func (suite *KeeperTestSuite) NameOf(contract common.Address) string {
	erc20 := contracts.GetERC20(suite.ctx.BlockHeight())
	res, err := suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20.ABI, types.ModuleAddress, contract, "name")
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	unpacked, err := erc20.ABI.Unpack("name", res.Ret)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(unpacked)

	return fmt.Sprintf("%v", unpacked[0])
}

func (suite *KeeperTestSuite) TransferFIP20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := contracts.GetERC20(suite.ctx.BlockHeight())
	transferData, err := erc20.ABI.Pack("transfer", to, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func InitEvmModuleParams(ctx sdk.Context, keeper erc20keeper.Keeper, dynamicTxFee bool) error {
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + fxcoretypes.EvmSupportBlock())
	defaultEvmParams := evmtypes.DefaultParams()
	defaultFeeMarketParams := feemarkettypes.DefaultParams()
	defaultErc20Params := types.DefaultParams()

	if dynamicTxFee {
		defaultFeeMarketParams.EnableHeight = fxcoretypes.EvmSupportBlock()
		defaultFeeMarketParams.NoBaseFee = false
	} else {
		defaultFeeMarketParams.NoBaseFee = true
	}

	if err := keeper.HandleInitEvmProposal(ctx, defaultErc20Params, defaultFeeMarketParams, defaultEvmParams, nil); err != nil {
		return err
	}
	return nil
}

func NewPriKey() cryptotypes.PrivKey {
	return secp256k1.GenPrivKey()
}
