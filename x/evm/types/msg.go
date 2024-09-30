package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/v8/contract"
)

var _ sdk.Msg = &MsgCallContract{}

func (m *MsgCallContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if err := contract.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("contract address")
	}
	if len(m.Data) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("data is empty")
	}
	return nil
}
