package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	volttypes "twilight-project/nyks/x/volt/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}

type VoltKeeper interface {
	SetBtcReserve(ctx context.Context, reserve *volttypes.BtcReserve) error
	GetBtcReserve(ctx context.Context, reserveId uint64) (*volttypes.BtcReserve, error)
	GetClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress) (*volttypes.ClearingAccount, bool)
	SetClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress, account *volttypes.ClearingAccount) error
	GetNextUnlockingReserve(ctx context.Context) (*uint64, *volttypes.BtcReserve, error)
	IterateBtcReserves(ctx context.Context, cb func([]byte, volttypes.BtcReserve) bool)
}
