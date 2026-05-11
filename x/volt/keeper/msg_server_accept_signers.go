package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/volt/types"
)

func (k msgServer) AcceptSigners(goCtx context.Context, msg *types.MsgAcceptSigners) (*types.MsgAcceptSignersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve the fragment from the store
	fragment, found := k.Keeper.GetFragment(ctx, msg.FragmentId)
	if !found {
		return nil, types.ErrFragmentNotFound.Wrapf("fragment %d not found", msg.FragmentId)
	}

	// Check if signer of the message is the judge of the fragment
	if fragment.JudgeAddress != msg.JudgeAddress {
		return nil, types.ErrJudgeMismatch.Wrapf("signer %s is not the judge of fragment %d", msg.JudgeAddress, msg.FragmentId)
	}

	// Add each signer application to the fragment
	for _, applicationId := range msg.SignerApplicationIds {
		// Check if the fragment already has the maximum number of signers
		if len(fragment.Signers) >= int(types.MaxSignersPerFragment) {
			return nil, types.ErrMaxSignersReached.Wrapf("A fragment can not have more than %d signers", types.MaxSignersPerFragment)
		}

		application, found := k.Keeper.GetSignerApplication(ctx, msg.FragmentId, applicationId)
		if !found {
			return nil, types.ErrApplicationNotFound.Wrapf("signer application %d not found", applicationId)
		}

		// check this signer already exists in any of the fragments
		exists := k.GetExistingSignerInFragments(ctx, application.SignerAddress)
		if exists {
			return nil, types.ErrSignerAlreadyExists.Wrapf("signer %s already exists in one of the fragments", application.SignerAddress)
		}

		newSigner := &types.FragmentSigners{
			FragmentID:           msg.FragmentId,
			SignerAddress:        application.SignerAddress,
			SignerStatus:         true,
			SignerBtcPublicKey:   application.BtcPubKey,
			SignerApplicationFee: application.ApplicationFee,
			SignerFeeBips:        application.FeeBips,
		}

		fragment.Signers = append(fragment.Signers, newSigner)

		// Return the signer application fee to the signer address from module account
		err := k.ReturnSignerApplicationFee(ctx, application.SignerAddress, application.ApplicationFee)
		if err != nil {
			return nil, types.ErrCouldNotReturnSignerApplicationFee.Wrap(fmt.Sprint(application.SignerAddress))
		}
	}

	// Check if the fragment now has the maximum number of signers and update the status
	if len(fragment.Signers) == int(types.MaxSignersPerFragment) {
		fragment.FragmentStatus = true
	}

	// Save the updated fragment back to the store
	err := k.Keeper.SetFragment(ctx, fragment)
	if err != nil {
		return nil, types.ErrCouldNotSetFragment.Wrap(fmt.Sprint(msg.FragmentId))
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventAcceptSigners{
			Message:      "MsgAcceptSigners",
			FragmentId:   msg.FragmentId,
			JudgeAddress: msg.JudgeAddress,
		},
	)
	return &types.MsgAcceptSignersResponse{}, nil
}
