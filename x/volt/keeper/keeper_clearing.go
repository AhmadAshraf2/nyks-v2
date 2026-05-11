package keeper

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	forkstypes "twilight-project/nyks/x/forks/types"
	"twilight-project/nyks/x/volt/types"
)

func (k Keeper) SetClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress, account *types.ClearingAccount) error {
	store := k.KVStore(ctx)
	aKey := types.GetClearingAccountKey(twilightAddress)
	store.Set(aKey, k.cdc.MustMarshal(account))
	return nil
}

func (k Keeper) SetBtcAddressForClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress, btcAddr string, depositIdentifer uint32) (*types.ClearingAccount, error) {
	account := &types.ClearingAccount{
		TwilightAddress:             twilightAddress.String(),
		BtcDepositAddress:           btcAddr,
		BtcDepositAddressIdentifier: depositIdentifer,
	}

	store := k.KVStore(ctx)
	aKey := types.GetClearingAccountKey(twilightAddress)
	store.Set(aKey, k.cdc.MustMarshal(account))

	return account, nil
}

func (k Keeper) SetBtcAddressForExistingClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress, btcAddr string, depositIdentifer uint32) (*types.ClearingAccount, error) {
	account, found := k.GetClearingAccount(ctx, twilightAddress)
	if !found {
		return nil, fmt.Errorf("clearing account not found: %s", twilightAddress)
	}

	account.BtcDepositAddress = btcAddr
	account.BtcDepositAddressIdentifier = depositIdentifer
	k.SetClearingAccount(ctx, twilightAddress, account)

	return account, nil
}

func (k Keeper) GetClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress) (*types.ClearingAccount, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetClearingAccountKey(twilightAddress)
	if !store.Has(aKey) {
		return nil, false
	}

	bz := store.Get(aKey)
	var clearingAccount types.ClearingAccount
	k.cdc.MustUnmarshal(bz, &clearingAccount)

	return &clearingAccount, true
}

func (k Keeper) GetAllClearingAccounts(ctx context.Context) ([]types.ClearingAccount, error) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.TwilightClearingAccountKey)
	defer iterator.Close()

	var clearingAccounts []types.ClearingAccount
	for ; iterator.Valid(); iterator.Next() {
		var clearingAccount types.ClearingAccount
		k.cdc.MustUnmarshal(iterator.Value(), &clearingAccount)
		clearingAccounts = append(clearingAccounts, clearingAccount)
	}

	return clearingAccounts, nil
}

func (k Keeper) GetAllClearingAccountsInaReserve(ctx context.Context, reserveId uint64) ([]types.ClearingAccount, bool) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.TwilightClearingAccountKey)
	defer iterator.Close()

	var clearingAccounts []types.ClearingAccount
	for ; iterator.Valid(); iterator.Next() {
		var clearingAccount types.ClearingAccount
		k.cdc.MustUnmarshal(iterator.Value(), &clearingAccount)

		for _, balance := range clearingAccount.ReserveAccountBalances {
			if balance.ReserveId == reserveId {
				clearingAccounts = append(clearingAccounts, clearingAccount)
				break
			}
		}
	}

	return clearingAccounts, true
}

func (k Keeper) DeductFromClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress, reserveId uint64, amount uint64) error {
	account, found := k.GetClearingAccount(ctx, twilightAddress)
	if !found {
		return fmt.Errorf("clearing account not found: %s", twilightAddress)
	}

	for _, balance := range account.ReserveAccountBalances {
		if balance.ReserveId == reserveId {
			if balance.Amount < amount {
				return fmt.Errorf("insufficient balance in reserve: %d", balance.Amount)
			}
			balance.Amount -= amount
		}
	}

	k.SetClearingAccount(ctx, twilightAddress, account)
	return nil
}

func (k Keeper) CheckClearingAccountBalance(ctx context.Context, twilightAddress sdk.AccAddress, reserveId uint64, amount uint64) error {
	account, found := k.GetClearingAccount(ctx, twilightAddress)
	if !found {
		return fmt.Errorf("clearing account not found: %s", twilightAddress)
	}
	for _, balance := range account.ReserveAccountBalances {
		if balance.ReserveId == reserveId {
			if balance.Amount < amount {
				return fmt.Errorf("insufficient balance in reserve: %d", balance.Amount)
			}
		}
	}
	return nil
}

func (k Keeper) UpdateTransfersInClearing(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	fromAccount, fromExists := k.GetClearingAccount(ctx, from)
	moduleAddr := k.accountKeeper.GetModuleAddress(forkstypes.ModuleName)
	isFromNyksModule := from.Equals(moduleAddr)

	if isFromNyksModule {
		return nil
	}

	if from.Equals(to) {
		return nil
	}

	if !fromExists {
		return fmt.Errorf("clearing account not found: %s", from)
	}

	toAccount, toExists := k.GetClearingAccount(ctx, to)
	if !toExists {
		toAccount = &types.ClearingAccount{
			TwilightAddress: to.String(),
		}
		k.SetClearingAccount(ctx, to, toAccount)
	}

	totalAmount := uint64(amount.AmountOf("sats").Int64())

	if fromExists {
		for _, balance := range fromAccount.ReserveAccountBalances {
			if balance.Amount == 0 {
				continue
			}

			transferAmount := balance.Amount
			if transferAmount > totalAmount {
				transferAmount = totalAmount
			}

			balance.Amount -= transferAmount

			found := false
			for j, toBalance := range toAccount.ReserveAccountBalances {
				if toBalance.ReserveId == balance.ReserveId {
					toBalance.Amount += transferAmount
					toAccount.ReserveAccountBalances[j] = toBalance
					found = true
					break
				}
			}

			if !found {
				toAccount.ReserveAccountBalances = append(toAccount.ReserveAccountBalances, &types.IndividualTwilightReserveAccountBalance{
					ReserveId: balance.ReserveId,
					Amount:    transferAmount,
				})
			}

			totalAmount -= transferAmount
			if totalAmount == 0 {
				break
			}
		}
	}

	if totalAmount != 0 {
		return fmt.Errorf("sender does not have enough funds in clearing")
	}

	k.SetClearingAccount(ctx, from, fromAccount)
	k.SetClearingAccount(ctx, to, toAccount)

	return nil
}

func (k Keeper) GetRefundTxSnapshot(ctx context.Context, reserveId uint64, roundId uint64) (*types.RefundTxSnapshot, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetRefundTxSnapshotKey(reserveId, roundId)
	if !store.Has(aKey) {
		return nil, false
	}

	bz := store.Get(aKey)
	var refundTxSnapshot types.RefundTxSnapshot
	k.cdc.MustUnmarshal(bz, &refundTxSnapshot)

	return &refundTxSnapshot, true
}

func (k Keeper) PruneRefundTxSnapshot(ctx context.Context, reserveId uint64, roundId uint64) {
	store := k.KVStore(ctx)
	key := types.GetRefundTxSnapshotKey(reserveId, roundId)
	store.Delete(key)
}
