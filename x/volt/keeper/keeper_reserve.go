package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	forkstypes "twilight-project/nyks/x/forks/types"
	"twilight-project/nyks/x/volt/types"
)

func (k Keeper) RegisterNewBtcReserve(ctx context.Context, judgeAddress sdk.AccAddress, reserveAddress string) (uint64, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	LastRegisteredReserve := k.GetLastRegisteredBtcReserve(ctx)
	reserveId := LastRegisteredReserve + 1

	if reserveId > types.BtcReserveMaxLimit {
		return 0, fmt.Errorf("btc max reserve limit reached: %d", types.BtcReserveMaxLimit)
	}

	res := &types.BtcReserve{
		ReserveId:             reserveId,
		ReserveAddress:        reserveAddress,
		JudgeAddress:          judgeAddress.String(),
		BtcRelayCapacityValue: 0,
		TotalValue:            0,
		PrivatePoolValue:      0,
		PublicValue:           0,
		FeePool:               0,
		UnlockHeight:          0,
		RoundId:               0,
	}

	err := k.SetBtcReserve(sdkCtx, res)
	if err != nil {
		return 0, fmt.Errorf("could not set reserve: %s", reserveAddress)
	}
	k.setLastRegisteredBtcReserve(sdkCtx, reserveId)

	withdrawPool := &types.ReserveWithdrawPool{
		ReserveID:                     reserveId,
		RoundID:                       0,
		ProcessingWithdrawIdentifiers: []uint32{},
		QueuedWithdrawIdentifiers:     []uint32{},
		CurrentProcessingIndex:        0,
	}

	err = k.SetReserveWithdrawPool(sdkCtx, withdrawPool)
	if err != nil {
		return 0, fmt.Errorf("could not set reserve withdraw pool: %w", err)
	}

	return reserveId, nil
}

func (k Keeper) SetBtcReserve(ctx context.Context, reserve *types.BtcReserve) error {
	store := k.KVStore(ctx)
	aKey := types.GetReserveKey(reserve.ReserveId)
	store.Set(aKey, k.cdc.MustMarshal(reserve))
	return nil
}

func (k Keeper) UpdateBtcReserveAfterMint(ctx context.Context, mintedValue uint64, twilightAddress sdk.AccAddress, reserveAddress string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	reserveId, err := k.GetBtcReserveIdByAddress(ctx, reserveAddress)
	if err != nil {
		return fmt.Errorf("btc reserve not found: %s", reserveAddress)
	}

	reserve, err := k.GetBtcReserve(ctx, reserveId)
	if err != nil {
		return fmt.Errorf("btc reserve not found: %s", reserveAddress)
	}

	reserve.TotalValue = reserve.TotalValue + mintedValue
	reserve.PublicValue = reserve.PublicValue + mintedValue

	clearingAccount, foundClearing := k.GetClearingAccount(ctx, twilightAddress)
	if !foundClearing || clearingAccount.BtcDepositAddress == "" {
		btcDeposit, found := k.GetBtcDepositAddressByTwilightAddress(ctx, twilightAddress)
		if !found {
			return fmt.Errorf("btc deposit address not found: %s", twilightAddress)
		}

		if btcDeposit.BtcSatoshiTestAmount != mintedValue {
			return fmt.Errorf("btc satoshi test amount not equal for: %s", twilightAddress)
		}

		if btcDeposit.TwilightStakingAmount > 0 {
			err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, forkstypes.ModuleName, twilightAddress, sdk.NewCoins(sdk.NewCoin("nyks", math.NewIntFromUint64(btcDeposit.TwilightStakingAmount))))
			if err != nil {
				return err
			}
		}

		err = k.SetBtcDepositConfirmed(sdkCtx, twilightAddress)
		if err != nil {
			return fmt.Errorf("clearing account not found: %s", twilightAddress)
		}

		depositIdentifier := k.IncrementCounter(ctx, DepositCounterKey)

		if !foundClearing {
			clearingAccount, err = k.SetBtcAddressForClearingAccount(ctx, twilightAddress, btcDeposit.BtcDepositAddress, depositIdentifier)
			if err != nil {
				return fmt.Errorf("clearing account not found: %s", twilightAddress)
			}
		} else if clearingAccount.BtcDepositAddress == "" {
			clearingAccount, err = k.SetBtcAddressForExistingClearingAccount(ctx, twilightAddress, btcDeposit.BtcDepositAddress, depositIdentifier)
			if err != nil {
				return fmt.Errorf("clearing account not found: %s", twilightAddress)
			}
		}
	}

	foundBalance := false
	for _, balance := range clearingAccount.ReserveAccountBalances {
		if balance.ReserveId == reserveId {
			balance.Amount += mintedValue
			foundBalance = true
			break
		}
	}

	if !foundBalance {
		clearingAccount.ReserveAccountBalances = append(clearingAccount.ReserveAccountBalances, &types.IndividualTwilightReserveAccountBalance{
			ReserveId: reserveId,
			Amount:    mintedValue,
		})
	}

	k.SetClearingAccount(ctx, twilightAddress, clearingAccount)

	store := k.KVStore(ctx)
	aKey := types.GetReserveKey(reserveId)
	store.Set(aKey, k.cdc.MustMarshal(reserve))

	return nil
}

func (k Keeper) UpdateBtcReserveAfterSweepProposal(ctx context.Context, reserveId uint64, reserveAddress string, judgeAddress string, btcBlockNumber uint64, btcRelayCapacityValue uint64, btcTxHash string, unlockHeight uint64, roundId uint64) error {
	reserve, err := k.GetBtcReserve(ctx, reserveId)
	if err != nil {
		return fmt.Errorf("btc reserve not found: %s", reserveAddress)
	}

	reserve.ReserveAddress = reserveAddress
	reserve.JudgeAddress = judgeAddress
	reserve.BtcRelayCapacityValue = btcRelayCapacityValue
	reserve.UnlockHeight = unlockHeight
	reserve.RoundId = roundId

	store := k.KVStore(ctx)
	aKey := types.GetReserveKey(reserveId)
	store.Set(aKey, k.cdc.MustMarshal(reserve))

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.setLastUnlockedReserve(sdkCtx, reserveId)

	return nil
}

func (k Keeper) GetBtcReserve(ctx context.Context, reserveId uint64) (*types.BtcReserve, error) {
	store := k.KVStore(ctx)
	aKey := types.GetReserveKey(reserveId)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil, fmt.Errorf("btc reserve not found: %d", reserveId)
	}
	reserve := &types.BtcReserve{}
	err := k.cdc.Unmarshal(bz, reserve)
	if err != nil {
		return nil, fmt.Errorf("btc reserve not found: %d", reserveId)
	}
	return reserve, nil
}

func (k Keeper) CheckBtcReserveExists(ctx context.Context, reserveId uint64) bool {
	store := k.KVStore(ctx)
	aKey := types.GetReserveKey(reserveId)
	return store.Has(aKey)
}

func (k Keeper) GetBtcReserveIdByAddress(ctx context.Context, reserveAddress string) (uint64, error) {
	store := k.KVStore(ctx)
	prefix := types.BtcReserveKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var res types.BtcReserve
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if res.ReserveAddress == reserveAddress {
			return res.ReserveId, nil
		}
	}

	return 0, fmt.Errorf("btc reserve not found: %s", reserveAddress)
}

func (k Keeper) setLastRegisteredBtcReserve(ctx context.Context, reserveId uint64) {
	store := k.KVStore(ctx)
	store.Set(types.LastRegisteredReserveKey, forkstypes.UInt64Bytes(reserveId))
}

func (k Keeper) GetLastRegisteredBtcReserve(ctx context.Context) uint64 {
	store := k.KVStore(ctx)
	bytes := store.Get(types.LastRegisteredReserveKey)
	if len(bytes) == 0 {
		return 0
	}
	return forkstypes.UInt64FromBytes(bytes)
}

func (k Keeper) IterateBtcReserves(ctx context.Context, cb func([]byte, types.BtcReserve) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcReserveKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var res types.BtcReserve
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) setLastUnlockedReserve(ctx context.Context, reserveId uint64) {
	store := k.KVStore(ctx)
	store.Set(types.LastUnlockedReserveKey, forkstypes.UInt64Bytes(reserveId))
}

func (k Keeper) GetLastUnlockedReserve(ctx context.Context) uint64 {
	store := k.KVStore(ctx)
	bytes := store.Get(types.LastUnlockedReserveKey)
	if len(bytes) == 0 {
		return 0
	}
	return forkstypes.UInt64FromBytes(bytes)
}

func (k Keeper) GetNextUnlockingReserve(ctx context.Context) (*uint64, *types.BtcReserve, error) {
	reserveId := uint64(1)
	nextReserveUnlockingId := reserveId
	if reserveId >= types.BtcReserveMaxLimit {
		nextReserveUnlockingId = 1
	}

	reserve, err := k.GetBtcReserve(ctx, nextReserveUnlockingId)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot find reserve: %w", err)
	}

	return &nextReserveUnlockingId, reserve, nil
}

func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	if len(prefix) == 0 {
		return nil, nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}
