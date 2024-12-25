package crosschain

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	evmtypes "github.com/pundiai/fx-core/v8/x/evm/types"
)

// Deprecated: After the upgrade to v8
type LegacyCrosschainMethod struct {
	*Keeper
	LegacyCrosschainABI
}

// Deprecated: After the upgrade to v8
func NewCrosschainMethod(keeper *Keeper) *LegacyCrosschainMethod {
	return &LegacyCrosschainMethod{
		Keeper:              keeper,
		LegacyCrosschainABI: NewCrosschainABI(),
	}
}

func (m *LegacyCrosschainMethod) IsReadonly() bool {
	return false
}

func (m *LegacyCrosschainMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *LegacyCrosschainMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *LegacyCrosschainMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	value := contract.Value()
	sender := contract.Caller()

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		fxTarget, err := crosschaintypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target))
		if err != nil {
			return err
		}
		crosschainKeeper, ok := m.router.GetRoute(fxTarget.GetModuleName())
		if !ok {
			return errors.New("invalid router")
		}
		if err = fxTarget.ValidateExternalAddr(args.Receipt); err != nil {
			return err
		}

		baseCoin := sdk.Coin{}
		totalAmount := big.NewInt(0).Add(args.Amount, args.Fee)

		isOriginToken := value.Sign() > 0
		if isOriginToken {
			if totalAmount.Cmp(value) != 0 {
				return errors.New("amount + fee not equal msg.value")
			}
			if !fxcontract.IsZeroEthAddress(args.Token) {
				return errors.New("token is not zero address")
			}

			baseCoin = sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(totalAmount))
			if err = m.bankKeeper.SendCoins(ctx, crosschainAddress.Bytes(), sender.Bytes(), sdk.NewCoins(baseCoin)); err != nil {
				return err
			}
		} else {
			baseCoin, err = m.EvmTokenToBaseCoin(ctx, evm, crosschainKeeper, sender, args.Token, totalAmount)
			if err != nil {
				return err
			}
		}

		amountCoin := sdk.NewCoin(baseCoin.Denom, sdkmath.NewIntFromBigInt(args.Amount))
		feeCoin := sdk.NewCoin(baseCoin.Denom, sdkmath.NewIntFromBigInt(args.Fee))
		if err = crosschainKeeper.CrosschainBaseCoin(ctx, sender.Bytes(), args.Receipt,
			amountCoin, feeCoin, fxTarget, args.Memo, isOriginToken); err != nil {
			return err
		}

		data, topic, err := m.NewCrosschainEvent(sender, args.Token, amountCoin.Denom, args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, crosschainAddress, data, topic)

		return nil
	}); err != nil {
		return nil, err
	}

	return m.PackOutput(true)
}

// Deprecated: After the upgrade to v8
type LegacyCrosschainABI struct {
	abi.Method
	abi.Event
}

// Deprecated: After the upgrade to v8
func NewCrosschainABI() LegacyCrosschainABI {
	return LegacyCrosschainABI{
		Method: crosschainABI.Methods["crossChain"],
		Event:  crosschainABI.Events["CrossChain"],
	}
}

func (m LegacyCrosschainABI) NewCrosschainEvent(sender, token common.Address, denom, receipt string, amount, fee *big.Int, target [32]byte, memo string) (data []byte, topic []common.Hash, err error) {
	return evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash(), token.Hash()}, denom, receipt, amount, fee, target, memo)
}

func (m LegacyCrosschainABI) UnpackInput(data []byte) (*fxcontract.CrosschainArgs, error) {
	args := new(fxcontract.CrosschainArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m LegacyCrosschainABI) PackInput(args fxcontract.CrosschainArgs) ([]byte, error) {
	data, err := m.Method.Inputs.Pack(args.Token, args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, data...), nil
}

func (m LegacyCrosschainABI) PackOutput(success bool) ([]byte, error) {
	return m.Method.Outputs.Pack(success)
}
