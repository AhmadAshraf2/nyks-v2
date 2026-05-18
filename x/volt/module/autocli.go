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
				{RpcMethod: "BtcReserve", Use: "btc-reserve", Short: "Query all BTC reserves"},
				{
					RpcMethod:      "ClearingAccount",
					Use:            "clearing-account [twilight-address]",
					Short:          "Query clearing account by twilight address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "twilightAddress"}},
				},
				{
					RpcMethod:      "ReserveClearingAccountsAll",
					Use:            "reserve-clearing-accounts-all [reserve-id]",
					Short:          "Query all clearing accounts in a reserve",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}},
				},
				{
					RpcMethod:      "ReserveWithdrawSnapshot",
					Use:            "reserve-withdraw-snapshot [reserve-id] [round-id]",
					Short:          "Query reserve withdraw snapshot",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{
					RpcMethod:      "RefundTxSnapshot",
					Use:            "refund-tx-snapshot [reserve-id] [round-id]",
					Short:          "Query refund tx snapshot",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{
					RpcMethod:      "BtcWithdrawRequest",
					Use:            "btc-withdraw-request [twilight-address] [reserve-id] [btc-address] [withdraw-amount]",
					Short:          "Query BTC withdraw request",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "twilightAddress"}, {ProtoField: "reserveId"}, {ProtoField: "btcAddress"}, {ProtoField: "withdrawAmount"}},
				},
				{
					RpcMethod:      "ReserveWithdrawPool",
					Use:            "reserve-withdraw-pool [reserve-id]",
					Short:          "Query reserve withdraw pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}},
				},
				{
					RpcMethod:      "FragmentById",
					Use:            "fragment-by-id [fragment-id]",
					Short:          "Query fragment by ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "fragmentId"}},
				},
				{RpcMethod: "GetAllFragments", Use: "get-all-fragments", Short: "Query all fragments"},
				{
					RpcMethod:      "SignerApplications",
					Use:            "signer-applications [fragment-id]",
					Short:          "Query signer applications for a fragment",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "fragmentId"}},
				},
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
