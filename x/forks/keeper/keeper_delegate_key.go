package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"twilight-project/nyks/x/forks/types"
)

func (k Keeper) SetDelegateAddresses(ctx sdk.Context, msg *types.MsgSetDelegateAddresses) error {
	val, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	store := k.KVStore(ctx)
	key := types.GetValidatorKey(val)

	bz := store.Get(key)
	if bz == nil {
		bz := k.cdc.MustMarshal(msg)
		store.Set(key, bz)
	} else {
		var existingMsg types.MsgSetDelegateAddresses
		k.cdc.MustUnmarshal(bz, &existingMsg)

		existingMsg.BtcOracleAddress = msg.BtcOracleAddress
		if msg.BtcPublicKey != "" {
			existingMsg.BtcPublicKey = msg.BtcPublicKey
		}
		if msg.ZkOracleAddress != "" {
			existingMsg.ZkOracleAddress = msg.ZkOracleAddress
		}

		bz := k.cdc.MustMarshal(&existingMsg)
		store.Set(key, bz)
	}

	return nil
}

func (k Keeper) GetDelegateAddresses(ctx sdk.Context, orchestratorAddress sdk.AccAddress) (*types.MsgSetDelegateAddresses, bool) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyValidator)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgSetDelegateAddresses
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		if msg.BtcOracleAddress == orchestratorAddress.String() {
			return &msg, true
		}
	}
	return nil, false
}

func (k Keeper) GetAllDelegateAddresses(ctx sdk.Context) ([]types.MsgSetDelegateAddresses, error) {
	store := k.KVStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyValidator)
	defer iterator.Close()

	var msgs []types.MsgSetDelegateAddresses
	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgSetDelegateAddresses
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (k Keeper) CheckOrchestratorValidatorInSet(goCtx context.Context, orchestrator string) (sdk.ValAddress, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, err := sdk.AccAddressFromBech32(orchestrator)
	if err != nil {
		return nil, fmt.Errorf("invalid orchestrator address: %w", err)
	}

	delegateAddresses, found := k.GetDelegateAddresses(ctx, accAddr)
	if !found {
		return nil, fmt.Errorf("invalid btc oracle account address")
	}

	valAddress, err := sdk.ValAddressFromBech32(delegateAddresses.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("validator stored address is invalid bech32")
	}

	val, err := k.StakingKeeper.Validator(ctx, valAddress)
	if err != nil || val == nil {
		return nil, fmt.Errorf("validator not found")
	}
	if !val.IsBonded() {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("validator not in active set")
	}

	return valAddress, nil
}

func (k Keeper) SetBtcPublicKeyForValidator(ctx sdk.Context, validator sdk.ValAddress, btcPk types.BtcPublicKey) ([]byte, error) {
	btcPkBytes, err := hex.DecodeString(btcPk.GetBtcPublicKey())
	if err != nil {
		return nil, fmt.Errorf("invalid btc public key hex encoding (%s): %w", btcPk.GetBtcPublicKey(), err)
	}
	store := k.KVStore(ctx)
	store.Set(types.GetBtcPublicKeyByValidatorKey(validator), btcPkBytes)
	store.Set(types.GetValidatorByBtcPublicKeyKey(btcPk), []byte(validator))

	return btcPkBytes, nil
}

func (k Keeper) GetBtcPublicKeyByValidator(ctx sdk.Context, validator sdk.ValAddress) (btcPublicKey *types.BtcPublicKey, found bool) {
	store := k.KVStore(ctx)
	btcPk := store.Get(types.GetBtcPublicKeyByValidatorKey(validator))
	if btcPk == nil {
		return nil, false
	}

	pk, err := types.NewBtcPublicKey(hex.EncodeToString(btcPk))
	if err != nil {
		return nil, false
	}
	return pk, true
}

func (k Keeper) GetDelegateKeys(ctx sdk.Context) ([]types.MsgSetDelegateAddresses, error) {
	store := k.KVStore(ctx)
	prefix := types.BtcPublicKeyByValidatorKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	btcPublicKeys := make(map[string]string)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()[len(types.BtcPublicKeyByValidatorKey):]
		value := iter.Value()
		btcPk, err := types.NewBtcPublicKey(hex.EncodeToString(value))
		if err != nil {
			return nil, fmt.Errorf("found invalid btcPk %v under key %v: %w", string(value), key, err)
		}
		valAddress := sdk.ValAddress(key)
		btcPublicKeys[valAddress.String()] = btcPk.GetBtcPublicKey()
	}

	store = k.KVStore(ctx)
	prefix = types.KeyOrchestratorAddress
	iter = store.Iterator(prefixRange(prefix))
	defer iter.Close()

	btcOracleAddresses := make(map[string]string)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()[len(types.KeyOrchestratorAddress):]
		value := iter.Value()
		orchAddress := sdk.AccAddress(key)
		valAddress := sdk.ValAddress(value)
		btcOracleAddresses[valAddress.String()] = orchAddress.String()
	}

	var result []types.MsgSetDelegateAddresses
	for valAddr, btcPk := range btcPublicKeys {
		oracle, ok := btcOracleAddresses[valAddr]
		if !ok {
			panic("Can't find address")
		}
		result = append(result, types.MsgSetDelegateAddresses{
			ValidatorAddress: valAddr,
			BtcOracleAddress: oracle,
			BtcPublicKey:     btcPk,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].BtcPublicKey < result[j].BtcPublicKey
	})

	return result, nil
}
