package jsonrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/btcsuite/btcutil/bech32"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

const DefGasLimit int64 = 200000

type jsonRPCCaller interface {
	Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (err error)
}

type NodeRPC struct {
	chainId string
	caller  jsonRPCCaller
	ctx     context.Context
}

func NewNodeRPC(caller jsonRPCCaller) *NodeRPC {
	return &NodeRPC{caller: caller, ctx: context.Background()}
}

func (c *NodeRPC) WithContext(ctx context.Context) *NodeRPC {
	c.ctx = ctx
	return c
}

// Custom API

func (c *NodeRPC) GetChainId() (chain string, err error) {
	res, err := c.Genesis()
	if err != nil {
		return
	}
	return res.Genesis.ChainID, nil
}

func (c *NodeRPC) GetBlockHeight() (int64, error) {
	status, err := c.Status()
	if err != nil {
		return 0, err
	}
	if status.SyncInfo.CatchingUp {
		return 0, errors.New("the node is catching up with the new block data")
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (c *NodeRPC) GetMintDenom() (denom string, err error) {
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

func (c *NodeRPC) GetAddressPrefix() (prefix string, err error) {
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
			Address string `json:"address"`
		} `json:"accounts"`
	}
	if err = json.Unmarshal(appState[authtypes.ModuleName], &authGen); err != nil {
		return
	}
	for _, account := range authGen.Accounts {
		prefix, _, err := bech32.Decode(account.Address)
		return prefix, err
	}
	return sdk.Bech32MainPrefix, nil
}

func (c *NodeRPC) BuildTx(privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error) {
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
		authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(price.Denom, price.Amount.MulRaw(int64(authInfo.Fee.GasLimit))))
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
		authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(price.Denom, price.Amount.MulRaw(int64(authInfo.Fee.GasLimit))))
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

func (c *NodeRPC) BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	if len(mode) > 0 {
		defaultMode = mode[0]
	}
	switch defaultMode {
	case tx.BroadcastMode_BROADCAST_MODE_SYNC:
		res, err := c.BroadcastTxSync(txBytes)
		if err != nil {
			return nil, err
		}
		return sdk.NewResponseFormatBroadcastTx(res), nil
	case tx.BroadcastMode_BROADCAST_MODE_ASYNC:
		res, err := c.BroadcastTxAsync(txBytes)
		if err != nil {
			return nil, err
		}
		return sdk.NewResponseFormatBroadcastTx(res), nil
	case tx.BroadcastMode_BROADCAST_MODE_BLOCK:
		res, err := c.BroadcastTxCommit(txBytes)
		if err != nil {
			return nil, err
		}
		return sdk.NewResponseFormatBroadcastTxCommit(res), nil
	default:
		return nil, fmt.Errorf("unsupported return type %s; supported types: sync, async, block", defaultMode)
	}
}

func (c *NodeRPC) BroadcastTxRawCommit(txRaw *tx.TxRaw) (*coreTypes.ResultBroadcastTxCommit, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	return c.BroadcastTxCommit(txBytes)
}

func (c *NodeRPC) TxByHash(txHash string) (*sdk.TxResponse, error) {
	hash, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}
	resultTx, err := c.Tx(hash)
	if err != nil {
		return nil, err
	}
	return sdk.NewResponseResultTx(resultTx, nil, ""), nil
}

func (c *NodeRPC) EstimatingGas(txBody *tx.TxBody, authInfo *tx.AuthInfo, sign []byte) (*sdk.GasInfo, error) {
	result, err := c.ABCIQueryIsOk("/app/simulate", nil)
	if err != nil {
		return nil, err
	}
	var resp = sdk.SimulationResponse{}
	return &resp.GasInfo, json.Unmarshal(result.Response.Value, &resp)
}

func (c *NodeRPC) AppVersion() (string, error) {
	result, err := c.ABCIQueryIsOk("/app/version", nil)
	if err != nil {
		return "", err
	}
	return string(result.Response.Value), nil
}

func (c *NodeRPC) ABCIQueryIsOk(path string, data tmBytes.HexBytes) (*coreTypes.ResultABCIQuery, error) {
	result := new(coreTypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": "0", "prove": false}
	err := c.caller.Call(c.ctx, "abci_query", params, result)
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

func (c *NodeRPC) Status() (*coreTypes.ResultStatus, error) {
	result := new(coreTypes.ResultStatus)
	err := c.caller.Call(c.ctx, "status", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Status")
	}
	return result, nil
}

func (c *NodeRPC) ABCIInfo() (*coreTypes.ResultABCIInfo, error) {
	result := new(coreTypes.ResultABCIInfo)
	err := c.caller.Call(c.ctx, "abci_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIInfo")
	}
	return result, nil
}

func (c *NodeRPC) ABCIQuery(path string, data tmBytes.HexBytes) (*coreTypes.ResultABCIQuery, error) {
	result := new(coreTypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": "0", "prove": false}
	err := c.caller.Call(c.ctx, "abci_query", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQuery")
	}
	return result, nil
}

func (c *NodeRPC) BroadcastTxCommit(tx types.Tx) (*coreTypes.ResultBroadcastTxCommit, error) {
	result := new(coreTypes.ResultBroadcastTxCommit)
	err := c.caller.Call(c.ctx, "broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, "broadcast_tx_commit")
	}
	return result, nil
}

func (c *NodeRPC) BroadcastTxAsync(tx types.Tx) (*coreTypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_async", tx)
}

func (c *NodeRPC) BroadcastTxSync(tx types.Tx) (*coreTypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_sync", tx)
}

func (c *NodeRPC) broadcastTX(route string, tx types.Tx) (*coreTypes.ResultBroadcastTx, error) {
	result := new(coreTypes.ResultBroadcastTx)
	err := c.caller.Call(c.ctx, route, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, route)
	}
	return result, nil
}

func (c *NodeRPC) UnconfirmedTxs(limit int) (*coreTypes.ResultUnconfirmedTxs, error) {
	result := new(coreTypes.ResultUnconfirmedTxs)
	err := c.caller.Call(c.ctx, "unconfirmed_txs", map[string]interface{}{"limit": limit}, result)
	if err != nil {
		return nil, errors.Wrap(err, "unconfirmed_txs")
	}
	return result, nil
}

func (c *NodeRPC) NumUnconfirmedTxs() (*coreTypes.ResultUnconfirmedTxs, error) {
	result := new(coreTypes.ResultUnconfirmedTxs)
	err := c.caller.Call(c.ctx, "num_unconfirmed_txs", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "num_unconfirmed_txs")
	}
	return result, nil
}

func (c *NodeRPC) NetInfo() (*coreTypes.ResultNetInfo, error) {
	result := new(coreTypes.ResultNetInfo)
	err := c.caller.Call(c.ctx, "net_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "NetInfo")
	}
	return result, nil
}

func (c *NodeRPC) DumpConsensusState() (*coreTypes.ResultDumpConsensusState, error) {
	result := new(coreTypes.ResultDumpConsensusState)
	err := c.caller.Call(c.ctx, "dump_consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "DumpConsensusState")
	}
	return result, nil
}

func (c *NodeRPC) ConsensusState() (*coreTypes.ResultConsensusState, error) {
	result := new(coreTypes.ResultConsensusState)
	err := c.caller.Call(c.ctx, "consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "ConsensusState")
	}
	return result, nil
}

func (c *NodeRPC) ConsensusParams(height int64) (*coreTypes.ResultConsensusParams, error) {
	result := new(coreTypes.ResultConsensusParams)
	params := map[string]interface{}{"height": height}
	if height <= 0 {
		params = map[string]interface{}{}
	}
	err := c.caller.Call(c.ctx, "consensus_params", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ConsensusParams")
	}
	return result, nil
}

func (c *NodeRPC) Health() (*coreTypes.ResultHealth, error) {
	result := new(coreTypes.ResultHealth)
	err := c.caller.Call(c.ctx, "health", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Health")
	}
	return result, nil
}

func (c *NodeRPC) BlockchainInfo(minHeight, maxHeight int64) (*coreTypes.ResultBlockchainInfo, error) {
	result := new(coreTypes.ResultBlockchainInfo)
	params := map[string]interface{}{"minHeight": minHeight, "maxHeight": maxHeight}
	err := c.caller.Call(c.ctx, "blockchain", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "BlockchainInfo")
	}
	return result, nil
}

func (c *NodeRPC) Genesis() (*coreTypes.ResultGenesis, error) {
	result := new(coreTypes.ResultGenesis)
	err := c.caller.Call(c.ctx, "genesis", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Genesis")
	}
	return result, nil
}

func (c *NodeRPC) Block(height int64) (*coreTypes.ResultBlock, error) {
	result := new(coreTypes.ResultBlock)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": strconv.FormatInt(height, 10)}
	}
	err := c.caller.Call(c.ctx, "block", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Block")
	}
	return result, nil
}

func (c *NodeRPC) BlockResults(height int64) (*coreTypes.ResultBlockResults, error) {
	result := new(coreTypes.ResultBlockResults)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": strconv.FormatInt(height, 10)}
	}
	err := c.caller.Call(c.ctx, "block_results", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Block Result")
	}
	return result, nil
}

func (c *NodeRPC) Commit(height int64) (*coreTypes.ResultCommit, error) {
	result := new(coreTypes.ResultCommit)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": height}
	}
	err := c.caller.Call(c.ctx, "commit", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Commit")
	}
	return result, nil
}

func (c *NodeRPC) Validators(height int64, page, perPage int) (*coreTypes.ResultValidators, error) {
	result := new(coreTypes.ResultValidators)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": height, "page": page, "per_page": perPage}
	}
	err := c.caller.Call(c.ctx, "validators", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Validators")
	}
	return result, nil
}

func (c *NodeRPC) Tx(hash []byte) (*coreTypes.ResultTx, error) {
	result := new(coreTypes.ResultTx)
	params := map[string]interface{}{"hash": hash, "prove": false}
	err := c.caller.Call(c.ctx, "tx", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Tx")
	}
	return result, nil
}

func (c *NodeRPC) TxSearch(query string, page, perPage int, orderBy string) (
	*coreTypes.ResultTxSearch, error) {
	result := new(coreTypes.ResultTxSearch)
	params := map[string]interface{}{"query": query, "prove": false, "page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "order_by": orderBy}
	err := c.caller.Call(c.ctx, "tx_search", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "TxSearch")
	}
	return result, nil
}

func (c *NodeRPC) BroadcastEvidence(ev types.Evidence) (*coreTypes.ResultBroadcastEvidence, error) {
	result := new(coreTypes.ResultBroadcastEvidence)
	err := c.caller.Call(c.ctx, "broadcast_evidence", map[string]interface{}{"evidence": ev}, result)
	if err != nil {
		return nil, errors.Wrap(err, "BroadcastEvidence")
	}
	return result, nil
}
