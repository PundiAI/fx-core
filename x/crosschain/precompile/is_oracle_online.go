package precompile

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/contract"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/evm/types"
)

type IsOracleOnlineMethod struct {
	*Keeper
	IsOracleOnlineABI
}

func NewIsOracleOnlineMethod(keeper *Keeper) *IsOracleOnlineMethod {
	return &IsOracleOnlineMethod{
		Keeper:            keeper,
		IsOracleOnlineABI: NewIsOracleOnlineABI(),
	}
}

func (i *IsOracleOnlineMethod) GetMethodId() []byte {
	return i.Method.ID
}

func (i *IsOracleOnlineMethod) RequiredGas() uint64 {
	return 1_000
}

func (i *IsOracleOnlineMethod) IsReadonly() bool {
	return true
}

func (i *IsOracleOnlineMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := i.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}
	stateDB := evm.StateDB.(types.ExtStateDB)

	router, has := i.Keeper.router.GetRoute(args.Chain)
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	oracleAddr, has := router.GetOracleAddrByExternalAddr(stateDB.Context(), crosschaintypes.ExternalAddrToStr(args.Chain, args.ExternalAddress.Bytes()))
	if !has {
		return i.PackOutput(false)
	}

	oracle, has := router.GetOracle(stateDB.Context(), oracleAddr)
	return i.PackOutput(has && oracle.Online)
}

type IsOracleOnlineABI struct {
	abi.Method
}

func NewIsOracleOnlineABI() IsOracleOnlineABI {
	return IsOracleOnlineABI{
		Method: crosschainABI.Methods["isOracleOnline"],
	}
}

func (i *IsOracleOnlineMethod) PackInput(args contract.IsOracleOnlineArgs) ([]byte, error) {
	arguments, err := i.Method.Inputs.Pack(args.Chain, args.ExternalAddress)
	if err != nil {
		return nil, err
	}
	return append(i.GetMethodId(), arguments...), nil
}

func (i *IsOracleOnlineMethod) UnpackInput(data []byte) (*contract.IsOracleOnlineArgs, error) {
	args := new(contract.IsOracleOnlineArgs)
	err := types.ParseMethodArgs(i.Method, args, data[4:])
	return args, err
}

func (i *IsOracleOnlineMethod) PackOutput(result bool) ([]byte, error) {
	return i.Method.Outputs.Pack(result)
}

func (i *IsOracleOnlineMethod) UnpackOutput(data []byte) (bool, error) {
	result, err := i.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return result[0].(bool), err
}
