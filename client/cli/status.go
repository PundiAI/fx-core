package cli

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/bytes"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/p2p"
)

// ValidatorInfo is info about the node's validator, same as Tendermint,
// except that we use our own PubKey.
type validatorInfo struct {
	Address     bytes.HexBytes
	PubKey      cryptotypes.PubKey
	VotingPower int64
}

type syncInfo struct {
	LatestBlockHash   bytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     bytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight string         `json:"latest_block_height"`
	LatestBlockTime   time.Time      `json:"latest_block_time"`

	EarliestBlockHash   bytes.HexBytes `json:"earliest_block_hash"`
	EarliestAppHash     bytes.HexBytes `json:"earliest_app_hash"`
	EarliestBlockHeight string         `json:"earliest_block_height"`
	EarliestBlockTime   time.Time      `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

// ResultStatus is node's info, same as Tendermint, except that we use our own
// PubKey.
type resultStatus struct {
	NodeInfo      p2p.DefaultNodeInfo
	SyncInfo      syncInfo
	ValidatorInfo validatorInfo
}

// StatusCommand returns the command to return the status of the network.
func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query remote node for status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			node, err := clientCtx.GetNode()
			if err != nil {
				return err
			}
			status, err := node.Status(context.Background())
			if err != nil {
				return err
			}

			// `status` has TM pubkeys, we need to convert them to our pubkeys.
			pk, err := cryptocodec.FromTmPubKeyInterface(status.ValidatorInfo.PubKey)
			if err != nil {
				return err
			}
			statusWithPk := resultStatus{
				NodeInfo: status.NodeInfo,
				SyncInfo: syncInfo{
					LatestBlockHash:     status.SyncInfo.LatestBlockHash,
					LatestAppHash:       status.SyncInfo.LatestAppHash,
					LatestBlockHeight:   strconv.FormatInt(status.SyncInfo.LatestBlockHeight, 10),
					LatestBlockTime:     status.SyncInfo.LatestBlockTime,
					EarliestBlockHash:   status.SyncInfo.EarliestBlockHash,
					EarliestAppHash:     status.SyncInfo.EarliestAppHash,
					EarliestBlockHeight: strconv.FormatInt(status.SyncInfo.EarliestBlockHeight, 10),
					EarliestBlockTime:   status.SyncInfo.EarliestBlockTime,
					CatchingUp:          status.SyncInfo.CatchingUp,
				},
				ValidatorInfo: validatorInfo{
					Address:     status.ValidatorInfo.Address,
					PubKey:      pk,
					VotingPower: status.ValidatorInfo.VotingPower,
				},
			}
			raw, err := json.Marshal(statusWithPk)
			if err != nil {
				return err
			}
			return clientCtx.PrintRaw(raw)
		},
	}

	cmd.Flags().StringP(tmcli.OutputFlag, "o", "json", "Output format (text|json)")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")

	return cmd
}
