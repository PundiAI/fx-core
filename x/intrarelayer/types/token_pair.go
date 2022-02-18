package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/functionx/fx-core/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// NewTokenPair returns an instance of TokenPair
func NewTokenPair(Fip20Address common.Address, denom string, enabled bool, contractOwner Owner) TokenPair {
	return TokenPair{
		Fip20Address:  Fip20Address.String(),
		Denom:         denom,
		Enabled:       true,
		ContractOwner: contractOwner,
	}
}

// GetID returns the SHA256 hash of the ERC20 address and denomination
func (tp TokenPair) GetID() []byte {
	id := tp.Fip20Address + "|" + tp.Denom
	return tmhash.Sum([]byte(id))
}

// GetFIP20Contract casts the hex string address of the FIP20 to common.Address
func (tp TokenPair) GetFIP20Contract() common.Address {
	return common.HexToAddress(tp.Fip20Address)
}

// Validate performs a stateless validation of a TokenPair
func (tp TokenPair) Validate() error {
	if err := sdk.ValidateDenom(tp.Denom); err != nil {
		return err
	}

	if err := ethermint.ValidateAddress(tp.Fip20Address); err != nil {
		return err
	}

	return nil
}

// IsNativeCoin returns true if the owner of the ERC20 contract is the
// intrarelayer module account
func (tp TokenPair) IsNativeCoin() bool {
	return tp.ContractOwner == OWNER_MODULE
}

// IsNativeERC20 returns true if the owner of the ERC20 contract not the
// intrarelayer module account
func (tp TokenPair) IsNativeERC20() bool {
	return tp.ContractOwner == OWNER_EXTERNAL
}
