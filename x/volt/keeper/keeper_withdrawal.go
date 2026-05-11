package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	forkstypes "twilight-project/nyks/x/forks/types"
	"twilight-project/nyks/x/volt/types"
)

func (k Keeper) SetReserveWithdrawPool(ctx context.Context, withdrawPool *types.ReserveWithdrawPool) error {
	store := k.KVStore(ctx)
	poolKey := types.GetReserveWithdrawPoolKey(withdrawPool.ReserveID)
	bz, err := k.cdc.Marshal(withdrawPool)
	if err != nil {
		return fmt.Errorf("could not marshal withdraw pool: %w", err)
	}
	store.Set(poolKey, bz)
	return nil
}

func (k Keeper) GetReserveWithdrawPool(ctx context.Context, reserveId uint64) (*types.ReserveWithdrawPool, bool) {
	store := k.KVStore(ctx)
	poolKey := types.GetReserveWithdrawPoolKey(reserveId)
	if !store.Has(poolKey) {
		return nil, false
	}
	bz := store.Get(poolKey)
	var pool types.ReserveWithdrawPool
	k.cdc.MustUnmarshal(bz, &pool)
	return &pool, true
}

func (k Keeper) SetBtcWithdrawRequest(ctx context.Context, twilightAddress sdk.AccAddress, reserveId uint64, withdrawAddress string, withdrawAmount uint64) (*uint32, error) {
	store := k.KVStore(ctx)

	withdrawIdentifier := k.IncrementCounter(ctx, WithdrawalCounterKey)

	aKey := types.GetBtcWithdrawRequestKeyInternal(twilightAddress, reserveId, withdrawAddress, withdrawAmount)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	withdrawRequest := &types.BtcWithdrawRequestInternal{
		WithdrawIdentifier:          withdrawIdentifier,
		WithdrawAddress:             withdrawAddress,
		WithdrawReserveId:           reserveId,
		WithdrawAmount:              withdrawAmount,
		TwilightAddress:             twilightAddress.String(),
		IsConfirmed:                 false,
		CreationTwilightBlockHeight: sdkCtx.BlockHeight(),
	}

	store.Set(aKey, k.cdc.MustMarshal(withdrawRequest))

	err := k.AddToReserveWithdrawPool(ctx, reserveId, withdrawIdentifier)
	if err != nil {
		return nil, err
	}

	withdrawCoin := sdk.NewCoin("sats", math.NewIntFromUint64(withdrawAmount))
	if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, twilightAddress, forkstypes.ModuleName, sdk.NewCoins(withdrawCoin)); err != nil {
		return nil, fmt.Errorf("failed to deduct coins from user account %s: %w", twilightAddress, err)
	}

	if err := k.BankKeeper.BurnCoins(ctx, forkstypes.ModuleName, sdk.NewCoins(withdrawCoin)); err != nil {
		return nil, fmt.Errorf("failed to burn coins %v: %w", withdrawCoin, err)
	}

	err = k.DeductFromClearingAccount(ctx, twilightAddress, reserveId, withdrawAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to deduct from ClearingAccount: %w", err)
	}

	return &withdrawIdentifier, nil
}

func (k Keeper) AddToReserveWithdrawPool(ctx context.Context, reserveId uint64, withdrawIdentifier uint32) error {
	store := k.KVStore(ctx)
	poolKey := types.GetReserveWithdrawPoolKey(reserveId)

	var pool types.ReserveWithdrawPool
	if store.Has(poolKey) {
		b := store.Get(poolKey)
		k.cdc.MustUnmarshal(b, &pool)
	} else {
		pool = types.ReserveWithdrawPool{
			ReserveID:                     reserveId,
			RoundID:                       0,
			ProcessingWithdrawIdentifiers: []uint32{},
			QueuedWithdrawIdentifiers:     []uint32{},
			CurrentProcessingIndex:        0,
		}
	}

	pool.QueuedWithdrawIdentifiers = append(pool.QueuedWithdrawIdentifiers, withdrawIdentifier)
	store.Set(poolKey, k.cdc.MustMarshal(&pool))
	return nil
}

func (k Keeper) GetBtcWithdrawRequest(ctx context.Context, twilightAddress sdk.AccAddress, reserveId uint64, withdrawAddress string, withdrawAmount uint64) (*types.BtcWithdrawRequestInternal, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetBtcWithdrawRequestKeyInternal(twilightAddress, reserveId, withdrawAddress, withdrawAmount)
	bz := store.Get(aKey)
	if bz == nil {
		return nil, false
	}
	var withdrawRequest types.BtcWithdrawRequestInternal
	k.cdc.MustUnmarshal(bz, &withdrawRequest)
	return &withdrawRequest, true
}

func (k Keeper) GetBtcWithdrawRequestByIdentifier(ctx context.Context, withdrawIdentifier uint32) (*types.BtcWithdrawRequestInternal, bool) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.BtcWithdrawRequestKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var withdrawRequest types.BtcWithdrawRequestInternal
		k.cdc.MustUnmarshal(iterator.Value(), &withdrawRequest)
		if withdrawRequest.WithdrawIdentifier == withdrawIdentifier {
			return &withdrawRequest, true
		}
	}
	return nil, false
}

func (k Keeper) ConfirmWithdrawRequestsAfterSweepConfirmation(ctx context.Context, reserveId uint64, roundId uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	snapshot, found := k.GetReserveWithdrawSnapshot(ctx, reserveId, roundId)
	if !found {
		sdkCtx.Logger().Error("reserve withdraw snapshot not found", "reserveId", reserveId, "roundId", roundId)
		return nil
	}

	pool, found := k.GetReserveWithdrawPool(ctx, reserveId)
	if !found {
		return fmt.Errorf("ReserveWithdrawPool not found for reserveId %d", reserveId)
	}

	processingMap := make(map[uint32]struct{})
	for _, id := range pool.ProcessingWithdrawIdentifiers {
		processingMap[id] = struct{}{}
	}

	for _, withdrawSnap := range snapshot.WithdrawRequests {
		withdrawRequest, found := k.GetBtcWithdrawRequestByIdentifier(ctx, withdrawSnap.WithdrawIdentifier)
		if !found {
			return fmt.Errorf("btc withdraw request not found for identifier %d", withdrawSnap.WithdrawIdentifier)
		}

		withdrawRequest.IsConfirmed = true
		err := k.SetBtcWithdrawRequestAfterSweepConfirmation(ctx, withdrawRequest)
		if err != nil {
			return fmt.Errorf("failed to update btc withdraw request: %w", err)
		}

		delete(processingMap, withdrawSnap.WithdrawIdentifier)
	}

	newProcessingIdentifiers := make([]uint32, 0, len(processingMap))
	for id := range processingMap {
		newProcessingIdentifiers = append(newProcessingIdentifiers, id)
	}
	pool.ProcessingWithdrawIdentifiers = newProcessingIdentifiers

	err := k.SetReserveWithdrawPool(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to update reserve withdraw pool: %w", err)
	}

	return nil
}

func (k Keeper) SetBtcWithdrawRequestAfterSweepConfirmation(ctx context.Context, withdrawRequest *types.BtcWithdrawRequestInternal) error {
	store := k.KVStore(ctx)

	twilightAddr, err := sdk.AccAddressFromBech32(withdrawRequest.TwilightAddress)
	if err != nil {
		return fmt.Errorf("invalid twilight address: %s: %w", withdrawRequest.TwilightAddress, err)
	}

	aKey := types.GetBtcWithdrawRequestKeyInternal(twilightAddr, withdrawRequest.WithdrawReserveId, withdrawRequest.WithdrawAddress, withdrawRequest.WithdrawAmount)
	store.Set(aKey, k.cdc.MustMarshal(withdrawRequest))

	return nil
}

func (k Keeper) SetNewSweepProposalReceived(ctx context.Context, reserveId uint64, roundId uint64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.KVStore(ctx)
	key := types.GetNewSweepProposalReceivedKey()

	proposalReceived := types.NewSweepProposalReceivedInternal{
		ReserveId:                   reserveId,
		RoundId:                     roundId,
		CreationTwilightBlockHeight: sdkCtx.BlockHeight(),
	}

	value := k.cdc.MustMarshal(&proposalReceived)
	store.Set(key, value)
}

func (k Keeper) CheckForNewSweepProposal(ctx context.Context) (bool, uint64, uint64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetNewSweepProposalReceivedKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposalReceived types.NewSweepProposalReceivedInternal
		k.cdc.MustUnmarshal(iterator.Value(), &proposalReceived)

		if proposalReceived.CreationTwilightBlockHeight == sdkCtx.BlockHeight()-1 {
			return true, proposalReceived.ReserveId, proposalReceived.RoundId
		}
	}
	return false, 0, 0
}

func (k Keeper) GetReserveWithdrawSnapshot(ctx context.Context, reserveId uint64, roundId uint64) (*types.ReserveWithdrawSnapshot, bool) {
	store := k.KVStore(ctx)
	key := types.GetReserveWithdrawSnapshotKey(reserveId, roundId)
	if !store.Has(key) {
		return nil, false
	}
	bz := store.Get(key)
	var snapshot types.ReserveWithdrawSnapshot
	k.cdc.MustUnmarshal(bz, &snapshot)
	return &snapshot, true
}

func (k Keeper) PruneReserveWithdrawSnapshot(ctx context.Context, reserveId uint64, roundId uint64) {
	store := k.KVStore(ctx)
	key := types.GetReserveWithdrawSnapshotKey(reserveId, roundId)
	store.Delete(key)
}
