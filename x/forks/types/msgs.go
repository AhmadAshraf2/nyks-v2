package types

import (
	"crypto/sha256"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BtcProposal represents a proposal of the latest btc chaintip
type BtcProposal interface {
	GetHeight() uint64
	GetProposarOrchestrator() sdk.AccAddress
	GetType() ProposalType
	ProposalHash() ([]byte, error)
}

func (msg MsgSeenBtcChainTip) GetProposarOrchestrator() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(msg.BtcOracleAddress)
	if err != nil {
		panic(err)
	}
	return val
}

func (msg *MsgSeenBtcChainTip) GetType() ProposalType {
	return PROPOSAL_TYPE_BTC_CHAINTIP
}

func (msg *MsgSeenBtcChainTip) ProposalHash() ([]byte, error) {
	path := fmt.Sprintf("%d/%s/", msg.Height, msg.Hash)
	hash := sha256.Sum256([]byte(path))
	return hash[:], nil
}
