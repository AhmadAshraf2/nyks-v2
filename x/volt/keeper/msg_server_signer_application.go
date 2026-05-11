package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/volt/types"
)

func (k msgServer) SignerApplication(goCtx context.Context, msg *types.MsgSignerApplication) (*types.MsgSignerApplicationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signerAddress, err := sdk.AccAddressFromBech32(msg.SignerAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrap("invalid signer address")
	}

	feeAmount := sdk.NewCoin("nyks", math.NewIntFromUint64(msg.ApplicationFee))

	// Get fragment with the id
	_, found := k.GetFragment(ctx, msg.FragmentId)
	if !found {
		return nil, types.ErrFragmentNotFound
	}

	// check if similar signer application already exists
	_, found = k.GetSignerApplicationBySignerAndFragmentId(ctx, msg.FragmentId, msg.SignerAddress, msg.BtcPubKey)
	if found {
		return nil, types.ErrSignerApplicationExists
	}

	// Check if signer has enough balance to pay the application fee
	balance := k.BankKeeper.GetBalance(ctx, signerAddress, "nyks")
	if balance.Amount.LT(math.NewIntFromUint64(msg.ApplicationFee)) {
		return nil, types.ErrInsufficientFunds
	}

	voltModuleAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	if voltModuleAcc == nil {
		return nil, types.ErrVoltModuleAccountNotFound
	}

	// Deduct the application fee from the signer's account
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, signerAddress, types.ModuleName, sdk.NewCoins(feeAmount))
	if err != nil {
		return nil, err
	}

	// Generate a new application ID
	lastApplicationID := k.Keeper.GetLastRegisteredApplicationId(ctx)
	newApplicationID := lastApplicationID + 1

	// Save the application data in the store
	signerApp := types.SignerApplicationData{
		ApplicationId:  newApplicationID,
		FragmentId:     msg.FragmentId,
		ApplicationFee: msg.ApplicationFee,
		FeeBips:        msg.FeeBips,
		BtcPubKey:      msg.BtcPubKey,
		SignerAddress:  msg.SignerAddress,
	}

	k.Keeper.SetSignerApplication(ctx, &signerApp)
	k.Keeper.setLastRegisteredApplicationId(ctx, newApplicationID)

	ctx.EventManager().EmitTypedEvent(
		&types.EventSignerApplication{
			Message:       "MsgSignerApplication",
			ApplicationId: newApplicationID,
		},
	)
	return &types.MsgSignerApplicationResponse{ApplicationId: newApplicationID}, nil
}
