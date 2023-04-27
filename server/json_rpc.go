package server

import (
	"net"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	ethlog "github.com/ethereum/go-ethereum/log"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/evmos/ethermint/rpc"
	"github.com/evmos/ethermint/server/config"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"golang.org/x/net/netutil"
)

// StartJSONRPC starts the JSON-RPC server
func StartJSONRPC(svrCtx *server.Context, clientCtx client.Context, tmRPCAddr, tmEndpoint string, config *config.Config, indexer ethermint.EVMTxIndexer) (*http.Server, error) {
	if len(config.JSONRPC.WsAddress) > 0 {
		svrCtx.Logger.Info("Starting JSON WebSocket server", "address", config.JSONRPC.WsAddress)

		// allocate separate WS connection to Tendermint
		tmWsClient, err := ConnectTmWS(tmRPCAddr, tmEndpoint, svrCtx.Logger)
		if err != nil {
			return nil, err
		}
		wsSrv := rpc.NewWebsocketsServer(clientCtx, svrCtx.Logger, tmWsClient, config)
		wsSrv.Start()
	}
	tmWsClient, err := ConnectTmWS(tmRPCAddr, tmEndpoint, svrCtx.Logger)
	if err != nil {
		return nil, err
	}

	logger := svrCtx.Logger.With("module", "geth")
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

	apis := rpc.GetRPCAPIs(svrCtx, clientCtx, tmWsClient, allowUnprotectedTxs, indexer, rpcAPIArr)
	for _, api := range apis {
		if err = rpcServer.RegisterName(api.Namespace, api.Service); err != nil {
			svrCtx.Logger.Error(
				"failed to register service in JSON RPC namespace",
				"namespace", api.Namespace,
				"service", api.Service,
			)
			return nil, err
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
	return httpSrv, nil
}

func ConnectTmWS(tmRPCAddr, tmEndpoint string, logger tmlog.Logger) (*rpcclient.WSClient, error) {
	tmWsClient, err := rpcclient.NewWS(tmRPCAddr, tmEndpoint,
		rpcclient.MaxReconnectAttempts(256),
		// rpcclient.ReadWait(120*time.Second),
		// rpcclient.WriteWait(120*time.Second),
		// rpcclient.PingPeriod(50*time.Second),
		rpcclient.OnReconnect(func() {
			logger.Debug("EVM RPC reconnects to Tendermint WS", "address", tmRPCAddr+tmEndpoint)
		}),
	)

	if err != nil {
		logger.Error(
			"Tendermint WS client could not be created",
			"address", tmRPCAddr+tmEndpoint,
			"error", err,
		)
		return nil, err
	} else if err := tmWsClient.OnStart(); err != nil {
		logger.Error(
			"Tendermint WS client could not start",
			"address", tmRPCAddr+tmEndpoint,
			"error", err,
		)
		return nil, err
	}

	return tmWsClient, nil
}

// Listen starts a net.Listener on the tcp network on the given address.
// If there is a specified MaxOpenConnections in the config, it will also set the limitListener.
func Listen(addr string, maxOpenConnections int) (net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	if maxOpenConnections > 0 {
		ln = netutil.LimitListener(ln, maxOpenConnections)
	}
	return ln, err
}
