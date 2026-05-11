package types

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "forks"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_forks")

var (
	// KeyOrchestratorAddress indexes the validator keys for an orchestrator
	KeyOrchestratorAddress = HashString("KeyOrchestratorAddress")

	// BtcPublicKeyByValidatorKey indexes cosmos validator account addresses
	BtcPublicKeyByValidatorKey = HashString("BtcPublicKeyValidatorKey")

	// ValidatorByBtcPublicKeyKey is used to index btc validator key
	ValidatorByBtcPublicKeyKey = HashString("ValidatorByBtcPublicKeyKey")

	// OracleAttestationKey attestation used in GetAttestationKey
	OracleAttestationKey = HashString("OracleAttestationKey")

	// LastBlockHeightByValidatorKey indexes latest block height by validator
	LastBlockHeightByValidatorKey = HashString("LastBlockHeightByValidatorKey")

	// LastObservedBlockHeightKey indexes the latest block height
	LastObservedBlockHeightKey = HashString("LastObservedBlockHeightKey")

	// KeyValidator returns the following key format
	KeyValidator = HashString("KeyValidator")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetValidatorKey returns key for storing delegate addresses by validator
func GetValidatorKey(val sdk.ValAddress) []byte {
	return AppendBytes(KeyValidator, val.Bytes())
}

// GetOrchestratorAddressKey returns key for orchestrator address mapping
func GetOrchestratorAddressKey(orc sdk.AccAddress) []byte {
	return AppendBytes(KeyOrchestratorAddress, orc.Bytes())
}

// GetBtcPublicKeyByValidatorKey returns key for BTC public key by validator
func GetBtcPublicKeyByValidatorKey(validator sdk.ValAddress) []byte {
	return AppendBytes(BtcPublicKeyByValidatorKey, validator.Bytes())
}

// GetValidatorByBtcPublicKeyKey returns key for validator by BTC public key
func GetValidatorByBtcPublicKeyKey(btcPk BtcPublicKey) []byte {
	return AppendBytes(ValidatorByBtcPublicKeyKey, []byte(btcPk.GetBtcPublicKey()))
}

// GetAttestationKey returns key for attestation by block height and proposal hash
func GetAttestationKey(btcBlockHeight uint64, proposalHash []byte) []byte {
	return AppendBytes(OracleAttestationKey, UInt64Bytes(btcBlockHeight), proposalHash)
}

// GetLastBlockHeightByValidatorKey returns key for latest block height by validator
func GetLastBlockHeightByValidatorKey(validator sdk.ValAddress) []byte {
	return AppendBytes(LastBlockHeightByValidatorKey, validator.Bytes())
}
