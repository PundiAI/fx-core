package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

const (
	ERC20EventTransfer = "Transfer"
)

type RelayTransfer struct {
	From          common.Address
	To            common.Address
	Amount        *big.Int
	TokenContract common.Address
	Validator     sdk.ValAddress
}

type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ParseTransferEvent transfer event ---> event Transfer(address indexed from, address indexed to, uint256 value);
func ParseTransferEvent(log *ethtypes.Log) (*TransferEvent, error) {
	// Note: the `Transfer` event contains 3 topics (id, from, to)
	if len(log.Topics) != 3 {
		return nil, nil
	}
	lpTokenABI := fxtypes.GetLPToken().ABI
	if log.Topics[0] != lpTokenABI.Events[ERC20EventTransfer].ID {
		return nil, nil
	}
	transferEvent := new(TransferEvent)
	if len(log.Data) > 0 {
		if err := lpTokenABI.UnpackIntoInterface(transferEvent, ERC20EventTransfer, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range lpTokenABI.Events[ERC20EventTransfer].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(transferEvent, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	return transferEvent, nil
}
