package types

import (
	"bytes"
	"encoding/binary"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	forkstypes "twilight-project/nyks/x/forks/types"
)

const (
	ModuleName = "volt"
	StoreKey   = ModuleName
	GovModuleName = "gov"

	BtcReserveMaxLimit          = uint64(10)
	FragmentMaxLimit            = uint64(10)
	MaxOutgoingBtcOutputs       = 2
	MaxSignersPerFragment       = 6
	MinSignersPerFragment       = 3
	FragmentSignersMinThreshold = 2
	MaxReservesPerFragment      = 1
)

var ParamsKey = collections.NewPrefix("p_volt")

var (
	TwilightClearingAccountKey          = forkstypes.HashString("TwilightClearingAccountKey")
	BtcReserveKey                       = forkstypes.HashString("BtcKeyReserve")
	LastRegisteredReserveKey            = forkstypes.HashString("LastRegisteredReserveKey")
	WithdrawPoolKey                     = forkstypes.HashString("WithdrawPoolKey")
	LastUnlockedReserveKey              = forkstypes.HashString("LastUnlockedReserveKey")
	BtcDepositKey                       = forkstypes.HashString("BtcDepositKey")
	BtcWithdrawRequestKey               = forkstypes.HashString("BtcWithdrawRequestKey")
	ReserveWithdrawPoolKey              = forkstypes.HashString("ReserveWithdrawPoolKey")
	NewSweepProposalReceivedKey         = forkstypes.HashString("NewSweepProposalReceivedKey")
	ReserveWithdrawSnapshotKey          = forkstypes.HashString("LastWithdrawSnapshotKey")
	RefundTxSnapshotKey                 = forkstypes.HashString("LastRefundTxSnapshotKey")
	SignerApplicationFeeKey             = KeyPrefix("SignerApplicationFeeKey")
	FragmentKey                         = forkstypes.HashString("FragmentKey")
	LastRegisteredFragmentKey           = forkstypes.HashString("LastRegisteredFragmentKey")
	LastRegisteredFragmentApplicationKey = forkstypes.HashString("LastRegisteredFragmentApplicationKey")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetClearingAccountKey(twilightAddress sdk.AccAddress) []byte {
	return forkstypes.AppendBytes(TwilightClearingAccountKey, twilightAddress.Bytes())
}

func GetReserveKey(reserveId uint64) []byte {
	reserveBufBytes := new(bytes.Buffer)
	binary.Write(reserveBufBytes, binary.LittleEndian, reserveId)
	return forkstypes.AppendBytes(BtcReserveKey, reserveBufBytes.Bytes())
}

func GetWithdrawPoolKey(reserveId uint64) []byte {
	reserveBufBytes := new(bytes.Buffer)
	binary.Write(reserveBufBytes, binary.LittleEndian, reserveId)
	return forkstypes.AppendBytes(WithdrawPoolKey, reserveBufBytes.Bytes())
}

func GetBtcDepositKey(twilightAddress sdk.AccAddress) []byte {
	return forkstypes.AppendBytes(BtcDepositKey, twilightAddress.Bytes())
}

func GetBtcWithdrawRequestKeyInternal(twilightAddress sdk.AccAddress, reserveId uint64, withdrawAddress string, withdrawAmount uint64) []byte {
	withdrawAmountBuf := new(bytes.Buffer)
	binary.Write(withdrawAmountBuf, binary.LittleEndian, withdrawAmount)
	reserveBufBytes := new(bytes.Buffer)
	binary.Write(reserveBufBytes, binary.LittleEndian, reserveId)
	return forkstypes.AppendBytes(BtcWithdrawRequestKey, twilightAddress.Bytes(), reserveBufBytes.Bytes(), []byte(withdrawAddress), withdrawAmountBuf.Bytes())
}

func GetReserveWithdrawPoolKey(reserveId uint64) []byte {
	reserveBufBytes := new(bytes.Buffer)
	binary.Write(reserveBufBytes, binary.LittleEndian, reserveId)
	return forkstypes.AppendBytes(ReserveWithdrawPoolKey, reserveBufBytes.Bytes())
}

func GetNewSweepProposalReceivedKey() []byte {
	return NewSweepProposalReceivedKey
}

func GetReserveWithdrawSnapshotKey(reserveId uint64, roundId uint64) []byte {
	reserveBufBytes := new(bytes.Buffer)
	binary.Write(reserveBufBytes, binary.LittleEndian, reserveId)
	roundBufBytes := new(bytes.Buffer)
	binary.Write(roundBufBytes, binary.LittleEndian, roundId)
	return forkstypes.AppendBytes(ReserveWithdrawSnapshotKey, reserveBufBytes.Bytes(), roundBufBytes.Bytes())
}

func GetRefundTxSnapshotKey(reserveId uint64, roundId uint64) []byte {
	reserveBufBytes := new(bytes.Buffer)
	binary.Write(reserveBufBytes, binary.LittleEndian, reserveId)
	roundBufBytes := new(bytes.Buffer)
	binary.Write(roundBufBytes, binary.LittleEndian, roundId)
	return forkstypes.AppendBytes(RefundTxSnapshotKey, reserveBufBytes.Bytes(), roundBufBytes.Bytes())
}

func GetSignerApplicationFeeKey(fragmentId uint64, applicationId uint64) []byte {
	fragmentIdBuf := new(bytes.Buffer)
	binary.Write(fragmentIdBuf, binary.LittleEndian, fragmentId)
	appIdBuf := new(bytes.Buffer)
	binary.Write(appIdBuf, binary.LittleEndian, applicationId)
	return forkstypes.AppendBytes(SignerApplicationFeeKey, fragmentIdBuf.Bytes(), appIdBuf.Bytes())
}

func GetFragmentKey(fragmentId uint64) []byte {
	fragmentBufBytes := new(bytes.Buffer)
	binary.Write(fragmentBufBytes, binary.LittleEndian, fragmentId)
	return forkstypes.AppendBytes(FragmentKey, fragmentBufBytes.Bytes())
}

func GetSignerApplicationFeePrefix(fragmentId uint64) []byte {
	fragmentIdBuf := new(bytes.Buffer)
	binary.Write(fragmentIdBuf, binary.LittleEndian, fragmentId)
	return forkstypes.AppendBytes(SignerApplicationFeeKey, fragmentIdBuf.Bytes())
}
