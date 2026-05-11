package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/zkos/types"
)

func (k msgServer) TransferTx(goCtx context.Context, msg *types.MsgTransferTx) (*types.MsgTransferTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.SetTransferTx(ctx, msg.TxId, msg.TxByteCode, msg.TxFee, msg.ZkOracleAddress)
	k.DeductFeeFromPrivatePool(ctx, msg.TxFee)

	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferTx{
			Message:         "MsgTransferTx",
			TxId:            msg.TxId,
			ZkOracleAddress: msg.ZkOracleAddress,
		},
	)

	return &types.MsgTransferTxResponse{}, nil
}
