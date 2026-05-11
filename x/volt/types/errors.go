package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/volt module sentinel errors
var (
	ErrInvalidSigner                      = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrBtcReserveMaxLimitReached          = errors.Register(ModuleName, 1, "Btc max reserve limit reached")
	ErrBtcReserveNotFound                 = errors.Register(ModuleName, 2, "Btc reserve not found")
	ErrInsufficientBtcValue               = errors.Register(ModuleName, 3, "Insufficient Btc value")
	ErrClearingAccountNotFound            = errors.Register(ModuleName, 4, "Clearing account not found")
	ErrCouldNotSetReserve                 = errors.Register(ModuleName, 5, "Could not set reserve")
	ErrCouldNotSetReserveWithdrawPool     = errors.Register(ModuleName, 6, "Could not set reserve withdraw pool")
	ErrCouldNotMarshalWithdrawPool        = errors.Register(ModuleName, 7, "Could not marshal withdraw pool")
	ErrBtcDepositAddressNotFound          = errors.Register(ModuleName, 8, "Btc deposit address not found")
	ErrCouldNotSetClearingAccount         = errors.Register(ModuleName, 9, "Could not set clearing account")
	ErrCouldNotReturnUserDeposit          = errors.Register(ModuleName, 10, "Could not return user deposit")
	ErrBtcSatoshiTestAmountNotEqual       = errors.Register(ModuleName, 11, "Btc satoshi test amount not equal")
	ErrInsufficientBalanceInReserve       = errors.Register(ModuleName, 12, "Insufficient balance in reserve")
	ErrSnapshotNotFound                   = errors.Register(ModuleName, 13, "Snapshot not found")
	ErrInvalid                            = errors.Register(ModuleName, 14, "Invalid")
	ErrCouldNotSetFragment                = errors.Register(ModuleName, 15, "Could not set fragment")
	ErrFragmentMaxLimitReached            = errors.Register(ModuleName, 16, "Fragment max limit reached")
	ErrFragmentNotFound                   = errors.Register(ModuleName, 17, "Fragment not found")
	ErrReserveAlreadyExists               = errors.Register(ModuleName, 18, "This reserve id already exists in the passed fragment id")
	ErrMaxSignersReached                  = errors.Register(ModuleName, 19, "Max signers reached")
	ErrMinSignersNotMet                   = errors.Register(ModuleName, 20, "Min signers not met")
	ErrApplicationNotFound                = errors.Register(ModuleName, 21, "Application not found")
	ErrSignerNotFound                     = errors.Register(ModuleName, 22, "Signer not found")
	ErrSignerAlreadyExists                = errors.Register(ModuleName, 23, "Signer already exists")
	ErrJudgeMismatch                      = errors.Register(ModuleName, 24, "Judge mismatch")
	ErrSignerApplicationExists            = errors.Register(ModuleName, 25, "Signer application or given btc pk already exists")
	ErrInsufficientFunds                  = errors.Register(ModuleName, 26, "Signer address has insufficient funds")
	ErrVoltModuleAccountNotFound          = errors.Register(ModuleName, 27, "Volt module account not found")
	ErrCouldNotReturnSignerApplicationFee = errors.Register(ModuleName, 28, "Could not return signer application fee")
)
