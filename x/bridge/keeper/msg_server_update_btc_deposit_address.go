package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) UpdateBtcDepositAddress(goCtx context.Context, msg *types.MsgUpdateBtcDepositAddress) (*types.MsgUpdateBtcDepositAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check that the sender is a bonded validator
	_, err := k.NyksKeeper.CheckOrchestratorValidatorInSet(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap("caller is not a bonded validator")
	}

	btcAddr, e1 := types.NewBtcAddress(msg.BtcDepositAddress)
	twilightAddress, e2 := sdk.AccAddressFromBech32(msg.TwilightAddress)
	if e1 != nil {
		return nil, types.ErrInvalid.Wrap(e1.Error())
	}
	if e2 != nil {
		return nil, types.ErrInvalid.Wrap(e2.Error())
	}

	errSetting := k.VoltKeeper.SetBtcDeposit(ctx, btcAddr.GetBtcAddress(), twilightAddress, msg.TwilightStakingAmount, msg.BtcSatoshiTestAmount)
	if errSetting != nil {
		return nil, errSetting
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterBtcDepositAddress{
			Message:        "MsgUpdateBtcDepositAddress",
			DepositAddress: btcAddr.GetBtcAddress(),
		},
	)

	return &types.MsgUpdateBtcDepositAddressResponse{}, nil
}
