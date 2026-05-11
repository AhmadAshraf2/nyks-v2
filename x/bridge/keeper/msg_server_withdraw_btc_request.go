package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) WithdrawBtcRequest(goCtx context.Context, msg *types.MsgWithdrawBtcRequest) (*types.MsgWithdrawBtcRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	twilightAddress, e1 := sdk.AccAddressFromBech32(msg.TwilightAddress)
	if e1 != nil {
		return nil, types.ErrInvalid.Wrap(e1.Error())
	}

	_, e2 := types.NewBtcAddress(msg.WithdrawAddress)
	if e2 != nil {
		return nil, types.ErrInvalid.Wrap(e2.Error())
	}

	found := k.VoltKeeper.CheckBtcReserveExists(ctx, msg.ReserveId)
	if !found {
		return nil, fmt.Errorf("btc reserve not found: %d", msg.ReserveId)
	}

	err := k.VoltKeeper.CheckClearingAccountBalance(ctx, twilightAddress, msg.ReserveId, msg.WithdrawAmount)
	if err != nil {
		return nil, types.ErrInsufficientBalance.Wrap("Insufficient balance in clearing account")
	}

	userBalance := k.BankKeeper.GetBalance(ctx, twilightAddress, "sats")
	if userBalance.Amount.LT(math.NewIntFromUint64(msg.WithdrawAmount)) {
		return nil, types.ErrInsufficientBalance.Wrap("Insufficient balance in bank")
	}

	withdrawIdentifier, err := k.VoltKeeper.SetBtcWithdrawRequest(ctx, twilightAddress, msg.ReserveId, msg.WithdrawAddress, msg.WithdrawAmount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventWithdrawBtcRequest{
			Message:         "MsgWithdrawBtcRequest",
			TwilightAddress: msg.TwilightAddress,
			ReserveId:       msg.ReserveId,
			WithdrawAddress: msg.WithdrawAddress,
			WithdrawAmount:  msg.WithdrawAmount,
		},
	)

	return &types.MsgWithdrawBtcRequestResponse{WithdrawIdentifer: *withdrawIdentifier}, nil
}
