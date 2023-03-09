package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" // nolint:staticcheck
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/gorilla/mux"
)

// RegisterSlashingRESTRoutes
// Deprecated
func RegisterSlashingRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := WithHTTPDeprecationHeaders(rtr)
	registerSlashingQueryRoutes(clientCtx, r)
}

func registerSlashingQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/slashing/validators/{validatorPubKey}/signing_info",
		signingInfoHandlerFn(clientCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/signing_infos",
		signingInfoHandlerListFn(clientCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/parameters",
		querySlashingParamsHandlerFn(clientCtx),
	).Methods("GET")
}

// Deprecated: http request handler to query signing info
func signingInfoHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pk, err := legacybech32.UnmarshalPubKey(legacybech32.ConsPK, vars["validatorPubKey"])
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.QuerySigningInfoRequest{ConsAddress: pk.Address().String()}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigningInfo)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// http request handler to query signing info
func signingInfoHandlerListFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := ParseHTTPArgsWithLimit(r, 0)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.NewQuerySigningInfosParams(page, limit)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckInternalServerError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigningInfos)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func querySlashingParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/parameters", types.QuerierRoute)

		res, height, err := clientCtx.QueryWithData(route, nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}
