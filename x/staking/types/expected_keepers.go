package types

import (
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

type AuthzKeeper interface {
	SaveGrant(ctx sdk.Context, grantee, granter sdk.AccAddress, authorization authz.Authorization, expiration *time.Time) error
	DeleteGrant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType string) error
	GetAuthorizations(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress) ([]authz.Authorization, error)
	GetAuthorization(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType string) (authz.Authorization, *time.Time)
}

type SlashingKeeper interface {
	AddPubkey(ctx sdk.Context, pubkey cryptotypes.PubKey) error
	DeleteConsensusPubKey(ctx sdk.Context, consAddr sdk.ConsAddress)
	GetValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress) (info slashingtypes.ValidatorSigningInfo, found bool)
	SetValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo)
	DeleteValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress)
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}
