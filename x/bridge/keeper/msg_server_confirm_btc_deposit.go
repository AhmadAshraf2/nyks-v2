package keeper

import (
	"context"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) ConfirmBtcDeposit(goCtx context.Context, msg *types.MsgConfirmBtcDeposit) (*types.MsgConfirmBtcDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := k.NyksKeeper.CheckOrchestratorValidatorInSet(ctx, msg.OracleAddress)
	if err != nil {
		return nil, fmt.Errorf("could not check orchestrator validator inset: %w", err)
	}

	// Check registered reserve address
	_, err = k.VoltKeeper.GetBtcReserveIdByAddress(ctx, msg.ReserveAddress)
	if err != nil {
		return nil, fmt.Errorf("btc reserve not found: %s", msg.ReserveAddress)
	}

	twilightDepositAddress, err := sdk.AccAddressFromBech32(msg.TwilightDepositAddress)
	if err != nil {
		return nil, err
	}

	_, found := k.VoltKeeper.GetClearingAccount(ctx, twilightDepositAddress)
	if !found {
		_, found := k.VoltKeeper.GetBtcDepositAddressByTwilightAddress(ctx, twilightDepositAddress)
		if !found {
			return nil, types.ErrClearingAccountDoesNotExist.Wrap("Clearing account for given twilight address doesn't exist")
		}
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, fmt.Errorf("could not check Any value: %w", err)
	}

	err = k.NyksKeeper.ClaimHandlerCommon(ctx, any, valAddr, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgConfirmBtcDepositResponse{}, nil
}
