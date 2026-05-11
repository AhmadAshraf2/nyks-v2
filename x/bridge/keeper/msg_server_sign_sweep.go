package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) SignSweep(goCtx context.Context, msg *types.MsgSignSweep) (*types.MsgSignSweepResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signerAddress, e1 := sdk.AccAddressFromBech32(msg.SignerAddress)
	if e1 != nil {
		return nil, types.ErrInvalid.Wrap(e1.Error())
	}

	sweepSigValid := types.ValidateSignatures(msg.SweepSignature)
	if !sweepSigValid {
		return nil, types.ErrInvalid.Wrap("invalid sweep signature")
	}

	_, found := k.GetBtcSignSweepMsgWithOracleAddress(ctx, msg.ReserveId, msg.RoundId, signerAddress)
	if found {
		return nil, types.ErrDuplicate.Wrap("Duplicate sweep Request")
	}

	err := k.SetBtcSignSweepMsg(ctx, signerAddress, msg.ReserveId, msg.RoundId, msg.SignerPublicKey, msg.SweepSignature)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventSignSweep{
			Message:         "MsgSignSweep",
			ReserveId:       msg.ReserveId,
			RoundId:         msg.RoundId,
			SignerPublicKey: msg.SignerPublicKey,
			SweepSignature:  msg.SweepSignature,
			SignerAddress:   msg.SignerAddress,
		},
	)

	return &types.MsgSignSweepResponse{}, nil
}
