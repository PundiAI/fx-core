package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

type BridgeCoinAmountMethod struct {
	*Keeper
	abi.Method
}

func NewBridgeCoinAmountMethod(keeper *Keeper) *BridgeCoinAmountMethod {
	return &BridgeCoinAmountMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["bridgeCoinAmount"],
	}
}

func (m *BridgeCoinAmountMethod) IsReadonly() bool {
	return true
}

func (m *BridgeCoinAmountMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *BridgeCoinAmountMethod) RequiredGas() uint64 {
	return 10_000
}

func (m *BridgeCoinAmountMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	ctx := stateDB.Context()

	pair, has := m.erc20Keeper.GetTokenPair(ctx, args.Token.Hex())
	if !has {
		return nil, fmt.Errorf("token not support: %s", args.Token.Hex())
	}
	// FX
	if fxcontract.IsZeroEthAddress(args.Token) {
		supply := m.bankKeeper.GetSupply(ctx, fxtypes.DefaultDenom)
		balance := m.bankKeeper.GetBalance(ctx, m.accountKeeper.GetModuleAddress(ethtypes.ModuleName), fxtypes.DefaultDenom)
		return m.PackOutput(supply.Amount.Sub(balance.Amount).BigInt())
	}
	// OriginDenom
	if m.erc20Keeper.IsOriginDenom(ctx, pair.GetDenom()) {
		erc20Call := fxcontract.NewERC20Call(evm, crosschaintypes.GetAddress(), args.Token, m.RequiredGas())
		supply, err := erc20Call.TotalSupply()
		if err != nil {
			return nil, err
		}
		return m.PackOutput(supply)
	}
	// one to one
	_, has = m.erc20Keeper.HasDenomAlias(ctx, pair.GetDenom())
	if !has && pair.GetDenom() != fxtypes.DefaultDenom {
		return m.PackOutput(
			m.bankKeeper.GetSupply(ctx, pair.GetDenom()).Amount.BigInt(),
		)
	}
	// many to one
	md, has := m.bankKeeper.GetDenomMetaData(ctx, pair.GetDenom())
	if !has {
		return nil, fmt.Errorf("denom not support: %s", pair.GetDenom())
	}
	denom := m.erc20Keeper.ToTargetDenom(
		ctx,
		pair.GetDenom(),
		md.GetBase(),
		md.GetDenomUnits()[0].GetAliases(),
		fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target)),
	)

	balance := m.bankKeeper.GetBalance(ctx, m.erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
	supply := m.bankKeeper.GetSupply(ctx, denom)
	if balance.Amount.LT(supply.Amount) {
		supply = balance
	}
	return m.PackOutput(supply.Amount.BigInt())
}

func (m *BridgeCoinAmountMethod) PackInput(args crosschaintypes.BridgeCoinAmountArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Token, args.Target)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *BridgeCoinAmountMethod) UnpackInput(data []byte) (*crosschaintypes.BridgeCoinAmountArgs, error) {
	args := new(crosschaintypes.BridgeCoinAmountArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *BridgeCoinAmountMethod) PackOutput(amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount)
}

func (m *BridgeCoinAmountMethod) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}
