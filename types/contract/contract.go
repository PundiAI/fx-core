package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ParseLogEvent todo: remove unused code
func ParseLogEvent(eventABI abi.ABI, log *ethtypes.Log, eventName string, res interface{}) error {
	if len(log.Data) > 0 {
		if err := eventABI.UnpackIntoInterface(res, eventName, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eventABI.Events[eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(res, indexed, log.Topics[1:]); err != nil {
		return err
	}
	return nil
}
