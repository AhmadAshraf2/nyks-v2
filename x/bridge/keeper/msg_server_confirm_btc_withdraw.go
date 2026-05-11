package keeper

import (
	"context"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) ConfirmBtcWithdraw(goCtx context.Context, msg *types.MsgConfirmBtcWithdraw) (*types.MsgConfirmBtcWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddress, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	found := k.CheckJudgeValidatorInSet(ctx, judgeAddress)
	if !found {
		return nil, types.ErrJudgeValidatorNotFound.Wrap("Could not check judge validator inset")
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, fmt.Errorf("could not check Any value: %w", err)
	}

	valAddr, err := k.GetValidatorAddressForJudgeAddress(ctx, judgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not find ValAddr for delegate key: %w", err)
	}

	err = k.NyksKeeper.ClaimHandlerCommon(ctx, any, valAddr, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgConfirmBtcWithdrawResponse{}, nil
}
