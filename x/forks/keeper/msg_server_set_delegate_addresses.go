package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"twilight-project/nyks/x/forks/types"
)

func (k msgServer) SetDelegateAddresses(goCtx context.Context, msg *types.MsgSetDelegateAddresses) (*types.MsgSetDelegateAddressesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	btcOracleAdd, err := sdk.AccAddressFromBech32(msg.BtcOracleAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap("btc oracle address is not in valid format")
	}

	// check if this btc oracle address is already registered
	_, found := k.GetDelegateAddresses(ctx, btcOracleAdd)
	if found {
		return nil, types.ErrInvalid.Wrap("btc oracle address is already registered")
	}

	val, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap(err.Error())
	}

	_, err = k.StakingKeeper.GetValidator(ctx, val)
	if err != nil {
		return nil, stakingtypes.ErrNoValidatorFound.Wrap(val.String())
	}

	// set delegate addresses
	err = k.Keeper.SetDelegateAddresses(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventSetDelegateAddresses{
			Message: "MsgSetDelegateAddresses",
			Address: msg.ValidatorAddress,
		},
	)

	return &types.MsgSetDelegateAddressesResponse{}, nil
}
