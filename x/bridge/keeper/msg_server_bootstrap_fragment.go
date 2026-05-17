package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) BootstrapFragment(goCtx context.Context, msg *types.MsgBootstrapFragment) (*types.MsgBootstrapFragmentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	judgeAddr, e := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if e != nil {
		return nil, types.ErrInvalid.Wrap(e.Error())
	}

	// If validatorAddress is empty, derive it from judgeAddress (same as old chain --from behavior)
	validatorAddress := msg.ValidatorAddress
	if validatorAddress == "" {
		validatorAddress = sdk.ValAddress(judgeAddr).String()
	}

	valAddr, err := sdk.ValAddressFromBech32(validatorAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap("invalid validator address")
	}

	_, err = k.StakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, types.ErrInvalid.Wrap("validator not found")
	}

	address, _ := k.GetJudgeAddressForValidatorAddress(ctx, valAddr)
	if address != nil {
		return nil, types.ErrInvalid.Wrap("validator already has judge address")
	}

	errSetting := k.SetJudgeAddressForValidatorAddress(ctx, judgeAddr, msg.NumOfSigners, msg.Threshold, msg.SignerApplicationFee, msg.ArbitraryData, valAddr)
	if errSetting != nil {
		return nil, errSetting
	}

	fragmentId, errSettingRes := k.VoltKeeper.RegisterNewFragment(ctx, judgeAddr, msg.Threshold, msg.SignerApplicationFee, msg.NumOfSigners, msg.FragmentFeeBips, msg.ArbitraryData)
	if errSettingRes != nil {
		return nil, errSettingRes
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventBootstrapFragmentAddress{
			Message:          "MsgBootstrapFragment",
			JudgeAddress:     judgeAddr.String(),
			FragmentId:       fragmentId,
			ValidatorAddress: valAddr.String(),
		},
	)

	return &types.MsgBootstrapFragmentResponse{FragmentId: strconv.FormatUint(fragmentId, 10), JudgeAddress: msg.JudgeAddress}, nil
}
