package keeper

import (
	"context"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) WithdrawTxSigned(_ context.Context, msg *types.MsgWithdrawTxSigned) (*types.MsgWithdrawTxSignedResponse, error) {
	// TODO: Handling the message
	_ = msg
	return &types.MsgWithdrawTxSignedResponse{}, nil
}
