package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) BroadcastTxSweep(goCtx context.Context, msg *types.MsgBroadcastTxSweep) (*types.MsgBroadcastTxSweepResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	_, foundDuplicate := k.GetBtcBroadcastTxSweepMsg(ctx, msg.ReserveId, msg.RoundId)
	if foundDuplicate {
		return nil, types.ErrDuplicate.Wrap("Duplicate broadcast sweep request")
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	errSet := k.SetBtcBroadcastTxSweepMsg(ctx, msg.ReserveId, msg.RoundId, judgeAddress, msg.SignedSweepTx)
	if errSet != nil {
		return nil, errSet
	}

	// Unlock the ProposeSweepAddress for the next sweep cycle
	k.UnlockProposeSweepAddress(ctx)

	ctx.EventManager().EmitTypedEvent(
		&types.EventBroadcastTxSweep{
			Message:       "MsgBroadcastTxSweep",
			ReserveId:     msg.ReserveId,
			RoundId:       msg.RoundId,
			SignedSweepTx: msg.SignedSweepTx,
			JudgeAddress:  msg.JudgeAddress,
		},
	)

	return &types.MsgBroadcastTxSweepResponse{}, nil
}
