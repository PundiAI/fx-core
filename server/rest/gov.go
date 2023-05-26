package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	gcutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/gorilla/mux"
)

// REST Variable names
const (
	RestParamsType     = "type"
	RestProposalID     = "proposal-id"
	RestDepositor      = "depositor"
	RestVoter          = "voter"
	RestProposalStatus = "status"
	// RestNumLimit       = "limit"

	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)

// RegisterGovRESTRoutes
// Deprecated
func RegisterGovRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	registerGovQueryRoutes(clientCtx, WithHTTPDeprecationHeaders(rtr))
}

func registerGovQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/gov/parameters/{%s}", RestParamsType), queryParamsHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/gov/proposals", queryProposalsWithParameterFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}", RestProposalID), queryProposalHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/proposer", RestProposalID), queryProposerHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits", RestProposalID), queryDepositsHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits/{%s}", RestProposalID, RestDepositor), queryDepositHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/tally", RestProposalID), queryTallyOnProposalHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes", RestProposalID), queryVotesOnProposalHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes/{%s}", RestProposalID, RestVoter), queryVoteHandlerFn(clientCtx)).Methods("GET")
}

func queryParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[RestParamsType]

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/gov/%s/%s", govv1beta1.QueryParams, paramType), nil)
		if CheckNotFoundError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := govv1beta1.NewQueryProposalParams(proposalID)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData("custom/gov/proposal", bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryDepositsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := govv1beta1.NewQueryProposalParams(proposalID)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData("custom/gov/proposal", bz)
		if CheckInternalServerError(w, err) {
			return
		}

		var proposal govv1beta1.Proposal
		if CheckInternalServerError(w, clientCtx.LegacyAmino.UnmarshalJSON(res, &proposal)) {
			return
		}

		// For inactive proposals we must query the txs directly to get the deposits
		// as they're no longer in state.
		propStatus := proposal.Status
		if !(propStatus == govv1beta1.StatusVotingPeriod || propStatus == govv1beta1.StatusDepositPeriod) {
			res, err = QueryDepositsByTxQuery(clientCtx, params)
		} else {
			res, _, err = clientCtx.QueryWithData("custom/gov/deposits", bz)
		}

		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// QueryDepositsByTxQuery will query for deposits via a direct txs tags query. It
// will fetch and build deposits directly from the returned txs and return a
// JSON marshalled result or any error that occurred.
//
// NOTE: SearchTxs is used to facilitate the txs query which does not currently
// support configurable pagination.
func QueryDepositsByTxQuery(clientCtx client.Context, params govv1beta1.QueryProposalParams) ([]byte, error) {
	var deposits []govv1beta1.Deposit

	// initial deposit was submitted with proposal, so must be queried separately
	initialDeposit, err := queryInitialDepositByTxQuery(clientCtx, params.ProposalID)
	if err != nil {
		return nil, err
	}

	if !initialDeposit.Amount.IsZero() {
		deposits = append(deposits, initialDeposit)
	}

	searchResult, err := combineEvents(
		clientCtx, defaultPage,
		// Query legacy Msgs event action
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgDeposit),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalDeposit, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
		},
		// Query proto Msgs event action
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgDeposit{})),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalDeposit, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
		},
	)
	if err != nil {
		return nil, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			if depMsg, ok := msg.(*govv1beta1.MsgDeposit); ok {
				deposits = append(deposits, govv1beta1.Deposit{
					Depositor:  depMsg.Depositor,
					ProposalId: params.ProposalID,
					Amount:     depMsg.Amount,
				})
			}
		}
	}

	bz, err := clientCtx.LegacyAmino.MarshalJSON(deposits)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// combineEvents queries txs by events with all events from each event group,
// and combines all those events together.
//
// Tx are indexed in tendermint via their Msgs `Type()`, which can be:
// - via legacy Msgs (amino or proto), their `Type()` is a custom string,
// - via ADR-031 proto msgs, their `Type()` is the protobuf FQ method name.
// In searching for events, we search for both `Type()`s, and we use the
// `combineEvents` function here to merge events.
func combineEvents(clientCtx client.Context, page int, eventGroups ...[]string) (*sdk.SearchTxsResult, error) {
	// Only the Txs field will be populated in the final SearchTxsResult.
	allTxs := []*sdk.TxResponse{}
	for _, events := range eventGroups {
		res, err := authtx.QueryTxsByEvents(clientCtx, events, page, defaultLimit, "")
		if err != nil {
			return nil, err
		}
		allTxs = append(allTxs, res.Txs...)
	}

	return &sdk.SearchTxsResult{Txs: allTxs}, nil
}

// queryInitialDepositByTxQuery will query for a initial deposit of a governance proposal by
// ID.
func queryInitialDepositByTxQuery(clientCtx client.Context, proposalID uint64) (govv1beta1.Deposit, error) {
	searchResult, err := combineEvents(
		clientCtx, defaultPage,
		// Query legacy Msgs event action
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgSubmitProposal),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeSubmitProposal, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", proposalID))),
		},
		// Query proto Msgs event action
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgSubmitProposal{})),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeSubmitProposal, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", proposalID))),
		},
	)
	if err != nil {
		return govv1beta1.Deposit{}, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			// there should only be a single proposal under the given conditions
			if subMsg, ok := msg.(*govv1beta1.MsgSubmitProposal); ok {
				return govv1beta1.Deposit{
					ProposalId: proposalID,
					Depositor:  subMsg.Proposer,
					Amount:     subMsg.InitialDeposit,
				}, nil
			}
		}
	}

	return govv1beta1.Deposit{}, errortypes.ErrNotFound.Wrapf("failed to find the initial deposit for proposalID %d", proposalID)
}

func queryProposerHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, err := gcutils.QueryProposerByTxQuery(clientCtx, proposalID)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

func queryDepositHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechDepositorAddr := vars[RestDepositor]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechDepositorAddr) == 0 {
			err := errors.New("depositor address required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := govv1beta1.NewQueryDepositParams(proposalID, depositorAddr)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData("custom/gov/deposit", bz)
		if CheckInternalServerError(w, err) {
			return
		}

		var deposit govv1beta1.Deposit
		if CheckBadRequestError(w, clientCtx.LegacyAmino.UnmarshalJSON(res, &deposit)) {
			return
		}

		// For an empty deposit, either the proposal does not exist or is inactive in
		// which case the deposit would be removed from state and should be queried
		// for directly via a txs query.
		if deposit.Empty() {
			bz, err := clientCtx.LegacyAmino.MarshalJSON(govv1beta1.NewQueryProposalParams(proposalID))
			if CheckBadRequestError(w, err) {
				return
			}

			res, _, err = clientCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := fmt.Errorf("proposalID %d does not exist", proposalID)
				WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			res, err = QueryDepositByTxQuery(clientCtx, params)
			if CheckInternalServerError(w, err) {
				return
			}
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// QueryDepositByTxQuery will query for a single deposit via a direct txs tags
// query.
func QueryDepositByTxQuery(clientCtx client.Context, params govv1beta1.QueryDepositParams) ([]byte, error) {
	// initial deposit was submitted with proposal, so must be queried separately
	initialDeposit, err := queryInitialDepositByTxQuery(clientCtx, params.ProposalID)
	if err != nil {
		return nil, err
	}

	if !initialDeposit.Amount.IsZero() {
		bz, err := clientCtx.Codec.MarshalJSON(&initialDeposit)
		if err != nil {
			return nil, err
		}

		return bz, nil
	}

	searchResult, err := combineEvents(
		clientCtx, defaultPage,
		// Query legacy Msgs event action
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgDeposit),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalDeposit, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Depositor.String())),
		},
		// Query proto Msgs event action
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgDeposit{})),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalDeposit, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Depositor.String())),
		},
	)
	if err != nil {
		return nil, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			// there should only be a single deposit under the given conditions
			if depMsg, ok := msg.(*govv1beta1.MsgDeposit); ok {
				deposit := govv1beta1.Deposit{
					Depositor:  depMsg.Depositor,
					ProposalId: params.ProposalID,
					Amount:     depMsg.Amount,
				}

				bz, err := clientCtx.Codec.MarshalJSON(&deposit)
				if err != nil {
					return nil, err
				}

				return bz, nil
			}
		}
	}

	return nil, fmt.Errorf("address '%s' did not deposit to proposalID %d", params.Depositor, params.ProposalID)
}

func queryVoteHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechVoterAddr := vars[RestVoter]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechVoterAddr) == 0 {
			err := errors.New("voter address required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := govv1beta1.NewQueryVoteParams(proposalID, voterAddr)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData("custom/gov/vote", bz)
		if CheckInternalServerError(w, err) {
			return
		}

		var vote govv1beta1.Vote
		if CheckBadRequestError(w, clientCtx.LegacyAmino.UnmarshalJSON(res, &vote)) {
			return
		}

		// For an empty vote, either the proposal does not exist or is inactive in
		// which case the vote would be removed from state and should be queried for
		// directly via a txs query.
		if vote.Empty() {
			bz, err := clientCtx.LegacyAmino.MarshalJSON(govv1beta1.NewQueryProposalParams(proposalID))
			if CheckBadRequestError(w, err) {
				return
			}

			res, _, err = clientCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := fmt.Errorf("proposalID %d does not exist", proposalID)
				WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			res, err = QueryVoteByTxQuery(clientCtx, params)
			if CheckInternalServerError(w, err) {
				return
			}
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// QueryVoteByTxQuery will query for a single vote via a direct txs tags query.
func QueryVoteByTxQuery(clientCtx client.Context, params govv1beta1.QueryVoteParams) ([]byte, error) {
	searchResult, err := combineEvents(
		clientCtx, defaultPage,
		// Query legacy Vote Msgs
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgVote),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Voter.String())),
		},
		// Query Vote proto Msgs
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgVote{})),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Voter.String())),
		},
		// Query legacy VoteWeighted Msgs
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgVoteWeighted),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Voter.String())),
		},
		// Query VoteWeighted proto Msgs
		[]string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgVoteWeighted{})),
			fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Voter.String())),
		},
	)
	if err != nil {
		return nil, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			// there should only be a single vote under the given conditions
			var vote *govv1beta1.Vote
			if voteMsg, ok := msg.(*govv1beta1.MsgVote); ok {
				vote = &govv1beta1.Vote{
					Voter:      voteMsg.Voter,
					ProposalId: params.ProposalID,
					Options:    govv1beta1.NewNonSplitVoteOption(voteMsg.Option),
				}
			}

			if voteWeightedMsg, ok := msg.(*govv1beta1.MsgVoteWeighted); ok {
				vote = &govv1beta1.Vote{
					Voter:      voteWeightedMsg.Voter,
					ProposalId: params.ProposalID,
					Options:    voteWeightedMsg.Options,
				}
			}

			if vote != nil {
				bz, err := clientCtx.Codec.MarshalJSON(vote)
				if err != nil {
					return nil, err
				}

				return bz, nil
			}
		}
	}

	return nil, fmt.Errorf("address '%s' did not vote on proposalID %d", params.Voter, params.ProposalID)
}

func queryVotesOnProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := ParseHTTPArgs(r)
		if CheckBadRequestError(w, err) {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(govv1beta1.NewQueryProposalParams(proposalID))
		if CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData("custom/gov/proposal", bz)
		if CheckInternalServerError(w, err) {
			return
		}

		var proposal govv1beta1.Proposal
		if CheckInternalServerError(w, clientCtx.LegacyAmino.UnmarshalJSON(res, &proposal)) {
			return
		}

		// For inactive proposals we must query the txs directly to get the votes
		// as they're no longer in state.
		params := govv1beta1.NewQueryProposalVotesParams(proposalID, page, limit)

		propStatus := proposal.Status
		if !(propStatus == govv1beta1.StatusVotingPeriod || propStatus == govv1beta1.StatusDepositPeriod) {
			res, err = QueryVotesByTxQuery(clientCtx, params)
		} else {
			bz, err = clientCtx.LegacyAmino.MarshalJSON(params)
			if CheckBadRequestError(w, err) {
				return
			}

			res, _, err = clientCtx.QueryWithData("custom/gov/votes", bz)
		}

		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponse(w, clientCtx, res)
	}
}

// QueryVotesByTxQuery will query for votes via a direct txs tags query. It
// will fetch and build votes directly from the returned txs and return a JSON
// marshalled result or any error that occurred.
func QueryVotesByTxQuery(clientCtx client.Context, params govv1beta1.QueryProposalVotesParams) ([]byte, error) {
	var (
		votes      []govv1beta1.Vote
		nextTxPage = defaultPage
		totalLimit = params.Limit * params.Page
	)

	// query interrupted either if we collected enough votes or tx indexer run out of relevant txs
	for len(votes) < totalLimit {
		// Search for both (legacy) votes and weighted votes.
		searchResult, err := combineEvents(
			clientCtx, nextTxPage,
			// Query legacy Vote Msgs
			[]string{
				fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgVote),
				fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			},
			// Query Vote proto Msgs
			[]string{
				fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgVote{})),
				fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			},
			// Query legacy VoteWeighted Msgs
			[]string{
				fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, govv1beta1.TypeMsgVoteWeighted),
				fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			},
			// Query VoteWeighted proto Msgs
			[]string{
				fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, sdk.MsgTypeURL(&govv1beta1.MsgVoteWeighted{})),
				fmt.Sprintf("%s.%s='%s'", govtypes.EventTypeProposalVote, govtypes.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
			},
		)
		if err != nil {
			return nil, err
		}

		for _, info := range searchResult.Txs {
			for _, msg := range info.GetTx().GetMsgs() {
				if voteMsg, ok := msg.(*govv1beta1.MsgVote); ok {
					votes = append(votes, govv1beta1.Vote{
						Voter:      voteMsg.Voter,
						ProposalId: params.ProposalID,
						Options:    govv1beta1.NewNonSplitVoteOption(voteMsg.Option),
					})
				}

				if voteWeightedMsg, ok := msg.(*govv1beta1.MsgVoteWeighted); ok {
					votes = append(votes, govv1beta1.Vote{
						Voter:      voteWeightedMsg.Voter,
						ProposalId: params.ProposalID,
						Options:    voteWeightedMsg.Options,
					})
				}
			}
		}
		if len(searchResult.Txs) != defaultLimit {
			break
		}

		nextTxPage++
	}
	start, end := client.Paginate(len(votes), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		votes = []govv1beta1.Vote{}
	} else {
		votes = votes[start:end]
	}

	bz, err := clientCtx.LegacyAmino.MarshalJSON(votes)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// HTTP request handler to query list of governance proposals
func queryProposalsWithParameterFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := ParseHTTPArgsWithLimit(r, 0)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var (
			voterAddr      sdk.AccAddress
			depositorAddr  sdk.AccAddress
			proposalStatus govv1beta1.ProposalStatus
		)

		if v := r.URL.Query().Get(RestVoter); len(v) != 0 {
			voterAddr, err = sdk.AccAddressFromBech32(v)
			if CheckBadRequestError(w, err) {
				return
			}
		}

		if v := r.URL.Query().Get(RestDepositor); len(v) != 0 {
			depositorAddr, err = sdk.AccAddressFromBech32(v)
			if CheckBadRequestError(w, err) {
				return
			}
		}

		if v := r.URL.Query().Get(RestProposalStatus); len(v) != 0 {
			proposalStatus, err = govv1beta1.ProposalStatusFromString(gcutils.NormalizeProposalStatus(v))
			if CheckBadRequestError(w, err) {
				return
			}
		}

		params := govv1beta1.NewQueryProposalsParams(page, limit, proposalStatus, voterAddr, depositorAddr)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", govtypes.QuerierRoute, govv1beta1.QueryProposals)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryTallyOnProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		clientCtx, ok = ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := govv1beta1.NewQueryProposalParams(proposalID)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData("custom/gov/tally", bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}
