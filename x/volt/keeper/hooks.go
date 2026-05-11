package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TrackBeforeSend tracks sats transfers to update clearing accounts.
// This is called via SendRestrictionFn registered on the BankKeeper in app.go
// instead of the old BankHooks pattern from the forked cosmos-sdk.
func (k Keeper) TrackBeforeSend(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error {
	for _, coin := range amount {
		if coin.Denom == "sats" {
			err := k.UpdateTransfersInClearing(ctx, from, to, amount)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
