package bridge

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"twilight-project/nyks/x/bridge/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Params", Use: "params", Short: "Shows the parameters of the module"},
				{RpcMethod: "RegisteredBtcDepositAddresses", Use: "registered-btc-deposit-addresses", Short: "Query all registered BTC deposit addresses"},
				{RpcMethod: "RegisteredReserveAddresses", Use: "registered-reserve-addresses", Short: "Query all registered reserve addresses"},
				{
					RpcMethod: "RegisteredBtcDepositAddress", Use: "registered-btc-deposit-address [deposit-address]", Short: "Query by deposit address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "depositAddress"}},
				},
				{
					RpcMethod: "RegisteredBtcDepositAddressByTwilightAddress", Use: "registered-btc-deposit-address-by-twilight-address [twilight-address]", Short: "Query by twilight address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "twilightDepositAddress"}},
				},
				{RpcMethod: "RegisteredJudges", Use: "registered-judges", Short: "Query all registered judges"},
				{
					RpcMethod: "RegisteredJudgeAddressByValidatorAddress", Use: "registered-judge-address-by-validator-address [validator-address]", Short: "Query judge by validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "validatorAddress"}},
				},
				{RpcMethod: "WithdrawBtcRequestAll", Use: "withdraw-btc-request-all", Short: "Query all withdraw requests"},
				{
					RpcMethod: "ProposeSweepAddress", Use: "propose-sweep-address [reserve-id] [round-id]", Short: "Query propose sweep address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{
					RpcMethod: "ProposeSweepAddressesAll", Use: "propose-sweep-addresses-all [limit]", Short: "Query all proposed sweep addresses",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "limit"}},
				},
				{
					RpcMethod: "UnsignedTxSweep", Use: "unsigned-tx-sweep [reserve-id] [round-id]", Short: "Query unsigned sweep tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{RpcMethod: "UnsignedTxSweepAll", Use: "unsigned-tx-sweep-all", Short: "Query all unsigned sweep txs"},
				{
					RpcMethod: "UnsignedTxRefund", Use: "unsigned-tx-refund [reserve-id] [round-id]", Short: "Query unsigned refund tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{RpcMethod: "UnsignedTxRefundAll", Use: "unsigned-tx-refund-all", Short: "Query all unsigned refund txs"},
				{
					RpcMethod: "SignRefund", Use: "sign-refund [reserve-id] [round-id]", Short: "Query sign refund",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{RpcMethod: "SignRefundAll", Use: "sign-refund-all", Short: "Query all sign refunds"},
				{
					RpcMethod: "SignSweep", Use: "sign-sweep [reserve-id] [round-id]", Short: "Query sign sweep",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{RpcMethod: "SignSweepAll", Use: "sign-sweep-all", Short: "Query all sign sweeps"},
				{
					RpcMethod: "BroadcastTxSweep", Use: "broadcast-tx-sweep [reserve-id] [round-id]", Short: "Query broadcast sweep tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{RpcMethod: "BroadcastTxSweepAll", Use: "broadcast-tx-sweep-all", Short: "Query all broadcast sweep txs"},
				{
					RpcMethod: "BroadcastTxRefund", Use: "broadcast-tx-refund [reserve-id] [round-id]", Short: "Query broadcast refund tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reserveId"}, {ProtoField: "roundId"}},
				},
				{RpcMethod: "BroadcastTxRefundAll", Use: "broadcast-tx-refund-all", Short: "Query all broadcast refund txs"},
				{RpcMethod: "ProposeRefundHashAll", Use: "propose-refund-hash-all", Short: "Query all proposed refund hashes"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "UpdateParams", Skip: true},
				{
					RpcMethod: "ConfirmBtcDeposit",
					Use:       "msg-confirm-btc-deposit [reserve-address] [deposit-amount] [block-height] [block-hash] [twilight-deposit-address]",
					Short:     "Confirm a BTC deposit",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveAddress"}, {ProtoField: "depositAmount"}, {ProtoField: "height"}, {ProtoField: "hash"}, {ProtoField: "twilightDepositAddress"},
					},
				},
				{
					RpcMethod: "RegisterBtcDepositAddress",
					Use:       "register-deposit-address [btc-deposit-address] [btc-satoshi-test-amount] [twilight-staking-amount]",
					Short:     "Register a BTC deposit address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "btcDepositAddress"}, {ProtoField: "btcSatoshiTestAmount"}, {ProtoField: "twilightStakingAmount"},
					},
				},
				{
					RpcMethod: "RegisterReserveAddress",
					Use:       "register-reserve-address [fragment-id] [reserve-script] [reserve-address] [judge-address]",
					Short:     "Register a reserve address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "fragmentId"}, {ProtoField: "reserveScript"}, {ProtoField: "reserveAddress"}, {ProtoField: "judgeAddress"},
					},
				},
				{
					RpcMethod: "BootstrapFragment",
					Use:       "bootstrap-fragment [judge-address] [num-of-signers] [threshold] [signer-application-fee] [fragment-fee-bips] [arbitrary-data]",
					Short:     "Bootstrap a new fragment",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "judgeAddress"}, {ProtoField: "numOfSigners"}, {ProtoField: "threshold"}, {ProtoField: "signerApplicationFee"}, {ProtoField: "fragmentFeeBips"}, {ProtoField: "arbitraryData"},
					},
				},
				{
					RpcMethod: "WithdrawBtcRequest",
					Use:       "withdraw-btc-request [withdraw-address] [reserve-id] [withdraw-amount]",
					Short:     "Request a BTC withdrawal",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "withdrawAddress"}, {ProtoField: "reserveId"}, {ProtoField: "withdrawAmount"},
					},
				},
				{
					RpcMethod: "SweepProposal",
					Use:       "sweep-proposal [reserve-id] [new-reserve-address] [judge-address] [btc-block-number] [btc-relay-capacity-value] [btc-tx-hash] [unlock-height] [round-id]",
					Short:     "Submit a sweep proposal",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveId"}, {ProtoField: "newReserveAddress"}, {ProtoField: "judgeAddress"}, {ProtoField: "BtcBlockNumber"}, {ProtoField: "btcRelayCapacityValue"}, {ProtoField: "btcTxHash"}, {ProtoField: "UnlockHeight"}, {ProtoField: "roundId"},
					},
				},
				{
					RpcMethod: "WithdrawTxSigned",
					Use:       "withdraw-tx-signed [validator-address] [btc-tx-signed]",
					Short:     "Submit a signed withdraw tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validatorAddress"}, {ProtoField: "btcTxSigned"},
					},
				},
				{
					RpcMethod: "WithdrawTxFinal",
					Use:       "withdraw-tx-final [judge-address] [btc-tx]",
					Short:     "Submit a final withdraw tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "judgeAddress"}, {ProtoField: "btcTx"},
					},
				},
				{
					RpcMethod: "SignRefund",
					Use:       "sign-refund [reserve-id] [round-id] [signer-public-key] [refund-signatures]",
					Short:     "Sign a refund transaction",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveId"}, {ProtoField: "roundId"}, {ProtoField: "signerPublicKey"}, {ProtoField: "refundSignature", Varargs: true},
					},
				},
				{
					RpcMethod: "BroadcastTxSweep",
					Use:       "broadcast-tx-sweep [reserve-id] [round-id] [signed-sweep-tx]",
					Short:     "Broadcast a signed sweep tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveId"}, {ProtoField: "roundId"}, {ProtoField: "signedSweepTx"},
					},
				},
				{
					RpcMethod: "SignSweep",
					Use:       "sign-sweep [reserve-id] [round-id] [signer-public-key] [sweep-signatures]",
					Short:     "Sign a sweep transaction",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveId"}, {ProtoField: "roundId"}, {ProtoField: "signerPublicKey"}, {ProtoField: "sweepSignature", Varargs: true},
					},
				},
				{
					RpcMethod: "ProposeRefundHash",
					Use:       "propose-refund-hash [refund-hash]",
					Short:     "Propose a refund hash",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "refundHash"},
					},
				},
				{
					RpcMethod: "ConfirmBtcWithdraw",
					Use:       "confirm-btc-withdraw [tx-hash] [height] [hash]",
					Short:     "Confirm a BTC withdrawal",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "txHash"}, {ProtoField: "height"}, {ProtoField: "hash"},
					},
				},
				{
					RpcMethod: "UnsignedTxSweep",
					Use:       "unsigned-tx-sweep [tx-id] [btc-unsigned-sweep-tx] [reserve-id] [round-id]",
					Short:     "Submit an unsigned sweep tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "txId"}, {ProtoField: "btcUnsignedSweepTx"}, {ProtoField: "reserveId"}, {ProtoField: "roundId"},
					},
				},
				{
					RpcMethod: "UnsignedTxRefund",
					Use:       "unsigned-tx-refund [reserve-id] [round-id] [btc-unsigned-refund-tx]",
					Short:     "Submit an unsigned refund tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveId"}, {ProtoField: "roundId"}, {ProtoField: "btcUnsignedRefundTx"},
					},
				},
				{
					RpcMethod: "BroadcastTxRefund",
					Use:       "broadcast-tx-refund [reserve-id] [round-id] [signed-refund-tx]",
					Short:     "Broadcast a signed refund tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "reserveId"}, {ProtoField: "roundId"}, {ProtoField: "signedRefundTx"},
					},
				},
				{
					RpcMethod: "ProposeSweepAddress",
					Use:       "propose-sweep-address [btc-address] [btc-script] [reserve-id] [round-id]",
					Short:     "Propose a sweep address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "btcAddress"}, {ProtoField: "btcScript"}, {ProtoField: "reserveId"}, {ProtoField: "roundId"},
					},
				},
				{
					RpcMethod: "UpdateBtcDepositAddress",
					Use:       "update-btc-deposit-address [btc-deposit-address] [btc-satoshi-test-amount] [twilight-staking-amount] [twilight-address]",
					Short:     "Update a BTC deposit address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "btcDepositAddress"}, {ProtoField: "btcSatoshiTestAmount"}, {ProtoField: "twilightStakingAmount"}, {ProtoField: "twilightAddress"},
					},
				},
			},
		},
	}
}
