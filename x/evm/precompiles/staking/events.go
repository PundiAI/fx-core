package staking

var (
	ApproveSharesEvent  = GetABI().Events[ApproveSharesEventName]
	DelegateEvent       = GetABI().Events[DelegateEventName]
	TransferSharesEvent = GetABI().Events[TransferSharesEventName]
	UndelegateEvent     = GetABI().Events[UndelegateEventName]
	WithdrawEvent       = GetABI().Events[WithdrawEventName]
	RedelegateEvent     = GetABI().Events[RedelegateEventName]
)
