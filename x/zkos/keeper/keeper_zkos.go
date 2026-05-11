package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	volttypes "twilight-project/nyks/x/volt/types"
	"twilight-project/nyks/x/zkos/types"
)

func (k Keeper) SetTransferTx(ctx context.Context, txId string, txByteCode string, txFee uint64, zkOracleAddress string) {
	ttx := &types.MsgTransferTx{
		TxId:            txId,
		TxByteCode:      txByteCode,
		TxFee:           txFee,
		ZkOracleAddress: zkOracleAddress,
	}

	store := k.KVStore(ctx)
	aKey := types.GetTransferTxKey(txId)
	store.Set(aKey, k.cdc.MustMarshal(ttx))
}

func (k Keeper) GetTransferTx(ctx context.Context, txId string) (types.MsgTransferTx, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetTransferTxKey(txId)
	if !store.Has(aKey) {
		return types.MsgTransferTx{}, false
	}
	bz := store.Get(aKey)
	var msg types.MsgTransferTx
	k.cdc.MustUnmarshal(bz, &msg)
	return msg, true
}

func (k Keeper) SetMintOrBurnTradingBtc(ctx context.Context, msg *types.MsgMintBurnTradingBtc) error {
	if k.HasUsedQqAccount(ctx, msg.QqAccount) {
		return types.ErrDuplicateQqAccount.Wrap("this QuisQuis account has already been used")
	}

	twilightAddress, err := sdk.AccAddressFromBech32(msg.TwilightAddress)
	if err != nil {
		return fmt.Errorf("invalid twilight address: %s", msg.TwilightAddress)
	}

	account, found := k.VoltKeeper.GetClearingAccount(ctx, twilightAddress)
	if !found {
		return fmt.Errorf("clearing account not found")
	}

	bankBalanceCoin := k.BankKeeper.GetBalance(ctx, twilightAddress, "sats")
	bankBalance := bankBalanceCoin.Amount.Uint64()

	totalBalance := uint64(0)
	for _, balance := range account.ReserveAccountBalances {
		totalBalance += balance.Amount
	}

	if bankBalance != totalBalance {
		return fmt.Errorf("bank balance not equal to reserve balance")
	}

	if msg.MintOrBurn { // Mint: move from public to private, burn user's sats
		if msg.BtcValue > totalBalance {
			return fmt.Errorf("not enough user balance in reserves: %s", msg.TwilightAddress)
		}

		remaining := msg.BtcValue
		for i, reserveBalance := range account.ReserveAccountBalances {
			if remaining == 0 {
				break
			}

			reserve, err := k.VoltKeeper.GetBtcReserve(ctx, reserveBalance.ReserveId)
			if err != nil {
				return fmt.Errorf("cannot find reserve while minting quisquis btc")
			}

			deduct := min(remaining, reserveBalance.Amount)

			if reserve.PublicValue < deduct {
				return fmt.Errorf("not enough balance in public pool for reserve %d", reserveBalance.ReserveId)
			}
			reserve.PublicValue -= deduct
			reserve.PrivatePoolValue += deduct

			reserveBalance.Amount -= deduct

			errBank := k.BankKeeper.SendCoinsFromAccountToModule(ctx, twilightAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin("sats", math.NewIntFromUint64(deduct))))
			if errBank != nil {
				return fmt.Errorf("failed to send coins to module account: %w", errBank)
			}

			errBurn := k.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("sats", math.NewIntFromUint64(deduct))))
			if errBurn != nil {
				return fmt.Errorf("failed to burn coins: %w", errBurn)
			}

			account.ReserveAccountBalances[i] = reserveBalance
			remaining -= deduct

			k.VoltKeeper.SetBtcReserve(ctx, reserve)
		}
	} else { // Burn: move from private to public, mint sats for user
		nextReserveUnlockingId, reserve, err := k.VoltKeeper.GetNextUnlockingReserve(ctx)
		if err != nil {
			return fmt.Errorf("next unlocking reserve not found: %w", err)
		}

		addBack := msg.BtcValue

		if reserve.PrivatePoolValue < addBack {
			// Not enough in this reserve, iterate all reserves to aggregate
			totalAggregated := uint64(0)
			k.VoltKeeper.IterateBtcReserves(ctx, func(_ []byte, res volttypes.BtcReserve) bool {
				if totalAggregated >= addBack {
					return true
				}
				if res.PrivatePoolValue > 0 {
					deduct := min(res.PrivatePoolValue, addBack-totalAggregated)
					res.PrivatePoolValue -= deduct
					res.PublicValue += deduct
					k.VoltKeeper.SetBtcReserve(ctx, &res)
					totalAggregated += deduct
				}
				return false
			})

			if totalAggregated < addBack {
				return fmt.Errorf("could not find enough value in private pool")
			}
		} else {
			reserve.PrivatePoolValue -= addBack
			reserve.PublicValue += addBack
		}

		// Update the clearing account
		foundBalance := false
		for _, balance := range account.ReserveAccountBalances {
			if balance.ReserveId == *nextReserveUnlockingId {
				balance.Amount += addBack
				foundBalance = true
				break
			}
		}

		if !foundBalance {
			account.ReserveAccountBalances = append(account.ReserveAccountBalances, &volttypes.IndividualTwilightReserveAccountBalance{
				ReserveId: *nextReserveUnlockingId,
				Amount:    addBack,
			})
		}

		// Mint sats and send to user
		errMint := k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("sats", math.NewIntFromUint64(addBack))))
		if errMint != nil {
			return fmt.Errorf("failed to mint new coins: %w", errMint)
		}

		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, twilightAddress, sdk.NewCoins(sdk.NewCoin("sats", math.NewIntFromUint64(addBack))))
		if err != nil {
			return fmt.Errorf("failed to send coins to user account: %w", err)
		}

		// Save the updated reserve
		k.VoltKeeper.SetBtcReserve(ctx, reserve)
	}

	// Save the updated clearing account
	k.VoltKeeper.SetClearingAccount(ctx, twilightAddress, account)

	store := k.KVStore(ctx)
	aKey := types.GetMintOrBurnTradingBtcKey(msg.TwilightAddress, msg.QqAccount)
	store.Set(aKey, k.cdc.MustMarshal(msg))

	k.MarkQqAccountAsUsed(ctx, msg.QqAccount)

	return nil
}

// DeductFeeFromPrivatePool deducts the fee from the private pool and adds to fee pool
func (k Keeper) DeductFeeFromPrivatePool(ctx context.Context, fee uint64) error {
	_, reserve, err := k.VoltKeeper.GetNextUnlockingReserve(ctx)
	if err != nil {
		return fmt.Errorf("next unlocking reserve not found: %w", err)
	}

	if reserve.PrivatePoolValue < fee {
		totalAggregated := uint64(0)
		k.VoltKeeper.IterateBtcReserves(ctx, func(_ []byte, res volttypes.BtcReserve) bool {
			if totalAggregated >= fee {
				return true
			}
			if res.PrivatePoolValue > 0 {
				deduct := min(res.PrivatePoolValue, fee-totalAggregated)
				res.PrivatePoolValue -= deduct
				res.FeePool += deduct
				k.VoltKeeper.SetBtcReserve(ctx, &res)
				totalAggregated += deduct
			}
			return false
		})

		if totalAggregated < fee {
			return fmt.Errorf("could not find enough value in private pool")
		}
	} else {
		reserve.PrivatePoolValue -= fee
		reserve.FeePool += fee
	}

	k.VoltKeeper.SetBtcReserve(ctx, reserve)
	return nil
}

func (k Keeper) MarkQqAccountAsUsed(ctx context.Context, QqAccount string) {
	store := k.KVStore(ctx)
	aKey := types.GetUsedQqAccountKey(QqAccount)
	store.Set(aKey, []byte{1})
}

func (k Keeper) HasUsedQqAccount(ctx context.Context, QqAccount string) bool {
	store := k.KVStore(ctx)
	aKey := types.GetUsedQqAccountKey(QqAccount)
	return store.Has(aKey)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
