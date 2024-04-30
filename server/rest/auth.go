package rest

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"
)

// RegisterAuthRESTRoutes registers the auth module REST routes.
// Deprecated
func RegisterAuthRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/auth/accounts/{address}", QueryAccountRequestHandlerFn(clientCtx)).Methods(MethodGet)
	r.HandleFunc("/auth/params", queryAuthParamsHandler(clientCtx)).Methods(MethodGet)
}

// QueryAccountRequestHandlerFn is the query accountREST Handler.
func QueryAccountRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		accGetter := authtypes.AccountRetriever{}

		account, height, err := accGetter.GetAccountWithHeight(clientCtx, addr)
		if err != nil {
			// Ref: https://github.com/cosmos/cosmos-sdk/issues/4923
			if err := accGetter.EnsureExists(clientCtx, addr); err != nil {
				clientCtx = clientCtx.WithHeight(height)
				PostProcessResponse(w, clientCtx, authtypes.BaseAccount{})
				return
			}

			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, account)
	}
}

func queryAuthParamsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", authtypes.QuerierRoute, authtypes.QueryParams)
		res, height, err := clientCtx.QueryWithData(route, nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// RegisterTxRESTRoutes registers all transaction routes on the provided router.
// Deprecated
func RegisterTxRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/txs/decode", DecodeTxRequestHandlerFn(clientCtx)).Methods("POST")
}

// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
func QueryTxRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		output, err := authtx.QueryTx(clientCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = packStdTxResponse(w, clientCtx, output)
		if err != nil {
			// Error is already returned by packStdTxResponse.
			return
		}

		if output.Empty() {
			WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		err = checkAminoMarshalError(clientCtx, output, "/cosmos/tx/v1beta1/txs/{txhash}")
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		PostProcessResponseBare(w, clientCtx, output)
	}
}

// QueryTxsRequestHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse query parameters: %s", err),
			)
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				QueryGenesisTxs(clientCtx, w)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			PostProcessResponseBare(w, clientCtx, txs)
			return
		}

		events, page, limit, err = ParseHTTPArgs(r)
		if CheckBadRequestError(w, err) {
			return
		}

		searchResult, err := authtx.QueryTxsByEvents(clientCtx, events, page, limit, "")
		if CheckInternalServerError(w, err) {
			return
		}

		for _, txRes := range searchResult.Txs {
			err = packStdTxResponse(w, clientCtx, txRes)
			if CheckInternalServerError(w, err) {
				return
			}
		}

		err = checkAminoMarshalError(clientCtx, searchResult, "/cosmos/tx/v1beta1/txs")
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		PostProcessResponseBare(w, clientCtx, searchResult)
	}
}

type (
	// DecodeReq defines a tx decoding request.
	DecodeReq struct {
		Tx string `json:"tx"`
	}

	// DecodeResp defines a tx decoding response.
	DecodeResp legacytx.StdTx
)

// DecodeTxRequestHandlerFn returns the decode tx REST handler. In particular,
// it takes base64-decoded bytes, decodes it from the Amino wire protocol,
// and responds with a json-formatted transaction.
func DecodeTxRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DecodeReq

		body, err := io.ReadAll(r.Body)
		if CheckBadRequestError(w, err) {
			return
		}

		// NOTE: amino is used intentionally here, don't migrate it
		err = clientCtx.LegacyAmino.UnmarshalJSON(body, &req)
		if CheckBadRequestError(w, err) {
			return
		}

		txBytes, err := base64.StdEncoding.DecodeString(req.Tx)
		if CheckBadRequestError(w, err) {
			return
		}

		stdTx, err := convertToStdTx(w, clientCtx, txBytes)
		if err != nil {
			// Error is already returned by convertToStdTx.
			return
		}

		response := DecodeResp(stdTx)

		err = checkAminoMarshalError(clientCtx, response, "/cosmos/tx/v1beta1/txs/decode")
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		PostProcessResponse(w, clientCtx, response)
	}
}

// packStdTxResponse takes a sdk.TxResponse, converts the Tx into a StdTx, and
// packs the StdTx again into the sdk.TxResponse Any. Amino then takes care of
// seamlessly JSON-outputting the Any.
func packStdTxResponse(w http.ResponseWriter, clientCtx client.Context, txRes *sdk.TxResponse) error {
	// We just unmarshalled from Tendermint, we take the proto Tx's raw
	// bytes, and convert them into a StdTx to be displayed.
	txBytes := txRes.Tx.Value
	stdTx, err := convertToStdTx(w, clientCtx, txBytes)
	if err != nil {
		return err
	}

	// Pack the amino stdTx into the TxResponse's Any.
	txRes.Tx = codectypes.UnsafePackAny(stdTx)

	return nil
}

// convertToStdTx converts tx proto binary bytes retrieved from Tendermint into
// a StdTx. Returns the StdTx, as well as a flag denoting if the function
// successfully converted or not.
func convertToStdTx(w http.ResponseWriter, clientCtx client.Context, txBytes []byte) (legacytx.StdTx, error) {
	txI, err := clientCtx.TxConfig.TxDecoder()(txBytes)
	if CheckBadRequestError(w, err) {
		return legacytx.StdTx{}, err
	}

	tx, ok := txI.(signing.Tx)
	if !ok {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("%+v is not backwards compatible with %T", tx, legacytx.StdTx{}))
		return legacytx.StdTx{}, errorsmod.Wrapf(errortypes.ErrInvalidType, "expected %T, got %T", (signing.Tx)(nil), txI)
	}

	stdTx, err := ConvertTxToStdTx(clientCtx.LegacyAmino, tx)
	if CheckBadRequestError(w, err) {
		return legacytx.StdTx{}, err
	}

	return stdTx, nil
}

// checkAminoMarshalError checks if there are errors with marshalling non-amino
// txs with amino.
func checkAminoMarshalError(ctx client.Context, resp interface{}, grpcEndPoint string) error {
	// LegacyAmino used intentionally here to handle the SignMode errors
	marshaler := ctx.LegacyAmino

	_, err := marshaler.MarshalJSON(resp)
	if err != nil {
		// If there's an unmarshalling error, we assume that it's because we're
		// using amino to unmarshal a non-amino tx.
		return fmt.Errorf("this transaction cannot be displayed via legacy REST endpoints, because it does not support"+
			" Amino serialization. Please either use CLI, gRPC, gRPC-gateway, or directly query the Tendermint RPC"+
			" endpoint to query this transaction. The new REST endpoint (via gRPC-gateway) is %s. Please also see the"+
			"REST endpoints migration guide at %s for more info", grpcEndPoint, DeprecationURL)
	}

	return nil
}
