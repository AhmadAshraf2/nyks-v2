package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateBtcDepositAddress{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposeSweepAddress{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBroadcastTxRefund{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnsignedTxRefund{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnsignedTxSweep{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConfirmBtcWithdraw{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposeRefundHash{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSignSweep{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBroadcastTxSweep{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSignRefund{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawTxFinal{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawTxSigned{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSweepProposal{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawBtcRequest{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBootstrapFragment{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterReserveAddress{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterBtcDepositAddress{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConfirmBtcDeposit{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
