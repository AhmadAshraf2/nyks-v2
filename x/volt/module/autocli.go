package volt

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"twilight-project/nyks/x/volt/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Params", Use: "params", Short: "Shows the parameters of the module"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "UpdateParams", Skip: true},
				{
					RpcMethod: "SignerApplication",
					Use:       "signer-application [fragment-id] [application-fee] [fee-bips] [btc-pub-key]",
					Short:     "Submit a signer application",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "fragmentId"}, {ProtoField: "applicationFee"}, {ProtoField: "feeBips"}, {ProtoField: "btcPubKey"},
					},
				},
				{
					RpcMethod: "AcceptSigners",
					Use:       "accept-signers [fragment-id] [signer-application-ids]",
					Short:     "Accept signers into a fragment",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "fragmentId"}, {ProtoField: "signerApplicationIds", Varargs: true},
					},
				},
			},
		},
	}
}
