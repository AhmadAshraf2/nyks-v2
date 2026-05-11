package keeper

import (
	"context"

	"twilight-project/nyks/x/zkos/types"
)

func (k msgServer) MintBurnTradingBtc(goCtx context.Context, msg *types.MsgMintBurnTradingBtc) (*types.MsgMintBurnTradingBtcResponse, error) {
	check, err := k.RevealCommitment(msg.QqAccount, msg.EncryptScalar, msg.GetBtcValue())
	if !check {
		if err != nil {
			return nil, types.ErrInvalidCommitment.Wrap(err.Error())
		}
		return nil, types.ErrInvalidCommitment
	}

	err = k.SetMintOrBurnTradingBtc(goCtx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgMintBurnTradingBtcResponse{}, nil
}
