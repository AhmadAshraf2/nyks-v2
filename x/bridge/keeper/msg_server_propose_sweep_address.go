package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) ProposeSweepAddress(goCtx context.Context, msg *types.MsgProposeSweepAddress) (*types.MsgProposeSweepAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	btcAddr, e1 := types.NewBtcAddress(msg.BtcAddress)
	if e1 != nil {
		return nil, types.ErrInvalid.Wrap(e1.Error())
	}

	_, foundDuplicate := k.GetProposeSweepAddress(ctx, msg.ReserveId, msg.RoundId)
	if foundDuplicate {
		return nil, types.ErrDuplicate.Wrap("A similar proposeSweepAddress already exists!")
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	lockStatus := k.IsProposeSweepAddressLocked(ctx)
	if lockStatus {
		return nil, types.ErrProposeSweepAddressIsLocked.Wrap("propose sweep address is locked")
	}
	k.LockProposeSweepAddress(ctx)

	_, errRes := k.VoltKeeper.GetBtcReserve(ctx, msg.ReserveId)
	if errRes != nil {
		return nil, fmt.Errorf("btc reserve not found: %d", msg.ReserveId)
	}

	errSet := k.SetProposeSweepAddress(ctx, *btcAddr, msg.BtcScript, msg.ReserveId, msg.RoundId, judgeAddress)
	if errSet != nil {
		return nil, fmt.Errorf("could not set propose sweep address: %w", errSet)
	}

	k.VoltKeeper.SetNewSweepProposalReceived(ctx, msg.ReserveId, msg.RoundId)

	ctx.EventManager().EmitTypedEvent(
		&types.EventProposeSweepAddress{
			Message:      "MsgProposeSweepAddress",
			BtcAddress:   msg.BtcAddress,
			BtcScript:    msg.BtcScript,
			ReserveId:    msg.ReserveId,
			JudgeAddress: msg.JudgeAddress,
		},
	)

	return &types.MsgProposeSweepAddressResponse{}, nil
}
