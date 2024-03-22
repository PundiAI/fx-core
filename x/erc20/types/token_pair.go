package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/functionx/fx-core/v7/contract"
)

// NewTokenPair returns an instance of TokenPair
func NewTokenPair(erc20Address common.Address, denom string, enabled bool, contractOwner Owner) TokenPair {
	return TokenPair{
		Erc20Address:  erc20Address.String(),
		Denom:         denom,
		Enabled:       enabled,
		ContractOwner: contractOwner,
	}
}

// GetID returns the SHA256 hash of the ERC20 address and denomination
func (tp *TokenPair) GetID() []byte {
	return tmhash.Sum([]byte(fmt.Sprintf("%s|%s", tp.Erc20Address, tp.Denom)))
}

// GetERC20Contract casts the hex string address of the ERC20 to common.Address
func (tp *TokenPair) GetERC20Contract() common.Address {
	return common.HexToAddress(tp.Erc20Address)
}

// Validate performs a stateless validation of a TokenPair
func (tp *TokenPair) Validate() error {
	if err := sdk.ValidateDenom(tp.Denom); err != nil {
		return err
	}

	if err := contract.ValidateEthereumAddress(tp.Erc20Address); err != nil {
		return err
	}

	return nil
}

// IsNativeCoin returns true if the owner of the ERC20 contract is the
// erc20 module account
func (tp *TokenPair) IsNativeCoin() bool {
	return tp.ContractOwner == OWNER_MODULE
}

// IsNativeERC20 returns true if the owner of the ERC20 contract not the
// erc20 module account
func (tp *TokenPair) IsNativeERC20() bool {
	return tp.ContractOwner == OWNER_EXTERNAL
}
