package keeper

import (
	"context"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/bridge/types"
)

func (k msgServer) SweepProposal(goCtx context.Context, msg *types.MsgSweepProposal) (*types.MsgSweepProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse judge address: %w", err)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, fmt.Errorf("could not check Any value: %w", err)
	}

	valAddr, err := k.NyksKeeper.CheckOrchestratorValidatorInSet(ctx, msg.OracleAddress)
	if err != nil {
		return nil, fmt.Errorf("could not check orchestrator validator inset: %w", err)
	}

	_, resErr := k.VoltKeeper.GetBtcReserve(ctx, msg.ReserveId)
	if resErr != nil {
		return nil, fmt.Errorf("btc reserve not found: %d", msg.ReserveId)
	}

	err = k.NyksKeeper.ClaimHandlerCommon(ctx, any, valAddr, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSweepProposalResponse{}, nil
}
