package rest

import (
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/gorilla/mux"
)

// RegisterUpgradeRESTRoutes registers REST routes for the upgrade module under the path specified by routeName.
func RegisterUpgradeRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := WithHTTPDeprecationHeaders(rtr)
	registerUpgradeQueryRoutes(clientCtx, r)
}

func registerUpgradeQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/upgrade/current", getCurrentPlanHandler(clientCtx),
	).Methods("GET")
	r.HandleFunc(
		"/upgrade/applied/{name}", getDonePlanHandler(clientCtx),
	).Methods("GET")
}

func getCurrentPlanHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		// ignore height for now
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryCurrent))
		if CheckInternalServerError(w, err) {
			return
		}
		if len(res) == 0 {
			http.NotFound(w, request)
			return
		}

		var plan types.Plan
		err = clientCtx.LegacyAmino.UnmarshalJSON(res, &plan)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponse(w, clientCtx, plan)
	}
}

func getDonePlanHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		params := types.QueryAppliedPlanRequest{Name: name}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryApplied), bz)
		if CheckBadRequestError(w, err) {
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}
		if len(res) != 8 {
			WriteErrorResponse(w, http.StatusInternalServerError, "unknown format for applied-upgrade")
		}

		applied := int64(binary.BigEndian.Uint64(res))
		fmt.Println(applied)
		PostProcessResponse(w, clientCtx, applied)
	}
}
