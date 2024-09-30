package types

import (
	"encoding/hex"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
)

type MemoPacket interface {
	proto.Message
	GetType() IbcCallType
	ValidateBasic() error
}

var _ MemoPacket = &IbcCallEvmPacket{}

func (icep IbcCallEvmPacket) GetType() IbcCallType {
	return IBC_CALL_TYPE_EVM
}

func (icep IbcCallEvmPacket) ValidateBasic() error {
	if err := contract.ValidateEthereumAddress(icep.To); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("to address: %s", err.Error())
	}
	if icep.Value.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf("value: %s", icep.Value.String())
	}
	if _, err := hex.DecodeString(icep.Data); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("data: %s", err.Error())
	}
	return nil
}

func (icep IbcCallEvmPacket) GetToAddress() *common.Address {
	to := common.HexToAddress(icep.To)
	return &to
}

func (icep IbcCallEvmPacket) MustGetData() []byte {
	bz, err := hex.DecodeString(icep.Data)
	if err != nil {
		panic(err)
	}
	return bz
}
