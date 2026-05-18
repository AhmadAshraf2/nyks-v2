package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSeenBtcChainTip{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetDelegateAddresses{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	// Register BtcProposal interface implementations for attestation unpacking
	registrar.RegisterInterface("twilightproject.nyks.forks.BtcProposal", (*BtcProposal)(nil))
	registrar.RegisterImplementations((*BtcProposal)(nil),
		&MsgSeenBtcChainTip{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
