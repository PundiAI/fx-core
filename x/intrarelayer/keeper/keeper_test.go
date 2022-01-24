package keeper_test

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
	"github.com/functionx/fx-core/x/intrarelayer/keeper"
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
	ethermint "github.com/functionx/fx-core/types"
	evm "github.com/functionx/fx-core/x/evm/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
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
	feemarketGenesis.BaseFee = sdk.NewInt(feemarketGenesis.Params.InitialBaseFee)
	suite.app = app.Setup(checkTx)

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
				ChainId:         "evmos_9000-1",
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{
		Height:          1,
		ChainID:         ethermint.EIP155ChainID().String(),
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
	suite.app.EvmKeeper.WithContext(suite.ctx)

	require.NoError(suite.T(), InitEvmModuleParams(suite.ctx, suite.app.EvmKeeper, suite.dynamicTxFee))
	require.NoError(suite.T(), InitIntrarelayerParams(suite.ctx, suite.app.IntrarelayerKeeper))
	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.IntrarelayerKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	//TODO update ethAccount 2021-12-02.
	//acc := &ethermint.EthAccount{
	//	BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
	//	CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	//}

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
	suite.DoSetupTest(suite.T())
}

func (suite *KeeperTestSuite) DeployContract(name string, symbol string, decimals uint8) common.Address {
	ctx := sdk.WrapSDKContext(suite.ctx)
	chainID := suite.app.EvmKeeper.ChainID()

	ctorArgs, err := contracts.ERC20RelayContract.ABI.Pack("", name, symbol, decimals)
	suite.Require().NoError(err)

	data := append(contracts.ERC20RelayContract.Bin, ctorArgs...)
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

	nonce := suite.app.EvmKeeper.GetNonce(suite.address)

	erc20DeployTx := evm.NewTxContract(
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

	erc20DeployTx.From = suite.address.Hex()
	err = erc20DeployTx.Sign(ethtypes.LatestSignerForChainID(chainID), suite.signer)
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
	suite.app.EvmKeeper.WithContext(suite.ctx)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) MintERC20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	transferData, err := contracts.ERC20RelayContract.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func (suite *KeeperTestSuite) BurnERC20Token(contractAddr, from common.Address, amount *big.Int) *evm.MsgEthereumTx {
	transferData, err := contracts.ERC20RelayContract.ABI.Pack("transfer", types.ModuleAddress, amount)
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

	nonce := suite.app.EvmKeeper.GetNonce(from)

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
	erc20 := contracts.ERC20RelayContract.ABI

	res, err := suite.app.IntrarelayerKeeper.CallEVMWithModule(suite.ctx, erc20, contract, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := erc20.Unpack("balanceOf", res.Ret)
	if len(unpacked) == 0 {
		return nil
	}

	return unpacked[0]
}

func (suite *KeeperTestSuite) NameOf(contract common.Address) interface{} {

	erc20 := contracts.ERC20RelayContract.ABI

	res, err := suite.app.IntrarelayerKeeper.CallEVMWithModule(suite.ctx, erc20, contract, "name")
	if err != nil {
		return nil
	}

	unpacked, err := erc20.Unpack("name", res.Ret)
	if len(unpacked) == 0 {
		return nil
	}

	return unpacked[0]
}

func (suite *KeeperTestSuite) TransferERC20Token(contractAddr, from, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	transferData, err := contracts.ERC20RelayContract.ABI.Pack("transfer", to, amount)
	suite.Require().NoError(err)
	return suite.sendTx(contractAddr, from, transferData)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func InitEvmModuleParams(ctx sdk.Context, keeper *evmkeeper.Keeper, dynamicTxFee bool) error {
	defaultEvmParams := evmtypes.DefaultParams()
	defaultFeeMarketParams := feemarkettypes.DefaultParams()

	if dynamicTxFee {
		defaultFeeMarketParams.EnableHeight = 1
		defaultFeeMarketParams.NoBaseFee = false
	} else {
		defaultFeeMarketParams.NoBaseFee = true
	}

	if err := keeper.HandleInitEvmParamsProposal(ctx, &evmtypes.InitEvmParamsProposal{
		Title:           "Init evm title",
		Description:     "Init emv module description",
		EvmParams:       &defaultEvmParams,
		FeemarketParams: &defaultFeeMarketParams,
	}); err != nil {
		return err
	}
	keeper.WithChainID(ctx)
	return nil
}

func InitIntrarelayerParams(ctx sdk.Context, keeper keeper.Keeper) error {
	defaultParams := types.DefaultParams()

	err := keeper.InitIntrarelayer(ctx, &types.InitIntrarelayerParamsProposal{
		Title:       "Init intrarelayer title",
		Description: "Init intrarelayer module description",
		Params:      &defaultParams,
	})

	return err
}

func EvmKeeperSetHook(evm *evmkeeper.Keeper, hooks evmkeeper.MultiEvmHooks) {
	evm.SetHooks(hooks)
}

func NewPriKey() cryptotypes.PrivKey {
	return secp256k1.GenPrivKey()
}
