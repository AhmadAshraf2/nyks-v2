package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) ProposeRefundHash(goCtx context.Context, msg *types.MsgProposeRefundHash) (*types.MsgProposeRefundHashResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, _ := sdk.AccAddressFromBech32(msg.JudgeAddress)

	_, found := k.GetBtcProposeRefundHashMsg(ctx, judgeAddress, msg.RefundHash)
	if found {
		return nil, types.ErrDuplicate.Wrap("Duplicate propose refund hash request")
	}

	err := k.SetBtcProposeRefundHashMsg(ctx, judgeAddress, msg.RefundHash)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventProposeRefundHash{
			Message:      "MsgProposeRefundHash",
			RefundHash:   msg.RefundHash,
			JudgeAddress: msg.JudgeAddress,
		},
	)

	return &types.MsgProposeRefundHashResponse{}, nil
}

// GetBtcProposeRefundHashMsg returns propose refund hash message
func (k Keeper) GetBtcProposeRefundHashMsg(ctx context.Context, judgeAddress sdk.AccAddress, refundHash string) (*types.MsgProposeRefundHash, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetBtcProposeRefundHashMsgKey(judgeAddress, refundHash)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var msg types.MsgProposeRefundHash
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}
