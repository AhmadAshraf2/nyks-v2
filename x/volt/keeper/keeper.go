package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/volt/types"
)

type Keeper struct {
	storeService  corestore.KVStoreService
	cdc           codec.Codec
	addressCodec  address.Codec
	authority     []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	accountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	BridgeKeeper  types.BridgeKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	bridgeKeeper types.BridgeKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		authority:     authority,
		accountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		BridgeKeeper:  bridgeKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// KVStore returns a KVStore adapter for the keeper's store service.
func (k Keeper) KVStore(ctx context.Context) storetypes.KVStore {
	return runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
}

// Store returns the KVStore (legacy name used by EndBlocker)
func (k Keeper) Store(ctx context.Context) storetypes.KVStore {
	return k.KVStore(ctx)
}

// Codec returns the codec for the volt module
func (k Keeper) Codec() codec.Codec {
	return k.cdc
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
