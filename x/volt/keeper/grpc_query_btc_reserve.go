package keeper

import (
	"context"

	"twilight-project/nyks/x/volt/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) BtcReserve(goCtx context.Context, req *types.QueryBtcReserveRequest) (*types.QueryBtcReserveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var reserves []types.BtcReserve
	q.k.IterateBtcReserves(goCtx, func(_ []byte, res types.BtcReserve) bool {
		reserves = append(reserves, res)
		return false
	})

	return &types.QueryBtcReserveResponse{BtcReserves: reserves}, nil
}
