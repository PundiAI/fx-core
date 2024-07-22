package client

import (
	"errors"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
)

const DefGasLimit uint64 = 200000

type buildTxClient interface {
	GetChainId() (string, error)
	QueryAccount(address string) (authtypes.AccountI, error)
	GetGasPrices() (sdk.Coins, error)
	GetAddressPrefix() (string, error)
	EstimatingGas(raw *tx.TxRaw) (*sdk.GasInfo, error)
}

func BuildTxRawWithCli(cli buildTxClient, privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasLimit, timeout uint64, memo string) (*tx.TxRaw, error) {
	if gasLimit == 0 {
		gasLimit = DefGasLimit
	}
	prefix, err := cli.GetAddressPrefix()
	if err != nil {
		return nil, err
	}
	from, err := bech32.ConvertAndEncode(prefix, privKey.PubKey().Address())
	if err != nil {
		return nil, err
	}
	account, chainId, gasPrice, err := GetChainInfo(cli, from)
	if err != nil {
		return nil, err
	}
	txRaw, err := BuildTxRaw(chainId, account.GetSequence(), account.GetAccountNumber(), privKey, msgs, gasPrice, gasLimit, timeout, memo)
	if err != nil {
		return nil, err
	}
	estimatingGas, err := cli.EstimatingGas(txRaw)
	if err != nil {
		return nil, err
	}
	if estimatingGas.GetGasUsed() > gasLimit {
		gasLimit = estimatingGas.GetGasUsed() + (estimatingGas.GetGasUsed())*2/10
	}
	return BuildTxRaw(chainId, account.GetSequence(), account.GetAccountNumber(), privKey, msgs, gasPrice, gasLimit, timeout, memo)
}

func BuildTxRaw(chainId string, sequence, accountNumber uint64, privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasPrice sdk.Coin, gasLimit, timeout uint64, memo string) (*tx.TxRaw, error) {
	txBody, err := NewTxBody(msgs, memo, timeout)
	if err != nil {
		return nil, err
	}
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		return nil, err
	}

	authInfo, err := NewAuthInfo(privKey.PubKey(), sequence, gasLimit, gasPrice)
	if err != nil {
		return nil, err
	}
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}

	signDoc := &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}
	signatures, err := proto.Marshal(signDoc)
	if err != nil {
		return nil, err
	}
	sign, err := privKey.Sign(signatures)
	if err != nil {
		return nil, err
	}
	return &tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}

type waitMinedClient interface {
	TxByHash(txHash string) (*sdk.TxResponse, error)
}

func WaitMined(cli waitMinedClient, txHash string, timeout, pollInterval time.Duration) (*sdk.TxResponse, error) {
	for i := int64(0); i < int64(timeout/pollInterval); i++ {
		time.Sleep(pollInterval)
		txResponse, err := cli.TxByHash(txHash)
		if err != nil && strings.Contains(err.Error(), "not found") {
			continue
		}
		if err != nil {
			return nil, err
		}
		if txResponse != nil {
			return txResponse, nil
		}
	}
	return nil, errors.New("waiting for tx timeout")
}

func GetChainInfo(cli buildTxClient, from string) (account authtypes.AccountI, chainId string, gasPrice sdk.Coin, err error) {
	account, err = cli.QueryAccount(from)
	if err != nil {
		return nil, "", sdk.Coin{}, err
	}
	chainId, err = cli.GetChainId()
	if err != nil {
		return nil, "", sdk.Coin{}, err
	}
	gasPrices, err := cli.GetGasPrices()
	if err != nil {
		return nil, "", sdk.Coin{}, err
	}
	if len(gasPrices) > 0 {
		gasPrice = gasPrices[0]
	}
	return account, chainId, gasPrice, nil
}

func NewTxBody(msgs []sdk.Msg, memo string, timeout uint64) (*tx.TxBody, error) {
	txBodyMessage := make([]*types.Any, 0)
	for i := 0; i < len(msgs); i++ {
		msgAnyValue, err := types.NewAnyWithValue(msgs[i])
		if err != nil {
			return nil, err
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)
	}

	txBody := &tx.TxBody{
		Messages:                    txBodyMessage,
		Memo:                        memo,
		TimeoutHeight:               timeout,
		ExtensionOptions:            nil,
		NonCriticalExtensionOptions: nil,
	}
	return txBody, nil
}

func NewAuthInfo(pubKey cryptotypes.PubKey, sequence, gasLimit uint64, gasPrice sdk.Coin) (*tx.AuthInfo, error) {
	pubAny, err := types.NewAnyWithValue(pubKey)
	if err != nil {
		return nil, err
	}
	return &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &tx.ModeInfo{
					Sum: &tx.ModeInfo_Single_{
						Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: sequence,
			},
		},
		Fee: &tx.Fee{
			Amount:   sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.MulRaw(int64(gasLimit)))),
			GasLimit: gasLimit,
			Payer:    "",
			Granter:  "",
		},
	}, nil
}
