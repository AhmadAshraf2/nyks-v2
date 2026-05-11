package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) SignRefund(goCtx context.Context, msg *types.MsgSignRefund) (*types.MsgSignRefundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signerAddress, e1 := sdk.AccAddressFromBech32(msg.SignerAddress)
	if e1 != nil {
		return nil, types.ErrInvalid.Wrap(e1.Error())
	}

	refundSigValid := types.ValidateSignatures(msg.RefundSignature)
	if !refundSigValid {
		return nil, types.ErrInvalid.Wrap("invalid refund signature")
	}

	_, found := k.GetBtcSignRefundMsgWithOracleAddress(ctx, msg.ReserveId, msg.RoundId, signerAddress)
	if found {
		return nil, types.ErrDuplicate.Wrap("Duplicate Refund Request")
	}

	err := k.SetBtcSignRefundMsg(ctx, signerAddress, msg.ReserveId, msg.RoundId, msg.SignerPublicKey, msg.RefundSignature)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventSignRefund{
			Message:         "MsgSignRefund",
			ReserveId:       msg.ReserveId,
			RoundId:         msg.RoundId,
			SignerPublicKey: msg.SignerPublicKey,
			RefundSignature: msg.RefundSignature,
			SignerAddress:   msg.SignerAddress,
		},
	)

	return &types.MsgSignRefundResponse{}, nil
}
