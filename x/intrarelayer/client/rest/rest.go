package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/functionx/fx-core/x/intrarelayer/types"
)

// InitIntrarelayerParamsProposalRequest defines a request for a new init intrarelayer params proposal.
type InitIntrarelayerParamsProposalRequest struct {
	BaseReq     rest.BaseReq         `json:"base_req" yaml:"base_req"`
	Title       string               `json:"title" yaml:"title"`
	Description string               `json:"description" yaml:"description"`
	Deposit     sdk.Coins            `json:"deposit" yaml:"deposit"`
	Params      *types.Params        `json:"params" yaml:"params"`
	Metadata    []banktypes.Metadata `json:"metadata" yaml:"metadata"`
}

// RegisterCoinProposalRequest defines a request for a new register coin proposal.
type RegisterCoinProposalRequest struct {
	BaseReq     rest.BaseReq       `json:"base_req" yaml:"base_req"`
	Title       string             `json:"title" yaml:"title"`
	Description string             `json:"description" yaml:"description"`
	Deposit     sdk.Coins          `json:"deposit" yaml:"deposit"`
	Metadata    banktypes.Metadata `json:"metadata" yaml:"metadata"`
}

// RegisterFIP20ProposalRequest defines a request for a new register FIP20 proposal.
type RegisterFIP20ProposalRequest struct {
	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title        string       `json:"title" yaml:"title"`
	Description  string       `json:"description" yaml:"description"`
	Deposit      sdk.Coins    `json:"deposit" yaml:"deposit"`
	FIP20Address string       `json:"fip20_address" yaml:"fip20_address"`
}

// ToggleTokenRelayProposalRequest defines a request for a toggle token relay proposal.
type ToggleTokenRelayProposalRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Token       string       `json:"token" yaml:"token"`
}

func InitIntrarelayerParamsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler:  newInitIntrarelayerParamsProposalHandler(clientCtx),
	}
}

func RegisterCoinProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler:  newRegisterCoinProposalHandler(clientCtx),
	}
}

func RegisterFIP20ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler:  newRegisterFIP20ProposalHandler(clientCtx),
	}
}

func ToggleTokenRelayRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: types.ModuleName,
		Handler:  newToggleTokenRelayHandler(clientCtx),
	}
}

// nolint: dupl
func newInitIntrarelayerParamsProposalHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InitIntrarelayerParamsProposalRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewInitIntrarelayerParamsProposal(req.Title, req.Description, req.Params, req.Metadata)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// nolint: dupl
func newRegisterCoinProposalHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterCoinProposalRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewRegisterCoinProposal(req.Title, req.Description, req.Metadata)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// nolint: dupl
func newRegisterFIP20ProposalHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterFIP20ProposalRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewRegisterFIP20Proposal(req.Title, req.Description, req.FIP20Address)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// nolint: dupl
func newToggleTokenRelayHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ToggleTokenRelayProposalRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewToggleTokenRelayProposal(req.Title, req.Description, req.Token)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
