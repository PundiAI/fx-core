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

func (m *IbcCallEvmPacket) GetType() IbcCallType {
	return IBC_CALL_TYPE_EVM
}

func (m *IbcCallEvmPacket) ValidateBasic() error {
	if err := contract.ValidateEthereumAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("to address: %s", err.Error())
	}
	if m.Value.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf("value: %s", m.Value.String())
	}
	if _, err := hex.DecodeString(m.Data); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("data: %s", err.Error())
	}
	return nil
}

func (m *IbcCallEvmPacket) GetToAddress() *common.Address {
	to := common.HexToAddress(m.To)
	return &to
}

func (m *IbcCallEvmPacket) MustGetData() []byte {
	bz, err := hex.DecodeString(m.Data)
	if err != nil {
		panic(err)
	}
	return bz
}
