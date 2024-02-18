package client

import (
	"context"
	"fmt"
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BuildEthTransaction(ctx context.Context, cli *ethclient.Client, priKey cryptotypes.PrivKey, to *common.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	sender := common.BytesToAddress(priKey.PubKey().Address().Bytes())

	chainId, err := cli.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := cli.NonceAt(ctx, sender, nil)
	if err != nil {
		return nil, err
	}
	head, err := cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	var gasTipCap, gasFeeCap, gasPrice *big.Int
	if head.BaseFee != nil {
		tip, err := cli.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, err
		}
		gasTipCap = tip
		gasFeeCap = new(big.Int).Add(tip, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))
		if gasFeeCap.Cmp(gasTipCap) < 0 {
			return nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
		}
	} else {
		gasPrice, err = cli.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
	}

	msg := ethereum.CallMsg{From: sender, To: to, GasPrice: gasPrice, GasTipCap: gasTipCap, GasFeeCap: gasFeeCap, Value: value, Data: data}
	gasLimit, err := cli.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas needed: %w", err)
	}
	gasLimit = gasLimit * 130 / 100
	if value == nil {
		value = big.NewInt(0)
	}

	var rawTx *ethtypes.Transaction
	if gasFeeCap == nil {
		baseTx := &ethtypes.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       to,
			Value:    value,
			Data:     data,
		}
		rawTx = ethtypes.NewTx(baseTx)
	} else {
		baseTx := &ethtypes.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     nonce,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Gas:       gasLimit,
			To:        to,
			Value:     value,
			Data:      data,
		}
		rawTx = ethtypes.NewTx(baseTx)
	}
	signer := ethtypes.NewLondonSigner(chainId)
	signature, err := priKey.Sign(signer.Hash(rawTx).Bytes())
	if err != nil {
		return nil, err
	}
	return rawTx.WithSignature(signer, signature)
}
