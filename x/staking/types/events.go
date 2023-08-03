package types

const (
	EventTypeApproveShares          = "approve_shares"
	EventTypeTransferShares         = "transfer_shares"
	EventTypeGrantPrivilege         = "grant_privilege"
	EventTypeEditConsensusPubKey    = "edit_consensus_pubkey"
	EventTypeEditingConsensusPubKey = "editing_consensus_pubkey"
	EventTypeEditedConsensusPubKey  = "edited_consensus_pubkey"

	AttributeKeyOwner     = "owner"
	AttributeKeySpender   = "spender"
	AttributeKeyShares    = "shares"
	AttributeKeyAmount    = "amount"
	AttributeKeyFrom      = "from"
	AttributeKeyRecipient = "recipient"
	AttributeKeyTo        = "to"
	AttributeKeyPubKey    = "pubkey"
	AttributeResult       = "result"
)
