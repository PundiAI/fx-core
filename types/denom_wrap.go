package types

const (
	PundixWrapDenom = "pundix"
	PundixChannel   = "channel-0"
	PundixPort      = "transfer"
)

const (
	MainnetPundixUnWrapDenom = "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38"
)

const (
	TestnetPundixUnWrapDenom = "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B"
)

func IsPundixChannel(port, channel string) bool {
	return port == PundixPort && channel == PundixChannel
}

func GetPundixUnWrapDenom(chainId string) string {
	if chainId == MainnetChainId {
		return MainnetPundixUnWrapDenom
	}
	return TestnetPundixUnWrapDenom
}
