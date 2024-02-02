package jsonrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcutil/bech32"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v6/client"
)

type jsonRPCCaller interface {
	Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (err error)
}

type NodeRPC struct {
	chainId    string
	addrPrefix string
	gasPrices  sdk.Coins
	height     int64
	ctx        context.Context
	caller     jsonRPCCaller
}

func NewNodeRPC(caller jsonRPCCaller, ctx ...context.Context) *NodeRPC {
	c := &NodeRPC{caller: caller, height: 0}
	if len(ctx) > 0 {
		c.ctx = ctx[0]
	} else {
		c.ctx = context.Background()
	}
	return c
}

func (c *NodeRPC) WithContext(ctx context.Context) *NodeRPC {
	return &NodeRPC{chainId: c.chainId, gasPrices: c.gasPrices, height: c.height, ctx: ctx, caller: c.caller}
}

func (c *NodeRPC) WithGasPrices(gasPrices sdk.Coins) *NodeRPC {
	return &NodeRPC{chainId: c.chainId, gasPrices: gasPrices, height: c.height, ctx: c.ctx, caller: c.caller}
}

func (c *NodeRPC) WithBlockHeight(height int64) *NodeRPC {
	return &NodeRPC{chainId: c.chainId, gasPrices: c.gasPrices, height: height, ctx: c.ctx, caller: c.caller}
}

func (c *NodeRPC) WithChainId(chainId string) *NodeRPC {
	return &NodeRPC{chainId: chainId, gasPrices: c.gasPrices, height: c.height, ctx: c.ctx, caller: c.caller}
}

// Custom API

func (c *NodeRPC) GetChainId() (chain string, err error) {
	if len(c.chainId) > 0 {
		return c.chainId, nil
	}
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
	if len(c.addrPrefix) > 0 {
		return c.addrPrefix, nil
	}
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
	if len(authGen.Accounts) == 0 {
		return sdk.Bech32MainPrefix, nil
	}
	c.addrPrefix, _, err = bech32.Decode(authGen.Accounts[0].Address)
	return c.addrPrefix, err
}

func (c *NodeRPC) GetStakeValidators(status stakingtypes.BondStatus) (stakingtypes.Validators, error) {
	data, err := json.Marshal(map[string]string{"Page": "1", "Limit": "200", "Status": status.String()})
	if err != nil {
		return nil, err
	}
	result, err := c.ABCIQueryIsOk("/custom/staking/validators", data)
	if err != nil {
		return nil, err
	}
	validators := make(stakingtypes.Validators, 0)
	if err := json.Unmarshal(result.Response.Value, &validators); err != nil {
		return nil, err
	}
	return validators, err
}

func (c *NodeRPC) GetValAddressByCons(consAddrStr string) (sdk.ValAddress, error) {
	consAddr, err := sdk.ConsAddressFromBech32(consAddrStr)
	if err != nil {
		consAddr, err = hex.DecodeString(consAddrStr)
		if err != nil {
			return nil, errors.New("expected hex or bech32 address")
		}
	}
	result, err := c.ABCIQueryIsOk("/store/staking/key", stakingtypes.GetValidatorByConsAddrKey(consAddr))
	if err != nil {
		return nil, err
	}
	if result.Response.Value == nil {
		return nil, fmt.Errorf("not found validator by consAddress: %s", consAddr.String())
	}
	return result.Response.Value, nil
}

func (c *NodeRPC) BuildTx(privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error) {
	return client.BuildTx(c, privKey, msgs)
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

func (c *NodeRPC) BroadcastTxRawCommit(txRaw *tx.TxRaw) (*ctypes.ResultBroadcastTxCommit, error) {
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

func (c *NodeRPC) EstimatingGas(raw *tx.TxRaw) (*sdk.GasInfo, error) {
	txBytes, err := proto.Marshal(raw)
	if err != nil {
		return nil, err
	}
	result, err := c.ABCIQueryIsOk("/app/simulate", txBytes)
	if err != nil {
		return nil, err
	}
	resp := struct {
		GasInfo struct {
			GasWanted string `json:"gas_wanted"`
			GasUsed   string `json:"gas_used"`
		} `json:"gas_info"`
	}{}
	if err = json.Unmarshal(result.Response.Value, &resp); err != nil {
		return nil, err
	}
	gasWanted, err := strconv.ParseUint(resp.GasInfo.GasWanted, 10, 64)
	if err != nil && len(resp.GasInfo.GasWanted) > 0 {
		return nil, err
	}
	gasUsed, err := strconv.ParseUint(resp.GasInfo.GasUsed, 10, 64)
	if err != nil && len(resp.GasInfo.GasUsed) > 0 {
		return nil, err
	}
	return &sdk.GasInfo{
		GasWanted: gasWanted,
		GasUsed:   gasUsed,
	}, nil
}

func (c *NodeRPC) AppVersion() (string, error) {
	result, err := c.ABCIQueryIsOk("/app/version", nil)
	if err != nil {
		return "", err
	}
	return string(result.Response.Value), nil
}

func (c *NodeRPC) ABCIQueryIsOk(path string, data tmbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": strconv.FormatInt(c.height, 10), "prove": false}
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

func (c *NodeRPC) Status() (*ctypes.ResultStatus, error) {
	result := new(ctypes.ResultStatus)
	err := c.caller.Call(c.ctx, "status", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Status")
	}
	return result, nil
}

func (c *NodeRPC) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	result := new(ctypes.ResultABCIInfo)
	err := c.caller.Call(c.ctx, "abci_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIInfo")
	}
	return result, nil
}

func (c *NodeRPC) ABCIQuery(path string, data tmbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": strconv.FormatInt(c.height, 10), "prove": false}
	err := c.caller.Call(c.ctx, "abci_query", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQuery")
	}
	return result, nil
}

func (c *NodeRPC) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	result := new(ctypes.ResultBroadcastTxCommit)
	err := c.caller.Call(c.ctx, "broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, "broadcast_tx_commit")
	}
	return result, nil
}

func (c *NodeRPC) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_async", tx)
}

func (c *NodeRPC) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_sync", tx)
}

func (c *NodeRPC) broadcastTX(route string, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	result := new(ctypes.ResultBroadcastTx)
	err := c.caller.Call(c.ctx, route, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, route)
	}
	return result, nil
}

func (c *NodeRPC) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	result := new(ctypes.ResultUnconfirmedTxs)
	params := map[string]interface{}{"limit": strconv.Itoa(limit)}
	err := c.caller.Call(c.ctx, "unconfirmed_txs", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "unconfirmed_txs")
	}
	if len(result.Txs) == 0 {
		result.Txs = make([]types.Tx, 0)
	}
	return result, nil
}

func (c *NodeRPC) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	result := new(ctypes.ResultUnconfirmedTxs)
	err := c.caller.Call(c.ctx, "num_unconfirmed_txs", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "num_unconfirmed_txs")
	}
	return result, nil
}

func (c *NodeRPC) NetInfo() (*ctypes.ResultNetInfo, error) {
	result := new(ctypes.ResultNetInfo)
	err := c.caller.Call(c.ctx, "net_info", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "NetInfo")
	}
	if len(result.Peers) == 0 {
		result.Peers = make([]ctypes.Peer, 0)
	}
	return result, nil
}

func (c *NodeRPC) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	result := new(ctypes.ResultDumpConsensusState)
	err := c.caller.Call(c.ctx, "dump_consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "DumpConsensusState")
	}
	if len(result.Peers) == 0 {
		result.Peers = make([]ctypes.PeerStateInfo, 0)
	}
	return result, nil
}

func (c *NodeRPC) ConsensusState() (*ctypes.ResultConsensusState, error) {
	result := new(ctypes.ResultConsensusState)
	err := c.caller.Call(c.ctx, "consensus_state", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "ConsensusState")
	}
	return result, nil
}

func (c *NodeRPC) ConsensusParams(height int64) (*ctypes.ResultConsensusParams, error) {
	result := new(ctypes.ResultConsensusParams)
	params := map[string]interface{}{"height": strconv.FormatInt(height, 10)}
	if height <= 0 {
		params = map[string]interface{}{}
	}
	err := c.caller.Call(c.ctx, "consensus_params", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "ConsensusParams")
	}
	return result, nil
}

func (c *NodeRPC) Health() (*ctypes.ResultHealth, error) {
	result := new(ctypes.ResultHealth)
	err := c.caller.Call(c.ctx, "health", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Health")
	}
	return result, nil
}

func (c *NodeRPC) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	result := new(ctypes.ResultBlockchainInfo)
	params := map[string]interface{}{
		"minHeight": strconv.FormatInt(minHeight, 10),
		"maxHeight": strconv.FormatInt(maxHeight, 10),
	}
	err := c.caller.Call(c.ctx, "blockchain", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "BlockchainInfo")
	}
	return result, nil
}

func (c *NodeRPC) Genesis() (*ctypes.ResultGenesis, error) {
	result := new(ctypes.ResultGenesis)
	err := c.caller.Call(c.ctx, "genesis", map[string]interface{}{}, result)
	if err != nil {
		return nil, errors.Wrap(err, "Genesis")
	}
	return result, nil
}

func (c *NodeRPC) Block(height int64) (*ctypes.ResultBlock, error) {
	result := new(ctypes.ResultBlock)
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

func (c *NodeRPC) BlockResults(height int64) (*ctypes.ResultBlockResults, error) {
	result := new(ctypes.ResultBlockResults)
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

func (c *NodeRPC) Commit(height int64) (*ctypes.ResultCommit, error) {
	result := new(ctypes.ResultCommit)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": strconv.FormatInt(height, 10)}
	}
	err := c.caller.Call(c.ctx, "commit", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Commit")
	}
	return result, nil
}

func (c *NodeRPC) Validators(height int64, page, perPage int) (*ctypes.ResultValidators, error) {
	result := new(ctypes.ResultValidators)
	params := map[string]interface{}{}
	if height > 0 {
		params = map[string]interface{}{"height": strconv.FormatInt(height, 10), "page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage)}
	}
	err := c.caller.Call(c.ctx, "validators", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Validators")
	}
	return result, nil
}

func (c *NodeRPC) Tx(hash []byte) (*ctypes.ResultTx, error) {
	result := new(ctypes.ResultTx)
	params := map[string]interface{}{"hash": hash, "prove": false}
	err := c.caller.Call(c.ctx, "tx", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "Tx")
	}
	return result, nil
}

func (c *NodeRPC) TxSearch(query string, page, perPage int, orderBy string) (
	*ctypes.ResultTxSearch, error,
) {
	result := new(ctypes.ResultTxSearch)
	params := map[string]interface{}{"query": query, "prove": false, "page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "order_by": orderBy}
	err := c.caller.Call(c.ctx, "tx_search", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "TxSearch")
	}
	return result, nil
}

func (c *NodeRPC) BlockSearch(query string, page, perPage int, orderBy string) (*ctypes.ResultBlockSearch, error) {
	result := new(ctypes.ResultBlockSearch)
	params := map[string]interface{}{"query": query, "prove": false, "page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "order_by": orderBy}
	err := c.caller.Call(c.ctx, "block_search", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "BlockSearch")
	}
	return result, nil
}

func (c *NodeRPC) BroadcastEvidence(ev types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	result := new(ctypes.ResultBroadcastEvidence)
	err := c.caller.Call(c.ctx, "broadcast_evidence", map[string]interface{}{"evidence": ev}, result)
	if err != nil {
		return nil, errors.Wrap(err, "BroadcastEvidence")
	}
	return result, nil
}
