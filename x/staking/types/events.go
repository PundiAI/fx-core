package types

const (
	EventTypeApproveShares            = "approve_shares"
	EventTypeTransferShares           = "transfer_shares"
	EventTypeGrantPrivilege           = "grant_privilege"
	EventTypeEditConsensusPubKey      = "edit_consensus_pubkey"
	EventTypeStartEditConsensusPubKey = "start_edit_consensus_pubkey"
	EventTypeEndEditConsensusPubKey   = "end_edit_consensus_pubkey"

	AttributeKeyOwner     = "owner"
	AttributeKeySpender   = "spender"
	AttributeKeyShares    = "shares"
	AttributeKeyAmount    = "amount"
	AttributeKeyFrom      = "from"
	AttributeKeyRecipient = "recipient"
	AttributeKeyTo        = "to"
	AttributeKeyPubKey    = "pubkey"
	AttributeOldConsAddr  = "old_cons_addr"
	AttributeNewConsAddr  = "new_cons_addr"
)
