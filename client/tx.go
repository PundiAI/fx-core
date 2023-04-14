package client

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/proto"
)

const DefGasLimit int64 = 200000

type buildTxClient interface {
	GetChainId() (string, error)
	QueryAccount(address string) (authtypes.AccountI, error)
	GetGasPrices() (sdk.Coins, error)
	EstimatingGas(raw *tx.TxRaw) (*sdk.GasInfo, error)
}

//gocyclo:ignore
func BuildTx(cli buildTxClient, privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error) {
	account, err := cli.QueryAccount(sdk.AccAddress(privKey.PubKey().Address()).String())
	if err != nil {
		return nil, err
	}
	chainId, err := cli.GetChainId()
	if err != nil {
		return nil, err
	}

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
		Memo:                        "",
		TimeoutHeight:               0,
		ExtensionOptions:            nil,
		NonCriticalExtensionOptions: nil,
	}
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		return nil, err
	}

	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}

	gasPrices, err := cli.GetGasPrices()
	if err != nil {
		return nil, err
	}
	var gasPrice sdk.Coin
	if len(gasPrices) > 0 {
		gasPrice = gasPrices[0]
	}

	authInfo := &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &tx.ModeInfo{
					Sum: &tx.ModeInfo_Single_{
						Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: account.GetSequence(),
			},
		},
		Fee: &tx.Fee{
			Amount:   sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.MulRaw(DefGasLimit))),
			GasLimit: uint64(DefGasLimit),
			Payer:    "",
			Granter:  "",
		},
	}

	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	signDoc := &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: account.GetAccountNumber(),
	}
	signatures, err := proto.Marshal(signDoc)
	if err != nil {
		return nil, err
	}
	sign, err := privKey.Sign(signatures)
	if err != nil {
		return nil, err
	}
	gasInfo, err := cli.EstimatingGas(&tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	})
	if err != nil {
		return nil, err
	}

	authInfo.Fee.GasLimit = gasInfo.GasUsed * 12 / 10
	authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.MulRaw(int64(authInfo.Fee.GasLimit))))

	signDoc.AuthInfoBytes, err = proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	signatures, err = proto.Marshal(signDoc)
	if err != nil {
		return nil, err
	}
	sign, err = privKey.Sign(signatures)
	if err != nil {
		return nil, err
	}
	return &tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}

func BuildTxV1(chainId string, sequence, accountNumber uint64, privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasPrice sdk.Coin, gasLimit int64, memo string, timeout uint64) (*tx.TxRaw, error) {
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
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		return nil, err
	}

	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}

	authInfo := &tx.AuthInfo{
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
			Amount:   sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.MulRaw(gasLimit))),
			GasLimit: uint64(gasLimit),
			Payer:    "",
			Granter:  "",
		},
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
