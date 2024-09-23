package jsonrpc

import (
	"fmt"
	"strconv"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/pkg/errors"

	"github.com/functionx/fx-core/v8/client"
)

func (c *NodeRPC) GetModuleAccounts() ([]sdk.AccountI, error) {
	data, err := proto.Marshal(&authtypes.QueryModuleAccountsRequest{})
	if err != nil {
		return nil, err
	}
	result, err := c.ABCIQueryIsOk("/cosmos.auth.v1beta1.Query/ModuleAccounts", data)
	if err != nil {
		return nil, err
	}
	response := new(authtypes.QueryModuleAccountsResponse)
	if err = proto.Unmarshal(result.Response.Value, response); err != nil {
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

func (c *NodeRPC) QueryAccount(address string) (sdk.AccountI, error) {
	data, err := proto.Marshal(&authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}
	result, err := c.ABCIQueryIsOk("/cosmos.auth.v1beta1.Query/Account", data)
	if err != nil {
		return nil, err
	}
	response := new(authtypes.QueryAccountResponse)
	if err = proto.Unmarshal(result.Response.Value, response); err != nil {
		return nil, err
	}
	var account sdk.AccountI
	if err = client.NewAccountCodec().UnpackAny(response.GetAccount(), &account); err != nil {
		return nil, err
	}
	return account, nil
}

func (c *NodeRPC) QueryBalance(address string, denom string) (sdk.Coin, error) {
	data, err := proto.Marshal(&banktypes.QueryBalanceRequest{Address: address, Denom: denom})
	if err != nil {
		return sdk.Coin{}, err
	}
	result, err := c.ABCIQueryIsOk("/cosmos.bank.v1beta1.Query/Balance", data)
	if err != nil {
		return sdk.Coin{}, err
	}
	response := new(banktypes.QueryBalanceResponse)
	if err = proto.Unmarshal(result.Response.Value, response); err != nil {
		return sdk.Coin{}, err
	}
	return *response.Balance, nil
}

func (c *NodeRPC) QueryBalances(address string) (sdk.Coins, error) {
	data, err := proto.Marshal(&banktypes.QueryAllBalancesRequest{Address: address})
	if err != nil {
		return nil, err
	}
	result, err := c.ABCIQueryIsOk("/cosmos.bank.v1beta1.Query/AllBalances", data)
	if err != nil {
		return nil, err
	}
	response := new(banktypes.QueryAllBalancesResponse)
	if err = proto.Unmarshal(result.Response.Value, response); err != nil {
		return nil, err
	}
	return response.Balances, nil
}

func (c *NodeRPC) QuerySupply() (sdk.Coins, error) {
	result, err := c.ABCIQueryIsOk("/cosmos.bank.v1beta1.Query/TotalSupply", nil)
	if err != nil {
		return nil, err
	}
	response := new(banktypes.QueryTotalSupplyResponse)
	if err = proto.Unmarshal(result.Response.Value, response); err != nil {
		return nil, err
	}
	return response.Supply, nil
}

func (c *NodeRPC) GetGasPrices() (sdk.Coins, error) {
	if len(c.gasPrices) > 0 {
		return c.gasPrices, nil
	}
	result, err := c.ABCIQueryIsOk("/cosmos.base.node.v1beta1.Service/Config", nil)
	if err != nil {
		return sdk.Coins{}, err
	}
	response := new(node.ConfigResponse)
	if err = proto.Unmarshal(result.Response.Value, response); err != nil {
		return nil, err
	}
	coins, err := sdk.ParseCoinsNormalized(response.GetMinimumGasPrice())
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (c *NodeRPC) AppVersion() (string, error) {
	result, err := c.ABCIQueryIsOk("/app/version", nil)
	if err != nil {
		return "", err
	}
	return string(result.Response.Value), nil
}

func (c *NodeRPC) Store(path string) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/store/"+path, nil)
}

func (c *NodeRPC) PeersByAddressPort(port string) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/addr/"+port, nil)
}

func (c *NodeRPC) PeersById(id string) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/id/"+id, nil)
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
