package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"twilight-project/nyks/x/volt/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ClearingAccount(goCtx context.Context, req *types.QueryClearingAccountRequest) (*types.QueryClearingAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	twilightAddr, err := sdk.AccAddressFromBech32(req.TwilightAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid twilight address: %s", err)
	}
	account, found := q.k.GetClearingAccount(goCtx, twilightAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "clearing account not found")
	}
	return &types.QueryClearingAccountResponse{ClearingAccount: *account}, nil
}

func (q queryServer) ReserveClearingAccountsAll(goCtx context.Context, req *types.QueryReserveClearingAccountsAllRequest) (*types.QueryReserveClearingAccountsAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	accounts, _ := q.k.GetAllClearingAccountsInaReserve(goCtx, req.ReserveId)
	return &types.QueryReserveClearingAccountsAllResponse{ReserveClearingAccountsAll: accounts}, nil
}

func (q queryServer) ReserveWithdrawSnapshot(goCtx context.Context, req *types.QueryReserveWithdrawSnapshotRequest) (*types.QueryReserveWithdrawSnapshotResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	snapshot, found := q.k.GetReserveWithdrawSnapshot(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "snapshot not found")
	}
	return &types.QueryReserveWithdrawSnapshotResponse{ReserveWithdrawSnapshot: *snapshot}, nil
}

func (q queryServer) RefundTxSnapshot(goCtx context.Context, req *types.QueryRefundTxSnapshotRequest) (*types.QueryRefundTxSnapshotResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	snapshot, found := q.k.GetRefundTxSnapshot(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "refund tx snapshot not found")
	}
	return &types.QueryRefundTxSnapshotResponse{RefundTxSnapshot: *snapshot}, nil
}

func (q queryServer) BtcWithdrawRequest(goCtx context.Context, req *types.QueryBtcWithdrawRequestRequest) (*types.QueryBtcWithdrawRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	twilightAddr, err := sdk.AccAddressFromBech32(req.TwilightAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid twilight address: %s", err)
	}
	withdrawRequest, found := q.k.GetBtcWithdrawRequest(goCtx, twilightAddr, req.ReserveId, req.BtcAddress, req.WithdrawAmount)
	if !found {
		return nil, status.Error(codes.NotFound, "withdraw request not found")
	}
	return &types.QueryBtcWithdrawRequestResponse{BtcWithdrawRequest: *withdrawRequest}, nil
}

func (q queryServer) ReserveWithdrawPool(goCtx context.Context, req *types.QueryReserveWithdrawPoolRequest) (*types.QueryReserveWithdrawPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	pool, found := q.k.GetReserveWithdrawPool(goCtx, req.ReserveId)
	if !found {
		return nil, status.Error(codes.NotFound, "reserve withdraw pool not found")
	}
	return &types.QueryReserveWithdrawPoolResponse{ReserveWithdrawPool: *pool}, nil
}

func (q queryServer) FragmentById(goCtx context.Context, req *types.QueryFragmentByIdRequest) (*types.QueryFragmentByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	fragment, found := q.k.GetFragment(goCtx, req.FragmentId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "fragment %d not found", req.FragmentId)
	}
	return &types.QueryFragmentByIdResponse{Fragment: *fragment}, nil
}

func (q queryServer) GetAllFragments(goCtx context.Context, req *types.QueryGetAllFragmentsRequest) (*types.QueryGetAllFragmentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var fragments []types.Fragment
	q.k.IterateFragments(goCtx, func(_ []byte, res types.Fragment) bool {
		fragments = append(fragments, res)
		return false
	})
	return &types.QueryGetAllFragmentsResponse{Fragments: fragments}, nil
}

func (q queryServer) SignerApplications(goCtx context.Context, req *types.QuerySignerApplicationsRequest) (*types.QuerySignerApplicationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	applications, found := q.k.GetSignerApplications(goCtx, req.FragmentId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no applications found for fragment %d", req.FragmentId)
	}
	return &types.QuerySignerApplicationsResponse{SignerApplications: applications}, nil
}
