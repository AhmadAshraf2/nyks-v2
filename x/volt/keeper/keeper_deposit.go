package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"twilight-project/nyks/x/volt/types"
)

func (k Keeper) SetBtcDeposit(ctx context.Context, btcDepositAddress string, twilightAddress sdk.AccAddress, twilightStakingAmount uint64, btcSatoshiTestAmount uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	btcDepositAddr := &types.BtcDepositAddress{
		BtcDepositAddress:           btcDepositAddress,
		BtcSatoshiTestAmount:        btcSatoshiTestAmount,
		TwilightStakingAmount:       twilightStakingAmount,
		TwilightAddress:             twilightAddress.String(),
		IsConfirmed:                 false,
		CreationTwilightBlockHeight: sdkCtx.BlockHeight(),
	}

	store := k.KVStore(ctx)
	aKey := types.GetBtcDepositKey(twilightAddress)
	store.Set(aKey, k.cdc.MustMarshal(btcDepositAddr))

	return nil
}

func (k Keeper) SetBtcDepositConfirmed(ctx context.Context, twilightDepositAddress sdk.AccAddress) error {
	btcDepositAddress, found := k.GetBtcDepositAddressByTwilightAddress(ctx, twilightDepositAddress)
	if !found {
		return types.ErrBtcDepositAddressNotFound.Wrap("A BtcDepositAddress msg doesn't exist with the given twilight address")
	}

	btcDepositAddress.IsConfirmed = true

	store := k.KVStore(ctx)
	aKey := types.GetBtcDepositKey(twilightDepositAddress)
	store.Set(aKey, k.cdc.MustMarshal(btcDepositAddress))

	return nil
}

func (k Keeper) GetBtcDepositAddressByTwilightAddress(ctx context.Context, twilightDepositAddress sdk.AccAddress) (btcDepositAddress *types.BtcDepositAddress, found bool) {
	store := k.KVStore(ctx)
	aKey := types.GetBtcDepositKey(twilightDepositAddress)
	if !store.Has(aKey) {
		return nil, false
	}

	bz := store.Get(aKey)
	var addr types.BtcDepositAddress
	k.cdc.MustUnmarshal(bz, &addr)

	return &addr, true
}

func (k Keeper) CheckBtcAddress(ctx context.Context, twilightAddress sdk.Address, btcAddress string, newSatoshiTestAmount uint64) bool {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.BtcDepositKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var btcDepositAddress types.BtcDepositAddress
		k.cdc.MustUnmarshal(iterator.Value(), &btcDepositAddress)

		if btcDepositAddress.TwilightAddress == twilightAddress.String() {
			continue
		}

		isSameAddressAndConfirmed := btcDepositAddress.BtcDepositAddress == btcAddress && btcDepositAddress.IsConfirmed
		isSameAddressAndMatchingSatoshiAmount := btcDepositAddress.BtcDepositAddress == btcAddress && btcDepositAddress.BtcSatoshiTestAmount == newSatoshiTestAmount

		if isSameAddressAndConfirmed || isSameAddressAndMatchingSatoshiAmount {
			return true
		}
	}

	return false
}

func (k Keeper) GetAllBtcRegisteredDepositAddresses(ctx context.Context) (btcDepositAddresses []types.BtcDepositAddress) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.BtcDepositKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var btcDepositAddress types.BtcDepositAddress
		k.cdc.MustUnmarshal(iterator.Value(), &btcDepositAddress)
		btcDepositAddresses = append(btcDepositAddresses, btcDepositAddress)
	}

	return btcDepositAddresses
}
