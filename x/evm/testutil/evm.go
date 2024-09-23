package testutil

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	fxevmkeeper "github.com/functionx/fx-core/v8/x/evm/keeper"
	fxevmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

type EVMSuite struct {
	*require.Assertions
	ctx          sdk.Context
	evmKeeper    *fxevmkeeper.Keeper
	from         common.Address
	contractAddr common.Address
	signer       *helpers.Signer
	gasPrice     *big.Int
}

func (s *EVMSuite) Init(ass *require.Assertions, ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, signer *helpers.Signer) *EVMSuite {
	s.Assertions = ass
	s.ctx = ctx
	s.evmKeeper = evmKeeper
	s.signer = signer
	return s
}

func (s *EVMSuite) GetContractAddr() *common.Address {
	return &s.contractAddr
}

func (s *EVMSuite) WithContractAddr(addr common.Address) {
	s.contractAddr = addr
}

func (s *EVMSuite) WithGasPrice(gasPrice *big.Int) {
	s.gasPrice = gasPrice
}

func (s *EVMSuite) WithSigner(signer *helpers.Signer) {
	s.signer = signer
}

func (s *EVMSuite) WithFrom(from common.Address) {
	s.from = from
}

func (s *EVMSuite) GetFrom() common.Address {
	from := s.from
	if contract.IsZeroEthAddress(from) && s.signer != nil {
		from = s.signer.Address()
	}
	return from
}

func (s *EVMSuite) HexAddr() common.Address {
	return s.signer.Address()
}

func (s *EVMSuite) AccAddr() sdk.AccAddress {
	return s.signer.AccAddress()
}

func (s *EVMSuite) Call(abi abi.ABI, method string, res interface{}, args ...interface{}) {
	err := s.evmKeeper.QueryContract(s.ctx, s.GetFrom(), s.contractAddr, abi, method, res, args...)
	s.NoError(err)
}

func (s *EVMSuite) CallEVM(data []byte, gasLimit uint64) *evmtypes.MsgEthereumTxResponse {
	tx, err := s.evmKeeper.CallEVM(s.ctx, s.GetFrom(), &s.contractAddr, nil, gasLimit, data, false)
	s.NoError(err)
	return tx
}

func (s *EVMSuite) Send(abi abi.ABI, method string, args ...interface{}) *evmtypes.MsgEthereumTxResponse {
	response, err := s.evmKeeper.ApplyContract(s.ctx, s.signer.Address(), s.contractAddr, nil, abi, method, args...)
	s.NoError(err)
	return response
}

func (s *EVMSuite) EthereumTx(to *common.Address, data []byte, value *big.Int, gasLimit uint64) (*evmtypes.MsgEthereumTxResponse, error) {
	chanId := s.evmKeeper.ChainID()
	s.Equal(fxtypes.EIP155ChainID(s.ctx.ChainID()), chanId)
	if value == nil {
		value = big.NewInt(0)
	}

	nonce := s.evmKeeper.GetNonce(s.ctx, s.signer.Address())
	tx := evmtypes.NewTx(
		chanId,
		nonce,
		to,
		value,
		gasLimit,
		s.gasPrice,
		nil,
		nil,
		data,
		nil,
	)
	tx.From = s.signer.Address().Bytes()
	s.NoError(tx.Sign(ethtypes.LatestSignerForChainID(chanId), s.signer))

	return s.evmKeeper.EthereumTx(s.ctx, tx)
}

func (s *EVMSuite) DeployUpgradableERC20Logic(symbol string) common.Address {
	erc20Contract := contract.GetFIP20()
	erc20ModuleAddress := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
	initializeArgs := []interface{}{symbol + " Token", symbol, uint8(18), erc20ModuleAddress}
	newContractAddr, err := s.evmKeeper.DeployUpgradableContract(s.ctx,
		s.signer.Address(), erc20Contract.Address, nil, &erc20Contract.ABI, initializeArgs...)
	s.NoError(err)
	s.contractAddr = newContractAddr
	return newContractAddr
}

func (s *EVMSuite) CallContract(data []byte) error {
	msg := &fxevmtypes.MsgCallContract{
		Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ContractAddress: s.contractAddr.String(),
		Data:            common.Bytes2Hex(data),
	}
	_, err := s.evmKeeper.CallContract(s.ctx, msg)
	return err
}
