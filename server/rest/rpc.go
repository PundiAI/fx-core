package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// RegisterRPCRoutes Register REST endpoints.
func RegisterRPCRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/node_info", NodeInfoRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/syncing", NodeSyncingRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/blocks/latest", LatestBlockRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}", BlockRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/latest", LatestValidatorSetRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/{height}", ValidatorSetRequestHandlerFn(clientCtx)).Methods("GET")
}

// NodeInfoResponse defines a response type that contains node status and version
// information.
type NodeInfoResponse struct {
	p2p.DefaultNodeInfo `json:"node_info"`

	ApplicationVersion version.Info `json:"application_version"`
}

// NodeInfoRequestHandlerFn REST handler for node info
func NodeInfoRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := getNodeStatus(clientCtx)
		if CheckInternalServerError(w, err) {
			return
		}

		resp := NodeInfoResponse{
			DefaultNodeInfo:    status.NodeInfo,
			ApplicationVersion: version.NewInfo(),
		}

		PostProcessResponseBare(w, clientCtx, resp)
	}
}

func getNodeStatus(clientCtx client.Context) (*ctypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}

	return node.Status(context.Background())
}

// SyncingResponse defines a response type that contains node syncing information.
type SyncingResponse struct {
	Syncing bool `json:"syncing"`
}

// NodeSyncingRequestHandlerFn REST handler for node syncing
func NodeSyncingRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := getNodeStatus(clientCtx)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponseBare(w, clientCtx, SyncingResponse{Syncing: status.SyncInfo.CatchingUp})
	}
}

// LatestBlockRequestHandlerFn REST handler to get the latest block
func LatestBlockRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlock(clientCtx, nil)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponseBare(w, clientCtx, output)
	}
}

func getBlock(clientCtx client.Context, height *int64) ([]byte, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(context.Background(), height)
	if err != nil {
		return nil, err
	}

	return legacy.Cdc.MarshalJSON(res)
}

// BlockRequestHandlerFn REST handler to get a block
func BlockRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest,
				"couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}

		chainHeight, err := GetChainHeight(clientCtx)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlock(clientCtx, &height)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponseBare(w, clientCtx, output)
	}
}

// GetChainHeight get the current blockchain height
func GetChainHeight(clientCtx client.Context) (int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return -1, err
	}

	status, err := node.Status(context.Background())
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// LatestValidatorSetRequestHandlerFn Latest Validator Set REST handler
func LatestValidatorSetRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse pagination parameters")
			return
		}

		output, err := GetValidators(r.Context(), clientCtx, nil, &page, &limit)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponse(w, clientCtx, output)
	}
}

// ValidatorOutput Validator output
type ValidatorOutput struct {
	Address          sdk.ConsAddress    `json:"address"`
	PubKey           cryptotypes.PubKey `json:"pub_key"`
	ProposerPriority int64              `json:"proposer_priority"`
	VotingPower      int64              `json:"voting_power"`
}

// ResultValidatorsOutput Validators at a certain height output in bech32 format
type ResultValidatorsOutput struct {
	BlockHeight int64             `json:"block_height"`
	Validators  []ValidatorOutput `json:"validators"`
	Total       uint64            `json:"total"`
}

func (rvo ResultValidatorsOutput) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("block height: %d\n", rvo.BlockHeight))
	b.WriteString(fmt.Sprintf("total count: %d\n", rvo.Total))

	for _, val := range rvo.Validators {
		b.WriteString(
			fmt.Sprintf(`
  Address:          %s
  Pubkey:           %s
  ProposerPriority: %d
  VotingPower:      %d
		`,
				val.Address, val.PubKey, val.ProposerPriority, val.VotingPower,
			),
		)
	}

	return b.String()
}

// GetValidators from client
func GetValidators(ctx context.Context, clientCtx client.Context, height *int64, page, limit *int) (ResultValidatorsOutput, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return ResultValidatorsOutput{}, err
	}

	validatorsRes, err := node.Validators(ctx, height, page, limit)
	if err != nil {
		return ResultValidatorsOutput{}, err
	}

	total := validatorsRes.Total
	if validatorsRes.Total < 0 {
		total = 0
	}
	out := ResultValidatorsOutput{
		BlockHeight: validatorsRes.BlockHeight,
		Validators:  make([]ValidatorOutput, len(validatorsRes.Validators)),
		Total:       uint64(total),
	}
	for i := 0; i < len(validatorsRes.Validators); i++ {
		out.Validators[i], err = validatorOutput(validatorsRes.Validators[i])
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func validatorOutput(validator *tmtypes.Validator) (ValidatorOutput, error) {
	pk, err := cryptocodec.FromTmPubKeyInterface(validator.PubKey)
	if err != nil {
		return ValidatorOutput{}, err
	}

	return ValidatorOutput{
		Address:          sdk.ConsAddress(validator.Address),
		PubKey:           pk,
		ProposerPriority: validator.ProposerPriority,
		VotingPower:      validator.VotingPower,
	}, nil
}

// ValidatorSetRequestHandlerFn Validator Set at a height REST handler
func ValidatorSetRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse pagination parameters")
			return
		}

		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "failed to parse block height")
			return
		}

		chainHeight, err := GetChainHeight(clientCtx)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}
		if height > chainHeight {
			WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := GetValidators(r.Context(), clientCtx, &height, &page, &limit)
		if CheckInternalServerError(w, err) {
			return
		}
		PostProcessResponse(w, clientCtx, output)
	}
}
