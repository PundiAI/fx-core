package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	libTypes "github.com/tendermint/tendermint/libs/bytes"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"
	"nhooyr.io/websocket"
)

var _ jsonRPCCaller = &WSClient{}

// WSClient implement jsonRPCCaller
type WSClient struct {
	ctx          context.Context
	conn         *websocket.Conn
	quit         chan struct{}
	mtx          sync.Mutex
	responseChan map[string]map[string]chan RPCResponse
	log.Logger
}

func NewWsClient(ctx context.Context, url string) (*WSClient, error) {
	split := strings.Split(url, "://")
	if len(split) > 1 && split[0] == "tcp" {
		url = fmt.Sprintf("http://%s", split[1])
	}

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(1024 * 1000 * 100)

	var ws = &WSClient{
		ctx:          ctx,
		conn:         conn,
		quit:         make(chan struct{}),
		responseChan: make(map[string]map[string]chan RPCResponse),
		Logger:       log.NewNopLogger(),
	}

	go func() {
		defer close(ws.quit)
		for {
			_, msg, err := ws.conn.Read(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "status = StatusNormalClosure") {
					ws.Logger.Debug("websocket normal closure")
					return
				}
				if strings.Contains(err.Error(), "context canceled") {
					ws.Logger.Debug("websocket context canceled")
					return
				}
				ws.Logger.Error("websocket read", "error", err.Error())
				return
			}
			if bytes.Equal(msg, []byte("{}")) {
				continue
			}
			var rpc RPCResponse
			if err = json.Unmarshal(msg, &rpc); err != nil {
				ws.Logger.Error("failed to unmarshal response", "error", err)
				continue
			}

			if bytes.Equal(rpc.Result, []byte("{}")) {
				continue
			}

			if rpc.Error != nil && rpc.Error.ServerExit() {
				ws.Logger.Error("websocket", "data", rpc.Error.Data)
				return
			}

			if _, ch := ws.response(rpc.ID); ch != nil {
				if cap(ch) == 0 {
					ch <- rpc
				} else {
					select {
					case ch <- rpc:
					default:
						ws.Logger.Error("wanted to publish response, but out channel is full. ", "rpc", rpc)
					}
				}
			} else {
				ws.Logger.Error("no found receive response chan.s", "rpc", rpc)
			}
		}
	}()
	return ws, nil
}

func (ws *WSClient) running() bool {
	select {
	case <-ws.quit:
		return false
	default:
		return true
	}
}

func (ws *WSClient) ExitCh() <-chan struct{} {
	return ws.quit
}

func (ws *WSClient) Close() {
	select {
	case <-ws.ctx.Done():
		return
	default:
		if !ws.running() {
			return
		}
		if err := ws.conn.Close(websocket.StatusNormalClosure, "close"); err != nil {
			ws.Logger.Debug("web socket close", "error", err.Error())
		}
	}
}

// SubscribeEvent Experiment
func (ws *WSClient) SubscribeEvent(ctx context.Context, query string, event chan<- coreTypes.ResultEvent) (err error) {
	id, err := ws.send(ctx, "subscribe", map[string]interface{}{"query": query})
	if err != nil {
		return
	}
	var response = make(chan RPCResponse, len(event))
	ws.addResponseChan(id, query, response)

	go func() {
		defer ws.Unsubscribe(id)
		for {
			select {
			case resp := <-response:
				if resp.Error != nil {
					ws.Logger.Error("response error", "code", resp.Error.Code, "data", resp.Error.Data, "msg", resp.Error.Message)
					continue
				}
				var res coreTypes.ResultEvent
				if err = tmjson.Unmarshal(resp.Result, &res); err != nil {
					ws.Logger.Error("Parse result event", "error", err)
				}
				event <- res
			case <-ws.quit:
				return
			case <-ctx.Done():
				ws.Logger.Debug("tendermint subscribe closed")
				return
			}
		}
	}()
	return nil
}

func (ws *WSClient) BlockResults(height int64) (*coreTypes.ResultBlockResults, error) {
	if !ws.running() {
		return nil, nil
	}
	result := new(coreTypes.ResultBlockResults)
	if err := ws.Call(ws.ctx, "block_results", map[string]interface{}{"height": fmt.Sprintf("%d", height)}, result); err != nil {
		ws.Logger.Error("Failed to Call block_results", "error", err.Error())
		return nil, err
	}
	return result, nil
}

/*
# js
var ws = new WebSocket("ws://localhost:26657/websocket")
ws.send(JSON.stringify({"jsonrpc":"2.0","id":"py-test","method":"subscribe","params":{"query":"tm.event='NewBlockHeader'"}}))
*/

func (ws *WSClient) Subscribe(ctx context.Context, query string, resp chan RPCResponse) (id string, err error) {
	id, err = ws.send(ctx, "subscribe", map[string]interface{}{"query": query})
	if err != nil {
		return
	}
	ws.addResponseChan(id, query, resp)
	return id, nil
}

func (ws *WSClient) Unsubscribe(subId string) {
	query, _ := ws.response(subId)
	if err := ws.delResponseChan(subId); err != nil {
		ws.Logger.Debug("Failed to delete response chan", "error", err.Error())
		return
	}
	if !ws.running() {
		return
	}
	var res = new(map[string]interface{})
	if err := ws.Call(ws.ctx, "unsubscribe", map[string]interface{}{"query": query}, res); err != nil {
		ws.Logger.Debug("Failed to unsubscribe", "error", err.Error())
	}
}

func (ws *WSClient) TxSearch(
	query string,
	prove bool,
	page,
	perPage *int,
	orderBy string,
) (*coreTypes.ResultTxSearch, error) {
	result := new(coreTypes.ResultTxSearch)
	params := map[string]interface{}{
		"query":    query,
		"prove":    prove,
		"order_by": orderBy,
	}

	if page != nil {
		params["page"] = page
	}
	if perPage != nil {
		params["per_page"] = perPage
	}

	err := ws.Call(ws.ctx, "tx_search", params, result)
	return result, err
}

func (ws *WSClient) ABCIQuery(path string, data libTypes.HexBytes) (*coreTypes.ResultABCIQuery, error) {
	result := new(coreTypes.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data, "height": "0", "prove": false}
	err := ws.Call(ws.ctx, "abci_query", params, result)
	return result, err
}

func (ws *WSClient) Status() (*coreTypes.ResultStatus, error) {
	result := new(coreTypes.ResultStatus)
	err := ws.Call(ws.ctx, "status", map[string]interface{}{}, result)
	return result, err
}

func (ws *WSClient) Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (err error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return
	}

	reqId := fmt.Sprintf("go-%s", RandomString(8))
	body, err := json.Marshal(NewRPCRequest(reqId, method, payload))
	if err != nil {
		return
	}

	// 添加接收结果的chan
	respChan := make(chan RPCResponse)
	ws.addResponseChan(reqId, "", respChan)
	defer func(ws *WSClient, id string) {
		err := ws.delResponseChan(id)
		if err != nil {
			ws.Logger.Debug("Failed to unsubscribe", "error", err.Error())
		}
	}(ws, reqId)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ws.ctx.Done():
		return ctx.Err()
	default:
		ws.Logger.Debug("Request web socket write ==>", "body", string(body))
		if err = ws.conn.Write(ctx, websocket.MessageText, body); err != nil {
			return
		}
	}

	resp := <-respChan
	if resp.Error != nil {
		return fmt.Errorf(resp.Error.String())
	}

	return tmjson.Unmarshal(resp.Result, result)
}

func (ws *WSClient) send(ctx context.Context, method string, params map[string]interface{}) (reqId string, err error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return
	}

	reqId = fmt.Sprintf("go-%s", RandomString(8))
	body, err := json.Marshal(NewRPCRequest(reqId, method, payload))
	if err != nil {
		return
	}

	ws.Logger.Debug("Request web socket write ==>", "body", string(body))
	return reqId, ws.conn.Write(ctx, websocket.MessageText, body)
}

func (ws *WSClient) addResponseChan(id, query string, respChan chan RPCResponse) {
	ws.mtx.Lock()
	defer ws.mtx.Unlock()
	ws.responseChan[id] = map[string]chan RPCResponse{query: respChan}
}

func (ws *WSClient) delResponseChan(id string) error {
	ws.mtx.Lock()
	defer ws.mtx.Unlock()
	if _, ok := ws.responseChan[id]; ok {
		delete(ws.responseChan, id)
		return nil
	} else {
		return errors.New("subscription not found")
	}
}

func (ws *WSClient) response(id string) (string, chan RPCResponse) {
	ws.mtx.Lock()
	defer ws.mtx.Unlock()
	for key, ch := range ws.responseChan[id] {
		return key, ch
	}
	return "", nil
}

type RPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"` // must be map[string]interface{} or []interface{}
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (err RPCError) ServerExit() bool {
	return err.Code == -32000 && err.Message == "Server error" && strings.Contains(err.Data, "Tendermint exited")
}

func (err RPCError) String() string {
	return fmt.Sprintf("code: %d, data: %s, msg: %s", err.Code, err.Data, err.Message)
}

func NewRPCRequest(id, method string, params json.RawMessage) RPCRequest {
	return RPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

var _ jsonRPCCaller = &FastClient{}

// FastClient implement jsonRPCCaller
type FastClient struct {
	Remote string
	cli    *http.Client
	log.Logger
}

func NewFastClient(remote string) *FastClient {
	return &FastClient{
		Remote: remote,
		cli:    http.DefaultClient,
		Logger: log.NewNopLogger(),
	}
}

func (cli *FastClient) SetTimeout(t time.Duration) {
	cli.cli.Timeout = t
}

func (cli *FastClient) Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (err error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return
	}

	reqId := fmt.Sprintf("go-%s", RandomString(8))
	body, err := json.Marshal(NewRPCRequest(reqId, method, payload))
	if err != nil {
		return
	}

	if method == "subscribe" {
		return errors.New("this method is not supported")
	}

	if strings.HasPrefix(cli.Remote, "tcp") {
		cli.Remote = strings.Replace(cli.Remote, "tcp", "http", 1)
	}

	cli.Logger.Debug("Request Post ==>", "remote", cli.Remote, "body", string(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cli.Remote, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/json")
	resp, err := cli.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	date, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cli.Logger.Debug("Response Body <==", "remote", cli.Remote, "body", string(date))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(date))
	}

	var rpcResp RPCResponse
	if err = json.Unmarshal(date, &rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("response code: %d, data: %s, msg: %s", rpcResp.Error.Code, rpcResp.Error.Data, rpcResp.Error.Message)
	}
	return tmjson.Unmarshal(rpcResp.Result, result)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	return string(b)
}
