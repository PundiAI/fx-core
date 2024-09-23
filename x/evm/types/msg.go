package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/contract"
)

var _ sdk.Msg = &MsgCallContract{}

func (m *MsgCallContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := contract.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return errorsmod.Wrap(err, "contract address")
	}
	if len(m.Data) == 0 {
		return errors.New("empty data")
	}
	return nil
}
