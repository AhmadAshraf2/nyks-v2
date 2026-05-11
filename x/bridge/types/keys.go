package types

import (
	"bytes"
	"encoding/binary"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	forkstypes "twilight-project/nyks/x/forks/types"
)

const (
	ModuleName    = "bridge"
	StoreKey      = ModuleName
	GovModuleName = "gov"
)

var ParamsKey = collections.NewPrefix("p_bridge")

var (
	BtcReserveAddressKey       = forkstypes.HashString("BtcReserveAddressKey")
	BtcReserveScriptKey        = forkstypes.HashString("BtcReserveScriptKey")
	JudgeAddressKey            = forkstypes.HashString("JudgeAddressKey")
	BtcWithdrawRequestKey      = forkstypes.HashString("BtcWithdrawRequest")
	BtcSignRefundMsgKey        = forkstypes.HashString("BtcSignRefundMsg")
	BtcSignSweepMsgKey         = forkstypes.HashString("BtcSignSweepMsg")
	BtcBroadcastTxSweepMsgKey  = forkstypes.HashString("BtcBroadcastTxSweepMsg")
	BtcBroadcastTxRefundMsgKey = forkstypes.HashString("BtcBroadcastTxRefundMsg")
	BtcProposeRefundHashMsgKey = forkstypes.HashString("BtcProposeRefundHashMsg")
	UnsignedTxSweepMsgKey      = forkstypes.HashString("UnsignedTxSweepMsg")
	UnsignedTxRefundMsgKey     = forkstypes.HashString("UnsignedTxRefundMsg")
	ProposeSweepAddressMsg     = forkstypes.HashString("ProposeSweepAddressMsg")
	ProposeSweepAddressLockKey = forkstypes.HashString("ProposeSweepAddressLockKey")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetBtcRegisterReserveAddressKey(judgeAddress sdk.AccAddress, reserveAddress BtcAddress) []byte {
	return forkstypes.AppendBytes(BtcReserveAddressKey, judgeAddress.Bytes(), []byte(reserveAddress.BtcAddress))
}

func GetBtcRegisterReserveScriptKey(judgeAddress sdk.AccAddress, reserveAddress BtcAddress) []byte {
	return forkstypes.AppendBytes(BtcReserveScriptKey, judgeAddress.Bytes(), []byte(reserveAddress.BtcAddress))
}

func GetBootstrapFragmentAddressKey(validatorAddress sdk.ValAddress) []byte {
	return forkstypes.AppendBytes(JudgeAddressKey, validatorAddress.Bytes())
}

func GetBtcProposeRefundHashMsgKey(judgeAddress sdk.AccAddress, refundHash string) []byte {
	return forkstypes.AppendBytes(BtcProposeRefundHashMsgKey, judgeAddress.Bytes(), []byte(refundHash))
}

func GetProposeSweepAddressMsgKey(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(ProposeSweepAddressMsg, reserveId, roundId)
}

func GetUnsignedTxSweepMsgKey(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(UnsignedTxSweepMsgKey, reserveId, roundId)
}

func GetUnsignedTxRefundMsgKey(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(UnsignedTxRefundMsgKey, reserveId, roundId)
}

func GetBtcSignRefundMsgKey(reserveId uint64, roundId uint64, btcOracleAddress sdk.AccAddress) []byte {
	msgKey := generateMsgKey(BtcSignRefundMsgKey, reserveId, roundId)
	return forkstypes.AppendBytes(msgKey, btcOracleAddress.Bytes())
}

func GetBtcSignSweepMsgKey(reserveId uint64, roundId uint64, btcOracleAddress sdk.AccAddress) []byte {
	msgKey := generateMsgKey(BtcSignSweepMsgKey, reserveId, roundId)
	return forkstypes.AppendBytes(msgKey, btcOracleAddress.Bytes())
}

func GetBtcSignRefundMsgPrefix(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(BtcSignRefundMsgKey, reserveId, roundId)
}

func GetBtcSignSweepMsgPrefix(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(BtcSignSweepMsgKey, reserveId, roundId)
}

func GetBtcBroadcastTxRefundMsgKey(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(BtcBroadcastTxRefundMsgKey, reserveId, roundId)
}

func GetBtcBroadcastTxSweepMsgKey(reserveId uint64, roundId uint64) []byte {
	return generateMsgKey(BtcBroadcastTxSweepMsgKey, reserveId, roundId)
}

func uint64ToBytes(value uint64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, value)
	return buf.Bytes()
}

func generateMsgKey(hashKey []byte, reserveId uint64, roundId uint64) []byte {
	return forkstypes.AppendBytes(hashKey, uint64ToBytes(reserveId), uint64ToBytes(roundId))
}
