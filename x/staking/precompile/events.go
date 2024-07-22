package precompile

var (
	ApproveSharesEvent  = GetABI().Events[ApproveSharesEventName]
	DelegateEvent       = GetABI().Events[DelegateEventName]
	DelegateV2Event     = GetABI().Events[DelegateV2EventName]
	TransferSharesEvent = GetABI().Events[TransferSharesEventName]
	UndelegateEvent     = GetABI().Events[UndelegateEventName]
	UndelegateV2Event   = GetABI().Events[UndelegateV2EventName]
	WithdrawEvent       = GetABI().Events[WithdrawEventName]
	RedelegateEvent     = GetABI().Events[RedelegateEventName]
	RedelegateV2Event   = GetABI().Events[RedelegateV2EventName]
)
