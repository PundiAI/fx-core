package evm

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	initAmount        = big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(1e18))
	recipientMnemonic = "rebel knee sight blush remember clog spy arch siren kitchen panther response crime moment margin metal awful mansion head pioneer puppy fence around win"
)

func TestTransactionTransfer(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client := NewClient(t, DefaultGRPCUrl, DefaultNodeRPCUrl, DefaultEthUrl, DefaultMnemonic, EthHDPath)

	addr := common.HexToAddress("0xaC58d3199775c12C77E26146ffafE28c94804502")

	balance := client.Balance(addr)
	t.Log("addr", addr.Hex(), "balance", balance.String())

	client.Transfer(addr, big.NewInt(1000))

	balance = client.Balance(addr)
	t.Log("addr", addr.Hex(), "balance", balance.String())
}

func TestERC20(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client := NewClient(t, DefaultGRPCUrl, DefaultNodeRPCUrl, DefaultEthUrl, DefaultMnemonic, EthHDPath)

	addr, tx, erc20, err := DeployERC20Token(client.TransactOpts(), client.ethClient, big.NewInt(1e18), "USDT", "USDT", 6)
	require.NoError(t, err)
	t.Log("erc20", addr.Hex(), "hash", tx.Hash())
	client.PendingTx(tx)

	balances, err := erc20.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("addr balance", balances.String())

	recipient := client.HexAddress(recipientMnemonic)
	t.Log("recipient", recipient.String())
	client.Transfer(recipient, big.NewInt(1e18))

	t.Log("transfer to receipt")
	tx, err = erc20.Transfer(client.TransactOpts(), recipient, big.NewInt(1e9))
	require.NoError(t, err)
	client.PendingTx(tx)

	t.Log("approve receipt")
	tx, err = erc20.Approve(client.TransactOpts(), recipient, big.NewInt(1e3))
	require.NoError(t, err)
	client.PendingTx(tx)

	allowance, err := erc20.Allowance(nil, client.HexAddress(), recipient)
	require.NoError(t, err)
	t.Log("allowance recipient", allowance.String())

	t.Log("receipt transferFrom")
	tx, err = erc20.TransferFrom(client.TransactOpts(recipientMnemonic), client.HexAddress(), recipient, big.NewInt(1))
	require.NoError(t, err)
	client.PendingTx(tx)

	balances, err = erc20.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("addr balance", balances.String())

	balances, err = erc20.BalanceOf(nil, recipient)
	require.NoError(t, err)
	t.Log("recipient balance", balances.String())
}

func TestFIP20(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client := NewClient(t, DefaultGRPCUrl, DefaultNodeRPCUrl, DefaultEthUrl, DefaultMnemonic, EthHDPath)

	token := client.Token("FX")
	t.Log("wfx token", token.Hex())

	recipient := client.HexAddress(recipientMnemonic)
	t.Log("recipient", recipient.String())

	client.Transfer(recipient, initAmount)

	wfx, err := NewERC20Token(token, client.ethClient)
	require.NoError(t, err)

	b1, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	client.ConvertERC20(token, sdk.NewIntFromBigInt(b1), client.HexAddress().Bytes())

	b2, err := wfx.BalanceOf(nil, recipient)
	require.NoError(t, err)
	client.SetKey(recipientMnemonic, EthHDPath).ConvertERC20(token, sdk.NewIntFromBigInt(b2), recipient.Bytes())
	client = client.SetKey(DefaultMnemonic, EthHDPath)

	queryBalance(client, token, client.HexAddress(), recipient)

	t.Log("convert FX to addr")
	client.ConvertCoin(client.HexAddress(), sdk.NewCoin("FX", sdk.NewInt(10)))

	queryBalance(client, token, client.HexAddress(), recipient)

	t.Log("convert erc20 WFX to recipient")
	client.ConvertERC20(token, sdk.NewInt(1), recipient.Bytes())

	queryBalance(client, token, client.HexAddress(), recipient)

	t.Log("wfx transfer to recipient")
	tx, err := wfx.Transfer(client.TransactOpts(), recipient, big.NewInt(1))
	require.NoError(t, err)
	client.PendingTx(tx)

	queryBalance(client, token, client.HexAddress(), recipient)

	t.Log("wfx approve to recipient")
	tx, err = wfx.Approve(client.TransactOpts(), recipient, big.NewInt(1e3))
	require.NoError(t, err)
	client.PendingTx(tx)

	allowance, err := wfx.Allowance(nil, client.HexAddress(), recipient)
	require.NoError(t, err)
	t.Log("wfx allowance recipient", allowance.String())

	t.Log("wfx transferFrom to recipient")
	tx, err = wfx.TransferFrom(client.TransactOpts(recipientMnemonic), client.HexAddress(), recipient, big.NewInt(1))
	require.NoError(t, err)
	client.PendingTx(tx)

	queryBalance(client, token, client.HexAddress(), recipient)

	allowance, err = wfx.Allowance(nil, client.HexAddress(), recipient)
	require.NoError(t, err)
	t.Log("allowance recipient", allowance.String())
}

func TestFIP20CrossChain(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client := NewClient(t, DefaultGRPCUrl, DefaultNodeRPCUrl, DefaultEthUrl, DefaultMnemonic, EthHDPath)

	recipient := client.HexAddress(recipientMnemonic)
	t.Log("recipient", recipient.String())

	token := client.Token("FX")
	t.Log("wfx token", token.Hex())

	wfx, err := NewERC20Token(token, client.ethClient)
	require.NoError(t, err)

	client.GravityInitialize()

	b1 := client.Balance(recipient)
	t.Log("recipient balance", b1.String())

	client.GravitySendToTx(recipient.Bytes(), common.HexToAddress(EthFXTokenContract), big.NewInt(100), "")

	b2 := client.Balance(recipient)
	t.Log("recipient balance", b2.String())

	b3, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx addr balance", b3.String())

	client.GravitySendToTx(client.AccAddress(), common.HexToAddress(EthFXTokenContract), big.NewInt(100), "module/evm")

	b4 := client.Balance(recipient)
	t.Log("recipient balance", b4.String())

	b5, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx addr balance", b5.String())

	client.TransferCrossChain(token, client.HexAddress().String(), big.NewInt(10), big.NewInt(10), "chain/gravity")

	b6, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx addr balance", b6.String())

	client.GravityCheckPoolTx()
}

func TestFIP20IBCTransfer(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client := NewClient(t, DefaultGRPCUrl, DefaultNodeRPCUrl, DefaultEthUrl, DefaultMnemonic, EthHDPath)

	recipient := client.HexAddress(recipientMnemonic)
	t.Log("recipient", recipient.String())

	token := client.Token("FX")
	t.Log("wfx token", token.Hex())
	wfx, err := NewERC20Token(token, client.ethClient)
	require.NoError(t, err)

	client.CheckIBCChannelState("transfer", "channel-0")

	amount := big.NewInt(100)

	t.Log("convert FX to addr")
	client.ConvertCoin(client.HexAddress(), sdk.NewCoin("FX", sdk.NewIntFromBigInt(amount)))

	b1, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx addr balance", b1.String())
	client.TransferCrossChain(token, "px13rvykxlacsvpa0564pg6v8vf9xxeknrzg9xugy", amount, big.NewInt(0), "ibc/px/transfer/channel-0")

	b2, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx addr balance", b2.String())

	t.Logf("run command ===> pundixd tx ibc-transfer transfer transfer channel-0 %s %sibc/37CA072246C3BCBB445AEC196645F5AAB7876C456D76BC96141D3A0D6E615D2E --ibc-fee=0 --ibc-router=\"erc20\" --from fx1 --node tcp://0.0.0.0:27757\n", client.HexAddress(), amount)

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		b3, err := wfx.BalanceOf(nil, client.HexAddress())
		require.NoError(t, err)
		if b3.Cmp(b2) == 0 {
			continue
		}
		t.Log("wfx addr balance", b3.String())
		break
	}
}

func TestWFX(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client := NewClient(t, DefaultGRPCUrl, DefaultNodeRPCUrl, DefaultEthUrl, DefaultMnemonic, EthHDPath)

	token := client.Token("FX")
	t.Log("wfx token", token.Hex())

	recipient := client.HexAddress(recipientMnemonic)
	t.Log("recipient", recipient.String())
	client.Transfer(recipient, initAmount)

	wfx, err := NewERC20Token(token, client.ethClient)
	require.NoError(t, err)

	b1, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	client.ConvertERC20(token, sdk.NewIntFromBigInt(b1), client.HexAddress().Bytes())

	b2, err := wfx.BalanceOf(nil, recipient)
	require.NoError(t, err)
	client.SetKey(recipientMnemonic, EthHDPath).ConvertERC20(token, sdk.NewIntFromBigInt(b2), recipient.Bytes())

	client = client.SetKey(DefaultMnemonic, EthHDPath)

	client.Deposit(token, big.NewInt(10))

	b3, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx balance", b3.String())

	b4 := client.Balance(recipient)
	t.Log("recipient balance", b4.String())

	client.Withdraw(token, recipient, big.NewInt(1))

	b5, err := wfx.BalanceOf(nil, client.HexAddress())
	require.NoError(t, err)
	t.Log("wfx addr balance", b5.String())

	b6 := client.Balance(recipient)
	t.Log("wfx recipient balance", b6.String())
}

func queryBalance(c *Client, token common.Address, addrs ...common.Address) {
	erc20, err := NewERC20Token(token, c.ethClient)
	require.NoError(c.t, err)
	buf := bytes.Buffer{}
	buf.WriteString("balance ===> ")
	for _, addr := range addrs {
		balance, err := erc20.BalanceOf(nil, addr)
		require.NoError(c.t, err)
		buf.WriteString(fmt.Sprintf("%s-%s ", addr.Hex(), balance.String()))
	}
	c.t.Log(buf.String())
}
