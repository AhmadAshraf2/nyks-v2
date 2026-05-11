package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) BroadcastTxRefund(goCtx context.Context, msg *types.MsgBroadcastTxRefund) (*types.MsgBroadcastTxRefundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	_, foundDuplicate := k.GetBtcBroadcastTxRefundMsg(ctx, msg.ReserveId, msg.RoundId)
	if foundDuplicate {
		return nil, types.ErrDuplicate.Wrap("Duplicate broadcast refund request")
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	errSet := k.SetBtcBroadcastTxRefundMsg(ctx, msg.ReserveId, msg.RoundId, judgeAddress, msg.SignedRefundTx)
	if errSet != nil {
		return nil, errSet
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventBroadcastTxRefund{
			Message:        "MsgBroadcastTxRefund",
			ReserveId:      msg.ReserveId,
			RoundId:        msg.RoundId,
			SignedRefundTx: msg.SignedRefundTx,
			JudgeAddress:   msg.JudgeAddress,
		},
	)

	return &types.MsgBroadcastTxRefundResponse{}, nil
}
