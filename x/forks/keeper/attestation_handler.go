package keeper

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bridgetypes "twilight-project/nyks/x/bridge/types"
	"twilight-project/nyks/x/forks/types"
)

// AttestationHandler processes observed Attestations
type AttestationHandler struct {
	Keeper *Keeper
}

func (a AttestationHandler) ValidateMembers() {
	if a.Keeper == nil {
		panic("Nil keeper!")
	}
}

func (a AttestationHandler) Handle(ctx context.Context, att types.Attestation, proposal types.BtcProposal) error {
	switch p := proposal.(type) {
	case *types.MsgSeenBtcChainTip:
		return nil // no post-processing needed
	case *bridgetypes.MsgConfirmBtcDeposit:
		return a.handleConfirmBtcDeposit(ctx, *p)
	case *bridgetypes.MsgSweepProposal:
		return a.handleSweepProposal(ctx, *p)
	case *bridgetypes.MsgConfirmBtcWithdraw:
		return nil // no post-processing needed currently
	default:
		panic(fmt.Sprintf("Invalid event type for attestations %d", proposal.GetType()))
	}
}

func (a AttestationHandler) handleConfirmBtcDeposit(goCtx context.Context, proposal bridgetypes.MsgConfirmBtcDeposit) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mintAmount := math.NewIntFromUint64(proposal.DepositAmount)
	denom := "sats"
	coin := sdk.NewCoin(denom, mintAmount)

	moduleAddr := a.Keeper.accountKeeper.GetModuleAddress(types.ModuleName)
	preMintBalance := a.Keeper.bankKeeper.GetBalance(ctx, moduleAddr, coin.Denom)

	prevSupply := a.Keeper.bankKeeper.GetSupply(ctx, coin.Denom)
	newSupply := new(big.Int).Add(prevSupply.Amount.BigInt(), math.NewIntFromUint64(proposal.DepositAmount).BigInt())
	if newSupply.BitLen() > 256 {
		return types.ErrAttestationOverflow.Wrap("invalid supply after deposit attestation")
	}

	coins := sdk.Coins{coin}
	if err := a.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return fmt.Errorf("unable to mint cosmos originated coins %v: %w", coins, err)
	}

	postMintBalance := a.Keeper.bankKeeper.GetBalance(ctx, moduleAddr, coin.Denom)
	if !postMintBalance.Sub(preMintBalance).Amount.Equal(math.NewIntFromUint64(proposal.DepositAmount)) {
		panic(fmt.Sprintf("Somehow minted incorrect amount! Previous balance %v Post-mint balance %v proposal amount %v",
			preMintBalance.String(), postMintBalance.String(), proposal.DepositAmount))
	}

	receiver, err := sdk.AccAddressFromBech32(proposal.TwilightDepositAddress)
	if err != nil {
		return fmt.Errorf("invalid twilight deposit address %s: %w", proposal.TwilightDepositAddress, err)
	}

	err = a.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins)
	if err != nil {
		a.Keeper.Logger(ctx).Error("Could not send coins to the account", "cause", err.Error())
	}

	err = a.Keeper.VoltKeeper.UpdateBtcReserveAfterMint(ctx, proposal.DepositAmount, receiver, proposal.ReserveAddress)
	if err != nil {
		return fmt.Errorf("could not update the reserve %s: %w", proposal.ReserveAddress, err)
	}

	types.MintedSatsCounter.Add(float64(mintAmount.Int64()))

	return err
}

func (a AttestationHandler) handleSweepProposal(goCtx context.Context, proposal bridgetypes.MsgSweepProposal) error {
	err := a.Keeper.VoltKeeper.UpdateBtcReserveAfterSweepProposal(goCtx, proposal.ReserveId, proposal.NewReserveAddress, proposal.JudgeAddress, proposal.BtcBlockNumber, proposal.BtcRelayCapacityValue, proposal.BtcTxHash, proposal.UnlockHeight, proposal.RoundId)
	if err != nil {
		return fmt.Errorf("could not update the reserve after sweep attestation %s: %w", proposal.NewReserveAddress, err)
	}

	err = a.Keeper.VoltKeeper.ConfirmWithdrawRequestsAfterSweepConfirmation(goCtx, proposal.ReserveId, proposal.RoundId)
	if err != nil {
		return fmt.Errorf("could not confirm withdraw requests after sweep attestation: %w", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	a.Keeper.VoltKeeper.PruneReserveWithdrawSnapshot(ctx, proposal.ReserveId, proposal.RoundId)
	a.Keeper.VoltKeeper.PruneRefundTxSnapshot(ctx, proposal.ReserveId, proposal.RoundId)

	return nil
}
