package types

import (
	"crypto/sha256"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	forkstypes "twilight-project/nyks/x/forks/types"
)

// Ensure bridge message types implement BtcProposal interface
var _ forkstypes.BtcProposal = &MsgConfirmBtcDeposit{}
var _ forkstypes.BtcProposal = &MsgSweepProposal{}
var _ forkstypes.BtcProposal = &MsgConfirmBtcWithdraw{}

// MsgConfirmBtcDeposit BtcProposal interface

func (msg MsgConfirmBtcDeposit) GetProposarOrchestrator() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		panic(err)
	}
	return val
}

func (msg *MsgConfirmBtcDeposit) GetType() forkstypes.ProposalType {
	return forkstypes.PROPOSAL_TYPE_BTC_DEPOSIT
}

func (msg *MsgConfirmBtcDeposit) ProposalHash() ([]byte, error) {
	path := fmt.Sprintf("%d/%s/", msg.Height, msg.Hash)
	hash := sha256.Sum256([]byte(path))
	return hash[:], nil
}

// MsgSweepProposal BtcProposal interface

func (msg MsgSweepProposal) GetProposarOrchestrator() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		panic(err)
	}
	return val
}

func (msg *MsgSweepProposal) GetType() forkstypes.ProposalType {
	return forkstypes.PROPOSAL_TYPE_SWEEP_PROPOSAL
}

func (msg *MsgSweepProposal) ProposalHash() ([]byte, error) {
	path := fmt.Sprintf("%s", msg.BtcTxHash)
	hash := sha256.Sum256([]byte(path))
	return hash[:], nil
}

// GetHeight for SweepProposal always returns 0 since sweep txs haven't been broadcasted yet
func (msg *MsgSweepProposal) GetHeight() uint64 {
	return 0
}

// MsgConfirmBtcWithdraw BtcProposal interface

func (msg MsgConfirmBtcWithdraw) GetProposarOrchestrator() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(msg.JudgeAddress)
	if err != nil {
		panic(err)
	}
	return val
}

func (msg *MsgConfirmBtcWithdraw) GetType() forkstypes.ProposalType {
	return forkstypes.PROPOSAL_TYPE_CONFIRM_WITHDRAW
}

func (msg *MsgConfirmBtcWithdraw) ProposalHash() ([]byte, error) {
	path := fmt.Sprintf("%d/%s/", msg.Height, msg.Hash)
	hash := sha256.Sum256([]byte(path))
	return hash[:], nil
}
