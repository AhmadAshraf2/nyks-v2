package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"twilight-project/nyks/x/forks/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) DelegateKeysByBtcOracleAddress(goCtx context.Context, req *types.QueryDelegateKeysByBtcOracleAddressRequest) (*types.QueryDelegateKeysByBtcOracleAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	accAddr, err := sdk.AccAddressFromBech32(req.BtcOracleAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid btc oracle address: %s", err.Error())
	}

	delegateAddresses, found := q.k.GetDelegateAddresses(ctx, accAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "delegate addresses not found for given btc oracle address")
	}

	return &types.QueryDelegateKeysByBtcOracleAddressResponse{
		Addresses: *delegateAddresses,
	}, nil
}

func (q queryServer) DelegateKeysAll(goCtx context.Context, req *types.QueryDelegateKeysAllRequest) (*types.QueryDelegateKeysAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	allAddresses, err := q.k.GetAllDelegateAddresses(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all delegate addresses: %s", err.Error())
	}

	return &types.QueryDelegateKeysAllResponse{
		Addresses: allAddresses,
	}, nil
}
