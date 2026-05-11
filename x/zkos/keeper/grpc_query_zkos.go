package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	forkstypes "twilight-project/nyks/x/forks/types"
	"twilight-project/nyks/x/zkos/types"
)

func (q queryServer) TransferTx(goCtx context.Context, req *types.QueryTransferTxRequest) (*types.QueryTransferTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ttx, found := q.k.GetTransferTx(goCtx, req.TxId)
	if !found {
		return nil, status.Error(codes.NotFound, "transfer tx not found")
	}
	return &types.QueryTransferTxResponse{TransferTx: ttx}, nil
}

func (q queryServer) MintOrBurnTradingBtc(goCtx context.Context, req *types.QueryMintOrBurnTradingBtcRequest) (*types.QueryMintOrBurnTradingBtcResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := q.k.KVStore(ctx)

	prefix := forkstypes.AppendBytes(types.KeyMintOrBurnTradingBtc, []byte(req.TwilightAddress))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var results []types.MsgMintBurnTradingBtc
	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgMintBurnTradingBtc
		q.k.cdc.MustUnmarshal(iterator.Value(), &msg)
		results = append(results, msg)
	}

	if len(results) == 0 {
		return nil, status.Error(codes.NotFound, "mint or burn trading btc not found")
	}

	return &types.QueryMintOrBurnTradingBtcResponse{MintOrBurnTradingBtc: results}, nil
}
