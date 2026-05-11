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

func (k Keeper) RegisterNewFragment(ctx context.Context, judgeAddress sdk.AccAddress, threshold uint64, applicationFee uint64, numOfSigners uint64, fragmentFeeBips uint64, arbitraryData string) (uint64, error) {
	LastRegisteredFragment := k.GetLastRegisteredFragment(ctx)
	fragmentId := LastRegisteredFragment + 1

	if fragmentId > types.FragmentMaxLimit {
		return 0, fmt.Errorf("fragment max limit reached: %d", types.FragmentMaxLimit)
	}

	fragment := &types.Fragment{
		FragmentId:           fragmentId,
		FragmentStatus:       false,
		JudgeAddress:         judgeAddress.String(),
		JudgeStatus:          true,
		Signers:              []*types.FragmentSigners{},
		SignerApplicationFee: applicationFee,
		Threshold:            threshold,
		FeePool:              0,
		FragmentFeeBips:      fragmentFeeBips,
		ArbitraryData:        arbitraryData,
		ReserveIds:           []uint64{},
	}

	err := k.SetFragment(ctx, fragment)
	if err != nil {
		return 0, fmt.Errorf("could not set fragment: %w", err)
	}
	k.setLastRegisteredFragment(ctx, fragmentId)

	return fragmentId, nil
}

func (k Keeper) SetFragment(ctx context.Context, fragment *types.Fragment) error {
	store := k.KVStore(ctx)
	aKey := types.GetFragmentKey(fragment.FragmentId)
	store.Set(aKey, k.cdc.MustMarshal(fragment))
	return nil
}

func (k Keeper) GetFragment(ctx context.Context, fragmentId uint64) (*types.Fragment, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetFragmentKey(fragmentId)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var fragment types.Fragment
	k.cdc.MustUnmarshal(bz, &fragment)
	return &fragment, true
}

func (k Keeper) setLastRegisteredFragment(ctx context.Context, fragmentId uint64) {
	store := k.KVStore(ctx)
	store.Set(types.LastRegisteredFragmentKey, forkstypes.UInt64Bytes(fragmentId))
}

func (k Keeper) GetLastRegisteredFragment(ctx context.Context) uint64 {
	store := k.KVStore(ctx)
	bytes := store.Get(types.LastRegisteredFragmentKey)
	if len(bytes) == 0 {
		return 0
	}
	return forkstypes.UInt64FromBytes(bytes)
}

func (k Keeper) UpdateFragmentReserves(ctx context.Context, fragmentId uint64, reserveId uint64) error {
	fragment, found := k.GetFragment(ctx, fragmentId)
	if !found {
		return fmt.Errorf("fragment ID %d not found", fragmentId)
	}

	for _, id := range fragment.ReserveIds {
		if id == reserveId {
			return fmt.Errorf("reserve ID %d already exists in fragment ID %d", reserveId, fragmentId)
		}
	}

	fragment.ReserveIds = append(fragment.ReserveIds, reserveId)
	return k.SetFragment(ctx, fragment)
}

func (k Keeper) setLastRegisteredApplicationId(ctx context.Context, applicationId uint64) {
	store := k.KVStore(ctx)
	store.Set(types.LastRegisteredFragmentApplicationKey, forkstypes.UInt64Bytes(applicationId))
}

func (k Keeper) GetLastRegisteredApplicationId(ctx context.Context) uint64 {
	store := k.KVStore(ctx)
	bytes := store.Get(types.LastRegisteredFragmentApplicationKey)
	if len(bytes) == 0 {
		return 0
	}
	return forkstypes.UInt64FromBytes(bytes)
}

func (k Keeper) SetSignerApplication(ctx context.Context, msg *types.SignerApplicationData) {
	store := k.KVStore(ctx)
	aKey := types.GetSignerApplicationFeeKey(msg.FragmentId, msg.ApplicationId)
	store.Set(aKey, k.cdc.MustMarshal(msg))
}

func (k Keeper) GetSignerApplication(ctx context.Context, fragmentId uint64, applicationId uint64) (*types.SignerApplicationData, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetSignerApplicationFeeKey(fragmentId, applicationId)
	bz := store.Get(aKey)
	if bz == nil {
		return nil, false
	}
	var application types.SignerApplicationData
	k.cdc.MustUnmarshal(bz, &application)
	return &application, true
}

func (k Keeper) GetSignerApplicationBySignerAndFragmentId(ctx context.Context, fragmentId uint64, signerAddress string, btcPubKey string) (*types.SignerApplicationData, bool) {
	store := k.KVStore(ctx)
	prefix := types.GetSignerApplicationFeePrefix(fragmentId)
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var application types.SignerApplicationData
		k.cdc.MustUnmarshal(iterator.Value(), &application)
		if application.SignerAddress == signerAddress || application.BtcPubKey == btcPubKey {
			return &application, true
		}
	}

	return nil, false
}

func (k Keeper) GetSignerApplications(ctx context.Context, fragmentId uint64) ([]types.SignerApplicationData, bool) {
	store := k.KVStore(ctx)
	prefix := types.GetSignerApplicationFeePrefix(fragmentId)
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var applications []types.SignerApplicationData
	for ; iterator.Valid(); iterator.Next() {
		var application types.SignerApplicationData
		k.cdc.MustUnmarshal(iterator.Value(), &application)
		applications = append(applications, application)
	}

	if len(applications) == 0 {
		return nil, false
	}
	return applications, true
}

func (k Keeper) GetExistingSignerInFragments(ctx context.Context, signerAddress string) bool {
	found := false
	k.IterateFragments(ctx, func(_ []byte, res types.Fragment) bool {
		for _, signer := range res.Signers {
			if signer.SignerAddress == signerAddress {
				found = true
				return true
			}
		}
		return false
	})
	return found
}

func (k Keeper) IterateFragments(ctx context.Context, cb func([]byte, types.Fragment) bool) {
	store := k.KVStore(ctx)
	iter := storetypes.KVStorePrefixIterator(store, types.FragmentKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var fragment types.Fragment
		k.cdc.MustUnmarshal(iter.Value(), &fragment)
		if cb(iter.Key(), fragment) {
			return
		}
	}
}

func (k Keeper) ReturnSignerApplicationFee(ctx context.Context, aSignerAddress string, applicationFee uint64) error {
	signerAddress, err := sdk.AccAddressFromBech32(aSignerAddress)
	if err != nil {
		return err
	}

	feeAmount := sdk.NewCoin("nyks", math.NewIntFromUint64(applicationFee))

	voltModuleAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	if voltModuleAcc == nil {
		return types.ErrVoltModuleAccountNotFound
	}

	balance := k.BankKeeper.GetBalance(ctx, voltModuleAcc.GetAddress(), "nyks")
	if balance.Amount.LT(math.NewIntFromUint64(applicationFee)) {
		return types.ErrInsufficientFunds
	}

	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, signerAddress, sdk.NewCoins(feeAmount))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetFragementForReserveId(ctx context.Context, reserveId uint64) (*types.Fragment, bool) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.FragmentKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var fragment types.Fragment
		k.cdc.MustUnmarshal(iterator.Value(), &fragment)
		for _, id := range fragment.ReserveIds {
			if id == reserveId {
				return &fragment, true
			}
		}
	}
	return nil, false
}

func (k Keeper) CheckSignerInFragment(ctx context.Context, reserveId uint64, signerAddress sdk.AccAddress) bool {
	fragment, found := k.GetFragementForReserveId(ctx, reserveId)
	if !found {
		return false
	}
	for _, signer := range fragment.Signers {
		if signer.SignerAddress == signerAddress.String() {
			return true
		}
	}
	return false
}

func (k Keeper) CheckReserveWithdrawSnapshot(ctx context.Context, btcTxHex string, reserveId uint64, roundId uint64) (bool, error) {
	_, found := k.GetReserveWithdrawSnapshot(ctx, reserveId, roundId)
	if !found {
		return false, fmt.Errorf("reserve withdraw snapshot not found")
	}
	return true, nil
}
