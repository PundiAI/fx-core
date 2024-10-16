package crosschain

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	crosschainv1 "github.com/functionx/fx-core/v8/api/fx/gravity/crosschain/v1"
)

type AutoCLIAppModule struct {
	ModuleName string
}

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AutoCLIAppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	// add chain_name flag to all rpc commands
	rpcCommandOptions := make([]*autocliv1.RpcCommandOptions, 0, len(crosschainv1.Query_ServiceDesc.Methods))
	for _, method := range crosschainv1.Query_ServiceDesc.Methods {
		// exclude QueryBridgeChainListRequest, because it does not have chain_name flag
		if method.MethodName == "BridgeChainList" {
			continue
		}

		rpcCommandOptions = append(rpcCommandOptions, &autocliv1.RpcCommandOptions{
			RpcMethod: method.MethodName,
			FlagOptions: map[string]*autocliv1.FlagOptions{
				"chain_name": {
					DefaultValue: am.ModuleName,
				},
			},
		})
	}

	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:              crosschainv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions:    rpcCommandOptions,
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: crosschainv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// skipped because deprecated
				{
					RpcMethod: "OracleSetConfirm",
					Skip:      true,
				},
				{
					RpcMethod: "ConfirmBatch",
					Skip:      true,
				},
				{
					RpcMethod: "BridgeCallConfirm",
					Skip:      true,
				},

				// skipped because authority gated
				{
					RpcMethod: "UpdateParams",
					Skip:      true,
				},
				{
					RpcMethod: "UpdateChainOracles",
					Skip:      true,
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}
