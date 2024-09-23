package types

import (
	"fmt"
	"math/big"
)

// mainnet
const (
	MainnetChainId       = "fxcore"
	mainnetEvmChainID    = 530
	MainnetGenesisHash   = "3d6cb5ebc05d42371581cd2f7bc23590e5623c3377424a89b1db982e1938fbad"
	MainnetBlockHeightV2 = 5_713_000
	MainnetBlockHeightV3 = 8_756_000
	MainnetBlockHeightV4 = 10_477_500
	MainnetBlockHeightV5 = 11_601_700
	MainnetBlockHeightV6 = 13_598_000
	MainnetBlockHeightV7 = 16_838_000
)

// testnet
const (
	TestnetChainId        = "dhobyghaut"
	testnetEvmChainID     = 90001
	TestnetGenesisHash    = "ec2bf940c025434d1fd17e2338a60b8803900310dc71f71ea55c185b24ddba23"
	TestnetBlockHeightV2  = 3_418_880
	TestnetBlockHeightV3  = 6_578_000
	TestnetBlockHeightV4  = 8_088_000
	TestnetBlockHeightV41 = 8_376_000 // v4.1
	TestnetBlockHeightV42 = 8_481_000 // v4.2
	TestnetBlockHeightV5  = 9_773_000
	TestnetBlockHeightV6  = 11_701_000
	TestnetBlockHeightV7  = 12_961_500
	TestnetBlockHeightV71 = 14_369_500 // v7.1
	TestnetBlockHeightV72 = 14_389_000 // v7.2
	TestnetBlockHeightV73 = 14_551_500 // v7.3
	TestnetBlockHeightV74 = 15_614_000 // v7.4
	TestnetBlockHeightV75 = 15_660_500 // v7.5
)

func EIP155ChainID(chainId string) *big.Int {
	if TestnetChainId == chainId {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
}

func ChainIdWithEIP155(chainId string) string {
	if TestnetChainId == chainId {
		return fmt.Sprintf("%s_%d-1", TestnetChainId, testnetEvmChainID)
	}
	return fmt.Sprintf("%s_%d-1", MainnetChainId, mainnetEvmChainID)
}
