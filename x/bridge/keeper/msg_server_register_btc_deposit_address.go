package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) RegisterBtcDepositAddress(goCtx context.Context, msg *types.MsgRegisterBtcDepositAddress) (*types.MsgRegisterBtcDepositAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	btcAddr, err := types.NewBtcAddress(msg.BtcDepositAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap(err.Error())
	}

	twilightAddress, err := sdk.AccAddressFromBech32(msg.TwilightAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap(err.Error())
	}

	// check if a btc address is already registered against this twilight address
	address, foundExistingBtcAddress := k.VoltKeeper.GetBtcDepositAddressByTwilightAddress(ctx, twilightAddress)
	if foundExistingBtcAddress {
		return nil, types.ErrResetBtcAddress.Wrap(address.BtcDepositAddress)
	}

	// check if a btc address is registered against *any other twilight address*
	checkBtcAddress := k.VoltKeeper.CheckBtcAddress(ctx, twilightAddress, btcAddr.GetBtcAddress(), msg.BtcSatoshiTestAmount)
	if checkBtcAddress {
		return nil, types.ErrBtcAddressAlreadyExists.Wrap(btcAddr.GetBtcAddress())
	}

	// Transfer the staking amount from the user's account to a module account if amount > 0
	if msg.TwilightStakingAmount > 0 {
		depositAmount := sdk.NewCoin("nyks", math.NewIntFromUint64(msg.TwilightStakingAmount))

		balance := k.BankKeeper.GetBalance(ctx, twilightAddress, "nyks")
		if balance.IsLT(depositAmount) {
			return nil, types.ErrInsufficientBalanceInBank.Wrapf("insufficient funds: %s < %s", balance, depositAmount)
		}

		errTakeStake := k.BankKeeper.SendCoinsFromAccountToModule(ctx, twilightAddress, types.ModuleName, sdk.NewCoins(depositAmount))
		if errTakeStake != nil {
			return nil, errTakeStake
		}
	}

	errSetting := k.VoltKeeper.SetBtcDeposit(ctx, btcAddr.GetBtcAddress(), twilightAddress, msg.TwilightStakingAmount, msg.BtcSatoshiTestAmount)
	if errSetting != nil {
		return nil, errSetting
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterBtcDepositAddress{
			Message:        "MsgRegisterBtcDepositAddress",
			DepositAddress: btcAddr.GetBtcAddress(),
		},
	)

	return &types.MsgRegisterBtcDepositAddressResponse{}, nil
}

// CheckandConfirmUserDeposit checks if a user has a deposit and confirms it
func (k Keeper) CheckandConfirmUserDeposit(ctx context.Context, twilightAddress sdk.AccAddress) error {
	_, found := k.VoltKeeper.GetBtcDepositAddressByTwilightAddress(ctx, twilightAddress)
	if !found {
		return fmt.Errorf("btc deposit address not found for twilight address: %s", twilightAddress)
	}
	return nil
}
