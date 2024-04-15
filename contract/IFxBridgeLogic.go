// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IFxBridgeLogicBridgeCallData is an auto generated low-level Go binding around an user-defined struct.
type IFxBridgeLogicBridgeCallData struct {
	Sender   common.Address
	Receiver common.Address
	To       common.Address
	Tokens   []common.Address
	Amounts  []*big.Int
	Message  []byte
	Value    *big.Int
	Timeout  *big.Int
	GasLimit *big.Int
}

// IFxBridgeLogicBridgeToken is an auto generated low-level Go binding around an user-defined struct.
type IFxBridgeLogicBridgeToken struct {
	Addr      common.Address
	Name      string
	Symbol    string
	Decimals  uint8
	TokenType uint8
}

// IFxBridgeLogicTokenStatus is an auto generated low-level Go binding around an user-defined struct.
type IFxBridgeLogicTokenStatus struct {
	IsOriginated bool
	IsActive     bool
	IsExist      bool
	TokenType    uint8
}

// IFxBridgeLogicMetaData contains all meta data concerning the IFxBridgeLogic contract.
var IFxBridgeLogicMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_tokenContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_channelIBC\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumIFxBridgeLogic.BridgeTokenType\",\"name\":\"_tokenType\",\"type\":\"uint8\"}],\"name\":\"AddBridgeTokenEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_dstChainId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"BridgeCallEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_newOracleSetNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"_powers\",\"type\":\"uint256[]\"}],\"name\":\"OracleSetUpdatedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_refundNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"RefundTokenExecutedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_tokenContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_targetIBC\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"SendToFxEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"name\":\"SubmitBridgeCallEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_batchNonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"TransactionBatchExecutedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"TransferOwnerEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddr\",\"type\":\"address\"}],\"name\":\"activeBridgeToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_channelIBC\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"_isOriginated\",\"type\":\"bool\"},{\"internalType\":\"enumIFxBridgeLogic.BridgeTokenType\",\"name\":\"_tokenType\",\"type\":\"uint8\"}],\"name\":\"addBridgeToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChainId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"bridgeCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_fxbridgeId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_methodName\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timeout\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_asset\",\"type\":\"bytes\"}],\"name\":\"bridgeCallCheckpoint\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddr\",\"type\":\"address\"}],\"name\":\"checkAssetStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_currentOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_currentPowers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_s\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_theHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_powerThreshold\",\"type\":\"uint256\"}],\"name\":\"checkOracleSignatures\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_erc20Address\",\"type\":\"address\"}],\"name\":\"convert_decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeTokenList\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"enumIFxBridgeLogic.BridgeTokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"}],\"internalType\":\"structIFxBridgeLogic.BridgeToken[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_erc20Address\",\"type\":\"address\"}],\"name\":\"lastBatchNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_powers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_oracleSetNonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_fxBridgeId\",\"type\":\"bytes32\"}],\"name\":\"makeCheckpoint\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_fxbridgeId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_methodName\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_oracleSetNonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_powers\",\"type\":\"uint256[]\"}],\"name\":\"oracleSetCheckpoint\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddr\",\"type\":\"address\"}],\"name\":\"pauseBridgeToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_currentOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_currentPowers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_s\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2]\",\"name\":\"_nonceArray\",\"type\":\"uint256[2]\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_timeout\",\"type\":\"uint256\"}],\"name\":\"refundBridgeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_destination\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_targetIBC\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"sendToFx\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"state_fxBridgeId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_erc20Address\",\"type\":\"address\"}],\"name\":\"state_lastBatchNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"state_lastEventNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"state_lastOracleSetCheckpoint\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"state_lastOracleSetNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"state_lastRefundNonce\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"state_powerThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_currentOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_currentPowers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_s\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"_destinations\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[2]\",\"name\":\"_nonceArray\",\"type\":\"uint256[2]\"},{\"internalType\":\"address\",\"name\":\"_tokenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_batchTimeout\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_feeReceive\",\"type\":\"address\"}],\"name\":\"submitBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_fxbridgeId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_methodName\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"_destinations\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_batchNonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_tokenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_batchTimeout\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_feeReceive\",\"type\":\"address\"}],\"name\":\"submitBatchCheckpoint\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_currentOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_currentPowers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_s\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2]\",\"name\":\"_nonceArray\",\"type\":\"uint256[2]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structIFxBridgeLogic.BridgeCallData\",\"name\":\"_input\",\"type\":\"tuple\"}],\"name\":\"submitBridgeCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddr\",\"type\":\"address\"}],\"name\":\"tokenStatus\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isOriginated\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isExist\",\"type\":\"bool\"},{\"internalType\":\"enumIFxBridgeLogic.BridgeTokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"}],\"internalType\":\"structIFxBridgeLogic.TokenStatus\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_newPowers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_newOracleSetNonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"_currentOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_currentPowers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_currentOracleSetNonce\",\"type\":\"uint256\"},{\"internalType\":\"uint8[]\",\"name\":\"_v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_s\",\"type\":\"bytes32[]\"}],\"name\":\"updateOracleSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IFxBridgeLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use IFxBridgeLogicMetaData.ABI instead.
var IFxBridgeLogicABI = IFxBridgeLogicMetaData.ABI

// IFxBridgeLogic is an auto generated Go binding around an Ethereum contract.
type IFxBridgeLogic struct {
	IFxBridgeLogicCaller     // Read-only binding to the contract
	IFxBridgeLogicTransactor // Write-only binding to the contract
	IFxBridgeLogicFilterer   // Log filterer for contract events
}

// IFxBridgeLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type IFxBridgeLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFxBridgeLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IFxBridgeLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFxBridgeLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IFxBridgeLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFxBridgeLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IFxBridgeLogicSession struct {
	Contract     *IFxBridgeLogic   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IFxBridgeLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IFxBridgeLogicCallerSession struct {
	Contract *IFxBridgeLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IFxBridgeLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IFxBridgeLogicTransactorSession struct {
	Contract     *IFxBridgeLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IFxBridgeLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type IFxBridgeLogicRaw struct {
	Contract *IFxBridgeLogic // Generic contract binding to access the raw methods on
}

// IFxBridgeLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IFxBridgeLogicCallerRaw struct {
	Contract *IFxBridgeLogicCaller // Generic read-only contract binding to access the raw methods on
}

// IFxBridgeLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IFxBridgeLogicTransactorRaw struct {
	Contract *IFxBridgeLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIFxBridgeLogic creates a new instance of IFxBridgeLogic, bound to a specific deployed contract.
func NewIFxBridgeLogic(address common.Address, backend bind.ContractBackend) (*IFxBridgeLogic, error) {
	contract, err := bindIFxBridgeLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogic{IFxBridgeLogicCaller: IFxBridgeLogicCaller{contract: contract}, IFxBridgeLogicTransactor: IFxBridgeLogicTransactor{contract: contract}, IFxBridgeLogicFilterer: IFxBridgeLogicFilterer{contract: contract}}, nil
}

// NewIFxBridgeLogicCaller creates a new read-only instance of IFxBridgeLogic, bound to a specific deployed contract.
func NewIFxBridgeLogicCaller(address common.Address, caller bind.ContractCaller) (*IFxBridgeLogicCaller, error) {
	contract, err := bindIFxBridgeLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicCaller{contract: contract}, nil
}

// NewIFxBridgeLogicTransactor creates a new write-only instance of IFxBridgeLogic, bound to a specific deployed contract.
func NewIFxBridgeLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*IFxBridgeLogicTransactor, error) {
	contract, err := bindIFxBridgeLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicTransactor{contract: contract}, nil
}

// NewIFxBridgeLogicFilterer creates a new log filterer instance of IFxBridgeLogic, bound to a specific deployed contract.
func NewIFxBridgeLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*IFxBridgeLogicFilterer, error) {
	contract, err := bindIFxBridgeLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicFilterer{contract: contract}, nil
}

// bindIFxBridgeLogic binds a generic wrapper to an already deployed contract.
func bindIFxBridgeLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IFxBridgeLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IFxBridgeLogic *IFxBridgeLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IFxBridgeLogic.Contract.IFxBridgeLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IFxBridgeLogic *IFxBridgeLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.IFxBridgeLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IFxBridgeLogic *IFxBridgeLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.IFxBridgeLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IFxBridgeLogic *IFxBridgeLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IFxBridgeLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IFxBridgeLogic *IFxBridgeLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IFxBridgeLogic *IFxBridgeLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.contract.Transact(opts, method, params...)
}

// BridgeTokens is a free data retrieval call binding the contract method 0xf8a06888.
//
// Solidity: function bridgeTokens() view returns(address[])
func (_IFxBridgeLogic *IFxBridgeLogicCaller) BridgeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "bridgeTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// BridgeTokens is a free data retrieval call binding the contract method 0xf8a06888.
//
// Solidity: function bridgeTokens() view returns(address[])
func (_IFxBridgeLogic *IFxBridgeLogicSession) BridgeTokens() ([]common.Address, error) {
	return _IFxBridgeLogic.Contract.BridgeTokens(&_IFxBridgeLogic.CallOpts)
}

// BridgeTokens is a free data retrieval call binding the contract method 0xf8a06888.
//
// Solidity: function bridgeTokens() view returns(address[])
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) BridgeTokens() ([]common.Address, error) {
	return _IFxBridgeLogic.Contract.BridgeTokens(&_IFxBridgeLogic.CallOpts)
}

// CheckAssetStatus is a free data retrieval call binding the contract method 0x474d561c.
//
// Solidity: function checkAssetStatus(address _tokenAddr) view returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) CheckAssetStatus(opts *bind.CallOpts, _tokenAddr common.Address) (bool, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "checkAssetStatus", _tokenAddr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckAssetStatus is a free data retrieval call binding the contract method 0x474d561c.
//
// Solidity: function checkAssetStatus(address _tokenAddr) view returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicSession) CheckAssetStatus(_tokenAddr common.Address) (bool, error) {
	return _IFxBridgeLogic.Contract.CheckAssetStatus(&_IFxBridgeLogic.CallOpts, _tokenAddr)
}

// CheckAssetStatus is a free data retrieval call binding the contract method 0x474d561c.
//
// Solidity: function checkAssetStatus(address _tokenAddr) view returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) CheckAssetStatus(_tokenAddr common.Address) (bool, error) {
	return _IFxBridgeLogic.Contract.CheckAssetStatus(&_IFxBridgeLogic.CallOpts, _tokenAddr)
}

// CheckOracleSignatures is a free data retrieval call binding the contract method 0x285a190a.
//
// Solidity: function checkOracleSignatures(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, bytes32 _theHash, uint256 _powerThreshold) pure returns()
func (_IFxBridgeLogic *IFxBridgeLogicCaller) CheckOracleSignatures(opts *bind.CallOpts, _currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _theHash [32]byte, _powerThreshold *big.Int) error {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "checkOracleSignatures", _currentOracles, _currentPowers, _v, _r, _s, _theHash, _powerThreshold)

	if err != nil {
		return err
	}

	return err

}

// CheckOracleSignatures is a free data retrieval call binding the contract method 0x285a190a.
//
// Solidity: function checkOracleSignatures(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, bytes32 _theHash, uint256 _powerThreshold) pure returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) CheckOracleSignatures(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _theHash [32]byte, _powerThreshold *big.Int) error {
	return _IFxBridgeLogic.Contract.CheckOracleSignatures(&_IFxBridgeLogic.CallOpts, _currentOracles, _currentPowers, _v, _r, _s, _theHash, _powerThreshold)
}

// CheckOracleSignatures is a free data retrieval call binding the contract method 0x285a190a.
//
// Solidity: function checkOracleSignatures(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, bytes32 _theHash, uint256 _powerThreshold) pure returns()
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) CheckOracleSignatures(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _theHash [32]byte, _powerThreshold *big.Int) error {
	return _IFxBridgeLogic.Contract.CheckOracleSignatures(&_IFxBridgeLogic.CallOpts, _currentOracles, _currentPowers, _v, _r, _s, _theHash, _powerThreshold)
}

// ConvertDecimals is a free data retrieval call binding the contract method 0x7d9a8ea6.
//
// Solidity: function convert_decimals(address _erc20Address) view returns(uint8)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) ConvertDecimals(opts *bind.CallOpts, _erc20Address common.Address) (uint8, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "convert_decimals", _erc20Address)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ConvertDecimals is a free data retrieval call binding the contract method 0x7d9a8ea6.
//
// Solidity: function convert_decimals(address _erc20Address) view returns(uint8)
func (_IFxBridgeLogic *IFxBridgeLogicSession) ConvertDecimals(_erc20Address common.Address) (uint8, error) {
	return _IFxBridgeLogic.Contract.ConvertDecimals(&_IFxBridgeLogic.CallOpts, _erc20Address)
}

// ConvertDecimals is a free data retrieval call binding the contract method 0x7d9a8ea6.
//
// Solidity: function convert_decimals(address _erc20Address) view returns(uint8)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) ConvertDecimals(_erc20Address common.Address) (uint8, error) {
	return _IFxBridgeLogic.Contract.ConvertDecimals(&_IFxBridgeLogic.CallOpts, _erc20Address)
}

// GetBridgeTokenList is a free data retrieval call binding the contract method 0x283040b4.
//
// Solidity: function getBridgeTokenList() view returns((address,string,string,uint8,uint8)[])
func (_IFxBridgeLogic *IFxBridgeLogicCaller) GetBridgeTokenList(opts *bind.CallOpts) ([]IFxBridgeLogicBridgeToken, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "getBridgeTokenList")

	if err != nil {
		return *new([]IFxBridgeLogicBridgeToken), err
	}

	out0 := *abi.ConvertType(out[0], new([]IFxBridgeLogicBridgeToken)).(*[]IFxBridgeLogicBridgeToken)

	return out0, err

}

// GetBridgeTokenList is a free data retrieval call binding the contract method 0x283040b4.
//
// Solidity: function getBridgeTokenList() view returns((address,string,string,uint8,uint8)[])
func (_IFxBridgeLogic *IFxBridgeLogicSession) GetBridgeTokenList() ([]IFxBridgeLogicBridgeToken, error) {
	return _IFxBridgeLogic.Contract.GetBridgeTokenList(&_IFxBridgeLogic.CallOpts)
}

// GetBridgeTokenList is a free data retrieval call binding the contract method 0x283040b4.
//
// Solidity: function getBridgeTokenList() view returns((address,string,string,uint8,uint8)[])
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) GetBridgeTokenList() ([]IFxBridgeLogicBridgeToken, error) {
	return _IFxBridgeLogic.Contract.GetBridgeTokenList(&_IFxBridgeLogic.CallOpts)
}

// LastBatchNonce is a free data retrieval call binding the contract method 0x011b2174.
//
// Solidity: function lastBatchNonce(address _erc20Address) view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) LastBatchNonce(opts *bind.CallOpts, _erc20Address common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "lastBatchNonce", _erc20Address)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBatchNonce is a free data retrieval call binding the contract method 0x011b2174.
//
// Solidity: function lastBatchNonce(address _erc20Address) view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicSession) LastBatchNonce(_erc20Address common.Address) (*big.Int, error) {
	return _IFxBridgeLogic.Contract.LastBatchNonce(&_IFxBridgeLogic.CallOpts, _erc20Address)
}

// LastBatchNonce is a free data retrieval call binding the contract method 0x011b2174.
//
// Solidity: function lastBatchNonce(address _erc20Address) view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) LastBatchNonce(_erc20Address common.Address) (*big.Int, error) {
	return _IFxBridgeLogic.Contract.LastBatchNonce(&_IFxBridgeLogic.CallOpts, _erc20Address)
}

// MakeCheckpoint is a free data retrieval call binding the contract method 0x71cbf381.
//
// Solidity: function makeCheckpoint(address[] _oracles, uint256[] _powers, uint256 _oracleSetNonce, bytes32 _fxBridgeId) pure returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) MakeCheckpoint(opts *bind.CallOpts, _oracles []common.Address, _powers []*big.Int, _oracleSetNonce *big.Int, _fxBridgeId [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "makeCheckpoint", _oracles, _powers, _oracleSetNonce, _fxBridgeId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MakeCheckpoint is a free data retrieval call binding the contract method 0x71cbf381.
//
// Solidity: function makeCheckpoint(address[] _oracles, uint256[] _powers, uint256 _oracleSetNonce, bytes32 _fxBridgeId) pure returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicSession) MakeCheckpoint(_oracles []common.Address, _powers []*big.Int, _oracleSetNonce *big.Int, _fxBridgeId [32]byte) ([32]byte, error) {
	return _IFxBridgeLogic.Contract.MakeCheckpoint(&_IFxBridgeLogic.CallOpts, _oracles, _powers, _oracleSetNonce, _fxBridgeId)
}

// MakeCheckpoint is a free data retrieval call binding the contract method 0x71cbf381.
//
// Solidity: function makeCheckpoint(address[] _oracles, uint256[] _powers, uint256 _oracleSetNonce, bytes32 _fxBridgeId) pure returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) MakeCheckpoint(_oracles []common.Address, _powers []*big.Int, _oracleSetNonce *big.Int, _fxBridgeId [32]byte) ([32]byte, error) {
	return _IFxBridgeLogic.Contract.MakeCheckpoint(&_IFxBridgeLogic.CallOpts, _oracles, _powers, _oracleSetNonce, _fxBridgeId)
}

// StateFxBridgeId is a free data retrieval call binding the contract method 0xf92367fd.
//
// Solidity: function state_fxBridgeId() view returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StateFxBridgeId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_fxBridgeId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateFxBridgeId is a free data retrieval call binding the contract method 0xf92367fd.
//
// Solidity: function state_fxBridgeId() view returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StateFxBridgeId() ([32]byte, error) {
	return _IFxBridgeLogic.Contract.StateFxBridgeId(&_IFxBridgeLogic.CallOpts)
}

// StateFxBridgeId is a free data retrieval call binding the contract method 0xf92367fd.
//
// Solidity: function state_fxBridgeId() view returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StateFxBridgeId() ([32]byte, error) {
	return _IFxBridgeLogic.Contract.StateFxBridgeId(&_IFxBridgeLogic.CallOpts)
}

// StateLastBatchNonces is a free data retrieval call binding the contract method 0xdf97174b.
//
// Solidity: function state_lastBatchNonces(address _erc20Address) view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StateLastBatchNonces(opts *bind.CallOpts, _erc20Address common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_lastBatchNonces", _erc20Address)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StateLastBatchNonces is a free data retrieval call binding the contract method 0xdf97174b.
//
// Solidity: function state_lastBatchNonces(address _erc20Address) view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StateLastBatchNonces(_erc20Address common.Address) (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StateLastBatchNonces(&_IFxBridgeLogic.CallOpts, _erc20Address)
}

// StateLastBatchNonces is a free data retrieval call binding the contract method 0xdf97174b.
//
// Solidity: function state_lastBatchNonces(address _erc20Address) view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StateLastBatchNonces(_erc20Address common.Address) (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StateLastBatchNonces(&_IFxBridgeLogic.CallOpts, _erc20Address)
}

// StateLastEventNonce is a free data retrieval call binding the contract method 0x73b20547.
//
// Solidity: function state_lastEventNonce() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StateLastEventNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_lastEventNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StateLastEventNonce is a free data retrieval call binding the contract method 0x73b20547.
//
// Solidity: function state_lastEventNonce() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StateLastEventNonce() (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StateLastEventNonce(&_IFxBridgeLogic.CallOpts)
}

// StateLastEventNonce is a free data retrieval call binding the contract method 0x73b20547.
//
// Solidity: function state_lastEventNonce() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StateLastEventNonce() (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StateLastEventNonce(&_IFxBridgeLogic.CallOpts)
}

// StateLastOracleSetCheckpoint is a free data retrieval call binding the contract method 0x70a0eb94.
//
// Solidity: function state_lastOracleSetCheckpoint() view returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StateLastOracleSetCheckpoint(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_lastOracleSetCheckpoint")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateLastOracleSetCheckpoint is a free data retrieval call binding the contract method 0x70a0eb94.
//
// Solidity: function state_lastOracleSetCheckpoint() view returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StateLastOracleSetCheckpoint() ([32]byte, error) {
	return _IFxBridgeLogic.Contract.StateLastOracleSetCheckpoint(&_IFxBridgeLogic.CallOpts)
}

// StateLastOracleSetCheckpoint is a free data retrieval call binding the contract method 0x70a0eb94.
//
// Solidity: function state_lastOracleSetCheckpoint() view returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StateLastOracleSetCheckpoint() ([32]byte, error) {
	return _IFxBridgeLogic.Contract.StateLastOracleSetCheckpoint(&_IFxBridgeLogic.CallOpts)
}

// StateLastOracleSetNonce is a free data retrieval call binding the contract method 0xbb83bf96.
//
// Solidity: function state_lastOracleSetNonce() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StateLastOracleSetNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_lastOracleSetNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StateLastOracleSetNonce is a free data retrieval call binding the contract method 0xbb83bf96.
//
// Solidity: function state_lastOracleSetNonce() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StateLastOracleSetNonce() (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StateLastOracleSetNonce(&_IFxBridgeLogic.CallOpts)
}

// StateLastOracleSetNonce is a free data retrieval call binding the contract method 0xbb83bf96.
//
// Solidity: function state_lastOracleSetNonce() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StateLastOracleSetNonce() (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StateLastOracleSetNonce(&_IFxBridgeLogic.CallOpts)
}

// StateLastRefundNonce is a free data retrieval call binding the contract method 0x0fa4f599.
//
// Solidity: function state_lastRefundNonce(uint256 _nonce) view returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StateLastRefundNonce(opts *bind.CallOpts, _nonce *big.Int) (bool, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_lastRefundNonce", _nonce)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StateLastRefundNonce is a free data retrieval call binding the contract method 0x0fa4f599.
//
// Solidity: function state_lastRefundNonce(uint256 _nonce) view returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StateLastRefundNonce(_nonce *big.Int) (bool, error) {
	return _IFxBridgeLogic.Contract.StateLastRefundNonce(&_IFxBridgeLogic.CallOpts, _nonce)
}

// StateLastRefundNonce is a free data retrieval call binding the contract method 0x0fa4f599.
//
// Solidity: function state_lastRefundNonce(uint256 _nonce) view returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StateLastRefundNonce(_nonce *big.Int) (bool, error) {
	return _IFxBridgeLogic.Contract.StateLastRefundNonce(&_IFxBridgeLogic.CallOpts, _nonce)
}

// StatePowerThreshold is a free data retrieval call binding the contract method 0xe5a2b5d2.
//
// Solidity: function state_powerThreshold() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) StatePowerThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "state_powerThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StatePowerThreshold is a free data retrieval call binding the contract method 0xe5a2b5d2.
//
// Solidity: function state_powerThreshold() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicSession) StatePowerThreshold() (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StatePowerThreshold(&_IFxBridgeLogic.CallOpts)
}

// StatePowerThreshold is a free data retrieval call binding the contract method 0xe5a2b5d2.
//
// Solidity: function state_powerThreshold() view returns(uint256)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) StatePowerThreshold() (*big.Int, error) {
	return _IFxBridgeLogic.Contract.StatePowerThreshold(&_IFxBridgeLogic.CallOpts)
}

// TokenStatus is a free data retrieval call binding the contract method 0x0acac942.
//
// Solidity: function tokenStatus(address _tokenAddr) view returns((bool,bool,bool,uint8))
func (_IFxBridgeLogic *IFxBridgeLogicCaller) TokenStatus(opts *bind.CallOpts, _tokenAddr common.Address) (IFxBridgeLogicTokenStatus, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "tokenStatus", _tokenAddr)

	if err != nil {
		return *new(IFxBridgeLogicTokenStatus), err
	}

	out0 := *abi.ConvertType(out[0], new(IFxBridgeLogicTokenStatus)).(*IFxBridgeLogicTokenStatus)

	return out0, err

}

// TokenStatus is a free data retrieval call binding the contract method 0x0acac942.
//
// Solidity: function tokenStatus(address _tokenAddr) view returns((bool,bool,bool,uint8))
func (_IFxBridgeLogic *IFxBridgeLogicSession) TokenStatus(_tokenAddr common.Address) (IFxBridgeLogicTokenStatus, error) {
	return _IFxBridgeLogic.Contract.TokenStatus(&_IFxBridgeLogic.CallOpts, _tokenAddr)
}

// TokenStatus is a free data retrieval call binding the contract method 0x0acac942.
//
// Solidity: function tokenStatus(address _tokenAddr) view returns((bool,bool,bool,uint8))
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) TokenStatus(_tokenAddr common.Address) (IFxBridgeLogicTokenStatus, error) {
	return _IFxBridgeLogic.Contract.TokenStatus(&_IFxBridgeLogic.CallOpts, _tokenAddr)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_IFxBridgeLogic *IFxBridgeLogicCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IFxBridgeLogic.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_IFxBridgeLogic *IFxBridgeLogicSession) Version() (string, error) {
	return _IFxBridgeLogic.Contract.Version(&_IFxBridgeLogic.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_IFxBridgeLogic *IFxBridgeLogicCallerSession) Version() (string, error) {
	return _IFxBridgeLogic.Contract.Version(&_IFxBridgeLogic.CallOpts)
}

// ActiveBridgeToken is a paid mutator transaction binding the contract method 0xdde65aea.
//
// Solidity: function activeBridgeToken(address _tokenAddr) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) ActiveBridgeToken(opts *bind.TransactOpts, _tokenAddr common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "activeBridgeToken", _tokenAddr)
}

// ActiveBridgeToken is a paid mutator transaction binding the contract method 0xdde65aea.
//
// Solidity: function activeBridgeToken(address _tokenAddr) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicSession) ActiveBridgeToken(_tokenAddr common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.ActiveBridgeToken(&_IFxBridgeLogic.TransactOpts, _tokenAddr)
}

// ActiveBridgeToken is a paid mutator transaction binding the contract method 0xdde65aea.
//
// Solidity: function activeBridgeToken(address _tokenAddr) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) ActiveBridgeToken(_tokenAddr common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.ActiveBridgeToken(&_IFxBridgeLogic.TransactOpts, _tokenAddr)
}

// AddBridgeToken is a paid mutator transaction binding the contract method 0x4557c080.
//
// Solidity: function addBridgeToken(address _tokenAddr, bytes32 _channelIBC, bool _isOriginated, uint8 _tokenType) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) AddBridgeToken(opts *bind.TransactOpts, _tokenAddr common.Address, _channelIBC [32]byte, _isOriginated bool, _tokenType uint8) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "addBridgeToken", _tokenAddr, _channelIBC, _isOriginated, _tokenType)
}

// AddBridgeToken is a paid mutator transaction binding the contract method 0x4557c080.
//
// Solidity: function addBridgeToken(address _tokenAddr, bytes32 _channelIBC, bool _isOriginated, uint8 _tokenType) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicSession) AddBridgeToken(_tokenAddr common.Address, _channelIBC [32]byte, _isOriginated bool, _tokenType uint8) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.AddBridgeToken(&_IFxBridgeLogic.TransactOpts, _tokenAddr, _channelIBC, _isOriginated, _tokenType)
}

// AddBridgeToken is a paid mutator transaction binding the contract method 0x4557c080.
//
// Solidity: function addBridgeToken(address _tokenAddr, bytes32 _channelIBC, bool _isOriginated, uint8 _tokenType) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) AddBridgeToken(_tokenAddr common.Address, _channelIBC [32]byte, _isOriginated bool, _tokenType uint8) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.AddBridgeToken(&_IFxBridgeLogic.TransactOpts, _tokenAddr, _channelIBC, _isOriginated, _tokenType)
}

// BridgeCall is a paid mutator transaction binding the contract method 0x251477c7.
//
// Solidity: function bridgeCall(string _dstChainId, uint256 _gasLimit, address _receiver, address _to, address[] _tokens, uint256[] _amounts, bytes _message, uint256 _value) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) BridgeCall(opts *bind.TransactOpts, _dstChainId string, _gasLimit *big.Int, _receiver common.Address, _to common.Address, _tokens []common.Address, _amounts []*big.Int, _message []byte, _value *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "bridgeCall", _dstChainId, _gasLimit, _receiver, _to, _tokens, _amounts, _message, _value)
}

// BridgeCall is a paid mutator transaction binding the contract method 0x251477c7.
//
// Solidity: function bridgeCall(string _dstChainId, uint256 _gasLimit, address _receiver, address _to, address[] _tokens, uint256[] _amounts, bytes _message, uint256 _value) returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) BridgeCall(_dstChainId string, _gasLimit *big.Int, _receiver common.Address, _to common.Address, _tokens []common.Address, _amounts []*big.Int, _message []byte, _value *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.BridgeCall(&_IFxBridgeLogic.TransactOpts, _dstChainId, _gasLimit, _receiver, _to, _tokens, _amounts, _message, _value)
}

// BridgeCall is a paid mutator transaction binding the contract method 0x251477c7.
//
// Solidity: function bridgeCall(string _dstChainId, uint256 _gasLimit, address _receiver, address _to, address[] _tokens, uint256[] _amounts, bytes _message, uint256 _value) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) BridgeCall(_dstChainId string, _gasLimit *big.Int, _receiver common.Address, _to common.Address, _tokens []common.Address, _amounts []*big.Int, _message []byte, _value *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.BridgeCall(&_IFxBridgeLogic.TransactOpts, _dstChainId, _gasLimit, _receiver, _to, _tokens, _amounts, _message, _value)
}

// BridgeCallCheckpoint is a paid mutator transaction binding the contract method 0x6a6c5d61.
//
// Solidity: function bridgeCallCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, address _sender, address _to, address _receiver, uint256 _value, uint256 _nonce, uint256 _gasLimit, uint256 _timeout, string _dstChain, bytes _message, bytes _asset) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) BridgeCallCheckpoint(opts *bind.TransactOpts, _fxbridgeId [32]byte, _methodName [32]byte, _sender common.Address, _to common.Address, _receiver common.Address, _value *big.Int, _nonce *big.Int, _gasLimit *big.Int, _timeout *big.Int, _dstChain string, _message []byte, _asset []byte) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "bridgeCallCheckpoint", _fxbridgeId, _methodName, _sender, _to, _receiver, _value, _nonce, _gasLimit, _timeout, _dstChain, _message, _asset)
}

// BridgeCallCheckpoint is a paid mutator transaction binding the contract method 0x6a6c5d61.
//
// Solidity: function bridgeCallCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, address _sender, address _to, address _receiver, uint256 _value, uint256 _nonce, uint256 _gasLimit, uint256 _timeout, string _dstChain, bytes _message, bytes _asset) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicSession) BridgeCallCheckpoint(_fxbridgeId [32]byte, _methodName [32]byte, _sender common.Address, _to common.Address, _receiver common.Address, _value *big.Int, _nonce *big.Int, _gasLimit *big.Int, _timeout *big.Int, _dstChain string, _message []byte, _asset []byte) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.BridgeCallCheckpoint(&_IFxBridgeLogic.TransactOpts, _fxbridgeId, _methodName, _sender, _to, _receiver, _value, _nonce, _gasLimit, _timeout, _dstChain, _message, _asset)
}

// BridgeCallCheckpoint is a paid mutator transaction binding the contract method 0x6a6c5d61.
//
// Solidity: function bridgeCallCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, address _sender, address _to, address _receiver, uint256 _value, uint256 _nonce, uint256 _gasLimit, uint256 _timeout, string _dstChain, bytes _message, bytes _asset) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) BridgeCallCheckpoint(_fxbridgeId [32]byte, _methodName [32]byte, _sender common.Address, _to common.Address, _receiver common.Address, _value *big.Int, _nonce *big.Int, _gasLimit *big.Int, _timeout *big.Int, _dstChain string, _message []byte, _asset []byte) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.BridgeCallCheckpoint(&_IFxBridgeLogic.TransactOpts, _fxbridgeId, _methodName, _sender, _to, _receiver, _value, _nonce, _gasLimit, _timeout, _dstChain, _message, _asset)
}

// OracleSetCheckpoint is a paid mutator transaction binding the contract method 0xa955665f.
//
// Solidity: function oracleSetCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, uint256 _oracleSetNonce, address[] _oracles, uint256[] _powers) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) OracleSetCheckpoint(opts *bind.TransactOpts, _fxbridgeId [32]byte, _methodName [32]byte, _oracleSetNonce *big.Int, _oracles []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "oracleSetCheckpoint", _fxbridgeId, _methodName, _oracleSetNonce, _oracles, _powers)
}

// OracleSetCheckpoint is a paid mutator transaction binding the contract method 0xa955665f.
//
// Solidity: function oracleSetCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, uint256 _oracleSetNonce, address[] _oracles, uint256[] _powers) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicSession) OracleSetCheckpoint(_fxbridgeId [32]byte, _methodName [32]byte, _oracleSetNonce *big.Int, _oracles []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.OracleSetCheckpoint(&_IFxBridgeLogic.TransactOpts, _fxbridgeId, _methodName, _oracleSetNonce, _oracles, _powers)
}

// OracleSetCheckpoint is a paid mutator transaction binding the contract method 0xa955665f.
//
// Solidity: function oracleSetCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, uint256 _oracleSetNonce, address[] _oracles, uint256[] _powers) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) OracleSetCheckpoint(_fxbridgeId [32]byte, _methodName [32]byte, _oracleSetNonce *big.Int, _oracles []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.OracleSetCheckpoint(&_IFxBridgeLogic.TransactOpts, _fxbridgeId, _methodName, _oracleSetNonce, _oracles, _powers)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) Pause() (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.Pause(&_IFxBridgeLogic.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) Pause() (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.Pause(&_IFxBridgeLogic.TransactOpts)
}

// PauseBridgeToken is a paid mutator transaction binding the contract method 0xa36a4ab0.
//
// Solidity: function pauseBridgeToken(address _tokenAddr) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) PauseBridgeToken(opts *bind.TransactOpts, _tokenAddr common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "pauseBridgeToken", _tokenAddr)
}

// PauseBridgeToken is a paid mutator transaction binding the contract method 0xa36a4ab0.
//
// Solidity: function pauseBridgeToken(address _tokenAddr) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicSession) PauseBridgeToken(_tokenAddr common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.PauseBridgeToken(&_IFxBridgeLogic.TransactOpts, _tokenAddr)
}

// PauseBridgeToken is a paid mutator transaction binding the contract method 0xa36a4ab0.
//
// Solidity: function pauseBridgeToken(address _tokenAddr) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) PauseBridgeToken(_tokenAddr common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.PauseBridgeToken(&_IFxBridgeLogic.TransactOpts, _tokenAddr)
}

// RefundBridgeToken is a paid mutator transaction binding the contract method 0x5e438fcf.
//
// Solidity: function refundBridgeToken(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[2] _nonceArray, address _receiver, address[] _tokens, uint256[] _amounts, uint256 _timeout) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) RefundBridgeToken(opts *bind.TransactOpts, _currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _nonceArray [2]*big.Int, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _timeout *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "refundBridgeToken", _currentOracles, _currentPowers, _v, _r, _s, _nonceArray, _receiver, _tokens, _amounts, _timeout)
}

// RefundBridgeToken is a paid mutator transaction binding the contract method 0x5e438fcf.
//
// Solidity: function refundBridgeToken(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[2] _nonceArray, address _receiver, address[] _tokens, uint256[] _amounts, uint256 _timeout) returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) RefundBridgeToken(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _nonceArray [2]*big.Int, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _timeout *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.RefundBridgeToken(&_IFxBridgeLogic.TransactOpts, _currentOracles, _currentPowers, _v, _r, _s, _nonceArray, _receiver, _tokens, _amounts, _timeout)
}

// RefundBridgeToken is a paid mutator transaction binding the contract method 0x5e438fcf.
//
// Solidity: function refundBridgeToken(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[2] _nonceArray, address _receiver, address[] _tokens, uint256[] _amounts, uint256 _timeout) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) RefundBridgeToken(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _nonceArray [2]*big.Int, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _timeout *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.RefundBridgeToken(&_IFxBridgeLogic.TransactOpts, _currentOracles, _currentPowers, _v, _r, _s, _nonceArray, _receiver, _tokens, _amounts, _timeout)
}

// SendToFx is a paid mutator transaction binding the contract method 0x6189d107.
//
// Solidity: function sendToFx(address _tokenContract, bytes32 _destination, bytes32 _targetIBC, uint256 _amount) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) SendToFx(opts *bind.TransactOpts, _tokenContract common.Address, _destination [32]byte, _targetIBC [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "sendToFx", _tokenContract, _destination, _targetIBC, _amount)
}

// SendToFx is a paid mutator transaction binding the contract method 0x6189d107.
//
// Solidity: function sendToFx(address _tokenContract, bytes32 _destination, bytes32 _targetIBC, uint256 _amount) returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) SendToFx(_tokenContract common.Address, _destination [32]byte, _targetIBC [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SendToFx(&_IFxBridgeLogic.TransactOpts, _tokenContract, _destination, _targetIBC, _amount)
}

// SendToFx is a paid mutator transaction binding the contract method 0x6189d107.
//
// Solidity: function sendToFx(address _tokenContract, bytes32 _destination, bytes32 _targetIBC, uint256 _amount) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) SendToFx(_tokenContract common.Address, _destination [32]byte, _targetIBC [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SendToFx(&_IFxBridgeLogic.TransactOpts, _tokenContract, _destination, _targetIBC, _amount)
}

// SubmitBatch is a paid mutator transaction binding the contract method 0x332caa1f.
//
// Solidity: function submitBatch(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256[2] _nonceArray, address _tokenContract, uint256 _batchTimeout, address _feeReceive) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) SubmitBatch(opts *bind.TransactOpts, _currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _amounts []*big.Int, _destinations []common.Address, _fees []*big.Int, _nonceArray [2]*big.Int, _tokenContract common.Address, _batchTimeout *big.Int, _feeReceive common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "submitBatch", _currentOracles, _currentPowers, _v, _r, _s, _amounts, _destinations, _fees, _nonceArray, _tokenContract, _batchTimeout, _feeReceive)
}

// SubmitBatch is a paid mutator transaction binding the contract method 0x332caa1f.
//
// Solidity: function submitBatch(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256[2] _nonceArray, address _tokenContract, uint256 _batchTimeout, address _feeReceive) returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) SubmitBatch(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _amounts []*big.Int, _destinations []common.Address, _fees []*big.Int, _nonceArray [2]*big.Int, _tokenContract common.Address, _batchTimeout *big.Int, _feeReceive common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SubmitBatch(&_IFxBridgeLogic.TransactOpts, _currentOracles, _currentPowers, _v, _r, _s, _amounts, _destinations, _fees, _nonceArray, _tokenContract, _batchTimeout, _feeReceive)
}

// SubmitBatch is a paid mutator transaction binding the contract method 0x332caa1f.
//
// Solidity: function submitBatch(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256[2] _nonceArray, address _tokenContract, uint256 _batchTimeout, address _feeReceive) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) SubmitBatch(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _amounts []*big.Int, _destinations []common.Address, _fees []*big.Int, _nonceArray [2]*big.Int, _tokenContract common.Address, _batchTimeout *big.Int, _feeReceive common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SubmitBatch(&_IFxBridgeLogic.TransactOpts, _currentOracles, _currentPowers, _v, _r, _s, _amounts, _destinations, _fees, _nonceArray, _tokenContract, _batchTimeout, _feeReceive)
}

// SubmitBatchCheckpoint is a paid mutator transaction binding the contract method 0x3d1e51f9.
//
// Solidity: function submitBatchCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256 _batchNonce, address _tokenContract, uint256 _batchTimeout, address _feeReceive) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) SubmitBatchCheckpoint(opts *bind.TransactOpts, _fxbridgeId [32]byte, _methodName [32]byte, _amounts []*big.Int, _destinations []common.Address, _fees []*big.Int, _batchNonce *big.Int, _tokenContract common.Address, _batchTimeout *big.Int, _feeReceive common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "submitBatchCheckpoint", _fxbridgeId, _methodName, _amounts, _destinations, _fees, _batchNonce, _tokenContract, _batchTimeout, _feeReceive)
}

// SubmitBatchCheckpoint is a paid mutator transaction binding the contract method 0x3d1e51f9.
//
// Solidity: function submitBatchCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256 _batchNonce, address _tokenContract, uint256 _batchTimeout, address _feeReceive) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicSession) SubmitBatchCheckpoint(_fxbridgeId [32]byte, _methodName [32]byte, _amounts []*big.Int, _destinations []common.Address, _fees []*big.Int, _batchNonce *big.Int, _tokenContract common.Address, _batchTimeout *big.Int, _feeReceive common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SubmitBatchCheckpoint(&_IFxBridgeLogic.TransactOpts, _fxbridgeId, _methodName, _amounts, _destinations, _fees, _batchNonce, _tokenContract, _batchTimeout, _feeReceive)
}

// SubmitBatchCheckpoint is a paid mutator transaction binding the contract method 0x3d1e51f9.
//
// Solidity: function submitBatchCheckpoint(bytes32 _fxbridgeId, bytes32 _methodName, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256 _batchNonce, address _tokenContract, uint256 _batchTimeout, address _feeReceive) returns(bytes32)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) SubmitBatchCheckpoint(_fxbridgeId [32]byte, _methodName [32]byte, _amounts []*big.Int, _destinations []common.Address, _fees []*big.Int, _batchNonce *big.Int, _tokenContract common.Address, _batchTimeout *big.Int, _feeReceive common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SubmitBatchCheckpoint(&_IFxBridgeLogic.TransactOpts, _fxbridgeId, _methodName, _amounts, _destinations, _fees, _batchNonce, _tokenContract, _batchTimeout, _feeReceive)
}

// SubmitBridgeCall is a paid mutator transaction binding the contract method 0xab9838a7.
//
// Solidity: function submitBridgeCall(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[2] _nonceArray, (address,address,address,address[],uint256[],bytes,uint256,uint256,uint256) _input) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) SubmitBridgeCall(opts *bind.TransactOpts, _currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _nonceArray [2]*big.Int, _input IFxBridgeLogicBridgeCallData) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "submitBridgeCall", _currentOracles, _currentPowers, _v, _r, _s, _nonceArray, _input)
}

// SubmitBridgeCall is a paid mutator transaction binding the contract method 0xab9838a7.
//
// Solidity: function submitBridgeCall(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[2] _nonceArray, (address,address,address,address[],uint256[],bytes,uint256,uint256,uint256) _input) returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) SubmitBridgeCall(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _nonceArray [2]*big.Int, _input IFxBridgeLogicBridgeCallData) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SubmitBridgeCall(&_IFxBridgeLogic.TransactOpts, _currentOracles, _currentPowers, _v, _r, _s, _nonceArray, _input)
}

// SubmitBridgeCall is a paid mutator transaction binding the contract method 0xab9838a7.
//
// Solidity: function submitBridgeCall(address[] _currentOracles, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[2] _nonceArray, (address,address,address,address[],uint256[],bytes,uint256,uint256,uint256) _input) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) SubmitBridgeCall(_currentOracles []common.Address, _currentPowers []*big.Int, _v []uint8, _r [][32]byte, _s [][32]byte, _nonceArray [2]*big.Int, _input IFxBridgeLogicBridgeCallData) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.SubmitBridgeCall(&_IFxBridgeLogic.TransactOpts, _currentOracles, _currentPowers, _v, _r, _s, _nonceArray, _input)
}

// TransferOwner is a paid mutator transaction binding the contract method 0x31678cf6.
//
// Solidity: function transferOwner(address _token, address _newOwner) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) TransferOwner(opts *bind.TransactOpts, _token common.Address, _newOwner common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "transferOwner", _token, _newOwner)
}

// TransferOwner is a paid mutator transaction binding the contract method 0x31678cf6.
//
// Solidity: function transferOwner(address _token, address _newOwner) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicSession) TransferOwner(_token common.Address, _newOwner common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.TransferOwner(&_IFxBridgeLogic.TransactOpts, _token, _newOwner)
}

// TransferOwner is a paid mutator transaction binding the contract method 0x31678cf6.
//
// Solidity: function transferOwner(address _token, address _newOwner) returns(bool)
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) TransferOwner(_token common.Address, _newOwner common.Address) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.TransferOwner(&_IFxBridgeLogic.TransactOpts, _token, _newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) Unpause() (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.Unpause(&_IFxBridgeLogic.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) Unpause() (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.Unpause(&_IFxBridgeLogic.TransactOpts)
}

// UpdateOracleSet is a paid mutator transaction binding the contract method 0x3a08e299.
//
// Solidity: function updateOracleSet(address[] _newOracles, uint256[] _newPowers, uint256 _newOracleSetNonce, address[] _currentOracles, uint256[] _currentPowers, uint256 _currentOracleSetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactor) UpdateOracleSet(opts *bind.TransactOpts, _newOracles []common.Address, _newPowers []*big.Int, _newOracleSetNonce *big.Int, _currentOracles []common.Address, _currentPowers []*big.Int, _currentOracleSetNonce *big.Int, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _IFxBridgeLogic.contract.Transact(opts, "updateOracleSet", _newOracles, _newPowers, _newOracleSetNonce, _currentOracles, _currentPowers, _currentOracleSetNonce, _v, _r, _s)
}

// UpdateOracleSet is a paid mutator transaction binding the contract method 0x3a08e299.
//
// Solidity: function updateOracleSet(address[] _newOracles, uint256[] _newPowers, uint256 _newOracleSetNonce, address[] _currentOracles, uint256[] _currentPowers, uint256 _currentOracleSetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s) returns()
func (_IFxBridgeLogic *IFxBridgeLogicSession) UpdateOracleSet(_newOracles []common.Address, _newPowers []*big.Int, _newOracleSetNonce *big.Int, _currentOracles []common.Address, _currentPowers []*big.Int, _currentOracleSetNonce *big.Int, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.UpdateOracleSet(&_IFxBridgeLogic.TransactOpts, _newOracles, _newPowers, _newOracleSetNonce, _currentOracles, _currentPowers, _currentOracleSetNonce, _v, _r, _s)
}

// UpdateOracleSet is a paid mutator transaction binding the contract method 0x3a08e299.
//
// Solidity: function updateOracleSet(address[] _newOracles, uint256[] _newPowers, uint256 _newOracleSetNonce, address[] _currentOracles, uint256[] _currentPowers, uint256 _currentOracleSetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s) returns()
func (_IFxBridgeLogic *IFxBridgeLogicTransactorSession) UpdateOracleSet(_newOracles []common.Address, _newPowers []*big.Int, _newOracleSetNonce *big.Int, _currentOracles []common.Address, _currentPowers []*big.Int, _currentOracleSetNonce *big.Int, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _IFxBridgeLogic.Contract.UpdateOracleSet(&_IFxBridgeLogic.TransactOpts, _newOracles, _newPowers, _newOracleSetNonce, _currentOracles, _currentPowers, _currentOracleSetNonce, _v, _r, _s)
}

// IFxBridgeLogicAddBridgeTokenEventIterator is returned from FilterAddBridgeTokenEvent and is used to iterate over the raw logs and unpacked data for AddBridgeTokenEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicAddBridgeTokenEventIterator struct {
	Event *IFxBridgeLogicAddBridgeTokenEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicAddBridgeTokenEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicAddBridgeTokenEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicAddBridgeTokenEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicAddBridgeTokenEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicAddBridgeTokenEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicAddBridgeTokenEvent represents a AddBridgeTokenEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicAddBridgeTokenEvent struct {
	TokenContract common.Address
	Name          string
	Symbol        string
	Decimals      uint8
	EventNonce    *big.Int
	ChannelIBC    [32]byte
	TokenType     uint8
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAddBridgeTokenEvent is a free log retrieval operation binding the contract event 0x2a9a8067c774124af41628bf45e0b4f0a3670160c40309af3fead56eaf270def.
//
// Solidity: event AddBridgeTokenEvent(address indexed _tokenContract, string _name, string _symbol, uint8 _decimals, uint256 _eventNonce, bytes32 _channelIBC, uint8 _tokenType)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterAddBridgeTokenEvent(opts *bind.FilterOpts, _tokenContract []common.Address) (*IFxBridgeLogicAddBridgeTokenEventIterator, error) {

	var _tokenContractRule []interface{}
	for _, _tokenContractItem := range _tokenContract {
		_tokenContractRule = append(_tokenContractRule, _tokenContractItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "AddBridgeTokenEvent", _tokenContractRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicAddBridgeTokenEventIterator{contract: _IFxBridgeLogic.contract, event: "AddBridgeTokenEvent", logs: logs, sub: sub}, nil
}

// WatchAddBridgeTokenEvent is a free log subscription operation binding the contract event 0x2a9a8067c774124af41628bf45e0b4f0a3670160c40309af3fead56eaf270def.
//
// Solidity: event AddBridgeTokenEvent(address indexed _tokenContract, string _name, string _symbol, uint8 _decimals, uint256 _eventNonce, bytes32 _channelIBC, uint8 _tokenType)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchAddBridgeTokenEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicAddBridgeTokenEvent, _tokenContract []common.Address) (event.Subscription, error) {

	var _tokenContractRule []interface{}
	for _, _tokenContractItem := range _tokenContract {
		_tokenContractRule = append(_tokenContractRule, _tokenContractItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "AddBridgeTokenEvent", _tokenContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicAddBridgeTokenEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "AddBridgeTokenEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddBridgeTokenEvent is a log parse operation binding the contract event 0x2a9a8067c774124af41628bf45e0b4f0a3670160c40309af3fead56eaf270def.
//
// Solidity: event AddBridgeTokenEvent(address indexed _tokenContract, string _name, string _symbol, uint8 _decimals, uint256 _eventNonce, bytes32 _channelIBC, uint8 _tokenType)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseAddBridgeTokenEvent(log types.Log) (*IFxBridgeLogicAddBridgeTokenEvent, error) {
	event := new(IFxBridgeLogicAddBridgeTokenEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "AddBridgeTokenEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicBridgeCallEventIterator is returned from FilterBridgeCallEvent and is used to iterate over the raw logs and unpacked data for BridgeCallEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicBridgeCallEventIterator struct {
	Event *IFxBridgeLogicBridgeCallEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicBridgeCallEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicBridgeCallEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicBridgeCallEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicBridgeCallEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicBridgeCallEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicBridgeCallEvent represents a BridgeCallEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicBridgeCallEvent struct {
	Sender     common.Address
	Receiver   common.Address
	To         common.Address
	Tokens     []common.Address
	Amounts    []*big.Int
	EventNonce *big.Int
	DstChainId string
	GasLimit   *big.Int
	Message    []byte
	Value      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBridgeCallEvent is a free log retrieval operation binding the contract event 0xed79a7bce50ea6c06eb0c0ce8d1ebcc6c792e869d54b505722b8795384ef9359.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address[] _tokens, uint256[] _amounts, uint256 _eventNonce, string _dstChainId, uint256 _gasLimit, bytes _message, uint256 _value)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterBridgeCallEvent(opts *bind.FilterOpts, _sender []common.Address, _receiver []common.Address, _to []common.Address) (*IFxBridgeLogicBridgeCallEventIterator, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "BridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicBridgeCallEventIterator{contract: _IFxBridgeLogic.contract, event: "BridgeCallEvent", logs: logs, sub: sub}, nil
}

// WatchBridgeCallEvent is a free log subscription operation binding the contract event 0xed79a7bce50ea6c06eb0c0ce8d1ebcc6c792e869d54b505722b8795384ef9359.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address[] _tokens, uint256[] _amounts, uint256 _eventNonce, string _dstChainId, uint256 _gasLimit, bytes _message, uint256 _value)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchBridgeCallEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicBridgeCallEvent, _sender []common.Address, _receiver []common.Address, _to []common.Address) (event.Subscription, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "BridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicBridgeCallEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "BridgeCallEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBridgeCallEvent is a log parse operation binding the contract event 0xed79a7bce50ea6c06eb0c0ce8d1ebcc6c792e869d54b505722b8795384ef9359.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address[] _tokens, uint256[] _amounts, uint256 _eventNonce, string _dstChainId, uint256 _gasLimit, bytes _message, uint256 _value)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseBridgeCallEvent(log types.Log) (*IFxBridgeLogicBridgeCallEvent, error) {
	event := new(IFxBridgeLogicBridgeCallEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "BridgeCallEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicOracleSetUpdatedEventIterator is returned from FilterOracleSetUpdatedEvent and is used to iterate over the raw logs and unpacked data for OracleSetUpdatedEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicOracleSetUpdatedEventIterator struct {
	Event *IFxBridgeLogicOracleSetUpdatedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicOracleSetUpdatedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicOracleSetUpdatedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicOracleSetUpdatedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicOracleSetUpdatedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicOracleSetUpdatedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicOracleSetUpdatedEvent represents a OracleSetUpdatedEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicOracleSetUpdatedEvent struct {
	NewOracleSetNonce *big.Int
	EventNonce        *big.Int
	Oracles           []common.Address
	Powers            []*big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterOracleSetUpdatedEvent is a free log retrieval operation binding the contract event 0x36c6022aad02313069de85ca9645431c7dd5e8e7a21685586461c4b25e2374b3.
//
// Solidity: event OracleSetUpdatedEvent(uint256 indexed _newOracleSetNonce, uint256 _eventNonce, address[] _oracles, uint256[] _powers)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterOracleSetUpdatedEvent(opts *bind.FilterOpts, _newOracleSetNonce []*big.Int) (*IFxBridgeLogicOracleSetUpdatedEventIterator, error) {

	var _newOracleSetNonceRule []interface{}
	for _, _newOracleSetNonceItem := range _newOracleSetNonce {
		_newOracleSetNonceRule = append(_newOracleSetNonceRule, _newOracleSetNonceItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "OracleSetUpdatedEvent", _newOracleSetNonceRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicOracleSetUpdatedEventIterator{contract: _IFxBridgeLogic.contract, event: "OracleSetUpdatedEvent", logs: logs, sub: sub}, nil
}

// WatchOracleSetUpdatedEvent is a free log subscription operation binding the contract event 0x36c6022aad02313069de85ca9645431c7dd5e8e7a21685586461c4b25e2374b3.
//
// Solidity: event OracleSetUpdatedEvent(uint256 indexed _newOracleSetNonce, uint256 _eventNonce, address[] _oracles, uint256[] _powers)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchOracleSetUpdatedEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicOracleSetUpdatedEvent, _newOracleSetNonce []*big.Int) (event.Subscription, error) {

	var _newOracleSetNonceRule []interface{}
	for _, _newOracleSetNonceItem := range _newOracleSetNonce {
		_newOracleSetNonceRule = append(_newOracleSetNonceRule, _newOracleSetNonceItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "OracleSetUpdatedEvent", _newOracleSetNonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicOracleSetUpdatedEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "OracleSetUpdatedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOracleSetUpdatedEvent is a log parse operation binding the contract event 0x36c6022aad02313069de85ca9645431c7dd5e8e7a21685586461c4b25e2374b3.
//
// Solidity: event OracleSetUpdatedEvent(uint256 indexed _newOracleSetNonce, uint256 _eventNonce, address[] _oracles, uint256[] _powers)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseOracleSetUpdatedEvent(log types.Log) (*IFxBridgeLogicOracleSetUpdatedEvent, error) {
	event := new(IFxBridgeLogicOracleSetUpdatedEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "OracleSetUpdatedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicRefundTokenExecutedEventIterator is returned from FilterRefundTokenExecutedEvent and is used to iterate over the raw logs and unpacked data for RefundTokenExecutedEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicRefundTokenExecutedEventIterator struct {
	Event *IFxBridgeLogicRefundTokenExecutedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicRefundTokenExecutedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicRefundTokenExecutedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicRefundTokenExecutedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicRefundTokenExecutedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicRefundTokenExecutedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicRefundTokenExecutedEvent represents a RefundTokenExecutedEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicRefundTokenExecutedEvent struct {
	Receiver    common.Address
	RefundNonce *big.Int
	EventNonce  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRefundTokenExecutedEvent is a free log retrieval operation binding the contract event 0x6dcbf583591c8ea7ca09e71708417169bd9d029f7ec9c7c23aeb204a697ee815.
//
// Solidity: event RefundTokenExecutedEvent(address indexed _receiver, uint256 indexed _refundNonce, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterRefundTokenExecutedEvent(opts *bind.FilterOpts, _receiver []common.Address, _refundNonce []*big.Int) (*IFxBridgeLogicRefundTokenExecutedEventIterator, error) {

	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _refundNonceRule []interface{}
	for _, _refundNonceItem := range _refundNonce {
		_refundNonceRule = append(_refundNonceRule, _refundNonceItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "RefundTokenExecutedEvent", _receiverRule, _refundNonceRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicRefundTokenExecutedEventIterator{contract: _IFxBridgeLogic.contract, event: "RefundTokenExecutedEvent", logs: logs, sub: sub}, nil
}

// WatchRefundTokenExecutedEvent is a free log subscription operation binding the contract event 0x6dcbf583591c8ea7ca09e71708417169bd9d029f7ec9c7c23aeb204a697ee815.
//
// Solidity: event RefundTokenExecutedEvent(address indexed _receiver, uint256 indexed _refundNonce, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchRefundTokenExecutedEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicRefundTokenExecutedEvent, _receiver []common.Address, _refundNonce []*big.Int) (event.Subscription, error) {

	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _refundNonceRule []interface{}
	for _, _refundNonceItem := range _refundNonce {
		_refundNonceRule = append(_refundNonceRule, _refundNonceItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "RefundTokenExecutedEvent", _receiverRule, _refundNonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicRefundTokenExecutedEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "RefundTokenExecutedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefundTokenExecutedEvent is a log parse operation binding the contract event 0x6dcbf583591c8ea7ca09e71708417169bd9d029f7ec9c7c23aeb204a697ee815.
//
// Solidity: event RefundTokenExecutedEvent(address indexed _receiver, uint256 indexed _refundNonce, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseRefundTokenExecutedEvent(log types.Log) (*IFxBridgeLogicRefundTokenExecutedEvent, error) {
	event := new(IFxBridgeLogicRefundTokenExecutedEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "RefundTokenExecutedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicSendToFxEventIterator is returned from FilterSendToFxEvent and is used to iterate over the raw logs and unpacked data for SendToFxEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicSendToFxEventIterator struct {
	Event *IFxBridgeLogicSendToFxEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicSendToFxEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicSendToFxEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicSendToFxEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicSendToFxEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicSendToFxEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicSendToFxEvent represents a SendToFxEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicSendToFxEvent struct {
	TokenContract common.Address
	Sender        common.Address
	Destination   [32]byte
	TargetIBC     [32]byte
	Amount        *big.Int
	EventNonce    *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSendToFxEvent is a free log retrieval operation binding the contract event 0x034c5b22dd525a50d0a6b15549df0a6ac83b833a6c3da57ea16890832c72507c.
//
// Solidity: event SendToFxEvent(address indexed _tokenContract, address indexed _sender, bytes32 indexed _destination, bytes32 _targetIBC, uint256 _amount, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterSendToFxEvent(opts *bind.FilterOpts, _tokenContract []common.Address, _sender []common.Address, _destination [][32]byte) (*IFxBridgeLogicSendToFxEventIterator, error) {

	var _tokenContractRule []interface{}
	for _, _tokenContractItem := range _tokenContract {
		_tokenContractRule = append(_tokenContractRule, _tokenContractItem)
	}
	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _destinationRule []interface{}
	for _, _destinationItem := range _destination {
		_destinationRule = append(_destinationRule, _destinationItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "SendToFxEvent", _tokenContractRule, _senderRule, _destinationRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicSendToFxEventIterator{contract: _IFxBridgeLogic.contract, event: "SendToFxEvent", logs: logs, sub: sub}, nil
}

// WatchSendToFxEvent is a free log subscription operation binding the contract event 0x034c5b22dd525a50d0a6b15549df0a6ac83b833a6c3da57ea16890832c72507c.
//
// Solidity: event SendToFxEvent(address indexed _tokenContract, address indexed _sender, bytes32 indexed _destination, bytes32 _targetIBC, uint256 _amount, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchSendToFxEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicSendToFxEvent, _tokenContract []common.Address, _sender []common.Address, _destination [][32]byte) (event.Subscription, error) {

	var _tokenContractRule []interface{}
	for _, _tokenContractItem := range _tokenContract {
		_tokenContractRule = append(_tokenContractRule, _tokenContractItem)
	}
	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _destinationRule []interface{}
	for _, _destinationItem := range _destination {
		_destinationRule = append(_destinationRule, _destinationItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "SendToFxEvent", _tokenContractRule, _senderRule, _destinationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicSendToFxEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "SendToFxEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSendToFxEvent is a log parse operation binding the contract event 0x034c5b22dd525a50d0a6b15549df0a6ac83b833a6c3da57ea16890832c72507c.
//
// Solidity: event SendToFxEvent(address indexed _tokenContract, address indexed _sender, bytes32 indexed _destination, bytes32 _targetIBC, uint256 _amount, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseSendToFxEvent(log types.Log) (*IFxBridgeLogicSendToFxEvent, error) {
	event := new(IFxBridgeLogicSendToFxEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "SendToFxEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicSubmitBridgeCallEventIterator is returned from FilterSubmitBridgeCallEvent and is used to iterate over the raw logs and unpacked data for SubmitBridgeCallEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicSubmitBridgeCallEventIterator struct {
	Event *IFxBridgeLogicSubmitBridgeCallEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicSubmitBridgeCallEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicSubmitBridgeCallEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicSubmitBridgeCallEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicSubmitBridgeCallEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicSubmitBridgeCallEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicSubmitBridgeCallEvent represents a SubmitBridgeCallEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicSubmitBridgeCallEvent struct {
	Sender     common.Address
	Receiver   common.Address
	To         common.Address
	Nonce      *big.Int
	EventNonce *big.Int
	Result     bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSubmitBridgeCallEvent is a free log retrieval operation binding the contract event 0x8b8629b0ea96056159172ee04000c0009e2ae7982b1e4cd9321ad4f74306316a.
//
// Solidity: event SubmitBridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, uint256 _nonce, uint256 _eventNonce, bool _result)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterSubmitBridgeCallEvent(opts *bind.FilterOpts, _sender []common.Address, _receiver []common.Address, _to []common.Address) (*IFxBridgeLogicSubmitBridgeCallEventIterator, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "SubmitBridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicSubmitBridgeCallEventIterator{contract: _IFxBridgeLogic.contract, event: "SubmitBridgeCallEvent", logs: logs, sub: sub}, nil
}

// WatchSubmitBridgeCallEvent is a free log subscription operation binding the contract event 0x8b8629b0ea96056159172ee04000c0009e2ae7982b1e4cd9321ad4f74306316a.
//
// Solidity: event SubmitBridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, uint256 _nonce, uint256 _eventNonce, bool _result)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchSubmitBridgeCallEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicSubmitBridgeCallEvent, _sender []common.Address, _receiver []common.Address, _to []common.Address) (event.Subscription, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "SubmitBridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicSubmitBridgeCallEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "SubmitBridgeCallEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSubmitBridgeCallEvent is a log parse operation binding the contract event 0x8b8629b0ea96056159172ee04000c0009e2ae7982b1e4cd9321ad4f74306316a.
//
// Solidity: event SubmitBridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, uint256 _nonce, uint256 _eventNonce, bool _result)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseSubmitBridgeCallEvent(log types.Log) (*IFxBridgeLogicSubmitBridgeCallEvent, error) {
	event := new(IFxBridgeLogicSubmitBridgeCallEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "SubmitBridgeCallEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicTransactionBatchExecutedEventIterator is returned from FilterTransactionBatchExecutedEvent and is used to iterate over the raw logs and unpacked data for TransactionBatchExecutedEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicTransactionBatchExecutedEventIterator struct {
	Event *IFxBridgeLogicTransactionBatchExecutedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicTransactionBatchExecutedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicTransactionBatchExecutedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicTransactionBatchExecutedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicTransactionBatchExecutedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicTransactionBatchExecutedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicTransactionBatchExecutedEvent represents a TransactionBatchExecutedEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicTransactionBatchExecutedEvent struct {
	BatchNonce *big.Int
	Token      common.Address
	EventNonce *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransactionBatchExecutedEvent is a free log retrieval operation binding the contract event 0x02c7e81975f8edb86e2a0c038b7b86a49c744236abf0f6177ff5afc6986ab708.
//
// Solidity: event TransactionBatchExecutedEvent(uint256 indexed _batchNonce, address indexed _token, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterTransactionBatchExecutedEvent(opts *bind.FilterOpts, _batchNonce []*big.Int, _token []common.Address) (*IFxBridgeLogicTransactionBatchExecutedEventIterator, error) {

	var _batchNonceRule []interface{}
	for _, _batchNonceItem := range _batchNonce {
		_batchNonceRule = append(_batchNonceRule, _batchNonceItem)
	}
	var _tokenRule []interface{}
	for _, _tokenItem := range _token {
		_tokenRule = append(_tokenRule, _tokenItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "TransactionBatchExecutedEvent", _batchNonceRule, _tokenRule)
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicTransactionBatchExecutedEventIterator{contract: _IFxBridgeLogic.contract, event: "TransactionBatchExecutedEvent", logs: logs, sub: sub}, nil
}

// WatchTransactionBatchExecutedEvent is a free log subscription operation binding the contract event 0x02c7e81975f8edb86e2a0c038b7b86a49c744236abf0f6177ff5afc6986ab708.
//
// Solidity: event TransactionBatchExecutedEvent(uint256 indexed _batchNonce, address indexed _token, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchTransactionBatchExecutedEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicTransactionBatchExecutedEvent, _batchNonce []*big.Int, _token []common.Address) (event.Subscription, error) {

	var _batchNonceRule []interface{}
	for _, _batchNonceItem := range _batchNonce {
		_batchNonceRule = append(_batchNonceRule, _batchNonceItem)
	}
	var _tokenRule []interface{}
	for _, _tokenItem := range _token {
		_tokenRule = append(_tokenRule, _tokenItem)
	}

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "TransactionBatchExecutedEvent", _batchNonceRule, _tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicTransactionBatchExecutedEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "TransactionBatchExecutedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransactionBatchExecutedEvent is a log parse operation binding the contract event 0x02c7e81975f8edb86e2a0c038b7b86a49c744236abf0f6177ff5afc6986ab708.
//
// Solidity: event TransactionBatchExecutedEvent(uint256 indexed _batchNonce, address indexed _token, uint256 _eventNonce)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseTransactionBatchExecutedEvent(log types.Log) (*IFxBridgeLogicTransactionBatchExecutedEvent, error) {
	event := new(IFxBridgeLogicTransactionBatchExecutedEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "TransactionBatchExecutedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFxBridgeLogicTransferOwnerEventIterator is returned from FilterTransferOwnerEvent and is used to iterate over the raw logs and unpacked data for TransferOwnerEvent events raised by the IFxBridgeLogic contract.
type IFxBridgeLogicTransferOwnerEventIterator struct {
	Event *IFxBridgeLogicTransferOwnerEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IFxBridgeLogicTransferOwnerEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFxBridgeLogicTransferOwnerEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IFxBridgeLogicTransferOwnerEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IFxBridgeLogicTransferOwnerEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFxBridgeLogicTransferOwnerEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFxBridgeLogicTransferOwnerEvent represents a TransferOwnerEvent event raised by the IFxBridgeLogic contract.
type IFxBridgeLogicTransferOwnerEvent struct {
	Token    common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferOwnerEvent is a free log retrieval operation binding the contract event 0xb0f1bf050fff9d249d22389b0f2673295260c8deca341a2755d95318f9fbc699.
//
// Solidity: event TransferOwnerEvent(address _token, address _newOwner)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) FilterTransferOwnerEvent(opts *bind.FilterOpts) (*IFxBridgeLogicTransferOwnerEventIterator, error) {

	logs, sub, err := _IFxBridgeLogic.contract.FilterLogs(opts, "TransferOwnerEvent")
	if err != nil {
		return nil, err
	}
	return &IFxBridgeLogicTransferOwnerEventIterator{contract: _IFxBridgeLogic.contract, event: "TransferOwnerEvent", logs: logs, sub: sub}, nil
}

// WatchTransferOwnerEvent is a free log subscription operation binding the contract event 0xb0f1bf050fff9d249d22389b0f2673295260c8deca341a2755d95318f9fbc699.
//
// Solidity: event TransferOwnerEvent(address _token, address _newOwner)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) WatchTransferOwnerEvent(opts *bind.WatchOpts, sink chan<- *IFxBridgeLogicTransferOwnerEvent) (event.Subscription, error) {

	logs, sub, err := _IFxBridgeLogic.contract.WatchLogs(opts, "TransferOwnerEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFxBridgeLogicTransferOwnerEvent)
				if err := _IFxBridgeLogic.contract.UnpackLog(event, "TransferOwnerEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferOwnerEvent is a log parse operation binding the contract event 0xb0f1bf050fff9d249d22389b0f2673295260c8deca341a2755d95318f9fbc699.
//
// Solidity: event TransferOwnerEvent(address _token, address _newOwner)
func (_IFxBridgeLogic *IFxBridgeLogicFilterer) ParseTransferOwnerEvent(log types.Log) (*IFxBridgeLogicTransferOwnerEvent, error) {
	event := new(IFxBridgeLogicTransferOwnerEvent)
	if err := _IFxBridgeLogic.contract.UnpackLog(event, "TransferOwnerEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
