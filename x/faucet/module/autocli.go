package blog

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "blog/api/blog/faucet"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "ListMinted",
					Use:            "list-minted ", //address
					Short:          "Query ListMinted Request",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},

				//{
				//	RpcMethod:      "ListPost",
				//	Use:            "list-post",
				//	Short:          "Query list-post",
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				//},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "Mint",
					Use:            "mint [amount] ",
					Short:          "Request Mint token tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				//{
				//	RpcMethod:      "UpdatePost",
				//	Use:            "update-post [title] [body] [id]",
				//	Short:          "Send a update-post tx",
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "body"}, {ProtoField: "id"}},
				//},
				//{
				//	RpcMethod:      "DeletePost",
				//	Use:            "delete-post [id]",
				//	Short:          "Send a delete-post tx",
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				//},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
