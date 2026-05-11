package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) RegisterReserveAddress(goCtx context.Context, msg *types.MsgRegisterReserveAddress) (*types.MsgRegisterReserveAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, e1 := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if e1 != nil {
		return nil, types.ErrInvalid.Wrap(e1.Error())
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	reserveScript, e2 := types.NewBtcScript(msg.ReserveScript)
	if e2 != nil {
		return nil, types.ErrInvalid.Wrap(e2.Error())
	}

	reserveAddress, e3 := types.NewBtcAddress(msg.ReserveAddress)
	if e3 != nil {
		return nil, types.ErrInvalid.Wrap(e3.Error())
	}

	k.SetReserveAddressForJudge(ctx, judgeAddress, *reserveScript, *reserveAddress)

	reserveId, errSettingRes := k.VoltKeeper.RegisterNewBtcReserve(ctx, judgeAddress, reserveAddress.BtcAddress)
	if errSettingRes != nil {
		return nil, errSettingRes
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterReserveAddress{
			Message:       "MsgRegisterReserveAddress",
			ReserveScript: msg.ReserveScript,
		},
	)

	return &types.MsgRegisterReserveAddressResponse{
		ReserveId:      strconv.FormatUint(reserveId, 10),
		ReserveAddress: msg.ReserveAddress,
	}, nil
}
