package types

import (
	"encoding/hex"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"

	"github.com/functionx/fx-core/v7/contract"
)

type MemoPacket interface {
	proto.Message
	GetType() IbcCallType
	ValidateBasic() error
}

var _ MemoPacket = &IbcCallEvmPacket{}

// DefaultRelativePacketTimeoutTimestamp is the default packet timeout timestamp (in nanoseconds)
// relative to the current block timestamp of the counterparty chain provided by the client
// state. The timeout is disabled when set to 0. The default is currently set to a 12-hour
// timeout.
var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(12) * time.Hour).Nanoseconds())

// NewFungibleTokenPacketData constructs a new FungibleTokenPacketData instance
func NewFungibleTokenPacketData(denom, amount, sender, receiver, router string, fee string,
) FungibleTokenPacketData {
	return FungibleTokenPacketData{
		Denom:    denom,
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
		Router:   router,
		Fee:      fee,
	}
}

// ValidateBasic is used for validating the token transfer.
// NOTE: The addresses formats are not validated as the sender and recipient can have different
// formats defined by their corresponding chains that are not known to IBC.
func (ftpd FungibleTokenPacketData) ValidateBasic() error {
	amount, ok := sdkmath.NewIntFromString(ftpd.Amount)
	if !ok {
		return errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount (%s) into sdkmath.Int", ftpd.Amount)
	}
	if !amount.IsPositive() {
		return errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "amount must be strictly positive: got %d", amount)
	}
	if strings.TrimSpace(ftpd.Sender) == "" {
		return errorsmod.Wrap(errortypes.ErrInvalidAddress, "sender address cannot be blank")
	}
	if strings.TrimSpace(ftpd.Receiver) == "" {
		return errorsmod.Wrap(errortypes.ErrInvalidAddress, "receiver address cannot be blank")
	}
	fee, ok := sdkmath.NewIntFromString(ftpd.Fee)
	if !ok {
		return errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer fee (%s) into sdkmath.Int", ftpd.Fee)
	}
	if fee.IsNegative() {
		return errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "fee must be strictly not negative: got %d", fee)
	}
	if err := transfertypes.ValidatePrefixedDenom(ftpd.Denom); err != nil {
		return err
	}
	return nil
}

// GetBytes is a helper for serialising
func (ftpd FungibleTokenPacketData) GetBytes() []byte {
	if ftpd.Router == "" {
		ftpd.Fee = ""
	}
	return sdk.MustSortJSON(mustProtoMarshalJSON(&ftpd))
}

// ToIBCPacketData is a helper for serialising
func (ftpd FungibleTokenPacketData) ToIBCPacketData() transfertypes.FungibleTokenPacketData {
	result := transfertypes.NewFungibleTokenPacketData(ftpd.Denom, ftpd.Amount, ftpd.Sender, ftpd.Receiver, ftpd.Memo)
	result.Memo = ftpd.Memo
	return result
}

func (icep IbcCallEvmPacket) GetType() IbcCallType {
	return IBC_CALL_TYPE_EVM
}

func (icep IbcCallEvmPacket) ValidateBasic() error {
	if err := contract.ValidateEthereumAddress(icep.To); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("to address: %s", err.Error())
	}
	if icep.Value.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrapf("value: %s", icep.Value.String())
	}
	if _, err := hex.DecodeString(icep.Data); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("data: %s", err.Error())
	}
	return nil
}

func (icep IbcCallEvmPacket) GetToAddress() *common.Address {
	to := common.HexToAddress(icep.To)
	return &to
}

func (icep IbcCallEvmPacket) MustGetMessage() []byte {
	bz, err := hex.DecodeString(icep.Data)
	if err != nil {
		panic(err)
	}
	return bz
}
