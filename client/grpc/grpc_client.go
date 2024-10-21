package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/google"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/functionx/fx-core/v8/client"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	migratetypes "github.com/functionx/fx-core/v8/x/migrate/types"
)

type Client struct {
	chainId    string
	addrPrefix string
	gasPrices  sdk.Coins
	ctx        context.Context
	grpc1.ClientConn
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
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return grpc.NewClient(_url, opts...)
}

func NewClient(conn grpc1.ClientConn, ctx ...context.Context) *Client {
	cli := &Client{ClientConn: conn}
	if len(ctx) > 0 {
		cli.ctx = ctx[0]
	} else {
		cli.ctx = context.Background()
	}
	return cli
}

func DailClient(rawUrl string, ctx ...context.Context) (*Client, error) {
	grpcConn, err := NewGrpcConn(rawUrl)
	if err != nil {
		return nil, err
	}
	return NewClient(grpcConn, ctx...), nil
}

func (cli *Client) WithContext(ctx context.Context) *Client {
	return &Client{chainId: cli.chainId, gasPrices: cli.gasPrices, ctx: ctx, ClientConn: cli.ClientConn}
}

func (cli *Client) WithGasPrices(gasPrices sdk.Coins) *Client {
	return &Client{chainId: cli.chainId, gasPrices: gasPrices, ctx: cli.ctx, ClientConn: cli.ClientConn}
}

func (cli *Client) WithBlockHeight(height int64) *Client {
	ctx := metadata.AppendToOutgoingContext(cli.ctx, grpctypes.GRPCBlockHeightHeader, strconv.FormatInt(height, 10))
	return &Client{chainId: cli.chainId, gasPrices: cli.gasPrices, ctx: ctx, ClientConn: cli.ClientConn}
}

func (cli *Client) WithChainId(chainId string) *Client {
	return &Client{chainId: chainId, gasPrices: cli.gasPrices, ctx: cli.ctx, ClientConn: cli.ClientConn}
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

func (cli *Client) GovQuery() govv1.QueryClient {
	return govv1.NewQueryClient(cli.ClientConn)
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

func (cli *Client) TMServiceClient() cmtservice.ServiceClient {
	return cmtservice.NewServiceClient(cli.ClientConn)
}

func (cli *Client) ERC20Query() erc20types.QueryClient {
	return erc20types.NewQueryClient(cli.ClientConn)
}

func (cli *Client) EVMQuery() evmtypes.QueryClient {
	return evmtypes.NewQueryClient(cli.ClientConn)
}

func (cli *Client) CrosschainQuery() crosschaintypes.QueryClient {
	return crosschaintypes.NewQueryClient(cli.ClientConn)
}

func (cli *Client) MigrateQuery() migratetypes.QueryClient {
	return migratetypes.NewQueryClient(cli.ClientConn)
}

func (cli *Client) AppVersion() (string, error) {
	info, err := cli.TMServiceClient().GetNodeInfo(cli.ctx, &cmtservice.GetNodeInfoRequest{})
	if err != nil {
		return "", err
	}
	return info.GetApplicationVersion().GetVersion(), nil
}

func (cli *Client) QueryAccount(address string) (sdk.AccountI, error) {
	response, err := cli.AuthQuery().Account(cli.ctx, &authtypes.QueryAccountRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	var account sdk.AccountI
	if err = client.NewAccountCodec().UnpackAny(response.GetAccount(), &account); err != nil {
		return nil, err
	}
	return account, nil
}

func (cli *Client) QueryBalance(address, denom string) (sdk.Coin, error) {
	response, err := cli.BankQuery().Balance(cli.ctx, &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})
	if err != nil {
		return sdk.Coin{}, err
	}
	return *response.GetBalance(), nil
}

func (cli *Client) QueryBalances(address string) (sdk.Coins, error) {
	response, err := cli.BankQuery().AllBalances(cli.ctx, &banktypes.QueryAllBalancesRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	return response.GetBalances(), nil
}

func (cli *Client) QuerySupply() (sdk.Coins, error) {
	response, err := cli.BankQuery().TotalSupply(cli.ctx, &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetSupply(), nil
}

func (cli *Client) GetMintDenom() (string, error) {
	response, err := cli.MintQuery().Params(cli.ctx, &minttypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}
	return response.GetParams().MintDenom, nil
}

func (cli *Client) GetStakingDenom() (string, error) {
	response, err := cli.StakingQuery().Params(cli.ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}
	return response.GetParams().BondDenom, nil
}

func (cli *Client) GetBlockHeight() (int64, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	return response.GetSdkBlock().GetHeader().Height, nil
}

func (cli *Client) GetChainId() (string, error) {
	if len(cli.chainId) > 0 {
		return cli.chainId, nil
	}
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return "", err
	}
	return response.GetSdkBlock().GetHeader().ChainID, nil
}

func (cli *Client) GetBlockTimeInterval() (time.Duration, error) {
	tmClient := cli.TMServiceClient()
	response1, err := tmClient.GetLatestBlock(cli.ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	if response1.GetSdkBlock().GetHeader().Height <= 1 {
		return 0, fmt.Errorf("please try again later, the current block height is less than 1")
	}
	response2, err := tmClient.GetBlockByHeight(cli.ctx, &cmtservice.GetBlockByHeightRequest{
		Height: response1.GetSdkBlock().GetHeader().Height - 1,
	})
	if err != nil {
		return 0, err
	}
	return response1.GetSdkBlock().GetHeader().Time.Sub(response2.GetSdkBlock().GetHeader().Time), nil
}

func (cli *Client) GetLatestBlock() (*cmtservice.Block, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetSdkBlock(), nil
}

func (cli *Client) GetBlockByHeight(blockHeight int64) (*cmtservice.Block, error) {
	response, err := cli.TMServiceClient().GetBlockByHeight(cli.ctx, &cmtservice.GetBlockByHeightRequest{
		Height: blockHeight,
	})
	if err != nil {
		return nil, err
	}
	return response.GetSdkBlock(), nil
}

func (cli *Client) GetGasPrices() (sdk.Coins, error) {
	if len(cli.gasPrices) > 0 {
		return cli.gasPrices, nil
	}
	response, err := node.NewServiceClient(cli).Config(cli.ctx, &node.ConfigRequest{})
	if err != nil {
		return nil, err
	}
	coins, err := sdk.ParseCoinsNormalized(response.GetMinimumGasPrice())
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (cli *Client) GetAddressPrefix() (string, error) {
	if len(cli.addrPrefix) > 0 {
		return cli.addrPrefix, nil
	}
	response, err := cli.TMServiceClient().GetLatestValidatorSet(cli.ctx, &cmtservice.GetLatestValidatorSetRequest{})
	if err != nil {
		return "", err
	}
	if len(response.GetValidators()) == 0 {
		return "", errors.New("no found validator")
	}
	prefix, _, err := bech32.DecodeAndConvert(response.GetValidators()[0].GetAddress())
	if err != nil {
		return "", err
	}
	valConPrefix := sdk.PrefixValidator + sdk.PrefixConsensus
	if strings.HasSuffix(prefix, valConPrefix) {
		cli.addrPrefix = prefix[:len(prefix)-len(valConPrefix)]
		return cli.addrPrefix, nil
	}
	return "", errors.New("no found address prefix")
}

func (cli *Client) GetModuleAccounts() ([]sdk.AccountI, error) {
	response, err := cli.AuthQuery().ModuleAccounts(cli.ctx, &authtypes.QueryModuleAccountsRequest{})
	if err != nil {
		return nil, err
	}
	accounts := make([]sdk.AccountI, 0, len(response.Accounts))
	for _, acc := range response.Accounts {
		var account sdk.AccountI
		if err = client.NewAccountCodec().UnpackAny(acc, &account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (cli *Client) GetSyncing() (bool, error) {
	response, err := cli.TMServiceClient().GetSyncing(cli.ctx, &cmtservice.GetSyncingRequest{})
	if err != nil {
		return false, err
	}
	return response.Syncing, nil
}

func (cli *Client) GetNodeInfo() (*cmtservice.VersionInfo, error) {
	response, err := cli.TMServiceClient().GetNodeInfo(cli.ctx, &cmtservice.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetApplicationVersion(), nil
}

func (cli *Client) CurrentPlan() (*upgradetypes.Plan, error) {
	response, err := cli.UpgradeQuery().CurrentPlan(cli.ctx, &upgradetypes.QueryCurrentPlanRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetPlan(), nil
}

func (cli *Client) GetValidators() ([]stakingtypes.Validator, error) {
	validators, err := cli.StakingQuery().Validators(cli.ctx, &stakingtypes.QueryValidatorsRequest{})
	if err != nil {
		return nil, err
	}
	return validators.GetValidators(), nil
}

func (cli *Client) GetConsensusValidators() ([]*cmtservice.Validator, error) {
	response, err := cli.TMServiceClient().GetLatestValidatorSet(cli.ctx, &cmtservice.GetLatestValidatorSetRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetValidators(), nil
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
	return response.GetGasInfo(), nil
}

func (cli *Client) BuildTxRaw(privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasLimit, timeout uint64, memo string) (*tx.TxRaw, error) {
	return client.BuildTxRawWithCli(cli, privKey, msgs, gasLimit, timeout, memo)
}

func (cli *Client) BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_SYNC
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
	return broadcastTxResponse.GetTxResponse(), nil
}

func (cli *Client) BroadcastTxBytes(txBytes []byte, mode ...tx.BroadcastMode) (*sdk.TxResponse, error) {
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_SYNC
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
	return broadcastTxResponse.GetTxResponse(), nil
}

func (cli *Client) TxByHash(txHash string) (*sdk.TxResponse, error) {
	resp, err := cli.ServiceClient().GetTx(cli.ctx, &tx.GetTxRequest{Hash: txHash})
	if err != nil {
		return nil, err
	}
	return resp.GetTxResponse(), nil
}

func (cli *Client) WaitMined(txHash string, timeout, pollInterval time.Duration) (*sdk.TxResponse, error) {
	ctx, cancelFunc := context.WithTimeout(cli.ctx, timeout)
	defer cancelFunc()
	newCli := cli.WithContext(ctx)
	return client.WaitMined(newCli, txHash, timeout, pollInterval)
}
