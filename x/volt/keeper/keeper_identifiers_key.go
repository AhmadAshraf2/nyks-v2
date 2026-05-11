package keeper

import (
	"context"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxUint16            = 65535
	WithdrawalCounterKey = "withdrawal_counter"
	DepositCounterKey    = "deposit_counter"
)

func (k Keeper) InitCounters(ctx context.Context) {
	k.InitCounter(ctx, WithdrawalCounterKey)
	k.InitCounter(ctx, DepositCounterKey)
}

func (k Keeper) InitCounter(ctx context.Context, counterKey string) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.KVStore(sdkCtx)
	if !store.Has([]byte(counterKey)) {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, 0)
		store.Set([]byte(counterKey), b)
	}
}

func (k Keeper) IncrementCounter(ctx context.Context, counterKey string) uint32 {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.KVStore(sdkCtx)
	var counter uint32

	if k.GetCounter(ctx, counterKey) >= MaxUint16 {
		counter = k.GetCounter(ctx, counterKey) + 1
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, counter)
		store.Set([]byte(counterKey), b)
	} else {
		counter = uint32(k.GetCounter(ctx, counterKey) + 1)
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(counter))
		store.Set([]byte(counterKey), b)
	}
	return counter
}

func (k Keeper) GetCounter(ctx context.Context, counterKey string) uint32 {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.KVStore(sdkCtx)
	b := store.Get([]byte(counterKey))
	if len(b) == 2 {
		return uint32(binary.LittleEndian.Uint16(b))
	} else if len(b) == 4 {
		return binary.LittleEndian.Uint32(b)
	}
	return 0
}
