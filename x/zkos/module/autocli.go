package zkos

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"twilight-project/nyks/x/zkos/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Params", Use: "params", Short: "Shows the parameters of the module"},
				{
					RpcMethod:      "TransferTx",
					Use:            "transfer-tx [tx-id]",
					Short:          "Query a transfer tx by ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "txId"}},
				},
				{
					RpcMethod:      "MintOrBurnTradingBtc",
					Use:            "mint-or-burn-trading-btc [twilight-address]",
					Short:          "Query mint/burn trading BTC by twilight address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "twilightAddress"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "UpdateParams", Skip: true},
				{
					RpcMethod: "TransferTx",
					Use:       "transfer-tx [tx-id] [tx-byte-code] [tx-fee]",
					Short:     "Submit a transfer transaction",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "txId"}, {ProtoField: "txByteCode"}, {ProtoField: "txFee"},
					},
				},
				{
					RpcMethod: "MintBurnTradingBtc",
					Use:       "mint-burn-trading-btc [mint-or-burn] [btc-value] [qq-account] [encrypt-scalar]",
					Short:     "Mint or burn trading BTC",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "mintOrBurn"}, {ProtoField: "btcValue"}, {ProtoField: "qqAccount"}, {ProtoField: "encryptScalar"},
					},
				},
			},
		},
	}
}
