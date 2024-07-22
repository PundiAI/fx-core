package network

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	ethlog "github.com/ethereum/go-ethereum/log"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/evmos/ethermint/app/ante"
	"github.com/evmos/ethermint/rpc"
	"github.com/evmos/ethermint/rpc/stream"
	ethermintserver "github.com/evmos/ethermint/server"
	"github.com/evmos/ethermint/server/config"
	etherminttypes "github.com/evmos/ethermint/types"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const MaxRetry = 6

type AppWithPendingTxStream interface {
	RegisterPendingTxListener(listener ante.PendingTxListener)
}

// StartJSONRPC starts the JSON-RPC server
//
//gocyclo:ignore
func StartJSONRPC(ctx *server.Context,
	clientCtx client.Context,
	config *config.Config,
	indexer etherminttypes.EVMTxIndexer,
	app AppWithPendingTxStream,
) (*http.Server, chan struct{}, error) {
	logger := ctx.Logger.With("module", "geth")

	evtClient, ok := clientCtx.Client.(rpcclient.EventsClient)
	if !ok {
		return nil, nil, fmt.Errorf("client %T does not implement EventsClient", clientCtx.Client)
	}

	var rpcStream *stream.RPCStream
	var err error
	for i := 0; i < MaxRetry; i++ {
		rpcStream, err = stream.NewRPCStreams(evtClient, logger, clientCtx.TxConfig.TxDecoder())
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create rpc streams: %w", err)
	}

	app.RegisterPendingTxListener(rpcStream.ListenPendingTx)

	ethlog.Root().SetHandler(ethlog.FuncHandler(func(r *ethlog.Record) error {
		switch r.Lvl {
		case ethlog.LvlTrace, ethlog.LvlDebug:
			logger.Debug(r.Msg, r.Ctx...)
		case ethlog.LvlInfo, ethlog.LvlWarn:
			logger.Info(r.Msg, r.Ctx...)
		case ethlog.LvlError, ethlog.LvlCrit:
			logger.Error(r.Msg, r.Ctx...)
		}
		return nil
	}))

	rpcServer := ethrpc.NewServer()

	allowUnprotectedTxs := config.JSONRPC.AllowUnprotectedTxs
	rpcAPIArr := config.JSONRPC.API

	apis := rpc.GetRPCAPIs(ctx, clientCtx, rpcStream, allowUnprotectedTxs, indexer, rpcAPIArr)

	for _, api := range apis {
		if err := rpcServer.RegisterName(api.Namespace, api.Service); err != nil {
			ctx.Logger.Error(
				"failed to register service in JSON RPC namespace",
				"namespace", api.Namespace,
				"service", api.Service,
			)
			return nil, nil, err
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/", rpcServer.ServeHTTP).Methods("POST")

	handlerWithCors := cors.Default()
	if config.API.EnableUnsafeCORS {
		handlerWithCors = cors.AllowAll()
	}

	httpSrv := &http.Server{
		Addr:              config.JSONRPC.Address,
		Handler:           handlerWithCors.Handler(r),
		ReadHeaderTimeout: config.JSONRPC.HTTPTimeout,
		ReadTimeout:       config.JSONRPC.HTTPTimeout,
		WriteTimeout:      config.JSONRPC.HTTPTimeout,
		IdleTimeout:       config.JSONRPC.HTTPIdleTimeout,
	}
	httpSrvDone := make(chan struct{}, 1)

	ln, err := ethermintserver.Listen(httpSrv.Addr, config)
	if err != nil {
		return nil, nil, err
	}

	errCh := make(chan error)
	go func() {
		ctx.Logger.Info("Starting JSON-RPC server", "address", config.JSONRPC.Address)
		if err = httpSrv.Serve(ln); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				close(httpSrvDone)
				return
			}

			ctx.Logger.Error("failed to start JSON-RPC server", "error", err.Error())
			errCh <- err
		}
	}()

	select {
	case err = <-errCh:
		ctx.Logger.Error("failed to boot JSON-RPC server", "error", err.Error())
		return nil, nil, err
	case <-time.After(time.Second): // assume JSON RPC server started successfully
	}

	ctx.Logger.Info("Starting JSON WebSocket server", "address", config.JSONRPC.WsAddress)

	wsSrv := rpc.NewWebsocketsServer(clientCtx, ctx.Logger, rpcStream, config)
	wsSrv.Start()
	return httpSrv, httpSrvDone, nil
}
