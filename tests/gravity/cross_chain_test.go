package gravity

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/bytes"

	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

// Test Gravity Afterr - crate second validator and bind orchestrator
// ft staking create-validator --commission-max-change-rate=0.01 --commission-max-rate=0.2 --commission-rate=0.01 --min-self-delegation=1 --amount=200000000000000000000FX --pubkey="fxvalconspub1zcjduepqh6g8kjavsnlwtn5gz8djmm7fhu9j523ths74fk9zy4lv6uzyr93q0s9amq" --from fx2
// ft gravity set-orchestrator-address fxvaloper16wvwsmpp4y4ttgzknyr6kqla877jud6u8yzl8y fx16wvwsmpp4y4ttgzknyr6kqla877jud6u04lqey 0x6f1D09Fed11115d65E1071CD2109eDb300D80A27 --from fx2

func TestOrchestratorChain(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client := NewClient(t)

	go signPendingValsetRequest(client)

	setOrchestratorAddress(client)

	// 2 ETH -> fxcore
	// 2.1 DepositClaim
	ethToFx(client)

	// 2.2 DepositClaim and ibc router pundix
	ethToFxAndIbcTransfer(client)

	// 3. query user balance
	showAllBalance(client, client.FxAddress())

	// 4. send-to-eth
	fxToEth(client, 20)

	// 5. request batch
	batchRequest(client)

	// 6. confirm batch
	confirmBatch(client)
}

func signPendingValsetRequest(c *Client) {
	c.t.Helper()
	defer func() {
		c.t.Logf("sign pending valset request defer ....\n")
		if err := recover(); err != nil {
			c.t.Fatal(err)
		}
	}()
	gravityId := queryGravityId(c)
	requestParams := &gravitytypes.QueryLastPendingValsetRequestByAddrRequest{Address: c.FxAddress().String()}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		queryResponse, err := c.gravityQueryClient.LastPendingValsetRequestByAddr(c.ctx, requestParams)
		if err != nil {
			c.t.Logf("query last pending valset request is err!params:%+v, errors:%v\n", requestParams, err)
			continue
		}
		valsets := queryResponse.Valsets
		if len(valsets) <= 0 {
			continue
		}
		for _, valset := range valsets {
			checkpoint := valset.GetCheckpoint(gravityId)
			c.t.Logf("need confirm valset: nonce:%v EthAddress:%v\n", valset.Nonce, c.ethAddress.Hex())
			signature, err := gravitytypes.NewEthereumSignature(checkpoint, c.ethPrivKey)
			if err != nil {
				c.t.Log(err)
				continue
			}
			c.BroadcastTx(&[]sdk.Msg{
				&gravitytypes.MsgValsetConfirm{
					Nonce:        valset.Nonce,
					Orchestrator: c.FxAddress().String(),
					EthAddress:   c.ethAddress.Hex(),
					Signature:    hex.EncodeToString(signature),
				},
			})
		}
	}
}

var (
	chainGravityId string
)

func queryGravityId(c *Client) string {
	c.t.Helper()
	once := &sync.Once{}
	once.Do(func() {
		abciQuery, err := c.fxRpc.ABCIQuery(c.ctx, "/custom/gravity/gravityID", bytes.HexBytes{})
		if err != nil {
			c.t.Fatal(err)
		}
		if abciQuery.Response.Code != 0 {
			c.t.Fatal(abciQuery.Response.String())
		}
		err = c.encodingConfig.Amino.UnmarshalJSON(abciQuery.Response.Value, &chainGravityId)
		if err != nil {
			c.t.Fatal(err)
		}
		c.t.Logf("abci query result:%v\n", chainGravityId)
	})
	return chainGravityId
}

func confirmBatch(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      Confirm Batch      ####################\n")
	gravityId := queryGravityId(c)
	orchestrator := c.FxAddress()
	for {
		lastPendingBatchRequestResponse, err := c.gravityQueryClient.LastPendingBatchRequestByAddr(c.ctx, &gravitytypes.QueryLastPendingBatchRequestByAddrRequest{Address: orchestrator.String()})
		if err != nil {
			c.t.Fatal(err)
		}
		outgoingTxBatch := lastPendingBatchRequestResponse.Batch
		if outgoingTxBatch == nil {
			break
		}
		checkpoint, err := outgoingTxBatch.GetCheckpoint(gravityId)
		if err != nil {
			c.t.Fatal(err)
		}
		signatureBytes, err := gravitytypes.NewEthereumSignature(checkpoint, c.ethPrivKey)
		if err != nil {
			c.t.Fatal(err)
		}

		err = gravitytypes.ValidateEthereumSignature(checkpoint, signatureBytes, c.ethAddress.Hex())
		if err != nil {
			c.t.Fatal(err)
		}
		c.BroadcastTx(&[]sdk.Msg{
			&gravitytypes.MsgConfirmBatch{
				Nonce:         outgoingTxBatch.BatchNonce,
				TokenContract: outgoingTxBatch.TokenContract,
				EthSigner:     c.ethAddress.Hex(),
				Orchestrator:  orchestrator.String(),
				Signature:     hex.EncodeToString(signatureBytes),
			},
		})
		c.t.Logf("\n")
		time.Sleep(2 * time.Second)
	}
}

func batchRequest(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      Batch Request      ####################\n")
	batchFeeResponse, err := c.gravityQueryClient.BatchFees(c.ctx, &gravitytypes.QueryBatchFeeRequest{})
	if err != nil {
		c.t.Fatal(err)
	}
	orchestrator := c.FxAddress()
	feeReceive := c.ethAddress.String()
	msgList := make([]sdk.Msg, 0, len(batchFeeResponse.BatchFees))
	for _, batchToken := range batchFeeResponse.BatchFees {
		if batchToken.TotalTxs >= 20 {
			denomResponse, err := c.gravityQueryClient.ERC20ToDenom(c.ctx, &gravitytypes.QueryERC20ToDenomRequest{
				Erc20: batchToken.TokenContract,
			})
			if err != nil {
				c.t.Fatal(err)
			}
			if strings.HasPrefix(denomResponse.Denom, batchToken.TokenContract) {
				c.t.Logf("warn!!! not found token contract, tokenContract:[%v], erc20ToDenom response:[%v]\n", batchToken.TokenContract, denomResponse.Denom)
				continue
			}
			c.t.Logf("Send MsgRequestBatch: token:[%v], totalTxCount:[%v], totalFees:[%v]\n", denomResponse.Denom, batchToken.TotalTxs, batchToken.TotalFees)
			msgList = append(msgList, gravitytypes.NewMsgRequestBatch(orchestrator, denomResponse.Denom, batchToken.TotalFees, feeReceive))
		}
	}
	if len(msgList) <= 0 {
		return
	}
	c.BroadcastTx(&msgList)
	c.t.Logf("\n")
}

// fx -> eth
func fxToEth(c *Client, count int) {
	c.t.Helper()
	c.t.Logf("\n####################      FX to ETH      ####################\n")
	msgList := make([]sdk.Msg, 0, count)
	denom := fmt.Sprintf("%s%s", "eth", ethTokenContract)
	for i := 0; i < count; i++ {
		msgSendToEth := gravitytypes.NewMsgSendToEth(
			c.FxAddress(), ethTokenContract, sdk.NewCoin(denom, sdk.NewInt(1000000000000000)),
			sdk.NewCoin(denom, sdk.NewInt(10001)))
		msgList = append(msgList, msgSendToEth)
	}
	c.BroadcastTx(&msgList)
}

func showAllBalance(c *Client, address sdk.AccAddress) {
	c.t.Helper()
	c.t.Logf("\n####################      Query Address Balance      ####################\n")
	queryAllBalancesResponse, err := c.bankQueryClient.AllBalances(c.ctx, banktypes.NewQueryAllBalancesRequest(address, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: true,
	}))
	if err != nil {
		c.t.Fatal(err)
	}
	c.t.Logf("address: [%v] all balance\n", address.String())
	for _, balance := range queryAllBalancesResponse.Balances {
		c.t.Logf("denom:%v, amount:%v\n", balance.Denom, balance.Amount.String())
	}
	c.t.Logf("\n")
}

// eth -> fx
func ethToFx(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      ETH to FX      ####################\n")
	depositClaimMsg := gravitytypes.NewMsgDepositClaim(c.QueryFxLastEventNonce(), 3, ethTokenContract,
		sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000)), c.ethAddress.Hex(), c.FxAddress().String(), "", c.FxAddress().String())
	c.BroadcastTx(&[]sdk.Msg{depositClaimMsg})
	c.t.Logf("\n")
}

// eth -> fx and ibc -> pundix
func ethToFxAndIbcTransfer(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      ETH to FX to PUNDIX      ####################\n")
	// fxAddr: fx1u66dz4r6yg4xugz3ej27ejpd73helayz5y0xwr
	// pundixAddr: px1u66dz4r6yg4xugz3ej27ejpd73helayz09l7gx
	depositClaimMsg := gravitytypes.NewMsgDepositClaim(c.QueryFxLastEventNonce(), 5, ethTokenContract,
		sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000)), c.ethAddress.Hex(),
		"fx1u66dz4r6yg4xugz3ej27ejpd73helayz5y0xwr", hex.EncodeToString([]byte("px/transfer/channel-0")), c.FxAddress().String())
	c.BroadcastTx(&[]sdk.Msg{depositClaimMsg})
	c.t.Logf("\n")
}

func setOrchestratorAddress(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      Validator SetOrchestratorAddress      ####################\n")
	fxAddress := c.FxAddress()

	if !gethcommon.IsHexAddress(c.ethAddress.Hex()) {
		c.t.Fatal("eth address is invalid")
	}
	queryOrchestratorResponse, err := c.gravityQueryClient.GetDelegateKeyByOrchestrator(c.ctx, &gravitytypes.QueryDelegateKeyByOrchestratorRequest{
		OrchestratorAddress: fxAddress.String(),
	})
	if queryOrchestratorResponse != nil && len(queryOrchestratorResponse.EthAddress) > 0 {
		c.t.Logf("already set orchestrator address! address:[%v], validatorAddress:[%v], ethAddress:[%v]\n",
			fxAddress.String(), queryOrchestratorResponse.ValidatorAddress, queryOrchestratorResponse.EthAddress)
		return
	}
	if err != nil {
		if !strings.Contains(err.Error(), "No validator") {
			c.t.Fatal(err)
		}
		c.t.Logf("not found validator!!error msg:%v\n", err.Error())
	}
	msgSetOrchestratorAddress := gravitytypes.NewMsgSetOrchestratorAddress(sdk.ValAddress(fxAddress), fxAddress, c.ethAddress.Hex())
	c.BroadcastTx(&[]sdk.Msg{msgSetOrchestratorAddress})
	c.t.Logf("\n")
}
