package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"

	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

func (c *JsonRPC) QueryAccount(address string) (authtypes.AccountI, error) {
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

func (c *JsonRPC) QueryBalance(address string, denom string) (sdk.Coin, error) {
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

func (c *JsonRPC) QueryBalances(address string) (sdk.Coins, error) {
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

func (c *JsonRPC) QuerySupply() (sdk.Coins, error) {
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

// Deprecated: GetGasPrices
func (c *JsonRPC) GetGasPrices() (sdk.Coins, error) {
	result, err := c.ABCIQueryIsOk("/custom/other/gasPrice", nil)
	if err != nil {
		return sdk.Coins{}, err
	}
	var gasPrice sdk.Coins
	if err = json.Unmarshal(result.Response.Value, &gasPrice); err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (c *JsonRPC) Store(path string) (*coreTypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/store/"+path, nil)
}

func (c *JsonRPC) PeersByAddressPort(port string) (*coreTypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/addr/"+port, nil)
}

func (c *JsonRPC) PeersById(id string) (*coreTypes.ResultABCIQuery, error) {
	return c.ABCIQueryIsOk("/p2p/filter/id/"+id, nil)
}

func (c *JsonRPC) GetGravityAttestation(cdc codec.Codec, id []byte) (*gravitytypes.Attestation, error) {
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

func (c *JsonRPC) GetGravityLastObservedEventNonce() (uint64, error) {
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
	return gravitytypes.UInt64FromBytes(query.Response.Value), nil
}
