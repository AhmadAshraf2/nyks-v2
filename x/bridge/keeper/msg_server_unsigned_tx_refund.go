package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) UnsignedTxRefund(goCtx context.Context, msg *types.MsgUnsignedTxRefund) (*types.MsgUnsignedTxRefundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	_, foundDuplicate := k.GetUnsignedTxRefundMsg(ctx, msg.ReserveId, msg.RoundId)
	if foundDuplicate {
		return nil, types.ErrDuplicate.Wrap("A similar unsignedTxRefund already exists!")
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	_, errRes := k.VoltKeeper.GetBtcReserve(ctx, msg.ReserveId)
	if errRes != nil {
		return nil, fmt.Errorf("btc reserve not found: %d", msg.ReserveId)
	}

	errSet := k.SetUnsignedTxRefundMsg(ctx, msg.ReserveId, msg.RoundId, msg.BtcUnsignedRefundTx, judgeAddress)
	if errSet != nil {
		return nil, fmt.Errorf("could not set the transaction refund: %w", errSet)
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventUnsignedTxRefund{
			Message:          "MsgUnsignedTxRefund",
			ReserveId:        msg.ReserveId,
			RoundId:          msg.RoundId,
			UnsignedRefundTx: msg.BtcUnsignedRefundTx,
			JudgeAddress:     msg.JudgeAddress,
		},
	)

	return &types.MsgUnsignedTxRefundResponse{}, nil
}
