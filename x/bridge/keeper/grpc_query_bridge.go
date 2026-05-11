package keeper

import (
	"context"
	"fmt"
	"sort"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"twilight-project/nyks/x/bridge/types"
)

const QUERY_LIMIT uint64 = 1000

// RegisteredBtcDepositAddresses returns all registered BTC deposit addresses
func (q queryServer) RegisteredBtcDepositAddresses(goCtx context.Context, req *types.QueryRegisteredBtcDepositAddressesRequest) (*types.QueryRegisteredBtcDepositAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	addresses := q.k.VoltKeeper.GetAllBtcRegisteredDepositAddresses(goCtx)
	return &types.QueryRegisteredBtcDepositAddressesResponse{Addresses: addresses}, nil
}

// RegisteredBtcDepositAddress returns a single deposit address by BTC address
func (q queryServer) RegisteredBtcDepositAddress(goCtx context.Context, req *types.QueryRegisteredBtcDepositAddressRequest) (*types.QueryRegisteredBtcDepositAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	allAddresses := q.k.VoltKeeper.GetAllBtcRegisteredDepositAddresses(goCtx)
	for _, addr := range allAddresses {
		if addr.BtcDepositAddress == req.DepositAddress {
			return &types.QueryRegisteredBtcDepositAddressResponse{
				DepositAddress:         addr.BtcDepositAddress,
				TwilightDepositAddress: addr.TwilightAddress,
			}, nil
		}
	}
	return nil, status.Error(codes.NotFound, "deposit address not found")
}

// RegisteredReserveAddresses returns all registered reserve addresses
func (q queryServer) RegisteredReserveAddresses(goCtx context.Context, req *types.QueryRegisteredReserveAddressesRequest) (*types.QueryRegisteredReserveAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var addresses []types.MsgRegisterReserveAddress
	q.k.IterateBtcReserveAddresses(goCtx, func(_ []byte, res types.MsgRegisterReserveAddress) bool {
		addresses = append(addresses, res)
		return false
	})
	return &types.QueryRegisteredReserveAddressesResponse{Addresses: addresses}, nil
}

// RegisteredBtcDepositAddressByTwilightAddress returns deposit address by twilight address
func (q queryServer) RegisteredBtcDepositAddressByTwilightAddress(goCtx context.Context, req *types.QueryRegisteredBtcDepositAddressByTwilightAddressRequest) (*types.QueryRegisteredBtcDepositAddressByTwilightAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	twilightAddr, err := sdk.AccAddressFromBech32(req.TwilightDepositAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid twilight address: %s", err)
	}
	deposit, found := q.k.VoltKeeper.GetBtcDepositAddressByTwilightAddress(goCtx, twilightAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "deposit address not found")
	}
	return &types.QueryRegisteredBtcDepositAddressByTwilightAddressResponse{
		DepositAddress:         deposit.BtcDepositAddress,
		TwilightDepositAddress: req.TwilightDepositAddress,
	}, nil
}

// RegisteredJudges returns all registered judges
func (q queryServer) RegisteredJudges(goCtx context.Context, req *types.QueryRegisteredJudgesRequest) (*types.QueryRegisteredJudgesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var judges []types.MsgBootstrapFragment
	q.k.IterateRegisteredJudges(goCtx, func(_ []byte, res types.MsgBootstrapFragment) bool {
		judges = append(judges, res)
		return false
	})
	return &types.QueryRegisteredJudgesResponse{Judges: judges}, nil
}

// RegisteredJudgeAddressByValidatorAddress returns judge address by validator address
func (q queryServer) RegisteredJudgeAddressByValidatorAddress(goCtx context.Context, req *types.QueryRegisteredJudgeAddressByValidatorAddressRequest) (*types.QueryRegisteredJudgeAddressByValidatorAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %s", err)
	}
	judgeAddr, err := q.k.GetJudgeAddressForValidatorAddress(goCtx, valAddr)
	if err != nil {
		return nil, status.Error(codes.NotFound, "judge address not found")
	}
	return &types.QueryRegisteredJudgeAddressByValidatorAddressResponse{
		JudgeAddress:     judgeAddr.String(),
		ValidatorAddress: req.ValidatorAddress,
	}, nil
}

// WithdrawBtcRequestAll returns all withdraw requests
func (q queryServer) WithdrawBtcRequestAll(goCtx context.Context, req *types.QueryWithdrawBtcRequestAllRequest) (*types.QueryWithdrawBtcRequestAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := q.k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.BtcWithdrawRequestKey)
	defer iterator.Close()

	var _ = fmt.Sprint // keep import
	// WithdrawRequests are stored in volt module, iterate via volt keeper prefix
	// For now return empty - needs volt keeper method to iterate withdraw requests
	return &types.QueryWithdrawBtcRequestAllResponse{}, nil
}

// ProposeSweepAddress returns propose sweep address by reserveId/roundId
func (q queryServer) ProposeSweepAddress(goCtx context.Context, req *types.QueryProposeSweepAddressRequest) (*types.QueryProposeSweepAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msg, found := q.k.GetProposeSweepAddress(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "propose sweep address not found")
	}
	return &types.QueryProposeSweepAddressResponse{ProposeSweepAddressMsg: *msg}, nil
}

// ProposeSweepAddressesAll returns all proposed sweep addresses with limit
func (q queryServer) ProposeSweepAddressesAll(goCtx context.Context, req *types.QueryProposeSweepAddressesAllRequest) (*types.QueryProposeSweepAddressesAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	limit := req.Limit
	if limit == 0 || limit > QUERY_LIMIT {
		limit = QUERY_LIMIT
	}
	msgs, err := q.k.GetAllProposedSweepAddresses(goCtx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed: %s", err)
	}
	return &types.QueryProposeSweepAddressesAllResponse{ProposeSweepAddressMsgs: msgs}, nil
}

// UnsignedTxSweep returns unsigned sweep tx by reserveId/roundId
func (q queryServer) UnsignedTxSweep(goCtx context.Context, req *types.QueryUnsignedTxSweepRequest) (*types.QueryUnsignedTxSweepResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msg, found := q.k.GetUnsignedTxSweepMsg(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "unsigned tx sweep not found")
	}
	return &types.QueryUnsignedTxSweepResponse{UnsignedTxSweepMsg: *msg}, nil
}

// UnsignedTxSweepAll returns all unsigned sweep txs with limit
func (q queryServer) UnsignedTxSweepAll(goCtx context.Context, req *types.QueryUnsignedTxSweepAllRequest) (*types.QueryUnsignedTxSweepAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	limit := req.Limit
	if limit == 0 || limit > QUERY_LIMIT {
		limit = QUERY_LIMIT
	}
	store := q.k.KVStore(goCtx)
	prefix := types.UnsignedTxSweepMsgKey
	iter := store.ReverseIterator(prefixRange(prefix))
	defer iter.Close()
	var msgs []types.MsgUnsignedTxSweep
	var count uint64
	for ; iter.Valid() && count < limit; iter.Next() {
		var msg types.MsgUnsignedTxSweep
		q.k.cdc.MustUnmarshal(iter.Value(), &msg)
		msgs = append(msgs, msg)
		count++
	}
	sort.Slice(msgs, func(i, j int) bool {
		if msgs[i].ReserveId != msgs[j].ReserveId {
			return msgs[i].ReserveId > msgs[j].ReserveId
		}
		return msgs[i].RoundId > msgs[j].RoundId
	})
	return &types.QueryUnsignedTxSweepAllResponse{UnsignedTxSweepMsgs: msgs}, nil
}

// UnsignedTxRefund returns unsigned refund tx by reserveId/roundId
func (q queryServer) UnsignedTxRefund(goCtx context.Context, req *types.QueryUnsignedTxRefundRequest) (*types.QueryUnsignedTxRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msg, found := q.k.GetUnsignedTxRefundMsg(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "unsigned tx refund not found")
	}
	return &types.QueryUnsignedTxRefundResponse{UnsignedTxRefundMsg: *msg}, nil
}

// UnsignedTxRefundAll returns all unsigned refund txs with limit
func (q queryServer) UnsignedTxRefundAll(goCtx context.Context, req *types.QueryUnsignedTxRefundAllRequest) (*types.QueryUnsignedTxRefundAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	limit := req.Limit
	if limit == 0 || limit > QUERY_LIMIT {
		limit = QUERY_LIMIT
	}
	store := q.k.KVStore(goCtx)
	prefix := types.UnsignedTxRefundMsgKey
	iter := store.ReverseIterator(prefixRange(prefix))
	defer iter.Close()
	var msgs []types.MsgUnsignedTxRefund
	var count uint64
	for ; iter.Valid() && count < limit; iter.Next() {
		var msg types.MsgUnsignedTxRefund
		q.k.cdc.MustUnmarshal(iter.Value(), &msg)
		msgs = append(msgs, msg)
		count++
	}
	sort.Slice(msgs, func(i, j int) bool {
		if msgs[i].ReserveId != msgs[j].ReserveId {
			return msgs[i].ReserveId > msgs[j].ReserveId
		}
		return msgs[i].RoundId > msgs[j].RoundId
	})
	return &types.QueryUnsignedTxRefundAllResponse{UnsignedTxRefundMsgs: msgs}, nil
}

// SignRefund returns sign refund by reserveId/roundId
func (q queryServer) SignRefund(goCtx context.Context, req *types.QuerySignRefundRequest) (*types.QuerySignRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msgs, found := q.k.GetBtcSignRefundMsg(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "sign refund not found")
	}
	return &types.QuerySignRefundResponse{SignRefundMsg: msgs}, nil
}

// SignRefundAll returns all sign refund messages
func (q queryServer) SignRefundAll(goCtx context.Context, req *types.QuerySignRefundAllRequest) (*types.QuerySignRefundAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var msgs []types.MsgSignRefund
	q.k.IterateRegisteredSignRefundMsgs(goCtx, func(_ []byte, res types.MsgSignRefund) bool {
		msgs = append(msgs, res)
		return false
	})
	return &types.QuerySignRefundAllResponse{SignRefundMsg: msgs}, nil
}

// SignSweep returns sign sweep by reserveId/roundId
func (q queryServer) SignSweep(goCtx context.Context, req *types.QuerySignSweepRequest) (*types.QuerySignSweepResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msgs, found := q.k.GetBtcSignSweepMsg(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "sign sweep not found")
	}
	return &types.QuerySignSweepResponse{SignSweepMsg: msgs}, nil
}

// SignSweepAll returns all sign sweep messages
func (q queryServer) SignSweepAll(goCtx context.Context, req *types.QuerySignSweepAllRequest) (*types.QuerySignSweepAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var msgs []types.MsgSignSweep
	q.k.IterateRegisteredSignSweepMsgs(goCtx, func(_ []byte, res types.MsgSignSweep) bool {
		msgs = append(msgs, res)
		return false
	})
	return &types.QuerySignSweepAllResponse{SignSweepMsg: msgs}, nil
}

// BroadcastTxSweep returns broadcast sweep tx by reserveId/roundId
func (q queryServer) BroadcastTxSweep(goCtx context.Context, req *types.QueryBroadcastTxSweepRequest) (*types.QueryBroadcastTxSweepResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msg, found := q.k.GetBtcBroadcastTxSweepMsg(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "broadcast tx sweep not found")
	}
	return &types.QueryBroadcastTxSweepResponse{BroadcastSweepMsg: *msg}, nil
}

// BroadcastTxSweepAll returns all broadcast sweep messages
func (q queryServer) BroadcastTxSweepAll(goCtx context.Context, req *types.QueryBroadcastTxSweepAllRequest) (*types.QueryBroadcastTxSweepAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var msgs []types.MsgBroadcastTxSweep
	q.k.IterateRegisteredBroadcastTxSweepMsgs(goCtx, func(_ []byte, res types.MsgBroadcastTxSweep) bool {
		msgs = append(msgs, res)
		return false
	})
	return &types.QueryBroadcastTxSweepAllResponse{BroadcastTxSweepMsg: msgs}, nil
}

// BroadcastTxRefund returns broadcast refund tx by reserveId/roundId
func (q queryServer) BroadcastTxRefund(goCtx context.Context, req *types.QueryBroadcastTxRefundRequest) (*types.QueryBroadcastTxRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	msg, found := q.k.GetBtcBroadcastTxRefundMsg(goCtx, req.ReserveId, req.RoundId)
	if !found {
		return nil, status.Error(codes.NotFound, "broadcast tx refund not found")
	}
	return &types.QueryBroadcastTxRefundResponse{BroadcastRefundMsg: *msg}, nil
}

// BroadcastTxRefundAll returns all broadcast refund messages
func (q queryServer) BroadcastTxRefundAll(goCtx context.Context, req *types.QueryBroadcastTxRefundAllRequest) (*types.QueryBroadcastTxRefundAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var msgs []types.MsgBroadcastTxRefund
	q.k.IterateRegisteredBroadcastTxRefundMsgs(goCtx, func(_ []byte, res types.MsgBroadcastTxRefund) bool {
		msgs = append(msgs, res)
		return false
	})
	return &types.QueryBroadcastTxRefundAllResponse{BroadcastTxRefundMsg: msgs}, nil
}

// ProposeRefundHashAll returns all propose refund hash messages
func (q queryServer) ProposeRefundHashAll(goCtx context.Context, req *types.QueryProposeRefundHashAllRequest) (*types.QueryProposeRefundHashAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var msgs []types.MsgProposeRefundHash
	q.k.IterateRegisteredProposeRefundHashMsgs(goCtx, func(_ []byte, res types.MsgProposeRefundHash) bool {
		msgs = append(msgs, res)
		return false
	})
	return &types.QueryProposeRefundHashAllResponse{ProposeRefundHashMsg: msgs}, nil
}

// Iteration methods for keeper that are needed by query handlers

func (k Keeper) IterateRegisteredSignRefundMsgs(ctx context.Context, cb func([]byte, types.MsgSignRefund) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcSignRefundMsgKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgSignRefund
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) IterateRegisteredSignSweepMsgs(ctx context.Context, cb func([]byte, types.MsgSignSweep) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcSignSweepMsgKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgSignSweep
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) IterateRegisteredBroadcastTxSweepMsgs(ctx context.Context, cb func([]byte, types.MsgBroadcastTxSweep) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcBroadcastTxSweepMsgKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgBroadcastTxSweep
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) IterateRegisteredBroadcastTxRefundMsgs(ctx context.Context, cb func([]byte, types.MsgBroadcastTxRefund) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcBroadcastTxRefundMsgKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgBroadcastTxRefund
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) IterateRegisteredProposeRefundHashMsgs(ctx context.Context, cb func([]byte, types.MsgProposeRefundHash) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcProposeRefundHashMsgKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgProposeRefundHash
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}
