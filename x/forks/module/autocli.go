package forks

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"twilight-project/nyks/x/forks/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Params", Use: "params", Short: "Shows the parameters of the module"},
				{RpcMethod: "GetAttestations", Use: "get-attestations", Short: "Get attestations with optional filters"},
				{
					RpcMethod:      "DelegateKeysByBtcOracleAddress",
					Use:            "delegate-keys-by-btc-oracle-address [btc-oracle-address]",
					Short:          "Query delegate keys by BTC oracle address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "btcOracleAddress"}},
				},
				{RpcMethod: "DelegateKeysAll", Use: "delegate-keys-all", Short: "Query all delegate keys"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "UpdateParams", Skip: true},
				{
					RpcMethod: "SetDelegateAddresses",
					Use:       "set-delegate-addresses [validator-address] [btc-oracle-address] [btc-public-key] [zk-oracle-address]",
					Short:     "Set delegate addresses",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validatorAddress"}, {ProtoField: "btcOracleAddress"}, {ProtoField: "btcPublicKey"}, {ProtoField: "zkOracleAddress"},
					},
				},
				{
					RpcMethod: "SeenBtcChainTip",
					Use:       "seen-btc-chain-tip [height] [hash] [btc-oracle-address]",
					Short:     "Report a seen BTC chain tip",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "height"}, {ProtoField: "hash"}, {ProtoField: "btcOracleAddress"},
					},
				},
			},
		},
	}
}
