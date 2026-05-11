package types

import (
	"cosmossdk.io/collections"
	forkstypes "twilight-project/nyks/x/forks/types"
)

const (
	ModuleName    = "zkos"
	StoreKey      = ModuleName
	GovModuleName = "gov"
)

var ParamsKey = collections.NewPrefix("p_zkos")

var (
	KeyTransferTx          = forkstypes.HashString("KeyTransferTx")
	KeyMintOrBurnTradingBtc = forkstypes.HashString("KeyMintOrBurnTradingBtc")
	KeyUsedQqAccount       = forkstypes.HashString("KeyUsedQqAccount")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetTransferTxKey(txId string) []byte {
	return forkstypes.AppendBytes(KeyTransferTx, []byte(txId))
}

func GetMintOrBurnTradingBtcKey(twilightAddress string, QqAccount string) []byte {
	return forkstypes.AppendBytes(KeyMintOrBurnTradingBtc, []byte(twilightAddress), []byte(QqAccount))
}

func GetUsedQqAccountKey(QqAccount string) []byte {
	return forkstypes.AppendBytes(KeyUsedQqAccount, []byte(QqAccount))
}
