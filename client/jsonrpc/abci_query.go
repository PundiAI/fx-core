package jsonrpc

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	etherminttypes "github.com/evmos/ethermint/types"
	"github.com/gogo/protobuf/proto"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	gravitytypes "github.com/functionx/fx-core/v7/x/gravity/types"
)

func (c *NodeRPC) QueryAccount(address string) (authtypes.AccountI, error) {
	result, err := c.ABCIQueryIsOk("/custom/auth/account", legacy.Cdc.MustMarshalJSON(authtypes.QueryAccountRequest{Address: address}))
	if err != nil {
		return nil, err
	}
	var account authtypes.AccountI
	if err = authtypes.ModuleCdc.UnmarshalInterfaceJSON(result.Response.Value, &account); err != nil {
		account = new(etherminttypes.EthAccount)
		if err1 := legacy.Cdc.UnmarshalJSON(result.Response.Value, account); err1 != nil {
			return nil, fmt.Errorf("%s: %s", err.Error(), err1.Error())
		}
	}
	return account, nil
}

func (c *NodeRPC) QueryBalance(address string, denom string) (sdk.Coin, error) {
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

func (c *NodeRPC) QueryBalances(address string) (sdk.Coins, error) {
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

func (c *NodeRPC) QuerySupply() (sdk.Coins, error) {
	result, err := c.ABCIQueryIsOk("/custom/bank/total_supply", legacy.Cdc.MustMarshalJSON(banktypes.QueryTotalSupplyRequest{}))
	if err != nil {
		return nil, err
	}
	var supplyRes banktypes.QueryTotalSupplyResponse
	if err = legacy.Cdc.UnmarshalJSON(result.Response.Value, &supplyRes); err != nil {
		return nil, err
	}
	return supplyRes.Supply, nil
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

func (c *NodeRPC) Store(path string) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/store/"+path, nil)
}

func (c *NodeRPC) PeersByAddressPort(port string) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/addr/"+port, nil)
}

func (c *NodeRPC) PeersById(id string) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/id/"+id, nil)
}

// Deprecated: GetGravityAttestation
func (c *NodeRPC) GetGravityAttestation(cdc codec.Codec, id []byte) (*gravitytypes.Attestation, error) {
	query, err := c.ABCIQuery("/store/gravity/key", id)
	if err != nil {
		return nil, err
	}
	if query.Response.Code != 0 {
		return nil, fmt.Errorf("abci query code %d, space %s, log %s", query.Response.Code, query.Response.Codespace, query.Response.Log)
	}
	var gravityAtt gravitytypes.Attestation
	cdc.MustUnmarshal(query.Response.Value, &gravityAtt)
	return &gravityAtt, nil
}

// Deprecated: GetGravityLastObservedEventNonce
func (c *NodeRPC) GetGravityLastObservedEventNonce() (uint64, error) {
	query, err := c.ABCIQuery("/store/gravity/key", []byte{0xc})
	if err != nil {
		return 0, err
	}
	if query.Response.Code != 0 {
		return 0, fmt.Errorf("abci query code %d, space %s, log %s", query.Response.Code, query.Response.Codespace, query.Response.Log)
	}
	if len(query.Response.Value) == 0 {
		return 0, nil
	}
	return sdk.BigEndianToUint64(query.Response.Value), nil
}
