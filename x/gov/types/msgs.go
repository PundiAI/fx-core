package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateStore{}
	_ sdk.Msg = &MsgUpdateSwitchParams{}
	_ sdk.Msg = &MsgUpdateCustomParams{}
)

func NewMsgUpdateStore(authority string, updateStores []UpdateStore) *MsgUpdateStore {
	return &MsgUpdateStore{Authority: authority, UpdateStores: updateStores}
}

func (m *MsgUpdateStore) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if len(m.UpdateStores) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("stores are empty")
	}
	for _, updateStore := range m.UpdateStores {
		if len(updateStore.Space) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("store space is empty")
		}
		if len(updateStore.Key) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("store key is empty")
		}
		if _, err := hex.DecodeString(updateStore.Key); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid store key")
		}
		if len(updateStore.OldValue) > 0 {
			if _, err := hex.DecodeString(updateStore.OldValue); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrap("invalid old store value")
			}
		}
		if len(updateStore.Value) > 0 {
			if _, err := hex.DecodeString(updateStore.Value); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrap("invalid store value")
			}
		}
	}
	return nil
}

func (us *UpdateStore) String() string {
	out, _ := json.Marshal(us)
	return string(out)
}

func (us *UpdateStore) KeyToBytes() []byte {
	b, err := hex.DecodeString(us.Key)
	if err != nil {
		panic(err)
	}
	return b
}

func (us *UpdateStore) OldValueToBytes() []byte {
	if len(us.OldValue) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(us.OldValue)
	if err != nil {
		panic(err)
	}
	return b
}

func (us *UpdateStore) ValueToBytes() []byte {
	if len(us.Value) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(us.Value)
	if err != nil {
		panic(err)
	}
	return b
}

func (m *MsgUpdateSwitchParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("params err: %s", err.Error())
	}
	return nil
}

// ValidateBasic performs basic validation on governance parameters.
func (p *CustomParams) ValidateBasic() error {
	if p.VotingPeriod == nil {
		return fmt.Errorf("voting period must not be nil: %d", p.VotingPeriod)
	}
	if p.VotingPeriod.Seconds() <= 0 {
		return fmt.Errorf("voting period must be positive: %s", p.VotingPeriod)
	}

	quorum, err := sdkmath.LegacyNewDecFromStr(p.Quorum)
	if err != nil {
		return fmt.Errorf("invalid quorum string: %w", err)
	}
	if quorum.IsNegative() {
		return fmt.Errorf("quorum cannot be negative: %s", quorum)
	}
	if quorum.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("quorum too large: %s", p.Quorum)
	}

	depositRatio, err := sdkmath.LegacyNewDecFromStr(p.DepositRatio)
	if err != nil {
		return fmt.Errorf("invalid depositRatio string: %w", err)
	}
	if depositRatio.IsNegative() {
		return fmt.Errorf("depositRatio cannot be negative: %s", depositRatio)
	}
	if depositRatio.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("depositRatio too large: %s", p.DepositRatio)
	}
	return nil
}
