package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/gorilla/mux"
)

// RegisterMintRESTRoutes registers minting module REST handlers on the provided router.
func RegisterMintRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := WithHTTPDeprecationHeaders(rtr)
	registerMintQueryRoutes(clientCtx, r)
}

func registerMintQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/minting/parameters",
		queryMintParamsHandlerFn(clientCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/inflation",
		queryInflationHandlerFn(clientCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/annual-provisions",
		queryAnnualProvisionsHandlerFn(clientCtx),
	).Methods("GET")
}

func queryMintParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryParameters)

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(route, nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryInflationHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryInflation)

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(route, nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryAnnualProvisionsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryAnnualProvisions)

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(route, nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}
