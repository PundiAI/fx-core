package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/x/authz"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	tenderminttypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/google"

	otherTypes "github.com/functionx/fx-core/x/other/types"
)

const DefGasLimit int64 = 200000

type Client struct {
	chainId   string
	gasPrices sdk.Coins
	ctx       context.Context
	*grpc.ClientConn
}

func NewGrpcConn(rawUrl string) (*grpc.ClientConn, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	_url := u.Host
	if u.Port() == "" {
		if u.Scheme == "https" {
			_url = u.Host + ":443"
		} else {
			_url = u.Host + ":80"
		}
	}
	var opts []grpc.DialOption
	if u.Scheme == "https" {
		opts = append(opts, grpc.WithCredentialsBundle(google.NewDefaultCredentials()))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(_url, opts...)
}

func NewClient(rawUrl string) (*Client, error) {
	grpcConn, err := NewGrpcConn(rawUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		ClientConn: grpcConn,
		ctx:        context.Background(),
	}, nil
}

func (cli *Client) WithContext(ctx context.Context) {
	cli.ctx = ctx
}

func (cli *Client) WithGasPrices(gasPrices sdk.Coins) {
	cli.gasPrices = gasPrices
}

func (cli *Client) AuthQuery() authtypes.QueryClient {
	return authtypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) AuthzQuery() authz.QueryClient {
	return authz.NewQueryClient(cli.ClientConn)
}
func (cli *Client) BankQuery() banktypes.QueryClient {
	return banktypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) DistrQuery() distrtypes.QueryClient {
	return distrtypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) EvidenceQuery() evidencetypes.QueryClient {
	return evidencetypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) FeegrantQuery() feegrant.QueryClient {
	return feegrant.NewQueryClient(cli.ClientConn)
}
func (cli *Client) GovQuery() govtypes.QueryClient {
	return govtypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) MintQuery() minttypes.QueryClient {
	return minttypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) SlashingQuery() slashingtypes.QueryClient {
	return slashingtypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) StakingQuery() stakingtypes.QueryClient {
	return stakingtypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) UpgradeQuery() upgradetypes.QueryClient {
	return upgradetypes.NewQueryClient(cli.ClientConn)
}
func (cli *Client) ServiceClient() tx.ServiceClient {
	return tx.NewServiceClient(cli.ClientConn)
}
func (cli *Client) TMServiceClient() tmservice.ServiceClient {
	return tmservice.NewServiceClient(cli.ClientConn)
}

func (cli *Client) AppVersion() (string, error) {
	info, err := cli.TMServiceClient().GetNodeInfo(cli.ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return "", err
	}
	return info.ApplicationVersion.Version, nil
}

func (cli *Client) QueryAccount(address string) (authtypes.AccountI, error) {
	response, err := cli.AuthQuery().Account(cli.ctx, &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}
	var account authtypes.AccountI
	if err = newInterfaceRegistry().UnpackAny(response.GetAccount(), &account); err != nil {
		return nil, err
	}
	return account, nil
}

func (cli *Client) QueryBalance(address string, denom string) (sdk.Coin, error) {
	response, err := cli.BankQuery().Balance(cli.ctx, &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})
	if err != nil {
		return sdk.Coin{}, err
	}
	return *response.Balance, nil
}

func (cli *Client) QueryBalances(address string) (sdk.Coins, error) {
	response, err := cli.BankQuery().AllBalances(cli.ctx, &banktypes.QueryAllBalancesRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	return response.Balances, nil
}

func (cli *Client) QuerySupply() (sdk.Coins, error) {
	response, err := cli.BankQuery().TotalSupply(cli.ctx, &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		return nil, err
	}
	return response.Supply, nil
}

func (cli *Client) GetMintDenom() (string, error) {
	response, err := cli.StakingQuery().Params(cli.ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}
	return response.Params.BondDenom, nil
}

func (cli *Client) GetBlockHeight() (int64, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	return response.Block.Header.Height, nil
}

func (cli *Client) GetChainId() (string, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return "", err
	}
	return response.Block.Header.ChainID, nil
}

func (cli *Client) GetBlockTimeInterval() (time.Duration, error) {
	tmClient := cli.TMServiceClient()
	response1, err := tmClient.GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	if response1.Block.Header.Height <= 1 {
		return 0, fmt.Errorf("please try again later, the current block height is less than 1")
	}
	response2, err := tmClient.GetBlockByHeight(cli.ctx, &tmservice.GetBlockByHeightRequest{
		Height: response1.Block.Header.Height - 1,
	})
	if err != nil {
		return 0, err
	}
	return response1.Block.Header.Time.Sub(response2.Block.Header.Time), nil
}

func (cli *Client) GetLatestBlock() (*tenderminttypes.Block, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return response.Block, nil
}

func (cli *Client) GetBlockByHeight(blockHeight int64) (*tenderminttypes.Block, error) {
	response, err := cli.TMServiceClient().GetBlockByHeight(cli.ctx, &tmservice.GetBlockByHeightRequest{
		Height: blockHeight,
	})
	if err != nil {
		return nil, err
	}
	return response.Block, nil
}

func (cli *Client) GetStatusByTx(txHash string) (*tx.GetTxResponse, error) {
	response, err := cli.ServiceClient().GetTx(cli.ctx, &tx.GetTxRequest{
		Hash: txHash,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Deprecated: GetGasPrices
func (cli *Client) GetGasPrices() (sdk.Coins, error) {
	response, err := otherTypes.NewQueryClient(cli).GasPrice(cli.ctx, &otherTypes.GasPriceRequest{})
	if err != nil {
		return nil, err
	}
	return response.GasPrices, nil
}

func (cli *Client) GetAddressPrefix() (string, error) {
	response, err := cli.TMServiceClient().GetValidatorSetByHeight(cli.ctx, &tmservice.GetValidatorSetByHeightRequest{Height: 1})
	if err != nil {
		return "", err
	}
	if len(response.Validators) <= 0 {
		return "", errors.New("no found validator")
	}
	prefix, _, err := bech32.DecodeAndConvert(response.Validators[0].Address)
	if err != nil {
		return "", err
	}
	valConPrefix := sdk.PrefixValidator + sdk.PrefixConsensus
	if strings.HasSuffix(prefix, valConPrefix) {
		return prefix[:len(prefix)-len(valConPrefix)], nil
	}
	return "", errors.New("no found address prefix")
}

func (cli *Client) EstimatingGas(raw *tx.TxRaw) (*sdk.GasInfo, error) {
	txBytes, err := proto.Marshal(raw)
	if err != nil {
		return nil, err
	}
	response, err := cli.ServiceClient().Simulate(cli.ctx, &tx.SimulateRequest{TxBytes: txBytes})
	if err != nil {
		return nil, err
	}
	return response.GasInfo, nil
}

func (cli *Client) BuildTx(privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error) {
	account, err := cli.QueryAccount(sdk.AccAddress(privKey.PubKey().Address()).String())
	if err != nil {
		return nil, err
	}
	if len(cli.chainId) <= 0 {
		chainId, err := cli.GetChainId()
		if err != nil {
			return nil, err
		}
		cli.chainId = chainId
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

	gasPrice := sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt())
	if len(cli.gasPrices) <= 0 {
		gasPrices, err := cli.GetGasPrices()
		if err != nil {
			return nil, err
		}
		if len(gasPrices) > 0 {
			gasPrice = gasPrices[0]
		}
	} else {
		gasPrice = cli.gasPrices[0]
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
		ChainId:       cli.chainId,
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

func (cli *Client) BroadcastTxOk(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error) {
	broadcastTx, err := cli.BroadcastTx(txRaw, mode...)
	if err != nil {
		return nil, err
	}
	if broadcastTx.Code != 0 {
		return nil, errors.New(broadcastTx.RawLog)
	}
	return broadcastTx, nil
}

func (cli *Client) BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	_, err = proto.Marshal(&tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	broadcastTxResponse, err := cli.ServiceClient().BroadcastTx(cli.ctx, &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	return broadcastTxResponse.TxResponse, nil
}

func (cli *Client) BroadcastTxBytes(txBytes []byte, mode ...tx.BroadcastMode) (*sdk.TxResponse, error) {
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	if len(mode) > 0 {
		defaultMode = mode[0]
	}
	_, err := proto.Marshal(&tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	broadcastTxResponse, err := cli.ServiceClient().BroadcastTx(cli.ctx, &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	return broadcastTxResponse.TxResponse, nil
}

func (cli *Client) TxByHash(txHash string) (*sdk.TxResponse, error) {
	resp, err := cli.ServiceClient().GetTx(cli.ctx, &tx.GetTxRequest{Hash: txHash})
	if err != nil {
		return nil, err
	}
	return resp.TxResponse, nil
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

// BuildTxV2 nolint
func BuildTxV2(chainId string, sequence, accountNumber uint64, privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasPrice sdk.Coin) (*tx.TxRaw, error) {
	return BuildTxV1(chainId, sequence, accountNumber, privKey, msgs, gasPrice, DefGasLimit, "", 0)
}

// BuildTxV3 nolint
func BuildTxV3(chainId string, sequence, accountNumber uint64, privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasPrice sdk.Coin, gasLimit int64) (*tx.TxRaw, error) {
	return BuildTxV1(chainId, sequence, accountNumber, privKey, msgs, gasPrice, gasLimit, "", 0)
}
