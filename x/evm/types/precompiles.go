package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	TypeAddress, _ = abi.NewType("address", "", nil)
	TypeUint256, _ = abi.NewType("uint256", "", nil)
	TypeString, _  = abi.NewType("string", "", nil)
	TypeBool, _    = abi.NewType("bool", "", nil)
	TypeBytes32, _ = abi.NewType("bytes32", "", nil)
)

func ParseMethodArgs(method abi.Method, v MethodArgs, data []byte) error {
	unpacked, err := method.Inputs.Unpack(data)
	if err != nil {
		return err
	}
	if err = method.Inputs.Copy(v, unpacked); err != nil {
		return err
	}
	return v.Validate()
}

func PackTopicData(event abi.Event, topics []common.Hash, args ...interface{}) ([]byte, []common.Hash, error) {
	data, err := event.Inputs.NonIndexed().Pack(args...)
	if err != nil {
		return nil, nil, fmt.Errorf("pack %s event error: %s", event.Name, err.Error())
	}
	newTopic := []common.Hash{event.ID}
	if len(topics) > 0 {
		newTopic = append(newTopic, topics...)
	}
	return data, newTopic, nil
}
