package types

import (
	"bytes"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CheckPubKey(pubKey cryptotypes.PubKey) error {
	if pubKey == nil {
		return nil
	}

	if bytes.Equal(pubKey.Bytes(), DisablePKBytes[:]) {
		return sdkerrors.ErrInvalidAddress.Wrap("account disabled")
	}

	return nil
}
