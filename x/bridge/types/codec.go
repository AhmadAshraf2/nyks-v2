package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	forkstypes "twilight-project/nyks/x/forks/types"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateBtcDepositAddress{},
		&MsgProposeSweepAddress{},
		&MsgBroadcastTxRefund{},
		&MsgUnsignedTxRefund{},
		&MsgUnsignedTxSweep{},
		&MsgConfirmBtcWithdraw{},
		&MsgProposeRefundHash{},
		&MsgSignSweep{},
		&MsgBroadcastTxSweep{},
		&MsgSignRefund{},
		&MsgWithdrawTxFinal{},
		&MsgWithdrawTxSigned{},
		&MsgSweepProposal{},
		&MsgWithdrawBtcRequest{},
		&MsgBootstrapFragment{},
		&MsgRegisterReserveAddress{},
		&MsgRegisterBtcDepositAddress{},
		&MsgConfirmBtcDeposit{},
		&MsgUpdateParams{},
	)

	// Register BtcProposal implementations for attestation unpacking
	registrar.RegisterImplementations((*forkstypes.BtcProposal)(nil),
		&MsgConfirmBtcDeposit{},
		&MsgSweepProposal{},
		&MsgConfirmBtcWithdraw{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
