package crosschain

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	crosschaintypes "github.com/functionx/fx-core/v4/x/crosschain/types"
	"github.com/functionx/fx-core/v4/x/evm/types"
)

var (
	CancelSendToExternalMethod = abi.NewMethod(
		CancelSendToExternalMethodName,
		CancelSendToExternalMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_chain", Type: types.TypeString},
			abi.Argument{Name: "_txID", Type: types.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: types.TypeBool},
		},
	)

	CancelSendToExternalEvent = abi.NewEvent(
		CancelSendToExternalEventName,
		CancelSendToExternalEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "chain", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "txID", Type: types.TypeUint256, Indexed: false},
		})
)

type CancelSendToExternalArgs struct {
	Chain string   `abi:"_chain"`
	TxID  *big.Int `abi:"_txID"`
}

func (c *Contract) CancelSendToExternal(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("cancel send to external method not readonly")
	}

	// args
	var args CancelSendToExternalArgs
	if err := ParseMethodParams(CancelSendToExternalMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	if err := crosschaintypes.ValidateModuleName(args.Chain); err != nil {
		return nil, err
	}
	if args.TxID.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid tx id: %s", args.TxID.String())
	}

	sender := contract.Caller()
	route, has := c.router.GetRoute(args.Chain)
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	originDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	// NOTE: must be get relation before cancel, cancel will delete it if relation exist
	hasRelation := c.erc20Keeper.HasOutgoingTransferRelation(ctx, args.TxID.Uint64())

	refundCoin, err := route.PrecompileCancelSendToExternal(ctx, args.TxID.Uint64(), sender.Bytes())
	if err != nil {
		return nil, err
	}
	if !hasRelation && refundCoin.Denom == originDenom {
		// add refund to sender in evm state db, because bank keeper add refund to sender
		evm.StateDB.AddBalance(sender, refundCoin.Amount.BigInt())
	}

	// add event log
	if err := c.AddLog(CancelSendToExternalEvent, []common.Hash{sender.Hash()}, args.Chain, args.TxID); err != nil {
		return nil, err
	}

	return CancelSendToExternalMethod.Outputs.Pack(true)
}
