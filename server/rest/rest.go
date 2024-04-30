package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/types"
)

// DeprecationURL is the URL for migrating deprecated REST endpoints to newer ones.
// https://github.com/cosmos/cosmos-sdk/issues/8019
const (
	DeprecationURL = "https://docs.cosmos.network/master/migrations/rest.html"

	DefaultPage    = 1
	DefaultLimit   = 30             // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
	TxMinHeightKey = "tx.minheight" // Inclusive minimum height filter
	TxMaxHeightKey = "tx.maxheight" // Inclusive maximum height filter

	MethodGet = "GET"
)

// addHTTPDeprecationHeaders is a mux middleware function for adding HTTP
// Deprecation headers to a http handler
func addHTTPDeprecationHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Link", "<"+DeprecationURL+">; rel=\"deprecation\"")
		w.Header().Set("Warning", "199 - \"this endpoint is deprecated and may not work as before, see deprecation link for more info\"")
		h.ServeHTTP(w, r)
	})
}

// WithHTTPDeprecationHeaders returns a new *mux.Router, identical to its input
// but with the addition of HTTP Deprecation headers. This is used to mark legacy
// amino REST endpoints as deprecated in the REST API.
func WithHTTPDeprecationHeaders(r *mux.Router) *mux.Router {
	subRouter := r.NewRoute().Subrouter()
	subRouter.Use(addHTTPDeprecationHeaders)
	return subRouter
}

// ParseQueryHeightOrReturnBadRequest sets the height to execute a query if set by the http request.
// It returns false if there was an error parsing the height.
func ParseQueryHeightOrReturnBadRequest(w http.ResponseWriter, clientCtx client.Context, r *http.Request) (client.Context, bool) {
	heightStr := r.FormValue("height")
	if heightStr != "" {
		height, err := strconv.ParseInt(heightStr, 10, 64)
		if CheckBadRequestError(w, err) {
			return clientCtx, false
		}

		if height < 0 {
			WriteErrorResponse(w, http.StatusBadRequest, "height must be equal or greater than zero")
			return clientCtx, false
		}

		if height > 0 {
			clientCtx = clientCtx.WithHeight(height)
		}
	} else {
		clientCtx = clientCtx.WithHeight(0)
	}

	return clientCtx, true
}

// CheckBadRequestError attaches an error message to an HTTP 400 BAD REQUEST response.
// Returns false when err is nil; it returns true otherwise.
func CheckBadRequestError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusBadRequest, err)
}

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(legacy.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
}

// CheckError takes care of writing an error response if err is not nil.
// Returns false when err is nil; it returns true otherwise.
func CheckError(w http.ResponseWriter, status int, err error) bool {
	if err != nil {
		WriteErrorResponse(w, status, err.Error())
		return true
	}

	return false
}

// CheckNotFoundError attaches an error message to an HTTP 404 NOT FOUND response.
// Returns false when err is nil; it returns true otherwise.
func CheckNotFoundError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusNotFound, err)
}

// ErrorResponse defines the attributes of a JSON error response.
type ErrorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, err string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err}
}

// PostProcessResponseBare post processes a body similar to PostProcessResponse
// except it does not wrap the body and inject the height.
func PostProcessResponseBare(w http.ResponseWriter, ctx client.Context, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch b := body.(type) {
	case []byte:
		resp = b

	default:
		resp, err = ctx.LegacyAmino.MarshalJSON(body)
		if CheckInternalServerError(w, err) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// CheckInternalServerError attaches an error message to an HTTP 500 INTERNAL SERVER ERROR response.
// Returns false when err is nil; it returns true otherwise.
func CheckInternalServerError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusInternalServerError, err)
}

// ParseHTTPArgs parses the request's URL and returns a slice containing all
// arguments pairs. It separates page and limit used for pagination.
func ParseHTTPArgs(r *http.Request) (tags []string, page, limit int, err error) {
	return ParseHTTPArgsWithLimit(r, DefaultLimit)
}

// ParseHTTPArgsWithLimit parses the request's URL and returns a slice containing
// all arguments pairs. It separates page and limit used for pagination where a
// default limit can be provided.
func ParseHTTPArgsWithLimit(r *http.Request, defaultLimit int) (tags []string, page, limit int, err error) {
	tags = make([]string, 0, len(r.Form))

	for key, values := range r.Form {
		if key == "page" || key == "limit" {
			continue
		}

		var value string
		value, err = url.QueryUnescape(values[0])

		if err != nil {
			return tags, page, limit, err
		}

		var tag string

		switch key {
		case types.TxHeightKey:
			tag = fmt.Sprintf("%s=%s", key, value)

		case TxMinHeightKey:
			tag = fmt.Sprintf("%s>=%s", types.TxHeightKey, value)

		case TxMaxHeightKey:
			tag = fmt.Sprintf("%s<=%s", types.TxHeightKey, value)

		default:
			tag = fmt.Sprintf("%s='%s'", key, value)
		}

		tags = append(tags, tag)
	}

	pageStr := r.FormValue("page")
	if pageStr == "" {
		page = DefaultPage
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return tags, page, limit, err
		} else if page <= 0 {
			return tags, page, limit, errors.New("page must greater than 0")
		}
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return tags, page, limit, err
		} else if limit <= 0 {
			return tags, page, limit, errors.New("limit must greater than 0")
		}
	}

	return tags, page, limit, nil
}

// ResponseWithHeight defines a response object type that wraps an original
// response with a height.
type ResponseWithHeight struct {
	Height int64           `json:"height"`
	Result json.RawMessage `json:"result"`
}

// NewResponseWithHeight creates a new ResponseWithHeight instance
func NewResponseWithHeight(height int64, result json.RawMessage) ResponseWithHeight {
	return ResponseWithHeight{
		Height: height,
		Result: result,
	}
}

// PostProcessResponse performs post processing for a REST response. The result
// returned to clients will contain two fields, the height at which the resource
// was queried at and the original result.
func PostProcessResponse(w http.ResponseWriter, ctx client.Context, resp interface{}) {
	var (
		result []byte
		err    error
	)

	if ctx.Height < 0 {
		WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("negative height in response").Error())
		return
	}

	// LegacyAmino used intentionally for REST
	marshaler := ctx.LegacyAmino

	switch res := resp.(type) {
	case []byte:
		result = res

	default:
		result, err = marshaler.MarshalJSON(resp)
		if CheckInternalServerError(w, err) {
			return
		}
	}

	wrappedResp := NewResponseWithHeight(ctx.Height, result)

	output, err := marshaler.MarshalJSON(wrappedResp)
	if CheckInternalServerError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(output)
}

// GasEstimateResponse defines a response definition for tx gas estimation.
type GasEstimateResponse struct {
	GasEstimate uint64 `json:"gas_estimate"`
}

// ParseUint64OrReturnBadRequest converts s to a uint64 value.
func ParseUint64OrReturnBadRequest(w http.ResponseWriter, s string) (n uint64, ok bool) {
	var err error

	n, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' is not a valid uint64", s))

		return n, false
	}

	return n, true
}
