package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// DefEthereumConfig returns an Ethereum ChainConfig for EVM state transitions.
// All the negative or nil values are converted to nil
func DefEthereumConfig(chainID *big.Int) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:                 chainID,
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            big.NewInt(0),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP150Hash:              common.Hash{},
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		MergeForkBlock:          big.NewInt(0),
		TerminalTotalDifficulty: nil,
		Ethash:                  nil,
		Clique:                  nil,
	}
}
