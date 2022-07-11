package jsonrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec/legacy"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

const DefGasLimit int64 = 200000

type CustomRPC interface {
	TendermintRPC
	GetChainId() (chain string, err error)
	GetBlockHeight() (int64, error)
	GetMintDenom() (denom string, err error)
	GetGasPrices() (sdk.DecCoins, error)
	GetAddressPrefix() (prefix string, err error)
	QueryAccount(address string) (authtypes.AccountI, error)
	QueryBalance(address string, denom string) (sdk.Coin, error)
	QueryAllBalances(address string) (sdk.Coins, error)
	BuildTx(privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error)
	EstimatingGas(txBody *tx.TxBody, authInfo *tx.AuthInfo, sign []byte) (*sdk.GasInfo, error)
	BroadcastTx(txRaw *tx.TxRaw, mode ...string) (*coreTypes.ResultBroadcastTx, error)
	BroadcastTxRawCommit(txRaw *tx.TxRaw) (*coreTypes.ResultBroadcastTxCommit, error)
	TxByHash(txHash string) (*coreTypes.ResultTx, error)
	ABCIQueryIsOk(path string, data tmBytes.HexBytes) (*coreTypes.ResultABCIQuery, error)
}

type TendermintRPC interface {
	Status() (*coreTypes.ResultStatus, error)
	ABCIInfo() (*coreTypes.ResultABCIInfo, error)
	ABCIQuery(path string, data tmBytes.HexBytes) (*coreTypes.ResultABCIQuery, error)

	BroadcastTxCommit(tx types.Tx) (*coreTypes.ResultBroadcastTxCommit, error)
	BroadcastTxAsync(tx types.Tx) (*coreTypes.ResultBroadcastTx, error)
	BroadcastTxSync(tx types.Tx) (*coreTypes.ResultBroadcastTx, error)

	UnconfirmedTxs(limit int) (*coreTypes.ResultUnconfirmedTxs, error)
	NumUnconfirmedTxs() (*coreTypes.ResultUnconfirmedTxs, error)
	NetInfo() (*coreTypes.ResultNetInfo, error)
	DumpConsensusState() (*coreTypes.ResultDumpConsensusState, error)
	ConsensusState() (*coreTypes.ResultConsensusState, error)
	ConsensusParams(height int64) (*coreTypes.ResultConsensusParams, error)
	Health() (*coreTypes.ResultHealth, error)

	Genesis() (*coreTypes.ResultGenesis, error)
	BlockchainInfo(minHeight, maxHeight int64) (*coreTypes.ResultBlockchainInfo, error)

	Block(height int64) (*coreTypes.ResultBlock, error)
	BlockResults(height int64) (*coreTypes.ResultBlockResults, error)
	Commit(height int64) (*coreTypes.ResultCommit, error)
	Validators(height int64, page, perPage int) (*coreTypes.ResultValidators, error)
	Tx(hash []byte) (*coreTypes.ResultTx, error)
	TxSearch(query string, page, perPage int, orderBy string) (*coreTypes.ResultTxSearch, error)

	BroadcastEvidence(ev types.Evidence) (*coreTypes.ResultBroadcastEvidence, error)
}

type JSONRPCCaller interface {
	Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (err error)
}

type jsonRPCClient struct {
	chainId string
	caller  JSONRPCCaller
}

func NewCustomRPC(caller JSONRPCCaller) CustomRPC {
	return &jsonRPCClient{caller: caller}
}

// Custom API

func (c *jsonRPCClient) GetChainId() (chain string, err error) {
	res, err := c.Genesis()
	if err != nil {
		return
	}
	return res.Genesis.ChainID, nil
}

func (c *jsonRPCClient) GetBlockHeight() (int64, error) {
	status, err := c.Status()
	if err != nil {
		return 0, err
	}
	if status.SyncInfo.CatchingUp {
		return 0, errors.New("the node is catching up with the new block data")
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (c *jsonRPCClient) GetMintDenom() (denom string, err error) {
	genesis, err := c.Genesis()
	if err != nil {
		return
	}

	var appState map[string]json.RawMessage
	if err = json.Unmarshal(genesis.Genesis.AppState, &appState); err != nil {
		return denom, err
	}

	var mintGenesis struct {
		Params struct {
			MintDenom string `json:"mint_denom"`
		} `json:"params"`
	}
	return mintGenesis.Params.MintDenom, json.Unmarshal(appState[minttypes.ModuleName], &mintGenesis)
}

func (c *jsonRPCClient) GetAddressPrefix() (prefix string, err error) {
	genesis, err := c.Genesis()
	if err != nil {
		return
	}
	var appState map[string]json.RawMessage
	if err = json.Unmarshal(genesis.Genesis.AppState, &appState); err != nil {
		return
	}

	var authGen struct {
		Accounts []struct {
			Value struct {
				Address string `json:"address"`
			} `json:"value"`
		} `json:"accounts"`
	}
	if err = json.Unmarshal(appState[authtypes.ModuleName], &authGen); err != nil {
		return
	}
	for _, account := range authGen.Accounts {
		prefix, _, err := bech32.Decode(account.Value.Address)
		return prefix, err
	}
	return sdk.Bech32MainPrefix, nil
}

func (c *jsonRPCClient) GetGasPrices() (sdk.DecCoins, error) {
	result, err := c.ABCIQueryIsOk("/custom/other/gasPrice", nil)
	if err != nil {
		return sdk.DecCoins{}, err
	}
	var gasPrice sdk.DecCoins
	if err = json.Unmarshal(result.Response.Value, &gasPrice); err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (c *jsonRPCClient) QueryAccount(address string) (authtypes.AccountI, error) {
	result, err := c.ABCIQueryIsOk("/custom/auth/account", legacy.Cdc.MustMarshalJSON(authtypes.QueryAccountRequest{Address: address}))
	if err != nil {
		return nil, err
	}
	var account authtypes.AccountI
	if err = legacy.Cdc.UnmarshalJSON(result.Response.Value, &account); err != nil {
		return nil, err
	}
	return account, nil
}

func (c *jsonRPCClient) QueryBalance(address string, denom string) (sdk.Coin, error) {
	result, err := c.ABCIQueryIsOk("/custom/bank/balance", legacy.Cdc.MustMarshalJSON(banktypes.QueryBalanceRequest{Address: address, Denom: denom}))
	if err != nil {
		return sdk.Coin{}, err
	}
	var coin sdk.Coin
	if err = legacy.Cdc.UnmarshalJSON(result.Response.Value, &coin); err != nil {
		return sdk.Coin{}, err
	}
	return coin, nil
}

func (c *jsonRPCClient) QueryAllBalances(address string) (sdk.Coins, error) {
	result, err := c.ABCIQueryIsOk("/custom/bank/all_balances", legacy.Cdc.MustMarshalJSON(banktypes.QueryAllBalancesRequest{Address: address}))
	if err != nil {
		return nil, err
	}
	var coins sdk.Coins
	if err = legacy.Cdc.UnmarshalJSON(result.Response.Value, &coins); err != nil {
		return nil, err
	}
	return coins, nil
}

func (c *jsonRPCClient) BuildTx(privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error) {
	account, err := c.QueryAccount(sdk.AccAddress(privKey.PubKey().Address().Bytes()).String())
	if err != nil {
		return nil, err
	}
	if len(c.chainId) <= 0 {
		chainId, err := c.GetChainId()
		if err != nil {
			return nil, err
		}
		c.chainId = chainId
	}

	txBodyMessage := make([]*codecTypes.Any, 0)
	for i := 0; i < len(msgs); i++ {
		msgAnyValue, err := codecTypes.NewAnyWithValue(msgs[i])
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

	pubAny, err := codecTypes.NewAnyWithValue(privKey.PubKey())
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
				Sequence: account.GetSequence(),
			},
		},
		Fee: &tx.Fee{
			Amount:   nil,
			GasLimit: uint64(DefGasLimit),
			Payer:    "",
			Granter:  "",
		},
	}

	prices, err := c.GetGasPrices()
	if err != nil {
		return nil, err
	}

	for _, price := range prices {
		authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(price.Denom, price.Amount.MulInt64(int64(authInfo.Fee.GasLimit)).RoundInt()))
		continue
	}

	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	signDoc := &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       c.chainId,
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
	gasInfo, err := c.EstimatingGas(txBody, authInfo, sign)
	if err != nil {
		return nil, err
	}

	authInfo.Fee.GasLimit = gasInfo.GasUsed * 12 / 10
	for _, price := range prices {
		authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(price.Denom, price.Amount.MulInt64(int64(authInfo.Fee.GasLimit)).RoundInt()))
		continue
	}

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

func (c *jsonRPCClient) BroadcastTx(txRaw *tx.TxRaw, mode ...string) (*coreTypes.ResultBroadcastTx, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	defaultMode := flags.BroadcastSync
	if len(mode) > 0 {
		defaultMode = mode[0]
	}
	switch defaultMode {
	case flags.BroadcastSync:
		return c.BroadcastTxSync(txBytes)
	case flags.BroadcastAsync:
		return c.BroadcastTxAsync(txBytes)
	case flags.BroadcastBlock:
		commit, err := c.BroadcastTxCommit(txBytes)
		if err != nil {
			return nil, err
		}
		return &coreTypes.ResultBroadcastTx{
			Code:      commit.DeliverTx.GetCode(),
			Data:      commit.DeliverTx.GetData(),
			Log:       commit.DeliverTx.Log,
			Codespace: commit.DeliverTx.Codespace,
			Hash:      commit.Hash,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported return type %s; supported types: sync, async, block", defaultMode)
	}
}

func (c *jsonRPCClient) BroadcastTxRawCommit(txRaw *tx.TxRaw) (*coreTypes.ResultBroadcastTxCommit, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	return c.BroadcastTxCommit(txBytes)
}

func (c *jsonRPCClient) TxByHash(txHash string) (*coreTypes.ResultTx, error) {
	hash, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}
	return c.Tx(hash)
}

func (c *jsonRPCClient) EstimatingGas(txBody *tx.TxBody, authInfo *tx.AuthInfo, sign []byte) (*sdk.GasInfo, error) {
	result, err := c.ABCIQueryIsOk("/app/simulate", nil)
	if err != nil {
		return nil, err
	}
	var resp = sdk.SimulationResponse{}
	return &resp.GasInfo, json.Unmarshal(result.Response.Value, &resp)
}

func (c *jsonRPCClient) AppVersion() (string, error) {
	result, err := c.ABCIQueryIsOk("/app/version", nil)
	if err != nil {
		return "", err
	}
	return string(result.Response.Value), nil
}

func (c *jsonRPCClient) Store(path string) (*coreTypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/store/"+path, nil)
}

func (c *jsonRPCClient) PeersByAddressPort(port string) (*coreTypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/addr/"+port, nil)
}

func (c *jsonRPCClient) PeersById(id string) (*coreTypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/id/"+id, nil)
}

func (c *jsonRPCClient) ABCIQueryIsOk(path string, data tmBytes.HexBytes) (*coreTypes.ResultABCIQuery, error) {
	result := new(coreTypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": "0", "prove": false}
	err := c.caller.Call(context.Background(), "abci_query", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQueryIsOk")
	}
	if result.Response.Code != 0 {
		return nil, fmt.Errorf("abci query response, space: %s, code: %d, log: %s",
			result.Response.Codespace, result.Response.Code, result.Response.Log)
	}
	return result, nil
}

// Tendermint API

func (c *jsonRPCClient) Status() (*coreTypes.ResultStatus, error) {
	result := new(coreTypes.ResultStatus)
	err := c.caller.Call(context.Background(), "status", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Status")
	}
	return result, nil
}

func (c *jsonRPCClient) ABCIInfo() (*coreTypes.ResultABCIInfo, error) {
	result := new(coreTypes.ResultABCIInfo)
	err := c.caller.Call(context.Background(), "abci_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIInfo")
	}
	return result, nil
}

func (c *jsonRPCClient) ABCIQuery(path string, data tmBytes.HexBytes) (*coreTypes.ResultABCIQuery, error) {
	result := new(coreTypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": "0", "prove": false}
	err := c.caller.Call(context.Background(), "abci_query", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQuery")
	}
	return result, nil
}

func (c *jsonRPCClient) BroadcastTxCommit(tx types.Tx) (*coreTypes.ResultBroadcastTxCommit, error) {
	result := new(coreTypes.ResultBroadcastTxCommit)
	err := c.caller.Call(context.Background(), "broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, "broadcast_tx_commit")
	}
	return result, nil
}

func (c *jsonRPCClient) BroadcastTxAsync(tx types.Tx) (*coreTypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_async", tx)
}

func (c *jsonRPCClient) BroadcastTxSync(tx types.Tx) (*coreTypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_sync", tx)
}

func (c *jsonRPCClient) broadcastTX(route string, tx types.Tx) (*coreTypes.ResultBroadcastTx, error) {
	result := new(coreTypes.ResultBroadcastTx)
	err := c.caller.Call(context.Background(), route, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, route)
	}
	return result, nil
}

func (c *jsonRPCClient) UnconfirmedTxs(limit int) (*coreTypes.ResultUnconfirmedTxs, error) {
	result := new(coreTypes.ResultUnconfirmedTxs)
	err := c.caller.Call(context.Background(), "unconfirmed_txs", map[string]interface{}{"limit": limit}, result)
	if err != nil {
		return nil, errors.Wrap(err, "unconfirmed_txs")
	}
	return result, nil
}

func (c *jsonRPCClient) NumUnconfirmedTxs() (*coreTypes.ResultUnconfirmedTxs, error) {
	result := new(coreTypes.ResultUnconfirmedTxs)
	err := c.caller.Call(context.Background(), "num_unconfirmed_txs", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "num_unconfirmed_txs")
	}
	return result, nil
}

func (c *jsonRPCClient) NetInfo() (*coreTypes.ResultNetInfo, error) {
	result := new(coreTypes.ResultNetInfo)
	err := c.caller.Call(context.Background(), "net_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "NetInfo")
	}
	return result, nil
}

func (c *jsonRPCClient) DumpConsensusState() (*coreTypes.ResultDumpConsensusState, error) {
	result := new(coreTypes.ResultDumpConsensusState)
	err := c.caller.Call(context.Background(), "dump_consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "DumpConsensusState")
	}
	return result, nil
}

func (c *jsonRPCClient) ConsensusState() (*coreTypes.ResultConsensusState, error) {
	result := new(coreTypes.ResultConsensusState)
	err := c.caller.Call(context.Background(), "consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "ConsensusState")
	}
	return result, nil
}

func (c *jsonRPCClient) ConsensusParams(height int64) (*coreTypes.ResultConsensusParams, error) {
	result := new(coreTypes.ResultConsensusParams)
	params := map[string]interface{}{"height": height}
	if height <= 0 {
		params = map[string]interface{}{}
	}
	err := c.caller.Call(context.Background(), "consensus_params", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ConsensusParams")
	}
	return result, nil
}

func (c *jsonRPCClient) Health() (*coreTypes.ResultHealth, error) {
	result := new(coreTypes.ResultHealth)
	err := c.caller.Call(context.Background(), "health", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Health")
	}
	return result, nil
}

func (c *jsonRPCClient) BlockchainInfo(minHeight, maxHeight int64) (*coreTypes.ResultBlockchainInfo, error) {
	result := new(coreTypes.ResultBlockchainInfo)
	params := map[string]interface{}{"minHeight": minHeight, "maxHeight": maxHeight}
	err := c.caller.Call(context.Background(), "blockchain", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "BlockchainInfo")
	}
	return result, nil
}

func (c *jsonRPCClient) Genesis() (*coreTypes.ResultGenesis, error) {
	result := new(coreTypes.ResultGenesis)
	err := c.caller.Call(context.Background(), "genesis", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Genesis")
	}
	return result, nil
}

func (c *jsonRPCClient) Block(height int64) (*coreTypes.ResultBlock, error) {
	result := new(coreTypes.ResultBlock)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": strconv.FormatInt(height, 10)}
	}
	err := c.caller.Call(context.Background(), "block", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Block")
	}
	return result, nil
}

func (c *jsonRPCClient) BlockResults(height int64) (*coreTypes.ResultBlockResults, error) {
	result := new(coreTypes.ResultBlockResults)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": strconv.FormatInt(height, 10)}
	}
	err := c.caller.Call(context.Background(), "block_results", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Block Result")
	}
	return result, nil
}

func (c *jsonRPCClient) Commit(height int64) (*coreTypes.ResultCommit, error) {
	result := new(coreTypes.ResultCommit)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": height}
	}
	err := c.caller.Call(context.Background(), "commit", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Commit")
	}
	return result, nil
}

func (c *jsonRPCClient) Validators(height int64, page, perPage int) (*coreTypes.ResultValidators, error) {
	result := new(coreTypes.ResultValidators)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": height, "page": page, "per_page": perPage}
	}
	err := c.caller.Call(context.Background(), "validators", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Validators")
	}
	return result, nil
}

func (c *jsonRPCClient) Tx(hash []byte) (*coreTypes.ResultTx, error) {
	result := new(coreTypes.ResultTx)
	params := map[string]interface{}{"hash": hash, "prove": false}
	err := c.caller.Call(context.Background(), "tx", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Tx")
	}
	return result, nil
}

func (c *jsonRPCClient) TxSearch(query string, page, perPage int, orderBy string) (
	*coreTypes.ResultTxSearch, error) {
	result := new(coreTypes.ResultTxSearch)
	params := map[string]interface{}{"query": query, "prove": false, "page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "order_by": orderBy}
	err := c.caller.Call(context.Background(), "tx_search", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "TxSearch")
	}
	return result, nil
}

func (c *jsonRPCClient) BroadcastEvidence(ev types.Evidence) (*coreTypes.ResultBroadcastEvidence, error) {
	result := new(coreTypes.ResultBroadcastEvidence)
	err := c.caller.Call(context.Background(), "broadcast_evidence", map[string]interface{}{"evidence": ev}, result)
	if err != nil {
		return nil, errors.Wrap(err, "BroadcastEvidence")
	}
	return result, nil
}
