package tests

import "github.com/cosmos/cosmos-sdk/crypto/hd"

const (
	GrpcUrl        = "localhost:9090"
	NodeJsonRpcUrl = "tcp://localhost:26657"
	RestUrl        = "tcp://localhost:1317"
	Web3Url        = "http://localhost:8545"
	Mnemonic       = "dune antenna hood magic kit blouse film video another pioneer dilemma hobby message rug sail gas culture upgrade twin flag joke people general aunt"
)

func GetGrpcUrl() string {
	return GrpcUrl
}

func GetNodeJsonRpcUrl() string {
	return NodeJsonRpcUrl
}

func GetRestUrl() string {
	return RestUrl
}

func GetAdminMnemonic() (string, hd.PubKeyType) {
	return Mnemonic, hd.Secp256k1Type
}

func GetWeb3Url() string {
	return Web3Url
}
