package types

import (
	"encoding/hex"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
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

// ToIBCPacketData is a helper for serializing
func (ftpd FungibleTokenPacketData) ToIBCPacketData() transfertypes.FungibleTokenPacketData {
	return transfertypes.NewFungibleTokenPacketData(ftpd.Denom, ftpd.Amount, ftpd.Sender, ftpd.Receiver, ftpd.Memo)
}

// ValidateBasic is used for validating the token transfer.
// NOTE: The addresses formats are not validated as the sender and recipient can have different
// formats defined by their corresponding chains that are not known to IBC.
func (ftpd FungibleTokenPacketData) ValidateBasic() error {
	amount, ok := sdkmath.NewIntFromString(ftpd.Amount)
	if !ok {
		return transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount (%s) into sdkmath.Int", ftpd.Amount)
	}
	if !amount.IsPositive() {
		return transfertypes.ErrInvalidAmount.Wrapf("amount must be strictly positive: got %s", amount)
	}
	if strings.TrimSpace(ftpd.Sender) == "" {
		return sdkerrors.ErrInvalidAddress.Wrap("sender address cannot be blank")
	}
	if strings.TrimSpace(ftpd.Receiver) == "" {
		return sdkerrors.ErrInvalidAddress.Wrap("receiver address cannot be blank")
	}
	fee, ok := sdkmath.NewIntFromString(ftpd.Fee)
	if !ok {
		return transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer fee (%s) into sdkmath.Int", ftpd.Fee)
	}
	if fee.IsNegative() {
		return transfertypes.ErrInvalidAmount.Wrapf("fee must be strictly not negative: got %s", fee)
	}
	return transfertypes.ValidatePrefixedDenom(ftpd.Denom)
}
