package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/cometbft/cometbft/crypto/tmhash"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
)

type (
	// CrossChainMsg cross msg must implement GetChainName interface.. using in router
	CrossChainMsg interface {
		GetChainName() string
	}
)

var (
	_ CrossChainMsg = &MsgBondedOracle{}
	_ CrossChainMsg = &MsgAddDelegate{}
	_ CrossChainMsg = &MsgReDelegate{}
	_ CrossChainMsg = &MsgEditBridger{}
	_ CrossChainMsg = &MsgWithdrawReward{}
	_ CrossChainMsg = &MsgUnbondedOracle{}
	_ CrossChainMsg = &MsgOracleSetConfirm{}
	_ CrossChainMsg = &MsgOracleSetUpdatedClaim{}
	_ CrossChainMsg = &MsgBridgeTokenClaim{}
	_ CrossChainMsg = &MsgSendToFxClaim{}
	_ CrossChainMsg = &MsgSendToExternal{}
	_ CrossChainMsg = &MsgCancelSendToExternal{}
	_ CrossChainMsg = &MsgIncreaseBridgeFee{}
	_ CrossChainMsg = &MsgSendToExternalClaim{}
	_ CrossChainMsg = &MsgRequestBatch{}
	_ CrossChainMsg = &MsgConfirmBatch{}
	_ CrossChainMsg = &MsgBridgeCallClaim{}
	_ CrossChainMsg = &MsgBridgeCall{}
	_ CrossChainMsg = &MsgBridgeCallConfirm{}
	_ CrossChainMsg = &MsgBridgeCallResultClaim{}
	_ CrossChainMsg = &MsgUpdateParams{}
	_ CrossChainMsg = &MsgUpdateChainOracles{}
)

var (
	_ sdk.Msg = &MsgBondedOracle{}
	_ sdk.Msg = &MsgAddDelegate{}
	_ sdk.Msg = &MsgReDelegate{}
	_ sdk.Msg = &MsgEditBridger{}
	_ sdk.Msg = &MsgWithdrawReward{}
	_ sdk.Msg = &MsgUnbondedOracle{}
	_ sdk.Msg = &MsgOracleSetConfirm{}
	_ sdk.Msg = &MsgOracleSetUpdatedClaim{}
	_ sdk.Msg = &MsgBridgeTokenClaim{}
	_ sdk.Msg = &MsgSendToFxClaim{}
	_ sdk.Msg = &MsgSendToExternal{}
	_ sdk.Msg = &MsgCancelSendToExternal{}
	_ sdk.Msg = &MsgIncreaseBridgeFee{}
	_ sdk.Msg = &MsgSendToExternalClaim{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgConfirmBatch{}
	_ sdk.Msg = &MsgBridgeCallClaim{}
	_ sdk.Msg = &MsgBridgeCall{}
	_ sdk.Msg = &MsgBridgeCallConfirm{}
	_ sdk.Msg = &MsgBridgeCallResultClaim{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgUpdateChainOracles{}
	_ sdk.Msg = &MsgClaim{}
)

func (m *MsgBondedOracle) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if !m.DelegateAmount.IsValid() || m.DelegateAmount.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid delegation amount")
	}
	if m.OracleAddress == m.BridgerAddress {
		return sdkerrors.ErrInvalidRequest.Wrap("same address")
	}
	return nil
}

func (m *MsgAddDelegate) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

func (m *MsgReDelegate) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if _, err = sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	return nil
}

func (m *MsgEditBridger) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if _, err = sdk.ValAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if m.OracleAddress == m.BridgerAddress {
		return sdkerrors.ErrInvalidRequest.Wrap("same address")
	}
	return nil
}

func (m *MsgWithdrawReward) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	return nil
}

func (m *MsgUnbondedOracle) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	return nil
}

type Confirm interface {
	GetChainName() string
	GetSignature() string
	GetBridgerAddress() string
	GetExternalAddress() string
}

var (
	_ Confirm = &MsgBridgeCallConfirm{}
	_ Confirm = &MsgConfirmBatch{}
	_ Confirm = &MsgOracleSetConfirm{}
)

func (m *MsgOracleSetConfirm) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

func (m *MsgSendToExternal) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Dest); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid dest address: %s", err)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.ErrInvalidRequest.Wrap("bridge fee denom not equal amount denom")
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	return nil
}

func (m *MsgRequestBatch) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.Denom) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty denom")
	}
	if m.MinimumFee.IsNil() || !m.MinimumFee.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid minimum fee")
	}
	if err = ValidateExternalAddr(m.ChainName, m.FeeReceive); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid fee receive address: %s", err)
	}
	if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid base fee")
	}
	return nil
}

func (m *MsgConfirmBatch) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

func (m *MsgBridgeCallConfirm) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

func (m *MsgCancelSendToExternal) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero transaction id")
	}
	return nil
}

func (m *MsgIncreaseBridgeFee) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero transaction id")
	}
	if !m.AddBridgeFee.IsValid() || !m.AddBridgeFee.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	return nil
}

// ExternalClaim represents a claim on ethereum state
type ExternalClaim interface {
	proto.Message
	// GetEventNonce All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// GetBlockHeight The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// GetClaimer the delegate address of the claimer, for MsgSendToExternalClaim and MsgSendToFxClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// GetType Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ ExternalClaim = &MsgSendToFxClaim{}
	_ ExternalClaim = &MsgBridgeCallClaim{}
	_ ExternalClaim = &MsgBridgeTokenClaim{}
	_ ExternalClaim = &MsgSendToExternalClaim{}
	_ ExternalClaim = &MsgOracleSetUpdatedClaim{}
	_ ExternalClaim = &MsgBridgeCallResultClaim{}
)

func UnpackAttestationClaim(cdc codectypes.AnyUnpacker, att *Attestation) (ExternalClaim, error) {
	var msg ExternalClaim
	err := cdc.UnpackAny(att.Claim, &msg)
	return msg, err
}

func (m *MsgClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if m.Claim == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("empty claim")
	}
	claim, ok := m.Claim.GetCachedValue().(ExternalClaim)
	if !ok {
		return sdkerrors.ErrInvalidRequest.Wrapf("expected claim type %T, got %T", new(ExternalClaim), m.Claim.GetCachedValue())
	}
	return claim.ValidateBasic()
}

func (m *MsgClaim) GetSigners() []sdk.AccAddress {
	claim, ok := m.Claim.GetCachedValue().(ExternalClaim)
	if !ok {
		panic(sdkerrors.ErrInvalidRequest.Wrapf("expected claim type %T, got %T", new(ExternalClaim), m.Claim.GetCachedValue()))
	}
	return []sdk.AccAddress{claim.GetClaimer()}
}

func (m *MsgSendToFxClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_FX
}

func (m *MsgSendToFxClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not decode hex targetIbc")
	}
	if m.EventNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m *MsgSendToFxClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgSendToFxClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d%s/%s/%s/%s/%s", m.BlockHeight, m.EventNonce, m.TokenContract, m.Sender, m.Amount.String(), m.Receiver, m.TargetIbc)
	return tmhash.Sum([]byte(path))
}

func (m *MsgBridgeCallClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_CALL
}

func (m *MsgBridgeCallClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if len(m.TokenContracts) != len(m.Amounts) {
		return sdkerrors.ErrInvalidRequest.Wrap("mismatched token contracts and amounts")
	}
	for _, contract := range m.TokenContracts {
		if err = ValidateExternalAddr(m.ChainName, contract); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
		}
	}
	return m.validateBasic()
}

func (m *MsgBridgeCallClaim) validateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to contract: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Refund); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid refund address: %s", err)
	}
	if m.Value.IsNil() || m.Value.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid value")
	}
	if len(m.Data) > 0 {
		if _, err = hex.DecodeString(m.Data); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid data")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero block height")
	}
	if err = ValidateExternalAddr(m.ChainName, m.TxOrigin); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid tx origin: %s", err)
	}
	if len(m.Memo) > 0 {
		if _, err = hex.DecodeString(m.Memo); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid memo")
		}
	}
	return nil
}

func (m *MsgBridgeCallClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgBridgeCallClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%s/%v/%v/%s", m.BlockHeight, m.EventNonce, m.Sender, m.Refund, m.To, m.TokenContracts, m.Amounts, m.Data, m.Value.String())
	return tmhash.Sum([]byte(path))
}

func (m *MsgBridgeCallClaim) GetSenderAddr() common.Address {
	return ExternalAddrToHexAddr(m.ChainName, m.Sender)
}

func (m *MsgBridgeCallClaim) GetRefundAddr() common.Address {
	return ExternalAddrToHexAddr(m.ChainName, m.Refund)
}

func (m *MsgBridgeCallClaim) GetToAddr() common.Address {
	return ExternalAddrToHexAddr(m.ChainName, m.To)
}

func (m *MsgBridgeCallClaim) MustData() []byte {
	if len(m.Data) == 0 {
		return []byte{}
	}
	bz, err := hex.DecodeString(m.Data)
	if err != nil {
		panic(err)
	}
	return bz
}

func (m *MsgBridgeCallClaim) MustMemo() []byte {
	if len(m.Memo) == 0 {
		return []byte{}
	}
	bz, err := hex.DecodeString(m.Memo)
	if err != nil {
		panic(err)
	}
	return bz
}

func (m *MsgBridgeCallClaim) GetTokensAddr() []common.Address {
	addrs := make([]common.Address, 0, len(m.TokenContracts))
	for _, token := range m.TokenContracts {
		addr := ExternalAddrToAccAddr(m.ChainName, token)
		addrs = append(addrs, common.BytesToAddress(addr))
	}
	return addrs
}

func (m *MsgBridgeCallClaim) GetAmounts() []*big.Int {
	amts := make([]*big.Int, 0, len(m.Amounts))
	for _, a := range m.Amounts {
		amts = append(amts, a.BigInt())
	}
	return amts
}

func (m *MsgBridgeCallResultClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_CALL_RESULT
}

func (m *MsgBridgeCallResultClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if m.Nonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero nonce")
	}
	if m.EventNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero block height")
	}
	if err = ValidateExternalAddr(m.ChainName, m.TxOrigin); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid tx origin: %s", err)
	}
	if len(m.Cause) > 0 {
		if _, err = hex.DecodeString(m.Cause); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid cause")
		}
	}
	return nil
}

func (m *MsgBridgeCallResultClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgBridgeCallResultClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

func (m *MsgBridgeCallResultClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%t/%s", m.BlockHeight, m.EventNonce, m.Nonce, m.Success, m.Cause)
	return tmhash.Sum([]byte(path))
}

func (m *MsgSendToExternalClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_EXTERNAL
}

func (m *MsgSendToExternalClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if m.EventNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero block height")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero batch nonce")
	}
	return nil
}

func (m *MsgSendToExternalClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%d/", m.BlockHeight, m.EventNonce, m.TokenContract, m.BatchNonce)
	return tmhash.Sum([]byte(path))
}

func (m *MsgSendToExternalClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgBridgeTokenClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not decode hex channelIbc string")
	}
	if len(m.Name) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty token name")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty token symbol")
	}
	if m.EventNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m *MsgBridgeTokenClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgBridgeTokenClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_TOKEN
}

func (m *MsgBridgeTokenClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d%s/%s/%s/%d/%s/", m.BlockHeight, m.EventNonce, m.TokenContract, m.Name, m.Symbol, m.Decimals, m.ChannelIbc)
	return tmhash.Sum([]byte(path))
}

func (m *MsgOracleSetUpdatedClaim) GetType() ClaimType {
	return CLAIM_TYPE_ORACLE_SET_UPDATED
}

func (m *MsgOracleSetUpdatedClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if len(m.Members) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty members")
	}
	for _, member := range m.Members {
		if err = ValidateExternalAddr(m.ChainName, member.ExternalAddress); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
		}
		if member.Power == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("zero power")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m *MsgOracleSetUpdatedClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgOracleSetUpdatedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%v/", m.BlockHeight, m.OracleSetNonce, m.EventNonce, m.Members)
	return tmhash.Sum([]byte(path))
}

func (m *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	return nil
}

func (m *MsgSetOrchestratorAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

func (m *MsgAddOracleDeposit) ValidateBasic() (err error) {
	return nil
}

func (m *MsgAddOracleDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority")
	}
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("params err: %s", err.Error())
	}
	if len(m.Params.Oracles) > 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("deprecated oracles")
	}
	return nil
}

func (m *MsgUpdateChainOracles) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if len(m.Oracles) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty oracles")
	}
	oraclesMap := make(map[string]bool)
	for _, addr := range m.Oracles {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrap("oracle address")
		}
		if oraclesMap[addr] {
			return sdkerrors.ErrInvalidRequest.Wrap("duplicate oracle address")
		}
		oraclesMap[addr] = true
	}
	return nil
}

func (m *MsgBridgeCall) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	return m.validateBasic()
}

func (m *MsgBridgeCall) validateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}
	if m.Value.Sign() != 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("value must be zero")
	}
	if err = m.Coins.Validate(); err != nil {
		return sdkerrors.ErrInvalidCoins.Wrap(err.Error())
	}
	if len(m.Coins) > 0 || len(m.Refund) > 0 {
		if _, err = sdk.AccAddressFromBech32(m.Refund); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid refund address: %s", err)
		}
	}
	if len(m.Data) > 0 {
		if _, err = hex.DecodeString(m.Data); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid data")
		}
	}
	if len(m.Memo) > 0 {
		if _, err = hex.DecodeString(m.Memo); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid memo")
		}
	}
	if m.Coins.Len() == 0 && len(m.Data) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("coins and data cannot be empty at the same time")
	}
	return nil
}

func (m *MsgBridgeCall) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

func (m *MsgBridgeCall) GetSenderAddr() common.Address {
	return common.BytesToAddress(sdk.MustAccAddressFromBech32(m.Sender).Bytes())
}

func (m *MsgBridgeCall) GetRefundAddr() common.Address {
	if len(m.Refund) == 0 {
		return common.Address{}
	}
	return common.BytesToAddress(sdk.MustAccAddressFromBech32(m.Refund).Bytes())
}

func (m *MsgBridgeCall) GetToAddr() common.Address {
	if m.To == "" {
		return common.Address{}
	}
	return ExternalAddrToHexAddr(m.ChainName, m.To)
}

func (m *MsgBridgeCall) MustData() []byte {
	if len(m.Data) == 0 {
		return []byte{}
	}
	bz, err := hex.DecodeString(m.Data)
	if err != nil {
		panic(err)
	}
	return bz
}

func (m *MsgBridgeCall) MustMemo() []byte {
	if len(m.Memo) == 0 {
		return []byte{}
	}
	bz, err := hex.DecodeString(m.Memo)
	if err != nil {
		panic(err)
	}
	return bz
}
