package keeper

import (
	"context"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/forks/types"
)

func (k msgServer) SeenBtcChainTip(goCtx context.Context, msg *types.MsgSeenBtcChainTip) (*types.MsgSeenBtcChainTipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := k.CheckOrchestratorValidatorInSet(ctx, msg.BtcOracleAddress)
	if err != nil {
		return nil, fmt.Errorf("could not check orchestrator validator inset: %w", err)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, fmt.Errorf("could not check Any value: %w", err)
	}

	// Set the telemetry gauge for this oracle to the block number
	types.OracleBlockGauge.WithLabelValues(msg.BtcOracleAddress).Set(float64(msg.Height))

	err = k.ClaimHandlerCommon(ctx, any, valAddr, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSeenBtcChainTipResponse{}, nil
}

func (k Keeper) ClaimHandlerCommon(ctx sdk.Context, msgAny *codectypes.Any, valAddr sdk.ValAddress, msg types.BtcProposal) error {
	_, err := k.Attest(ctx, msg, valAddr, msgAny)
	if err != nil {
		return fmt.Errorf("err while creating an attestation: %w", err)
	}
	hash, err := msg.ProposalHash()
	if err != nil {
		return fmt.Errorf("unable to compute proposal hash: %w", err)
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventProposal{
			Message:       string(msg.GetType()),
			ProposalHash:  string(hash),
			AttestationId: string(types.GetAttestationKey(msg.GetHeight(), hash)),
		},
	)

	return nil
}
