package keeper

import (
	"context"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) WithdrawTxFinal(_ context.Context, msg *types.MsgWithdrawTxFinal) (*types.MsgWithdrawTxFinalResponse, error) {
	// TODO: Handling the message
	_ = msg
	return &types.MsgWithdrawTxFinalResponse{}, nil
}
