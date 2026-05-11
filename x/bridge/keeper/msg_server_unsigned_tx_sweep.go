package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) UnsignedTxSweep(goCtx context.Context, msg *types.MsgUnsignedTxSweep) (*types.MsgUnsignedTxSweepResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	_, foundDuplicate := k.GetUnsignedTxSweepMsg(ctx, msg.ReserveId, msg.RoundId)
	if foundDuplicate {
		return nil, types.ErrDuplicate.Wrap("A similar unsignedTxSweep already exists!")
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	_, errRes := k.VoltKeeper.GetBtcReserve(ctx, msg.ReserveId)
	if errRes != nil {
		return nil, fmt.Errorf("btc reserve not found: %d", msg.ReserveId)
	}

	errSet := k.SetUnsignedTxSweepMsg(ctx, msg.TxId, msg.BtcUnsignedSweepTx, msg.ReserveId, msg.RoundId, judgeAddress)
	if errSet != nil {
		return nil, fmt.Errorf("could not set the transaction sweep: %w", errSet)
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventUnsignedTxSweep{
			Message:         "MsgUnsignedTxSweep",
			TxId:            msg.TxId,
			ReserveId:       msg.ReserveId,
			RoundId:         msg.RoundId,
			UnsignedSweepTx: msg.BtcUnsignedSweepTx,
			JudgeAddress:    msg.JudgeAddress,
		},
	)
	return &types.MsgUnsignedTxSweepResponse{}, nil
}
